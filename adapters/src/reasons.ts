// Connection failure reasons (D02 E): one code and its member-language text; never a key, token or credential.
import { PrivateFolderError } from "./private-files.js";
import { approvedAttempt } from "./login.js";

export type Reason =
  | "expired"
  | "cancelled"
  | "unconfirmed"
  | "private_folder"
  | "unreachable"
  | "rate_limited"
  | "agent_rate"
  | "service_capacity";
// signed_out is a web-only reason: the client never holds a member session.
const TEXT: Record<"ko" | "en", Record<Reason, string>> = {
  ko: {
    expired: "연결 시간이 지났습니다. 다시 연결하면 새 링크를 만듭니다.",
    cancelled: "연결을 취소했습니다. 다시 연결하면 새 링크를 만듭니다.",
    unconfirmed:
      "브라우저 동의는 끝났지만 키 저장을 확인하지 못했습니다. 다시 연결하면 먼저 이 시도를 확인합니다.",
    private_folder:
      "키를 보관할 비공개 폴더를 준비하지 못했습니다(<path>). 기본 위치로 다시 연결하세요.",
    unreachable:
      "KnowsLink 서버에 연결할 수 없습니다. 인터넷 연결을 확인하고 다시 연결하세요.",
    rate_limited: "요청이 많아 잠시 멈췄습니다. <time> 뒤 다시 연결하세요.",
    agent_rate:
      "오늘 새 에이전트를 많이 만들었습니다. <time> 뒤 다시 연결하세요. 기존 에이전트는 그대로 쓸 수 있습니다.",
    service_capacity:
      "지금은 서비스 용량 때문에 새 에이전트를 만들 수 없습니다. 기존 에이전트는 그대로 쓸 수 있습니다. 나중에 다시 연결하세요.",
  },
  en: {
    expired: "The connection timed out. Connect again to get a new link.",
    cancelled: "You cancelled the connection. Connect again to get a new link.",
    unconfirmed:
      "You agreed in the browser, but key storage wasn't confirmed. Connect again; it checks this attempt first.",
    private_folder:
      "Couldn't prepare a private folder for the key (<path>). Connect again using the default location.",
    unreachable:
      "Can't reach the KnowsLink server. Check your internet connection and connect again.",
    rate_limited: "Too many requests. Connect again after <time>.",
    agent_rate:
      "You've made many new agents today. Connect again after <time>. Your existing agents still work.",
    service_capacity:
      "The service can't add new agents right now. Your existing agents still work. Try again later.",
  },
};
export const lang = (): "ko" | "en" =>
  process.env.KNOWSLINK_LANG === "en" ? "en" : "ko";

export type Failure = {
  state: "failed" | "unconfirmed";
  reason: Reason;
  message: string;
  retry_at?: string;
  agent?: string;
  name?: string;
};
export function failure(reason: Reason, path = "", retryAt?: string): Failure {
  const en = lang() === "en";
  const time = retryAt
    ? `${retryAt.slice(0, 16).replace("T", " ")} UTC`
    : en
      ? "a moment"
      : "잠시";
  return {
    state: reason === "unconfirmed" ? "unconfirmed" : "failed",
    reason,
    message: TEXT[en ? "en" : "ko"][reason]
      .replace("<path>", path)
      .replace("<time>", time),
    ...(retryAt ? { retry_at: retryAt } : {}),
  };
}

// Maps a connection error to its reason. An approved attempt whose completion failed stays unconfirmed with its agent and
// name (F-07) unless the relay refused the new agent itself; protocol surprises read as a server problem.
export async function failureOf(e: unknown, folder: string): Promise<Failure> {
  const m = String((e as Error)?.message ?? e);
  const retry = / retry_at=(\S+)$/.exec(m)?.[1];
  if (/^relay HTTP 429 agent_rate/.test(m))
    return failure("agent_rate", "", retry);
  if (/^relay HTTP 409 service_capacity/.test(m))
    return failure("service_capacity");
  const lost = await approvedAttempt(folder);
  if (lost)
    return {
      ...failure("unconfirmed"),
      agent: lost.agent ?? "",
      name: typeof lost.name === "string" ? lost.name.slice(0, 80) : "",
    };
  if (
    e instanceof PrivateFolderError ||
    typeof (e as NodeJS.ErrnoException)?.code === "string"
  )
    return failure("private_folder", folder);
  if (/^relay HTTP 429/.test(m)) return failure("rate_limited", "", retry);
  if (m === "login expired") return failure("expired");
  if (m === "login cancelled") return failure("cancelled");
  return failure("unreachable");
}
