// After default 10s body timeout: socket closed by client, no leftover timers, no unhandled rejection; then a normal request succeeds.
import { createServer } from "node:http";
import { generateKeyPairSync } from "node:crypto";
const { TestTransport } = await import(process.argv[2] + "/test-transport.js");
const pem = generateKeyPairSync("ed25519").privateKey.export({ format: "pem", type: "pkcs8" }).toString();
let unhandled = 0; process.on("unhandledRejection", () => { unhandled += 1; });
let closedAt = null; const start = Date.now();
const server = createServer((q, r) => {
  r.writeHead(200, { "content-type": "application/json" });
  if (q.url.endsWith("/registry")) { q.socket.once("close", () => { closedAt = Date.now() - start; }); r.write('{"partial":true}'); return; }
  r.end('{"ok":true}');
});
await new Promise((ok) => server.listen(0, "127.0.0.1", ok));
const t = new TestTransport(`http://127.0.0.1:${server.address().port}`, "c", "a", "k", pem, "p");
const out = await t.request("/v1/registry").then(() => "resolved", (e) => e.name);
const ms = Date.now() - start;
await new Promise((ok) => setImmediate(ok)); await new Promise((ok) => setTimeout(ok, 200));
const timers = process.getActiveResourcesInfo().filter((n) => n === "Timeout").length;
const next = await t.request("/v1/pull", {});
console.log(JSON.stringify({ out, ms, socketClosedAt: closedAt, leftoverTimers: timers, unhandled, next }));
server.closeAllConnections(); server.close();
