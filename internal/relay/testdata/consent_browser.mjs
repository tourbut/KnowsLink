// Drives a real Chromium through synthetic remote consent with Node's built-in WebSocket; prints no URL query.
import { spawn } from "node:child_process";

const [chrome, rules, start, dir] = process.argv.slice(2);
const proc = spawn(chrome, ["--headless=new", "--no-sandbox", "--no-first-run", "--no-proxy-server", "--disable-gpu", "--disable-dev-shm-usage", "--disable-extensions", "--password-store=basic", "--ignore-certificate-errors", `--host-resolver-rules=${rules}`, `--user-data-dir=${dir}`, "--remote-debugging-port=0", "about:blank"], { stdio: ["ignore", "ignore", "pipe"] });
const redact = (s) => s.replace(/\?[^\s'"]*/g, "?<redacted>");
const sleep = (ms) => new Promise((r) => setTimeout(r, ms));
const result = { stage: "launch", steps: [], csp: [] };
const exited = new Promise((r) => proc.on("exit", r));
const finish = async (code) => {
  proc.kill();
  await Promise.race([exited, sleep(5000)]);
  console.log(JSON.stringify(result));
  process.exit(code);
};
setTimeout(() => {
  result.error = "timeout";
  finish(1);
}, 45000).unref();
try {
  const endpoint = await new Promise((resolve, reject) => {
    let err = "";
    proc.stderr.on("data", (d) => {
      err += d;
      const m = err.match(/DevTools listening on (ws:\/\/\S+)/);
      if (m) resolve(m[1]);
    });
    exited.then(() => reject(new Error("browser exited")));
  });
  result.stage = "target";
  let page;
  for (let i = 0; !page && i < 50; i++) {
    page = (await (await fetch(`http://${new URL(endpoint).host}/json/list`)).json()).find((p) => p.type === "page");
    if (!page) await sleep(100);
  }
  const ws = new WebSocket(page.webSocketDebuggerUrl);
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
  result.stage = "consent";
  await send("Log.enable");
  await send("Page.navigate", { url: start });
  result.steps.push(await until((l) => l.includes("/remote/consent/")));
  result.stage = "google";
  await evaluate(`document.querySelector('form[action="/auth/google"] button').click()`);
  result.steps.push(await until((l) => l.endsWith("/auth/google/callback")));
  result.stage = "continue";
  await evaluate(`document.querySelector('a[href^="/remote/consent/"]').click()`);
  result.steps.push(await until((l) => l.includes("/remote/consent/")));
  result.stage = "approve";
  await evaluate(`document.querySelector('input[name=name]').value = '합성 원격 도우미'`);
  await evaluate(`document.querySelector('button[value=approve]').click()`);
  result.steps.push(await until((l) => l.startsWith("https://callback.test")));
  result.stage = "done";
  await sleep(300);
  ws.close();
} catch (e) {
  result.error = String(e);
}
await finish(0);
