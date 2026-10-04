// Trial client checks held/config guards, fixed remote headers, streamed response bounds and timeout without contacting remote services.
import assert from "node:assert/strict";
import { createServer } from "node:http";
import { once } from "node:events";
import { mkdtemp, writeFile, rm } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join } from "node:path";
import { generateKeyPairSync } from "node:crypto";
import { TestTransport, testTransport } from "./test-transport.js";
async function main() {
  const originalEnv = { ...process.env };
  const originalFetch = globalThis.fetch;
  const directory = await mkdtemp(
    join(tmpdir(), "knowslink-trial-boundaries-"),
  );
  const { privateKey } = generateKeyPairSync("ed25519");
  const key = join(directory, "key.pem");
  const pem = privateKey.export({ format: "pem", type: "pkcs8" }).toString();
  await writeFile(key, pem, { mode: 0o600 });
  try {
    Object.assign(process.env, {
      KNOWSLINK_MODE: "held",
      RELAY_URL: "https://link.knowslog.com",
      AGENT_ID: "trial_codex",
      AGENT_KID: "key1",
      AGENT_CREDENTIAL: "test-only",
      AGENT_KEY_FILE: "/missing",
      KNOWSLINK_TEST_PEER: "trial_grok",
      CF_ACCESS_CLIENT_ID: "test-client",
      CF_ACCESS_CLIENT_SECRET: "test-secret",
    });
    assert.equal(await testTransport(), null);
    process.env.KNOWSLINK_MODE = "test-remote";
    for (const url of [
      "http://link.knowslog.com",
      "https://example.com",
      "https://127.0.0.1",
      "https://link.knowslog.com/path",
      "https://link.knowslog.com?x=1",
      "https://user@link.knowslog.com",
      "https://link.knowslog.com:444",
    ]) {
      process.env.RELAY_URL = url;
      await assert.rejects(testTransport(), /origin|URL/);
    }
    process.env.RELAY_URL = "https://link.knowslog.com";
    delete process.env.CF_ACCESS_CLIENT_SECRET;
    assert.equal(
      await testTransport(),
      null,
      "missing Access must not read key",
    );
    process.env.CF_ACCESS_CLIENT_SECRET = "test-secret";
    process.env.AGENT_KEY_FILE = key;
    const remote = await testTransport();
    assert(remote);
    let called = 0;
    globalThis.fetch = async (input, init) => {
      called++;
      assert.equal(input, "https://link.knowslog.com/v1/test/send");
      assert.equal(init?.redirect, "error");
      assert.equal(
        (init?.headers as Record<string, string>)["CF-Access-Client-Secret"],
        "test-secret",
      );
      assert.equal(
        (init?.headers as Record<string, string>).Authorization,
        "Bearer test-only",
      );
      const envelope = JSON.parse(init!.body as string);
      assert.equal(envelope.to, "trial_grok");
      assert.equal(envelope.intent, "relay.test.message");
      return new Response(JSON.stringify({ id: envelope.id }));
    };
    await remote.sendText("trial text", "trial-fixed-idempotent-key");
    assert.equal(
      called,
      1,
      "mocked HTTPS request only; no real remote evidence",
    );
    await assert.rejects(
      remote.sendText("x".repeat(4097), "trial-fixed-idempotent-key"),
    );
    await assert.rejects(
      remote.sendText("\ud800", "trial-fixed-idempotent-key"),
    );
    await assert.rejects(remote.sendText("ok", "short"));
    globalThis.fetch = originalFetch;
    const listener = createServer((request, response) => {
      if (request.url?.endsWith("registry")) {
        response.writeHead(200);
        response.write(" ");
        return;
      }
      if (request.url?.endsWith("send")) {
        response.end("x".repeat(65537));
        return;
      }
      response.writeHead(302, { Location: "http://127.0.0.1:1/leak" });
      response.end();
    });
    listener.listen(0, "127.0.0.1");
    await once(listener, "listening");
    const address = listener.address();
    assert(address && typeof address !== "string");
    const transport = new TestTransport(
      `http://127.0.0.1:${address.port}`,
      "test-only",
      "trial_codex",
      "key1",
      pem,
      "trial_grok",
      {},
      100,
    );
    try {
      const start = Date.now();
      await assert.rejects(transport.request("/v1/registry"));
      assert(
        Date.now() - start < 2000,
        "stream timeout must bound response body too",
      );
      await assert.rejects(transport.request("/v1/send", {}), /too large/);
      await assert.rejects(transport.request("/v1/pull", {}));
    } finally {
      listener.closeAllConnections();
      listener.close();
      await once(listener, "close");
    }
    console.info(
      "PASS: trial held, fixed HTTPS origin, Access+agent headers, text bound, redirect, streamed 64KiB bound, 100ms timeout; remote roundtrip unverified",
    );
  } finally {
    globalThis.fetch = originalFetch;
    process.env = originalEnv;
    await rm(directory, { recursive: true });
  }
}
main().catch(() => {
  console.error("Trial boundary check failed");
  process.exitCode = 1;
});
