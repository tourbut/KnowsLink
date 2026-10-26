// Installer-driven Google connection (SAR-AUTOSETUP-001 AS-05/AS-11): begins, resumes or confirms one attempt in the
// installation's private folder and reports only public status (URL, fingerprint, agent ID/name) as JSON lines.
import { createHash } from "node:crypto";
import { access, readFile, readdir, rename } from "node:fs/promises";
import { join } from "node:path";
import { setTimeout } from "node:timers/promises";
import { relayBase } from "./core.js";
import { complete } from "./connect.js";
import { approvedAttempt, beginLogin, finishLogin } from "./login.js";
import { failureOf } from "./reasons.js";

export type SetupEvent = Record<string, string | undefined>;
type Outcome = "connected" | "failed" | "unconfirmed";

const exists = (path: string) =>
  access(path).then(
    () => true,
    () => false,
  );
// A missing or unreadable attempt file is a dead leftover (R-01), not a failure.
async function json(path: string): Promise<Record<string, string> | null> {
  try {
    return JSON.parse(await readFile(path, "utf8"));
  } catch {
    return null;
  }
}
const sha = (data: string | Buffer) =>
  createHash("sha256").update(data).digest("base64url");
// No HTTP answer, timeouts, 5xx and rate limits say nothing about the attempt itself; keep polling.
const transient = (e: unknown) =>
  /unreachable|timed out|HTTP 5\d\d|rate_limited|empty relay response/.test(
    String((e as Error)?.message ?? e),
  );

// Every failure carries a D02 E reason and its text; an approved attempt that could not be confirmed keeps its agent and name (F-07).
export async function setupConnect(
  base: string,
  folder: string,
  fresh: boolean,
  emit: (event: SetupEvent) => void,
  pollMs = 5000,
): Promise<Outcome> {
  try {
    return await connect(base, folder, fresh, emit, pollMs);
  } catch (e) {
    const f = await failureOf(e, folder);
    emit(f);
    return f.state;
  }
}

async function connect(
  base: string,
  folder: string,
  fresh: boolean,
  emit: (event: SetupEvent) => void,
  pollMs: number,
): Promise<Outcome> {
  base = relayBase(base);
  // agent.json is never read here; its existence alone means the key is already stored (AS-10).
  if (await exists(join(folder, "agent.json"))) {
    emit({ state: "connected", existing: "true" });
    return "connected";
  }
  // An empty folder is used as it is (R-01); anything else is kept aside, never deleted (AS-11: no automatic cleanup).
  const aside = async () => {
    if ((await readdir(folder).catch(() => [])).length > 0)
      await rename(folder, `${folder}.unused-${Date.now()}`);
  };
  if (fresh) await aside();
  const login = await json(join(folder, "login.json"));
  // An approved attempt whose credential was not saved must be confirmed before any new attempt (F-07).
  // A cancelled or relay-expired attempt is never reused (AS-11); only a waiting or unanswered one resumes.
  let probe = "dead";
  if (await approvedAttempt(folder)) probe = "approved";
  else if (
    login &&
    login.state !== "consumed" &&
    Date.parse(login.exp ?? "") > Date.now()
  )
    probe = await finishLogin(folder).catch((e: unknown) =>
      transient(e) ? "waiting" : "dead",
    );
  if (probe === "dead" && (await approvedAttempt(folder))) probe = "approved";
  const of = (p: Record<string, string>) =>
    `SHA256:${sha(Buffer.from(p.public!, "base64url"))}`;
  let fingerprint: string;
  let expires: string;
  if (probe === "dead") {
    await aside();
    const started = await beginLogin(base, folder);
    ({ fingerprint, expires } = started);
    emit({ state: "waiting", url: started.url, fingerprint, expires });
  } else {
    const p = probe === "approved" ? (await approvedAttempt(folder))! : login!;
    fingerprint = of(p);
    expires = p.exp!;
    if (probe === "waiting")
      emit({
        state: "waiting",
        resumed: "true",
        url: `${base}/connect/${sha(p.token!)}`,
        fingerprint,
        expires,
      });
  }
  const deadline = Date.parse(expires) || Date.now();
  while (probe !== "connected") {
    try {
      if (await approvedAttempt(folder)) await complete(folder);
      else if ((await finishLogin(folder)) === "waiting") {
        await setTimeout(pollMs);
        continue;
      }
      break;
    } catch (e) {
      if (!transient(e) || Date.now() >= deadline + 60000) throw e;
      await setTimeout(pollMs);
    }
  }
  const { agent, name } = JSON.parse(
    await readFile(join(folder, "agent.json"), "utf8"),
  ) as { agent: string; name?: string };
  emit({ state: "connected", agent, name: name ?? "", fingerprint });
  return "connected";
}
