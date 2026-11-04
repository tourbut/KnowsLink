// Explicit trial-message transport for Codex and Grok; fixed HTTPS remote origin, paired peer, no business effects.
import { randomBytes, sign } from "node:crypto";
import { readFile } from "node:fs/promises";
import {
  AGENT,
  Adapter,
  relayBase,
  signingBytes,
  uuid7,
  validText,
  type Envelope,
} from "./core.js";

export type TrialMessage = {
  id: string;
  from: string;
  to: string;
  text: string;
  exp: string;
  untrusted: true;
};
export const trialMode = () =>
  ["test-loopback", "test-remote"].includes(process.env.KNOWSLINK_MODE ?? "");

export class TestTransport extends Adapter {
  constructor(
    base: string,
    credential: string,
    agent: string,
    kid: string,
    key: string,
    private readonly peer: string,
    headers: Record<string, string> = {},
    timeoutMs = 10000,
  ) {
    super(base, credential, agent, kid, key, headers, timeoutMs);
  }
  override async request<T>(path: string, body?: unknown): Promise<T> {
    return super.request<T>(path.replace(/^\/v1\//, "/v1/test/"), body);
  }
  async sendText(
    text: string,
    idempotencyKey: string,
  ): Promise<{ id: string }> {
    if (!validText(text)) throw new Error("invalid trial text");
    if (!/^[\x20-\x7E]{16,128}$/.test(idempotencyKey))
      throw new Error("invalid idempotency key");
    const envelope: Envelope = {
      v: "relay.v1",
      id: uuid7(),
      from: this.agent,
      to: this.peer,
      intent: "relay.test.message",
      body: { text },
      deliver: "agent",
      exp: new Date(Date.now() + 180000)
        .toISOString()
        .replace(/\.\d{3}Z$/, "Z"),
      idempotency_key: idempotencyKey,
      sig: { alg: "Ed25519", kid: this.kid, value: "" },
    };
    envelope.sig.value = sign(
      null,
      signingBytes(envelope),
      this.privateKey,
    ).toString("base64url");
    return this.request("/v1/send", envelope);
  }
  async receive(): Promise<TrialMessage | null> {
    await this.verifyRegistry();
    const lease = await this.request<{
      envelope: Envelope;
      lease_token: string;
    } | null>("/v1/pull", {});
    if (!lease) return null;
    const message = lease.envelope;
    if (
      message.intent !== "relay.test.message" ||
      message.from !== this.peer ||
      message.deliver !== "agent" ||
      message.reply_to !== undefined ||
      Object.keys(message.body).join() !== "text" ||
      typeof message.body.text !== "string" ||
      !validText(message.body.text) ||
      !Number.isFinite(Date.parse(message.exp)) ||
      Date.parse(message.exp) <= Date.now()
    )
      throw new Error("invalid trial message");
    await this.verifyMessage(message);
    const delivery = { id: message.id, token: lease.lease_token };
    await this.request("/v1/persist", delivery);
    await this.request("/v1/ack", delivery);
    await this.request("/v1/claim", { id: message.id });
    return {
      id: message.id,
      from: message.from,
      to: message.to,
      text: message.body.text,
      exp: message.exp,
      untrusted: true,
    };
  }
}

export async function testTransport(): Promise<TestTransport | null> {
  if (!trialMode()) return null;
  const {
    RELAY_URL,
    AGENT_CREDENTIAL,
    AGENT_ID,
    AGENT_KID,
    AGENT_KEY_FILE,
    KNOWSLINK_TEST_PEER,
    CF_ACCESS_CLIENT_ID,
    CF_ACCESS_CLIENT_SECRET,
  } = process.env;
  if (
    !RELAY_URL ||
    !AGENT_CREDENTIAL ||
    !AGENT_ID ||
    !AGENT_KID ||
    !AGENT_KEY_FILE ||
    !KNOWSLINK_TEST_PEER
  )
    return null;
  if (
    !AGENT.test(AGENT_ID) ||
    !AGENT.test(KNOWSLINK_TEST_PEER) ||
    AGENT_ID === KNOWSLINK_TEST_PEER ||
    !/^[A-Za-z0-9_:-]+$/.test(AGENT_KID)
  )
    throw new Error("invalid trial identity");
  const origin = relayBase(RELAY_URL);
  const headers: Record<string, string> = {};
  if (process.env.KNOWSLINK_MODE === "test-remote") {
    if (origin !== "https://link.knowslog.com")
      throw new Error("unapproved remote origin");
    if (!CF_ACCESS_CLIENT_ID || !CF_ACCESS_CLIENT_SECRET) return null;
    headers["CF-Access-Client-Id"] = CF_ACCESS_CLIENT_ID;
    headers["CF-Access-Client-Secret"] = CF_ACCESS_CLIENT_SECRET;
  } else if (!origin.startsWith("http:"))
    throw new Error("trial loopback required");
  return new TestTransport(
    origin,
    AGENT_CREDENTIAL,
    AGENT_ID,
    AGENT_KID,
    await readFile(AGENT_KEY_FILE, "utf8"),
    KNOWSLINK_TEST_PEER,
    headers,
  );
}
export const newTrialKey = () => randomBytes(24).toString("base64url");
