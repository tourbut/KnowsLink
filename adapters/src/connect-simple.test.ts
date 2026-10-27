// SAR-CONNECT-SIMPLE-001 client: E reasons (KO/EN), retry folders (R-01/R-03), release defaults without requests before
// consent (T-03), host on start (A-04) and display names in tool results (A-06), against a fake relay only.
import assert from "node:assert/strict";
import { generateKeyPairSync, sign } from "node:crypto";
import {
  copyFile,
  mkdir,
  mkdtemp,
  readdir,
  readFile,
  rm,
  writeFile,
} from "node:fs/promises";
import { createServer } from "node:http";
import { tmpdir } from "node:os";
import { join } from "node:path";
import { fileURLToPath } from "node:url";
import { Client } from "@modelcontextprotocol/sdk/client/index.js";
import { StdioClientTransport } from "@modelcontextprotocol/sdk/client/stdio.js";
import { uuid7 } from "./core.js";
import { beginLogin, finishLogin } from "./login.js";
import { failureOf } from "./reasons.js";
import { setupConnect, type SetupEvent } from "./setup-connect.js";
import { textSigningBytes, type TextEnvelope } from "./text.js";

const root = await mkdtemp(join(tmpdir(), "kl-simple-"));
const peer = generateKeyPairSync("ed25519");
const starts: Record<string, string>[] = [];
const kids = new Map<string, string>();
let requests = 0;
let decision: "wait" | "approve" | "cancel" = "wait";
// One injected answer per path: [status, body].
const inject = new Map<string, [number, object]>();
const queue: TextEnvelope[] = [];
const relay = createServer(async (req, res) => {
  requests++;
  let raw = "";
  for await (const chunk of req) raw += chunk;
  res.setHeader("Content-Type", "application/json");
  const url = req.url ?? "";
  const reply = (status: number, body: unknown) => {
    res.statusCode = status;
    res.end(JSON.stringify(body));
  };
  const injected = inject.get(url);
  if (injected) return reply(...injected);
  if (url.startsWith("/v1/contacts/")) {
    const agent = decodeURIComponent(url.slice(13));
    // A connected agent renamed on the web; unnamed ones come back as the relay's label.
    const name = /^agent_\d+$/.test(agent)
      ? "Grok Bot 웹 이름"
      : agent === "agent_peer"
        ? "도우미"
        : `이름 없는 에이전트 (${agent.slice(-4)})`;
    return reply(200, {
      agent,
      name,
      client: "Grok",
      friend: "친구 A",
      short: agent.slice(-4),
    });
  }
  if (url.startsWith("/v1/keys/"))
    return reply(200, {
      public: peer.publicKey
        .export({ format: "der", type: "spki" })
        .subarray(-32)
        .toString("base64url"),
    });
  if (url === "/v1/text/pull")
    return reply(
      200,
      queue[0] ? { envelope: queue[0], lease_token: "lease" } : null,
    );
  const body = JSON.parse(raw || "{}") as Record<string, string>;
  if (url === "/v1/text/send")
    return reply(200, { id: body.id, transport: "queued" });
  if (url === "/v1/text/ack") queue.shift();
  if (!url.startsWith("/v1/connect/")) return reply(200, { state: "ok" });
  if (url === "/v1/connect/start") {
    starts.push(body);
    kids.set(body.token!, body.kid!);
  }
  const agent = `agent_${starts.findIndex((s) => s.token === body.token)}`;
  const exp = new Date(Date.now() + 600000).toISOString();
  if (decision === "cancel" && url !== "/v1/connect/start")
    return reply(401, { error: "invalid_auth" });
  if (url === "/v1/connect/complete")
    return reply(200, {
      agent,
      kid: kids.get(body.token!),
      credential: "c".repeat(43),
    });
  reply(
    200,
    url === "/v1/connect/start"
      ? { state: "requested", exp }
      : decision === "approve"
        ? { state: "approved", exp, owner: "owner", agent, name: "Grok Bot" }
        : { state: "prepared", exp },
  );
});
await new Promise<void>((r) => relay.listen(0, "127.0.0.1", r));
const address = relay.address();
assert(address && typeof address === "object");
const base = `http://127.0.0.1:${address.port}`;

const run = async (folder: string) => {
  const events: SetupEvent[] = [];
  const outcome = await setupConnect(
    base,
    folder,
    false,
    (e) => {
      events.push(e);
    },
    10,
  );
  return { outcome, last: events.at(-1)! };
};
type Tool = { state?: string; [key: string]: unknown };
async function mcp(plugin: string, env: Record<string, string>) {
  const client = new Client({ name: "connect-simple-test", version: "1.0.0" });
  await client.connect(
    new StdioClientTransport({
      command: process.execPath,
      args: [plugin],
      env: { LOCALAPPDATA: join(root, "localappdata"), ...env },
      stderr: "ignore",
    }),
  );
  const call = async (name: string, args: Record<string, unknown> = {}) => {
    const r = await client.callTool({ name, arguments: args });
    const text = (r.content as { text: string }[])[0]!.text;
    for (const secret of ["c".repeat(43), "PRIVATE KEY", ...kids.keys()])
      assert(!text.includes(secret), `${name} output must not carry secrets`);
    return JSON.parse(text) as Tool;
  };
  return { call, close: () => client.close() };
}

try {
  // E: each reason has its code and the exact D02 text, KO by default and EN with KNOWSLINK_LANG=en.
  process.env.KNOWSLINK_HOST = "claude";
  decision = "cancel";
  let r = await run(join(root, "cancel"));
  assert.deepEqual(
    [r.outcome, r.last.reason, r.last.message],
    [
      "failed",
      "cancelled",
      "연결을 취소했습니다. 다시 연결하면 새 링크를 만듭니다.",
    ],
  );
  assert.equal(starts.at(-1)!.host, "claude", "known host goes to start");
  delete process.env.KNOWSLINK_HOST;
  process.env.KNOWSLINK_HOST = "vim";
  decision = "wait";
  await beginLogin(base, join(root, "expired"));
  assert.equal(starts.at(-1)!.host, undefined, "unknown host is not sent");
  delete process.env.KNOWSLINK_HOST;
  const login = JSON.parse(
    await readFile(join(root, "expired", "login.json"), "utf8"),
  );
  await writeFile(
    join(root, "expired", "login.json"),
    JSON.stringify({ ...login, exp: new Date(Date.now() - 1).toISOString() }),
  );
  const expiredError = await finishLogin(join(root, "expired")).catch(
    (e: unknown) => e,
  );
  assert.equal(
    (await failureOf(expiredError, join(root, "expired"))).message,
    "연결 시간이 지났습니다. 다시 연결하면 새 링크를 만듭니다.",
  );

  process.env.KNOWSLINK_LANG = "en";
  inject.set("/v1/connect/start", [
    429,
    { error: "rate_limited", retry_at: "2026-10-26T12:34:56Z" },
  ]);
  r = await run(join(root, "rate"));
  assert.deepEqual(
    [r.last.state, r.last.reason, r.last.message, r.last.retry_at],
    [
      "failed",
      "rate_limited",
      "Too many requests. Connect again after 2026-10-26 12:34 UTC.",
      "2026-10-26T12:34:56Z",
    ],
  );
  inject.clear();
  decision = "approve";
  inject.set("/v1/connect/complete", [
    429,
    { error: "agent_rate", retry_at: "2026-10-27T01:02:03Z" },
  ]);
  r = await run(join(root, "agent-rate"));
  assert.deepEqual(
    [r.last.state, r.last.reason, r.last.message],
    [
      "failed",
      "agent_rate",
      "You've made many new agents today. Connect again after 2026-10-27 01:02 UTC. Your existing agents still work.",
    ],
  );
  inject.set("/v1/connect/complete", [409, { error: "service_capacity" }]);
  r = await run(join(root, "capacity"));
  assert.deepEqual(
    [r.last.reason, r.last.message],
    [
      "service_capacity",
      "The service can't add new agents right now. Your existing agents still work. Try again later.",
    ],
  );
  // Completion refused for good after approval: F-07 keeps agent and name.
  inject.set("/v1/connect/complete", [401, { error: "invalid_auth" }]);
  r = await run(join(root, "lost"));
  assert.deepEqual(
    [r.outcome, r.last.reason, r.last.agent, r.last.name],
    ["unconfirmed", "unconfirmed", `agent_${starts.length - 1}`, "Grok Bot"],
  );
  assert.equal(
    r.last.message,
    "You agreed in the browser, but key storage wasn't confirmed. Connect again; it checks this attempt first.",
  );
  inject.clear();
  // A-02/F-07: the approved attempt still completes after its link expired.
  const lostLogin = JSON.parse(
    await readFile(join(root, "lost", "pending.json"), "utf8"),
  );
  await writeFile(
    join(root, "lost", "login.json"),
    JSON.stringify({
      ...lostLogin,
      exp: new Date(Date.now() - 1).toISOString(),
    }),
  );
  assert.equal(await finishLogin(join(root, "lost")), "connected");
  delete process.env.KNOWSLINK_LANG;
  await writeFile(join(root, "not-a-folder"), "x");
  r = await run(join(root, "not-a-folder"));
  assert.deepEqual(
    [r.last.reason, r.last.message],
    [
      "private_folder",
      `키를 보관할 비공개 폴더를 준비하지 못했습니다(${join(root, "not-a-folder")}). 기본 위치로 다시 연결하세요.`,
    ],
  );
  const closed = createServer();
  await new Promise<void>((done) => closed.listen(0, "127.0.0.1", done));
  const closedPort = (closed.address() as { port: number }).port;
  await new Promise((done) => closed.close(done));
  r = await (async () => {
    const events: SetupEvent[] = [];
    const outcome = await setupConnect(
      `http://127.0.0.1:${closedPort}`,
      join(root, "offline"),
      false,
      (e) => {
        events.push(e);
      },
      10,
    );
    return { outcome, last: events.at(-1)! };
  })();
  assert.deepEqual(
    [r.last.reason, r.last.message],
    [
      "unreachable",
      "KnowsLink 서버에 연결할 수 없습니다. 인터넷 연결을 확인하고 다시 연결하세요.",
    ],
  );

  // T-03: a release Grok install (defaults.json beside dist/) needs no environment and makes no request before consent.
  decision = "wait";
  const release = join(root, "prefix", "knowslink");
  await mkdir(join(release, "dist"), { recursive: true });
  await copyFile(
    fileURLToPath(new URL("./plugin.js", import.meta.url)),
    join(release, "dist", "plugin.js"),
  );
  const connections = join(root, "prefix", "connections");
  await writeFile(
    join(release, "defaults.json"),
    JSON.stringify({
      KNOWSLINK_MODE: "public-node",
      RELAY_URL: base,
      KNOWSLINK_HOST: "grok",
      KNOWSLINK_AGENT_ROOT: connections,
    }),
  );
  const plugin = join(release, "dist", "plugin.js");
  // R-01: an expired attempt left in the default folder moves aside and is kept.
  const folder = join(connections, "grok-agent");
  await beginLogin(base, folder);
  const left = JSON.parse(await readFile(join(folder, "login.json"), "utf8"));
  await writeFile(
    join(folder, "login.json"),
    JSON.stringify({ ...left, exp: new Date(Date.now() - 1).toISOString() }),
  );
  const before = requests;
  let m = await mcp(plugin, {});
  assert.equal(
    (await m.call("knowslink_status")).state,
    "public_node_configured_unverified",
  );
  const stale = await m.call("knowslink_connect_status");
  assert.deepEqual(
    [stale.state, stale.reason, typeof stale.next],
    ["failed", "expired", "string"],
  );
  assert.equal(requests, before, "no relay request before explicit consent");
  const waiting = await m.call("knowslink_connect", { confirmed: true });
  assert.equal(waiting.state, "waiting");
  assert.equal(starts.at(-1)!.host, "grok");
  const aside = (await readdir(connections)).filter((n) =>
    n.startsWith("grok-agent.unused-"),
  );
  assert.equal(aside.length, 1);
  assert.deepEqual((await readdir(join(connections, aside[0]!))).sort(), [
    "login.json",
    "private.pem",
  ]);
  await m.close();

  // R-03: a restarted MCP process finds the same automatic folder; A-06 names in connect, status, send and receive.
  decision = "approve";
  m = await mcp(plugin, {});
  const done = await m.call("knowslink_connect_status");
  assert.deepEqual(done, {
    state: "connected",
    me: { agent: `agent_${starts.length - 1}`, name: "Grok Bot" },
  });
  const me = done.me as { agent: string };
  const status = await m.call("knowslink_status");
  assert.deepEqual(status.me, done.me);
  assert.equal(
    status.sharedAgent,
    "이 컴퓨터에서 기본 설정을 쓰는 Grok Bot은 모두 이 에이전트를 함께 씁니다.",
  );
  const sent = await m.call("knowslink_text_send", {
    peer: "agent_peer",
    text: "연결 확인",
    idempotency_key: "connect-simple-test-1",
    confirmed: true,
  });
  assert.deepEqual(
    { me: sent.me, peer: sent.peer },
    {
      me: { agent: me.agent, name: "Grok Bot 웹 이름" },
      peer: { agent: "agent_peer", name: "도우미" },
    },
  );
  const e: TextEnvelope = {
    v: "knowslink.text.v1",
    id: uuid7(),
    from: "agent_xyz9",
    to: me.agent,
    text: "안녕",
    exp: new Date(Date.now() + 120000).toISOString(),
    idempotency_key: "connect-simple-test-2",
    sig: { alg: "Ed25519", kid: "k1", value: "" },
  };
  e.sig.value = sign(null, textSigningBytes(e), peer.privateKey).toString(
    "base64url",
  );
  queue.push(e);
  // A failed own lookup falls back to the local name without breaking the result.
  inject.set(`/v1/contacts/${me.agent}`, [503, { error: "unavailable" }]);
  const received = await m.call("knowslink_text_receive");
  inject.clear();
  assert.deepEqual(
    { state: received.state, me: received.me, peer: received.peer },
    {
      state: "received",
      me: done.me,
      peer: { agent: "agent_xyz9", name: "이름 없는 에이전트 (xyz9)" },
    },
  );
  await m.close();
  // A set variable wins over release defaults (existing installs).
  m = await mcp(plugin, { KNOWSLINK_MODE: "held" });
  assert.equal((await m.call("knowslink_status")).state, "held");
  await m.close();

  // R-01: a pre-created empty folder is used as it is; R-02 private_folder reaches the MCP result with its next action.
  const dist = fileURLToPath(new URL("./plugin.js", import.meta.url));
  const empty = join(root, "empty");
  await mkdir(empty);
  m = await mcp(dist, {
    KNOWSLINK_MODE: "public-node",
    RELAY_URL: base,
    KNOWSLINK_AGENT_FOLDER: empty,
  });
  decision = "wait";
  assert.equal(
    (await m.call("knowslink_connect", { confirmed: true })).state,
    "waiting",
  );
  assert(!(await readdir(root)).some((n) => n.startsWith("empty.unused-")));
  await m.close();
  m = await mcp(dist, {
    KNOWSLINK_MODE: "public-node",
    RELAY_URL: base,
    KNOWSLINK_AGENT_FOLDER: join(root, "not-a-folder"),
    KNOWSLINK_LANG: "en",
  });
  const bad = await m.call("knowslink_connect", { confirmed: true });
  assert.equal(bad.reason, "private_folder");
  assert.equal(
    bad.message,
    `Couldn't prepare a private folder for the key (${join(root, "not-a-folder")}). Connect again using the default location.`,
  );
  assert.match(String(bad.next), /knowslink_connect/);
  await m.close();
  // Windows R-03 default: %LOCALAPPDATA%\KnowsLink\<host>-agent.
  if (process.platform === "win32") {
    m = await mcp(dist, {
      KNOWSLINK_MODE: "public-node",
      RELAY_URL: base,
      KNOWSLINK_HOST: "claude",
    });
    assert.equal(
      (await m.call("knowslink_connect", { confirmed: true })).state,
      "waiting",
    );
    await readFile(
      join(root, "localappdata", "KnowsLink", "claude-agent", "login.json"),
    );
    await m.close();
  }
  process.stdout.write(
    "PASS: E reasons KO/EN with code/text/next, approved completion after expiry, host on start, release defaults with 0 requests before consent, env wins, leftovers aside and kept, empty folder used, default folder across restart, display names in connect/status/send/receive\n",
  );
} finally {
  relay.closeAllConnections();
  relay.close();
  await rm(root, { recursive: true, force: true });
}
