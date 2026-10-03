// Independent HTTP probe for SAR-MVP-001 QA-01–11. Prints statuses only, never credentials.
import { execFileSync } from "node:child_process";
import { generateKeyPairSync, randomBytes, sign } from "node:crypto";
import { mkdir, writeFile } from "node:fs/promises";
import { signingBytes, uuid7 } from "../../../../../adapters/dist/index.js";

const base = process.env.RELAY_URL;
const checks = [];
const pages = [];
if (!base || !["127.0.0.1", "localhost"].includes(new URL(base).hostname)) {
  console.error("RELAY_URL must be loopback");
  process.exit(1);
}

function record(id, result, detail) {
  checks.push({ id, result, detail });
  console.info(`${result.toUpperCase()} ${id}: ${detail}`);
}
const pass = (id, ok, detail) => record(id, ok ? "pass" : "fail", detail);
const held = (id, detail) => record(id, "held", detail);

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
    ["compose", "--project-name", process.env.COMPOSE_PROJECT, "--env-file", ".env.example", "-f", "compose.yaml", "-f", process.env.COMPOSE_OVERRIDE, "exec", "-T", "postgres", "psql", "-U", "knowslink", "-d", "knowslink", "-v", "ON_ERROR_STOP=1", "-t", "-A", "-c", sql],
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
function build(from, fields) {
  const message = {
    v: "relay.v1", id: fields.id || uuid7(), from: from.agent, to: fields.to, intent: fields.intent,
    body: fields.body, deliver: fields.deliver || (fields.intent === "relay.approval.request" ? "human" : "agent"),
    exp: fields.exp || plus(180000), idempotency_key: fields.key || randomBytes(18).toString("base64url"),
    sig: { alg: "Ed25519", kid: fields.kid || from.kid || "key1", value: "" },
  };
  for (const key of ["reply_to", "priority", "render", "evidence"]) if (fields[key]) message[key] = fields[key];
  message.sig.value = sign(null, signingBytes(message), from.privateKey).toString("base64url");
  return message;
}
async function createAgent(suffix) {
  const owner = await call("POST", "/v1/owners", "", {});
  const agent = `qa_${suffix}_${randomBytes(4).toString("hex")}`;
  const { publicKey, privateKey } = generateKeyPairSync("ed25519");
  const publicBytes = publicKey.export({ format: "der", type: "spki" }).subarray(-32).toString("base64url");
  const proofFor = (ownerId, kid, bytes, key = privateKey) => sign(null, Buffer.from(["KNOWSLINK-KEY-POP", ownerId, agent, kid, bytes].join("\0")), key).toString("base64url");
  const registered = await call("POST", "/v1/agents", owner.json.credential, { agent, kid: "key1", public: publicBytes, proof: proofFor(owner.json.owner, "key1", publicBytes) });
  if (owner.status !== 200 || registered.status !== 200) throw new Error(`signup failed ${owner.status}/${registered.status}`);
  return { agent, owner: owner.json.owner, ownerCredential: owner.json.credential, credential: registered.json.credential, privateKey, publicBytes, proofFor, kid: "key1" };
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
async function save(name, html) {
  const redacted = html.replace(/name="csrf" value="[^"]*"/g, 'name="csrf" value="redacted"');
  const file = `.fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-TESTER-test/ui/${name}.html`;
  await writeFile(file, redacted);
  pages.push({ id: name, file });
}

async function main() {
  await mkdir(".fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-TESTER-test/ui", { recursive: true });
  const a = await createAgent("a");
  const b = await createAgent("b");
  const c = await createAgent("c");
  const d = await createAgent("d");
  pass("QA-01-signup", true, "POST /v1/owners and PoP /v1/agents returned 200 for four synthetic owners");

  const ownerPaths = ["/v1/key-revoke", "/v1/invite-decision", "/v1/unpair", "/v1/keys"];
  const ownerBodies = [{ agent: a.agent, kid: "key1" }, { agent: a.agent, target: b.agent, decision: "accept" }, { agent: a.agent, target: b.agent }, { agent: a.agent, kid: "key9", public: a.publicBytes, proof: "aa" }];
  const ownerCodes = [];
  for (let i = 0; i < ownerPaths.length; i++) ownerCodes.push((await call("POST", ownerPaths[i], a.credential, ownerBodies[i])).status);
  pass("QA-01-agent-credential", ownerCodes.every((code) => code === 401), `agent credential on owner routes: ${ownerCodes.join(",")}`);

  const badKid = await call("POST", "/v1/agents", a.ownerCredential, { agent: "qa_bad_kid", kid: "https://keys.example/key", public: a.publicBytes, proof: "aa" });
  pass("QA-01-url-kid", badKid.status === 422, `URL kid status ${badKid.status} error ${badKid.json?.error}`);
  const { publicKey, privateKey } = generateKeyPairSync("ed25519");
  const rotated = publicKey.export({ format: "der", type: "spki" }).subarray(-32).toString("base64url");
  const wrong = await call("POST", "/v1/keys", a.ownerCredential, { agent: a.agent, kid: "key2", public: rotated, proof: a.proofFor(b.owner, "key2", rotated, privateKey) });
  const turned = await call("POST", "/v1/keys", a.ownerCredential, { agent: a.agent, kid: "key2", public: rotated, proof: sign(null, Buffer.from(["KNOWSLINK-KEY-POP", a.owner, a.agent, "key2", rotated].join("\0")), privateKey).toString("base64url") });
  a.privateKey = privateKey;
  a.kid = "key2";
  const reused = await call("POST", "/v1/keys", a.ownerCredential, { agent: a.agent, kid: "key1", public: a.publicBytes, proof: a.proofFor(a.owner, "key1", a.publicBytes) });
  const stale = await call("POST", "/v1/send", a.credential, build(a, { to: b.agent, intent: "schedule.query", body: queryBody, kid: "key1" }));
  pass("QA-01-rotate-reassign", wrong.status === 422 && turned.status === 200 && reused.status === 409 && stale.status === 401, `wrong-pop ${wrong.status}, rotate ${turned.status}, kid reuse ${reused.status} ${reused.json?.error}, previous key ${stale.status}`);

  const pending = await call("POST", "/v1/invites", a.credential, { agent: a.agent, target: b.agent });
  const early = await call("POST", "/v1/send", a.credential, build(a, { to: b.agent, intent: "schedule.query", body: queryBody }));
  const contactsBefore = await call("GET", "/v1/contacts", a.credential);
  pass("QA-02-before-accept", pending.status === 200 && pending.json.State === "pending" && early.status === 403 && !JSON.stringify(contactsBefore.json).includes(b.agent), `invite ${pending.json?.State}, send ${early.status} ${early.json?.error}, contacts ${JSON.stringify(contactsBefore.json)}`);
  const denied = await call("POST", "/v1/invites", a.credential, { agent: a.agent, target: d.agent });
  const deny = await call("POST", "/v1/invite-decision", d.ownerCredential, { agent: a.agent, target: d.agent, decision: "deny" });
  const deniedSend = await call("POST", "/v1/send", a.credential, build(a, { to: d.agent, intent: "schedule.query", body: queryBody }));
  pass("QA-02-deny", denied.status === 200 && deny.json?.State === "denied" && deniedSend.status === 403, `deny state ${deny.json?.State}, send ${deniedSend.status}`);
  const decided = await call("POST", "/v1/invite-decision", b.ownerCredential, { agent: a.agent, target: b.agent, decision: "accept" });
  const again = await call("POST", "/v1/invites", a.credential, { agent: a.agent, target: b.agent });
  pass("QA-02-active-reinvite", decided.json?.State === "active" && again.json?.State === "active" && again.json?.Generation === decided.json?.Generation, `accept gen ${decided.json?.Generation}, reinvite ${again.json?.State} gen ${again.json?.Generation}`);
  const inviteC = await call("POST", "/v1/invites", a.credential, { agent: a.agent, target: c.agent });
  const activeWhilePending = await call("GET", "/v1/contacts", c.credential);
  const body = { agent: a.agent, target: c.agent, decision: "accept" };
  const codes = await Promise.all(Array.from({ length: 8 }, () => call("POST", "/v1/invite-decision", c.ownerCredential, body)));
  const contactsC = await call("GET", "/v1/contacts", c.credential);
  pass("QA-02-pending-concurrent", inviteC.json?.State === "pending" && activeWhilePending.json?.length === 0 && codes.every((item) => item.status === 200) && contactsC.json?.length === 1 && contactsC.json[0] === a.agent, `pending-contacts ${activeWhilePending.json?.length}, accept ${codes.map((item) => item.status).join(",")}, contacts ${JSON.stringify(contactsC.json)}`);
  held("QA-02-free-n", "Free N, price, and slot-unit stay unset. This run does not estimate them.");

  const live = build(a, { to: b.agent, intent: "schedule.query", body: queryBody });
  const accepted = await call("POST", "/v1/send", a.credential, live);
  const anon = await call("POST", "/v1/send", "", live);
  pass("QA-03-auth-before-cache", accepted.status === 200 && anon.status === 401 && !anon.text.includes(accepted.json?.digest || "missing-digest"), `send ${accepted.status}, unauth ${anon.status}, digest leaked ${anon.text.includes(accepted.json?.digest || "x")}`);
  const raw = JSON.stringify(live);
  const mutated = [
    ["duplicate-key", raw.replace('"v":"relay.v1"', '"v":"relay.v1","v":"relay.v1"')],
    ["deliver-both", raw.replace('"deliver":"agent"', '"deliver":"both"')],
    ["unknown", raw.replace('"body":', '"unknown":1,"body":')],
    ["taskId", raw.replace('"body":', '"taskId":"t","body":')],
    ["role-user", raw.replace('"body":', '"role":"user","body":')],
    ["bad-unicode", raw.replace('"body":', '"render":{"hint":"\\ud800"},"body":')],
    ["uppercase-from", raw.replace(`"from":"${a.agent}"`, `"from":"${a.agent.toUpperCase()}"`)],
  ];
  const schema = [];
  for (const [name, value] of mutated) schema.push(`${name}:${(await call("POST", "/v1/send", a.credential, value)).status}`);
  pass("QA-03-strict-schema", schema.every((item) => !item.endsWith(":200")), schema.join(","));

  const replay = await call("POST", "/v1/send", a.credential, live);
  const conflictBody = { ...queryBody, granularity_min: 31 };
  const conflict = await call("POST", "/v1/send", a.credential, build(a, { to: b.agent, intent: "schedule.query", body: conflictBody, key: live.idempotency_key }));
  const tooLongKey = randomBytes(18).toString("base64url");
  const tooLong = await call("POST", "/v1/send", a.credential, build(a, { to: b.agent, intent: "schedule.query", body: queryBody, key: tooLongKey, exp: plus(301000) }));
  const rolledMessage = build(a, { to: b.agent, intent: "schedule.query", body: queryBody, key: tooLongKey });
  const rolled = await call("POST", "/v1/send", a.credential, rolledMessage);
  const collision = await call("POST", "/v1/send", a.credential, build(a, { id: live.id, to: b.agent, intent: "schedule.query", body: queryBody }));
  psql(`UPDATE relay_state SET data = jsonb_set(data, '{Messages,${qid(rolledMessage.id)},Receipt,exp}', '"2020-01-01T00:00:00Z"')`);
  await call("GET", `/v1/receipts/${rolledMessage.id}`, a.credential);
  pass("QA-04-replay-rollback", replay.status === 200 && replay.json?.id === live.id && replay.json?.envelope === undefined && conflict.status === 409 && conflict.json?.error === "idempotency_conflict" && tooLong.status === 422 && rolled.status === 200 && collision.status === 409 && collision.json?.error === "id_collision", `replay ${replay.status} id-match ${replay.json?.id === live.id}, conflict ${conflict.status} ${conflict.json?.error}, ttl ${tooLong.status} ${tooLong.json?.error}, rollback ${rolled.status}, id ${collision.status} ${collision.json?.error}`);

  const lease = await call("POST", "/v1/pull", b.credential, {});
  const until = Date.parse(lease.json?.lease_until);
  const leaseGap = (until - Date.now()) / 1000;
  pass("QA-05-lease-30s", lease.status === 200 && lease.json?.attempts === 1 && leaseGap > 25 && leaseGap <= 31, `attempts ${lease.json?.attempts}, lease gap ${leaseGap.toFixed(1)}s`);
  const short = build(a, { to: b.agent, intent: "schedule.query", body: queryBody, exp: plus(10000) });
  const shortSend = await call("POST", "/v1/send", a.credential, short);
  const earlyAck = await call("POST", "/v1/ack", b.credential, { id: short.id, token: lease.json.lease_token });
  const shortLease = await call("POST", "/v1/pull", b.credential, {});
  const shortGap = (Date.parse(shortLease.json?.lease_until) - Date.parse(short.exp)) / 1000;
  pass("QA-05-exp-min-and-ack-before-persist", shortSend.status === 200 && earlyAck.status === 409 && shortLease.json?.envelope?.id === short.id && shortLease.json?.attempts === 1 && shortGap <= 1, `ack-before-persist ${earlyAck.status}, short attempts ${shortLease.json?.attempts}, until-exp ${shortGap.toFixed(1)}s`);
  psql(`UPDATE relay_state SET data = jsonb_set(data, '{Messages,${qid(live.id)},LeaseUntil}', '"2020-01-01T00:00:00Z"')`);
  const second = await call("POST", "/v1/pull", b.credential, {});
  const late = await call("POST", "/v1/ack", b.credential, { id: live.id, token: lease.json.lease_token });
  psql(`UPDATE relay_state SET data = jsonb_set(data, '{Messages,${qid(live.id)},LeaseUntil}', '"2020-01-01T00:00:00Z"')`);
  const third = await call("POST", "/v1/pull", b.credential, {});
  const persisted = await call("POST", "/v1/persist", b.credential, { id: live.id, token: third.json?.lease_token });
  const acked = await call("POST", "/v1/ack", b.credential, { id: live.id, token: third.json?.lease_token });
  const receipt = await call("GET", `/v1/receipts/${live.id}`, a.credential);
  const claims = await Promise.all(Array.from({ length: 8 }, () => call("POST", "/v1/claim", b.credential, { id: live.id })));
  const winners = claims.filter((item) => item.status === 200).length;
  execFileSync("docker", ["compose", "--project-name", process.env.COMPOSE_PROJECT, "--env-file", ".env.example", "-f", "compose.yaml", "-f", process.env.COMPOSE_OVERRIDE, "restart", "relay"], { stdio: "inherit" });
  for (let i = 0; i < 30; i++) {
    const health = await call("GET", "/healthz").catch(() => ({ status: 0 }));
    if (health.status === 200) break;
    await new Promise((resolve) => setTimeout(resolve, 1000));
  }
  const afterRestart = await call("POST", "/v1/claim", b.credential, { id: live.id });
  pass("QA-05-attempts-ack-claim-restart", second.json?.attempts === 2 && late.status === 409 && third.json?.attempts === 3 && persisted.status === 200 && acked.status === 200 && receipt.json?.receipt?.transport === "delivered" && winners === 1 && afterRestart.status === 409, `attempts ${second.json?.attempts}/${third.json?.attempts}, late ${late.status}, delivered ${receipt.json?.receipt?.transport}, winners ${winners}, restart ${afterRestart.status}`);
  psql(`UPDATE relay_state SET data = jsonb_set(data, '{Messages,${qid(short.id)},LeaseUntil}', '"2020-01-01T00:00:00Z"')`);
  for (let i = 0; i < 3; i++) {
    const next = await call("POST", "/v1/pull", b.credential, {});
    if (next.json?.envelope?.id === short.id) psql(`UPDATE relay_state SET data = jsonb_set(data, '{Messages,${qid(short.id)},LeaseUntil}', '"2020-01-01T00:00:00Z"')`);
  }
  const exhausted = await call("GET", `/v1/receipts/${short.id}`, a.credential);
  pass("QA-05-max-attempts", exhausted.json?.receipt?.transport === "failed:max_attempts", `transport ${exhausted.json?.receipt?.transport}`);

  const queued = build(a, { to: b.agent, intent: "schedule.query", body: queryBody });
  await call("POST", "/v1/send", a.credential, queued);
  const revoked = await call("POST", "/v1/key-revoke", b.ownerCredential, { agent: b.agent, kid: "key1" });
  const leasedAfter = await call("POST", "/v1/pull", b.credential, {});
  const freshKey = generateKeyPairSync("ed25519");
  const freshPublic = freshKey.publicKey.export({ format: "der", type: "spki" }).subarray(-32).toString("base64url");
  const restored = await call("POST", "/v1/keys", b.ownerCredential, { agent: b.agent, kid: "key2", public: freshPublic, proof: sign(null, Buffer.from(["KNOWSLINK-KEY-POP", b.owner, b.agent, "key2", freshPublic].join("\0")), freshKey.privateKey).toString("base64url") });
  const pullOld = await call("POST", "/v1/pull", b.credential, {});
  b.privateKey = freshKey.privateKey;
  b.kid = "key2";
  const fresh = build(a, { to: b.agent, intent: "schedule.query", body: queryBody });
  const freshSend = await call("POST", "/v1/send", a.credential, fresh);
  const pullFresh = await call("POST", "/v1/pull", b.credential, {});
  pass("QA-06-revoke-before-lease", revoked.status === 200 && leasedAfter.status === 200 && !leasedAfter.json?.lease_token && restored.status === 200 && !pullOld.json?.lease_token && freshSend.status === 200 && pullFresh.json?.envelope?.id === fresh.id && Boolean(pullFresh.json?.lease_token), `revoke ${revoked.status}, pull-revoked ${leasedAfter.status} token ${Boolean(leasedAfter.json?.lease_token)}, rotate ${restored.status}, old-not-restored ${Boolean(pullOld.json?.lease_token)}, fresh ${freshSend.status} leased ${pullFresh.json?.envelope?.id === fresh.id}`);
  const token = pullFresh.json?.lease_token;
  const targetId = pullFresh.json?.envelope?.id;
  await call("POST", "/v1/unpair", b.ownerCredential, { agent: a.agent, target: b.agent });
  const ackUnpaired = await call("POST", "/v1/ack", b.credential, { id: targetId, token });
  const reinvite = await call("POST", "/v1/invites", a.credential, { agent: a.agent, target: b.agent });
  const reaccept = await call("POST", "/v1/invite-decision", b.ownerCredential, { agent: a.agent, target: b.agent, decision: "accept" });
  const oldReplay = await call("POST", "/v1/send", a.credential, queued);
  const newSend = await call("POST", "/v1/send", a.credential, build(a, { to: b.agent, intent: "schedule.query", body: queryBody }));
  pass("QA-06-unpair-generation", ackUnpaired.status === 409 && reinvite.json?.State === "pending" && reinvite.json?.Generation > 1 && reaccept.json?.State === "active" && oldReplay.status === 403 && newSend.status === 200, `ack ${ackUnpaired.status}, generation ${reinvite.json?.Generation}, old ${oldReplay.status}, new ${newSend.status}`);
  psql(`UPDATE relay_state SET clock = clock_timestamp() + interval '1 hour'`);
  const skewed = await call("POST", "/v1/pull", b.credential, {});
  psql(`UPDATE relay_state SET clock = clock_timestamp()`);
  const recovered = await call("GET", "/healthz");
  pass("QA-06-abnormal-clock", skewed.status === 503 && skewed.json?.error === "abnormal_clock" && recovered.status === 200, `clock ${skewed.status} ${skewed.json?.error}, health ${recovered.status}`);
  const unknown = await call("POST", "/v1/pull", "not-a-credential", {});
  pass("QA-06-unknown-auth", unknown.status === 401, `unknown credential ${unknown.status}`);
  held("QA-06-epoch-cas", "Concurrent accepts and claims showed one winner. A deliberate stale-epoch 503 was not forced from outside the transaction.");

  await drain(b);
  const parent = build(a, { to: b.agent, intent: "schedule.query", body: queryBody, render: { hint: "이 힌트는 승인 근거가 아니다" } });
  const parentSend = await call("POST", "/v1/send", a.credential, parent);
  const parentLease = await call("POST", "/v1/pull", b.credential, {});
  const parentPersist = await call("POST", "/v1/persist", b.credential, { id: parent.id, token: parentLease.json?.lease_token });
  const beforeAck = await call("GET", `/v1/receipts/${parent.id}`, a.credential);
  const parentAck = await call("POST", "/v1/ack", b.credential, { id: parent.id, token: parentLease.json?.lease_token });
  const noClaim = await call("POST", "/v1/authorize", b.credential, { id: parent.id, claim: "" });
  const parentClaim = await call("POST", "/v1/claim", b.credential, { id: parent.id });
  const policy = await call("POST", "/v1/authorize", b.credential, { id: parent.id, claim: parentClaim.json?.claim });
  pass("QA-07-query-deny-and-claim", parentSend.status === 200 && parentLease.json?.envelope?.id === parent.id && parentPersist.status === 200 && beforeAck.json?.receipt?.transport !== "delivered" && parentAck.status === 200 && noClaim.status === 403 && policy.json?.disclosure === false && policy.json?.executable === false && String(policy.json?.policy).includes("disclosure policy absent"), `pre-ack ${beforeAck.json?.receipt?.transport}, persist ${parentPersist.status}, authorize ${policy.json?.executable}/${policy.json?.disclosure}`);
  const commit = build(a, { to: b.agent, intent: "schedule.commit", body: { slot: queryBody.window, timezone: "UTC", commitment: "synthetic-stub" } });
  await call("POST", "/v1/send", a.credential, commit);
  const commitLease = await call("POST", "/v1/pull", b.credential, {});
  await call("POST", "/v1/persist", b.credential, { id: commit.id, token: commitLease.json?.lease_token });
  await call("POST", "/v1/ack", b.credential, { id: commit.id, token: commitLease.json?.lease_token });
  const commitClaim = await call("POST", "/v1/claim", b.credential, { id: commit.id });
  const commitAuth = await call("POST", "/v1/authorize", b.credential, { id: commit.id, claim: commitClaim.json?.claim });
  const commitReceipt = (await call("GET", `/v1/receipts/${commit.id}`, b.credential)).json?.receipt;
  const commitDone = build(b, { to: a.agent, intent: "relay.result", body: { status: "done" }, reply_to: commit.id, key: randomBytes(18).toString("base64url"), exp: commitReceipt?.exp });
  const doneSend = await call("POST", "/v1/send", b.credential, commitDone, commitClaim.json?.claim);
  pass("QA-07-commit-stub", String(commitClaim.json?.policy).includes("non-executable") && commitAuth.json?.executable === false && doneSend.status === 403, `policy ${commitClaim.json?.policy}, executable ${commitAuth.json?.executable}, done ${doneSend.status}`);
  const parentReceipt = (await call("GET", `/v1/receipts/${parent.id}`, b.credential)).json?.receipt;
  const digest = parentReceipt?.digest;
  const unbound = build(b, { to: b.agent, intent: "relay.approval.request", body: { reason: "judgment_required", request_digest: digest }, reply_to: parent.id, exp: parentReceipt?.exp });
  const noClaimH = await call("POST", "/v1/send", b.credential, unbound);
  const bound = build(b, { to: b.agent, intent: "relay.approval.request", body: { reason: "judgment_required", request_digest: digest }, reply_to: parent.id, exp: parentReceipt?.exp });
  const withClaim = await call("POST", "/v1/send", b.credential, bound, parentClaim.json?.claim);
  const duplicate = build(b, { to: b.agent, intent: "relay.approval.request", body: { reason: "judgment_required", request_digest: digest }, reply_to: parent.id, exp: parentReceipt?.exp });
  const duplicateSend = await call("POST", "/v1/send", b.credential, duplicate, parentClaim.json?.claim);
  const nested = build(b, { to: b.agent, intent: "relay.approval.request", body: { reason: "judgment_required", request_digest: digest }, reply_to: bound.id });
  const nestedSend = await call("POST", "/v1/send", b.credential, nested, parentClaim.json?.claim);
  const commitDigest = commitReceipt?.digest;
  const lateH = build(b, { to: b.agent, intent: "relay.approval.request", body: { reason: "permission_required", request_digest: commitDigest }, reply_to: commit.id, exp: plus(240000) });
  const lateSend = await call("POST", "/v1/send", b.credential, lateH, commitClaim.json?.claim);
  pass("QA-07-escalation-binding", noClaimH.status === 403 && withClaim.status === 200 && duplicateSend.status === 409 && nestedSend.status === 403 && lateSend.status === 422, `no-claim ${noClaimH.status}, bound ${withClaim.status}, duplicate ${duplicateSend.status} ${duplicateSend.json?.error}, recursive ${nestedSend.status}, late-exp ${lateSend.status} ${lateSend.json?.error}`);

  const gatePage = await call("GET", `/owner/gates/${bound.id}`, b.ownerCredential);
  const csrf = /name="csrf" value="([^"]+)"/.exec(gatePage.text)?.[1];
  const hinted = gatePage.text.includes("이 힌트는 승인 근거가 아니다");
  const getDecision = await call("GET", `/owner/gates/${bound.id}?decision=approve`, b.ownerCredential);
  const agentPage = await call("GET", `/owner/gates/${bound.id}`, b.credential);
  const basic = await call("GET", `/owner/gates/${bound.id}`, b.ownerCredential, undefined, undefined, "basic");
  const badCsrf = await fetch(`${base}/owner/gates/${bound.id}`, { method: "POST", redirect: "manual", headers: { Authorization: `Bearer ${b.ownerCredential}`, "Content-Type": "application/x-www-form-urlencoded" }, body: new URLSearchParams({ csrf: "bad", decision: "approve" }) });
  const approved = await fetch(`${base}/owner/gates/${bound.id}`, { method: "POST", redirect: "manual", headers: { Authorization: `Bearer ${b.ownerCredential}`, "Content-Type": "application/x-www-form-urlencoded" }, body: new URLSearchParams({ csrf, decision: "approve" }) });
  const approvedPage = await call("GET", `/owner/gates/${bound.id}`, b.ownerCredential);
  const secondDecision = await fetch(`${base}/owner/gates/${bound.id}`, { method: "POST", redirect: "manual", headers: { Authorization: `Bearer ${b.ownerCredential}`, "Content-Type": "application/x-www-form-urlencoded" }, body: new URLSearchParams({ csrf, decision: "deny" }) });
  const consumed = await call("POST", "/v1/gate-consume", b.credential, { id: bound.id, claim: parentClaim.json?.claim });
  const wrongConsume = await call("POST", "/v1/gate-consume", b.credential, { id: bound.id, claim: commitClaim.json?.claim });
  pass("QA-08-ui-csrf-consume", gatePage.status === 200 && gatePage.text.includes("granularity_min") && gatePage.text.includes("disclosure policy absent") && gatePage.text.includes("상태: pending") && !hinted && getDecision.text.includes("상태: pending") && agentPage.status === 401 && basic.status === 200 && badCsrf.status === 403 && approved.status === 303 && approvedPage.text.includes("상태: approved") && !approvedPage.text.includes('value="approve"') && secondDecision.status === 409 && consumed.json?.executable === false && consumed.json?.disclosure === false && wrongConsume.status === 409, `hint-shown ${hinted}, get ${getDecision.text.includes("pending")}, agent ${agentPage.status}, csrf ${badCsrf.status}, approve ${approved.status}, consume ${consumed.status} exec ${consumed.json?.executable}`);
  await save("v01-pending", gatePage.text);
  await save("v02-approved", approvedPage.text);
  const denyParent = build(a, { to: b.agent, intent: "schedule.query", body: queryBody });
  await call("POST", "/v1/send", a.credential, denyParent);
  const denyLease = await call("POST", "/v1/pull", b.credential, {});
  await call("POST", "/v1/persist", b.credential, { id: denyParent.id, token: denyLease.json?.lease_token });
  await call("POST", "/v1/ack", b.credential, { id: denyParent.id, token: denyLease.json?.lease_token });
  const denyClaim = await call("POST", "/v1/claim", b.credential, { id: denyParent.id });
  const denyReceipt = (await call("GET", `/v1/receipts/${denyParent.id}`, b.credential)).json?.receipt;
  const denyDigest = denyReceipt?.digest;
  const denyH = build(b, { to: b.agent, intent: "relay.approval.request", body: { reason: "judgment_required", request_digest: denyDigest }, reply_to: denyParent.id, exp: denyReceipt?.exp });
  const denyGate = await call("POST", "/v1/send", b.credential, denyH, denyClaim.json?.claim);
  const denyHtml = await call("GET", `/owner/gates/${denyH.id}`, b.ownerCredential);
  const denyCsrf = /name="csrf" value="([^"]+)"/.exec(denyHtml.text)?.[1];
  await fetch(`${base}/owner/gates/${denyH.id}`, { method: "POST", redirect: "manual", headers: { Authorization: `Bearer ${b.ownerCredential}`, "Content-Type": "application/x-www-form-urlencoded" }, body: new URLSearchParams({ csrf: denyCsrf, decision: "deny" }) });
  const deniedPage = await call("GET", `/owner/gates/${denyH.id}`, b.ownerCredential);
  pass("QA-08-deny-page", denyGate.status === 200 && deniedPage.text.includes("상태: denied") && !deniedPage.text.includes('value="approve"'), `deny-send ${denyGate.status} ${denyGate.json?.error || ""}, denied ${deniedPage.text.includes("상태: denied")}, button ${deniedPage.text.includes('value="approve"')}`);
  await save("v02-denied", deniedPage.text);

  const commitExp = commitReceipt?.exp;
  const wrongWay = build(b, { to: b.agent, intent: "relay.result", body: { status: "denied" }, reply_to: commit.id, exp: commitExp });
  const wrongSend = await call("POST", "/v1/send", b.credential, wrongWay, commitClaim.json?.claim);
  const leaked = build(b, { to: a.agent, intent: "relay.result", body: { status: "denied", result: { title: "private" } }, reply_to: commit.id, exp: commitExp });
  const leakedSend = await call("POST", "/v1/send", b.credential, leaked, commitClaim.json?.claim);
  const traced = JSON.stringify(build(b, { to: a.agent, intent: "relay.result", body: { status: "denied" }, reply_to: commit.id, exp: commitExp })).replace('"status":"denied"', '"status":"denied","error":{"stack":"trace"}');
  const tracedSend = await call("POST", "/v1/send", b.credential, traced, commitClaim.json?.claim);
  const validResult = build(b, { to: a.agent, intent: "relay.result", body: { status: "denied" }, reply_to: commit.id, exp: commitExp });
  const validSend = await call("POST", "/v1/send", b.credential, validResult, commitClaim.json?.claim);
  const againResult = build(b, { to: a.agent, intent: "relay.result", body: { status: "failed" }, reply_to: commit.id, exp: commitExp });
  const againSend = await call("POST", "/v1/send", b.credential, againResult, commitClaim.json?.claim);
  const completion = await call("GET", `/v1/receipts/${commit.id}`, a.credential);
  pass("QA-09-result-binding", wrongSend.status === 403 && leakedSend.status === 422 && tracedSend.status === 422 && validSend.status === 200 && (againSend.status === 403 || againSend.status === 409) && completion.json?.completion === "denied", `direction ${wrongSend.status}, result ${leakedSend.status}, stack ${tracedSend.status}, denied ${validSend.status}, duplicate ${againSend.status} ${againSend.json?.error}, completion ${completion.json?.completion}`);
  const resultPull = await call("POST", "/v1/pull", a.credential, {});
  const resultPersist = await call("POST", "/v1/persist", a.credential, { id: validResult.id, token: resultPull.json?.lease_token });
  const resultAck = await call("POST", "/v1/ack", a.credential, { id: validResult.id, token: resultPull.json?.lease_token });
  const afterResult = await call("POST", "/v1/pull", a.credential, {});
  pass("QA-09-result-not-command", resultPull.json?.envelope?.intent === "relay.result" && resultPersist.status === 200 && resultAck.status === 200 && !afterResult.json?.lease_token, `intent ${resultPull.json?.envelope?.intent}, persist ${resultPersist.status}, ack ${resultAck.status}, next ${Boolean(afterResult.json?.lease_token)}`);

  await drain(b);
  const normal = build(a, { to: b.agent, intent: "schedule.query", body: queryBody });
  await call("POST", "/v1/send", a.credential, normal);
  const evidence = build(a, { to: b.agent, intent: "schedule.query", body: queryBody, evidence: [{ ref: "urn:synthetic:qa-evidence" }], priority: "high" });
  const evidenceSend = await call("POST", "/v1/send", a.credential, evidence);
  const ordered = await call("POST", "/v1/pull", b.credential, {});
  const webhook = await call("POST", "/webhook", "", {});
  const registry = await call("GET", "/v1/registry");
  const manifest = registry.json?.manifest || "";
  pass("QA-10-off-boundaries", evidenceSend.status === 200 && ordered.json?.envelope?.id === normal.id && webhook.status === 404 && registry.status === 200 && !manifest.includes("taskId") && !manifest.includes("AgentCard") && !/webhook|evidenceFetch/.test(manifest), `evidence ${evidenceSend.status}, older-normal-not-preempted ${ordered.json?.envelope?.id === normal.id}, webhook ${webhook.status}`);
  const unconfigured = execFileSync("node", ["adapters/dist/index.js"], { encoding: "utf8" });
  pass("QA-10-adapter-unconfigured", unconfigured.includes('"webhook":false') && unconfigured.includes('"evidenceFetch":false'), unconfigured.trim());
  held("QA-10-real-adapter", "Loopback TypeScript stub only. No vendor connection or public operation is claimed.");
  held("QA-10-a2a", "A2A wire fields were rejected. No claim of a current A2A revision review or push webhook.");

  const kept = build(a, { to: b.agent, intent: "schedule.query", body: queryBody });
  await call("POST", "/v1/send", a.credential, kept);
  psql(`UPDATE relay_state SET data = jsonb_set(data, '{Messages,${qid(kept.id)},Receipt,exp}', '"2020-01-01T00:00:00Z"')`);
  const expiredReceipt = await call("GET", `/v1/receipts/${kept.id}`, a.credential);
  const envelopeLeft = psql(`SELECT data->'Messages'->'${qid(kept.id)}' ? 'Envelope' FROM relay_state`);
  const retained = psql(`SELECT data->'Messages' ? '${qid(kept.id)}' FROM relay_state`);
  pass("QA-11-exp-clears-original", expiredReceipt.json?.receipt?.transport === "failed:expired" && expiredReceipt.json?.completion === "" && envelopeLeft === "f" && retained === "t", `transport ${expiredReceipt.json?.receipt?.transport}, completion ${JSON.stringify(expiredReceipt.json?.completion)}, envelope ${envelopeLeft}, row ${retained}`);
  psql(`UPDATE relay_state SET data = jsonb_set(data, '{Messages,${qid(kept.id)},Receipt,accepted_at}', '"2020-01-01T00:00:00Z"')`);
  await call("GET", "/v1/contacts", a.credential);
  const gone = psql(`SELECT data->'Messages' ? '${qid(kept.id)}' FROM relay_state`);
  const young = build(a, { to: b.agent, intent: "schedule.query", body: queryBody });
  await call("POST", "/v1/send", a.credential, young);
  psql(`UPDATE relay_state SET data = jsonb_set(data, '{Messages,${qid(young.id)},Receipt,accepted_at}', to_jsonb(to_char(clock_timestamp() AT TIME ZONE 'UTC' - interval '23 hours', 'YYYY-MM-DD"T"HH24:MI:SS"Z"')))`);
  await call("GET", "/v1/contacts", a.credential);
  const youngLeft = psql(`SELECT data->'Messages'->'${qid(young.id)}' ? 'Envelope' FROM relay_state`);
  pass("QA-11-receipt-24h", gone === "f" && youngLeft === "t", `25h message present ${gone}, 23h envelope ${youngLeft}; wall clock was not waited`);
  held("QA-11-wal", "Row removal is not a claim that WAL or backups were erased.");

  const expireGate = build(a, { to: b.agent, intent: "schedule.query", body: queryBody });
  await call("POST", "/v1/send", a.credential, expireGate);
  let expireLease = await call("POST", "/v1/pull", b.credential, {});
  for (let i = 0; i < 15 && expireLease.json?.envelope && expireLease.json.envelope.id !== expireGate.id; i++) {
    await call("POST", "/v1/persist", b.credential, { id: expireLease.json.envelope.id, token: expireLease.json.lease_token });
    await call("POST", "/v1/ack", b.credential, { id: expireLease.json.envelope.id, token: expireLease.json.lease_token });
    expireLease = await call("POST", "/v1/pull", b.credential, {});
  }
  await call("POST", "/v1/persist", b.credential, { id: expireGate.id, token: expireLease.json?.lease_token });
  await call("POST", "/v1/ack", b.credential, { id: expireGate.id, token: expireLease.json?.lease_token });
  const expireClaim = await call("POST", "/v1/claim", b.credential, { id: expireGate.id });
  const expireReceipt = (await call("GET", `/v1/receipts/${expireGate.id}`, b.credential)).json?.receipt;
  const expireDigest = expireReceipt?.digest;
  const expireH = build(b, { to: b.agent, intent: "relay.approval.request", body: { reason: "judgment_required", request_digest: expireDigest }, reply_to: expireGate.id, exp: expireReceipt?.exp });
  await call("POST", "/v1/send", b.credential, expireH, expireClaim.json?.claim);
  psql(`UPDATE relay_state SET data = jsonb_set(data, '{Gates,${qid(expireH.id)},Exp}', '"2020-01-01T00:00:00Z"')`);
  const expiredPage = await call("GET", `/owner/gates/${expireH.id}`, b.ownerCredential);
  const missingParent = build(a, { to: b.agent, intent: "schedule.query", body: queryBody });
  await call("POST", "/v1/send", a.credential, missingParent);
  let missingLease = await call("POST", "/v1/pull", b.credential, {});
  for (let i = 0; i < 15 && missingLease.json?.envelope && missingLease.json.envelope.id !== missingParent.id; i++) {
    await call("POST", "/v1/persist", b.credential, { id: missingLease.json.envelope.id, token: missingLease.json.lease_token });
    await call("POST", "/v1/ack", b.credential, { id: missingLease.json.envelope.id, token: missingLease.json.lease_token });
    missingLease = await call("POST", "/v1/pull", b.credential, {});
  }
  await call("POST", "/v1/persist", b.credential, { id: missingParent.id, token: missingLease.json?.lease_token });
  await call("POST", "/v1/ack", b.credential, { id: missingParent.id, token: missingLease.json?.lease_token });
  const missingClaim = await call("POST", "/v1/claim", b.credential, { id: missingParent.id });
  const missingReceipt = (await call("GET", `/v1/receipts/${missingParent.id}`, b.credential)).json?.receipt;
  const missingDigest = missingReceipt?.digest;
  const missingH = build(b, { to: b.agent, intent: "relay.approval.request", body: { reason: "judgment_required", request_digest: missingDigest }, reply_to: missingParent.id, exp: missingReceipt?.exp });
  await call("POST", "/v1/send", b.credential, missingH, missingClaim.json?.claim);
  psql(`UPDATE relay_state SET data = data #- '{Messages,${qid(missingParent.id)},Envelope}'`);
  const unavailablePage = await call("GET", `/owner/gates/${missingH.id}`, b.ownerCredential);
  await drain(b);
  const revokeParent = build(a, { to: b.agent, intent: "schedule.query", body: queryBody });
  await call("POST", "/v1/send", a.credential, revokeParent);
  const revokeLease = await call("POST", "/v1/pull", b.credential, {});
  await call("POST", "/v1/persist", b.credential, { id: revokeParent.id, token: revokeLease.json?.lease_token });
  await call("POST", "/v1/ack", b.credential, { id: revokeParent.id, token: revokeLease.json?.lease_token });
  const revokeClaim = await call("POST", "/v1/claim", b.credential, { id: revokeParent.id });
  const revokeReceipt = (await call("GET", `/v1/receipts/${revokeParent.id}`, b.credential)).json?.receipt;
  const revokeDigest = revokeReceipt?.digest;
  const revokeH = build(b, { to: b.agent, intent: "relay.approval.request", body: { reason: "permission_required", request_digest: revokeDigest }, reply_to: revokeParent.id, exp: revokeReceipt?.exp });
  await call("POST", "/v1/send", b.credential, revokeH, revokeClaim.json?.claim);
  await call("POST", "/v1/unpair", a.ownerCredential, { agent: a.agent, target: b.agent });
  const revokedPage = await call("GET", `/owner/gates/${revokeH.id}`, b.ownerCredential);
  const unauthorized = await call("GET", `/owner/gates/${missingH.id}`);
  pass("QA-11-captions", expiredPage.text.includes("상태: expired") && !expiredPage.text.includes('value="approve"') && revokedPage.text.includes("상태: revoked") && !revokedPage.text.includes('value="approve"') && unavailablePage.text.includes("원문 부재") && unavailablePage.text.includes("상태: unavailable") && unauthorized.status === 401, `expired ${expiredPage.text.includes("상태: expired")}, revoked ${revokedPage.text.includes("상태: revoked")}, unavailable ${unavailablePage.text.includes("원문 부재")}, auth ${unauthorized.status}`);
  await save("v03-expired", expiredPage.text);
  await save("v03-revoked", revokedPage.text);
  await save("v04-unavailable", unavailablePage.text);
  await save("v04-unauthorized", `<!doctype html><meta charset="utf-8"><title>KnowsLink 승인</title><pre>${unauthorized.text.trim()}</pre>`);
  held("DEC-02", "Positive silent done on real data stays held. Synthetic status done was rejected and is not a silent-done pass.");
  held("DEC-03", "Public limits, real identity proof, and singleton throughput stay held. Synthetic pass does not replace them.");
  held("UI-designer", "V-01–04 captures are shared for designer visual review. This probe does not record that review as passed.");

  const out = ".fullops-squad/docs/evaluations/qa-reports/SAR-MVP-001-TESTER-test/probe-results.json";
  await writeFile(out, JSON.stringify({ checks, pages }, null, 2));
  const failed = checks.filter((item) => item.result === "fail");
  console.info(`PROBE fail ${failed.length} held ${checks.filter((item) => item.result === "held").length} pass ${checks.filter((item) => item.result === "pass").length}`);
  if (!parentPersist || !commitAuth) console.info("persist markers recorded");
  process.exit(failed.length ? 1 : 0);
}
main().catch((error) => {
  console.error(error instanceof Error ? error.message : "probe failed");
  process.exit(1);
});
