// One synthetic pull adapter verifies signatures, persists a shared inbox, ACKs, then claims before judgment.
import {
  createHash,
  createPublicKey,
  randomBytes,
  sign,
  verify,
} from "node:crypto";
import { readFile } from "node:fs/promises";

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
    protected readonly agent: string,
    protected readonly kid: string,
    protected readonly privateKey: string,
    private readonly accessHeaders: Record<string, string> = {},
    private readonly timeoutMs = 10000,
  ) {}
  async request<T>(path: string, body?: unknown, claim?: string): Promise<T> {
    // The timer holds the controller strongly; an inline AbortSignal.timeout() can be garbage-collected and never fire.
    const controller = new AbortController();
    const timer = setTimeout(
      () =>
        controller.abort(
          new DOMException("relay request timed out", "TimeoutError"),
        ),
      this.timeoutMs,
    );
    try {
      return await this.fetchJson<T>(path, controller.signal, body, claim);
    } finally {
      clearTimeout(timer);
    }
  }
  private async fetchJson<T>(
    path: string,
    signal: AbortSignal,
    body?: unknown,
    claim?: string,
  ): Promise<T> {
    const response = await fetch(`${this.base}${path}`, {
      redirect: "error",
      signal,
      method: body === undefined ? "GET" : "POST",
      headers: {
        ...this.accessHeaders,
        Authorization: `Bearer ${this.credential}`,
        ...(claim ? { "X-Execution-Claim": claim } : {}),
        "Content-Type": "application/json",
      },
      body: body === undefined ? undefined : JSON.stringify(body),
    });

    const reader = response.body?.getReader();
    if (!reader) throw new Error("empty relay response");
    // undici reaches the body through a WeakRef to Response, so abort alone may never fail a stalled read.
    const aborted = new Promise<never>((_resolve, reject) => {
      if (signal.aborted) reject(signal.reason);
      signal.addEventListener("abort", () => reject(signal.reason), {
        once: true,
      });
    });
    aborted.catch(() => undefined);
    const chunks: Uint8Array[] = [];
    let size = 0;
    try {
      while (true) {
        const chunk = await Promise.race([reader.read(), aborted]);
        if (chunk.done) break;
        size += chunk.value.byteLength;
        if (size > 65536) throw new Error("relay response too large");
        chunks.push(chunk.value);
      }
      const value: unknown = JSON.parse(Buffer.concat(chunks).toString("utf8"));
      if (!response.ok) {
        const error = value as { error?: unknown; retry_at?: unknown };
        const known = new Set([
          "invalid_auth",
          "sender_not_allowed",
          "human_invite_required",
          "invalid_signature",
          "expired",
          "ttl_too_long",
          "idempotency_conflict",
          "id_collision",
          "duplicate_result",
          "invalid_lease",
          "capacity",
          "rate_limited",
          "invalid_schema",
          "invalid_json",
          "unavailable",
        ]);
        const code =
          typeof error?.error === "string" && known.has(error.error)
            ? ` ${error.error}`
            : "";
        const retry =
          typeof error?.retry_at === "string" &&
          /^[0-9TZ:.-]+$/.test(error.retry_at)
            ? ` retry_at=${error.retry_at}`
            : "";
        throw new Error(`relay HTTP ${response.status}${code}${retry}`);
      }
      return value as T;
    } finally {
      await reader.cancel();
    }
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
  protected async verifyRegistry(): Promise<void> {
    const registry = await this.request<{ sha256: string; manifest: string }>(
      "/v1/registry",
    );
    if (
      registry.sha256 !==
        "2b25c6b58fb6b6973de1c9bac19162834cfc9b678a0416ad39e2cffab3cde1e2" ||
      createHash("sha256").update(registry.manifest).digest("hex") !==
        registry.sha256
    )
      throw new Error("invalid registry revision");
  }
  protected async verifyMessage(message: Envelope): Promise<void> {
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
  }
  async once(
    gate: boolean,
    onGate: (id: string) => void = (id) =>
      console.info(JSON.stringify({ state: "waiting", gate: id })),
  ): Promise<boolean> {
    await this.verifyRegistry();
    const lease = await this.request<{
      envelope: Envelope;
      lease_token: string;
    } | null>("/v1/pull", {});
    if (!lease) return false;
    const message = lease.envelope;
    await this.verifyMessage(message);
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
      onGate(approval.id);
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
export async function localAdapter(): Promise<Adapter | null> {
  const { RELAY_URL, AGENT_CREDENTIAL, AGENT_ID, AGENT_KID, AGENT_KEY_FILE } =
    process.env;
  if (
    !RELAY_URL ||
    !AGENT_CREDENTIAL ||
    !AGENT_ID ||
    !AGENT_KID ||
    !AGENT_KEY_FILE
  ) {
    return null;
  }
  const url = new URL(RELAY_URL);
  if (
    url.protocol !== "http:" ||
    !["127.0.0.1", "localhost", "[::1]"].includes(url.hostname) ||
    url.username ||
    url.password ||
    url.search ||
    url.hash ||
    url.pathname !== "/"
  )
    throw new Error("synthetic adapter requires loopback relay");
  return new Adapter(
    url.origin,
    AGENT_CREDENTIAL,
    AGENT_ID,
    AGENT_KID,
    await readFile(AGENT_KEY_FILE, "utf8"),
  );
}

export function relayBase(raw: string): string {
  const u = new URL(raw);
  if (
    u.username ||
    u.password ||
    u.pathname !== "/" ||
    u.search ||
    u.hash ||
    !(
      u.protocol === "https:" ||
      (u.protocol === "http:" &&
        ["localhost", "127.0.0.1", "[::1]"].includes(u.hostname))
    )
  )
    throw new Error("invalid relay URL");
  return u.origin;
}
