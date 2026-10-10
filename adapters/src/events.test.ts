// Real HTTP arrival streams exercise push latency, reconnect catch-up, coalescing, serialized pull, fallback and cancellation.
import assert from "node:assert/strict";
import { generateKeyPairSync, sign } from "node:crypto";
import { mkdtemp, readdir, rm } from "node:fs/promises";
import { createServer, type ServerResponse } from "node:http";
import { tmpdir } from "node:os";
import { join } from "node:path";
import { AutoReceiver, Inbox } from "./inbox.js";
import { uuid7 } from "./core.js";
import { TextTransport, textSigningBytes, type TextEnvelope } from "./text.js";

const peer = generateKeyPairSync("ed25519");
const queue: TextEnvelope[] = [];
const streams = new Set<ServerResponse>();
let connections = 0;
let pulls = 0;
let active = 0;
let maxActive = 0;
let unavailable = false;
let malformed = false;
const relay = createServer(async (req, res) => {
  assert.equal(req.headers.authorization, `Bearer ${"c".repeat(43)}`);
  if (req.url === "/v1/text/events") {
    connections++;
    if (unavailable) {
      res.writeHead(503);
      res.end();
      return;
    }
    res.writeHead(200, { "Content-Type": "text/event-stream" });
    if (malformed) {
      res.end('event: arrival\ndata: {"text":"bad"}\n\n');
      return;
    }
    streams.add(res);
    res.on("close", () => streams.delete(res));
    res.write("event: rea");
    res.write("dy\ndata: {}\n\n");
    return;
  }
  let raw = "";
  for await (const part of req) raw += part;
  res.setHeader("Content-Type", "application/json");
  if (req.url === "/v1/text/pull") {
    pulls++;
    active++;
    maxActive = Math.max(maxActive, active);
    await new Promise((r) => setTimeout(r, 25));
    active--;
    res.end(
      JSON.stringify(
        queue[0] ? { envelope: queue[0], lease_token: "test-lease" } : null,
      ),
    );
  } else if (req.url?.startsWith("/v1/keys/")) {
    res.end(
      JSON.stringify({
        public: peer.publicKey
          .export({ format: "der", type: "spki" })
          .subarray(-32)
          .toString("base64url"),
      }),
    );
  } else {
    if (req.url === "/v1/text/ack") {
      const { id } = JSON.parse(raw) as { id: string };
      if (queue[0]?.id === id) queue.shift();
    }
    res.end("{}");
  }
});
await new Promise<void>((r) => relay.listen(0, "127.0.0.1", r));
const address = relay.address();
assert(address && typeof address === "object");
const transport = new TextTransport(
  `http://127.0.0.1:${address.port}`,
  "c".repeat(43),
  "agent_b",
  "key1",
  "",
);
const root = await mkdtemp(join(tmpdir(), "kl-events-"));
const inbox = await Inbox.open(root);
const received: string[] = [];
const auto = new AutoReceiver(
  transport,
  inbox,
  async (m) => {
    received.push(m.id);
  },
  null,
  10000,
  "sse",
);
function enqueue(): string {
  const e: TextEnvelope = {
    v: "knowslink.text.v1",
    id: uuid7(),
    from: "agent_a",
    to: "agent_b",
    text: "untrusted test body",
    exp: new Date(Date.now() + 60000).toISOString(),
    idempotency_key: uuid7(),
    sig: { alg: "Ed25519", kid: "key1", value: "" },
  };
  e.sig.value = sign(null, textSigningBytes(e), peer.privateKey).toString(
    "base64url",
  );
  queue.push(e);
  return e.id;
}
async function wait(ok: () => boolean, ms = 2500): Promise<void> {
  const end = Date.now() + ms;
  while (!ok()) {
    assert(Date.now() < end, "arrival test timed out");
    await new Promise((r) => setTimeout(r, 10));
  }
}
try {
  const initial = enqueue();
  auto.start();
  auto.start();
  await wait(
    () => received.includes(initial) && active === 0 && streams.size === 1,
  );
  assert.equal(auto.status.signalState, "connected");
  const pushed = enqueue();
  const count = pulls;
  for (const stream of streams)
    stream.write("event: arrival\ndata: {}\n\n".repeat(100));
  await Promise.all([
    auto.serial(() => transport.receive((m) => inbox.store(m))),
    wait(() => queue.length === 0),
  ]);
  await wait(() => active === 0);
  assert((await readdir(join(root, "inbox"))).includes(`${pushed}.json`));
  assert(pulls - count < 6, "duplicate doorbells must coalesce");
  assert.equal(maxActive, 1);
  for (const stream of streams) stream.end();
  const caughtUp = enqueue();
  await wait(() => connections >= 2 && received.includes(caughtUp));
  auto.stop();
  await wait(() => streams.size === 0);
  const stoppedConnections = connections;
  await new Promise((r) => setTimeout(r, 1100));
  assert.equal(connections, stoppedConnections, "stop must cancel reconnect");

  unavailable = true;
  const fallback = new AutoReceiver(
    transport,
    inbox,
    async (m) => {
      received.push(m.id);
    },
    null,
    100,
    "sse",
  );
  fallback.start();
  await wait(() => fallback.status.signalState === "backoff");
  const repaired = enqueue();
  await wait(() => received.includes(repaired));
  fallback.stop();
  unavailable = false;
  malformed = true;
  await assert.rejects(
    transport.listen(
      () => assert.fail("payload event accepted"),
      new AbortController().signal,
    ),
    /invalid arrival signal/,
  );
  assert.equal(auto.status.receiveMode, "sse");
  assert(!JSON.stringify(auto.status).includes("c".repeat(43)));
  process.stdout.write(
    "PASS: SSE initial/push/reconnect/coalesce/serial/fallback/cancel and strict metadata frames\n",
  );
} finally {
  auto.stop();
  relay.closeAllConnections();
  await new Promise<void>((r) => relay.close(() => r()));
  await rm(root, { recursive: true, force: true });
}
