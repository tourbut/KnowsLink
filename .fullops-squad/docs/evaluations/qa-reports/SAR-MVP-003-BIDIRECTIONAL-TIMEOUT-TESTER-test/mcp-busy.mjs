// Stalled trial receive fails near 10s, an overlapping call is busy, and the next call is not busy.
import { createServer } from "node:http";
import { generateKeyPairSync } from "node:crypto";
import { mkdtemp, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join } from "node:path";
import { pathToFileURL } from "node:url";

const root = process.argv[2];
const sdk = pathToFileURL(join(root, "node_modules/@modelcontextprotocol/sdk/dist/esm/client/index.js")).href;
const stdio = pathToFileURL(join(root, "node_modules/@modelcontextprotocol/sdk/dist/esm/client/stdio.js")).href;
const { Client } = await import(sdk);
const { StdioClientTransport } = await import(stdio);
const directory = await mkdtemp(join(tmpdir(), "kl-timeout-busy-"));
const key = join(directory, "key.pem");
await writeFile(key, generateKeyPairSync("ed25519").privateKey.export({ format: "pem", type: "pkcs8" }), { mode: 0o600 });
let hits = 0;
const server = createServer((_request, response) => {
  hits += 1;
  response.writeHead(200, { "content-type": "application/json" });
  if (hits === 1) {
    response.write('{"partial":true}');
    return;
  }
  response.end("{}");
});
await new Promise((ok) => server.listen(0, "127.0.0.1", ok));
const env = {
  PATH: process.env.PATH || "",
  TMPDIR: directory,
  HOME: directory,
  KNOWSLINK_MODE: "test-loopback",
  RELAY_URL: `http://127.0.0.1:${server.address().port}`,
  AGENT_CREDENTIAL: "test-only",
  AGENT_ID: "trial_codex",
  AGENT_KID: "key1",
  AGENT_KEY_FILE: key,
  KNOWSLINK_TEST_PEER: "trial_grok",
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
const client = new Client({ name: "timeout-busy-check", version: "0" });
await client.connect(transport);
const stateOf = (result) => {
  const text = result.content?.[0]?.text || "";
  try {
    return JSON.parse(text).state ?? text;
  } catch {
    return text.slice(0, 80);
  }
};
const call = (options) =>
  client.callTool({ name: "knowslink_test_receive", arguments: {} }, undefined, options).then(stateOf);
const started = Date.now();
const firstPromise = call({ timeout: 30000 }).then((state) => ({ state, ms: Date.now() - started }));
await new Promise((ok) => setTimeout(ok, 300));
const concurrent = await call();
const first = await firstPromise;
const after = await call();
const stderrHasKey = stderr.includes("PRIVATE KEY");
const pass =
  first.state === "failed" &&
  first.ms >= 9900 &&
  first.ms < 12000 &&
  concurrent === "busy" &&
  after !== "busy" &&
  !stderrHasKey;
console.log(JSON.stringify({ first, concurrent, after, hits, stderrHasKey, pass }));
await client.close();
server.closeAllConnections();
server.close();
process.exit(pass ? 0 : 1);
