// Explicit member connection-check CLI/MCP transport: private onboarding folder, signed short text, manual receive and no automatic reply.
import { createPublicKey, sign, verify } from "node:crypto";
import { readFile } from "node:fs/promises";
import { join, basename } from "node:path";
import { pathToFileURL } from "node:url";
import { Adapter, canonical, relayBase, uuid7, validText } from "./core.js";
import { privatePath } from "./private-files.js";
import { AutoReceiver, Inbox } from "./inbox.js";
import { GrokWake } from "./grok-wake.js";
import { claimSession, holdsSession, releaseSession } from "./session-lock.js";

export type TextEnvelope = {
  v: "knowslink.text.v1";
  id: string;
  from: string;
  to: string;
  text: string;
  exp: string;
  idempotency_key: string;
  reply_to?: string;
  sig: { alg: "Ed25519"; kid: string; value: string };
};
export type ReceivedText = {
  id: string;
  from: string;
  to: string;
  text: string;
  exp: string;
  reply_to?: string;
  untrusted: true;
  authorship?: "owner" | "agent";
};
// HP-05 public label of a friend's agent: friend-chosen names, shown only as data (O-03).
export type Contact = {
  agent: string;
  name: string;
  client: string;
  friend: string;
};
export type OwnerTextEnvelope = {
  v: "knowslink.owner-text.v1";
  id: string;
  from: string;
  to: string;
  text: string;
  exp: string;
  idempotency_key: string;
  authorship: "owner";
};
export function textSigningBytes(e: TextEnvelope): Buffer {
  return Buffer.from(
    `KNOWSLINK-TEXT\0knowslink.text.v1\0Ed25519\0${e.sig.kid}\0${canonical({ v: e.v, id: e.id, from: e.from, to: e.to, text: e.text, exp: e.exp, idempotency_key: e.idempotency_key, reply_to: e.reply_to ?? "" })}`,
  );
}
export class TextTransport extends Adapter {
  async listen(onSignal: () => void, signal: AbortSignal): Promise<void> {
    const response = await this.openEvents(signal);
    if (!response.ok) {
      await response.body?.cancel();
      throw new Error(`relay HTTP ${response.status}`);
    }
    if (
      !response.headers.get("content-type")?.startsWith("text/event-stream")
    ) {
      await response.body?.cancel();
      throw new Error("invalid arrival stream");
    }
    const reader = response.body?.getReader();
    if (!reader) throw new Error("empty arrival stream");
    const decoder = new TextDecoder("utf-8", { fatal: true });
    let buffer = "";
    let bytes = 0;
    const aborted = new Promise<never>((_resolve, reject) => {
      if (signal.aborted) reject(signal.reason);
      signal.addEventListener("abort", () => reject(signal.reason), {
        once: true,
      });
    });
    aborted.catch(() => undefined);
    try {
      while (true) {
        const chunk = await Promise.race([reader.read(), aborted]);
        if (chunk.done) return;
        bytes += chunk.value.byteLength;
        buffer += decoder.decode(chunk.value, { stream: true });
        if (bytes > 65536 || buffer.length > 4096)
          throw new Error("arrival stream too large");
        buffer = buffer.replace(/\r\n/g, "\n");
        let end: number;
        while ((end = buffer.indexOf("\n\n")) >= 0) {
          const frame = buffer.slice(0, end);
          buffer = buffer.slice(end + 2);
          if (/^event: (ready|arrival)\ndata: \{\}$/.test(frame)) onSignal();
          else if (!frame.startsWith(":"))
            throw new Error("invalid arrival signal");
        }
      }
    } finally {
      await reader.cancel();
    }
  }
  async sendText(
    peer: string,
    text: string,
    key: string,
    confirmed: boolean,
    replyTo?: string,
  ): Promise<{ id: string; transport: string }> {
    if (
      !confirmed ||
      !validText(text) ||
      !/^[\x20-\x7E]{16,128}$/.test(key) ||
      !/^[a-z][a-z0-9_:-]*$/.test(peer) ||
      peer === this.agent ||
      (replyTo !== undefined &&
        !/^[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/.test(
          replyTo,
        ))
    )
      throw new Error("invalid explicit text request");
    let expiry = Date.now() + 180000;
    if (replyTo) {
      const parent = await this.receipt(replyTo);
      if (
        parent.receipt.to !== this.agent ||
        parent.receipt.from !== peer ||
        parent.receipt.transport !== "delivered"
      )
        throw new Error("invalid reply parent");
      expiry = Math.min(expiry, Date.parse(parent.receipt.exp));
    }
    const e: TextEnvelope = {
      v: "knowslink.text.v1",
      id: uuid7(),
      from: this.agent,
      to: peer,
      text,
      exp: new Date(expiry).toISOString().replace(/\.\d{3}Z$/, "Z"),
      idempotency_key: key,
      sig: { alg: "Ed25519", kid: this.kid, value: "" },
    };
    if (replyTo) e.reply_to = replyTo;
    e.sig.value = sign(null, textSigningBytes(e), this.privateKey).toString(
      "base64url",
    );
    return this.request("/v1/text/send", e);
  }
  receipt(id: string): Promise<{
    receipt: {
      id: string;
      from: string;
      to: string;
      exp: string;
      transport: string;
    };
    completion: string;
    reply_id: string;
    reply_to: string;
  }> {
    if (!/^[0-9a-f-]{36}$/.test(id)) throw new Error("invalid receipt ID");
    return this.request(`/v1/receipts/${encodeURIComponent(id)}`);
  }
  async contact(agent: string): Promise<Contact> {
    if (!/^[a-z][a-z0-9_:-]*$/.test(agent)) throw new Error("invalid agent");
    const c = await this.request<Record<string, unknown>>(
      `/v1/contacts/${encodeURIComponent(agent)}`,
    );
    const field = (v: unknown) => (typeof v === "string" ? v.slice(0, 80) : "");
    return {
      agent,
      name: field(c.name),
      client: field(c.client),
      friend: field(c.friend),
    };
  }
  // Tells the sender's receipt that this work session showed the text (SP-12); no text leaves.
  seen(id: string): Promise<unknown> {
    return this.request("/v1/text/seen", { id });
  }
  // store runs after signature verification and before relay persist/ACK, so a local copy exists before the relay erases text.
  async receive(
    store?: (message: ReceivedText) => Promise<unknown>,
  ): Promise<ReceivedText | null> {
    const lease = await this.request<{
      envelope: TextEnvelope | OwnerTextEnvelope;
      lease_token: string;
    } | null>("/v1/text/pull", {});
    if (!lease) return null;
    const e = lease.envelope;
    const ownerWritten = e.v === "knowslink.owner-text.v1";
    const allowed = [
      "v",
      "id",
      "from",
      "to",
      "text",
      "exp",
      "idempotency_key",
      "reply_to",
      ...(ownerWritten ? ["authorship"] : ["sig"]),
    ];
    if (
      Object.keys(e).some((k) => !allowed.includes(k)) ||
      (e.v !== "knowslink.text.v1" && !ownerWritten) ||
      !/^[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/.test(
        e.id,
      ) ||
      !/^[a-z][a-z0-9_:-]*$/.test(e.from) ||
      e.to !== this.agent ||
      typeof e.text !== "string" ||
      !validText(e.text) ||
      (ownerWritten ? e.authorship !== "owner" : e.sig?.alg !== "Ed25519") ||
      !Number.isFinite(Date.parse(e.exp)) ||
      Date.parse(e.exp) <= Date.now()
    )
      throw new Error("invalid text envelope");
    // Owner text comes only from the authenticated relay inbox; it carries no agent signature.
    if (e.v === "knowslink.text.v1") {
      const key = await this.request<{ public: string }>(
        `/v1/keys/${encodeURIComponent(e.from)}/${encodeURIComponent(e.sig.kid)}`,
      );
      const publicKey = createPublicKey({
        key: Buffer.concat([
          Buffer.from("302a300506032b6570032100", "hex"),
          Buffer.from(key.public, "base64url"),
        ]),
        type: "spki",
        format: "der",
      });
      if (
        !verify(
          null,
          textSigningBytes(e),
          publicKey,
          Buffer.from(e.sig.value, "base64url"),
        )
      )
        throw new Error("invalid text signature");
    }
    const result: ReceivedText = {
      id: e.id,
      from: e.from,
      to: e.to,
      text: e.text,
      exp: e.exp,
      untrusted: true,
      authorship: ownerWritten ? "owner" : "agent",
    };
    if (e.v === "knowslink.text.v1" && e.reply_to) result.reply_to = e.reply_to;
    if (store) await store(result);
    const delivery = { id: e.id, token: lease.lease_token };
    await this.request("/v1/text/persist", delivery);
    await this.request("/v1/text/ack", delivery);
    return result;
  }
}
export async function memberTransport(folder: string): Promise<TextTransport> {
  for (const name of ["", "agent.json", "private.pem"]) {
    await privatePath(join(folder, name), name === "");
  }
  const c: { relay: string; agent: string; kid: string; credential: string } =
    JSON.parse(await readFile(join(folder, "agent.json"), "utf8"));
  if (
    !/^[a-z][a-z0-9_:-]*$/.test(c.agent) ||
    !/^[A-Za-z0-9_:-]{1,128}$/.test(c.kid) ||
    !/^[A-Za-z0-9_-]{43}$/.test(c.credential)
  )
    throw new Error("invalid member connection");
  return new TextTransport(
    relayBase(c.relay),
    c.credential,
    c.agent,
    c.kid,
    await readFile(join(folder, "private.pem"), "utf8"),
  );
}
export async function configuredMemberTransport(): Promise<TextTransport | null> {
  if (
    process.env.KNOWSLINK_MODE !== "public-node" ||
    !process.env.KNOWSLINK_AGENT_FOLDER
  )
    return null;
  return memberTransport(process.env.KNOWSLINK_AGENT_FOLDER);
}
// Always-on receive into the same private inbox the MCP server reads; prints metadata only, never text.
// With a session alias the watcher is one work session's receiver: it binds the agent folder (SP-02) and ends once another session takes it (SP-09).
export async function watch(
  folder: string,
  session?: string,
  take = false,
): Promise<void> {
  const say = (event: object): void => {
    process.stdout.write(JSON.stringify(event) + "\n");
  };
  const transport = await memberTransport(folder);
  const inbox = await Inbox.open(folder);
  if (session) {
    const holder = await claimSession(folder, session, take);
    if (holder) {
      say({
        event: "knowslink_session_busy",
        session: holder.session,
        since: holder.since,
      });
      process.exitCode = 2;
      return;
    }
    say({
      event: "knowslink_session_bound",
      agent: transport.agentId,
      session,
      pid: process.pid,
      pending: (await inbox.pending()).length, // Unread copies kept from before this session count from the start (SP-08).
    });
  }
  const auto = new AutoReceiver(
    transport,
    inbox,
    async (m, pending) => {
      const c = await transport.contact(m.from).catch(() => null);
      say({
        event: "knowslink_text_received",
        id: m.id,
        from: m.from,
        pending,
        ...(c ? { friend: c.friend, name: c.name } : {}),
      });
    },
    GrokWake.fromEnv(process.env),
  );
  auto.start();
  if (!session) {
    setInterval(() => {}, 1 << 30);
    return;
  }
  let relay = "ok";
  const guard = setInterval(() => {
    // Relay health rides the same tick: unreachable (outage) and error (key, friend or other refusal) need different advice (SP-11, V07).
    const error = auto.status.lastError ?? "";
    const now =
      auto.status.state !== "backoff"
        ? "ok"
        : /^relay (unreachable|HTTP 5)/.test(error)
          ? "unreachable"
          : "error";
    if (now !== relay) {
      relay = now;
      say({
        event: "knowslink_relay",
        state: now,
        ...(now === "ok" ? {} : { error }),
      });
    }
    void holdsSession(folder, session).then((held) => {
      if (held) return;
      clearInterval(guard);
      auto.stop();
      say({ event: "knowslink_session_replaced", session });
    });
  }, 10000);
  const quit = () => {
    clearInterval(guard);
    auto.stop();
    void releaseSession(folder, session).finally(() => process.exit(0));
  };
  process.once("SIGTERM", quit);
  process.once("SIGINT", quit);
}
// `contact <folder> <agent>`: one JSON line with the contact label, for the Claude Code send dialog.
export async function printContact(folder: string, agent: string) {
  const c = await (await memberTransport(folder)).contact(agent);
  process.stdout.write(JSON.stringify(c) + "\n");
}
// `watch <folder> [--session <alias>] [--take]` from either CLI entry.
export function watchFromArgv(argv: string[]): Promise<void> {
  const folder = argv[0];
  if (!folder) throw new Error("usage");
  const at = argv.indexOf("--session");
  return watch(
    folder,
    at >= 0 ? argv[at + 1] : undefined,
    argv.includes("--take"),
  );
}
async function main(): Promise<void> {
  const [, , action, folder, peer, key, confirmed, replyTo] = process.argv;
  if (!folder) throw new Error("usage");
  const t = await memberTransport(folder);
  let result: unknown;
  if (action === "send" && peer && key && confirmed === "--confirmed") {
    let input = "";
    for await (const chunk of process.stdin) {
      input += String(chunk);
      if (Buffer.byteLength(input) > 4096) throw new Error("text too large");
    }
    result = await t.sendText(peer, input, key, true, replyTo);
  } else if (action === "watch") {
    await watchFromArgv(process.argv.slice(3));
    return;
  } else if (action === "receive") {
    const inbox = await Inbox.open(folder);
    result = await inbox.take();
    if (!result) {
      await t.receive((m) => inbox.store(m));
      result = await inbox.take();
    }
  } else if (action === "receipt" && peer) result = await t.receipt(peer);
  else throw new Error("usage");
  process.stdout.write(JSON.stringify(result) + "\n");
}
if (
  process.argv[1] &&
  basename(process.argv[1]) === "text.js" &&
  import.meta.url === pathToFileURL(process.argv[1]).href
) {
  main().catch((e: unknown) => {
    const safe =
      e instanceof Error &&
      /^(relay HTTP [0-9]{3}(?: [a-z_]+)?(?: retry_at=[0-9TZ:.-]+)?|relay unreachable)$/.test(
        e.message,
      )
        ? e.message
        : "text request failed";
    process.stderr.write(
      /^relay (unreachable|HTTP 5)/.test(safe)
        ? `${safe}. relay와 통신할 수 없습니다. 키·친구 문제가 아닙니다. 잠시 뒤 다시 시도하세요. 불확실한 송신은 같은 key·내용으로 재시도하세요.\n`
        : `${safe}. 수신 성공으로 처리하지 마세요. 현재 키·관계·기한을 확인하세요. 불확실한 송신은 같은 key·내용으로 재시도하세요.\n`,
    );
    process.exitCode = 1;
  });
}
