// Installer-driven Google connection (SAR-AUTOSETUP-001 AS-05/AS-11): begins, resumes or confirms one attempt in the
// installation's private folder and reports only public status (URL, fingerprint, agent ID/name) as JSON lines.
import { createHash } from "node:crypto";
import { access, readFile, rename } from "node:fs/promises";
import { join } from "node:path";
import { setTimeout } from "node:timers/promises";
import { relayBase } from "./core.js";
import { complete } from "./connect.js";
import { beginLogin, finishLogin } from "./login.js";

export type SetupEvent = Record<string, string>;
type Outcome = "connected" | "failed" | "unconfirmed";

const exists = (path: string) =>
  access(path).then(
    () => true,
    () => false,
  );
async function json(path: string): Promise<Record<string, string> | null> {
  return (await exists(path)) ? JSON.parse(await readFile(path, "utf8")) : null;
}
const sha = (data: string | Buffer) =>
  createHash("sha256").update(data).digest("base64url");
// No HTTP answer, timeouts, 5xx and rate limits say nothing about the attempt itself; keep polling.
const transient = (e: unknown) =>
  /unreachable|timed out|HTTP 5\d\d|rate_limited|empty relay response/.test(
    String((e as Error)?.message ?? e),
  );

export async function setupConnect(
  base: string,
  folder: string,
  fresh: boolean,
  emit: (event: SetupEvent) => void,
  pollMs = 5000,
): Promise<Outcome> {
  base = relayBase(base);
  // agent.json is never read here; its existence alone means the key is already stored (AS-10).
  if (await exists(join(folder, "agent.json"))) {
    emit({ state: "connected", existing: "true" });
    return "connected";
  }
  const aside = async () => {
    if (await exists(folder))
      await rename(folder, `${folder}.unused-${Date.now()}`);
  };
  const pending = () => json(join(folder, "pending.json"));
  let approved = await pending();
  if (fresh) {
    await aside();
    approved = null;
  }
  // An approved attempt whose credential was not saved must be confirmed before any new attempt (F-07).
  const unsaved = (p: Record<string, string> | null) =>
    p !== null && p.state !== "consumed" && p.mode === "register";
  const login = await json(join(folder, "login.json"));
  // A cancelled or relay-expired attempt is never reused (AS-11); only a waiting or unanswered one resumes.
  let probe = "dead";
  if (unsaved(approved)) probe = "approved";
  else if (
    login &&
    login.state !== "consumed" &&
    Date.parse(login.exp ?? "") > Date.now()
  )
    probe = await finishLogin(folder).catch((e: unknown) =>
      transient(e) ? "waiting" : "dead",
    );
  if (probe === "dead" && unsaved(await pending())) probe = "approved";
  const of = (p: Record<string, string>) =>
    `SHA256:${sha(Buffer.from(p.public!, "base64url"))}`;
  let fingerprint: string;
  let expires: string;
  if (probe === "dead") {
    // Expired, cancelled or empty leftovers are kept aside, never deleted (AS-11: no automatic cleanup).
    await aside();
    const started = await beginLogin(base, folder);
    ({ fingerprint, expires } = started);
    emit({ state: "waiting", url: started.url, fingerprint, expires });
  } else {
    const p = probe === "approved" ? (await pending())! : login!;
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
    const p = await pending();
    try {
      if (unsaved(p)) await complete(folder);
      else if ((await finishLogin(folder)) === "waiting") {
        await setTimeout(pollMs);
        continue;
      }
      break;
    } catch (e) {
      const approvedNow = unsaved(await pending());
      if (transient(e) && Date.now() < deadline + 60000) {
        await setTimeout(pollMs);
        continue;
      }
      if (approvedNow) {
        const lost = (await pending())!;
        emit({
          state: "unconfirmed",
          agent: lost.agent!,
          name: typeof lost.name === "string" ? lost.name.slice(0, 80) : "",
        });
        return "unconfirmed";
      }
      emit({
        state: "failed",
        reason: transient(e) ? "unreachable" : "cancelled_or_expired",
      });
      return "failed";
    }
  }
  const { agent, name } = JSON.parse(
    await readFile(join(folder, "agent.json"), "utf8"),
  ) as { agent: string; name?: string };
  emit({ state: "connected", agent, name: name ?? "", fingerprint });
  return "connected";
}
