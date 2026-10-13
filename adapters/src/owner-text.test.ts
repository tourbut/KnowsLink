// The authenticated relay can deliver explicitly owner-written text without inventing an agent signature.
import assert from "node:assert/strict";
import { createServer } from "node:http";
import { uuid7 } from "./core.js";
import { TextTransport, type OwnerTextEnvelope } from "./text.js";

let envelope: Record<string, unknown> = {};
const operations: string[] = [];
let stored = false;
const server = createServer(async (request, response) => {
  assert.equal(request.headers.authorization, `Bearer ${"c".repeat(43)}`);
  for await (const chunk of request) void chunk;
  response.setHeader("Content-Type", "application/json");
  const path = request.url ?? "";
  if (path === "/v1/text/pull") {
    response.end(JSON.stringify({ envelope, lease_token: "synthetic-lease" }));
    return;
  }
  operations.push(path);
  assert(stored, "local copy must precede persist/ACK");
  response.end("{}");
});
await new Promise<void>((resolve) => server.listen(0, "127.0.0.1", resolve));
const address = server.address();
assert(address && typeof address === "object");
const transport = new TextTransport(
  `http://127.0.0.1:${address.port}`,
  "c".repeat(43),
  "agent_b",
  "key1",
  "",
);
const valid: OwnerTextEnvelope = {
  v: "knowslink.owner-text.v1",
  id: uuid7(),
  from: "agent_a",
  to: "agent_b",
  text: "synthetic owner text",
  exp: new Date(Date.now() + 60000).toISOString(),
  idempotency_key: "web:synthetic-owner-attempt",
  authorship: "owner",
};
try {
  envelope = { ...valid };
  const result = await transport.receive(async (message) => {
    assert.equal(message.authorship, "owner");
    assert.equal(message.from, "agent_a");
    assert.equal(message.untrusted, true);
    stored = true;
  });
  assert.equal(result?.text, valid.text);
  assert.deepEqual(operations, ["/v1/text/persist", "/v1/text/ack"]);
  for (const change of [
    { authorship: "agent" },
    { sig: { alg: "Ed25519", kid: "key1", value: "fake" } },
    { to: "agent_c" },
    { from: "invalid sender" },
    { text: " " },
    { text: "한".repeat(1366) },
    { text: "\ud800" },
    { text: "\udfff" },
    { exp: new Date(0).toISOString() },
    { v: "knowslink.text.v1" },
    { id: "not-a-message" },
  ]) {
    envelope = { ...valid, ...change };
    stored = false;
    await assert.rejects(() => transport.receive(), /invalid text envelope/);
    assert.equal(
      operations.length,
      2,
      "invalid envelopes must never persist/ACK",
    );
  }
  console.log(
    "PASS: owner authorship, recipient binding, schema refusal and durable copy before ACK",
  );
} finally {
  await new Promise<void>((resolve, reject) =>
    server.close((error) => (error ? reject(error) : resolve())),
  );
}
