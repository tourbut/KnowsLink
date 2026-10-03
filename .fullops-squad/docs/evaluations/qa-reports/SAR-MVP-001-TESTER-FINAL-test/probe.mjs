// Narrow HTTP probe for the 78b1d92 legacy-claim boundary.
// Statuses and error codes only. Credentials, CSRF, and claim tokens stay out of the log.
import { createHash, createPrivateKey, randomBytes, sign } from "node:crypto";
import { readFile, writeFile } from "node:fs/promises";
import { execFileSync } from "node:child_process";
import { signingBytes, uuid7 } from "../../../../../adapters/dist/index.js";

const base = process.env.RELAY_URL;
const human = "0199a3f2-4c10-7a11-8b22-334455667788";
const gated = "0199a3f2-4c10-7a11-8b22-334455667799";
const unrouted = "0199a3f2-4c10-7a11-8b22-3344556677cc";
const legacyGate = "0199a3f2-4c10-7a11-8b22-3344556677aa";
const checks = [];
if (!base || !["127.0.0.1", "localhost"].includes(new URL(base).hostname)) {
  console.error("RELAY_URL must be loopback");
  process.exit(1);
}

function record(id, result, detail) {
  checks.push({ id, result, detail });
  console.info(`${result.toUpperCase()} ${id}: ${detail}`);
}
const pass = (id, ok, detail) => record(id, ok ? "pass" : "fail", detail);
function code(res) {
  const error = res.json && typeof res.json.error === "string" ? res.json.error : "";
  return error ? `${res.status} ${error}` : String(res.status);
}
function psql(sql) {
  try {
    return execFileSync(
      "docker",
      ["compose", "--project-name", process.env.COMPOSE_PROJECT, "--env-file", ".env.example", "-f", "compose.yaml", "-f", process.env.COMPOSE_OVERRIDE, "exec", "-T", "postgres", "psql", "-U", "knowslink", "-d", "knowslink", "-v", "ON_ERROR_STOP=1", "-q", "-t", "-A", "-c", sql],
      { encoding: "utf8", stdio: ["pipe", "pipe", "pipe"] },
    ).trim();
  } catch (error) {
    throw new Error(`psql failed ${error.status ?? "error"}`);
  }
}
function keyFromGo(encoded) {
  const seed = Buffer.from(encoded, "base64").subarray(0, 32);
  const pkcs8 = Buffer.concat([Buffer.from("302e020100300506032b657004220420", "hex"), seed]);
  return createPrivateKey({ key: pkcs8, format: "der", type: "pkcs8" });
}
async function call(method, path, token, body, claim) {
  const headers = { "Content-Type": body === undefined ? "text/plain" : "application/json" };
  if (token) headers.Authorization = `Bearer ${token}`;
  if (claim) headers["X-Execution-Claim"] = claim;
  const response = await fetch(`${base}${path}`, {
    method,
    headers,
    body: body === undefined ? undefined : typeof body === "string" ? body : JSON.stringify(body),
    redirect: "manual",
  });
  const text = await response.text();
  let json = null;
  try {
    json = JSON.parse(text);
  } catch {
    json = null;
  }
  return { status: response.status, text, json };
}
function build(privateKey, fields) {
  const message = {
    v: "relay.v1", id: fields.id || uuid7(), from: fields.from, to: fields.to, intent: fields.intent,
    body: fields.body, deliver: fields.deliver, exp: fields.exp,
    idempotency_key: randomBytes(18).toString("base64url"),
    sig: { alg: "Ed25519", kid: "key1", value: "" },
  };
  if (fields.reply_to) message.reply_to = fields.reply_to;
  message.sig.value = sign(null, signingBytes(message), privateKey).toString("base64url");
  return message;
}
function plus(ms) {
  return new Date(Date.now() + ms).toISOString().replace(/\.\d{3}Z$/, "Z");
}
function cell(expr) {
  return psql(`SELECT (${expr}) FROM relay_state`);
}
async function waitHealth() {
  for (let i = 0; i < 30; i++) {
    const health = await call("GET", "/healthz").catch(() => ({ status: 0 }));
    if (health.status === 200) return;
    await new Promise((resolve) => setTimeout(resolve, 1000));
  }
  throw new Error("relay health failed");
}
async function main() {
  await waitHealth();
  const fixture = JSON.parse(await readFile("internal/relay/testdata/legacy_claims.json", "utf8"));
  const state = fixture.State;
  const hash = (token) => createHash("sha256").update(token).digest("base64url");
  const agentA = "legacy_agent_a";
  const agentB = "legacy_agent_b";
  const ownerB = "legacy_owner_b";
  if (hash(agentA) !== state.Agents.agent_a.Credential || hash(agentB) !== state.Agents.agent_b.Credential || hash(ownerB) !== state.Owners.owner_b.Credential) {
    throw new Error("fixture credential mismatch");
  }
  const privateA = keyFromGo(fixture.PrivateA);
  const privateB = keyFromGo(fixture.PrivateB);
  const soon = plus(240000);
  for (const message of Object.values(state.Messages)) message.Receipt.exp = soon;
  for (const gate of Object.values(state.Gates)) gate.Exp = soon;
  const payload = JSON.stringify(state);
  if (payload.includes("$qa$") || payload.includes("'") || payload.includes('"Deliver"')) throw new Error("fixture payload rejected");
  const loaded = psql(`WITH u AS (UPDATE relay_state SET data = $qa$${payload}$qa$::jsonb, epoch = epoch + 1, clock = clock_timestamp() RETURNING data) SELECT (SELECT count(*) FROM u)::text || '|' || (SELECT count(*) FROM u, jsonb_each(data->'Messages'))::text || '|' || (SELECT bool_and(NOT jsonb_exists(value, 'Deliver')) FROM u, jsonb_each(data->'Messages'))::text || '|' || (SELECT bool_and(value->'Receipt'->>'transport' = 'delivered' AND value->>'Claimed' = 'true' AND length(COALESCE(value->>'ClaimToken','')) > 0) FROM u, jsonb_each(data->'Messages') WHERE key IN ('${human}','${gated}','${unrouted}'))::text || '|' || (SELECT data->'Gates'->'${legacyGate}'->>'State' FROM u) || '|' || (SELECT COALESCE(data->'Gates'->'${legacyGate}'->>'Consumed','false') FROM u)`);
  pass("FIX-load-fixture", loaded === "1|4|true|true|approved|false", loaded);
  const auth = await call("POST", "/v1/owner-revoke", agentB, {});
  pass("REG-current-auth", auth.status === 401 && auth.json?.error === "invalid_auth", code(auth));
  const queryBody = { window: { start: "2026-10-03T10:00:00Z", end: "2026-10-03T11:00:00Z" }, granularity_min: 30 };
  async function boundaries(name, id) {
    const claim = state.Messages[id].ClaimToken;
    const digest = state.Messages[id].Receipt.digest;
    const authorized = await call("POST", "/v1/authorize", agentB, { id, claim });
    const approval = await call("POST", "/v1/send", agentB, build(privateB, { from: "agent_b", to: "agent_b", intent: "relay.approval.request", deliver: "human", reply_to: id, exp: plus(60000), body: { reason: "judgment_required", request_digest: digest } }), claim);
    const result = await call("POST", "/v1/send", agentB, build(privateB, { from: "agent_b", to: "agent_a", intent: "relay.result", deliver: "agent", reply_to: id, exp: plus(60000), body: { status: "denied" } }), claim);
    const denied = [authorized, approval, result].every((res) => res.status === 403 && res.json?.error === "sender_not_allowed");
    pass(name, denied, `authorize ${code(authorized)}, H ${code(approval)}, R ${code(result)}`);
  }
  await boundaries("FIX-legacy-human", human);
  await boundaries("FIX-legacy-unrouted", unrouted);
  const consume = await call("POST", "/v1/gate-consume", agentB, { id: legacyGate, claim: state.Messages[gated].ClaimToken });
  pass("FIX-legacy-consume", consume.status === 403 && consume.json?.error === "sender_not_allowed", code(consume));
  const unchanged = [human, unrouted].every((id) => cell(`COALESCE(data->'Messages'->'${id}'->>'Completion','')`) === "" && cell(`COALESCE(data->'Messages'->'${id}'->>'Deliver','')`) === "" && cell(`COALESCE(data->'Messages'->'${id}'->>'ClaimToken','') <> ''`) === "t");
  pass("FIX-legacy-unchanged", unchanged && cell(`COALESCE(data->'Gates'->'${legacyGate}'->>'Consumed','false')`) === "false" && cell(`data->'Gates'->'${legacyGate}'->>'State'`) === "approved", `unchanged ${unchanged}`);
  const sent = await call("POST", "/v1/send", agentA, build(privateA, { from: "agent_a", to: "agent_b", intent: "schedule.query", deliver: "agent", exp: plus(180000), body: queryBody }));
  const lease = await call("POST", "/v1/pull", agentB, {});
  const live = sent.json?.id;
  const persisted = await call("POST", "/v1/persist", agentB, { id: live, token: lease.json?.lease_token });
  const acked = await call("POST", "/v1/ack", agentB, { id: live, token: lease.json?.lease_token });
  const claimedLive = await call("POST", "/v1/claim", agentB, { id: live });
  const authorized = await call("POST", "/v1/authorize", agentB, { id: live, claim: claimedLive.json?.claim });
  pass("REG-agent-delivery", sent.status === 200 && lease.json?.lease_token && persisted.status === 200 && acked.status === 200 && claimedLive.status === 200 && authorized.status === 200 && authorized.json?.executable === false && authorized.json?.disclosure === false, `send ${sent.status}, persist ${persisted.status}, ack ${acked.status}, claim ${claimedLive.status}, authorize ${authorized.status}`);
  const approval = await call("POST", "/v1/send", agentB, build(privateB, { from: "agent_b", to: "agent_b", intent: "relay.approval.request", deliver: "human", reply_to: live, exp: plus(120000), body: { reason: "judgment_required", request_digest: sent.json?.digest } }), claimedLive.json?.claim);
  const gatePage = await call("GET", `/owner/gates/${approval.json?.id}`, ownerB);
  const agentPage = await call("GET", `/owner/gates/${approval.json?.id}`, agentB);
  const csrf = /name="csrf" value="([^"]+)"/.exec(gatePage.text)?.[1];
  const bad = await fetch(`${base}/owner/gates/${approval.json?.id}`, { method: "POST", redirect: "manual", headers: { Authorization: `Bearer ${ownerB}`, "Content-Type": "application/x-www-form-urlencoded" }, body: new URLSearchParams({ csrf: "bad", decision: "approve" }) });
  const approved = await fetch(`${base}/owner/gates/${approval.json?.id}`, { method: "POST", redirect: "manual", headers: { Authorization: `Bearer ${ownerB}`, "Content-Type": "application/x-www-form-urlencoded" }, body: new URLSearchParams({ csrf, decision: "approve" }) });
  const consumed = await call("POST", "/v1/gate-consume", agentB, { id: approval.json?.id, claim: claimedLive.json?.claim });
  const done = await call("POST", "/v1/send", agentB, build(privateB, { from: "agent_b", to: "agent_a", intent: "relay.result", deliver: "agent", reply_to: live, exp: plus(120000), body: { status: "denied" } }), claimedLive.json?.claim);
  const completion = await call("GET", `/v1/receipts/${live}`, agentA);
  pass("REG-owner-approval-result", approval.status === 200 && gatePage.status === 200 && gatePage.text.includes("상태: pending") && gatePage.text.includes("disclosure policy absent") && agentPage.status === 401 && bad.status === 403 && approved.status === 303 && consumed.status === 200 && consumed.json?.executable === false && consumed.json?.disclosure === false && done.status === 200 && completion.json?.completion === "denied", `H ${approval.status}, page ${gatePage.status}, agent ${agentPage.status}, csrf ${bad.status}, approve ${approved.status}, consume ${consumed.status}, result ${done.status}, completion ${completion.json?.completion}`);
  const still = await call("POST", "/v1/authorize", agentB, { id: human, claim: state.Messages[human].ClaimToken });
  pass("FIX-legacy-still-denied", still.status === 403 && still.json?.error === "sender_not_allowed", code(still));
  await writeFile(new URL("./probe-results.json", import.meta.url), JSON.stringify({ checks }, null, 2) + "\n");
  const failed = checks.filter((item) => item.result === "fail").length;
  console.info(`PROBE fail ${failed} pass ${checks.length - failed}`);
  if (failed) process.exit(1);
}
main().catch((error) => {
  console.error(`PROBE error ${error.message}`);
  process.exit(1);
});
