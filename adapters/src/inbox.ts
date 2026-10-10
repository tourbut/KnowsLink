// Automatic member text receive: one serialized pull loop per process, private local inbox written before relay persist/ACK, metadata-only host notice.
import {
  open,
  readdir,
  readFile,
  rename,
  rm,
  stat,
  writeFile,
} from "node:fs/promises";
import { join } from "node:path";
import { privateDirectory, privatePath } from "./private-files.js";
import type { ReceivedText, TextTransport } from "./text.js";

const ID =
  /^[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/;
const KEEP_READ_MS = 24 * 3600 * 1000; // Relay TTL is 180s, so a read marker only needs to outlive any re-lease of the same ID.
const exists = (path: string) =>
  stat(path).then(
    () => true,
    () => false,
  );

export class Inbox {
  private constructor(private readonly dir: string) {}
  static async open(folder: string): Promise<Inbox> {
    const dir = join(folder, "inbox");
    try {
      await privateDirectory(dir);
    } catch (e) {
      if ((e as NodeJS.ErrnoException).code !== "EEXIST") throw e;
      await privatePath(dir, true);
    }
    return new Inbox(dir);
  }
  // Returns false for an ID already stored or already read, so a re-leased message is ACKed again but never re-announced.
  async store(m: ReceivedText): Promise<boolean> {
    if (!ID.test(m.id)) throw new Error("invalid message ID");
    const file = join(this.dir, `${m.id}.json`);
    if ((await exists(file)) || (await exists(join(this.dir, `${m.id}.read`))))
      return false;
    const tmp = `${file}.${process.pid}.tmp`;
    const handle = await open(tmp, "wx", 0o600);
    try {
      await handle.writeFile(JSON.stringify(m));
      await handle.sync();
    } finally {
      await handle.close();
    }
    await rename(tmp, file);
    return true;
  }
  async pending(): Promise<string[]> {
    return (await readdir(this.dir))
      .filter((name) => /^[0-9a-f-]{36}\.json$/.test(name))
      .sort(); // UUIDv7 names sort by send time.
  }
  // Oldest unread message. The rename is the claim, so two MCP processes never return the same message.
  async take(): Promise<(ReceivedText & { expired: boolean }) | null> {
    for (const name of await this.pending()) {
      const marker = join(this.dir, name.replace(/\.json$/, ".read"));
      try {
        await rename(join(this.dir, name), marker);
      } catch {
        continue;
      }
      const m = JSON.parse(await readFile(marker, "utf8")) as ReceivedText;
      await writeFile(marker, "", { mode: 0o600 }); // Keep only the ID marker after display.
      await this.prune();
      return { ...m, expired: Date.parse(m.exp) <= Date.now() };
    }
    return null;
  }
  private async prune(): Promise<void> {
    for (const name of await readdir(this.dir)) {
      if (!/\.(read|tmp)$/.test(name)) continue;
      const path = join(this.dir, name);
      const age = Date.now() - (await stat(path)).mtimeMs;
      if (age > (name.endsWith(".tmp") ? 3600000 : KEEP_READ_MS))
        await rm(path, { force: true });
    }
  }
}

export function safeError(e: unknown): string {
  return e instanceof Error &&
    /^(relay HTTP [0-9]{3}(?: [a-z_]+)?(?: retry_at=[0-9TZ:.-]+)?|invalid text (envelope|signature))$/.test(
      e.message,
    )
    ? e.message
    : "failed";
}

export type AutoStatus = {
  state: "running" | "backoff";
  lastSuccessAt?: string;
  lastError?: string;
  lastErrorAt?: string;
  nextPollAt?: string;
  lastReceivedId?: string;
  // sent_unverified: the MCP notification left this process; it does not prove the host showed it or started a turn.
  hostNotice: "none" | "sent_unverified" | "failed";
};

export class AutoReceiver {
  readonly status: AutoStatus = { state: "running", hostNotice: "none" };
  private chain: Promise<unknown> = Promise.resolve();
  private failures = 0;
  constructor(
    private readonly transport: TextTransport,
    readonly inbox: Inbox,
    private readonly notify: (
      m: ReceivedText,
      pending: number,
    ) => Promise<void>,
    private readonly idleMs = 10000,
  ) {}
  // Every relay pull in this process goes through here, so auto ticks and manual receive never lease in parallel.
  serial<T>(f: () => Promise<T>): Promise<T> {
    const run = this.chain.then(f, f);
    this.chain = run.catch(() => {});
    return run;
  }
  start(): void {
    this.schedule(0);
  }
  private schedule(ms: number): void {
    this.status.nextPollAt = new Date(Date.now() + ms).toISOString();
    setTimeout(() => void this.tick(), ms).unref();
  }
  private async tick(): Promise<void> {
    let delay = this.idleMs;
    try {
      const fresh: { m?: ReceivedText } = {};
      const leased = await this.serial(() =>
        this.transport.receive(async (m) => {
          if (await this.inbox.store(m)) fresh.m = m;
        }),
      );
      this.failures = 0;
      this.status.state = "running";
      this.status.lastSuccessAt = new Date().toISOString();
      if (leased) delay = 0; // Drain queued messages; each one was ACKed, so this ends.
      if (fresh.m) {
        this.status.lastReceivedId = fresh.m.id;
        try {
          await this.notify(fresh.m, (await this.inbox.pending()).length);
          this.status.hostNotice = "sent_unverified";
        } catch {
          this.status.hostNotice = "failed"; // The message stays in the local inbox for manual receive.
        }
      }
    } catch (e) {
      this.failures++;
      this.status.state = "backoff";
      this.status.lastError = safeError(e);
      this.status.lastErrorAt = new Date().toISOString();
      delay = Math.min(300000, this.idleMs * 2 ** this.failures);
      const retryAt = /retry_at=(\S+)$/.exec(this.status.lastError)?.[1];
      if (retryAt)
        delay = Math.max(delay, Date.parse(retryAt) - Date.now() || 0);
    }
    this.schedule(delay);
  }
}
