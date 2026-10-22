// SAR-CODEX-RECV-001-TESTER approval-boundary matrix.
// A scripted MCP client plays Codex: it answers elicitations the way a TUI owner, an Esc-cancel, a full-auto host or a client without elicitation would.
// Each case spawns a fresh installed MCP process (= a new Codex session). Probe A is a plain (non-Codex) connected member that reads what actually reached the relay.
// usage: node approval-matrix.mjs <tools-dir> <out.jsonl>
import { readFileSync, existsSync, appendFileSync, writeFileSync } from "node:fs";
import { pathToFileURL } from "node:url";
import { join } from "node:path";

const [tools, outFile] = process.argv.slice(2);
const sdk = (p) => import(pathToFileURL(join(tools, "src/adapters/node_modules/@modelcontextprotocol/sdk/dist/esm", p)).href);
const { Client } = await sdk("client/index.js");
const { StdioClientTransport, getDefaultEnvironment } = await sdk("client/stdio.js");
const { ElicitRequestSchema } = await sdk("types.js");
writeFileSync(outFile, "");
const log = (o) => appendFileSync(outFile, JSON.stringify({ t: new Date().toISOString(), ...o }) + "\n");

const fx = JSON.parse(readFileSync(join(tools, "fixture-agents.json"), "utf8"));
const cfg = (inst, file) => JSON.parse(readFileSync(join(tools, inst, "knowslink", file), "utf8")).mcpServers.knowslink;
const B = cfg("inst-b", ".mcp.json");
const A = cfg("inst-a", "config/local.mcp.json");
const lockOf = (folder) => join(folder, "session.lock");
const sleep = (ms) => new Promise((r) => setTimeout(r, ms));

// answer: (message) => {action, content} | "none" (client without elicitation) ; an async function may wait on a gate
async function open(spec, { elicit = true, answer } = {}) {
  const transport = new StdioClientTransport({ command: spec.command, args: spec.args, env: { ...getDefaultEnvironment(), ...spec.env }, stderr: "ignore" });
  const client = new Client({ name: "qa-codex-sim", version: "0" }, { capabilities: elicit ? { elicitation: {} } : {} });
  const dialogs = [];
  if (elicit)
    client.setRequestHandler(ElicitRequestSchema, async (req) => {
      const message = req.params.message;
      const code = /(\d{4}) 를 입력하세요/.exec(message)?.[1] ?? null;
      const reply = await answer(message, code);
      dialogs.push({ message, code, reply });
      return reply;
    });
  await client.connect(transport);
  const call = async (name, args) => {
    try {
      const r = await client.callTool({ name, arguments: args }, undefined, { timeout: 30000 });
      const text = r.content?.[0]?.text ?? "";
      let json = null;
      try { json = JSON.parse(text); } catch {}
      return { isError: !!r.isError, text, json };
    } catch (e) {
      return { isError: true, text: String(e.message ?? e), json: null, thrown: true };
    }
  };
  return { client, call, dialogs, close: () => client.close() };
}
const typed = (code) => ({ action: "accept", content: { code } });
const owners = {
  decline: () => ({ action: "decline" }),
  cancel: () => ({ action: "cancel" }),
  blindAccept: () => ({ action: "accept", content: {} }), // full-auto host: accept without showing
  blindAcceptNoContent: () => ({ action: "accept" }),
  wrongCode: (_m, code) => typed(String(Number(code) === 9999 ? 1000 : Number(code) + 1)),
  emptyCode: () => typed(""),
  spaces: () => typed("    "),
  right: (_m, code) => typed(code),
  rightSpaced: (_m, code) => typed(` ${code} `),
};
let results = [];
const check = (id, name, ok, detail) => {
  results.push({ id, name, ok });
  log({ check: id, name, ok, detail });
  console.log(`${ok ? "PASS" : "FAIL"} ${id} ${name}`);
};
const hasLock = (folder) => existsSync(lockOf(folder));
const lockInfo = (folder) => (hasLock(folder) ? JSON.parse(readFileSync(lockOf(folder), "utf8")) : null);

// Probe A: plain member that pulls whatever is queued for A at the relay.
async function drainA() {
  const a = await open(A, { elicit: false });
  const got = [];
  for (let i = 0; i < 6; i++) {
    const r = await a.call("knowslink_text_receive", {});
    if (r.json?.state === "received") got.push(r.json.message); else break;
  }
  await a.close();
  return got;
}
const sendFromA = async (text) => {
  const a = await open(A, { elicit: false });
  const r = await a.call("knowslink_text_send", { peer: fx.b.agent, text, idempotency_key: "qa-a2b-" + Math.random().toString(16).slice(2) + "xxxxxxxx", confirmed: true });
  await a.close();
  return r;
};
const mkKey = () => "qa-b2a-" + Math.random().toString(16).slice(2) + "yyyyyyyy";
const sendB = (b, text, extra = {}) => b.call("knowslink_text_send", { peer: fx.a.agent, text, idempotency_key: mkKey(), confirmed: true, ...extra });
const cleanLock = async () => { /* every case ends with off or close; assert below */ };


const { readdirSync } = await import("node:fs");
const folder = B.env.KNOWSLINK_AGENT_FOLDER;
const inboxFiles = () => { try { return readdirSync(join(folder, "inbox")).length; } catch { return -1; } };
await drainA();
let found = null;
for (let i = 0; i < 8 && !found; i++) {
  const b = await open(B, { answer: owners.right });
  const [, on] = await Promise.all([b.call("knowslink_session", { action: "status" }), b.call("knowslink_session", { action: "on" })]);
  const st = await b.call("knowslink_session", { action: "status" });
  log({ round: i, on: on.json?.state, status: st.json?.state });
  if (on.json?.state === "bound" && st.json?.state !== "bound") {
    found = b;
    const off = await b.call("knowslink_session", { action: "off" });
    const before = inboxFiles();
    const lockAfterOff = lockInfo(folder);
    await sendFromA("ORPHAN probe: sent after the owner turned the session off");
    await sleep(15000);
    const after = inboxFiles();
    const rv = await b.call("knowslink_text_receive", {});
    log({ ORPHAN: true, offText: off.text, lockAfterOff, inboxBefore: before, inboxAfter: after, receiveText: rv.text.slice(0, 160) });
    console.log(JSON.stringify({ offText: off.text, lockAfterOff, inboxBefore: before, inboxAfter: after, receive: rv.json?.state }));
  }
  await b.close();
  await sleep(500);
}
if (!found) console.log("no orphan in 8 rounds");
await drainA();
