// Local Google onboarding verifies device/complete PoP, automatic private storage, expiry and independent client keys.
import assert from "node:assert/strict";
import { createPublicKey, verify } from "node:crypto";
import { mkdtemp, readFile, rm, writeFile } from "node:fs/promises";
import { createServer } from "node:http";
import { tmpdir } from "node:os";
import { join } from "node:path";
import { beginLogin, finishLogin } from "./login.js";
import { memberTransport } from "./text.js";

const root = await mkdtemp(join(tmpdir(), "kl-login-"));
const devices = new Map<string, Record<string, string>>();
let approved = false;
let completed = 0;
const server = createServer(async (req, res) => {
  let raw = "";
  for await (const chunk of req) raw += chunk;
  const body: Record<string, string> = JSON.parse(raw);
  if (req.url === "/v1/connect/start") devices.set(body.token!, body);
  const p = devices.get(body.token!);
  assert(p);
  const key = createPublicKey({
    key: Buffer.concat([
      Buffer.from("302a300506032b6570032100", "hex"),
      Buffer.from(p.public!, "base64url"),
    ]),
    format: "der",
    type: "spki",
  });
  const complete = req.url === "/v1/connect/complete";
  const agent = `agent_${[...devices.keys()].indexOf(body.token!)}`;
  const fields = complete
    ? [
        "KNOWSLINK-CONNECT",
        body.token!,
        "owner",
        agent,
        "node-local",
        "register",
        p.kid!,
        p.public!,
      ]
    : ["KNOWSLINK-DEVICE", body.token!, "node-local", p.kid!, p.public!];
  assert(
    verify(
      null,
      Buffer.from(fields.join("\0")),
      key,
      Buffer.from(body.proof!, "base64url"),
    ),
  );
  res.setHeader("Content-Type", "application/json");
  if (complete) {
    assert(approved);
    completed++;
    res.end(
      JSON.stringify({
        agent,
        kid: p.kid,
        credential: "c".repeat(43),
        state: "connected",
      }),
    );
  } else
    res.end(
      JSON.stringify({
        state:
          req.url === "/v1/connect/start"
            ? "requested"
            : approved
              ? "approved"
              : "prepared",
        exp: new Date(Date.now() + 600000).toISOString(),
        owner: "owner",
        agent,
      }),
    );
});
try {
  await new Promise<void>((resolve) => server.listen(0, "127.0.0.1", resolve));
  const address = server.address();
  assert(address && typeof address === "object");
  const base = `http://127.0.0.1:${address.port}`;
  const first = join(root, "new-parent", "first");
  const second = join(root, "second");
  const login = await beginLogin(base, first);
  assert(login.url.startsWith(`${base}/connect/`));
  assert(login.fingerprint.startsWith("SHA256:"));
  for (const secret of devices.keys())
    assert(!JSON.stringify(login).includes(secret));
  assert.equal(await finishLogin(first), "waiting");
  assert.equal(completed, 0);
  await assert.rejects(beginLogin(base, first));
  const other = await beginLogin(base, second);
  assert.notEqual(login.fingerprint, other.fingerprint);
  approved = true;
  assert.equal(await finishLogin(first), "connected");
  assert.equal(completed, 1);
  await memberTransport(first);
  await assert.rejects(finishLogin(first));
  const config = JSON.parse(await readFile(join(first, "agent.json"), "utf8"));
  assert.equal(config.agent, "agent_0");
  assert.equal(config.credential, "c".repeat(43));
  assert(
    !(await readFile(join(first, "login.json"), "utf8")).includes(
      [...devices.keys()][0]!,
    ),
  );
  const pending = JSON.parse(
    await readFile(join(second, "login.json"), "utf8"),
  );
  pending.exp = new Date(Date.now() - 1).toISOString();
  await writeFile(join(second, "login.json"), JSON.stringify(pending), {
    mode: 0o600,
  });
  await assert.rejects(finishLogin(second), /expired/);
  assert.equal(completed, 1);
  process.stdout.write(
    "PASS: browser URL excludes private token, key-bound poll/complete, explicit approval, automatic local storage, distinct client keys, no overwrite, expiry\n",
  );
} finally {
  server.closeAllConnections();
  server.close();
  await rm(root, { recursive: true, force: true });
}
