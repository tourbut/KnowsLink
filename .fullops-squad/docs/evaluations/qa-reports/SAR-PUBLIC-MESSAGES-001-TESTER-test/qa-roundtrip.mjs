// Independent local Node CLI and MCP processes for one member's two agents.
// Prints IDs only. Does not claim a vendor account, real mail, or production PS-08.
import assert from "node:assert/strict";
import { spawn } from "node:child_process";
import { join } from "node:path";
import { pathToFileURL } from "node:url";

const clone = process.env.QA_CLONE;
const [folderA, folderB, agentA, agentB] = process.argv.slice(2);
assert.ok(clone && folderA && folderB && agentA && agentB);
// execFile's input option does not finish this CLI. Close stdin from spawn.
function run(command, args, options = {}) {
  return new Promise((resolve, reject) => {
    const child = spawn(command, args, { stdio: ["pipe", "pipe", "pipe"] });
    let stdout = "";
    let stderr = "";
    child.stdout.on("data", (chunk) => {
      stdout += chunk;
    });
    child.stderr.on("data", (chunk) => {
      stderr += chunk;
    });
    child.on("error", reject);
    child.on("close", (code) => {
      if (code) {
        const error = new Error("command failed");
        error.code = code;
        error.stdout = stdout;
        error.stderr = stderr;
        reject(error);
        return;
      }
      resolve({ stdout, stderr });
    });
    child.stdin.end(options.input ?? "");
  });
}
const textCLI = join(clone, "adapters/dist/text.js");
const plugin = join(clone, "adapters/dist/plugin.js");
const { memberTransport } = await import(pathToFileURL(join(clone, "adapters/dist/text.js")).href);
const { Client } = await import(
  pathToFileURL(join(clone, "adapters/node_modules/@modelcontextprotocol/sdk/dist/esm/client/index.js")).href
);
const { StdioClientTransport } = await import(
  pathToFileURL(join(clone, "adapters/node_modules/@modelcontextprotocol/sdk/dist/esm/client/stdio.js")).href
);

const transportA = await memberTransport(folderA);
await assert.rejects(transportA.sendText(agentB, "unconfirmed", "qa-unconfirmed-key", false));

const tooBig = await run(process.execPath, [textCLI, "send", folderA, agentB, "qa-cli-too-big-key", "--confirmed"], {
  input: "x".repeat(4097),
}).catch((error) => error);
assert.equal(tooBig.code, 1);
assert.equal(tooBig.stdout ?? "", "");
assert.doesNotMatch(String(tooBig.stderr ?? ""), /BEGIN |credential|PRIVATE/i);

const held = new Client({ name: "qa-held", version: "1.0.0" });
await held.connect(new StdioClientTransport({ command: process.execPath, args: [plugin], stderr: "pipe" }));
const heldResult = await held.callTool({ name: "knowslink_status", arguments: {} });
const heldBody = JSON.parse(heldResult.content[0].text);
assert.equal(heldBody.state, "held");
await held.close();

const clients = [];
try {
  for (const folder of [folderA, folderB]) {
    const client = new Client({ name: "qa-public-node", version: "1.0.0" });
    clients.push(client);
    await client.connect(
      new StdioClientTransport({
        command: process.execPath,
        args: [plugin],
        env: { KNOWSLINK_MODE: "public-node", KNOWSLINK_AGENT_FOLDER: folder },
        stderr: "pipe",
      }),
    );
  }
  async function call(index, name, args = {}) {
    const result = await clients[index].callTool({ name, arguments: args });
    assert.equal(result.isError, false, `${name} failed`);
    return JSON.parse(result.content[0].text);
  }
  const sent = await call(0, "knowslink_text_send", {
    peer: agentB,
    text: "QA connection check",
    idempotency_key: "qa-mcp-outbound-01",
    confirmed: true,
  });
  assert.equal(sent.transport, "queued");
  const incoming = await call(1, "knowslink_text_receive");
  assert.equal(incoming.state, "received");
  assert.equal(incoming.message.id, sent.id);
  assert.equal(incoming.message.from, agentA);
  assert.equal(incoming.message.to, agentB);
  assert.equal(incoming.message.untrusted, true);
  assert.equal(incoming.message.text, "QA connection check");
  const replay = await call(0, "knowslink_text_send", {
    peer: agentB,
    text: "QA connection check",
    idempotency_key: "qa-mcp-outbound-01",
    confirmed: true,
  });
  assert.equal(replay.id, sent.id);
  assert.equal((await call(1, "knowslink_text_receive")).state, "empty");
  const reply = await call(1, "knowslink_text_send", {
    peer: agentA,
    text: "QA related reply",
    reply_to: sent.id,
    idempotency_key: "qa-mcp-reply-00001",
    confirmed: true,
  });
  const returned = await call(0, "knowslink_text_receive");
  assert.equal(returned.message.id, reply.id);
  assert.equal(returned.message.reply_to, sent.id);
  assert.equal(returned.message.from, agentB);
  assert.equal(returned.message.untrusted, true);
  const receipt = await call(0, "knowslink_text_receipt", { id: sent.id });
  assert.equal(receipt.reply_id, reply.id);
  assert.equal(receipt.completion, "reply_received");
  assert.equal(Object.hasOwn(receipt, "text"), false);
  const conflict = await clients[0].callTool({
    name: "knowslink_text_send",
    arguments: {
      peer: agentB,
      text: "different body",
      idempotency_key: "qa-mcp-outbound-01",
      confirmed: true,
    },
  });
  assert.equal(conflict.isError, true);
  assert.match(JSON.stringify(conflict.content), /idempotency_conflict/);
  const exact = await run(process.execPath, [textCLI, "send", folderA, agentB, "qa-cli-exact-4096", "--confirmed"], {
    input: "a".repeat(4096),
  });
  const exactBody = JSON.parse(exact.stdout);
  assert.equal(exactBody.transport, "queued");
  const exactReplay = await run(process.execPath, [textCLI, "send", folderA, agentB, "qa-cli-exact-4096", "--confirmed"], {
    input: "a".repeat(4096),
  });
  assert.equal(JSON.parse(exactReplay.stdout).id, exactBody.id);
  const receivedExact = await run(process.execPath, [textCLI, "receive", folderB]);
  const exactMessage = JSON.parse(receivedExact.stdout);
  assert.equal(exactMessage.id, exactBody.id);
  assert.equal(Buffer.byteLength(exactMessage.text), 4096);
  assert.equal(exactMessage.untrusted, true);
  const cliReceipt = JSON.parse((await run(process.execPath, [textCLI, "receipt", folderA, sent.id])).stdout);
  assert.equal(cliReceipt.reply_id, reply.id);
  assert.equal(Object.hasOwn(cliReceipt, "text"), false);
  process.stdout.write(
    JSON.stringify({
      state: "qa_local_node_mcp_pass",
      request: sent.id,
      received: incoming.message.id,
      reply: reply.id,
      replyReceived: returned.message.id,
      exact4096: exactBody.id,
      held: "default_mcp",
      actualGrok: "unverified",
      actualDadot: "unverified",
      realMail: "unverified",
    }) + "\n",
  );
} finally {
  await Promise.all(clients.map((client) => client.close()));
}
