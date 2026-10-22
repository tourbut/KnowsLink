// SAR-CODEX-RECV-001-DEV V02 driver: an App Server client (the protocol a --remote TUI uses) for one Codex session.
// usage: node drive.mjs <name> <install-dir> <workdir>; append {"turn":"...","answers":["accept"]} lines to cmd-<name>.jsonl; events go to log-<name>.jsonl.
// The synthetic owner's dialog answers come only from the command file, never from the model.
import { spawn } from "node:child_process";
import { appendFileSync, existsSync, readFileSync, writeFileSync } from "node:fs";
import { join } from "node:path";

const [name, inst, cwd] = process.argv.slice(2);
const dir = import.meta.dirname;
const cmdFile = join(dir, `cmd-${name}.jsonl`);
const logFile = join(dir, `log-${name}.jsonl`);
writeFileSync(cmdFile, "");
const log = (o) =>
  appendFileSync(logFile, JSON.stringify({ t: new Date().toISOString(), ...o }) + "\n");
const mcp = JSON.parse(readFileSync(join(inst, "knowslink", ".mcp.json"), "utf8")).mcpServers.knowslink;
const fwd = (p) => p.replaceAll("\\", "/");
const env = Object.entries(mcp.env)
  .map(([k, v]) => `${k}='${fwd(v)}'`)
  .join(", ");
const codexJs = "C:/Users/shin/AppData/Roaming/npm/node_modules/@openai/codex/bin/codex.js";
const child = spawn(
  process.execPath,
  [
    codexJs,
    "app-server",
    "-c", `mcp_servers.knowslink.command='${fwd(mcp.command)}'`,
    "-c", `mcp_servers.knowslink.args=['${fwd(mcp.args[0])}']`,
    "-c", `mcp_servers.knowslink.env={${env}}`,
    "-c", `model_reasoning_effort='low'`,
  ],
  { stdio: ["pipe", "pipe", "pipe"] },
);
child.stderr.on("data", (d) => log({ stderr: String(d).slice(0, 2000) }));
let next = 1;
const pending = new Map();
const send = (o) => child.stdin.write(JSON.stringify(o) + "\n");
const request = (method, params) =>
  new Promise((resolve, reject) => {
    const id = next++;
    pending.set(id, { resolve, reject });
    send({ id, method, params });
  });
let answers = [];
let turnDone = null;
let buf = "";
child.stdout.on("data", (d) => {
  buf += d;
  let i;
  while ((i = buf.indexOf("\n")) >= 0) {
    const line = buf.slice(0, i);
    buf = buf.slice(i + 1);
    if (!line.trim()) continue;
    const m = JSON.parse(line);
    if (m.id !== undefined && !m.method) {
      const p = pending.get(m.id);
      pending.delete(m.id);
      m.error ? p?.reject(m.error) : p?.resolve(m.result);
      continue;
    }
    if (m.id !== undefined && m.method) {
      // Server request: dialogs. The scripted owner answers KnowsLink dialogs; anything else is logged and declined.
      if (m.method === "mcpServer/elicitation/request") {
        // Codex's own MCP tool-call approval is a separate dialog: the owner allows the call; KnowsLink dialogs take the scripted answers.
        const own = m.params.message?.startsWith("KnowsLink");
        const answer = own ? (answers.shift() ?? "cancel") : "accept";
        log({ elicitation: m.params.message, serverName: m.params.serverName, knowslinkDialog: own, answer });
        // The synthetic owner reads the one-time code in the dialog and types it; "blind" is an accept without input (full-auto host).
        const code = /(\d{4}) 를 입력하세요/.exec(m.params.message ?? "")?.[1];
        const content = answer === "accept" ? (own ? { code } : {}) : null;
        send({ id: m.id, result: { action: answer === "blind" ? "accept" : answer, content } });
      } else {
        log({ serverRequest: m.method, params: m.params });
        send({ id: m.id, result: { decision: "denied" } });
      }
      continue;
    }
    const keep = ["item/completed", "turn/completed", "turn/started", "error", "thread/status/changed"];
    if (keep.includes(m.method)) {
      const item = m.params?.item;
      if (m.method === "item/completed" && item && !["mcpToolCall", "agentMessage", "userMessage"].includes(item.type)) continue;
      log({ event: m.method, params: m.params });
    }
    if (m.method === "turn/completed") turnDone?.();
  }
});
await request("initialize", { clientInfo: { name: "knowslink-v02-driver", version: "0.0.1" } });
send({ method: "initialized" });
const started = await request("thread/start", { cwd, approvalPolicy: "on-request", sandbox: "read-only" });
const threadId = started.thread.id;
log({ thread: threadId });
let done = 0;
for (;;) {
  await new Promise((r) => setTimeout(r, 1000));
  const lines = existsSync(cmdFile) ? readFileSync(cmdFile, "utf8").split("\n").filter(Boolean) : [];
  if (lines.length <= done) continue;
  const c = JSON.parse(lines[done++]);
  if (c.exit) break;
  answers = c.answers ?? [];
  log({ turn: c.turn, answers, approvalPolicy: c.approvalPolicy });
  const finished = new Promise((r) => (turnDone = r));
  await request("turn/start", { threadId, input: [{ type: "text", text: c.turn }], ...(c.approvalPolicy ? { approvalPolicy: c.approvalPolicy } : {}) });
  await finished;
  log({ turnEnd: true, unusedAnswers: answers });
}
child.kill();
process.exit(0);
