// Codex/manual trial CLI uses the same bounded transport as MCP; configuration stays in process environment.
import { newTrialKey, testTransport } from "./test-transport.js";
async function main() {
  const transport = await testTransport();
  if (!transport) throw new Error("held or unconfigured");
  const [action, key] = process.argv.slice(2);
  if (action === "receive") {
    console.info(
      JSON.stringify({
        state: "received_or_empty",
        message: await transport.receive(),
      }),
    );
  } else if (action === "send") {
    const chunks: Buffer[] = [];
    let size = 0;
    for await (const value of process.stdin) {
      const chunk = Buffer.from(value);
      size += chunk.length;
      if (size > 4096) throw new Error("trial input too large");
      chunks.push(chunk);
    }
    console.info(
      JSON.stringify(
        await transport.sendText(
          Buffer.concat(chunks).toString("utf8"),
          key ?? newTrialKey(),
        ),
      ),
    );
  } else throw new Error("use send [idempotency-key] or receive");
}
main().catch(() => {
  console.error(
    "KnowsLink trial failed: check mode, identity, credentials, TTL and Access; no automatic retry",
  );
  process.exitCode = 1;
});
