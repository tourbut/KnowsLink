// Tester-owned CDP driver (SAR-SSE-INSTALL-001-OAUTH-RETURN-TESTER). Node 22 built-in WebSocket; prints no URL query or secret.
// argv: chrome hostRules startUrl userDataDir scenarioJSON outDir. stdout: optional "READY" gate line, then one result JSON line.
import { spawn } from "node:child_process";
import { writeFileSync } from "node:fs";
import { createInterface } from "node:readline";

const [chrome, rules, start, dir, scenarioJson, outDir] = process.argv.slice(2);
const sc = JSON.parse(scenarioJson);
const proc = spawn(chrome, ["--headless=new", "--no-sandbox", "--no-first-run", "--no-proxy-server", "--disable-gpu", "--disable-dev-shm-usage", "--disable-extensions", "--password-store=basic", "--ignore-certificate-errors", `--host-resolver-rules=${rules}`, `--user-data-dir=${dir}`, "--remote-debugging-port=0", "about:blank"], { stdio: ["ignore", "ignore", "pipe"] });
const redact = (s) => s.replace(/\?[^\s'"]*/g, "?<redacted>").replace(/\/remote\/consent\/[\w-]+/g, "/remote/consent/:id");
const sleep = (ms) => new Promise((r) => setTimeout(r, ms));
const result = { name: sc.name, stage: "launch", steps: [], csp: [], blocked: [], docs: [], arrivals: [], shots: [] };
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
}, 90000).unref();
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
  let page;
  for (let i = 0; !page && i < 50; i++) {
    page = (await (await fetch(`http://${new URL(endpoint).host}/json/list`)).json()).find((p) => p.type === "page");
    if (!page) await sleep(100);
  }
  const ws = new WebSocket(page.webSocketDebuggerUrl);
  await new Promise((r) => ws.addEventListener("open", r, { once: true }));
  let next = 0;
  const pending = new Map();
  const methods = new Map();
  let heldDone = false;
  let paused = 0;
  const send = (method, params = {}) =>
    new Promise((resolve) => {
      const id = ++next;
      pending.set(id, resolve);
      setTimeout(() => resolve({}), 5000);
      ws.send(JSON.stringify({ id, method, params }));
    });
  ws.addEventListener("message", async (e) => {
    const m = JSON.parse(e.data);
    if (m.id) pending.get(m.id)?.(m);
    const p = m.params;
    if (m.method === "Log.entryAdded" && /Content Security Policy/.test(p.entry.text)) result.csp.push(redact(p.entry.text));
    if (m.method === "Network.requestWillBeSent") methods.set(p.requestId, p.request.method);
    if (m.method === "Network.loadingFailed" && p.blockedReason) result.blocked.push(p.blockedReason);
    if (m.method === "Network.responseReceived" && p.type === "Document") {
      const h = Object.fromEntries(Object.entries(p.response.headers).map(([k, v]) => [k.toLowerCase(), v]));
      result.docs.push({ method: methods.get(p.requestId), status: p.response.status, url: redact(p.response.url), headers: { "cache-control": h["cache-control"], "content-security-policy": h["content-security-policy"], "referrer-policy": h["referrer-policy"], "x-content-type-options": h["x-content-type-options"], "content-type": h["content-type"] } });
    }
    if (m.method === "Fetch.requestPaused") {
      const u = new URL(p.request.url);
      paused++;
      const held = sc.hold && !heldDone;
      result.arrivals.push({ host: u.host, path: u.pathname, keys: [...u.searchParams.keys()].sort(), held });
      if (held) {
        heldDone = true;
        await send("Fetch.fulfillRequest", { requestId: p.requestId, responseCode: 204, responseHeaders: [] });
      } else {
        await send("Fetch.continueRequest", { requestId: p.requestId });
      }
    }
  });
  const evaluate = async (expression) => (await send("Runtime.evaluate", { expression, returnByValue: true })).result?.result?.value;
  const where = async () => (await evaluate("document.readyState === 'complete' && location.origin + location.pathname")) || "";
  const until = async (test, tries = 80) => {
    for (let i = 0; i < tries; i++) {
      const loc = await where();
      if (loc && test(loc)) return loc;
      await sleep(100);
    }
    return (await where()) || "unknown";
  };
  const shot = async (name) => {
    const r = await send("Page.captureScreenshot", { format: "png" });
    if (r.result?.data) {
      const file = `${outDir}/${sc.name}-${name}.png`;
      writeFileSync(file, Buffer.from(r.result.data, "base64"));
      result.shots.push(file);
    }
  };
  // Mobile uses real touch input; desktop uses a real mouse press at the element center.
  const click = async (selector) => {
    const c = await evaluate(`(() => { const b = document.querySelector(${JSON.stringify(selector)}).getBoundingClientRect(); return [b.x + b.width / 2, b.y + b.height / 2]; })()`);
    if (sc.mobile) {
      await send("Input.dispatchTouchEvent", { type: "touchStart", touchPoints: [{ x: c[0], y: c[1] }] });
      await send("Input.dispatchTouchEvent", { type: "touchEnd", touchPoints: [] });
    } else {
      await send("Input.dispatchMouseEvent", { type: "mouseMoved", x: c[0], y: c[1] });
      await send("Input.dispatchMouseEvent", { type: "mousePressed", x: c[0], y: c[1], button: "left", clickCount: 1 });
      await send("Input.dispatchMouseEvent", { type: "mouseReleased", x: c[0], y: c[1], button: "left", clickCount: 1 });
    }
  };
  await send("Log.enable");
  await send("Network.enable");
  await send("Fetch.enable", { patterns: [{ urlPattern: "https://callback.test/*", resourceType: "Document", requestStage: "Request" }] });
  const view = sc.mobile ? { width: 390, height: 844, deviceScaleFactor: 2, mobile: true } : { width: 1280, height: 800, deviceScaleFactor: 1, mobile: false };
  await send("Emulation.setDeviceMetricsOverride", view);
  if (sc.mobile) await send("Emulation.setTouchEmulationEnabled", { enabled: true });
  result.view = view;
  result.stage = "consent";
  await send("Page.navigate", { url: start });
  result.steps.push(redact(await until((l) => l.includes("/remote/consent/"))));
  const consentUrl = await evaluate("location.href");
  result.stage = "google";
  await click('form[action="/auth/google"] button');
  result.steps.push(redact(await until((l) => l.endsWith("/auth/google/callback"))));
  await click('a[href^="/remote/consent/"]');
  result.steps.push(redact(await until((l) => l.includes("/remote/consent/"))));
  await shot("consent");
  if (sc.gate) {
    console.log("READY");
    await new Promise((r) => createInterface({ input: process.stdin }).once("line", r));
  }
  result.stage = "decide";
  if (sc.clearCookies) await send("Network.clearBrowserCookies");
  if (sc.tamperCsrf) await evaluate(`document.querySelector('input[name=csrf]').value = 'tampered'`);
  await click(`button[value=${sc.decision || "approve"}]`);
  const arrived = async (afterClick = false) => {
    for (let i = 0; i < 80; i++) {
      const w = await where();
      if (w.startsWith("https://callback.test")) return true;
      if (/원격 연결을 계속할 수 없습니다/.test((await evaluate("document.querySelector('h1')?.innerText")) || "")) return false;
      if (sc.hold && !afterClick && paused) return null;
      await sleep(100);
    }
    return false;
  };
  const first = await arrived();
  if (sc.hold && first === null) {
    await sleep(300);
    result.stage = "return-page";
    result.returnPage = JSON.parse((await evaluate(`JSON.stringify((() => { const a = document.querySelector('a'); const r = a.getBoundingClientRect(); const m = document.querySelector('meta[http-equiv=refresh]'); return { h1: document.querySelector('h1').innerText, linkText: a.innerText, rect: [Math.round(r.x), Math.round(r.y), Math.round(r.width), Math.round(r.height)], innerWidth, innerHeight, scrollWidth: document.documentElement.scrollWidth, refreshZero: !!m && /^0;url=https:\\/\\/callback\\.test\\/oauth\\?/.test(m.content), origin: location.origin }; })())`)) || "null");
    await shot("return");
    result.stage = "link-click";
    await click("a");
    result.arrivedAfterClick = await arrived(true);
  }
  result.finalLocation = redact(await where());
  if (/원격 연결을 계속할 수 없습니다/.test((await evaluate("document.querySelector('h1')?.innerText")) || "")) {
    result.refusal = { text: await evaluate("document.body.innerText.slice(0, 400)"), links: await evaluate("[...document.querySelectorAll('a')].map((a) => a.getAttribute('href'))") };
    await shot("refusal");
    if (sc.retry) {
      result.stage = "retry";
      await send("Page.navigate", { url: consentUrl });
      await until((l) => l.includes("/remote/consent/"));
      await click("button[value=approve]");
      result.steps.push(redact(await until((l) => l.startsWith("https://callback.test"))));
      result.finalLocation = redact(await where());
    }
  }
  if (sc.revisit) {
    result.stage = "revisit";
    await send("Page.navigate", { url: consentUrl });
    await sleep(800);
    result.revisit = { text: await evaluate("document.body.innerText.slice(0, 400)") };
    await shot("revisit");
  }
  result.history = ((await send("Page.getNavigationHistory")).result?.entries || []).map((e) => redact(e.url));
  result.stage = "done";
  ws.close();
} catch (e) {
  result.error = String(e);
}
await finish(0);
