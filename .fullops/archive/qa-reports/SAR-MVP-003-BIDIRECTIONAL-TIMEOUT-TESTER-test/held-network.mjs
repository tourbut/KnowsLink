// Held MCP calls a canary relay URL and must open no socket.
import { createServer } from "node:http";
import { generateKeyPairSync } from "node:crypto";
import { mkdtemp, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join } from "node:path";
import { pathToFileURL } from "node:url";

const root = process.argv[2];
const { Client } = await import(pathToFileURL(join(root, "node_modules/@modelcontextprotocol/sdk/dist/esm/client/index.js")).href);
const { StdioClientTransport } = await import(pathToFileURL(join(root, "node_modules/@modelcontextprotocol/sdk/dist/esm/client/stdio.js")).href);
const directory = await mkdtemp(join(tmpdir(), "kl-timeout-held-"));
const key = join(directory, "key.pem");
await writeFile(key, generateKeyPairSync("ed25519").privateKey.export({ format: "pem", type: "pkcs8" }), { mode: 0o600 });
let hits = 0;
const server = createServer((_request, response) => {
  hits += 1;
  response.end("no");
});
await new Promise((ok) => server.listen(0, "127.0.0.1", ok));
const env = {
  PATH: process.env.PATH || "",
  TMPDIR: directory,
  HOME: directory,
  KNOWSLINK_MODE: "held",
  RELAY_URL: `http://127.0.0.1:${server.address().port}`,
  AGENT_CREDENTIAL: "held-not-used",
  AGENT_ID: "trial_codex",
  AGENT_KID: "key1",
  AGENT_KEY_FILE: key,
  KNOWSLINK_TEST_PEER: "trial_grok",
  CF_ACCESS_CLIENT_ID: "held-client",
  CF_ACCESS_CLIENT_SECRET: "example-access",
};
let stderr = "";
const transport = new StdioClientTransport({
  command: process.execPath,
  args: [join(root, "dist/plugin.js")],
  env,
  stderr: "pipe",
});
transport.stderr.on("data", (chunk) => {
  stderr += chunk.toString();
});
const client = new Client({ name: "timeout-held-check", version: "0" });
await client.connect(transport);
const call = async (name, arguments_) => {
  const result = await client.callTool({ name, arguments: arguments_ });
  const text = result.content?.[0]?.text || "";
  return JSON.parse(text);
};
const status = await call("knowslink_status", {});
const send = await call("knowslink_test_send", { text: "should-not-send", idempotency_key: "qa-held-no-network1" });
const receive = await call("knowslink_test_receive", {});
const pull = await call("knowslink_pull_once", {});
const bodies = [status, send, receive, pull];
const held = bodies.every((body) => body.state === "held" && body.actualConnection === "held" && body.webhook === false && body.evidenceFetch === false);
const stderrHasKey = stderr.includes("PRIVATE KEY") || stderr.includes("example-access");
const pass = hits === 0 && held && !stderrHasKey;
console.log(JSON.stringify({ hits, states: bodies.map((body) => body.state), held, stderrHasKey, pass }));
await client.close();
server.closeAllConnections();
server.close();
process.exit(pass ? 0 : 1);
