// Google browser consent connects one local key; polling and credential storage never expose secrets through stdout or MCP.
import {
  createHash,
  generateKeyPairSync,
  randomBytes,
  sign,
} from "node:crypto";
import { mkdir, readFile, writeFile } from "node:fs/promises";
import { join, basename, dirname } from "node:path";
import { pathToFileURL } from "node:url";
import { setTimeout } from "node:timers/promises";
import { Adapter, relayBase } from "./core.js";
import { complete } from "./connect.js";
import { privateDirectory, privatePath } from "./private-files.js";

type Login = {
  base: string;
  token: string;
  client: string;
  kid: string;
  public: string;
  exp: string;
};
function deviceProof(p: Login, key: string): string {
  return sign(
    null,
    Buffer.from(
      ["KNOWSLINK-DEVICE", p.token, p.client, p.kid, p.public].join("\0"),
    ),
    key,
  ).toString("base64url");
}
export async function beginLogin(base: string, folder: string) {
  base = relayBase(base);
  await mkdir(dirname(folder), { recursive: true, mode: 0o700 });
  await privateDirectory(folder);
  const keys = generateKeyPairSync("ed25519");
  const key = keys.privateKey
    .export({ format: "pem", type: "pkcs8" })
    .toString();
  const publicKey = keys.publicKey
    .export({ format: "der", type: "spki" })
    .subarray(-32);
  const p: Login = {
    base,
    token: randomBytes(32).toString("base64url"),
    client: "node-local",
    kid: `key_${randomBytes(16).toString("hex")}`,
    public: publicKey.toString("base64url"),
    exp: new Date(Date.now() + 600000).toISOString(),
  };
  await writeFile(join(folder, "private.pem"), key, {
    mode: 0o600,
    flag: "wx",
  });
  await writeFile(join(folder, "login.json"), JSON.stringify(p), {
    mode: 0o600,
    flag: "wx",
  });
  const response = await new Adapter(base, "", "", "", "").request<{
    exp: string;
    state: string;
  }>("/v1/connect/start", {
    token: p.token,
    client: p.client,
    kid: p.kid,
    public: p.public,
    proof: deviceProof(p, key),
  });
  if (
    response.state !== "requested" ||
    !Number.isFinite(Date.parse(response.exp)) ||
    Date.parse(response.exp) <= Date.now()
  )
    throw new Error("invalid login response");
  return {
    state: "waiting",
    url: `${base}/connect/${createHash("sha256").update(p.token).digest("base64url")}`,
    fingerprint: `SHA256:${createHash("sha256").update(publicKey).digest("base64url")}`,
    expires: response.exp,
  };
}

export async function finishLogin(
  folder: string,
): Promise<"waiting" | "connected"> {
  await privatePath(folder, true);
  for (const name of ["login.json", "private.pem"])
    await privatePath(join(folder, name));
  const p: Login = JSON.parse(
    await readFile(join(folder, "login.json"), "utf8"),
  );
  if (!Number.isFinite(Date.parse(p.exp)) || Date.parse(p.exp) <= Date.now())
    throw new Error("login expired");
  const key = await readFile(join(folder, "private.pem"), "utf8");
  const response = await new Adapter(
    relayBase(p.base),
    "",
    "",
    "",
    "",
  ).request<{
    state: string;
    owner?: string;
    agent?: string;
  }>("/v1/connect/poll", {
    token: p.token,
    client: p.client,
    proof: deviceProof(p, key),
  });
  if (response.state === "requested" || response.state === "prepared")
    return "waiting";
  if (
    response.state !== "approved" ||
    !response.owner ||
    !response.agent ||
    !/^[a-z][a-z0-9_:-]*$/.test(response.agent)
  )
    throw new Error("invalid approval");
  await writeFile(
    join(folder, "pending.json"),
    JSON.stringify({ ...p, ...response, mode: "register" }),
    { mode: 0o600, flag: "wx" },
  );
  await complete(folder);
  await writeFile(
    join(folder, "login.json"),
    JSON.stringify({ state: "consumed" }),
    { mode: 0o600 },
  );
  return "connected";
}

export async function waitForLogin(folder: string): Promise<void> {
  while ((await finishLogin(folder)) === "waiting") await setTimeout(10000);
}
async function main(): Promise<void> {
  const [, , base, folder] = process.argv;
  if (!base || !folder)
    throw new Error("usage: login <service URL> <new private folder>");
  const login = await beginLogin(base, folder);
  process.stdout.write(
    `자기 브라우저에서 Google 로그인: ${login.url}\n지문: ${login.fingerprint}\n로그인 후 같은 지문을 확인하고 승인하세요. 이 명령이 연결을 자동 저장합니다.\n`,
  );
  await waitForLogin(folder);
  process.stdout.write("연결 완료. 이 클라이언트에만 자격을 저장했습니다.\n");
}
if (
  process.argv[1] &&
  basename(process.argv[1]) === "login.js" &&
  import.meta.url === pathToFileURL(process.argv[1]).href
) {
  main().catch(() => {
    process.stderr.write(
      "연결 실패·만료. 기존 연결은 그대로입니다. 새 비공개 폴더로 다시 시작하세요. 승인 뒤 응답을 잃었다면 홈에서 미사용 agent를 철회하세요.\n",
    );
    process.exitCode = 1;
  });
}
