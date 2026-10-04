// Retest the fa168938 held.mjs stall: real HTTP partial body, default 6-arg TestTransport.
import { createServer } from "node:http";
import { generateKeyPairSync } from "node:crypto";
import { pathToFileURL } from "node:url";
import { join } from "node:path";

const dist = process.argv[2];
const limit = Number(process.argv[3] || "10500");
const run = Number(process.argv[4] || "1");
const { TestTransport } = await import(pathToFileURL(join(dist, "test-transport.js")).href);
const pem = generateKeyPairSync("ed25519").privateKey.export({ format: "pem", type: "pkcs8" }).toString();
const server = createServer((_request, response) => {
  response.writeHead(200, { "content-type": "application/json" });
  response.write('{"partial":true}');
});
await new Promise((ok) => server.listen(0, "127.0.0.1", ok));
const transport = new TestTransport(
  `http://127.0.0.1:${server.address().port}`,
  "held-not-used",
  "trial_codex",
  "key1",
  pem,
  "trial_grok",
);
let timer;
const started = Date.now();
const outcome = await Promise.race([
  transport.request("/v1/registry").then(
    () => ({ kind: "resolved" }),
    (error) => ({ kind: "error", name: error?.name || "Error" }),
  ),
  new Promise((resolve) => {
    timer = setTimeout(() => resolve({ kind: "hung" }), limit);
  }),
]);
const ms = Date.now() - started;
clearTimeout(timer);
server.closeAllConnections();
server.close();
const pass =
  outcome.kind === "error" &&
  outcome.name === "TimeoutError" &&
  transport.timeoutMs === 10000 &&
  TestTransport.length === 6 &&
  ms >= 9900 &&
  ms < limit;
console.log(JSON.stringify({ run, storedTimeoutMs: transport.timeoutMs, ctorLength: TestTransport.length, outcome, ms, limit, pass }));
process.exit(pass ? 0 : 1);
