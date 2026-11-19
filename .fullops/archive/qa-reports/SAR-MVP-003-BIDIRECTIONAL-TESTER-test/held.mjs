// Held mode, fixed origin, redirect, timeout, response limit, and private launcher checks.
import { spawnSync } from "node:child_process";
import { mkdtemp, readFile, rm, writeFile, symlink } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join } from "node:path";
import { pathToFileURL } from "node:url";
import { baseEnv, call, check, clone, listen, Mcp, restoreEnv } from "./session.mjs";

export async function held(plugin) {
  let hits = 0;
  const { server, url } = await listen((_request, response) => {
    hits += 1;
    response.end("no");
  });
  const directory = await mkdtemp(join(tmpdir(), "knowslink-qa-held-"));
  try {
    const core = await readFile(join(clone, "adapters/src/core.ts"), "utf8");
    check("source timeout 10s", core.includes("timeoutMs = 10000"));
    check("source redirect error", core.includes('redirect: "error"'));
    check("source 64KiB", core.includes("size > 65536"));
    const { testTransport, TestTransport } = await import(pathToFileURL(join(clone, "adapters/dist/test-transport.js")).href);
    const { privateKey } = (await import("node:crypto")).generateKeyPairSync("ed25519");
    const pem = privateKey.export({ format: "pem", type: "pkcs8" }).toString();
    const key = join(directory, "key.pem");
    await writeFile(key, pem, { mode: 0o600 });
    Object.assign(process.env, {
      KNOWSLINK_MODE: "held",
      RELAY_URL: "https://link.knowslog.com",
      AGENT_ID: "trial_codex",
      AGENT_KID: "key1",
      AGENT_CREDENTIAL: "held-not-used",
      AGENT_KEY_FILE: "/missing",
      KNOWSLINK_TEST_PEER: "trial_grok",
      CF_ACCESS_CLIENT_ID: "held-client",
      CF_ACCESS_CLIENT_SECRET: "example-access",
    });
    check("held transport absent", (await testTransport()) === null);
    process.env.KNOWSLINK_MODE = "test-remote";
    delete process.env.CF_ACCESS_CLIENT_SECRET;
    check("remote missing access reads no key", (await testTransport()) === null);
    process.env.CF_ACCESS_CLIENT_SECRET = "example-access";
    process.env.AGENT_KEY_FILE = key;
    let remoteFetches = 0;
    const original = globalThis.fetch;
    globalThis.fetch = async (...args) => {
      remoteFetches += 1;
      return original(...args);
    };
    try {
      check("approved origin configures", Boolean(await testTransport()));
      for (const bad of ["http://link.knowslog.com", "https://example.com", "https://127.0.0.1", "https://link.knowslog.com/path", "https://link.knowslog.com?x=1", "https://user@link.knowslog.com", "https://link.knowslog.com:444"]) {
        process.env.RELAY_URL = bad;
        let rejected = false;
        try {
          await testTransport();
        } catch {
          rejected = true;
        }
        check(`reject origin ${bad}`, rejected);
      }
      check("origin checks made no remote fetch", remoteFetches === 0, String(remoteFetches));
    } finally {
      globalThis.fetch = original;
    }
    process.env.KNOWSLINK_MODE = "test-loopback";
    process.env.RELAY_URL = "https://link.knowslog.com";
    let loopRejected = false;
    try {
      await testTransport();
    } catch {
      loopRejected = true;
    }
    check("loopback rejects public origin", loopRejected);
    const { server: edge, url: edgeURL } = await listen((request, response) => {
      if (request.url?.endsWith("/send")) {
        response.end("x".repeat(65537));
        return;
      }
      if (request.url?.endsWith("/pull")) {
        response.writeHead(302, { Location: "http://127.0.0.1:1/leak" });
        response.end();
        return;
      }
      response.writeHead(200, { "content-type": "application/json" });
      response.write('{"partial":true}');
    });
    try {
      const transport = new TestTransport(edgeURL, "held-not-used", "trial_codex", "key1", pem, "trial_grok");
      let tooLarge = false;
      try {
        await transport.request("/v1/send", {});
      } catch (error) {
        tooLarge = String(error.message).includes("too large");
      }
      check("response above 64KiB rejected", tooLarge);
      let redirected = false;
      try {
        await transport.request("/v1/pull", {});
      } catch {
        redirected = true;
      }
      check("redirect rejected", redirected);
      const started = Date.now();
      const outcome = await Promise.race([
        transport.request("/v1/registry").then(() => "resolved", (error) => error),
        new Promise((resolve) => setTimeout(() => resolve("hung"), 15000)),
      ]);
      const elapsed = Date.now() - started;
      check("body timeout rejects", outcome instanceof Error, outcome === "hung" ? "hung" : outcome.name || "error");
      check("body timeout is about 10s", elapsed >= 9000 && elapsed < 15000, String(elapsed));
    } finally {
      edge.closeAllConnections();
      edge.close();
    }
    const probe = new Mcp(plugin, baseEnv({ KNOWSLINK_MODE: "held", RELAY_URL: url, AGENT_CREDENTIAL: "held-not-used", AGENT_ID: "trial_codex", AGENT_KID: "key1", AGENT_KEY_FILE: key, KNOWSLINK_TEST_PEER: "trial_grok" }));
    await probe.init();
    const listed = await probe.request("tools/list", {});
    const names = (listed.result?.tools || []).map((tool) => tool.name).sort();
    check("held discovery has four tools", names.join() === "knowslink_pull_once,knowslink_status,knowslink_test_receive,knowslink_test_send", names.join());
    const status = await call(probe, "knowslink_status");
    check("default status held", !status.isError && status.body.state === "held" && status.body.actualConnection === "held" && status.body.webhook === false && status.body.evidenceFetch === false);
    for (const name of ["knowslink_test_send", "knowslink_test_receive", "knowslink_pull_once"]) {
      const args = name === "knowslink_test_send" ? { text: "should not send", idempotency_key: "qa-held-should-not-send" } : {};
      const result = await call(probe, name, args);
      check(`${name} held no effect`, result.isError && result.body.state === "held" && result.body.actualConnection === "held");
    }
    await probe.close();
    check("held made no canary connection", hits === 0, String(hits));
    check("held stderr has no key", !probe.stderr.includes("PRIVATE KEY") && !probe.stderr.includes("example-access"));
    const sentinel = `sentinel-${Date.now()}-not-a-credential`;
    const publicConfig = join(directory, "public.json");
    await writeFile(publicConfig, JSON.stringify({ KNOWSLINK_MODE: "held", AGENT_CREDENTIAL: sentinel }), { mode: 0o644 });
    const denied = spawnSync("python3", [join(clone, "scripts/run_trial.py"), "--config", publicConfig, "receive"], { encoding: "utf8", cwd: "/tmp", env: baseEnv() });
    const deniedText = `${denied.stdout || ""}\n${denied.stderr || ""}`;
    check("mode 0644 launcher rejected", denied.status === 1, String(denied.status));
    check("rejected launcher hides sentinel", !deniedText.includes(sentinel));
    const link = join(directory, "link.json");
    await symlink(publicConfig, link);
    const linked = spawnSync("python3", [join(clone, "scripts/run_trial.py"), "--config", link, "receive"], { encoding: "utf8", cwd: "/tmp", env: baseEnv() });
    check("symlink config rejected", linked.status === 1, String(linked.status));
    const privateConfig = join(directory, "private.json");
    await writeFile(privateConfig, JSON.stringify({ KNOWSLINK_MODE: "held", AGENT_CREDENTIAL: "config-value" }), { mode: 0o600 });
    const wrapper = join(directory, "wrap.js");
    await writeFile(wrapper, "process.stdout.write(JSON.stringify({credential:process.env.AGENT_CREDENTIAL||'',access:process.env.CF_ACCESS_CLIENT_SECRET||''}))\n");
    const parentSecret = "parent-agent-credential";
    const parentAccess = "parent-access-secret";
    const launched = spawnSync("python3", [join(clone, "scripts/run_trial.py"), "--config", privateConfig, "--node", process.execPath, "--bundle", wrapper], {
      encoding: "utf8",
      cwd: "/tmp",
      env: baseEnv({ AGENT_CREDENTIAL: parentSecret, CF_ACCESS_CLIENT_SECRET: parentAccess }),
    });
    const launchedText = `${launched.stdout || ""}\n${launched.stderr || ""}`;
    check("0600 launcher replaces parent credential", launched.status === 0 && launched.stdout.includes('"credential":"config-value"') && launched.stdout.includes('"access":""'), String(launched.status));
    check("launcher output hides parent secrets", !launchedText.includes(parentSecret) && !launchedText.includes(parentAccess));
  } finally {
    restoreEnv();
    server.closeAllConnections();
    server.close();
    await rm(directory, { recursive: true, force: true });
  }
}
