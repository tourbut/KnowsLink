// SAR-AUTOSETUP-001-QA fallback: drive the MCP servers that the installer registered (read from the installed plugin config) with a scripted client.
// This replaces the model/host UI only; it does NOT prove the host's own session or approval dialog (reported as unverified).
import { readFileSync, appendFileSync, writeFileSync } from "node:fs";
import { pathToFileURL } from "node:url";
import { join } from "node:path";
const root = process.argv[2], out = join(root, "work", "mcp-b-"+process.argv[3]+".jsonl");
const sdk = (p) => import(pathToFileURL(join(root, "src/adapters/node_modules/@modelcontextprotocol/sdk/dist/esm", p)).href);
const { Client } = await sdk("client/index.js");
const { StdioClientTransport, getDefaultEnvironment } = await sdk("client/stdio.js");
const { ElicitRequestSchema } = await sdk("types.js");
writeFileSync(out, "");
const log = (o) => { appendFileSync(out, JSON.stringify({ t: new Date().toISOString(), ...o }) + "\n"); console.log(JSON.stringify(o).slice(0, 600)); };
const cfg = (box, inst, file) => JSON.parse(readFileSync(join(root, "work", box, "local/KnowsLink", inst, "knowslink", file), "utf8")).mcpServers.knowslink;
const B = cfg("codex-B", "codex-01", ".mcp.json");
const agentOf = (box, inst) => JSON.parse(readFileSync(join(root, "work", box, "local/KnowsLink", inst + "-private/agent.json"), "utf8")).agent;
const idB = agentOf("codex-B", "codex-01"); const idA = "agent_784f29eabfa9016f153394";
async function open(spec, answer) {
  const transport = new StdioClientTransport({ command: spec.command, args: spec.args, env: { ...getDefaultEnvironment(), ...spec.env }, stderr: "ignore" });
  const client = new Client({ name: "qa-scripted", version: "0" }, { capabilities: answer ? { elicitation: {} } : {} });
  if (answer) client.setRequestHandler(ElicitRequestSchema, async (req) => { const m = req.params.message; const code = /(\d{4}) 를 입력하세요/.exec(m)?.[1] ?? null; log({ dialog: m.slice(0, 500), code }); return { action: "accept", content: { code } }; });
  await client.connect(transport);
  const call = async (name, args) => { const r = await client.callTool({ name, arguments: args }, undefined, { timeout: 30000 }); const text = r.content?.[0]?.text ?? ""; let json = null; try { json = JSON.parse(text); } catch {} log({ tool: name, args: JSON.stringify(args).slice(0, 200), isError: !!r.isError, text: text.slice(0, 500) }); return json ?? text; };
  return { call, close: () => client.close() };
}

const b = await open(B, true);
await b.call("knowslink_session", { action: "on" });
const key = () => "qa-" + Math.random().toString(16).slice(2) + "-0123456789";
if (process.argv[3] === "recv") await b.call("knowslink_text_receive", {});
if (process.argv[3] === "send") { await b.call("knowslink_text_send", { peer: idA, text: "AUTOSETUP-QA2 합성->Claude 첫 발송", idempotency_key: key(), confirmed: true }); }
await b.close();
