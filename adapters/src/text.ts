// Explicit member connection-check CLI/MCP transport: private onboarding folder, signed short text, manual receive and no automatic reply.
import { createPublicKey, sign, verify } from "node:crypto";
import { readFile } from "node:fs/promises";
import { join, basename } from "node:path";
import { pathToFileURL } from "node:url";
import { Adapter, canonical, uuid7 } from "./core.js";
import { relayBase } from "./core.js";
import { privatePath } from "./private-files.js";

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
export function textSigningBytes(e: TextEnvelope): Buffer {
  return Buffer.from(
    `KNOWSLINK-TEXT\0knowslink.text.v1\0Ed25519\0${e.sig.kid}\0${canonical({ v: e.v, id: e.id, from: e.from, to: e.to, text: e.text, exp: e.exp, idempotency_key: e.idempotency_key, reply_to: e.reply_to ?? "" })}`,
  );
}
function validText(text: string): boolean {
  return (
    !!text.trim() &&
    Buffer.byteLength(text) <= 4096 &&
    !/[\uD800-\uDFFF]/u.test(text)
  );
}
export class TextTransport extends Adapter {
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
  async receive(): Promise<{
    id: string;
    from: string;
    to: string;
    text: string;
    exp: string;
    reply_to?: string;
    untrusted: true;
  } | null> {
    const lease = await this.request<{
      envelope: TextEnvelope;
      lease_token: string;
    } | null>("/v1/text/pull", {});
    if (!lease) return null;
    const e = lease.envelope;
    const allowed = [
      "v",
      "id",
      "from",
      "to",
      "text",
      "exp",
      "idempotency_key",
      "reply_to",
      "sig",
    ];
    if (
      Object.keys(e).some((k) => !allowed.includes(k)) ||
      e.v !== "knowslink.text.v1" ||
      e.to !== this.agent ||
      typeof e.text !== "string" ||
      !validText(e.text) ||
      e.sig.alg !== "Ed25519" ||
      !Number.isFinite(Date.parse(e.exp)) ||
      Date.parse(e.exp) <= Date.now()
    )
      throw new Error("invalid text envelope");
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
    const delivery = { id: e.id, token: lease.lease_token };
    await this.request("/v1/text/persist", delivery);
    await this.request("/v1/text/ack", delivery);
    const result: {
      id: string;
      from: string;
      to: string;
      text: string;
      exp: string;
      reply_to?: string;
      untrusted: true;
    } = {
      id: e.id,
      from: e.from,
      to: e.to,
      text: e.text,
      exp: e.exp,
      untrusted: true,
    };
    if (e.reply_to) result.reply_to = e.reply_to;
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
  } else if (action === "receive") result = await t.receive();
  else if (action === "receipt" && peer) result = await t.receipt(peer);
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
      /^relay HTTP [0-9]{3}(?: [a-z_]+)?(?: retry_at=[0-9TZ:.-]+)?$/.test(
        e.message,
      )
        ? e.message
        : "text request failed";
    process.stderr.write(
      `${safe}. 수신 성공으로 처리하지 마세요. 현재 키·관계·기한을 확인하세요. 불확실한 송신은 같은 key·내용으로 재시도하세요.\n`,
    );
    process.exitCode = 1;
  });
}
