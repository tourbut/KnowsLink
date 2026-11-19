// Independent probe: one stdio process, overlapping pull_once must be busy and must not touch the relay.
import assert from "node:assert/strict";
import { createServer } from "node:http";
import { once } from "node:events";
import { mkdtemp, rm, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join } from "node:path";
import { Client } from "/tmp/sar-mvp-002-clone/adapters/node_modules/@modelcontextprotocol/sdk/dist/esm/client/index.js";
import { StdioClientTransport } from "/tmp/sar-mvp-002-clone/adapters/node_modules/@modelcontextprotocol/sdk/dist/esm/client/stdio.js";

const plugin = process.argv[2];
if (!plugin) {
  console.error("usage: concurrent.mjs <plugin.js>");
  process.exit(2);
}

const body = (state) => [
  {
    type: "text",
    text: JSON.stringify({
      state,
      transport: "pull",
      actualConnection: "held",
      webhook: false,
      evidenceFetch: false,
    }),
  },
];

let releaseHold;
const hold = new Promise((resolve) => {
  releaseHold = resolve;
});
let seen = 0;
const server = createServer((_request, response) => {
  seen += 1;
  hold.then(() => {
    response.writeHead(500, { "content-type": "application/json" });
    response.end("{}");
  });
});
server.listen(0, "127.0.0.1");
await once(server, "listening");
const address = server.address();
assert(address && typeof address !== "string");
const directory = await mkdtemp(join(tmpdir(), "knowslink-busy-"));
const keyFile = join(directory, "key.pem");
await writeFile(keyFile, "synthetic-probe-marker\n", { mode: 0o600 });

const transport = new StdioClientTransport({
  command: process.execPath,
  args: [plugin],
  env: {
    KNOWSLINK_MODE: "synthetic-loopback",
    RELAY_URL: `http://127.0.0.1:${address.port}`,
    AGENT_CREDENTIAL: "synthetic-probe",
    AGENT_ID: "synthetic_b",
    AGENT_KID: "key1",
    AGENT_KEY_FILE: keyFile,
  },
  stderr: "pipe",
});
const client = new Client({ name: "busy-probe", version: "1.0.0" });
let stderr = "";
transport.stderr?.on("data", (chunk) => {
  stderr += chunk.toString();
});

try {
  await client.connect(transport);
  const first = client.callTool(
    { name: "knowslink_pull_once", arguments: {} },
    undefined,
    { timeout: 20000 },
  );
  const started = Date.now();
  while (seen < 1) {
    if (Date.now() - started > 5000) break;
    await new Promise((resolve) => setTimeout(resolve, 20));
  }
  assert.equal(seen, 1, "first pull did not reach the loopback relay");
  const status = await client.callTool(
    { name: "knowslink_status", arguments: {} },
    undefined,
    { timeout: 5000 },
  );
  const secondStarted = Date.now();
  const second = await client.callTool(
    { name: "knowslink_pull_once", arguments: {} },
    undefined,
    { timeout: 5000 },
  );
  const secondMs = Date.now() - secondStarted;
  assert.equal(second.isError, true);
  assert.deepEqual(second.content, body("busy"));
  assert.equal(seen, 1, "busy pull contacted the relay");
  assert.ok(secondMs < 2000, `busy response took ${secondMs}ms`);
  releaseHold();
  const firstResult = await first;
  assert.equal(firstResult.isError, true);
  assert.deepEqual(firstResult.content, body("failed"));
  const third = await client.callTool(
    { name: "knowslink_pull_once", arguments: {} },
    undefined,
    { timeout: 5000 },
  );
  assert.equal(third.isError, true);
  assert.deepEqual(third.content, body("failed"));
  assert.equal(seen, 2, "the pull after busy cleared did not reach the relay");
  assert.equal(status.isError, false);
  assert.deepEqual(status.content, body("synthetic_only"));
  const visible = JSON.stringify({ status, second, firstResult, third, stderr });
  assert.equal(visible.includes("synthetic-probe"), false);
  assert.equal(visible.includes("synthetic-probe-marker"), false);
  assert.equal(stderr, "");
  console.info(
    `PASS: overlap busy in ${secondMs}ms without a second relay request; status synthetic_only; later pull failed; seen=${seen}`,
  );
} finally {
  releaseHold();
  await client.close().catch(() => {});
  server.close();
  await once(server, "close");
  await rm(directory, { recursive: true, force: true });
}
