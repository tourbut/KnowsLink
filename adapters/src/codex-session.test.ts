// Codex session mode over real stdio MCP with an elicitation-capable client standing in for the owner: bind, single receiver, takeover and per-send approval.
import assert from "node:assert/strict";
import { generateKeyPairSync, sign } from "node:crypto";
import { mkdtemp, writeFile } from "node:fs/promises";
import { createServer } from "node:http";
import { tmpdir } from "node:os";
import { join } from "node:path";
import { fileURLToPath } from "node:url";
import { Client } from "@modelcontextprotocol/sdk/client/index.js";
import { StdioClientTransport } from "@modelcontextprotocol/sdk/client/stdio.js";
import { ElicitRequestSchema } from "@modelcontextprotocol/sdk/types.js";
import { uuid7 } from "./core.js";
import { privateDirectory } from "./private-files.js";
import { textSigningBytes, type TextEnvelope } from "./text.js";

const peer = generateKeyPairSync("ed25519");
const mine = generateKeyPairSync("ed25519");
const queue: TextEnvelope[] = [];
const sent: string[] = [];
let pulls = 0;
const relay = createServer(async (req, res) => {
  let raw = "";
  for await (const chunk of req) raw += chunk;
  res.setHeader("Content-Type", "application/json");
  const url = req.url ?? "";
  if (url === "/v1/text/pull") {
    pulls++;
    const head = queue[0];
    return res.end(
      JSON.stringify(head ? { envelope: head, lease_token: "lease" } : null),
    );
  }
  if (url.startsWith("/v1/keys/"))
    return res.end(
      JSON.stringify({
        public: peer.publicKey
          .export({ format: "der", type: "spki" })
          .subarray(-32)
          .toString("base64url"),
      }),
    );
  if (url.startsWith("/v1/contacts/"))
    return res.end(
      JSON.stringify({ name: "도우미", client: "Codex", friend: "친구 A" }),
    );
  const body = JSON.parse(raw) as { id: string; text?: string };
  if (url === "/v1/text/send") {
    sent.push(body.text ?? "");
    return res.end(JSON.stringify({ id: body.id, transport: "queued" }));
  }
  if (url === "/v1/text/ack" && queue[0]?.id === body.id) queue.shift();
  res.end(JSON.stringify({ state: "ok" }));
});
await new Promise<void>((r) => relay.listen(0, "127.0.0.1", r));
const address = relay.address();
assert(address && typeof address === "object");
const base = `http://127.0.0.1:${address.port}`;
const folder = join(await mkdtemp(join(tmpdir(), "kl-codex-")), "agent");
await privateDirectory(folder);
await writeFile(
  join(folder, "agent.json"),
  JSON.stringify({
    relay: base,
    agent: "agent_b",
    kid: "k1",
    credential: "c".repeat(43),
  }),
  { mode: 0o600 },
);
await writeFile(
  join(folder, "private.pem"),
  mine.privateKey.export({ format: "pem", type: "pkcs8" }),
  { mode: 0o600 },
);
function envelope(text: string): TextEnvelope {
  const e: TextEnvelope = {
    v: "knowslink.text.v1",
    id: uuid7(),
    from: "agent_a",
    to: "agent_b",
    text,
    exp: new Date(Date.now() + 120000).toISOString(),
    idempotency_key: "codex-session-test-" + text.length,
    sig: { alg: "Ed25519", kid: "k1", value: "" },
  };
  e.sig.value = sign(null, textSigningBytes(e), peer.privateKey).toString(
    "base64url",
  );
  return e;
}

// answers: the owner's reply to each dialog in order; asked: what the dialog showed.
// accept types the shown code; blind is a full-auto host's accept without input; wrong types another code.
async function open(elicit = true) {
  const answers: ("accept" | "blind" | "wrong" | "decline")[] = [];
  const asked: string[] = [];
  const client = new Client(
    { name: "codex-session-test", version: "1.0.0" },
    { capabilities: elicit ? { elicitation: {} } : {} },
  );
  if (elicit)
    client.setRequestHandler(ElicitRequestSchema, (r) => {
      asked.push(r.params.message);
      const answer = answers.shift() ?? "cancel";
      const code = /확인 코드 (\d{4})/.exec(r.params.message)?.[1];
      if (answer === "blind") return { action: "accept" };
      if (answer === "wrong")
        return {
          action: "accept",
          content: { code: code === "1000" ? "1001" : "1000" },
        };
      return answer === "accept"
        ? { action: "accept", content: { code: ` ${code} ` } }
        : { action: answer };
    });
  await client.connect(
    new StdioClientTransport({
      command: process.execPath,
      args: [fileURLToPath(new URL("./plugin.js", import.meta.url))],
      env: {
        KNOWSLINK_MODE: "public-node",
        KNOWSLINK_AGENT_FOLDER: folder,
        KNOWSLINK_AUTO_RECEIVE: "off",
        KNOWSLINK_SESSION_HOST: "codex",
      },
    }),
  );
  const call = async (name: string, args: Record<string, unknown> = {}) => {
    const r = await client.callTool({ name, arguments: args });
    const content = r.content as { text: string }[];
    return JSON.parse(content[0]!.text) as Record<string, any>; // eslint-disable-line @typescript-eslint/no-explicit-any
  };
  return { client, answers, asked, call };
}
const send = (text: string) => ({
  peer: "agent_a",
  text,
  idempotency_key: "codex-session-send-" + text.length,
  confirmed: true,
});

const a = await open();
try {
  // 1. Installed and connected is not bound: no pull, every text tool denied, nothing asked (SP-01, SP-12).
  assert(
    (await a.client.listTools()).tools.some(
      (t) => t.name === "knowslink_session",
    ),
  );
  queue.push(envelope("before bind"));
  await new Promise((r) => setTimeout(r, 300));
  assert.equal(pulls, 0, "an unbound Codex session never pulls");
  assert.equal((await a.call("knowslink_text_receive")).state, "denied");
  assert.equal(
    (await a.call("knowslink_text_send", send("hi"))).state,
    "denied",
  );
  assert.equal(a.asked.length, 0);
  assert.equal(
    (await a.call("knowslink_session", { action: "status" })).state,
    "unbound",
  );

  // 2. Decline, a blind accept (Codex full-auto mode) and a wrong code keep it unbound; the typed code binds and starts receive.
  for (const answer of ["decline", "blind", "wrong"] as const) {
    a.answers.push(answer);
    const r = await a.call("knowslink_session", { action: "on" });
    assert.equal(r.state, "declined", answer);
  }
  assert.equal(pulls, 0);
  a.answers.push("accept");
  const bound = await a.call("knowslink_session", { action: "on" });
  assert.equal(bound.state, "bound");
  assert.match(
    a.asked[3]!,
    /이 Codex 세션\(구분값 codex-[0-9a-f]{6}\)을 내 agent agent_b/,
  );
  const got = await a.call("knowslink_text_receive");
  assert.equal(got.state, "received");
  assert.equal(got.message.text, "before bind");
  assert.equal(got.message.untrusted, true);

  // 3. confirmed:true from the model is not approval; the owner sees the whole body, invisible code points named.
  const body = "보낼 본문\u200Bend";
  a.answers.push("decline");
  assert.equal(
    (await a.call("knowslink_text_send", send(body))).state,
    "denied",
  );
  assert.equal(sent.length, 0);
  assert.match(a.asked[4]!, /친구 친구 A · agent 도우미 \(agent_a\)/);
  assert(a.asked[4]!.includes("「보낼 본문⟨U+200B⟩end」"));
  assert(a.asked[4]!.includes("보이지 않는 문자 1개"));
  a.answers.push("blind");
  assert.equal(
    (await a.call("knowslink_text_send", send(body))).state,
    "denied",
  );
  assert.equal(sent.length, 0, "a full-auto accept never sends");
  a.answers.push("accept");
  assert.equal(
    (await a.call("knowslink_text_send", send(body))).transport,
    "queued",
  );
  assert.deepEqual(sent, [body]);

  // 4. A second live session is refused without asking; take moves the address and the first session loses it (SP-02, SP-09).
  const b = await open();
  try {
    assert.equal(
      (await b.call("knowslink_session", { action: "on" })).state,
      "busy",
    );
    assert.equal(b.asked.length, 0);
    b.answers.push("accept");
    assert.equal(
      (await b.call("knowslink_session", { action: "on", take: true })).state,
      "bound",
    );
    assert.match(b.asked[0]!, /다른 세션 codex-[0-9a-f]{6} 의 수신은 끊깁니다/);
    const lost = await a.call("knowslink_text_receive");
    assert.equal(lost.state, "denied");
    assert.match(lost.text, /다른 세션이 이 agent를 가져가/);
    assert.equal(
      (await b.call("knowslink_session", { action: "off" })).state,
      "unbound",
    );
    assert.equal(
      (await b.call("knowslink_text_receipt", { id: uuid7() })).state,
      "denied",
    );
  } finally {
    await b.client.close();
  }

  // 5. A host without elicitation can never bind or send.
  const c = await open(false);
  try {
    assert.equal(
      (await c.call("knowslink_session", { action: "on" })).state,
      "declined",
    );
    assert.equal(
      (await c.call("knowslink_text_send", send("x"))).state,
      "denied",
    );
  } finally {
    await c.client.close();
  }
  assert.deepEqual(sent, [body], "only the one approved send left");
} finally {
  await a.client.close();
  relay.close();
}
process.stdout.write(
  "PASS: codex session bind/decline/single receiver/take/per-send approval/no elicitation\n",
);
