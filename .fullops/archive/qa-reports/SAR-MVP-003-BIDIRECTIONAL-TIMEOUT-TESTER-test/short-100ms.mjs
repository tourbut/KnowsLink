// Short related checks: 100ms stalled stream after headers, and a redirect header rejection.
import { createServer } from "node:http";
import { generateKeyPairSync } from "node:crypto";
import { pathToFileURL } from "node:url";
import { join } from "node:path";

const dist = process.argv[2];
const { TestTransport } = await import(pathToFileURL(join(dist, "test-transport.js")).href);
const pem = generateKeyPairSync("ed25519").privateKey.export({ format: "pem", type: "pkcs8" }).toString();
const server = createServer((request, response) => {
  if (request.url?.endsWith("/pull")) {
    response.writeHead(302, { Location: "http://127.0.0.1:1/leak" });
    response.end();
    return;
  }
  response.writeHead(200, { "content-type": "application/json" });
  response.write('{"partial":true}');
});
await new Promise((ok) => server.listen(0, "127.0.0.1", ok));
const base = `http://127.0.0.1:${server.address().port}`;
const transport = new TestTransport(base, "held-not-used", "trial_codex", "key1", pem, "trial_grok", {}, 100);
const started = Date.now();
const stream = await transport.request("/v1/registry").then(
  () => ({ kind: "resolved" }),
  (error) => ({ kind: "error", name: error?.name || "Error" }),
);
const streamMs = Date.now() - started;
const headerStart = Date.now();
const header = await transport.request("/v1/pull", {}).then(
  () => ({ kind: "resolved" }),
  (error) => ({ kind: "error", name: error?.name || "Error" }),
);
const headerMs = Date.now() - headerStart;
server.closeAllConnections();
server.close();
const streamPass = stream.kind === "error" && stream.name === "TimeoutError" && streamMs >= 50 && streamMs < 500;
const headerPass = header.kind === "error" && headerMs < 1000;
const pass = streamPass && headerPass && transport.timeoutMs === 100;
console.log(JSON.stringify({ storedTimeoutMs: transport.timeoutMs, stream, streamMs, header, headerMs, pass }));
process.exit(pass ? 0 : 1);
