// Real local CLI and independent MCP processes check member text against the isolated Go/Postgres service; no vendor account claims.
import assert from "node:assert/strict";
import { execFile } from "node:child_process";
import { promisify } from "node:util";
import { fileURLToPath } from "node:url";
import { Client } from "@modelcontextprotocol/sdk/client/index.js";
import { StdioClientTransport } from "@modelcontextprotocol/sdk/client/stdio.js";
import { memberTransport } from "./text.js";

const [, , folderA, folderB, agentA, agentB] = process.argv;
assert.ok(folderA && folderB && agentA && agentB);
const a = await memberTransport(folderA),
  b = await memberTransport(folderB);
await assert.rejects(
  a.sendText(agentB, "ping", "public-explicit-key-01", false),
);
await assert.rejects(
  a.sendText(agentB, "한".repeat(1366), "public-too-long-key", true),
);
const clients: Client[] = [];
try {
  for (const folder of [folderA, folderB]) {
    const c = new Client({ name: "public-local-check", version: "1.0.0" });
    clients.push(c);
    await c.connect(
      new StdioClientTransport({
        command: process.execPath,
        args: [fileURLToPath(new URL("./plugin.js", import.meta.url))],
        env: { KNOWSLINK_MODE: "public-node", KNOWSLINK_AGENT_FOLDER: folder },
        stderr: "pipe",
      }),
    );
  }
  async function call(
    index: number,
    name: string,
    args: Record<string, unknown> = {},
  ) {
    const r = await clients[index]!.callTool({ name, arguments: args });
    assert.equal(r.isError, false);
    const content = r.content as { text: string }[];
    return JSON.parse(content[0]!.text);
  }
  const sent = await call(0, "knowslink_text_send", {
    peer: agentB,
    text: "Public connection check",
    idempotency_key: "public-mcp-outbound-01",
    confirmed: true,
  });
  assert.equal(sent.transport, "queued");
  const incoming = await call(1, "knowslink_text_receive");
  assert.equal(incoming.message.id, sent.id);
  assert.equal(incoming.message.from, agentA);
  assert.equal(incoming.message.untrusted, true);
  assert.equal(incoming.message.text, "Public connection check");
  const replay = await call(0, "knowslink_text_send", {
    peer: agentB,
    text: "Public connection check",
    idempotency_key: "public-mcp-outbound-01",
    confirmed: true,
  });
  assert.equal(replay.id, sent.id);
  assert.equal((await call(1, "knowslink_text_receive")).state, "empty");
  const reply = await call(1, "knowslink_text_send", {
    peer: agentA,
    text: "Public related reply",
    reply_to: sent.id,
    idempotency_key: "public-mcp-reply-0001",
    confirmed: true,
  });
  const returned = await call(0, "knowslink_text_receive");
  assert.equal(returned.message.reply_to, sent.id);
  assert.equal(returned.message.id, reply.id);
  assert.equal(returned.message.from, agentB);
  assert.equal(returned.message.text, "Public related reply");
  const receipt = await call(0, "knowslink_text_receipt", { id: sent.id });
  assert.equal(receipt.reply_id, reply.id);
  assert.equal(receipt.completion, "reply_received");
  const conflict = await clients[0]!.callTool({
    name: "knowslink_text_send",
    arguments: {
      peer: agentB,
      text: "different",
      idempotency_key: "public-mcp-outbound-01",
      confirmed: true,
    },
  });
  assert.equal(conflict.isError, true);
  assert.match(JSON.stringify(conflict.content), /idempotency_conflict/);
  // Use actual CLI processes too; stdin carries text, argv contains no credential or private key.
  const run = promisify(execFile);
  const textCLI = fileURLToPath(new URL("./text.js", import.meta.url));
  const receivedCLI = await run(process.execPath, [
    textCLI,
    "receive",
    folderA,
  ]);
  assert.equal(JSON.parse(receivedCLI.stdout), null);
  const cliReceipt = await run(process.execPath, [
    textCLI,
    "receipt",
    folderA,
    sent.id,
  ]);
  assert.equal(JSON.parse(cliReceipt.stdout).reply_id, reply.id);
  console.info(
    JSON.stringify({
      state: "public_node_process_roundtrip_pass",
      request: sent.id,
      received: incoming.message.id,
      reply: reply.id,
      replyReceived: returned.message.id,
      actualGrok: "unverified",
      actualDadot: "unverified",
    }),
  );
} finally {
  await Promise.all(clients.map((c) => c.close()));
}
// Keep both loaded clients used and verify no queued data remains after the roundtrip.
assert.equal(await b.receive(), null);
