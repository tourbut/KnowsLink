// Inbox housekeeping failures must preserve unread text; disappearing markers are harmless concurrent activity.
import assert from "node:assert/strict";
import fs, { mkdtemp, readFile, stat, writeFile } from "node:fs/promises";
import { syncBuiltinESMExports } from "node:module";
import { tmpdir } from "node:os";
import { join } from "node:path";
import { mock } from "node:test";
import { Inbox } from "./inbox.js";
import { uuid7 } from "./core.js";
import type { ReceivedText } from "./text.js";

const root = await mkdtemp(join(tmpdir(), "kl-inbox-files-"));
const inbox = await Inbox.open(root);
const marker = join(root, "inbox", "concurrent.tmp");
await writeFile(marker, "", { mode: 0o600 });
let failure: string | undefined;
const originalStat = stat;
mock.method(fs, "stat", async (path: Parameters<typeof stat>[0]) => {
  if (path === marker && failure)
    throw Object.assign(new Error("synthetic housekeeping failure"), {
      code: failure,
    });
  return originalStat(path);
});
syncBuiltinESMExports();

try {
  const message: ReceivedText = {
    id: uuid7(),
    from: "agent_a",
    to: "agent_b",
    text: "synthetic unread text",
    exp: new Date(Date.now() + 60000).toISOString(),
    untrusted: true,
  };
  await inbox.store(message);
  failure = "EACCES";
  await assert.rejects(inbox.take(), { code: "EACCES" });
  assert.deepEqual(await inbox.pending(), [`${message.id}.json`]);
  assert.equal(
    JSON.parse(
      await readFile(join(root, "inbox", `${message.id}.json`), "utf8"),
    ).text,
    message.text,
    "housekeeping failure must leave the original unread body available",
  );
  failure = "ENOENT";
  assert.equal((await inbox.take())?.id, message.id);
  assert.equal(await inbox.take(), null);
  assert.equal(await inbox.store(message), false);
  process.stdout.write(
    "PASS: housekeeping failure preserves text; missing marker permits receive; duplicate stays read\n",
  );
} finally {
  mock.restoreAll();
  syncBuiltinESMExports();
  process.stdout.write(`Preserved synthetic fixture: ${root}\n`);
}
