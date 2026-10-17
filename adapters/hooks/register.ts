// KnowsLink work-session link for Claude Code: /knowslink-session on|off|status, arrival toasts and host-side send approval (SP-01/02/05/06/07/09).
// The mod owns one watcher child per session; the child binds the agent folder and stores text locally. The model never sees text until it calls knowslink_text_receive.
import type { Register } from "claude-code";

type Engine = Parameters<Parameters<Register>[0]>[0] extends infer H
  ? H extends (...args: infer A) => unknown
    ? A[0]
    : never
  : never;
type Link = {
  alias: string;
  agent: string;
  pending: number;
  relay: string; // ok | unreachable | error, from the watcher's knowslink_relay events
  stop: () => void;
};
const TOOL = /knowslink_text_(send|receive|receipt)$/;
const CONTROL = /[\p{Cc}\p{Cf}]/gu;
// Variation selectors, bidi and other format or default-ignorable code points draw as nothing; the dialog names each one (review #8).
const INVISIBLE = /[\p{Cf}\p{Default_Ignorable_Code_Point}]/gu;
// Claude Code 2.1.302 draws an AskUserQuestion question up to its first 2000 characters and ends a longer one with an ellipsis (O-RT-01).
const ASK_MAX = 2000;
let link: Link | null = null;

const visible = (text: string): string =>
  text
    .replace(
      INVISIBLE,
      (c) =>
        `⟨U+${c.codePointAt(0)!.toString(16).toUpperCase().padStart(4, "0")}⟩`,
    )
    .replace(CONTROL, " ");
const line = (l: Link): string =>
  l.relay === "unreachable"
    ? `KnowsLink ${l.alias} 통신 불가 · 재시도 중`
    : l.relay === "error"
      ? `KnowsLink ${l.alias} 수신 오류 · 키·친구 확인`
      : l.pending
        ? `KnowsLink ${l.alias} 대기 ${l.pending}`
        : `KnowsLink ${l.alias} 연결`;
const summary = (l: Link): string =>
  `세션 ${l.alias} · agent ${l.agent} · 대기 ${l.pending}건${l.relay === "unreachable" ? " · relay 통신 불가(자동 재시도 중)" : l.relay === "error" ? " · 수신 오류(키·친구 확인)" : ""}`;

async function readConfig(
  $: Engine,
): Promise<{ argv: string[]; folder: string }> {
  const text = await $.fs.read(`${$.plugin.root}/config/local.mcp.json`);
  const server = (
    JSON.parse(text) as {
      mcpServers: {
        knowslink: {
          command: string;
          args: string[];
          env?: Record<string, string>;
        };
      };
    }
  ).mcpServers.knowslink;
  const env = server.env ?? {};
  if (env.KNOWSLINK_MODE !== "public-node" || !env.KNOWSLINK_AGENT_FOLDER)
    throw new Error(
      "이 설치는 held 상태입니다. package_plugin.py --host claude --install-dir <새 경로> --connect 로 준비한 뒤 Google 연결을 마치세요.",
    );
  return {
    argv: [server.command, ...server.args],
    folder: env.KNOWSLINK_AGENT_FOLDER,
  };
}

function parse(line: string): Record<string, unknown> | null {
  try {
    const event = JSON.parse(line) as unknown;
    return event && typeof event === "object"
      ? (event as Record<string, unknown>)
      : null;
  } catch {
    return null;
  }
}

// Friend and agent names are friend-chosen text: shown cut and with invisible characters named, beside the agent ID (O-03).
const named = (friend: unknown, name: unknown, id: string): string =>
  typeof friend === "string" && typeof name === "string" && friend && name
    ? `친구 ${visible(friend.slice(0, 40))} · agent ${visible(name.slice(0, 40))} (${id})`
    : `agent ${id}`;

// The send dialog's label for a peer ID: one short `plugin.js contact` run; on any failure the ID alone.
async function contact($: Engine, peer: string): Promise<string> {
  try {
    const { argv, folder } = await readConfig($);
    const child = $.process.spawn({ argv: [...argv, "contact", folder, peer] });
    const timer = $.clock.after(5000, () => {
      void child.return({ code: null, signal: null });
    });
    let out = "";
    try {
      for (;;) {
        const piece = await child.next();
        if (piece.done) break;
        if (piece.value.stream === "stdout") out += piece.value.text;
      }
    } finally {
      timer.cancel();
    }
    const c = parse(out.trim());
    return named(c?.friend, c?.name, peer);
  } catch {
    return named(null, null, peer);
  }
}

async function turnOn($: Engine, take: boolean): Promise<string> {
  if (link)
    return `KnowsLink 이미 연결됨: ${summary(link)}. 끄려면 /knowslink-session off`;
  let config: { argv: string[]; folder: string };
  try {
    config = await readConfig($);
  } catch (e) {
    return `KnowsLink 연결 불가: ${e instanceof Error ? e.message : "설정을 읽지 못했습니다"}`;
  }
  const alias =
    (await $.session.id()).replace(/[^A-Za-z0-9_-]/g, "").slice(0, 8) ||
    "session";
  const argv = [...config.argv, "watch", config.folder, "--session", alias];
  if (take) argv.push("--take");
  const child = $.process.spawn({ argv });
  let resolveFirst: (text: string) => void = () => {};
  const first = new Promise<string>((resolve) => (resolveFirst = resolve));
  const timer = $.clock.after(20000, () =>
    resolveFirst(
      "KnowsLink 연결 확인 응답이 없습니다. /knowslink-session status 로 다시 확인하세요.",
    ),
  );
  const current: Link = {
    alias,
    agent: "",
    pending: 0,
    relay: "ok",
    stop: () => void child.return({ code: null, signal: null }),
  };
  void (async () => {
    let buffer = "";
    let ended: { code: number | null } = { code: null };
    try {
      for (;;) {
        const piece = await child.next();
        if (piece.done) {
          ended = piece.value ?? { code: null }; // a stopped child ends with no result
          break;
        }
        if (piece.value.stream !== "stdout") continue;
        buffer += piece.value.text;
        let cut: number;
        while ((cut = buffer.indexOf("\n")) >= 0) {
          const event = parse(buffer.slice(0, cut));
          buffer = buffer.slice(cut + 1);
          if (!event) continue;
          switch (event.event) {
            case "knowslink_session_bound":
              current.agent = String(event.agent ?? "");
              current.pending = Number(event.pending ?? 0);
              link = current;
              $.ui.status(line(current));
              resolveFirst(
                `KnowsLink 연결됨: 세션 ${alias} · 내 agent ${current.agent} · 대기 ${current.pending}건. 도착은 toast로 알립니다. 본문은 요청 때 knowslink_text_receive 로 외부 데이터로 봅니다.`,
              );
              break;
            case "knowslink_session_busy":
              resolveFirst(
                `KnowsLink 연결 보류: 다른 세션 ${String(event.session)} 이(가) ${String(event.since)} 부터 이 agent를 받고 있습니다. 전환하려면 /knowslink-session on --take`,
              );
              break;
            case "knowslink_text_received": {
              current.pending = Number(event.pending ?? 0);
              const from = visible(String(event.from ?? "").slice(0, 64));
              $.ui.toast(
                `KnowsLink 외부 메시지 도착 · 발신 ${named(event.friend, event.name, from)} · 대기 ${current.pending}건 · 확인: knowslink_text_receive`,
                { timeoutMs: 12000 },
              );
              if (link === current) $.ui.status(line(current));
              break;
            }
            case "knowslink_relay": {
              const state = String(event.state ?? "");
              if (link !== current || state === current.relay) break;
              current.relay = state;
              $.ui.status(line(current));
              $.ui.toast(
                state === "unreachable"
                  ? "KnowsLink relay와 통신할 수 없습니다. 키·친구 문제가 아닙니다. 자동으로 다시 시도합니다."
                  : state === "error"
                    ? `KnowsLink 수신 오류(${String(event.error ?? "")
                        .replace(CONTROL, "")
                        .slice(0, 80)}). 키·친구 관계를 확인하세요.`
                    : "KnowsLink relay 통신이 복구됐습니다.",
              );
              break;
            }
            case "knowslink_session_replaced":
              // Release now, not when the child exits: until then the old session's text tools would still pass.
              if (link === current) {
                link = null;
                $.ui.status(undefined);
                current.stop();
              }
              $.ui.toast(
                "KnowsLink 연결 해제: 다른 세션이 이 agent를 가져갔습니다.",
              );
              break;
          }
        }
      }
    } catch {
      ended = { code: -1 };
    } finally {
      timer.cancel();
      if (link === current) {
        link = null;
        $.ui.status(undefined);
        if (ended.code !== null && ended.code !== 0)
          $.ui.toast(
            `KnowsLink 연결 종료 (exit ${ended.code}). /knowslink-session on 으로 다시 연결하세요.`,
          );
      }
      resolveFirst(`KnowsLink 연결 실패 (exit ${String(ended.code)}).`);
    }
  })();
  return first;
}

function turnOff($: Engine): string {
  if (!link) return "KnowsLink: 이 세션은 연결되어 있지 않습니다.";
  const { alias, stop } = link;
  link = null;
  stop();
  $.ui.status(undefined);
  return `KnowsLink 연결 해제: 세션 ${alias}. 이 세션은 더 받거나 보내지 않습니다.`;
}

export const register: Register = (on) => {
  on("session.start", async ($, e, next) => {
    $.ui.status(undefined); // A reload starts the module over: its watcher child is gone, so the owner runs on again.
    await $.command.register({
      name: "knowslink-session",
      description: "KnowsLink 세션 연결 켜기·끄기·상태",
      argumentHint: "on [--take] | off | status",
    });
    return next(e);
  });
  on("command.run", { command: "knowslink-session" }, async ($, e) => {
    const words = e.args.trim().split(/\s+/);
    if (words[0] === "on")
      return { text: await turnOn($, words.includes("--take")) };
    if (words[0] === "off") return { text: turnOff($) };
    return {
      text: link
        ? `KnowsLink 연결됨: ${summary(link)}`
        : "KnowsLink: 이 세션은 연결되어 있지 않습니다. /knowslink-session on",
    };
  });
  // Text tools work only in the session the owner switched on; a send needs the owner's own approval here, never the model's or a received text's.
  on("tool.call", async ($, e, next) => {
    const kind = TOOL.exec(e.tool)?.[1];
    if (!kind) return next(e);
    // A subagent or teammate loop of this session (agentId set) would share the owner's link and inbox; only the main conversation may use it (SP-04).
    if (e.agentId)
      return {
        deny: "KnowsLink: 하위 에이전트에서는 KnowsLink 메시지를 쓸 수 없습니다. 메인 대화에서 사용하세요.",
      };
    if (!link)
      return {
        deny: "KnowsLink: 이 세션은 연결되지 않았습니다. /knowslink-session on 으로 먼저 연결하세요.",
      };
    if (kind === "send") {
      const args = e as Record<string, unknown>;
      if (typeof args.peer !== "string" || typeof args.text !== "string")
        return {
          deny: "KnowsLink: peer와 text가 필요합니다. 승인 전 인자를 확인하세요.",
        };
      const peer = await contact($, visible(args.peer.slice(0, 128)));
      const kindLine =
        typeof args.reply_to === "string" && args.reply_to
          ? `답장(reply_to ${visible(args.reply_to.slice(0, 64))})`
          : "새 메시지";
      const body = visible(args.text);
      const hidden = args.text.match(INVISIBLE)?.length ?? 0;
      if (!link)
        return {
          deny: "KnowsLink: 이 세션은 연결되지 않았습니다. /knowslink-session on 으로 먼저 연결하세요.",
        };
      const question = `KnowsLink 발송 승인: 내 agent ${link.agent} → ${peer} · ${kindLine} · 본문 전체 ${args.text.length}자${hidden ? ` · 보이지 않는 문자 ${hidden}개를 ⟨U+…⟩로 표시` : ""} 「${body}」 을 그대로 보낼까요?`;
      // The owner approves exactly what is sent: the whole body, never a cut preview.
      if (question.length > ASK_MAX)
        return {
          deny: `KnowsLink: 승인 창은 ${ASK_MAX}자까지만 보여 이 본문(${args.text.length}자)의 끝을 확인할 수 없습니다. 발송을 막았습니다. 본문을 약 ${Math.max(0, ASK_MAX - (question.length - body.length))}자 이하로 나눠 각각 승인받으세요.`,
        };
      const held = link;
      const answer = await $.ui
        .ask(question, ["보내기", "취소"])
        .catch(() => "취소");
      // The link may end while the dialog is open (replaced or off); a released session never sends (SP-09, review #7).
      if (link !== held)
        return {
          deny: "KnowsLink: 승인하는 동안 이 세션의 연결이 해제돼 발송을 막았습니다. /knowslink-session status 로 확인하세요.",
        };
      if (answer !== "보내기")
        return {
          deny: "KnowsLink: 소유자가 이 발송을 승인하지 않았습니다. 답장은 소유자의 별도 승인이 있어야 보냅니다.",
        };
    }
    const ran = await next(e);
    if (kind === "receive" && link && !ran.deny && !ran.isError) {
      link.pending = Math.max(0, link.pending - 1);
      $.ui.status(line(link));
    }
    return ran;
  }).catch(($, e, next) =>
    next.called
      ? next(e)
      : { deny: "KnowsLink: 승인 확인에 실패해 이 호출을 막았습니다." },
  );
};
