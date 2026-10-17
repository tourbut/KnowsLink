// A second session never silently takes an agent address; a dead or explicitly replaced holder loses it.
import assert from "node:assert/strict";
import { mkdtemp, readFile, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join } from "node:path";
import {
  claimSession,
  holdsSession,
  lockPath,
  readLock,
  releaseSession,
} from "./session-lock.js";

const folder = await mkdtemp(join(tmpdir(), "kl-session-lock-"));
assert.equal(await claimSession(folder, "sessionA"), null);
assert.equal(await holdsSession(folder, "sessionA"), true);
assert.equal(await holdsSession(folder, "sessionB"), false);
// Same session claiming again (mod re-ran on) keeps the lock.
assert.equal(await claimSession(folder, "sessionA"), null);
// A live holder of another alias blocks the claim without --take: this process is the live holder.
const busy = await claimSession(folder, "sessionB");
assert.equal(busy?.session, "sessionA");
assert.equal(
  JSON.parse(await readFile(lockPath(folder), "utf8")).pid,
  process.pid,
);
// Explicit takeover replaces the lock; the old holder no longer holds it.
assert.equal(await claimSession(folder, "sessionB", true), null);
assert.equal(await holdsSession(folder, "sessionA"), false);
assert.equal(await holdsSession(folder, "sessionB"), true);
// Release removes only our own lock.
await releaseSession(folder, "sessionA");
assert.equal((await readLock(folder))?.session, "sessionB");
await releaseSession(folder, "sessionB");
assert.equal(await readLock(folder), null);
// A dead holder (pid that cannot exist) is taken over without --take.
await writeFile(
  lockPath(folder),
  JSON.stringify({
    session: "ghost",
    pid: 2147483646,
    since: "2026-01-01T00:00:00Z",
  }),
);
assert.equal(await claimSession(folder, "sessionC"), null);
assert.equal(await holdsSession(folder, "sessionC"), true);
// Corrupt lock content is replaced, never trusted.
await writeFile(lockPath(folder), "not json");
assert.equal(await readLock(folder), null);
assert.equal(await claimSession(folder, "sessionD"), null);
await assert.rejects(
  claimSession(folder, "bad alias!"),
  /invalid session alias/,
);
process.stdout.write(
  "PASS: session lock claim/busy/take/release/dead holder/corrupt\n",
);
