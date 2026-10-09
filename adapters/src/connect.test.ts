// Local onboarding tests verify PoP bytes, private-file permissions, one-use completion and safe URL selection.
import assert from "node:assert/strict";
import { execFileSync } from "node:child_process";
import { createPublicKey, verify } from "node:crypto";
import { chmod, lstat, mkdtemp, readFile, rm } from "node:fs/promises";
import { createServer } from "node:http";
import { tmpdir } from "node:os";
import { join } from "node:path";
import { complete, prepare, relayBase } from "./connect.js";
import { privatePath } from "./private-files.js";

async function assertPrivate(path: string, mode: number) {
  if (process.platform === "win32") await privatePath(path, mode === 0o700);
  else assert.equal((await lstat(path)).mode & 0o777, mode);
}
async function expose(path: string, open: boolean) {
  if (process.platform === "win32")
    execFileSync(
      "icacls.exe",
      [path, open ? "/grant" : "/remove:g", open ? "*S-1-1-0:(R)" : "*S-1-1-0"],
      { stdio: "ignore" },
    );
  else await chmod(path, open ? 0o644 : 0o600);
}

for (const url of [
  "http://example.com",
  "https://user:secret@example.com",
  "https://example.com/path",
  "https://example.com/?token=x",
  "file:///tmp/key",
])
  assert.throws(() => relayBase(url));
assert.equal(relayBase("https://example.com/"), "https://example.com");
const root = await mkdtemp(join(tmpdir(), "kl-connect-"));
const token = "a".repeat(43);
let pending: Record<string, string> = {};
let consumed = false;
const server = createServer(async (req, res) => {
  let raw = "";
  for await (const chunk of req) raw += chunk;
  const body: Record<string, string> = JSON.parse(raw);
  assert.equal(body.token, token);
  assert.equal(body.client, "node-local");
  res.setHeader("Content-Type", "application/json");
  if (req.url === "/v1/connect/info") {
    res.end(
      JSON.stringify({
        owner: "owner",
        agent: "agent_a",
        client: "node-local",
        mode: "register",
        state: "waiting",
        exp: new Date(Date.now() + 600000).toISOString(),
      }),
    );
  } else {
    if (req.url === "/v1/connect/prepare") pending = body;
    const key = createPublicKey({
      key: Buffer.concat([
        Buffer.from("302a300506032b6570032100", "hex"),
        Buffer.from(pending.public!, "base64url"),
      ]),
      format: "der",
      type: "spki",
    });
    assert(
      verify(
        null,
        Buffer.from(
          [
            "KNOWSLINK-CONNECT",
            token,
            "owner",
            "agent_a",
            "node-local",
            "register",
            pending.kid!,
            pending.public!,
          ].join("\0"),
        ),
        key,
        Buffer.from(body.proof!, "base64url"),
      ),
    );
    if (req.url === "/v1/connect/complete") {
      assert(!consumed);
      consumed = true;
      res.end(
        JSON.stringify({
          agent: "agent_a",
          kid: pending.kid,
          credential: "c".repeat(43),
        }),
      );
    } else res.end(JSON.stringify({ state: "prepared" }));
  }
});
try {
  await new Promise<void>((resolve) => server.listen(0, "127.0.0.1", resolve));
  const address = server.address();
  assert(address && typeof address === "object");
  const folder = join(root, "client");
  const result = await prepare(
    `http://127.0.0.1:${address.port}`,
    folder,
    token,
  );
  assert(result.fingerprint.startsWith("SHA256:"));
  assert(!JSON.stringify(result).includes(token));
  await assertPrivate(folder, 0o700);
  await assertPrivate(join(folder, "private.pem"), 0o600);
  await assert.rejects(
    prepare(`http://127.0.0.1:${address.port}`, folder, token),
  );
  await expose(join(folder, "private.pem"), true);
  await assert.rejects(complete(folder));
  assert(!consumed);
  await expose(join(folder, "private.pem"), false);
  await complete(folder);
  assert(consumed);
  const config = JSON.parse(await readFile(join(folder, "agent.json"), "utf8"));
  assert.equal(config.agent, "agent_a");
  assert.equal(config.credential, "c".repeat(43));
  await assertPrivate(join(folder, "agent.json"), 0o600);
  assert(
    !(await readFile(join(folder, "pending.json"), "utf8")).includes(token),
  );
  await assert.rejects(complete(folder));
  process.stdout.write(
    "PASS: local connection PoP, private directory/files, no overwrite, credential storage, one-use completion, URL boundary\n",
  );
} finally {
  server.close();
  await rm(root, { recursive: true, force: true });
}
