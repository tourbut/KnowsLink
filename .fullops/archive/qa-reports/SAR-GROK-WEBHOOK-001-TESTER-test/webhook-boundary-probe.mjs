// Independent tester probe for SAR-GROK-WEBHOOK-001 (b86f397). Real fetch against loopback servers; nothing leaves the host.
import http from "node:http";
import net from "node:net";
import { spawnSync } from "node:child_process";
import { chmodSync, mkdtempSync, readFileSync, readdirSync, symlinkSync, writeFileSync, mkdirSync } from "node:fs";
import { tmpdir } from "node:os";
import { join } from "node:path";

const D = process.argv[2]; // .../repo/adapters/dist
const { GrokWake } = await import(join(D, "grok-wake.js"));
const { Inbox, AutoReceiver } = await import(join(D, "inbox.js"));

const results = [];
const ok = (name, cond, detail = "") => {
  results.push({ name, pass: !!cond, detail: String(detail) });
  console.log(`${cond ? "PASS" : "FAIL"} ${name}${detail ? " :: " + detail : ""}`);
};
const info = (name, v) => console.log(`INFO ${name} :: ${JSON.stringify(v)}`);

const KEY = "sk_probe_SECRET_KEY_0123456789";
const RID = "rt_PROBE-routine_Id9";
const URL_OK = `https://api2.cursor.sh/automations/webhook/${RID}`;
const root = mkdtempSync(join(tmpdir(), "kl-wh-probe-"));
const cfg = join(root, "webhook.json");
const writeCfg = (c, mode = 0o600) => { writeFileSync(cfg, typeof c === "string" ? c : JSON.stringify(c), { mode }); chmodSync(cfg, mode); };
writeCfg({ url: URL_OK, key: KEY });

// loopback servers
const mk = (handler) => new Promise((res) => {
  const log = [];
  const s = http.createServer((req, rsp) => {
    let b = ""; req.on("data", (c) => (b += c));
    req.on("end", () => { log.push({ method: req.method, url: req.url, headers: req.headers, body: b }); handler(req, rsp, log); });
  });
  s.listen(0, "127.0.0.1", () => res({ s, log, port: s.address().port }));
});
const B = await mk((_q, r) => { r.writeHead(200); r.end("B"); }); // redirect target
let mode = "ok";
const A = await mk((_q, r) => {
  if (mode === "ok") { r.writeHead(200); r.end("echo " + KEY); }
  else if (mode === "302") { r.writeHead(302, { location: `http://127.0.0.1:${B.port}/steal` }); r.end(); }
  else if (mode === "hang") { /* never reply */ }
  else if (mode === "reset") { r.socket.destroy(); }
  else { r.writeHead(Number(mode)); r.end("x"); }
});
const local = (port) => (u, init) => fetch(`http://127.0.0.1:${port}${new URL(u).pathname}`, init);
const wakeTo = (send, ms = 8000, file = cfg) => new GrokWake("", 0, "", ms, file, send);
const UUIDS = ["0198c2a4-1b2c-7d3e-8f40-123456789abc", "0198c2a4-1b2c-7d3e-8f40-123456789abd"];
const all = [];
const ring = async (w, ids = UUIDS, n = 2) => { const r = await w.ring(ids, n); all.push(r); return r; };

// 1 happy path through real fetch
mode = "ok";
let r = await ring(wakeTo(local(A.port)));
ok("2xx -> accepted_unverified, response body (with key echo) not surfaced", r.state === "accepted_unverified" && !JSON.stringify(r).includes(KEY), JSON.stringify(r));
let q = A.log.at(-1);
ok("exactly one POST, fixed path, JSON", A.log.length === 1 && q.method === "POST" && q.url === `/automations/webhook/${RID}` && q.headers["content-type"] === "application/json");
ok("auth headers both carry the key", q.headers.authorization === `Bearer ${KEY}` && q.headers["x-automation-key"] === KEY);
const body = JSON.parse(q.body);
ok("body keys exactly event/ids/pending/next", JSON.stringify(Object.keys(body).sort()) === '["event","ids","next","pending"]' && body.event === "knowslink_text_arrived" && body.pending === 2 && body.ids.length === 2);
ok("body carries no key/url", !q.body.includes(KEY) && !q.body.includes(RID));
r = await ring(wakeTo(local(A.port)), [...Array(15)].map((_, i) => `0198c2a4-1b2c-7d3e-8f40-12345678${String(i).padStart(4, "0")}`).concat(["not-a-uuid", "../x"]), 17);
ok("ids capped to 10, non-UUID dropped", JSON.parse(A.log.at(-1).body).ids.length === 10 && JSON.parse(A.log.at(-1).body).ids.every((i) => /^[0-9a-f-]{36}$/.test(i)));

// 2 redirect not followed, secret not forwarded
mode = "302"; const bBefore = B.log.length;
r = await ring(wakeTo(local(A.port)));
ok("302 -> rejected, not followed, nothing reaches redirect target", r.state === "rejected" && r.error === "webhook HTTP 302" && B.log.length === bBefore, JSON.stringify(r));

// 3 status mapping
for (const [m, st] of [["401", "rejected"], ["429", "rejected"], ["500", "uncertain"], ["503", "uncertain"]]) {
  mode = m; r = await ring(wakeTo(local(A.port)));
  ok(`HTTP ${m} -> ${st}`, r.state === st && r.error === `webhook HTTP ${m}`, JSON.stringify(r));
}
// 4 single attempt on timeout / reset
mode = "hang"; let n0 = A.log.length; let t0 = Date.now();
r = await ring(wakeTo(local(A.port), 300));
ok("silent peer -> uncertain after own timeout, one attempt", r.state === "uncertain" && r.error === "webhook_no_response" && A.log.length === n0 + 1 && Date.now() - t0 < 3000, `${JSON.stringify(r)} ${Date.now() - t0}ms`);
mode = "reset"; n0 = A.log.length;
r = await ring(wakeTo(local(A.port)));
ok("connection reset after request -> uncertain, one attempt", r.state === "uncertain" && A.log.length === n0 + 1, JSON.stringify(r));
// 5 refused / dns
const closed = await new Promise((res) => { const s = net.createServer().listen(0, "127.0.0.1", () => { const p = s.address().port; s.close(() => res(p)); }); });
r = await ring(wakeTo(local(closed)));
ok("ECONNREFUSED -> retry webhook_unreachable", r.state === "retry" && r.error === "webhook_unreachable", JSON.stringify(r));
r = await ring(wakeTo((u, i) => fetch(`https://nonexistent-kl-probe.invalid${new URL(u).pathname}`, i)));
ok("real DNS NXDOMAIN (.invalid) -> retry", r.state === "retry", JSON.stringify(r));
for (const host of ["[2001:db8::1]", "192.0.2.1"]) {
  const t = Date.now();
  const x = await ring(wakeTo((u, i) => fetch(`https://${host}${new URL(u).pathname}`, i), 1500));
  info(`unroutable ${host} classification`, { ...x, ms: Date.now() - t });
}

// 6 untrusted configs never POST
const bad = [
  ["uppercase host", { url: "https://API2.cursor.sh/automations/webhook/x1", key: KEY }],
  ["trailing slash", { url: URL_OK + "/", key: KEY }],
  ["129-char id", { url: "https://api2.cursor.sh/automations/webhook/" + "a".repeat(129), key: KEY }],
  ["dot in id", { url: "https://api2.cursor.sh/automations/webhook/a.b", key: KEY }],
  ["port", { url: "https://api2.cursor.sh:443/automations/webhook/x1", key: KEY }],
  ["fragment", { url: URL_OK + "#f", key: KEY }],
  ["newline suffix", { url: URL_OK + "\n", key: KEY }],
  ["encoded slash", { url: "https://api2.cursor.sh/automations/webhook/a%2Fb", key: KEY }],
  ["other path", { url: "https://api2.cursor.sh/automations/other/x1", key: KEY }],
  ["subdomain", { url: "https://x.api2.cursor.sh/automations/webhook/x1", key: KEY }],
  ["key 4097", { url: URL_OK, key: "k".repeat(4097) }],
  ["key unicode", { url: URL_OK, key: "keykey키키키키" }],
  ["key CRLF", { url: URL_OK, key: "keykey12\r\nX-Evil: 1" }],
  ["url not string", { url: 5, key: KEY }],
  ["array", [URL_OK, KEY]],
  ["not json", "{nope"],
];
mode = "ok"; const before = A.log.length;
for (const [name, c] of bad) {
  writeCfg(c); const x = await ring(wakeTo(local(A.port)));
  const chk = await wakeTo(local(A.port)).check();
  ok(`bad config (${name}) -> retry/no POST, check failed`, x.state === "retry" && x.error === "webhook_config_unavailable" && chk.state === "failed" && chk.mode === "webhook" && chk.error === "webhook_config_invalid", JSON.stringify(x));
}
writeCfg({ url: URL_OK, key: KEY }, 0o640); let x = await ring(wakeTo(local(A.port)));
ok("group-readable file -> no POST", x.error === "webhook_config_unavailable");
writeCfg({ url: URL_OK, key: KEY }, 0o604); x = await ring(wakeTo(local(A.port)));
ok("world-readable file -> no POST", x.error === "webhook_config_unavailable");
writeCfg(JSON.stringify({ url: URL_OK, key: KEY, pad: "p".repeat(9000) }));
x = await ring(wakeTo(local(A.port)));
ok("oversize (>8192) file -> no POST", x.error === "webhook_config_unavailable");
writeCfg({ url: URL_OK, key: KEY });
const link = join(root, "link.json"); symlinkSync(cfg, link);
x = await ring(wakeTo(local(A.port), 8000, link));
ok("symlinked config -> no POST", x.error === "webhook_config_unavailable");
x = await ring(wakeTo(local(A.port), 8000, join(root, "missing.json")));
ok("missing config -> retry no POST", x.state === "retry" && x.error === "webhook_config_unavailable");
ok("no POST for any bad/untrusted config", A.log.length === before, `${A.log.length - before} extra`);
writeCfg({ url: URL_OK, key: KEY, extra: "ignored" });
x = await ring(wakeTo(local(A.port)));
ok("valid config with extra field still works (key rotation re-read per ring)", x.state === "accepted_unverified");
writeCfg({ url: URL_OK, key: "rotated_KEY_987654321" }); const nA = A.log.length;
await ring(wakeTo(local(A.port)));
ok("rotated key used on next ring without restart", A.log.at(-1).headers.authorization === "Bearer rotated_KEY_987654321" && A.log.length === nA + 1);
writeCfg({ url: URL_OK, key: KEY });

// 7 env gate / legacy preservation
ok("unset -> null", GrokWake.fromEnv({}) === null);
ok("empty webhook var behaves as unset", GrokWake.fromEnv({ KNOWSLINK_GROK_WEBHOOK_FILE: "" }) === null);
ok("relative webhook path -> invalid", GrokWake.fromEnv({ KNOWSLINK_GROK_WEBHOOK_FILE: "w.json" }) === "invalid");
const gw = await mk((_q, r) => { r.writeHead(200); r.end("{}"); });
const gwTok = join(root, "gateway.json"); writeFileSync(gwTok, JSON.stringify({ token: "gw_token_0123456789" }), { mode: 0o600 });
const AGENT = "0198c2a4-1b2c-7d3e-8f40-0000000000aa";
const envBoth = { KNOWSLINK_GROK_WEBHOOK_FILE: cfg, KNOWSLINK_GROK_WAKE_AGENT: AGENT, KNOWSLINK_GROK_GATEWAY_PORT: String(gw.port), KNOWSLINK_GROK_GATEWAY_FILE: gwTok };
const wBoth = GrokWake.fromEnv(envBoth);
// route webhook through loopback by wrapping send after construction
wBoth.send = local(A.port);
mode = "ok"; const nb = A.log.length;
x = await ring(wBoth);
ok("webhook+agent both set -> webhook only, gateway untouched", x.state === "accepted_unverified" && gw.log.length === 0 && A.log.length === nb + 1);
mode = "500";
x = await ring(wBoth);
ok("webhook failure does not fall back to sendPrompt", x.state === "uncertain" && gw.log.length === 0);
const wLegacy = GrokWake.fromEnv({ KNOWSLINK_GROK_WAKE_AGENT: AGENT, KNOWSLINK_GROK_GATEWAY_PORT: String(gw.port), KNOWSLINK_GROK_GATEWAY_FILE: gwTok });
x = await ring(wLegacy);
ok("legacy only -> sendPrompt to loopback unchanged", x.state === "accepted_unverified" && gw.log.length === 1 && gw.log[0].url === "/api/sendPrompt" && gw.log[0].headers.authorization === "Bearer gw_token_0123456789" && JSON.parse(gw.log[0].body).agentId === AGENT, JSON.stringify(x));
ok("legacy prompt has IDs only, no webhook fields", !gw.log[0].body.includes("knowslink_text_arrived") && !gw.log[0].body.includes(KEY));

// 8 AutoReceiver marker semantics with webhook mode
const folder = join(root, "conn"); mkdirSync(folder, { mode: 0o700 });
const inbox = await Inbox.open(folder);
const dir = join(folder, "inbox");
const MID = (n) => `0198c2a4-1b2c-7d3e-8f40-12345678${String(n).padStart(4, "0")}`;
const addMsg = (n) => writeFileSync(join(dir, `${MID(n)}.json`), JSON.stringify({ id: MID(n), text: "SECRET MESSAGE BODY do-not-send" }), { mode: 0o600 });
const recv = new AutoReceiver({}, inbox, async () => {}, wakeTo(local(A.port)), 10000, "poll");
const markers = () => readdirSync(dir).filter((f) => f.endsWith(".wake")).map((f) => readFileSync(join(dir, f), "utf8"));
mode = "ok"; addMsg(1); addMsg(2); let na = A.log.length;
await recv.ring(recv.wake);
ok("tick 1: one POST for the batch, markers written", A.log.length === na + 1 && markers().length === 2 && markers().every((m) => JSON.parse(m).state === "accepted_unverified"));
ok("tick 1 POST body has no message text", !A.log.at(-1).body.includes("SECRET MESSAGE"));
await recv.ring(recv.wake);
ok("tick 2: duplicate marker -> no second POST", A.log.length === na + 1);
addMsg(3); await recv.ring(recv.wake);
ok("new ID rings once; body ids only the new one, pending counts all", A.log.length === na + 2 && JSON.parse(A.log.at(-1).body).ids.join() === MID(3) && JSON.parse(A.log.at(-1).body).pending === 3);
// uncertain: marker kept, never resent
mode = "500"; addMsg(4); na = A.log.length; await recv.ring(recv.wake); await recv.ring(recv.wake);
ok("5xx -> uncertain marker kept, one POST only", A.log.length === na + 1 && JSON.parse(readFileSync(join(dir, `${MID(4)}.wake`), "utf8")).state === "uncertain");
// retry: unreachable releases marker, resumes later
const recvDown = new AutoReceiver({}, inbox, async () => {}, wakeTo(local(closed)), 10000, "poll");
addMsg(5); await recvDown.ring(recvDown.wake);
ok("unreachable -> marker released (will ring again)", !readdirSync(dir).includes(`${MID(5)}.wake`) && recvDown.status.hostWake.state === "retry" && recvDown.status.hostWake.lastError === "webhook_unreachable");
mode = "ok"; na = A.log.length; await recv.ring(recv.wake);
ok("after recovery the released ID rings", A.log.length === na + 1 && JSON.parse(A.log.at(-1).body).ids.join() === MID(5));
// no secrets in statuses/markers/results
const dump = JSON.stringify([all, recv.status, recvDown.status, markers()]);
ok("no key / routine id / message text in results, status, markers", !dump.includes(KEY) && !dump.includes(RID) && !dump.includes("SECRET MESSAGE"));

// 9 wake-check CLI
const plugin = join(D, "plugin.js");
const cli = (env) => spawnSync(process.execPath, [plugin, "wake-check"], { env: { PATH: process.env.PATH, ...env }, encoding: "utf8", timeout: 20000 });
let c = cli({ KNOWSLINK_GROK_WEBHOOK_FILE: cfg });
ok("wake-check ok file -> {state:ok,mode:webhook} exit 0", c.status === 0 && c.stdout.trim() === '{"state":"ok","mode":"webhook"}', `${c.status} ${c.stdout.trim()}`);
writeCfg({ url: "http://api2.cursor.sh/x", key: KEY });
c = cli({ KNOWSLINK_GROK_WEBHOOK_FILE: cfg });
ok("wake-check bad file -> failed exit 1, no secret", c.status === 1 && c.stdout.includes("webhook_config_invalid") && !(c.stdout + c.stderr).includes(KEY), `${c.status} ${c.stdout.trim()}`);
c = cli({ KNOWSLINK_GROK_WEBHOOK_FILE: "relative.json" });
ok("wake-check relative path -> invalid_config exit 1", c.status === 1 && c.stdout.includes("invalid_config"), `${c.status} ${c.stdout.trim()}`);
writeCfg({ url: URL_OK, key: KEY });
c = cli({ KNOWSLINK_GROK_WEBHOOK_FILE: cfg, KNOWSLINK_GROK_WAKE_AGENT: AGENT, KNOWSLINK_GROK_GATEWAY_PORT: String(gw.port), KNOWSLINK_GROK_GATEWAY_FILE: gwTok });
ok("wake-check with both vars checks webhook file only (gateway untouched)", c.status === 0 && c.stdout.includes('"mode":"webhook"') && gw.log.length === 1, `${c.stdout.trim()} gwlog=${gw.log.length}`);

const failed = results.filter((r) => !r.pass);
console.log(`\nSUMMARY pass=${results.length - failed.length} fail=${failed.length}`);
for (const s of [A, B, gw]) s.s.close();
process.exit(failed.length ? 1 : 0);
