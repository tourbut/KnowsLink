// Repro: default 10000ms TestTransport against a real local server that stalls body after first chunk.
import { createServer } from "node:http";
import { generateKeyPairSync } from "node:crypto";
const dist = process.argv[2];
const { TestTransport } = await import(dist + "/test-transport.js");
const pem = generateKeyPairSync("ed25519").privateKey.export({ format: "pem", type: "pkcs8" }).toString();
const server = createServer((_q, r) => { r.writeHead(200, { "content-type": "application/json" }); r.write('{"partial":true}'); });
await new Promise((ok) => server.listen(0, "127.0.0.1", ok));
const tm = process.argv[4] ? Number(process.argv[4]) : undefined;
const t = tm ? new TestTransport(`http://127.0.0.1:${server.address().port}`, "c", "a", "k", pem, "p", {}, tm) : new TestTransport(`http://127.0.0.1:${server.address().port}`, "c", "a", "k", pem, "p");
const gc = process.argv[3] === "gc" && globalThis.gc;
if (gc) setInterval(() => gc(), 200).unref();
const start = Date.now();
const out = await Promise.race([
  t.request("/v1/registry").then(() => "resolved", (e) => e.name),
  new Promise((ok) => setTimeout(() => ok("hung"), (tm ?? 10000) + 5000)),
]);
console.log(JSON.stringify({ gc: Boolean(gc), out, ms: Date.now() - start }));
server.closeAllConnections(); server.close();
