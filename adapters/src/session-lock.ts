// One receiving work session per agent folder (SP-02/SP-09): the lock file is the claim, and a replaced or dead holder never receives again.
import { open, readFile, rm, writeFile } from "node:fs/promises";
import { join } from "node:path";

export type SessionLock = { session: string; pid: number; since: string };
const SESSION = /^[A-Za-z0-9_-]{1,64}$/;
export const lockPath = (folder: string): string =>
  join(folder, "session.lock");

export function alive(pid: number): boolean {
  try {
    process.kill(pid, 0);
    return true;
  } catch (e) {
    return (e as NodeJS.ErrnoException).code === "EPERM";
  }
}
export async function readLock(folder: string): Promise<SessionLock | null> {
  try {
    const l = JSON.parse(
      await readFile(lockPath(folder), "utf8"),
    ) as SessionLock;
    return SESSION.test(l.session) && Number.isInteger(l.pid) ? l : null;
  } catch {
    return null;
  }
}
// Resolves the holder that blocks the claim, or null once this process holds the lock.
export async function claimSession(
  folder: string,
  session: string,
  take = false,
): Promise<SessionLock | null> {
  if (!SESSION.test(session)) throw new Error("invalid session alias");
  const mine = JSON.stringify({
    session,
    pid: process.pid,
    since: new Date().toISOString(),
  });
  try {
    const handle = await open(lockPath(folder), "wx", 0o600);
    try {
      await handle.writeFile(mine);
    } finally {
      await handle.close();
    }
    return null;
  } catch (e) {
    if ((e as NodeJS.ErrnoException).code !== "EEXIST") throw e;
  }
  const held = await readLock(folder);
  // A live holder of another session keeps the address until the owner explicitly takes it over.
  if (held && held.session !== session && alive(held.pid) && !take) return held;
  await writeFile(lockPath(folder), mine, { mode: 0o600 });
  return null;
}
export async function holdsSession(
  folder: string,
  session: string,
): Promise<boolean> {
  const l = await readLock(folder);
  return l?.session === session && l.pid === process.pid;
}
export async function releaseSession(
  folder: string,
  session: string,
): Promise<void> {
  if (await holdsSession(folder, session))
    await rm(lockPath(folder), { force: true });
}
