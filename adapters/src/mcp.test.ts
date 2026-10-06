// Exercise the actual SDK stdio handshake, default hold, URL guard, and redirect rejection.
import assert from "node:assert/strict";
import { createServer } from "node:http";
import { once } from "node:events";
import { Client } from "@modelcontextprotocol/sdk/client/index.js";
import { StdioClientTransport } from "@modelcontextprotocol/sdk/client/stdio.js";
import { fileURLToPath } from "node:url";
import { readFile } from "node:fs/promises";
import { createHash, generateKeyPairSync, sign } from "node:crypto";
import { Adapter, signingBytes, uuid7, type Envelope } from "./core.js";

async function failureBoundaries() {
  const manifest = await readFile(
    new URL("../../internal/relay/registry.json", import.meta.url),
    "utf8",
  );
  const { publicKey, privateKey } = generateKeyPairSync("ed25519");
  const message: Envelope = {
    v: "relay.v1",
    id: uuid7(),
    from: "synthetic_a",
    to: "synthetic_b",
    intent: "schedule.query",
    body: {
      window: { start: "2026-10-03T10:00:00Z", end: "2026-10-03T11:00:00Z" },
      granularity_min: 30,
    },
    deliver: "agent",
    exp: new Date(Date.now() + 60000).toISOString(),
    idempotency_key: "synthetic-test-key",
    sig: { alg: "Ed25519", kid: "key1", value: "" },
  };
  const signature = sign(null, signingBytes(message), privateKey).toString(
    "base64url",
  );
  const originalFetch = globalThis.fetch;
  try {
    for (const failure of [
      "signature",
      "/v1/persist",
      "/v1/ack",
      "/v1/claim",
    ]) {
      const paths: string[] = [];
      message.sig.value =
        failure === "signature"
          ? Buffer.alloc(64).toString("base64url")
          : signature;
      globalThis.fetch = async (input) => {
        const url =
          typeof input === "string"
            ? input
            : input instanceof URL
              ? input.href
              : input.url;
        const path = new URL(url).pathname;
        paths.push(path);
        if (path === failure) return new Response("{}", { status: 403 });
        let data: unknown = {};
        if (path === "/v1/registry")
          data = {
            manifest,
            sha256: createHash("sha256").update(manifest).digest("hex"),
          };
        if (path === "/v1/pull")
          data = { envelope: message, lease_token: "synthetic-lease" };
        if (path.startsWith("/v1/keys/"))
          data = {
            public: publicKey
              .export({ format: "der", type: "spki" })
              .subarray(-32)
              .toString("base64url"),
          };
        return new Response(JSON.stringify(data));
      };
      const adapter = new Adapter(
        "http://127.0.0.1",
        "synthetic-test",
        "synthetic_b",
        "key1",
        privateKey.export({ format: "pem", type: "pkcs8" }).toString(),
      );
      await assert.rejects(adapter.once(false));
      const expected = [
        "/v1/registry",
        "/v1/pull",
        "/v1/keys/synthetic_a/key1",
      ];
      if (failure !== "signature") {
        for (const path of ["/v1/persist", "/v1/ack", "/v1/claim"]) {
          expected.push(path);
          if (path === failure) break;
        }
      }
      assert.deepEqual(
        paths,
        expected,
        "verification/persist/ACK/claim failure must stop before judgment or result",
      );
    }
  } finally {
    globalThis.fetch = originalFetch;
  }
}

async function check(
  environment: Record<string, string>,
  expected: string,
  error: boolean,
) {
  const transport = new StdioClientTransport({
    command: process.execPath,
    args: [
      process.argv[2] ?? fileURLToPath(new URL("./plugin.js", import.meta.url)),
    ],
    env: environment,
    stderr: "pipe",
  });
  const client = new Client({ name: "synthetic-test", version: "1.0.0" });
  let stderr = "";
  transport.stderr?.on("data", (chunk: Buffer) => {
    stderr += chunk.toString();
  });
  try {
    await client.connect(transport);
    const tools = await client.listTools();
    assert.deepEqual(tools.tools.map((tool) => tool.name).sort(), [
      "knowslink_pull_once",
      "knowslink_status",
      "knowslink_test_receive",
      "knowslink_test_send",
      "knowslink_text_receipt",
      "knowslink_text_receive",
      "knowslink_text_send",
    ]);
    const status = await client.callTool({
      name: "knowslink_status",
      arguments: {},
    });
    assert.equal(status.isError, false);
    const result = await client.callTool({
      name: "knowslink_pull_once",
      arguments: {},
    });
    assert.equal(result.isError, error);
    assert.deepEqual(result.content, [
      {
        type: "text",
        text: JSON.stringify({
          state: expected,
          transport: "pull",
          actualConnection: "held",
          webhook: false,
          evidenceFetch: false,
        }),
      },
    ]);
    assert.equal(stderr, "");
  } finally {
    await client.close();
  }
}

async function main() {
  await failureBoundaries();
  let requests = 0;
  const listener = createServer((_request, response) => {
    requests++;
    response.writeHead(302, { Location: "/leak" });
    response.end();
  });
  listener.listen(0, "127.0.0.1");
  await once(listener, "listening");
  const address = listener.address();
  assert(address && typeof address !== "string");
  const base = `http://127.0.0.1:${address.port}`;
  const configured = {
    RELAY_URL: base,
    AGENT_CREDENTIAL: "synthetic-test",
    AGENT_ID: "synthetic_b",
    AGENT_KID: "key1",
    AGENT_KEY_FILE: "/nonexistent/synthetic-test.pem",
  };
  try {
    await check(configured, "held", true);
    await check({ ...configured, KNOWSLINK_MODE: "production" }, "held", true);
    await check({ KNOWSLINK_MODE: "synthetic-loopback" }, "unconfigured", true);
    for (const url of [
      "https://example.com",
      "file:///tmp/test",
      `${base}/path`,
      `${base}?query=1`,
      `http://user:pass@127.0.0.1:${address.port}`,
    ]) {
      await check(
        { ...configured, KNOWSLINK_MODE: "synthetic-loopback", RELAY_URL: url },
        "failed",
        true,
      );
    }
    assert.equal(
      requests,
      0,
      "held and invalid configuration must not contact relay",
    );
    const adapter = new Adapter(
      base,
      "synthetic-test",
      "synthetic_b",
      "key1",
      "unused",
    );
    await assert.rejects(adapter.request("/v1/registry"));
    assert.equal(requests, 1, "redirect must not follow even a local target");
  } finally {
    listener.close();
    await once(listener, "close");
  }
  console.info(
    "PASS: MCP initialize/discovery, held without network/key reads, invalid config, redirect rejection; no secrets in tool output",
  );
}
main().catch(() => {
  console.error("MCP synthetic boundary check failed");
  process.exitCode = 1;
});
