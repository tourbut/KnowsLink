// Independent HTTP probe for the 4262d02 deliver:human boundary and related regressions.
// Prints statuses and error codes only. Credentials, CSRF, and claim tokens stay out of the log.
import { execFileSync } from "node:child_process";
import { generateKeyPairSync, randomBytes, sign } from "node:crypto";
import { writeFile } from "node:fs/promises";
import { signingBytes, uuid7 } from "../../../../../adapters/dist/index.js";

const base = process.env.RELAY_URL;
const checks = [];
if (!base || !["127.0.0.1", "localhost"].includes(new URL(base).hostname)) {
  console.error("RELAY_URL must be loopback");
  process.exit(1);
}

function record(id, result, detail) {
  checks.push({ id, result, detail });
  console.info(`${result.toUpperCase()} ${id}: ${detail}`);
}
const pass = (id, ok, detail) => record(id, ok ? "pass" : "fail", detail);

async function call(method, path, token, body, claim, auth = "bearer") {
  const headers = { "Content-Type": body === undefined ? "text/plain" : "application/json" };
  if (token && auth === "bearer") headers.Authorization = `Bearer ${token}`;
  if (token && auth === "basic") headers.Authorization = `Basic ${Buffer.from(`owner:${token}`).toString("base64")}`;
  if (claim) headers["X-Execution-Claim"] = claim;
  const response = await fetch(`${base}${path}`, {
    method,
    headers,
    body: body === undefined ? undefined : typeof body === "string" ? body : JSON.stringify(body),
    redirect: "manual",
  });
  const text = await response.text();
  let json = null;
  try {
    json = JSON.parse(text);
  } catch {
    json = null;
  }
  return { status: response.status, text, json };
}
function psql(sql) {
  return execFileSync(
    "docker",
    ["compose", "--project-name", process.env.COMPOSE_PROJECT, "--env-file", ".env.example", "-f", "compose.yaml", "-f", process.env.COMPOSE_OVERRIDE, "exec", "-T", "postgres", "psql", "-U", "knowslink", "-d", "knowslink", "-v", "ON_ERROR_STOP=1", "-q", "-t", "-A", "-c", sql],
    { encoding: "utf8" },
  ).trim();
}
function qid(id) {
  if (!/^[0-9a-f-]+$/.test(id)) throw new Error("unexpected id");
  return id;
}
function plus(ms) {
  return new Date(Date.now() + ms).toISOString().replace(/\.\d{3}Z$/, "Z");
}
const queryBody = { window: { start: "2026-10-03T10:00:00Z", end: "2026-10-03T11:00:00Z" }, granularity_min: 30 };
const commitBody = { slot: queryBody.window, timezone: "UTC", commitment: "synthetic-stub" };
function build(from, fields) {
  const message = {
    v: "relay.v1", id: fields.id || uuid7(), from: from.agent, to: fields.to, intent: fields.intent,
    body: fields.body, deliver: fields.deliver || (fields.intent === "relay.approval.request" ? "human" : "agent"),
    exp: fields.exp || plus(180000), idempotency_key: fields.key || randomBytes(18).toString("base64url"),
    sig: { alg: "Ed25519", kid: fields.kid || from.kid || "key1", value: "" },
  };
  if (fields.reply_to) message.reply_to = fields.reply_to;
  if (fields.render) message.render = fields.render;
  message.sig.value = sign(null, signingBytes(message), from.privateKey).toString("base64url");
  return message;
}
async function createAgent(suffix) {
  const owner = await call("POST", "/v1/owners", "", {});
  const agent = `qa_${suffix}_${randomBytes(4).toString("hex")}`;
  const { publicKey, privateKey } = generateKeyPairSync("ed25519");
  const publicBytes = publicKey.export({ format: "der", type: "spki" }).subarray(-32).toString("base64url");
  const proof = sign(null, Buffer.from(["KNOWSLINK-KEY-POP", owner.json.owner, agent, "key1", publicBytes].join("\0")), privateKey).toString("base64url");
  const registered = await call("POST", "/v1/agents", owner.json.credential, { agent, kid: "key1", public: publicBytes, proof });
  if (owner.status !== 200 || registered.status !== 200) throw new Error(`signup failed ${owner.status}/${registered.status}`);
  return { agent, owner: owner.json.owner, ownerCredential: owner.json.credential, credential: registered.json.credential, privateKey, publicBytes, kid: "key1" };
}
async function drain(agent) {
  for (let i = 0; i < 15; i++) {
    const lease = await call("POST", "/v1/pull", agent.credential, {});
    if (!lease.json?.lease_token) return;
    const id = lease.json.envelope.id;
    const token = lease.json.lease_token;
    await call("POST", "/v1/persist", agent.credential, { id, token });
    await call("POST", "/v1/ack", agent.credential, { id, token });
  }
}
function present(id) {
  return psql(`SELECT data->'Messages' ? '${qid(id)}' FROM relay_state`);
}
function gates() {
  return psql(`SELECT (SELECT count(*) FROM jsonb_object_keys(COALESCE(data->'Gates', '{}'::jsonb))) FROM relay_state`);
}
async function waitHealth() {
  for (let i = 0; i < 30; i++) {
    const health = await call("GET", "/healthz").catch(() => ({ status: 0 }));
    if (health.status === 200) return health;
    await new Promise((resolve) => setTimeout(resolve, 1000));
  }
  return { status: 0 };
}

async function main() {
  const a = await createAgent("a");
  const b = await createAgent("b");
  const ownerDenied = await call("POST", "/v1/key-revoke", a.credential, { agent: a.agent, kid: "key1" });
  pass("REG-current-auth-agent-credential", ownerDenied.status === 401, `agent credential on owner revoke ${ownerDenied.status}`);

  const invited = await call("POST", "/v1/invites", a.credential, { agent: a.agent, target: b.agent });
  const accepted = await call("POST", "/v1/invite-decision", b.ownerCredential, { agent: a.agent, target: b.agent, decision: "accept" });
  pass("REG-pair", invited.status === 200 && accepted.json?.State === "active", `invite ${invited.status}, accept ${accepted.json?.State}`);

  const humanQuery = build(a, { to: b.agent, intent: "schedule.query", body: queryBody, deliver: "human" });
  const humanSend = await call("POST", "/v1/send", a.credential, humanQuery);
  const humanPull = await call("POST", "/v1/pull", b.credential, {});
  const humanClaim = await call("POST", "/v1/claim", b.credential, { id: humanQuery.id });
  const humanReceipt = await call("GET", `/v1/receipts/${humanQuery.id}`, a.credential);
  pass(
    "FIX-direct-human-query",
    humanSend.status === 403 && humanSend.json?.error === "sender_not_allowed" && present(humanQuery.id) === "f" && humanPull.status === 200 && humanPull.text.trim() === "null" && humanClaim.status === 403 && humanReceipt.status === 403 && gates() === "0",
    `send ${humanSend.status} ${humanSend.json?.error}, stored ${present(humanQuery.id)}, pull ${humanPull.status} body-null ${humanPull.text.trim() === "null"}, claim ${humanClaim.status}, receipt ${humanReceipt.status}, gates ${gates()}`,
  );
  const humanCommit = build(a, { to: b.agent, intent: "schedule.commit", body: commitBody, deliver: "human" });
  const commitSend = await call("POST", "/v1/send", a.credential, humanCommit);
  pass(
    "FIX-direct-human-commit",
    commitSend.status === 403 && commitSend.json?.error === "sender_not_allowed" && present(humanCommit.id) === "f",
    `send ${commitSend.status} ${commitSend.json?.error}, stored ${present(humanCommit.id)}`,
  );
  const humanResult = build(b, { to: a.agent, intent: "relay.result", body: { status: "denied" }, deliver: "human", reply_to: humanQuery.id });
  const resultSchema = await call("POST", "/v1/send", b.credential, humanResult);
  pass("FIX-result-human-stays-schema", resultSchema.status === 422 && resultSchema.json?.error === "invalid_schema", `result deliver human ${resultSchema.status} ${resultSchema.json?.error}`);

  const stored = build(a, { to: b.agent, intent: "schedule.query", body: queryBody });
  const storedSend = await call("POST", "/v1/send", a.credential, stored);
  psql(`UPDATE relay_state SET data = jsonb_set(data, '{Messages,${qid(stored.id)},Deliver}', '"human"')`);
  const storedPull = await call("POST", "/v1/pull", b.credential, {});
  psql(`UPDATE relay_state SET data = jsonb_set(jsonb_set(jsonb_set(data, '{Messages,${qid(stored.id)},Receipt,transport}', '"leased"'), '{Messages,${qid(stored.id)},LeaseToken}', '"forged-lease"'), '{Messages,${qid(stored.id)},LeaseUntil}', '"2099-01-01T00:00:00Z"')`);
  const forgedPersist = await call("POST", "/v1/persist", b.credential, { id: stored.id, token: "forged-lease" });
  const forgedAck = await call("POST", "/v1/ack", b.credential, { id: stored.id, token: "forged-lease" });
  psql(`UPDATE relay_state SET data = jsonb_set(jsonb_set(jsonb_set(data, '{Messages,${qid(stored.id)},Receipt,transport}', '"delivered"'), '{Messages,${qid(stored.id)},Persisted}', 'true'), '{Messages,${qid(stored.id)},Inbox}', data->'Messages'->'${qid(stored.id)}'->'Envelope')`);
  const storedClaim = await call("POST", "/v1/claim", b.credential, { id: stored.id });
  const storedClaimed = psql(`SELECT COALESCE(data->'Messages'->'${qid(stored.id)}'->>'Claimed', 'null') FROM relay_state`);
  pass(
    "FIX-stored-human-transport",
    storedSend.status === 200 && storedPull.status === 200 && storedPull.text.trim() === "null" && forgedPersist.status === 409 && forgedAck.status === 409 && storedClaim.status === 403 && storedClaimed !== "true" && gates() === "0",
    `send ${storedSend.status}, pull-null ${storedPull.text.trim() === "null"}, persist ${forgedPersist.status} ${forgedPersist.json?.error}, ack ${forgedAck.status} ${forgedAck.json?.error}, claim ${storedClaim.status}, claimed ${storedClaimed}, gates ${gates()}`,
  );

  const live = build(a, { to: b.agent, intent: "schedule.query", body: queryBody });
  const liveSend = await call("POST", "/v1/send", a.credential, live);
  const lease = await call("POST", "/v1/pull", b.credential, {});
  const leaseGap = (Date.parse(lease.json?.lease_until) - Date.now()) / 1000;
  const earlyAck = await call("POST", "/v1/ack", b.credential, { id: live.id, token: lease.json?.lease_token });
  const persisted = await call("POST", "/v1/persist", b.credential, { id: live.id, token: lease.json?.lease_token });
  const acked = await call("POST", "/v1/ack", b.credential, { id: live.id, token: lease.json?.lease_token });
  const receipt = await call("GET", `/v1/receipts/${live.id}`, a.credential);
  const claim = await call("POST", "/v1/claim", b.credential, { id: live.id });
  const claimAgain = await call("POST", "/v1/claim", b.credential, { id: live.id });
  pass(
    "REG-agent-delivery",
    liveSend.status === 200 && lease.json?.attempts === 1 && leaseGap > 25 && leaseGap <= 31 && earlyAck.status === 409 && persisted.status === 200 && acked.status === 200 && receipt.json?.receipt?.transport === "delivered" && claim.status === 200 && Boolean(claim.json?.claim) && claimAgain.status === 409,
    `send ${liveSend.status}, attempts ${lease.json?.attempts}, lease gap ${Number.isFinite(leaseGap) ? leaseGap.toFixed(1) : "none"}s, ack-before ${earlyAck.status}, persist ${persisted.status}, ack ${acked.status}, transport ${receipt.json?.receipt?.transport}, claim ${claim.status}, again ${claimAgain.status}`,
  );
  execFileSync("docker", ["compose", "--project-name", process.env.COMPOSE_PROJECT, "--env-file", ".env.example", "-f", "compose.yaml", "-f", process.env.COMPOSE_OVERRIDE, "restart", "relay"], { stdio: "inherit" });
  const health = await waitHealth();
  const afterRestart = await call("POST", "/v1/claim", b.credential, { id: live.id });
  const fresh = build(a, { to: b.agent, intent: "schedule.query", body: queryBody });
  const freshSend = await call("POST", "/v1/send", a.credential, fresh);
  const freshPull = await call("POST", "/v1/pull", b.credential, {});
  pass(
    "REG-agent-restart",
    health.status === 200 && afterRestart.status === 409 && freshSend.status === 200 && freshPull.json?.envelope?.id === fresh.id,
    `health ${health.status}, reclaim ${afterRestart.status}, fresh ${freshSend.status}, leased-id-match ${freshPull.json?.envelope?.id === fresh.id}`,
  );
  if (freshPull.json?.lease_token) {
    await call("POST", "/v1/persist", b.credential, { id: fresh.id, token: freshPull.json.lease_token });
    await call("POST", "/v1/ack", b.credential, { id: fresh.id, token: freshPull.json.lease_token });
  }

  const legacyShape = psql(`UPDATE relay_state SET data = jsonb_set(data #- '{Messages,${qid(live.id)},Deliver}', '{Messages,${qid(live.id)},Envelope,deliver}', '"human"') RETURNING (data->'Messages'->'${qid(live.id)}'->'Envelope'->>'deliver') || ':' || ((data->'Messages'->'${qid(live.id)}' ? 'Deliver')::text)`);
  const legacyClaim = await call("POST", "/v1/claim", b.credential, { id: live.id });
  const legacyAuth = await call("POST", "/v1/authorize", b.credential, { id: live.id, claim: claim.json?.claim });
  const legacyReceipt = await call("GET", `/v1/receipts/${live.id}`, b.credential);
  const legacyResult = build(b, { to: a.agent, intent: "relay.result", body: { status: "denied" }, reply_to: live.id, exp: legacyReceipt.json?.receipt?.exp });
  const legacySend = await call("POST", "/v1/send", b.credential, legacyResult, claim.json?.claim);
  const legacyCompletion = await call("GET", `/v1/receipts/${live.id}`, a.credential);
  pass("FIX-legacy-new-claim-denied", legacyShape === "human:false" && legacyClaim.status === 403, `shape ${legacyShape}, new claim ${legacyClaim.status}`);
  pass(
    "FIX-legacy-claim-authorize-result",
    legacyShape === "human:false" && legacyAuth.status === 403 && legacySend.status === 403 && gates() === "0",
    `authorize ${legacyAuth.status} ${legacyAuth.json?.error || ""}, result ${legacySend.status} ${legacySend.json?.error || ""}, completion ${legacyCompletion.json?.completion ?? ""}, gates ${gates()}`,
  );
  await drain(a);

  const tooLongKey = randomBytes(18).toString("base64url");
  const tooLong = await call("POST", "/v1/send", a.credential, build(a, { to: b.agent, intent: "schedule.query", body: queryBody, key: tooLongKey, exp: plus(301000) }));
  const rolled = await call("POST", "/v1/send", a.credential, build(a, { to: b.agent, intent: "schedule.query", body: queryBody, key: tooLongKey }));
  pass("REG-ttl", tooLong.status === 422 && tooLong.json?.error === "ttl_too_long" && rolled.status === 200, `ttl ${tooLong.status} ${tooLong.json?.error}, rollback ${rolled.status}`);

  await drain(b);
  const attempts = build(a, { to: b.agent, intent: "schedule.query", body: queryBody });
  const attemptsSend = await call("POST", "/v1/send", a.credential, attempts);
  for (let i = 0; i < 3; i++) {
    const next = await call("POST", "/v1/pull", b.credential, {});
    if (next.json?.envelope?.id === attempts.id) psql(`UPDATE relay_state SET data = jsonb_set(data, '{Messages,${qid(attempts.id)},LeaseUntil}', '"2020-01-01T00:00:00Z"')`);
  }
  const exhausted = await call("GET", `/v1/receipts/${attempts.id}`, a.credential);
  pass("REG-lease-max-attempts", attemptsSend.status === 200 && exhausted.json?.receipt?.transport === "failed:max_attempts", `send ${attemptsSend.status}, transport ${exhausted.json?.receipt?.transport}`);

  const queued = build(a, { to: b.agent, intent: "schedule.query", body: queryBody });
  await call("POST", "/v1/send", a.credential, queued);
  const revoked = await call("POST", "/v1/key-revoke", b.ownerCredential, { agent: b.agent, kid: "key1" });
  const leasedAfter = await call("POST", "/v1/pull", b.credential, {});
  const freshKey = generateKeyPairSync("ed25519");
  const freshPublic = freshKey.publicKey.export({ format: "der", type: "spki" }).subarray(-32).toString("base64url");
  const restored = await call("POST", "/v1/keys", b.ownerCredential, { agent: b.agent, kid: "key2", public: freshPublic, proof: sign(null, Buffer.from(["KNOWSLINK-KEY-POP", b.owner, b.agent, "key2", freshPublic].join("\0")), freshKey.privateKey).toString("base64url") });
  b.privateKey = freshKey.privateKey;
  b.kid = "key2";
  const pullOld = await call("POST", "/v1/pull", b.credential, {});
  const rotated = build(a, { to: b.agent, intent: "schedule.query", body: queryBody });
  const rotatedSend = await call("POST", "/v1/send", a.credential, rotated);
  const pullFresh = await call("POST", "/v1/pull", b.credential, {});
  pass(
    "REG-revoke",
    revoked.status === 200 && leasedAfter.status === 200 && !leasedAfter.json?.lease_token && restored.status === 200 && !pullOld.json?.lease_token && rotatedSend.status === 200 && pullFresh.json?.envelope?.id === rotated.id,
    `revoke ${revoked.status}, revoked-token ${Boolean(leasedAfter.json?.lease_token)}, rotate ${restored.status}, old-restored ${Boolean(pullOld.json?.lease_token)}, fresh-leased ${pullFresh.json?.envelope?.id === rotated.id}`,
  );
  const unpaired = await call("POST", "/v1/unpair", b.ownerCredential, { agent: a.agent, target: b.agent });
  const ackUnpaired = await call("POST", "/v1/ack", b.credential, { id: rotated.id, token: pullFresh.json?.lease_token });
  const reinvite = await call("POST", "/v1/invites", a.credential, { agent: a.agent, target: b.agent });
  const reaccept = await call("POST", "/v1/invite-decision", b.ownerCredential, { agent: a.agent, target: b.agent, decision: "accept" });
  const oldReplay = await call("POST", "/v1/send", a.credential, queued);
  const newSend = await call("POST", "/v1/send", a.credential, build(a, { to: b.agent, intent: "schedule.query", body: queryBody }));
  pass(
    "REG-unpair-generation",
    unpaired.status === 200 && ackUnpaired.status === 409 && reinvite.json?.Generation > 1 && reaccept.json?.State === "active" && oldReplay.status === 403 && newSend.status === 200,
    `unpair ${unpaired.status}, ack ${ackUnpaired.status}, generation ${reinvite.json?.Generation}, old ${oldReplay.status}, new ${newSend.status}`,
  );
  psql(`UPDATE relay_state SET clock = clock_timestamp() + interval '1 hour'`);
  const skewed = await call("POST", "/v1/pull", b.credential, {});
  psql(`UPDATE relay_state SET clock = clock_timestamp()`);
  const recovered = await call("GET", "/healthz");
  const unknown = await call("POST", "/v1/pull", "not-a-credential", {});
  pass("REG-clock-unknown-auth", skewed.status === 503 && skewed.json?.error === "abnormal_clock" && recovered.status === 200 && unknown.status === 401, `clock ${skewed.status} ${skewed.json?.error}, health ${recovered.status}, unknown ${unknown.status}`);

  await drain(b);
  const parent = build(a, { to: b.agent, intent: "schedule.query", body: queryBody, render: { hint: "이 힌트는 승인 근거가 아니다" } });
  const parentSend = await call("POST", "/v1/send", a.credential, parent);
  const parentLease = await call("POST", "/v1/pull", b.credential, {});
  const beforePersist = await call("POST", "/v1/ack", b.credential, { id: parent.id, token: parentLease.json?.lease_token });
  await call("POST", "/v1/persist", b.credential, { id: parent.id, token: parentLease.json?.lease_token });
  const beforeAck = await call("GET", `/v1/receipts/${parent.id}`, a.credential);
  await call("POST", "/v1/ack", b.credential, { id: parent.id, token: parentLease.json?.lease_token });
  const parentClaim = await call("POST", "/v1/claim", b.credential, { id: parent.id });
  const parentReceipt = (await call("GET", `/v1/receipts/${parent.id}`, b.credential)).json?.receipt;
  const unbound = build(b, { to: b.agent, intent: "relay.approval.request", body: { reason: "judgment_required", request_digest: parentReceipt?.digest }, reply_to: parent.id, exp: parentReceipt?.exp });
  const noClaimH = await call("POST", "/v1/send", b.credential, unbound);
  const bound = build(b, { to: b.agent, intent: "relay.approval.request", body: { reason: "judgment_required", request_digest: parentReceipt?.digest }, reply_to: parent.id, exp: parentReceipt?.exp });
  const withClaim = await call("POST", "/v1/send", b.credential, bound, parentClaim.json?.claim);
  const humanPullH = await call("POST", "/v1/pull", b.credential, {});
  const hDeliver = psql(`SELECT data->'Messages'->'${qid(bound.id)}'->>'Deliver' FROM relay_state`);
  pass(
    "REG-owner-gate-open",
    parentSend.status === 200 && parentLease.json?.envelope?.id === parent.id && beforePersist.status === 409 && beforeAck.json?.receipt?.transport !== "delivered" && noClaimH.status === 403 && withClaim.status === 200 && hDeliver === "human" && humanPullH.text.trim() === "null",
    `parent ${parentSend.status}, ack-before ${beforePersist.status}, pre-ack ${beforeAck.json?.receipt?.transport}, H-without-claim ${noClaimH.status}, H ${withClaim.status} deliver ${hDeliver}, H-not-pulled ${humanPullH.text.trim() === "null"}`,
  );
  const gatePage = await call("GET", `/owner/gates/${bound.id}`, b.ownerCredential);
  const csrf = /name="csrf" value="([^"]+)"/.exec(gatePage.text)?.[1];
  const agentPage = await call("GET", `/owner/gates/${bound.id}`, b.credential);
  const basic = await call("GET", `/owner/gates/${bound.id}`, b.ownerCredential, undefined, undefined, "basic");
  const badCsrf = await fetch(`${base}/owner/gates/${bound.id}`, { method: "POST", redirect: "manual", headers: { Authorization: `Bearer ${b.ownerCredential}`, "Content-Type": "application/x-www-form-urlencoded" }, body: new URLSearchParams({ csrf: "bad", decision: "approve" }) });
  const getDecision = await call("GET", `/owner/gates/${bound.id}?decision=approve`, b.ownerCredential);
  const approved = await fetch(`${base}/owner/gates/${bound.id}`, { method: "POST", redirect: "manual", headers: { Authorization: `Bearer ${b.ownerCredential}`, "Content-Type": "application/x-www-form-urlencoded" }, body: new URLSearchParams({ csrf, decision: "approve" }) });
  const approvedPage = await call("GET", `/owner/gates/${bound.id}`, b.ownerCredential);
  const secondDecision = await fetch(`${base}/owner/gates/${bound.id}`, { method: "POST", redirect: "manual", headers: { Authorization: `Bearer ${b.ownerCredential}`, "Content-Type": "application/x-www-form-urlencoded" }, body: new URLSearchParams({ csrf, decision: "deny" }) });
  const consumed = await call("POST", "/v1/gate-consume", b.credential, { id: bound.id, claim: parentClaim.json?.claim });
  pass(
    "REG-gate-csrf",
    gatePage.status === 200 && gatePage.text.includes("granularity_min") && gatePage.text.includes("disclosure policy absent") && !gatePage.text.includes("이 힌트는 승인 근거가 아니다") && gatePage.text.includes("상태: pending") && agentPage.status === 401 && basic.status === 200 && badCsrf.status === 403 && getDecision.text.includes("상태: pending") && approved.status === 303 && approvedPage.text.includes("상태: approved") && !approvedPage.text.includes('value="approve"') && secondDecision.status === 409 && consumed.status === 200 && consumed.json?.executable === false && consumed.json?.disclosure === false,
    `page ${gatePage.status}, agent ${agentPage.status}, basic ${basic.status}, csrf ${badCsrf.status}, get ${getDecision.text.includes("pending")}, approve ${approved.status}, redecision ${secondDecision.status}, consume ${consumed.status} exec ${consumed.json?.executable} disclosure ${consumed.json?.disclosure}`,
  );

  const commit = build(a, { to: b.agent, intent: "schedule.commit", body: commitBody });
  await call("POST", "/v1/send", a.credential, commit);
  const commitLease = await call("POST", "/v1/pull", b.credential, {});
  await call("POST", "/v1/persist", b.credential, { id: commit.id, token: commitLease.json?.lease_token });
  await call("POST", "/v1/ack", b.credential, { id: commit.id, token: commitLease.json?.lease_token });
  const commitClaim = await call("POST", "/v1/claim", b.credential, { id: commit.id });
  const commitAuth = await call("POST", "/v1/authorize", b.credential, { id: commit.id, claim: commitClaim.json?.claim });
  const commitReceipt = (await call("GET", `/v1/receipts/${commit.id}`, b.credential)).json?.receipt;
  const done = build(b, { to: a.agent, intent: "relay.result", body: { status: "done" }, reply_to: commit.id, exp: commitReceipt?.exp });
  const doneSend = await call("POST", "/v1/send", b.credential, done, commitClaim.json?.claim);
  const wrong = build(b, { to: b.agent, intent: "relay.result", body: { status: "denied" }, reply_to: commit.id, exp: commitReceipt?.exp });
  const wrongSend = await call("POST", "/v1/send", b.credential, wrong, commitClaim.json?.claim);
  const leaked = build(b, { to: a.agent, intent: "relay.result", body: { status: "denied", result: { title: "private" } }, reply_to: commit.id, exp: commitReceipt?.exp });
  const leakedSend = await call("POST", "/v1/send", b.credential, leaked, commitClaim.json?.claim);
  const valid = build(b, { to: a.agent, intent: "relay.result", body: { status: "denied" }, reply_to: commit.id, exp: commitReceipt?.exp });
  const validSend = await call("POST", "/v1/send", b.credential, valid, commitClaim.json?.claim);
  const again = build(b, { to: a.agent, intent: "relay.result", body: { status: "failed" }, reply_to: commit.id, exp: commitReceipt?.exp });
  const againSend = await call("POST", "/v1/send", b.credential, again, commitClaim.json?.claim);
  const completion = await call("GET", `/v1/receipts/${commit.id}`, a.credential);
  const resultPull = await call("POST", "/v1/pull", a.credential, {});
  const resultPersist = await call("POST", "/v1/persist", a.credential, { id: valid.id, token: resultPull.json?.lease_token });
  const resultAck = await call("POST", "/v1/ack", a.credential, { id: valid.id, token: resultPull.json?.lease_token });
  const afterResult = await call("POST", "/v1/pull", a.credential, {});
  pass(
    "REG-approval-result",
    String(commitClaim.json?.policy).includes("non-executable") && commitAuth.json?.executable === false && doneSend.status === 403 && wrongSend.status === 403 && leakedSend.status === 422 && validSend.status === 200 && againSend.status === 403 && againSend.json?.error === "sender_not_allowed" && completion.json?.completion === "denied" && resultPull.json?.envelope?.intent === "relay.result" && resultPersist.status === 200 && resultAck.status === 200 && !afterResult.json?.lease_token,
    `done ${doneSend.status}, direction ${wrongSend.status}, optional ${leakedSend.status}, denied ${validSend.status}, duplicate ${againSend.status} ${againSend.json?.error}, completion ${completion.json?.completion}, result-ack ${resultAck.status}, next ${Boolean(afterResult.json?.lease_token)}`,
  );

  const out = ".fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-TESTER-FIX-test/probe-results.json";
  const failed = checks.filter((item) => item.result === "fail");
  await writeFile(out, JSON.stringify({ product: "4262d02fdd7b0b57a804d9e550597852950ffeae", checks }, null, 2));
  console.info(`PROBE fail ${failed.length} pass ${checks.filter((item) => item.result === "pass").length}`);
  process.exit(failed.length ? 1 : 0);
}
main().catch(async (error) => {
  console.error(error instanceof Error ? error.message : "probe failed");
  await writeFile(".fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-TESTER-FIX-test/probe-results.json", JSON.stringify({ product: "4262d02fdd7b0b57a804d9e550597852950ffeae", checks, error: "probe aborted" }, null, 2)).catch(() => {});
  process.exit(1);
});
