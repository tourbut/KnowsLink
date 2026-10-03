// Run one loopback synthetic pull; library imports have no CLI side effects.
import { pathToFileURL } from "node:url";
import { localAdapter } from "./core.js";
export * from "./core.js";

async function main(): Promise<void> {
  const adapter = await localAdapter();
  if (!adapter) {
    console.info(
      JSON.stringify({
        state: "unconfigured",
        transport: "pull",
        webhook: false,
        evidenceFetch: false,
      }),
    );
    return;
  }
  await adapter.once(process.env.ADAPTER_GATE === "1");
}
if (
  process.argv[1] &&
  import.meta.url === pathToFileURL(process.argv[1]).href
) {
  main().catch(() => {
    console.error("synthetic adapter failed");
    process.exitCode = 1;
  });
}
