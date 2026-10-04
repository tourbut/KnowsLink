// Two independent stdio MCP agents exchange trial messages over the real isolated Postgres relay.
import assert from "node:assert/strict";
import { mkdtemp, writeFile, rm } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join } from "node:path";
import { fileURLToPath } from "node:url";
import { Client } from "@modelcontextprotocol/sdk/client/index.js";
import { StdioClientTransport } from "@modelcontextprotocol/sdk/client/stdio.js";
import { setupTrial } from "./trial-setup.js";
import { TestTransport } from "./test-transport.js";

async function main() {
  const base = process.argv[2]!;
  const { a, b } = await setupTrial(base);
  const directory = await mkdtemp(join(tmpdir(), "knowslink-trial-"));
  const clients: Client[] = [];
  try {
    for (const [own, peer] of [
      [a, b],
      [b, a],
    ]) {
      const key = join(directory, own!.agent + ".pem");
      await writeFile(key, own!.pem, { mode: 0o600 });
      const client = new Client({ name: own!.agent, version: "1.0.0" });
      clients.push(client);
      await client.connect(
        new StdioClientTransport({
          command: process.execPath,
          args: [fileURLToPath(new URL("./plugin.js", import.meta.url))],
          env: {
            KNOWSLINK_MODE: "test-loopback",
            RELAY_URL: base,
            AGENT_CREDENTIAL: own!.credential,
            AGENT_ID: own!.agent,
            AGENT_KID: "key1",
            AGENT_KEY_FILE: key,
            KNOWSLINK_TEST_PEER: peer!.agent,
          },
          stderr: "pipe",
        }),
      );
    }
    async function call(
      client: Client,
      name: string,
      args: Record<string, string> = {},
    ) {
      const result = await client.callTool({ name, arguments: args });
      assert.equal(result.isError, false);
      const content = result.content as { text: string }[];
      return JSON.parse(content[0]!.text);
    }
    const outbound = await call(clients[0]!, "knowslink_test_send", {
      text: "Codex trial ping 003",
      idempotency_key: "trial-roundtrip-codex-003",
    });
    const replay = await call(clients[0]!, "knowslink_test_send", {
      text: "Codex trial ping 003",
      idempotency_key: "trial-roundtrip-codex-003",
    });
    assert.equal(outbound.id, replay.id);
    const received = await call(clients[1]!, "knowslink_test_receive");
    assert.equal(received.message.id, outbound.id);
    assert.equal(received.message.text, "Codex trial ping 003");
    assert.equal(received.message.from, a.agent);
    assert.equal(received.message.untrusted, true);
    const reply = await call(clients[1]!, "knowslink_test_send", {
      text: `Grok trial reply to ${outbound.id}`,
      idempotency_key: "trial-roundtrip-grok-003",
    });
    const returned = await call(clients[0]!, "knowslink_test_receive");
    assert.equal(returned.message.id, reply.id);
    assert.equal(returned.message.text, `Grok trial reply to ${outbound.id}`);
    assert.equal(returned.message.from, b.agent);
    assert.equal(
      (await call(clients[0]!, "knowslink_test_receive")).state,
      "empty",
    );
    assert.equal(
      (await call(clients[1]!, "knowslink_test_receive")).state,
      "empty",
    );
    const transport = new TestTransport(
      base,
      a.credential,
      a.agent,
      "key1",
      a.pem,
      b.agent,
    );
    await assert.rejects(transport.request("/v1/owners", {}));
    const receipt = await transport.request<{ receipt: { transport: string } }>(
      `/v1/receipts/${outbound.id}`,
    );
    assert.equal(receipt.receipt.transport, "delivered");
    console.info(
      JSON.stringify({
        state: "local_two_agent_roundtrip_pass",
        codex_send: outbound.id,
        grok_received: received.message.id,
        grok_reply: reply.id,
        codex_received: returned.message.id,
        actualGrok: "unverified",
        duplicate: "same_receipt",
        second_pull: "empty",
      }),
    );
  } finally {
    await Promise.all(clients.map((client) => client.close()));
    await rm(directory, { recursive: true });
  }
}
main().catch(() => {
  console.error("Trial roundtrip failed; no credentials or payload logged");
  process.exitCode = 1;
});
