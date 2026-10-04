// Run one independent observation phase and write a secret-free result file.
import { writeFile } from "node:fs/promises";
import { join } from "node:path";
import { checks, check, clone, failed, scrub } from "./session.mjs";
import { held } from "./held.mjs";
import { closed, live } from "./live.mjs";

const phase = process.argv[2];
const plugin = join(clone, "adapters/dist/plugin.js");
const runner = phase === "held" ? held : phase === "live" ? live : closed;
const ids = await runner(plugin).catch((error) => {
  check("phase exception", false, scrub(error && error.message, []).slice(0, 180));
  return {};
});
await writeFile(process.env.RESULT_PATH, JSON.stringify({ phase, ok: !failed, ids: ids || {}, checks }, null, 2), { mode: 0o600 });
process.exitCode = failed ? 1 : 0;
