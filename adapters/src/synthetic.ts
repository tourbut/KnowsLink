// Synthetic QA seed and end-to-end adapter check; credentials stay in ignored local files only.
import assert from "node:assert/strict";
import { generateKeyPairSync, randomBytes, sign } from "node:crypto";
import { mkdir, mkdtemp, rm, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join } from "node:path";
import { fileURLToPath } from "node:url";
import { Client } from "@modelcontextprotocol/sdk/client/index.js";
import { StdioClientTransport } from "@modelcontextprotocol/sdk/client/stdio.js";
import { Adapter, signingBytes, uuid7, type Envelope } from "./index.js";

const base = process.argv[2] ?? "http://127.0.0.1:8080";
const seed = process.argv.includes("--seed");
assert(["127.0.0.1", "localhost"].includes(new URL(base).hostname));
async function api<T>(
  path: string,
  credential = "",
  body?: unknown,
  claim?: string,
): Promise<T> {
  const response = await fetch(`${base}${path}`, {
    method: body === undefined ? "GET" : "POST",
    headers: {
      Authorization: `Bearer ${credential}`,
      "Content-Type": "application/json",
      ...(claim ? { "X-Execution-Claim": claim } : {}),
    },
    body: body === undefined ? undefined : JSON.stringify(body),
  });
  assert.equal(
    response.status,
    200,
    `${path}: ${await response.clone().text()}`,
  );
  return (await response.json()) as T;
}
async function createAgent(suffix: string) {
  const owner = await api<{ owner: string; credential: string }>(
    "/v1/owners",
    "",
    {},
  );
  const agent = `synthetic_${suffix}_${randomBytes(5).toString("hex")}`;
  const { publicKey, privateKey } = generateKeyPairSync("ed25519");
  const publicBytes = publicKey
    .export({ format: "der", type: "spki" })
    .subarray(-32)
    .toString("base64url");
  const proof = sign(
    null,
    Buffer.from(
      ["KNOWSLINK-KEY-POP", owner.owner, agent, "key1", publicBytes].join("\0"),
    ),
    privateKey,
  ).toString("base64url");
  const registered = await api<{ credential: string }>(
    "/v1/agents",
    owner.credential,
    { agent, kid: "key1", public: publicBytes, proof },
  );
  const pem = privateKey.export({ type: "pkcs8", format: "pem" }).toString();
  return { agent, owner, credential: registered.credential, pem };
}
async function main() {
  const a = await createAgent("a");
  const b = await createAgent("b");
  await api("/v1/invites", a.credential, { agent: a.agent, target: b.agent });
  await api("/v1/invite-decision", b.owner.credential, {
    agent: a.agent,
    target: b.agent,
    decision: "accept",
  });
  const adapter = new Adapter(base, b.credential, b.agent, "key1", b.pem);
  async function sendQuery() {
    const message: Envelope = {
      v: "relay.v1",
      id: uuid7(),
      from: a.agent,
      to: b.agent,
      intent: "schedule.query",
      body: {
        window: { start: "2026-10-03T10:00:00Z", end: "2026-10-03T11:00:00Z" },
        granularity_min: 30,
      },
      deliver: "agent",
      exp: new Date(Date.now() + 180000)
        .toISOString()
        .replace(/\.\d{3}Z$/, "Z"),
      idempotency_key: randomBytes(24).toString("base64url"),
      sig: { alg: "Ed25519", kid: "key1", value: "" },
    };
    message.sig.value = sign(null, signingBytes(message), a.pem).toString(
      "base64url",
    );
    await api("/v1/send", a.credential, message);
    return message;
  }
  const first = await sendQuery();
  assert.equal(await adapter.once(false), true);
  const receipt = await api<{
    receipt: { transport: string };
    completion: string;
  }>(`/v1/receipts/${first.id}`, a.credential);
  assert.equal(receipt.receipt.transport, "delivered");
  assert.equal(receipt.completion, "denied");
  const source = new Adapter(base, a.credential, a.agent, "key1", a.pem);
  assert.equal(await source.once(false), true);
  assert.equal(await source.once(false), false);
  const second = await sendQuery();
  let pending: Promise<boolean>;
  let mcp: Client | undefined;
  if (seed) {
    const lease = await api<{ lease_token: string }>(
      "/v1/pull",
      b.credential,
      {},
    );
    await api("/v1/persist", b.credential, {
      id: second.id,
      token: lease.lease_token,
    });
    await api("/v1/ack", b.credential, {
      id: second.id,
      token: lease.lease_token,
    });
    const claim = await api<{ claim: string }>("/v1/claim", b.credential, {
      id: second.id,
    });
    const receipt = await api<{ receipt: { digest: string } }>(
      `/v1/receipts/${second.id}`,
      b.credential,
    );
    await adapter.send(
      second,
      "relay.approval.request",
      { reason: "judgment_required", request_digest: receipt.receipt.digest },
      claim.claim,
    );
    await mkdir("build", { recursive: true });
    await writeFile(
      "build/qa-fixture.json",
      JSON.stringify(
        { base, a, b, parent: second.id, claim: claim.claim },
        null,
        2,
      ),
      { mode: 0o600 },
    );
    console.info(
      "QA credentials: build/qa-fixture.json (0600, ignored); owner UI: /owner",
    );
    pending = Promise.resolve(true);
  } else {
    const keyDirectory = await mkdtemp(join(tmpdir(), "knowslink-mcp-test-"));
    const keyFile = join(keyDirectory, "key.pem");
    await writeFile(keyFile, b.pem, { mode: 0o600 });
    mcp = new Client({ name: "synthetic-relay", version: "1.0.0" });
    await mcp.connect(
      new StdioClientTransport({
        command: process.execPath,
        args: [fileURLToPath(new URL("./plugin.js", import.meta.url))],
        env: {
          KNOWSLINK_MODE: "synthetic-loopback",
          RELAY_URL: base,
          AGENT_CREDENTIAL: b.credential,
          AGENT_ID: b.agent,
          AGENT_KID: "key1",
          AGENT_KEY_FILE: keyFile,
        },
        stderr: "pipe",
      }),
    );
    pending = mcp
      .callTool({ name: "knowslink_pull_once", arguments: {} }, undefined, {
        timeout: 240000,
      })
      .then((result) => {
        assert.equal(result.isError, false);
        assert.deepEqual(result.content, [
          {
            type: "text",
            text: JSON.stringify({
              state: "processed",
              transport: "pull",
              actualConnection: "held",
              webhook: false,
              evidenceFetch: false,
            }),
          },
        ]);
        return true;
      })
      .finally(async () => {
        await mcp?.close();
        await rm(keyDirectory, { recursive: true });
      });
  }
  let gate = "";
  for (let attempt = 0; attempt < 40 && !gate; attempt++) {
    const response = await fetch(`${base}/owner`, {
      headers: { Authorization: `Bearer ${b.owner.credential}` },
    });
    gate =
      /href="\/owner\/gates\/([^"]+)"/.exec(await response.text())?.[1] ?? "";
    if (!gate) await new Promise((resolve) => setTimeout(resolve, 100));
  }
  assert(gate, "pending gate absent");
  const response = await fetch(`${base}/owner/gates/${gate}`, {
    headers: { Authorization: `Bearer ${b.owner.credential}` },
  });
  const html = await response.text();
  assert(
    html.includes("granularity_min") &&
      html.includes("pending") &&
      html.includes("disclosure policy absent"),
  );
  const csrf = /name="csrf" value="([^"]+)"/.exec(html)?.[1];
  assert(csrf);
  if (!seed) {
    const decided = await fetch(`${base}/owner/gates/${gate}`, {
      method: "POST",
      redirect: "manual",
      headers: {
        Authorization: `Bearer ${b.owner.credential}`,
        "Content-Type": "application/x-www-form-urlencoded",
      },
      body: new URLSearchParams({ csrf, decision: "approve" }),
    });
    assert.equal(decided.status, 303);
    assert.equal(await pending, true);
    const result = await api<{ completion: string }>(
      `/v1/receipts/${second.id}`,
      a.credential,
    );
    assert.equal(result.completion, "denied");
    console.info(
      "PASS: TS verify/persist/ACK/claim, policy deny, owner approve and denied result; no effects",
    );
  } else {
    await pending;
    console.info(
      `UI candidate: ${base}/owner/gates/${gate}; expires ${second.exp}`,
    );
  }
}
main().catch((error: unknown) => {
  console.error(
    error instanceof Error ? error.message : "synthetic check failed",
  );
  process.exitCode = 1;
});
