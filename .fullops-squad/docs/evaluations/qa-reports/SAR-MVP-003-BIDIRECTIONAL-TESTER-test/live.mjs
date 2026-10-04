// Two independent MCP processes, Codex stdin, and the Grok command launcher over real Postgres.
import { spawn, spawnSync } from "node:child_process";
import { mkdtemp, readFile, rm, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join } from "node:path";
import { pathToFileURL } from "node:url";
import { baseEnv, call, check, clone, Mcp, post, relay, signed } from "./session.mjs";

export async function live(plugin) {
  const directory = await mkdtemp(join(tmpdir(), "knowslink-qa-live-"));
  const secrets = [];
  const clients = [];
  try {
    const { setupTrial } = await import(pathToFileURL(join(clone, "adapters/dist/trial-setup.js")).href);
    const core = await import(pathToFileURL(join(clone, "adapters/dist/core.js")).href);
    const { a, b } = await setupTrial(relay);
    secrets.push(a.credential, b.credential, a.owner.credential, b.owner.credential, a.pem, b.pem);
    async function open(own, peer) {
      const key = join(directory, `${own.agent}.pem`);
      await writeFile(key, own.pem, { mode: 0o600 });
      const mcp = new Mcp(plugin, baseEnv({
        KNOWSLINK_MODE: "test-loopback",
        RELAY_URL: relay,
        AGENT_CREDENTIAL: own.credential,
        AGENT_ID: own.agent,
        AGENT_KID: "key1",
        AGENT_KEY_FILE: key,
        KNOWSLINK_TEST_PEER: peer.agent,
      }));
      clients.push(mcp);
      await mcp.init();
      return { mcp, key };
    }
    const left = await open(a, b);
    const right = await open(b, a);
    check("two mcp pids differ", left.mcp.child.pid !== right.mcp.child.pid, `${left.mcp.child.pid}/${right.mcp.child.pid}`);
    const status = await call(left.mcp, "knowslink_status");
    check("trial status is unverified", status.body.state === "trial_configured_unverified" && status.body.actualConnection === "held");
    const text = "Codex independent ping";
    const outbound = await call(left.mcp, "knowslink_test_send", { text, idempotency_key: "qa-mcp-codex-round-1" });
    const replay = await call(left.mcp, "knowslink_test_send", { text, idempotency_key: "qa-mcp-codex-round-1" });
    check("duplicate idempotency returns same id", !outbound.isError && outbound.body.state === "queued" && outbound.body.id === replay.body.id, outbound.body.state || "");
    const conflict = await call(left.mcp, "knowslink_test_send", { text: "Codex independent changed", idempotency_key: "qa-mcp-codex-round-1" });
    check("idempotency conflict fails closed", conflict.isError && conflict.body.state === "failed");
    const untouched = await call(left.mcp, "knowslink_pull_once");
    check("trial pull_once stays held", untouched.isError && untouched.body.state === "held");
    const early = await call(left.mcp, "knowslink_test_receive");
    check("sender does not receive own message", early.body.state === "empty");
    const received = await call(right.mcp, "knowslink_test_receive");
    check("receive id from to text", received.body.state === "received" && received.body.message.id === outbound.body.id && received.body.message.from === "trial_codex" && received.body.message.to === "trial_grok" && received.body.message.text === text && received.body.message.untrusted === true);
    const replyText = `Grok independent reply to ${outbound.body.id}`;
    const reply = await call(right.mcp, "knowslink_test_send", { text: replyText, idempotency_key: "qa-mcp-grok-round-1" });
    const returned = await call(left.mcp, "knowslink_test_receive");
    check("reply id from to text", returned.body.state === "received" && returned.body.message.id === reply.body.id && returned.body.message.from === "trial_grok" && returned.body.message.to === "trial_codex" && returned.body.message.text === replyText && returned.body.message.untrusted === true);
    check("queues empty after claim", (await call(left.mcp, "knowslink_test_receive")).body.state === "empty" && (await call(right.mcp, "knowslink_test_receive")).body.state === "empty");
    for (const mcp of clients) await mcp.close();
    clients.length = 0;
    async function configFor(own, peer) {
      const key = join(directory, `${own.agent}.pem`);
      const path = join(directory, `${own.agent}.json`);
      await writeFile(path, JSON.stringify({
        KNOWSLINK_MODE: "test-loopback",
        RELAY_URL: relay,
        AGENT_ID: own.agent,
        AGENT_CREDENTIAL: own.credential,
        AGENT_KID: "key1",
        AGENT_KEY_FILE: key,
        KNOWSLINK_TEST_PEER: peer.agent,
      }), { mode: 0o600 });
      return path;
    }
    const codexConfig = await configFor(a, b);
    const grokConfig = await configFor(b, a);
    const cliText = "Codex cli ping";
    const sent = spawnSync("python3", [join(clone, "scripts/run_trial.py"), "--config", codexConfig, "send", "qa-cli-codex-round-1"], { input: cliText, encoding: "utf8", cwd: "/tmp", env: baseEnv() });
    let sentBody = {};
    if (sent.status === 0) {
      try {
        sentBody = JSON.parse(sent.stdout);
      } catch {
        sentBody = {};
      }
    }
    check("codex stdin send", sent.status === 0 && typeof sentBody.id === "string", String(sent.status));
    const launched = spawn("python3", [join(clone, "scripts/run_trial.py"), "--config", grokConfig, "--node", process.execPath, "--bundle", plugin], { cwd: "/tmp", env: baseEnv(), stdio: ["pipe", "pipe", "pipe"] });
    const grokMcp = Mcp.fromChild(launched);
    clients.push(grokMcp);
    await grokMcp.init();
    const cliReceived = await call(grokMcp, "knowslink_test_receive");
    check("grok command receive matches codex send", cliReceived.body.message?.id === sentBody.id && cliReceived.body.message?.from === "trial_codex" && cliReceived.body.message?.to === "trial_grok" && cliReceived.body.message?.text === cliText);
    const cliReplyText = `Grok cli reply to ${sentBody.id}`;
    const cliReply = await call(grokMcp, "knowslink_test_send", { text: cliReplyText, idempotency_key: "qa-cli-grok-round-1" });
    await grokMcp.close();
    const got = spawnSync("python3", [join(clone, "scripts/run_trial.py"), "--config", codexConfig, "receive"], { encoding: "utf8", cwd: "/tmp", env: baseEnv() });
    let gotBody = {};
    if (got.status === 0) {
      try {
        gotBody = JSON.parse(got.stdout);
      } catch {
        gotBody = {};
      }
    }
    check("codex receive matches grok reply", got.status === 0 && gotBody.message?.id === cliReply.body.id && gotBody.message?.from === "trial_grok" && gotBody.message?.to === "trial_codex" && gotBody.message?.text === cliReplyText, String(got.status));
    const replayHttp = await post("/v1/test/send", a.credential, signed(a, b.agent, { core, text, key: "qa-mcp-codex-round-1", expMs: 60000 }), secrets);
    check("http duplicate same id", replayHttp.status === 200 && replayHttp.id === outbound.body.id, String(replayHttp.status));
    const conflictHttp = await post("/v1/test/send", a.credential, signed(a, b.agent, { core, text: "different body", key: "qa-mcp-codex-round-1", expMs: 60000 }), secrets);
    check("http idempotency conflict", conflictHttp.status === 409 && conflictHttp.error === "idempotency_conflict", `${conflictHttp.status}:${conflictHttp.error}`);
    const badSig = await post("/v1/test/send", a.credential, signed(a, b.agent, { core, text: "bad signature", key: "qa-bad-signature-01", expMs: 60000, corrupt: true }), secrets);
    check("bad signature rejected", badSig.status === 422 && badSig.error === "invalid_signature", `${badSig.status}:${badSig.error}`);
    const longTtl = await post("/v1/test/send", a.credential, signed(a, b.agent, { core, text: "too long", key: "qa-ttl-too-long-01", expMs: 310000 }), secrets);
    check("ttl above 300s rejected", longTtl.status === 422 && longTtl.error === "ttl_too_long", `${longTtl.status}:${longTtl.error}`);
    const expired = await post("/v1/test/send", a.credential, signed(a, b.agent, { core, text: "expired", key: "qa-ttl-expired-01", expMs: -5000 }), secrets);
    check("expired trial rejected", expired.status === 422 && expired.error === "expired", `${expired.status}:${expired.error}`);
    const cross = await post("/v1/test/send", a.credential, signed(a, "trial_other", { core, text: "cross agent", key: "qa-cross-agent-0001", expMs: 60000, to: "trial_other" }), secrets);
    check("cross agent rejected", cross.status === 403 && cross.error === "sender_not_allowed", `${cross.status}:${cross.error}`);
    const business = await post("/v1/test/send", a.credential, signed(a, b.agent, { core, key: "qa-business-send-01", expMs: 60000, intent: "schedule.query", body: { window: { start: "2026-10-03T10:00:00Z", end: "2026-10-03T11:00:00Z" }, granularity_min: 30 } }), secrets);
    check("machine path rejects business intent", business.status === 403 && business.error === "sender_not_allowed", `${business.status}:${business.error}`);
    const anonymous = await post("/v1/test/send", "", signed(a, b.agent, { core, text: "no auth", key: "qa-no-auth-send-01", expMs: 60000 }), secrets);
    check("missing auth rejected", anonymous.status === 401 && anonymous.error === "invalid_auth", `${anonymous.status}:${anonymous.error}`);
    for (const path of ["/v1/test/owners", "/v1/test/agents", "/v1/test/invites", "/v1/test/authorize", "/v1/test/gate-consume"]) {
      const blocked = await post(path, a.credential, {}, secrets);
      check(`blocked ${path}`, blocked.status === 403 && blocked.error === "sender_not_allowed", `${blocked.status}:${blocked.error}`);
    }
    const ids = { codex_send: outbound.body.id, grok_receive: received.body.message.id, grok_reply: reply.body.id, codex_receive: returned.body.message.id, cli_send: sentBody.id, cli_reply: cliReply.body.id };
    await writeFile(process.env.STATE_PATH, JSON.stringify({ credential: a.credential, pem: a.pem, agent: a.agent, peer: b.agent }), { mode: 0o600 });
    const noisy = [sent.stdout, sent.stderr, got.stdout, got.stderr, grokMcp.stderr].join("\n");
    check("command output hides credentials", secrets.every((secret) => secret && !noisy.includes(secret)));
    return ids;
  } finally {
    for (const mcp of clients) await mcp.close();
    await rm(directory, { recursive: true, force: true });
  }
}

export async function closed(plugin) {
  const state = JSON.parse(await readFile(process.env.STATE_PATH, "utf8"));
  const core = await import(pathToFileURL(join(clone, "adapters/dist/core.js")).href);
  const result = await post("/v1/test/send", state.credential, signed({ agent: state.agent, pem: state.pem }, state.peer, { core, text: "after allowlist cleared", key: "qa-allowlist-off-01", expMs: 60000 }), [state.credential, state.pem]);
  check("empty allowlist rejects signed trial", result.status === 403 && result.error === "sender_not_allowed", `${result.status}:${result.error}`);
  await rm(process.env.STATE_PATH, { force: true });
  return {};
}
