// Release defaults (T-03) and the automatic agent folder (R-03), applied before any other module reads the environment.
import { readFileSync } from "node:fs";
import { join } from "node:path";

const env = process.env;
// No credential may ever be listed here: defaults.json holds only modes, the relay URL and paths.
export const DEFAULT_KEYS = [
  "KNOWSLINK_MODE",
  "RELAY_URL",
  "KNOWSLINK_HOST",
  "KNOWSLINK_AGENT_ROOT",
  // Orca fixed install (SAR-ORCA-HOST-001): one existing folder, the bot to ring, and MCP reading without its own receive loop.
  "KNOWSLINK_AGENT_FOLDER",
  "KNOWSLINK_ORCA_BOT",
  "KNOWSLINK_AUTO_RECEIVE",
] as const;
// A release install (install_bot_mcp.sh) writes defaults.json beside dist/; a set variable always wins, so existing installs keep their values.
try {
  const d = JSON.parse(
    readFileSync(new URL("../defaults.json", import.meta.url), "utf8"),
  ) as Record<string, unknown>;
  for (const k of DEFAULT_KEYS)
    if (env[k] === undefined && typeof d[k] === "string") env[k] = d[k];
} catch {
  // No defaults.json: dev builds and host plugins stay held unless configured.
}
export const root =
  env.KNOWSLINK_AGENT_ROOT ??
  (process.platform === "win32" && env.LOCALAPPDATA
    ? join(env.LOCALAPPDATA, "KnowsLink")
    : undefined);
const host = ["claude", "codex", "grok", "orca"].includes(
  env.KNOWSLINK_HOST ?? "",
)
  ? env.KNOWSLINK_HOST
  : "local";
// One fixed folder per host: after an MCP restart status and receive find the same agent, and a retry reuses or sets it aside (R-01).
export const defaultFolder =
  env.KNOWSLINK_MODE === "public-node" && !env.KNOWSLINK_AGENT_FOLDER && !!root;
if (defaultFolder) env.KNOWSLINK_AGENT_FOLDER = join(root!, `${host}-agent`);
