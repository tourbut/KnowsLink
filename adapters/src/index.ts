// One synthetic pull adapter verifies signatures, persists a shared inbox, ACKs, then claims before judgment.
import {
  createHash,
  createPublicKey,
  randomBytes,
  sign,
  verify,
} from "node:crypto";
import { readFile } from "node:fs/promises";
import { pathToFileURL } from "node:url";

type JSONValue =
  | null
  | boolean
  | number
  | string
  | JSONValue[]
  | { [key: string]: JSONValue };
export function canonical(value: JSONValue): string {
  if (value === null || typeof value !== "object") return JSON.stringify(value);
  if (Array.isArray(value)) return `[${value.map(canonical).join(",")}]`;
  return `{${Object.keys(value)
    .sort()
    .map((key) => `${JSON.stringify(key)}:${canonical(value[key]!)}`)
    .join(",")}}`;
}
export function uuid7(): string {
  const bytes = randomBytes(16);
  bytes.writeUIntBE(Date.now(), 0, 6);
  bytes[6] = (bytes[6]! & 15) | 112;
  bytes[8] = (bytes[8]! & 63) | 128;
  const hex = bytes.toString("hex");
  return `${hex.slice(0, 8)}-${hex.slice(8, 12)}-${hex.slice(12, 16)}-${hex.slice(16, 20)}-${hex.slice(20)}`;
}
export type Envelope = {
  v: string;
  id: string;
  from: string;
  to: string;
  intent: string;
  body: Record<string, JSONValue>;
  deliver: string;
  exp: string;
  idempotency_key: string;
  reply_to?: string;
  sig: { alg: string; kid: string; value: string };
};
export function signingBytes(envelope: Envelope): Buffer {
  const { sig, ...unsigned } = envelope;
  return Buffer.from(
    `SILENT-AGENT-RELAY\0relay.v1\0Ed25519\0${sig.kid}\0${canonical(unsigned)}`,
  );
}
export class Adapter {
  constructor(
    private readonly base: string,
    private readonly credential: string,
    private readonly agent: string,
    private readonly kid: string,
    private readonly privateKey: string,
  ) {}
  async request<T>(path: string, body?: unknown, claim?: string): Promise<T> {
    const response = await fetch(`${this.base}${path}`, {
      method: body === undefined ? "GET" : "POST",
      headers: {
        Authorization: `Bearer ${this.credential}`,
        ...(claim ? { "X-Execution-Claim": claim } : {}),
        "Content-Type": "application/json",
      },
      body: body === undefined ? undefined : JSON.stringify(body),
    });
    if (!response.ok) throw new Error(`relay HTTP ${response.status}`);
    return (await response.json()) as T;
  }
  async send(
    parent: Envelope,
    intent: string,
    body: Record<string, JSONValue>,
    claim: string,
  ): Promise<{ id: string }> {
    const envelope: Envelope = {
      v: "relay.v1",
      id: uuid7(),
      from: this.agent,
      to: intent === "relay.approval.request" ? this.agent : parent.from,
      intent,
      body,
      deliver: intent === "relay.approval.request" ? "human" : "agent",
      exp: parent.exp,
      idempotency_key: randomBytes(24).toString("base64url"),
      reply_to: parent.id,
      sig: { alg: "Ed25519", kid: this.kid, value: "" },
    };
    envelope.sig.value = sign(
      null,
      signingBytes(envelope),
      this.privateKey,
    ).toString("base64url");
    return this.request("/v1/send", envelope, claim);
  }
  async once(gate: boolean): Promise<boolean> {
    const registry = await this.request<{ sha256: string; manifest: string }>(
      "/v1/registry",
    );
    if (
      registry.sha256 !==
        "b9759a1ec4035281c704d23f35f78921d8b54d232fa4cd0f27002b183e5eafda" ||
      createHash("sha256").update(registry.manifest).digest("hex") !==
        registry.sha256
    )
      throw new Error("invalid registry revision");
    const lease = await this.request<{
      envelope: Envelope;
      lease_token: string;
    } | null>("/v1/pull", {});
    if (!lease) return false;
    const message = lease.envelope;
    if (
      message.to !== this.agent ||
      message.v !== "relay.v1" ||
      message.sig.alg !== "Ed25519"
    )
      throw new Error("invalid envelope");
    const { public: key } = await this.request<{ public: string }>(
      `/v1/keys/${encodeURIComponent(message.from)}/${encodeURIComponent(message.sig.kid)}`,
    );
    const publicKey = createPublicKey({
      key: Buffer.concat([
        Buffer.from("302a300506032b6570032100", "hex"),
        Buffer.from(key, "base64url"),
      ]),
      type: "spki",
      format: "der",
    });
    if (
      !verify(
        null,
        signingBytes(message),
        publicKey,
        Buffer.from(message.sig.value, "base64url"),
      )
    )
      throw new Error("invalid signature");
    const delivery = { id: message.id, token: lease.lease_token };
    // Relay's shared Postgres inbox stores the verified immutable bytes before ACK.
    await this.request("/v1/persist", delivery);
    await this.request("/v1/ack", delivery);
    if (message.intent === "relay.result") return true; // Never execute or auto-reply to control results.
    const claim = await this.request<{ claim: string }>("/v1/claim", {
      id: message.id,
    });
    if (gate) {
      const receipt = await this.request<{ receipt: { digest: string } }>(
        `/v1/receipts/${message.id}`,
      );
      const approval = await this.send(
        message,
        "relay.approval.request",
        {
          reason: "judgment_required",
          request_digest: receipt.receipt.digest,
        },
        claim.claim,
      );
      console.info(JSON.stringify({ state: "waiting", gate: approval.id }));
      while (Date.now() < Date.parse(message.exp)) {
        const result = await this.request<{ state: string }>(
          `/v1/gates/${approval.id}`,
        );
        if (result.state === "approved") {
          await this.request("/v1/gate-consume", {
            id: approval.id,
            claim: claim.claim,
          });
          break;
        }
        if (result.state !== "pending") break;
        await new Promise((resolve) => setTimeout(resolve, 2000));
      }
    }
    await this.request("/v1/authorize", { id: message.id, claim: claim.claim });
    // No disclosure policy and commit stub: approve cannot turn either into an effect.
    await this.send(message, "relay.result", { status: "denied" }, claim.claim);
    return true;
  }
}
async function main(): Promise<void> {
  const { RELAY_URL, AGENT_CREDENTIAL, AGENT_ID, AGENT_KID, AGENT_KEY_FILE } =
    process.env;
  if (
    !RELAY_URL ||
    !AGENT_CREDENTIAL ||
    !AGENT_ID ||
    !AGENT_KID ||
    !AGENT_KEY_FILE
  ) {
    console.info(
      JSON.stringify({
        state: "unconfigured",
        transport: "pull",
        webhook: false,
        evidenceFetch: false,
      }),
    );
    return;
  }
  const url = new URL(RELAY_URL);
  if (!["127.0.0.1", "localhost", "[::1]"].includes(url.hostname))
    throw new Error("synthetic adapter requires loopback relay");
  const adapter = new Adapter(
    RELAY_URL,
    AGENT_CREDENTIAL,
    AGENT_ID,
    AGENT_KID,
    await readFile(AGENT_KEY_FILE, "utf8"),
  );
  await adapter.once(process.env.ADAPTER_GATE === "1");
}
if (
  process.argv[1] &&
  import.meta.url === pathToFileURL(process.argv[1]).href
) {
  main().catch(() => {
    console.error("synthetic adapter failed");
    process.exitCode = 1;
  });
}
