// Loopback owner-gate checks in real Chrome. Credentials stay in memory and are never printed.
// Trace, HAR, and screenshots are not recorded. The external Playwright core path is SAR_BETA_PW.
import { spawnSync } from "node:child_process";
import { createHash } from "node:crypto";
import { chmodSync, readFileSync, statSync, writeFileSync } from "node:fs";
import { createRequire } from "node:module";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const require = createRequire(import.meta.url);
const DEPLOY = "/home/shin/deploy/knowslink";
const FIXTURE = `${DEPLOY}/build/qa-fixture.json`;
const BASE = "http://127.0.0.1:8080";
const EXPECT_SHA = "28bd1bb2bdd4b90d0c6f5d2a8d4ba3e00f630ff2";
const PRODUCT = "78b1d92c8aa626245d3349ffaf7367d28f1dd3ef";
const outDir = dirname(fileURLToPath(import.meta.url));
const uuidRe = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/;
const checks = [];
const meta = { user_cache_observed: false, reproduction: "isolated-browser-contexts" };

function check(id, ok, extra = {}) {
  checks.push({ id, ok, ...extra });
  console.log(`${ok ? "PASS" : "FAIL"} ${id}`);
  return ok;
}
function git(args) {
  const result = spawnSync("git", ["-C", DEPLOY, ...args], { encoding: "utf8" });
  if (result.status !== 0) throw new Error(`git ${args[0]} failed`);
  return result.stdout.trim();
}
function psql(sql) {
  const result = spawnSync(
    "docker",
    ["exec", "-i", "knowslink-postgres-1", "psql", "-U", "knowslink", "-d", "knowslink", "-At", "-v", "ON_ERROR_STOP=1"],
    { input: sql, encoding: "utf8" },
  );
  if (result.status !== 0) throw new Error(redact(result.stderr || "psql failed"));
  return result.stdout.trim();
}
function redact(text) {
  return String(text || "").replace(/[A-Za-z0-9_-]{24,}/g, "[redacted]").slice(0, 400);
}
function sha256(bytes) {
  return createHash("sha256").update(bytes).digest("hex");
}
function credsOf(side) {
  return { username: side.owner.owner, password: side.owner.credential };
}
function sqlId(value, pattern) {
  if (!pattern.test(value)) throw new Error("unexpected identifier");
  return value;
}
function expiredGateIds(ownerId) {
  const owner = sqlId(ownerId, /^owner_[A-Za-z0-9_-]{20,}$/);
  const out = psql(
    `SELECT g.key FROM relay_state CROSS JOIN LATERAL jsonb_each(data->'Gates') g WHERE g.value->>'Owner' = '${owner}' AND g.value->>'State' = 'expired' ORDER BY g.key;`,
  );
  return out ? out.split("\n").filter((id) => uuidRe.test(id)) : [];
}
function gateMeta(id) {
  const gate = sqlId(id, uuidRe);
  const line = psql(
    `SELECT COALESCE(g.value->>'State','') || '|' || COALESCE((SELECT m.value->>'Completion' FROM jsonb_each(data->'Messages') m WHERE m.key = g.value->>'Parent'),'') || '|' || COALESCE((SELECT m.value->'Receipt'->>'intent' FROM jsonb_each(data->'Messages') m WHERE m.key = g.value->>'Parent'),'') FROM relay_state CROSS JOIN LATERAL jsonb_each(data->'Gates') g WHERE g.key = '${gate}';`,
  );
  const [state, completion, intent] = line.split("|");
  for (const part of [state, completion, intent]) {
    if (!/^[A-Za-z0-9:._-]{0,80}$/.test(part || "")) throw new Error("unexpected gate meta");
  }
  return { state, completion, intent };
}
function parseDisplayed(value) {
  const match = /^(\d{4}-\d{2}-\d{2}) (\d{2}:\d{2}:\d{2})(?:\.(\d+))? ([+-]\d{4})/.exec(value.trim());
  if (!match) return Date.parse(value);
  const fraction = (match[3] || "0").slice(0, 3).padEnd(3, "0");
  return Date.parse(`${match[1]}T${match[2]}.${fraction}${match[4].slice(0, 3)}:${match[4].slice(3)}`);
}
function seed(label) {
  const result = spawnSync("node", ["adapters/dist/synthetic.js", BASE, "--seed"], { cwd: DEPLOY, encoding: "utf8" });
  const line = (result.stdout || "").split("\n").find((item) => item.startsWith("UI candidate:")) || "";
  const match = /^UI candidate: http:\/\/127\.0\.0\.1:8080\/owner\/gates\/([0-9a-f-]+); expires (\S+)$/.exec(line.trim());
  const ok = result.status === 0 && Boolean(match);
  check(`seed_${label}`, ok, { exit: result.status ?? 1, gate: match?.[1] || "" });
  if (!ok) {
    console.log(`SEED_FAIL ${label} ${redact(result.stderr)}`);
    throw new Error(`seed ${label} failed`);
  }
  const fixture = JSON.parse(readFileSync(FIXTURE, "utf8"));
  return { gate: match[1], exp: match[2], ownerB: credsOf(fixture.b), ownerA: credsOf(fixture.a) };
}
async function readGate(page, exp) {
  const heading = await page.locator("h2").first().innerText();
  const pre = await page.locator("pre").innerText();
  const intent = await page.locator("dd").nth(3).innerText();
  const displayed = await page.locator("dd").nth(4).innerText();
  const policy = await page.locator("dd").nth(5).innerText();
  const text = await page.locator("body").innerText();
  return {
    state: heading.replace("상태:", "").trim(),
    bodyOk: pre.includes("granularity_min") && pre.includes("2026-10-03T10:00:00Z"),
    intentOk: intent.trim() === "schedule.query",
    policyOk: policy.includes("disclosure policy absent") && policy.includes("deny even after gate approve"),
    expOk: Math.abs(parseDisplayed(displayed) - Date.parse(exp)) < 2000,
    buttons: await page.locator("button").count(),
    notice: text.includes("정보 공개나 일정 실행을 허용하지 않습니다."),
    inactive: text.includes("승인·거절 버튼 비활성:"),
    text,
  };
}
async function decide(browser, seeded, decision, state) {
  const context = await browser.newContext({ httpCredentials: seeded.ownerB });
  const page = await context.newPage();
  page.setDefaultTimeout(20000);
  const posts = [];
  page.on("response", (response) => {
    if (response.request().method() === "POST") posts.push(response.status());
  });
  try {
    const dashboard = await page.goto(`${BASE}/owner`, { waitUntil: "domcontentloaded" });
    const dashText = await page.locator("body").innerText();
    const link = await page.locator(`a[href="/owner/gates/${seeded.gate}"]`).count();
    check(`${decision}_dashboard`, dashboard.status() === 200 && dashText.includes("KnowsLink Owner 작업 화면") && link === 1, { status: dashboard.status() });
    const opened = await page.goto(`${BASE}/owner/gates/${seeded.gate}`, { waitUntil: "domcontentloaded" });
    const fresh = await readGate(page, seeded.exp);
    check(`${decision}_fresh`, opened.status() === 200 && fresh.state === "pending" && fresh.bodyOk && fresh.intentOk && fresh.policyOk && fresh.expOk && fresh.buttons === 2 && fresh.notice, { status: opened.status() });
    const name = decision === "approve" ? "Approve 승인" : "Deny 거절";
    await page.getByRole("button", { name }).click();
    await page.waitForLoadState("domcontentloaded");
    const decided = await readGate(page, seeded.exp);
    check(`${decision}_decided`, decided.state === state && decided.buttons === 0 && decided.text.includes(`승인·거절 버튼 비활성: ${state}`) && decided.notice && decided.policyOk && decided.bodyOk, { posts: posts.join(",") });
    await page.reload({ waitUntil: "domcontentloaded" });
    const reloaded = await readGate(page, seeded.exp);
    check(`${decision}_reload`, reloaded.state === state && reloaded.buttons === 0 && reloaded.text.includes(`승인·거절 버튼 비활성: ${state}`));
    const effect = gateMeta(seeded.gate);
    check(`${decision}_no_schedule_effect`, effect.state === state && effect.completion === "" && effect.intent === "schedule.query", { completion: effect.completion || "empty" });
  } finally {
    await context.close();
  }
}
async function expectDenied(browser, creds, gate, id) {
  const context = await browser.newContext({ httpCredentials: creds });
  const page = await context.newPage();
  page.setDefaultTimeout(20000);
  try {
    const response = await page.goto(`${BASE}/owner/gates/${gate}`, { waitUntil: "domcontentloaded" });
    const text = await page.locator("body").innerText();
    const buttons = await page.locator("button").count();
    check(id, response.status() === 403 && text.includes("sender_not_allowed") && !text.includes("granularity_min") && buttons === 0, { status: response.status() });
  } finally {
    await context.close();
  }
}
async function recover(browser, seeded) {
  const context = await browser.newContext({ httpCredentials: seeded.ownerB });
  const page = await context.newPage();
  page.setDefaultTimeout(20000);
  try {
    const response = await page.goto(`${BASE}/owner/gates/${seeded.gate}`, { waitUntil: "domcontentloaded" });
    const text = await page.locator("body").innerText();
    check("fresh_context_200", response.status() === 200 && text.includes("상태: pending") && text.includes("granularity_min") && !text.includes("sender_not_allowed"), { status: response.status() });
  } finally {
    await context.close();
  }
}
async function readExpired(browser, creds, gate, id) {
  const context = await browser.newContext({ httpCredentials: creds });
  const page = await context.newPage();
  page.setDefaultTimeout(20000);
  const posts = [];
  page.on("response", (response) => {
    if (response.request().method() === "POST") posts.push(response.status());
  });
  try {
    const response = await page.goto(`${BASE}/owner/gates/${gate}`, { waitUntil: "domcontentloaded" });
    const text = await page.locator("body").innerText();
    const buttons = await page.locator("button").count();
    const ok = response.status() === 200 && text.includes("상태: expired") && text.includes("승인·거절 버튼 비활성: expired") && buttons === 0 && posts.length === 0 && !text.includes("Approve 승인");
    check(id, ok, { status: response.status(), buttons, posts: posts.length });
    return ok;
  } finally {
    await context.close();
  }
}
async function quietExpired(browser, creds, gate) {
  const context = await browser.newContext({ httpCredentials: creds });
  const page = await context.newPage();
  page.setDefaultTimeout(20000);
  try {
    const response = await page.goto(`${BASE}/owner/gates/${gate}`, { waitUntil: "domcontentloaded" });
    const text = await page.locator("body").innerText();
    const buttons = await page.locator("button").count();
    return response.status() === 200 && text.includes("상태: expired") && text.includes("승인·거절 버튼 비활성: expired") && buttons === 0 && !text.includes("Approve 승인");
  } finally {
    await context.close();
  }
}
async function waitExpiry(browser) {
  const seeded = seed("expiry");
  const waitMs = Math.max(0, Date.parse(seeded.exp) + 2000 - Date.now());
  console.log(`WAIT_EXPIRY_MS ${waitMs}`);
  await new Promise((resolve) => setTimeout(resolve, waitMs));
  let extra = 0;
  let seen = await quietExpired(browser, seeded.ownerB, seeded.gate);
  for (let attempt = 0; attempt < 30 && !seen; attempt += 1) {
    extra += 2000;
    await new Promise((resolve) => setTimeout(resolve, 2000));
    seen = await quietExpired(browser, seeded.ownerB, seeded.gate);
  }
  const recorded = await readExpired(browser, seeded.ownerB, seeded.gate, "expiry_waited");
  meta.expiry_mode = "waited";
  meta.wait_ms = waitMs;
  meta.extra_ms = extra;
  meta.expiry_gate = seeded.gate;
  return seen && recorded;
}

let browser;
let backup;
let mode = 0;
let exitCode = 2;
try {
  const version = spawnSync("/usr/bin/google-chrome", ["--version"], { encoding: "utf8" });
  meta.chrome = (version.stdout || "").trim();
  const head = git(["rev-parse", "HEAD"]);
  check("deploy_sha", head === EXPECT_SHA);
  check("deploy_clean_before", git(["status", "--porcelain"]) === "");
  const productDiff = git(["diff", "--name-only", PRODUCT, "HEAD", "--", "cmd", "internal", "adapters", "db", "Dockerfile", "compose.yaml", "Makefile", "scripts"]);
  check("product_tree", productDiff === "");
  const health = await fetch(`${BASE}/healthz`);
  check("healthz_precondition", health.status === 200, { status: health.status });
  backup = readFileSync(FIXTURE);
  mode = statSync(FIXTURE).mode & 0o777;
  meta.fixture_sha256 = sha256(backup);
  meta.fixture_mode = mode.toString(8);
  meta.db_before = psql("SELECT (SELECT count(*) FROM relay_state)::text || ' ' || (SELECT COALESCE(sum(octet_length(data::text)),0) FROM relay_state)::text;");
  const playwright = require(process.env.SAR_BETA_PW || "/tmp/sar-beta-002-browser/node_modules/playwright-core");
  meta.playwright_core = playwright.chromium ? "loaded" : "missing";
  browser = await playwright.chromium.launch({
    executablePath: "/usr/bin/google-chrome",
    headless: true,
    args: ["--no-sandbox", "--disable-dev-shm-usage", "--disable-gpu"],
  });
  check("browser_launched", true, { headless: true });
  const original = JSON.parse(backup.toString("utf8"));
  const existing = expiredGateIds(original.b.owner.owner);
  meta.existing_expired_count = existing.length;
  let expiredOk = false;
  if (existing.length > 0) {
    meta.expiry_mode = "existing";
    meta.expiry_gate = existing[0];
    expiredOk = await readExpired(browser, credsOf(original.b), existing[0], "expiry_existing");
  }
  const approved = seed("approve");
  meta.approve_gate = approved.gate;
  await decide(browser, approved, "approve", "approved");
  const denied = seed("deny");
  meta.deny_gate = denied.gate;
  await expectDenied(browser, approved.ownerB, denied.gate, "previous_owner_403");
  await expectDenied(browser, denied.ownerA, denied.gate, "other_owner_403");
  await recover(browser, denied);
  await decide(browser, denied, "deny", "denied");
  if (!expiredOk) await waitExpiry(browser);
  exitCode = checks.some((item) => !item.ok) ? 1 : 0;
} catch (error) {
  console.log(`BROWSER_FAIL ${redact(error.stack || error.message)}`);
  exitCode = checks.some((item) => item.id === "browser_launched" && item.ok) ? 1 : 2;
} finally {
  if (browser) await browser.close();
  if (backup) {
    writeFileSync(FIXTURE, backup, { mode: 0o600 });
    chmodSync(FIXTURE, mode);
    const restored = sha256(readFileSync(FIXTURE)) === meta.fixture_sha256 && (statSync(FIXTURE).mode & 0o777) === mode;
    check("fixture_restored", restored);
    meta.db_after = psql("SELECT (SELECT count(*) FROM relay_state)::text || ' ' || (SELECT COALESCE(sum(octet_length(data::text)),0) FROM relay_state)::text;");
    const [rowsBefore, bytesBefore] = (meta.db_before || "").split(" ");
    const [rowsAfter, bytesAfter] = meta.db_after.split(" ");
    check("db_single_row", rowsBefore === "1" && rowsAfter === "1", { bytes_before: Number(bytesBefore), bytes_after: Number(bytesAfter) });
    check("deploy_clean_after", git(["status", "--porcelain"]) === "");
    check("deploy_sha_after", git(["rev-parse", "HEAD"]) === EXPECT_SHA);
    if (checks.some((item) => !item.ok)) exitCode = exitCode === 2 ? 2 : 1;
  }
  writeFileSync(join(outDir, "result.json"), `${JSON.stringify({ ...meta, exit: exitCode, checks }, null, 2)}\n`);
  console.log(`EXIT ${exitCode}`);
}
process.exit(exitCode);
