// Local beta operator prepares two trial identities, pairing and private files; owner credentials never reach the remote agent.
import { generateKeyPairSync, sign } from "node:crypto";
import { mkdir, writeFile } from "node:fs/promises";
import { resolve, join } from "node:path";
import { pathToFileURL } from "node:url";
import { Adapter } from "./core.js";

export async function setupTrial(base: string) {
  const url = new URL(base);
  if (
    url.protocol !== "http:" ||
    !["127.0.0.1", "localhost", "[::1]"].includes(url.hostname) ||
    url.pathname !== "/" ||
    url.username ||
    url.password ||
    url.search ||
    url.hash
  )
    throw new Error("local operator URL required");
  const api = new Adapter(url.origin, "", "unused", "key1", "");
  const agents = [];
  for (const agent of ["trial_codex", "trial_grok"]) {
    const owner = await api.request<{ owner: string; credential: string }>(
      "/v1/owners",
      {},
    );
    const { privateKey, publicKey } = generateKeyPairSync("ed25519");
    const publicBytes = publicKey
      .export({ type: "spki", format: "der" })
      .subarray(-32)
      .toString("base64url");
    const proof = sign(
      null,
      Buffer.from(
        ["KNOWSLINK-KEY-POP", owner.owner, agent, "key1", publicBytes].join(
          "\0",
        ),
      ),
      privateKey,
    ).toString("base64url");
    const ownerAPI = new Adapter(
      url.origin,
      owner.credential,
      agent,
      "key1",
      "",
    );
    const registration = await ownerAPI.request<{ credential: string }>(
      "/v1/agents",
      { agent, kid: "key1", public: publicBytes, proof },
    );
    agents.push({
      agent,
      owner,
      credential: registration.credential,
      pem: privateKey.export({ type: "pkcs8", format: "pem" }).toString(),
    });
  }
  const a = agents[0]!,
    b = agents[1]!;
  await new Adapter(url.origin, a.credential, a.agent, "key1", a.pem).request(
    "/v1/invites",
    { agent: a.agent, target: b.agent },
  );
  await new Adapter(
    url.origin,
    b.owner.credential,
    b.agent,
    "key1",
    b.pem,
  ).request("/v1/invite-decision", {
    agent: a.agent,
    target: b.agent,
    decision: "accept",
  });
  return { a, b };
}
async function main() {
  const [base, directory] = process.argv.slice(2);
  if (!base || !directory || !directory.startsWith("/"))
    throw new Error("absolute destination required");
  // Refuse reuse before registering agents. Any partial files remain private for operator recovery.
  await mkdir(directory, { mode: 0o700 });
  const { a, b } = await setupTrial(base);
  for (const [own, peer] of [
    [a, b],
    [b, a],
  ]) {
    const folder = join(directory, own!.agent);
    await mkdir(folder, { mode: 0o700 });
    const key = join(folder, "key.pem");
    await writeFile(key, own!.pem, { mode: 0o600, flag: "wx" });
    await writeFile(
      join(folder, "environment.json"),
      JSON.stringify(
        {
          KNOWSLINK_MODE: "test-remote",
          RELAY_URL: "https://link.knowslog.com",
          AGENT_ID: own!.agent,
          AGENT_CREDENTIAL: own!.credential,
          AGENT_KID: "key1",
          AGENT_KEY_FILE: key,
          KNOWSLINK_TEST_PEER: peer!.agent,
        },
        null,
        2,
      ),
      { mode: 0o600, flag: "wx" },
    );
  }
  await writeFile(
    join(directory, "owners.json"),
    JSON.stringify({ a: a.owner, b: b.owner }),
    { mode: 0o600, flag: "wx" },
  );
  console.info(
    "Prepared trial_codex and trial_grok in private destination; provision separate Access service tokens before remote use. Owner records stay local.",
  );
}
if (
  process.argv[1] &&
  import.meta.url === pathToFileURL(resolve(process.argv[1])).href
)
  main().catch(() => {
    console.error(
      "Trial setup failed; inspect private state locally, do not rerun blindly or print credentials",
    );
    process.exitCode = 1;
  });
