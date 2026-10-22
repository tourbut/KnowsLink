// KnowsLink work-session link for Codex (SAR-CODEX-RECV-001): Codex has no host mod, so the MCP process Codex starts for one session is the session.
// Binding and every send need the owner's answer to an MCP elicitation drawn by Codex; the model's tool arguments never bind or approve (SP-01/05/06).
import { randomBytes } from "node:crypto";
import { AutoReceiver, Inbox } from "./inbox.js";
import {
  alive,
  claimSession,
  holdsSession,
  readLock,
  releaseSession,
} from "./session-lock.js";
import type { ReceivedText, TextTransport } from "./text.js";

// Resolves true only when the owner accepts the dialog; decline, cancel, timeout or a client without elicitation is false.
export type Ask = (message: string) => Promise<boolean>;
export type Notice = (m: ReceivedText, pending: number) => Promise<void>;
const CONTROL = /[\p{Cc}\p{Cf}]/gu;
// Same rule as the Claude Code dialog (hooks/register.ts): invisible code points are named, never drawn as nothing.
const INVISIBLE = /[\p{Cf}\p{Default_Ignorable_Code_Point}]/gu;
export const visible = (text: string): string =>
  text
    .replace(
      INVISIBLE,
      (c) =>
        `⟨U+${c.codePointAt(0)!.toString(16).toUpperCase().padStart(4, "0")}⟩`,
    )
    .replace(CONTROL, " ");

export class CodexSession {
  // Non-secret distinguishing value shown to the owner (SP-01); one per MCP process, so a new Codex session never inherits it.
  readonly alias = `codex-${randomBytes(3).toString("hex")}`;
  private auto: AutoReceiver | null = null;
  private guard: ReturnType<typeof setInterval> | null = null;
  private ended = "";

  constructor(
    private readonly folder: string,
    private readonly transport: TextTransport,
    private readonly ask: Ask,
    private readonly notice: Notice,
  ) {}

  get bound(): boolean {
    return this.auto !== null;
  }
  get inbox(): Inbox | null {
    return this.auto?.inbox ?? null;
  }
  serial<T>(task: () => Promise<T>): Promise<T> {
    return this.auto ? this.auto.serial(task) : task();
  }

  async on(take: boolean): Promise<{ state: string; text: string }> {
    const agent = this.transport.agentId;
    if (this.auto)
      return {
        state: "bound",
        text: `KnowsLink 이미 연결됨: ${await this.summary()}`,
      };
    const held = await readLock(this.folder);
    const other = held && held.session !== this.alias && alive(held.pid);
    if (other && !take)
      return {
        state: "busy",
        text: `KnowsLink 연결 보류: 다른 세션 ${held.session} 이(가) ${held.since} 부터 agent ${agent} 를 받고 있습니다. 전환하려면 소유자가 take로 다시 요청하세요.`,
      };
    const question = `KnowsLink 세션 연결: 이 Codex 세션(구분값 ${this.alias})을 내 agent ${agent} 의 수신 세션으로 연결할까요?${other ? ` 다른 세션 ${held.session} 의 수신은 끊깁니다.` : ""} 친구 메시지는 외부 데이터로만 보이고, 발송은 매번 이 창에서 따로 승인합니다.`;
    if (!(await this.ask(question)))
      return {
        state: "declined",
        text: "KnowsLink: 소유자가 이 세션 연결을 승인하지 않았습니다.",
      };
    const blocker = await claimSession(this.folder, this.alias, take);
    if (blocker)
      return {
        state: "busy",
        text: `KnowsLink 연결 보류: 다른 세션 ${blocker.session} 이(가) ${blocker.since} 부터 받고 있습니다.`,
      };
    const inbox = await Inbox.open(this.folder);
    this.auto = new AutoReceiver(this.transport, inbox, this.notice);
    this.auto.start();
    this.ended = "";
    // Another session's takeover or a lost lock ends this link at the next tick (SP-09); no automatic re-claim.
    this.guard = setInterval(() => {
      void holdsSession(this.folder, this.alias).then((held) => {
        if (!held) this.stop("replaced");
      });
    }, 10000);
    this.guard.unref();
    return {
      state: "bound",
      text: `KnowsLink 연결됨: ${await this.summary()}. 본문은 요청 때 knowslink_text_receive 로 외부 데이터로 봅니다.`,
    };
  }

  async off(): Promise<{ state: string; text: string }> {
    if (!this.auto)
      return {
        state: "unbound",
        text: "KnowsLink: 이 세션은 연결되어 있지 않습니다.",
      };
    this.stop("off");
    await releaseSession(this.folder, this.alias);
    return {
      state: "unbound",
      text: `KnowsLink 연결 해제: Codex 세션 ${this.alias}. 이 세션은 더 받거나 보내지 않습니다.`,
    };
  }

  async status(): Promise<{ state: string; text: string }> {
    if (this.auto)
      return {
        state: "bound",
        text: `KnowsLink 연결됨: ${await this.summary()}`,
      };
    return {
      state: "unbound",
      text:
        this.ended === "replaced"
          ? "KnowsLink 연결 해제: 다른 세션이 이 agent를 가져갔습니다."
          : `KnowsLink: 이 Codex 세션(구분값 ${this.alias})은 연결되어 있지 않습니다. knowslink_session on 을 소유자가 승인해야 받습니다.`,
    };
  }

  // Null when this session may use text tools now; otherwise the deny text.
  async refusal(): Promise<string | null> {
    if (this.auto && !(await holdsSession(this.folder, this.alias)))
      this.stop("replaced");
    return this.auto
      ? null
      : this.ended === "replaced"
        ? "KnowsLink: 다른 세션이 이 agent를 가져가 이 세션의 연결이 해제됐습니다."
        : "KnowsLink: 이 세션은 연결되지 않았습니다. 소유자가 knowslink_session on 을 승인한 뒤 사용하세요.";
  }

  // Shows the whole body as it will be sent; null means the owner approved exactly this one attempt.
  async approveSend(
    peer: string,
    text: string,
    replyTo?: string,
  ): Promise<string | null> {
    const deny = await this.refusal();
    if (deny) return deny;
    const c = await this.transport.contact(peer).catch(() => null);
    const target =
      c?.friend && c.name
        ? `친구 ${visible(c.friend.slice(0, 40))} · agent ${visible(c.name.slice(0, 40))} (${visible(peer)})`
        : `agent ${visible(peer)}`;
    const hidden = text.match(INVISIBLE)?.length ?? 0;
    const question = `KnowsLink 발송 승인: 내 agent ${this.transport.agentId} → ${target} · ${replyTo ? `답장(reply_to ${visible(replyTo.slice(0, 64))})` : "새 메시지"} · 본문 전체 ${text.length}자${hidden ? ` · 보이지 않는 문자 ${hidden}개를 ⟨U+…⟩로 표시` : ""} 「${visible(text)}」 을 그대로 보낼까요?`;
    const held = this.auto;
    const ok = await this.ask(question);
    // A link released while the dialog was open never sends (SP-09).
    if (held !== this.auto || (await this.refusal()))
      return "KnowsLink: 승인하는 동안 이 세션의 연결이 해제돼 발송을 막았습니다.";
    return ok
      ? null
      : "KnowsLink: 소유자가 이 발송을 승인하지 않았습니다. 답장은 소유자의 별도 승인이 있어야 보냅니다.";
  }

  stop(reason: string): void {
    if (this.guard) clearInterval(this.guard);
    this.guard = null;
    this.auto?.stop();
    this.auto = null;
    this.ended = reason;
  }

  private async summary(): Promise<string> {
    const pending = this.auto ? (await this.auto.inbox.pending()).length : 0;
    return `Codex 세션 ${this.alias} · 내 agent ${this.transport.agentId} · 대기 ${pending}건`;
  }
}
