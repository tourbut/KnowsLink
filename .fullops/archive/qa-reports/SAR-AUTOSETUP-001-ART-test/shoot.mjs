// SAR-AUTOSETUP-001-ART: drives headless Chromium over CDP against the isolated loopback relay (synthetic email member,
// loopback SMTP sink for the login code). Measures each screen x width x language and writes cropped PNGs of #use.
// usage: node shoot.mjs <relay origin> <mail dir> <out dir>
import { spawn } from "node:child_process";
import { mkdirSync, readFileSync, writeFileSync, rmSync } from "node:fs";
import { tmpdir } from "node:os";
import { join } from "node:path";

const [ORIGIN, MAIL, OUT] = process.argv.slice(2);
const CHROME = process.env.LOCALAPPDATA + "\\ms-playwright\\chromium-1248\\chrome-win64\\chrome.exe";
const profile = join(tmpdir(), "art-chrome-" + process.pid);
mkdirSync(OUT, { recursive: true });
const chrome = spawn(CHROME, ["--headless=new", "--remote-debugging-port=9333", "--user-data-dir=" + profile, "--no-first-run", "--hide-scrollbars=false", "about:blank"], { stdio: "ignore" });
const sleep = (ms) => new Promise((r) => setTimeout(r, ms));
let ws, seq = 0;
const pending = new Map(), events = [];
async function connect() {
  for (let i = 0; i < 50; i++) {
    try {
      const list = await (await fetch("http://127.0.0.1:9333/json")).json();
      const page = list.find((t) => t.type === "page");
      if (page) { ws = new WebSocket(page.webSocketDebuggerUrl); break; }
    } catch { /* chrome not ready */ }
    await sleep(200);
  }
  await new Promise((r) => (ws.onopen = r));
  ws.onmessage = (m) => {
    const d = JSON.parse(m.data);
    if (d.id && pending.has(d.id)) { pending.get(d.id)(d); pending.delete(d.id); } else events.push(d);
  };
}
const send = (method, params = {}) => new Promise((res, rej) => {
  const id = ++seq;
  pending.set(id, (d) => (d.error ? rej(new Error(method + ": " + d.error.message)) : res(d.result)));
  ws.send(JSON.stringify({ id, method, params }));
});
const ev = async (expr) => (await send("Runtime.evaluate", { expression: expr, returnByValue: true, awaitPromise: true })).result.value;
async function go(url) {
  events.length = 0;
  await send("Page.navigate", { url });
  for (let i = 0; i < 100; i++) { if (events.some((e) => e.method === "Page.loadEventFired")) break; await sleep(100); }
  await sleep(250);
}

const MEASURE = String.raw`(() => {
  const de = document.documentElement, iw = innerWidth;
  const vis = (e) => { const r = e.getBoundingClientRect(), s = getComputedStyle(e); return r.width > 0 && r.height > 0 && s.visibility !== "hidden" && s.display !== "none"; };
  const over = [...document.querySelectorAll("body *")].filter(vis).filter((e) => { const r = e.getBoundingClientRect(); return r.right > iw + 1 || r.left < -1; }).map((e) => e.tagName.toLowerCase() + "." + e.className);
  const rgb = (c) => { const m = c.match(/[\d.]+/g).map(Number); return { r: m[0], g: m[1], b: m[2], a: m[3] === undefined ? 1 : m[3] }; };
  const bgOf = (e) => { for (; e; e = e.parentElement) { const c = rgb(getComputedStyle(e).backgroundColor); if (c.a > 0) return c; } return { r: 255, g: 255, b: 255, a: 1 }; };
  const lum = (c) => { const f = (v) => { v /= 255; return v <= 0.03928 ? v / 12.92 : ((v + 0.055) / 1.055) ** 2.4; }; return 0.2126 * f(c.r) + 0.7152 * f(c.g) + 0.0722 * f(c.b); };
  const contrast = (e) => { const a = lum(rgb(getComputedStyle(e).color)), b = lum(bgOf(e)); return +((Math.max(a, b) + 0.05) / (Math.min(a, b) + 0.05)).toFixed(2); };
  const use = document.querySelector("#use");
  const guide = document.querySelector("#use-guide");
  const ta = document.querySelector("#setup-command");
  const out = { iw, sw: de.scrollWidth, docH: de.scrollHeight, lang: de.lang, overflow: over };
  const r = use.getBoundingClientRect();
  out.use = { x: r.left + scrollX, y: r.top + scrollY, w: r.width, h: r.height };
  out.choices = [...document.querySelectorAll("#use .choice")].map((l) => l.innerText.replace(/\s+/g, " ").trim());
  out.guideH3 = guide && guide.querySelector("h3") ? guide.querySelector("h3").innerText : null;
  out.rows = guide ? [...guide.querySelectorAll("dl.kv > dt")].map((dt) => [dt.innerText.trim(), dt.nextElementSibling.innerText.replace(/\s+/g, " ").trim()]) : [];
  out.note = guide && guide.querySelector("p.note") ? guide.querySelector("p.note").innerText : null;
  out.details = guide && guide.querySelector("details") ? { summary: guide.querySelector("details summary").innerText, href: (guide.querySelector("details a") || {}).href, linkText: (guide.querySelector("details a") || {}).innerText } : null;
  if (ta) {
    const cs = getComputedStyle(ta), lbl = document.querySelector("label[for=setup-command]"), help = document.querySelector("#setup-command-help");
    ta.focus(); ta.select();
    const lh = parseFloat(cs.lineHeight) || parseFloat(cs.fontSize) * 1.2;
    out.cmd = { value: ta.value, hasNewline: /\n/.test(ta.value), readOnly: ta.readOnly, spellcheck: ta.spellcheck, rows: ta.rows, clientW: ta.clientWidth, scrollW: ta.scrollWidth, clientH: ta.clientHeight, scrollH: ta.scrollHeight, wFrac: +(ta.getBoundingClientRect().width / use.getBoundingClientRect().width).toFixed(2),
      selectedAll: ta.value.slice(ta.selectionStart, ta.selectionEnd) === ta.value, font: cs.fontFamily, fontPx: cs.fontSize, textLines: Math.round((ta.scrollHeight - parseFloat(cs.paddingTop) - parseFloat(cs.paddingBottom)) / lh),
      label: lbl && lbl.innerText, help: help && help.innerText, contrastText: contrast(ta), contrastHelp: help ? contrast(help) : null, contrastLabel: lbl ? contrast(lbl) : null, describedby: ta.getAttribute("aria-describedby") };
  }
  const muted = guide ? [...guide.querySelectorAll(".muted")].filter(vis) : [];
  out.mutedMinContrast = muted.length ? Math.min(...muted.map(contrast)) : null;
  out.hasCommandUI = !!guide && !!guide.querySelector("textarea");
  out.buttonsInGuide = guide ? [...guide.querySelectorAll("button")].map((b) => b.innerText) : [];
  out.hangulInEn = document.documentElement.lang === "en" && use ? ((use.innerText.match(/[가-힣]+/g)) || []) : null;
  return JSON.stringify(out);
})()`;

const WIDTHS = { 1440: [1440, 900], 390: [390, 844], 320: [320, 640] };
const results = [];
async function shot(name, w, clip) {
  const { data } = await send("Page.captureScreenshot", { format: "png", captureBeyondViewport: true, clip: { x: 0, y: clip.y, width: w, height: Math.min(clip.h, 6000), scale: 1 } });
  writeFileSync(join(OUT, name + ".png"), Buffer.from(data, "base64"));
}
try {
  await connect();
  await send("Page.enable");
  // real email login against the loopback relay; the code comes from the sink, never from a real mailbox
  await send("Emulation.setDeviceMetricsOverride", { width: 1440, height: 900, deviceScaleFactor: 1, mobile: false });
  await go(ORIGIN + "/");
  rmSync(join(MAIL, "last.txt"), { force: true });
  await ev(`(()=>{const f=document.querySelector('input[name=email]');f.value='art-${Date.now()}@autosetup-art.invalid';f.form.submit();})()`);
  await sleep(1500);
  let code = null;
  for (let i = 0; i < 30 && !code; i++) { try { code = readFileSync(join(MAIL, "last.txt"), "utf8").match(/\b(\d{6,8})\b/)?.[1]; } catch { /* wait */ } await sleep(300); }
  if (!code) throw new Error("no login code in sink");
  await ev(`(()=>{const f=document.querySelector('input[name=code]');f.value='${code}';f.form.submit();})()`);
  await sleep(1500);
  if ((await ev("location.pathname")) === "/home/name") {
    await ev(`(()=>{const f=document.querySelector('input[name=name]');f.value='합성 회원 ART';f.form.submit();})()`);
    await sleep(1500);
  }
  console.log("after login:", await ev("location.pathname"));

  const SCREENS = [["list", ""], ["claude", "use=claude&start=first"], ["codex", "use=codex&start=first"]];
  for (const lang of ["ko", "en"]) {
    for (const [wname, [w, h]] of Object.entries(WIDTHS)) {
      await send("Emulation.setDeviceMetricsOverride", { width: w, height: h, deviceScaleFactor: 1, mobile: false });
      for (const [screen, q] of SCREENS) {
        await go(`${ORIGIN}/home/agents?${q ? q + "&" : ""}lang=${lang}#use-guide`);
        const m = JSON.parse(await ev(MEASURE));
        m.key = `${screen}-${wname}-${lang}`;
        results.push(m);
        await shot(m.key, w, m.use);
      }
    }
  }
  // text-only extras at 1440: continue variants and the screens that must carry no command (AC-10)
  await send("Emulation.setDeviceMetricsOverride", { width: 1440, height: 900, deviceScaleFactor: 1, mobile: false });
  for (const lang of ["ko", "en"]) {
    for (const q of ["use=claude&start=continue", "use=codex&start=continue", "use=web&start=first", "use=grok&start=first", "use=dots&start=first", "use=other&start=first"]) {
      await go(`${ORIGIN}/home/agents?${q}&lang=${lang}#use-guide`);
      const m = JSON.parse(await ev(MEASURE));
      m.key = `extra-${q.replace(/[=&]/g, "_")}-${lang}`;
      results.push(m);
    }
  }
  writeFileSync(join(OUT, "measure.json"), JSON.stringify(results, null, 1));
  console.log("measured", results.length);
} catch (e) {
  console.error("FAILED:", e.stack);
  process.exitCode = 1;
} finally {
  chrome.kill();
  await sleep(500);
  rmSync(profile, { recursive: true, force: true });
  process.exit();
}
