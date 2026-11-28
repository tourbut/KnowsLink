// Orca bot host on this PC (SAR-ORCA-HOST-001/002): one watch per agent folder, fixed install and per-user autostart.
// Nothing here prints keys or message text: only paths, IDs, pids and states.
import { execFile, spawn } from "node:child_process";
import { readFileSync, rmSync } from "node:fs";
import {
  copyFile,
  mkdir,
  open,
  readFile,
  rm,
  stat,
  writeFile,
} from "node:fs/promises";
import { uptime } from "node:os";
import { basename, dirname, join, resolve } from "node:path";
import { promisify } from "node:util";
import { root } from "./defaults.js";
import { alive } from "./session-lock.js";
import { memberTransport } from "./text.js";

export const execAsync = promisify(execFile);
// No control characters and no leading dash, so the value can only be one argv operand.
export const argValue = (v: string | undefined, max: number): v is string =>
  !!v && [...v].length <= max && !/[\p{Cc}]/u.test(v) && !v.startsWith("-");
export const flag = (argv: string[], name: string) => {
  const at = argv.indexOf(name);
  return at >= 0 ? argv[at + 1] : undefined;
};
export const say = (event: unknown) =>
  process.stdout.write(JSON.stringify(event) + "\n");
const quiet = { shell: false, windowsHide: true, timeout: 20000 };

// The image path and command line of a pid, or null when it is gone or cannot be read.
export type ProcessInfo = { exe: string; cmd: string };
export async function processInfo(pid: number): Promise<ProcessInfo | null> {
  try {
    if (process.platform === "win32") {
      const { stdout } = await execAsync(
        "powershell.exe",
        [
          "-NoProfile",
          "-NonInteractive",
          "-Command",
          `Get-CimInstance Win32_Process -Filter 'ProcessId=${pid}' | Select-Object ExecutablePath,CommandLine | ConvertTo-Json -Compress`,
        ],
        quiet,
      );
      const p = JSON.parse(stdout) as {
        ExecutablePath?: string;
        CommandLine?: string;
      };
      return p.ExecutablePath && p.CommandLine
        ? { exe: p.ExecutablePath, cmd: p.CommandLine }
        : null;
    }
    const { stdout } = await execAsync(
      "ps",
      ["-o", "command=", "-p", String(pid)],
      quiet,
    );
    const cmd = stdout.trim();
    return cmd ? { exe: cmd.split(/\s+/)[0]!, cmd } : null;
  } catch {
    return null;
  }
}
// A pid is a watch only if its image is node and its command line runs `orca watch`.
// shortcut: the folder argument is not matched (relative or differently cased paths); match it if one PC runs watches for several folders.
export const isWatch = (p: ProcessInfo | null) =>
  !!p &&
  /^node(\.exe)?$/i.test(basename(p.exe.replace(/\\/g, "/"))) &&
  /\borca"?\s+"?watch\b/.test(p.cmd);

// One watch per agent folder (SAR-ORCA-HOST-001): the watch is the folder's only automatic receiver.
// A lock from before this boot, of a dead pid, or of a pid now running something else is stale: never kill it.
type WatchLock = { pid: number; since: string };
export const watchLock = (folder: string) => join(folder, "orca-watch.lock");
export async function liveWatch(
  folder: string,
  probe: (pid: number) => Promise<ProcessInfo | null> = processInfo,
): Promise<WatchLock | null> {
  try {
    const l = JSON.parse(
      await readFile(watchLock(folder), "utf8"),
    ) as WatchLock;
    const boot = Date.now() - uptime() * 1000;
    return Number.isInteger(l.pid) &&
      Date.parse(l.since) > boot &&
      alive(l.pid) &&
      (l.pid === process.pid || isWatch(await probe(l.pid)))
      ? l
      : null;
  } catch {
    return null;
  }
}
// Resolves the live holder that blocks this watch, or null once this process holds the lock.
// The lock goes with this process on a normal exit or Ctrl+C; a killed watch leaves it, and stop removes it.
export async function claimWatch(
  folder: string,
  probe?: (pid: number) => Promise<ProcessInfo | null>,
): Promise<WatchLock | null> {
  const mine = JSON.stringify({
    pid: process.pid,
    since: new Date().toISOString(),
  });
  for (let i = 0; i < 2; i++) {
    try {
      await writeFile(watchLock(folder), mine, { flag: "wx", mode: 0o600 });
      process.once("exit", () => releaseWatch(folder));
      for (const sig of ["SIGINT", "SIGTERM"] as const)
        process.once(sig, () => process.exit(130));
      return null;
    } catch (e) {
      if ((e as NodeJS.ErrnoException).code !== "EEXIST") throw e;
    }
    const held = await liveWatch(folder, probe);
    if (held) return held;
    await rm(watchLock(folder), { force: true });
  }
  throw new Error("orca watch already running");
}
// Removes the lock only while this process holds it; sync, so it also runs inside an exit handler.
export function releaseWatch(folder: string) {
  try {
    const l = JSON.parse(readFileSync(watchLock(folder), "utf8")) as WatchLock;
    if (l.pid === process.pid) rmSync(watchLock(folder), { force: true });
  } catch {
    // No lock, or not ours: nothing to release.
  }
}

// Fixed install (SAR-ORCA-HOST-001): bundle, node and defaults.json under <root>\orca, so nothing depends on a worktree.
// defaults.json carries only modes, the relay URL, the folder path and the bot name; keys stay in the agent folder.
export function orcaDefaults(folder: string, bot: string, relay: string) {
  return {
    KNOWSLINK_MODE: "public-node",
    RELAY_URL: relay,
    KNOWSLINK_HOST: "orca",
    KNOWSLINK_AGENT_FOLDER: folder,
    KNOWSLINK_ORCA_BOT: bot,
    // The watch is the one automatic receiver; the bot's MCP reads the shared inbox and pulls only when asked.
    KNOWSLINK_AUTO_RECEIVE: "off",
  };
}
// shortcut: a path with spaces is quoted, and Orca's tools-add parsing of quoted paths is unchecked; check it if a profile path has spaces.
const quote = (p: string) => (/\s/.test(p) ? `"${p}"` : p);
// The one line for the bot chat: a command already holding double quotes goes in single quotes (PowerShell and sh literal).
export const toolsAdd = (mcp: string) =>
  `orca bot tools add ${mcp.includes('"') ? `'${mcp}'` : `"${mcp}"`} --name knowslink`;
export const installedCommand = (dir: string, ...args: string[]) =>
  [join(dir, "node.exe"), join(dir, "dist", "plugin.js"), ...args]
    .map(quote)
    .join(" ");
const RUN_KEY = "HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Run";
export const RUN_VALUE = "KnowsLinkOrcaWatch";
// Per-user logon start without admin rights: reg.exe argv for the HKCU Run value, or its removal.
export const autostartArgs = (command?: string) =>
  command
    ? ["add", RUN_KEY, "/v", RUN_VALUE, "/t", "REG_SZ", "/d", command, "/f"]
    : ["delete", RUN_KEY, "/v", RUN_VALUE, "/f"];
const installed = () => dirname(dirname(resolve(process.argv[1]!)));
export const configured = () => {
  const folder = process.env.KNOWSLINK_AGENT_FOLDER;
  const bot = process.env.KNOWSLINK_ORCA_BOT;
  if (!folder || !argValue(bot, 80)) throw new Error("usage");
  return { folder, bot };
};

// Starts the configured watch detached and hidden, unless one already holds the folder.
export async function start() {
  const { folder, bot } = configured();
  const held = await liveWatch(folder);
  if (held) return say({ state: "running", pid: held.pid });
  // shortcut: watch.log only grows (IDs, kinds and sizes per message); rotate it if it ever gets large.
  const log = await open(join(installed(), "watch.log"), "a", 0o600);
  try {
    const child = spawn(
      process.execPath,
      [process.argv[1]!, "orca", "watch", folder, "--bot", bot],
      { detached: true, windowsHide: true, stdio: ["ignore", log.fd, log.fd] },
    );
    child.unref();
    say({ state: "started", pid: child.pid });
  } finally {
    await log.close();
  }
}

export async function stop(folder: string) {
  const held = await liveWatch(folder);
  if (!held) return "not_running";
  process.kill(held.pid);
  await rm(watchLock(folder), { force: true });
  return "stopped";
}

export async function install(folder: string, argv: string[]) {
  if (process.platform !== "win32")
    throw new Error("orca install is Windows only");
  const bot = flag(argv, "--bot");
  if (!argValue(bot, 80)) throw new Error("invalid orca bot or handle");
  // A connected folder only: memberTransport checks the private files; the relay URL comes from the folder itself.
  const t = await memberTransport(folder);
  const { relay } = JSON.parse(
    await readFile(join(folder, "agent.json"), "utf8"),
  ) as { relay: string };
  const dir = flag(argv, "--dir") ?? (root && join(root, "orca"));
  if (!dir) throw new Error("usage");
  for (const f of new Set([folder, process.env.KNOWSLINK_AGENT_FOLDER]))
    if (f) await stop(f);
  await mkdir(join(dir, "dist"), { recursive: true });
  const plugin = join(dir, "dist", "plugin.js");
  if (resolve(process.argv[1]!) !== resolve(plugin))
    await copyFile(process.argv[1]!, plugin);
  const node = join(dir, "node.exe");
  // A running bot MCP keeps node.exe open, so the same build is not copied again.
  const size = (p: string) =>
    stat(p).then(
      (s) => s.size,
      () => -1,
    );
  if (
    resolve(process.execPath) !== resolve(node) &&
    (await size(process.execPath)) !== (await size(node))
  )
    await copyFile(process.execPath, node);
  await writeFile(
    join(dir, "defaults.json"),
    JSON.stringify(orcaDefaults(resolve(folder), bot, relay), null, 2) + "\n",
  );
  await execAsync(
    "reg.exe",
    autostartArgs(installedCommand(dir, "orca", "start")),
    quiet,
  );
  const started = (await execAsync(node, [plugin, "orca", "start"], quiet)) as {
    stdout: string;
  };
  const mcp = installedCommand(dir);
  say({
    installed: dir,
    agent: t.agentId,
    autostart: RUN_VALUE,
    watch: JSON.parse(started.stdout.trim()) as unknown,
    mcp,
    toolsAdd: toolsAdd(mcp),
  });
}

// Off also ends the running watch, so nothing receives until autostart or start runs again.
export async function autostart(mode: "on" | "off") {
  const { folder } = configured();
  await execAsync(
    "reg.exe",
    autostartArgs(
      mode === "on"
        ? installedCommand(installed(), "orca", "start")
        : undefined,
    ),
    quiet,
  );
  say({
    autostart: mode,
    watch: mode === "off" ? await stop(folder) : undefined,
  });
}
