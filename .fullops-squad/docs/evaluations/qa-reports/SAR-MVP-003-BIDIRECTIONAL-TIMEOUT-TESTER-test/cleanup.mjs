// After the default 10s body timeout: socket closes, no leftover timers, no unhandled rejection, next body succeeds.
import { createServer } from "node:http";
import { generateKeyPairSync } from "node:crypto";
import { pathToFileURL } from "node:url";
import { join } from "node:path";

const dist = process.argv[2];
const limit = Number(process.argv[3] || "10500");
const { TestTransport } = await import(pathToFileURL(join(dist, "test-transport.js")).href);
const pem = generateKeyPairSync("ed25519").privateKey.export({ format: "pem", type: "pkcs8" }).toString();
let unhandled = 0;
process.on("unhandledRejection", () => {
  unhandled += 1;
});
let closedAt = null;
const started = Date.now();
const server = createServer((request, response) => {
  response.writeHead(200, { "content-type": "application/json" });
  if (request.url?.endsWith("/registry")) {
    request.socket.once("close", () => {
      closedAt = Date.now() - started;
    });
    response.write('{"partial":true}');
    return;
  }
  response.end('{"ok":true}');
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
const waitUntil = Date.now() + 1000;
while (closedAt === null && Date.now() < waitUntil) {
  await new Promise((ok) => setTimeout(ok, 20));
}
const resources = process.getActiveResourcesInfo();
const leftoverTimers = resources.filter((name) => name === "Timeout").length;
const next = await transport.request("/v1/pull", {}).then(
  (body) => body,
  (error) => ({ kind: "error", name: error?.name || "Error" }),
);
server.closeAllConnections();
server.close();
const pass =
  outcome.kind === "error" &&
  outcome.name === "TimeoutError" &&
  ms >= 9900 &&
  ms < limit &&
  closedAt !== null &&
  closedAt < 12000 &&
  leftoverTimers === 0 &&
  unhandled === 0 &&
  next?.ok === true;
console.log(JSON.stringify({ outcome, ms, socketClosedAt: closedAt, leftoverTimers, resources, unhandled, next, pass }));
process.exit(pass ? 0 : 1);
