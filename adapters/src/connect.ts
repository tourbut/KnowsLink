// Local Node onboarding keeps private keys and credentials in a new 0700 directory; stdout contains only public status.
import {
  createHash,
  createPublicKey,
  generateKeyPairSync,
  sign,
} from "node:crypto";
import { mkdir, readFile, writeFile, lstat } from "node:fs/promises";
import { join } from "node:path";
import { pathToFileURL } from "node:url";
import { Adapter, relayBase } from "./core.js";

type Grant = {
  owner: string;
  agent: string;
  client: string;
  mode: string;
  state: string;
  exp: string;
};
type Pending = Grant & {
  base: string;
  token: string;
  kid: string;
  public: string;
};
export { relayBase } from "./core.js";
function proof(p: Pending, privateKey: string): string {
  return sign(
    null,
    Buffer.from(
      [
        "KNOWSLINK-CONNECT",
        p.token,
        p.owner,
        p.agent,
        p.client,
        p.mode,
        p.kid,
        p.public,
      ].join("\0"),
    ),
    privateKey,
  ).toString("base64url");
}
function api(base: string): Adapter {
  return new Adapter(base, "", "", "", "");
}
export async function prepare(
  base: string,
  folder: string,
  token: string,
): Promise<{ fingerprint: string; connection: string }> {
  base = relayBase(base);
  if (!/^[A-Za-z0-9_-]{43}$/.test(token))
    throw new Error("invalid connection input");
  const grant = await api(base).request<Grant>("/v1/connect/info", {
    token,
    client: "node-local",
  });
  if (
    grant.client !== "node-local" ||
    grant.state !== "waiting" ||
    !["register", "rotate"].includes(grant.mode) ||
    !/^[a-z][a-z0-9_:-]*$/.test(grant.agent) ||
    !Number.isFinite(Date.parse(grant.exp)) ||
    Date.parse(grant.exp) <= Date.now()
  )
    throw new Error("invalid grant");
  await mkdir(folder, { mode: 0o700 }); // Existing directories fail; never overwrite a working key.
  const keys = generateKeyPairSync("ed25519");
  const privateKey = keys.privateKey
    .export({ format: "pem", type: "pkcs8" })
    .toString();
  const publicKey = createPublicKey(privateKey)
    .export({ format: "der", type: "spki" })
    .subarray(-32);
  const p: Pending = {
    ...grant,
    base,
    token,
    kid: `key_${Date.now()}`,
    public: publicKey.toString("base64url"),
  };
  await writeFile(join(folder, "private.pem"), privateKey, {
    mode: 0o600,
    flag: "wx",
  });
  await writeFile(join(folder, "pending.json"), JSON.stringify(p), {
    mode: 0o600,
    flag: "wx",
  });
  await api(base).request("/v1/connect/prepare", {
    token,
    client: p.client,
    kid: p.kid,
    public: p.public,
    proof: proof(p, privateKey),
  });
  return {
    fingerprint: `SHA256:${createHash("sha256").update(publicKey).digest("base64url")}`,
    connection: `${base}/home/connections/${createHash("sha256").update(token).digest("base64url")}`,
  };
}
export async function complete(folder: string): Promise<void> {
  const directory = await lstat(folder);
  if (
    !directory.isDirectory() ||
    directory.isSymbolicLink() ||
    (directory.mode & 0o077) !== 0 ||
    directory.uid !== process.getuid?.()
  )
    throw new Error("private directory required");
  for (const name of ["private.pem", "pending.json"]) {
    const st = await lstat(join(folder, name));
    if (
      !st.isFile() ||
      st.isSymbolicLink() ||
      (st.mode & 0o077) !== 0 ||
      st.uid !== process.getuid?.()
    )
      throw new Error("private file required");
  }
  const p: Pending = JSON.parse(
    await readFile(join(folder, "pending.json"), "utf8"),
  );
  const privateKey = await readFile(join(folder, "private.pem"), "utf8");
  const connected = await api(relayBase(p.base)).request<{
    agent: string;
    kid: string;
    credential: string;
  }>("/v1/connect/complete", {
    token: p.token,
    client: p.client,
    proof: proof(p, privateKey),
  });
  await writeFile(
    join(folder, "agent.json"),
    JSON.stringify({
      relay: p.base,
      ...connected,
      keyFile: join(folder, "private.pem"),
    }),
    { mode: 0o600, flag: "wx" },
  );
  await writeFile(
    join(folder, "pending.json"),
    JSON.stringify({ state: "consumed" }),
    { mode: 0o600 },
  );
}
async function main(): Promise<void> {
  const [, , action, first, folder] = process.argv;
  if (action === "prepare" && first && folder) {
    let input = "";
    for await (const chunk of process.stdin) {
      input += String(chunk);
      if (input.length > 128) throw new Error("input too large");
    }
    const v = await prepare(first, folder, input.trim());
    process.stdout.write(
      `상태: 공개키 준비 완료\n지문: ${v.fingerprint}\n자기 브라우저에서 확인: ${v.connection}\n`,
    );
  } else if (action === "complete" && first) {
    await complete(first);
    process.stdout.write(
      "상태: 연결 완료. agent.json·private.pem은 이 클라이언트에만 보관하세요. idle pull은 10초 이상 간격입니다.\n",
    );
  } else throw new Error("usage");
}
if (
  process.argv[1] &&
  import.meta.url === pathToFileURL(process.argv[1]).href
) {
  main().catch(() => {
    process.stderr.write(
      "연결 실패. 아직 사용할 새 자격이 없습니다. 브라우저 상태·기한을 확인하세요. 완료 응답을 잃었다면 새 연결·회전으로 복구하세요.\n",
    );
    process.exitCode = 1;
  });
}
