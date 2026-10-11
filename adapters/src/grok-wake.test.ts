// Webhook wake: config gate, metadata-only body, auth headers, result mapping and no secret in results.
import assert from "node:assert/strict";
import { chmod, mkdtemp, rm, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join } from "node:path";
import { GrokWake } from "./grok-wake.js";
import { privateDirectory } from "./private-files.js";

const root = await mkdtemp(join(tmpdir(), "knowslink-webhook-"));
try {
  const dir = join(root, "private");
  await privateDirectory(dir);
  const file = join(dir, "webhook.json");
  const key = "sk_test_DO_NOT_LEAK_0123456789";
  const url = "https://api2.cursor.sh/automations/webhook/rt_Abc-123";
  const config = (c: unknown) =>
    writeFile(file, JSON.stringify(c), { mode: 0o600 });
  const id = "0198c2a4-1b2c-7d3e-8f40-123456789abc";
  const calls: { url: string; init: RequestInit }[] = [];
  let reply: () => Promise<Response> = async () =>
    new Response("routine secret echo " + key, { status: 200 });
  const stub = (async (u: string | URL | Request, init?: RequestInit) => {
    calls.push({ url: String(u), init: init! });
    return reply();
  }) as typeof fetch;
  const wake = (ms = 8000) => new GrokWake("", 0, "", ms, file, stub);

  // Unset keeps the current wake; a webhook file replaces sendPrompt instead of adding a second ring.
  assert.equal(GrokWake.fromEnv({}), null);
  assert.equal(
    GrokWake.fromEnv({ KNOWSLINK_GROK_WEBHOOK_FILE: "webhook.json" }),
    "invalid",
  );
  const both = GrokWake.fromEnv({
    KNOWSLINK_GROK_WEBHOOK_FILE: file,
    KNOWSLINK_GROK_WAKE_AGENT: id,
  });
  assert(both instanceof GrokWake && both.agentId === "");

  await config({ url, key });
  const results: unknown[] = [];
  const ring = async (w = wake()) => {
    const r = await w.ring([id, "../not-an-id"], 3);
    results.push(r);
    return r;
  };
  assert.deepEqual(await ring(), { state: "accepted_unverified" });
  const sent = calls[0]!;
  assert.equal(sent.url, url);
  assert.equal(sent.init.redirect, "manual");
  assert.deepEqual(sent.init.headers, {
    authorization: `Bearer ${key}`,
    "x-automation-key": key,
    "content-type": "application/json",
  });
  const body = JSON.parse(String(sent.init.body)) as Record<string, unknown>;
  assert.deepEqual(Object.keys(body).sort(), [
    "event",
    "ids",
    "next",
    "pending",
  ]);
  assert.deepEqual(body.ids, [id], "only validated IDs leave the host");
  assert.equal(body.pending, 3);

  // 3xx is not followed, 4xx is final, 5xx and timeouts may have acted, refused/unresolved certainly did not arrive.
  for (const [status, state] of [
    [302, "rejected"],
    [401, "rejected"],
    [503, "uncertain"],
  ] as const) {
    reply = async () => new Response(null, { status });
    assert.deepEqual(await ring(), {
      state,
      error: `webhook HTTP ${status}`,
    });
  }
  reply = async () => {
    throw Object.assign(new TypeError("fetch failed"), {
      cause: { code: "ENOTFOUND" },
    });
  };
  assert.deepEqual(await ring(), {
    state: "retry",
    error: "webhook_unreachable",
  });
  // A silent peer: only the wake's own timeout ends it (the held timer stands in for the open socket).
  const slow = (async (_u: string | URL | Request, init?: RequestInit) =>
    new Promise<Response>((_, reject) => {
      const hold = setTimeout(() => {}, 5000);
      init!.signal!.addEventListener("abort", () => {
        clearTimeout(hold);
        reject(init!.signal!.reason);
      });
    })) as typeof fetch;
  assert.deepEqual(
    await ring(new GrokWake("", 0, "", 50, file, slow)),
    { state: "uncertain", error: "webhook_no_response" },
    "timeout is a single uncertain attempt",
  );

  // Untrusted or malformed config never sends anything.
  const before = calls.length;
  for (const c of [
    { url: "http://api2.cursor.sh/automations/webhook/rt1", key },
    { url: "https://api2.cursor.sh.evil.test/automations/webhook/rt1", key },
    { url: "https://api2.cursor.sh@evil.test/automations/webhook/rt1", key },
    { url: "https://api2.cursor.sh/automations/webhook/rt1?x=1", key },
    { url: "https://api2.cursor.sh/automations/webhook/../x", key },
    { url, key: "short" },
    { url, key: "has space inside" },
    { url },
  ]) {
    await config(c);
    assert.deepEqual(await ring(), {
      state: "retry",
      error: "webhook_config_unavailable",
    });
    assert.deepEqual(await wake().check(), {
      state: "failed",
      mode: "webhook",
      error: "webhook_config_invalid",
    });
  }
  if (process.platform !== "win32") {
    await config({ url, key });
    await chmod(file, 0o644);
    assert.equal((await ring()).error, "webhook_config_unavailable");
    await chmod(file, 0o600);
  }
  assert.equal(calls.length, before, "invalid config must not POST");
  await config({ url, key });
  assert.deepEqual(await wake().check(), { state: "ok", mode: "webhook" });
  assert.equal(calls.length, before, "check must not POST");

  const shown = JSON.stringify(results);
  assert(!shown.includes(key) && !shown.includes("rt_Abc-123"));
  process.stdout.write(
    "PASS: webhook wake config gate, single path, metadata-only body, headers, 2xx/3xx/4xx/5xx/unreachable/timeout, untrusted config, check without POST, no secret in results\n",
  );
} finally {
  await rm(root, { recursive: true, force: true });
}
