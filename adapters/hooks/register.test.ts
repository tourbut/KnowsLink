// V05: without /knowslink-session on no text tool runs, and a send needs the owner's approval in this session (0 automatic replies).
import { expect, test } from "claude-code/testing";

const SEND = "mcp__plugin_knowslink_knowslink__knowslink_text_send";
const RECEIVE = "mcp__plugin_knowslink_knowslink__knowslink_text_receive";
const reason = (ran: { deny?: string; text?: string }) => ran.deny ?? ran.text;

test("text tools are denied until the owner runs /knowslink-session on", async ($, on) => {
  let reached = 0;
  on("tool.call", () => {
    reached += 1;
    return { result: { content: [] } };
  });
  for (const tool of [SEND, RECEIVE]) {
    const ran = await $.tool.call({
      tool,
      peer: "friend_agent",
      text: "hello",
      idempotency_key: "0123456789abcdef",
      confirmed: true,
    });
    expect(reason(ran)).toMatch(/knowslink-session on/);
  }
  expect(reached).toBe(0);
});

test("/knowslink-session status and off report an unconnected session", async ($, on) => {
  on("session.start", (_$, e) => ({ cwd: e.cwd }));
  on("command.register", (_$, e) => ({ value: { command: e.name } }));
  await $.session.start({ cwd: ".", surface: "terminal", isInteractive: true });
  const status = await $.command.run({
    command: "knowslink-session",
    args: "status",
  });
  expect(status.text).toMatch(/연결되어 있지 않습니다/);
  const off = await $.command.run({
    command: "knowslink-session",
    args: "off",
  });
  expect(off.text).toMatch(/연결되어 있지 않습니다/);
});

test("Bash and other tools pass through untouched", async ($, on) => {
  on("tool.call", { tool: "Bash" }, () => ({
    result: { stdout: "ok", stderr: "" },
  }));
  const ran = await $.tool.call({ tool: "Bash", command: "echo ok" });
  expect(ran.isError).toBeUndefined();
});
