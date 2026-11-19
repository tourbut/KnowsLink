// Independent filesystem error matrix; fresh synthetic fixtures stay on disk.
import assert from 'node:assert/strict';
import fs, { mkdtemp, readFile, stat, writeFile } from 'node:fs/promises';
import { syncBuiltinESMExports } from 'node:module';
import { tmpdir } from 'node:os';
import { join, resolve } from 'node:path';
import { pathToFileURL } from 'node:url';
import { mock } from 'node:test';

const product = resolve(process.argv[2]);
const { Inbox } = await import(pathToFileURL(join(product, 'adapters/dist/inbox.js')));
const root = await mkdtemp(join(tmpdir(), 'kl-qa-refactor002-'));
const inbox = await Inbox.open(root);
const marker = join(root, 'inbox', 'fresh.tmp');
await writeFile(marker, '', { mode: 0o600 });
const messages = [1, 2].map(n => ({
  id: `00000000-0000-7000-8000-${String(n).padStart(12, '0')}`,
  from: 'agent_a', to: 'agent_b', text: `independent unread ${n}`, untrusted: true,
  exp: new Date(Date.now() + 60000).toISOString(),
}));
let failure;
let hits = 0;
const originalStat = stat;
mock.method(fs, 'stat', async path => {
  if (path === marker && failure) {
    hits++;
    throw Object.assign(new Error('synthetic stat failure'), { code: failure });
  }
  return originalStat(path);
});
syncBuiltinESMExports();
try {
  for (const message of messages) assert.equal(await inbox.store(message), true);
  for (const code of ['EACCES', 'EIO', 'EPERM', 'EACCES']) {
    failure = code;
    await assert.rejects(inbox.take(), { code });
    assert.equal((await inbox.pending()).length, 2);
    for (const message of messages) {
      assert.deepEqual(JSON.parse(await readFile(join(root, 'inbox', `${message.id}.json`), 'utf8')), message);
    }
  }
  failure = 'ENOENT';
  assert.equal((await inbox.take()).text, messages[0].text);
  failure = undefined;
  assert.equal((await inbox.take()).text, messages[1].text);
  for (const message of messages) {
    assert.equal(await inbox.store(message), false);
    assert.equal(await readFile(join(root, 'inbox', `${message.id}.read`), 'utf8'), '');
  }
  failure = 'EIO';
  const prior = hits;
  assert.equal(await inbox.take(), null);
  assert.equal(hits, prior, 'empty inbox performs no housekeeping');
  console.log('PASS: repeated EACCES/EIO/EPERM preserve both bodies; ENOENT receives; retry order; empty markers; duplicate; empty inbox');
} finally {
  mock.restoreAll();
  syncBuiltinESMExports();
  console.log(`Preserved synthetic fixture: ${root}`);
}
