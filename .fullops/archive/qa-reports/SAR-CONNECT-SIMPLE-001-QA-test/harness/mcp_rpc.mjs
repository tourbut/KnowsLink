// Line-protocol MCP stdio client: stdin lines `{"tool":..,"args":{..}}` or `{"list":true}` -> one JSON line each. argv: plugin.js, env JSON.
import { Client } from "file:///C:/Users/shin/orca/workspaces/KnowsLink/tester/adapters/node_modules/@modelcontextprotocol/sdk/dist/esm/client/index.js";
import { StdioClientTransport } from "file:///C:/Users/shin/orca/workspaces/KnowsLink/tester/adapters/node_modules/@modelcontextprotocol/sdk/dist/esm/client/stdio.js";
import { createInterface } from "node:readline";
const [, , plugin, envJson] = process.argv;
const env = { PATH: process.env.PATH ?? "", SystemRoot: process.env.SystemRoot ?? "", ...JSON.parse(envJson || "{}") };
const t = new StdioClientTransport({ command: process.execPath, args: [plugin], env });
const c = new Client({ name: "qa", version: "0" });
await c.connect(t);
const rl = createInterface({ input: process.stdin });
for await (const line of rl) {
  if (!line.trim()) continue;
  const q = JSON.parse(line);
  try {
    const r = q.list ? await c.listTools() : await c.callTool({ name: q.tool, arguments: q.args ?? {} });
    console.log(JSON.stringify(r));
  } catch (e) { console.log(JSON.stringify({ rpcError: String(e.message ?? e) })); }
}
await c.close();
