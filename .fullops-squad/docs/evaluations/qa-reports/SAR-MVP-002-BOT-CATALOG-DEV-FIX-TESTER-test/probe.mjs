// Independent stdio probe outside the repo: empty environment, status only.
import { spawn } from "node:child_process";
import { readFileSync, readlinkSync, writeFileSync } from "node:fs";

const [nodePath, pluginPath, cwd, outputPath] = process.argv.slice(2);
const child = spawn("env", ["-i", nodePath, pluginPath], {
  cwd,
  stdio: ["pipe", "pipe", "pipe"],
});
let stdout = Buffer.alloc(0);
let stderr = Buffer.alloc(0);
let environ = "";
try {
  environ = readFileSync(`/proc/${child.pid}/environ`).toString("utf8");
} catch {
  environ = "";
}
const cwdLink = (() => {
  try {
    return readlinkSync(`/proc/${child.pid}/cwd`);
  } catch {
    return "";
  }
})();
child.stdout.on("data", (chunk) => {
  stdout = Buffer.concat([stdout, chunk]);
});
child.stderr.on("data", (chunk) => {
  stderr = Buffer.concat([stderr, chunk]);
});
const send = (message) => child.stdin.write(`${JSON.stringify(message)}\n`);
send({
  jsonrpc: "2.0",
  id: 1,
  method: "initialize",
  params: {
    protocolVersion: "2025-11-25",
    capabilities: {},
    clientInfo: { name: "qa-independent", version: "0" },
  },
});
send({ jsonrpc: "2.0", method: "notifications/initialized" });
send({ jsonrpc: "2.0", id: 2, method: "tools/list", params: {} });
send({
  jsonrpc: "2.0",
  id: 3,
  method: "tools/call",
  params: { name: "knowslink_status", arguments: {} },
});
child.stdin.end();
const code = await new Promise((resolve) => {
  const timer = setTimeout(() => {
    child.kill("SIGKILL");
  }, 20000);
  child.on("close", (value) => {
    clearTimeout(timer);
    resolve(value);
  });
});
const lines = stdout.toString("utf8").split("\n").filter((line) => line.trim());
const nonjson = [];
const messages = [];
for (const line of lines) {
  try {
    messages.push(JSON.parse(line));
  } catch {
    nonjson.push(line);
  }
}
const byId = Object.fromEntries(
  messages.filter((message) => message.id != null).map((message) => [message.id, message]),
);
const names = (byId[2]?.result?.tools ?? []).map((tool) => tool.name).sort();
const status = byId[3]?.result;
const text = status?.content?.[0]?.text ?? "";
const expected =
  '{"state":"held","transport":"pull","actualConnection":"held","webhook":false,"evidenceFetch":false}';
const envKeys = environ.split("\0").filter(Boolean).map((item) => item.split("=", 1)[0]);
const report = {
  code,
  stderr_len: stderr.length,
  stderr: stderr.toString("utf8"),
  nonjson,
  names,
  text,
  isError: status?.isError ?? null,
  protocol: byId[1]?.result?.protocolVersion ?? null,
  cwd,
  cwd_link: cwdLink,
  env_keys: envKeys,
  stdout_len: stdout.length,
};
report.ok =
  code === 0 &&
  nonjson.length === 0 &&
  names.join(",") === "knowslink_pull_once,knowslink_status" &&
  status?.isError === false &&
  text === expected &&
  stderr.length === 0 &&
  envKeys.length === 0;
writeFileSync(outputPath, `${JSON.stringify(report, null, 2)}\n`);
process.exit(report.ok ? 0 : 1);
