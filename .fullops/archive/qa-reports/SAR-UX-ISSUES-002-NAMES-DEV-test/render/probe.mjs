// Minimal CDP probe: page overflow, clipped text boxes and 44px targets at fixed widths; writes PNGs.
import { spawn } from "node:child_process";
import { writeFileSync, mkdtempSync } from "node:fs";
import { tmpdir } from "node:os";
import { join } from "node:path";
const [dir, out, ...files] = process.argv.slice(2);
const chrome = spawn("C:/Program Files/Google/Chrome/Application/chrome.exe", ["--headless=new", "--remote-debugging-port=9333", "--user-data-dir=" + mkdtempSync(join(tmpdir(), "cdp-")), "about:blank"], { stdio: "ignore" });
await new Promise(r => setTimeout(r, 2500));
const targets = await (await fetch("http://127.0.0.1:9333/json")).json();
const ws = new WebSocket(targets.find(t => t.type === "page").webSocketDebuggerUrl);
await new Promise(r => ws.onopen = r);
let id = 0; const wait = new Map();
ws.onmessage = e => { const m = JSON.parse(e.data); if (wait.has(m.id)) { wait.get(m.id)(m); wait.delete(m.id); } };
const send = (method, params = {}) => new Promise(r => { const n = ++id; wait.set(n, r); ws.send(JSON.stringify({ id: n, method, params })); });
await send("Page.enable");
const probe = `(() => { const W = innerWidth, r = { overflow: document.documentElement.scrollWidth > W, clipped: [], small: [] };
 for (const el of document.querySelectorAll("label,span,strong,p,h1,h2,h3,code,dd,legend,textarea,input,select,a,button,summary")) {
  const s = getComputedStyle(el); if (s.display === "none" || el.closest("details:not([open])") && el.tagName !== "SUMMARY") continue;
  const b = el.getBoundingClientRect(); if (!b.width) continue;
  if (b.right > W + 0.5) r.clipped.push(el.tagName + ":" + (el.id || el.textContent.trim().slice(0, 30)));
  if (/^(TEXTAREA|INPUT)$/.test(el.tagName) && el.scrollWidth > el.clientWidth + 1) r.clipped.push("scroll:" + el.tagName + "#" + el.id);
  if (/^(A|BUTTON|SUMMARY|INPUT|SELECT|TEXTAREA)$/.test(el.tagName) && el.type !== "hidden" && b.height < 44 && el.type !== "radio" && el.type !== "checkbox") r.small.push(el.tagName + ":" + el.textContent.trim().slice(0, 20) + ":" + Math.round(b.height));
 } return JSON.stringify(r); })()`;
const results = {};
for (const f of files) for (const w of [1440, 390, 320]) {
  await send("Emulation.setDeviceMetricsOverride", { width: w, height: 900, deviceScaleFactor: 1, mobile: false });
  const loaded = new Promise(r => { const h = e => { if (JSON.parse(e.data).method === "Page.loadEventFired") { ws.removeEventListener("message", h); r(); } }; ws.addEventListener("message", h); });
  await send("Page.navigate", { url: "file:///" + dir + "/" + f + ".html" }); await loaded;
  const v = await send("Runtime.evaluate", { expression: probe, returnByValue: true });
  results[f + "@" + w] = JSON.parse(v.result.result.value);
  const shot = await send("Page.captureScreenshot", { format: "png", captureBeyondViewport: true });
  writeFileSync(join(out, f + "-" + w + ".png"), Buffer.from(shot.result.data, "base64"));
}
writeFileSync(join(out, "probe.json"), JSON.stringify(results, null, 1));
console.log(JSON.stringify(results));
ws.close(); chrome.kill();
