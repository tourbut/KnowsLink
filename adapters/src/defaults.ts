// Release defaults (T-03) and the automatic agent folder (R-03), applied before any other module reads the environment.
import { readFileSync } from "node:fs";
import { join } from "node:path";

const env = process.env;
// A release install (install_bot_mcp.sh) writes defaults.json beside dist/; a set variable always wins, so existing installs keep their values.
try {
  const d = JSON.parse(
    readFileSync(new URL("../defaults.json", import.meta.url), "utf8"),
  ) as Record<string, unknown>;
  for (const k of [
    "KNOWSLINK_MODE",
    "RELAY_URL",
    "KNOWSLINK_HOST",
    "KNOWSLINK_AGENT_ROOT",
  ])
    if (env[k] === undefined && typeof d[k] === "string") env[k] = d[k];
} catch {
  // No defaults.json: dev builds and host plugins stay held unless configured.
}
const root =
  env.KNOWSLINK_AGENT_ROOT ??
  (process.platform === "win32" && env.LOCALAPPDATA
    ? join(env.LOCALAPPDATA, "KnowsLink")
    : undefined);
const host = ["claude", "codex", "grok"].includes(env.KNOWSLINK_HOST ?? "")
  ? env.KNOWSLINK_HOST
  : "local";
// One fixed folder per host: after an MCP restart status and receive find the same agent, and a retry reuses or sets it aside (R-01).
export const defaultFolder =
  env.KNOWSLINK_MODE === "public-node" && !env.KNOWSLINK_AGENT_FOLDER && !!root;
if (defaultFolder) env.KNOWSLINK_AGENT_FOLDER = join(root!, `${host}-agent`);
