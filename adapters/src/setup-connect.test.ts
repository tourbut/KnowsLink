// Installer connection step: new/resumed/cancelled attempts, lost completion (F-07) and the stored-key short cut.
import assert from "node:assert/strict";
import { mkdtemp, readdir, rm } from "node:fs/promises";
import { createServer } from "node:http";
import { tmpdir } from "node:os";
import { join } from "node:path";
import { beginLogin } from "./login.js";
import { setupConnect, type SetupEvent } from "./setup-connect.js";

const root = await mkdtemp(join(tmpdir(), "kl-setup-"));
const tokens: string[] = [];
const kids = new Map<string, string>();
let decision: "wait" | "approve" | "cancel" = "wait";
let loseComplete = false;
const consumed = new Set<string>();
let calls = 0;
const server = createServer(async (req, res) => {
  calls++;
  res.setHeader("Content-Type", "application/json");
  let raw = "";
  for await (const chunk of req) raw += chunk;
  const body: Record<string, string> = JSON.parse(raw);
  if (req.url === "/v1/connect/start") {
    tokens.push(body.token!);
    kids.set(body.token!, body.kid!);
  }
  const agent = `agent_${tokens.indexOf(body.token!)}`;
  const exp = new Date(Date.now() + 600000).toISOString();
  // The relay keeps cancelled and consumed attempts unusable (invalid_auth).
  if (decision === "cancel" && req.url !== "/v1/connect/start")
    consumed.add(body.token!);
  if (consumed.has(body.token!)) {
    res.statusCode = 401;
    res.end(JSON.stringify({ error: "invalid_auth" }));
  } else if (req.url === "/v1/connect/complete") {
    consumed.add(body.token!);
    if (loseComplete) res.destroy();
    else
      res.end(
        JSON.stringify({
          agent,
          kid: kids.get(body.token!),
          credential: "c".repeat(43),
        }),
      );
  } else
    res.end(
      JSON.stringify(
        req.url === "/v1/connect/start"
          ? { state: "requested", exp }
          : decision === "approve"
            ? {
                state: "approved",
                exp,
                owner: "owner",
                agent,
                name: "Claude 세션",
              }
            : { state: "prepared", exp },
      ),
    );
});
const run = async (
  folder: string,
  fresh = false,
  until?: (e: SetupEvent) => void,
) => {
  const events: SetupEvent[] = [];
  const outcome = await setupConnect(
    base,
    folder,
    fresh,
    (e) => {
      events.push(e);
      until?.(e);
    },
    10,
  );
  return { outcome, events };
};
let base = "";
try {
  await new Promise<void>((resolve) => server.listen(0, "127.0.0.1", resolve));
  const address = server.address();
  assert(address && typeof address === "object");
  base = `http://127.0.0.1:${address.port}`;
  const folder = join(root, "claude-01-private");

  // Cancel: no key stored; the next run starts a new attempt and keeps the old one aside.
  decision = "cancel";
  const cancelled = await run(folder);
  assert.equal(cancelled.outcome, "failed");
  assert.equal(cancelled.events.at(-1)!.reason, "cancelled");
  assert.equal(
    cancelled.events.at(-1)!.message,
    "연결을 취소했습니다. 다시 연결하면 새 링크를 만듭니다.",
  );
  const firstPrint = cancelled.events[0]!.fingerprint;
  assert(!JSON.stringify(cancelled.events).includes(tokens[0]!));

  // Resume: an interrupted waiting attempt keeps its URL and fingerprint.
  decision = "wait";
  let started: SetupEvent | undefined;
  const pending = run(folder, false, (e) => {
    if (e.state === "waiting") started = e;
  });
  while (!started) await new Promise((r) => setTimeout(r, 5));
  assert.notEqual(started.fingerprint, firstPrint);
  assert.equal(tokens.length, 2);
  assert(
    (await readdir(root)).some((n) =>
      n.startsWith("claude-01-private.unused-"),
    ),
  );
  decision = "approve";
  const done = await pending;
  if (done.outcome !== "connected")
    throw new Error(JSON.stringify(done.events));
  assert.equal(done.events.at(-1)!.agent, "agent_1");
  assert.equal(done.events.at(-1)!.name, "Claude 세션");
  assert.equal(done.events.at(-1)!.fingerprint, started.fingerprint);
  assert(!JSON.stringify(done.events).includes("c".repeat(43)));

  // Stored key: no relay call, no read of agent.json.
  const before = calls;
  const again = await run(folder);
  assert.equal(again.outcome, "connected");
  assert.equal(again.events[0]!.existing, "true");
  assert.equal(calls, before);

  // Resume: an installer stopped while waiting leaves a live attempt; the next run reuses its URL and fingerprint.
  const left = join(root, "other-private");
  decision = "wait";
  const interrupted = await beginLogin(base, left);
  const startsBefore = tokens.length;
  const resumed = await run(left, false, (e) => {
    if (e.resumed) decision = "approve";
  });
  assert.equal(resumed.outcome, "connected");
  assert.equal(resumed.events[0]!.url, interrupted.url);
  assert.equal(resumed.events.at(-1)!.fingerprint, interrupted.fingerprint);
  assert.equal(tokens.length, startsBefore);

  // F-07: approval happened, completion reached the relay but the answer was lost.
  const lost = join(root, "codex-01-private");
  decision = "approve";
  loseComplete = true;
  const unconfirmed = await run(lost);
  assert.equal(unconfirmed.outcome, "unconfirmed");
  assert.equal(unconfirmed.events.at(-1)!.agent, "agent_3");
  assert.equal(unconfirmed.events.at(-1)!.name, "Claude 세션");
  loseComplete = false;
  const starts = tokens.length;
  const recheck = await run(lost);
  assert.equal(
    recheck.outcome,
    "unconfirmed",
    "re-run checks the saved result first",
  );
  assert.equal(tokens.length, starts, "no new attempt before confirmation");
  const fresh = await run(lost, true);
  assert.equal(fresh.outcome, "connected");
  assert.equal(tokens.length, starts + 1);
  process.stdout.write(
    "PASS: setup-connect new/resume/cancel/F-07 confirm-first/fresh restart/stored key without relay call\n",
  );
} finally {
  server.closeAllConnections();
  server.close();
  await rm(root, { recursive: true, force: true });
}
