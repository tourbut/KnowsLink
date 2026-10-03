// Grok Bot MCP connector: held by default; only explicit loopback synthetic pulls can run.
import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import { localAdapter } from "./core.js";

const server = new McpServer({ name: "knowslink", version: "0.1.0" });
const synthetic = process.env.KNOWSLINK_MODE === "synthetic-loopback";
let busy = false;
const result = (state: string, isError = false) => ({
  content: [
    {
      type: "text" as const,
      text: JSON.stringify({
        state,
        transport: "pull",
        actualConnection: "held",
        webhook: false,
        evidenceFetch: false,
      }),
    },
  ],
  isError,
});
server.registerTool(
  "knowslink_status",
  {
    description:
      "Report KnowsLink connector readiness without reading credentials or contacting relay. Actual Grok Bot connection remains held.",
    inputSchema: {},
    annotations: { readOnlyHint: true, openWorldHint: false },
  },
  async () => result(synthetic ? "synthetic_only" : "held"),
);
server.registerTool(
  "knowslink_pull_once",
  {
    description:
      "Process one synthetic loopback relay message: verify, durable persist, ACK, shared claim, KnowsLink owner gate, minimal denied result. No calendar effects or model-visible payload. Default held.",
    inputSchema: {},
    annotations: {
      readOnlyHint: false,
      destructiveHint: false,
      idempotentHint: false,
      openWorldHint: false,
    },
  },
  async () => {
    if (!synthetic) return result("held", true);
    if (busy) return result("busy", true);
    busy = true;
    try {
      const adapter = await localAdapter();
      if (!adapter) return result("unconfigured", true);
      const processed = await adapter.once(true, () => {});
      return result(processed ? "processed" : "empty");
    } catch {
      // Do not expose credentials, signed bodies, claim tokens, or server error detail to the model.
      return result("failed", true);
    } finally {
      busy = false;
    }
  },
);
server.connect(new StdioServerTransport()).catch(() => {
  console.error("KnowsLink MCP startup failed");
  process.exitCode = 1;
});
