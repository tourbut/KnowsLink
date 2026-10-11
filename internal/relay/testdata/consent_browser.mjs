// Drives a real Chromium through synthetic remote consent with Node's built-in WebSocket; prints no URL query.
import { spawn } from "node:child_process";

const [chrome, rules, start, dir] = process.argv.slice(2);
const proc = spawn(chrome, ["--headless=new", "--no-sandbox", "--no-first-run", "--disable-gpu", "--ignore-certificate-errors", `--host-resolver-rules=${rules}`, `--user-data-dir=${dir}`, "--remote-debugging-port=0", "about:blank"], { stdio: ["ignore", "ignore", "pipe"] });
const redact = (s) => s.replace(/\?[^\s'"]*/g, "?<redacted>");
const sleep = (ms) => new Promise((r) => setTimeout(r, ms));
const result = { steps: [], csp: [] };
setTimeout(() => {
  proc.kill();
  console.log(JSON.stringify({ ...result, error: "timeout" }));
  process.exit(1);
}, 30000).unref();
try {
  const endpoint = await new Promise((resolve, reject) => {
    let err = "";
    proc.stderr.on("data", (d) => {
      err += d;
      const m = err.match(/DevTools listening on (ws:\/\/\S+)/);
      if (m) resolve(m[1]);
    });
    proc.on("exit", () => reject(new Error("browser exited")));
  });
  const pages = await (await fetch(`http://${new URL(endpoint).host}/json/list`)).json();
  const ws = new WebSocket(pages.find((p) => p.type === "page").webSocketDebuggerUrl);
  await new Promise((r) => ws.addEventListener("open", r, { once: true }));
  let next = 0;
  const pending = new Map();
  ws.addEventListener("message", (e) => {
    const m = JSON.parse(e.data);
    if (m.id) pending.get(m.id)?.(m);
    if (m.method === "Log.entryAdded" && /Content Security Policy/.test(m.params.entry.text)) result.csp.push(redact(m.params.entry.text));
  });
  const send = (method, params = {}) =>
    new Promise((resolve) => {
      const id = ++next;
      pending.set(id, resolve);
      setTimeout(() => resolve({}), 2000);
      ws.send(JSON.stringify({ id, method, params }));
    });
  const evaluate = async (expression) => (await send("Runtime.evaluate", { expression, returnByValue: true })).result?.result?.value;
  const until = async (test) => {
    for (let i = 0; i < 50; i++) {
      const loc = await evaluate("document.readyState === 'complete' && location.origin + location.pathname");
      if (loc && test(loc)) return loc;
      await sleep(100);
    }
    return (await evaluate("location.origin + location.pathname")) || "unknown";
  };
  await send("Log.enable");
  await send("Page.navigate", { url: start });
  result.steps.push(await until((l) => l.includes("/remote/consent/")));
  await evaluate(`document.querySelector('form[action="/auth/google"] button').click()`);
  result.steps.push(await until((l) => l.endsWith("/auth/google/callback")));
  await evaluate(`document.querySelector('a[href^="/remote/consent/"]').click()`);
  result.steps.push(await until((l) => l.includes("/remote/consent/")));
  await evaluate(`document.querySelector('button[value=approve]').click()`);
  result.steps.push(await until((l) => l.startsWith("https://callback.test")));
  await sleep(300);
  ws.close();
} catch (e) {
  result.error = String(e);
} finally {
  proc.kill();
  await sleep(300);
}
console.log(JSON.stringify(result));
