// SAR-SCHEDULE-001-VISUAL-FIX copy: also records text-only summaries whose text does not start beside the marker (V2).
// Renders synthetic fixture HTML at 1440/390/320 in headless Chrome (non-mobile viewport) and records page overflow.
// Usage: node probe.mjs <chrome.exe> <fixture-html-dir> <out-dir> <name>...
import { spawn } from 'node:child_process';
import { mkdtempSync, writeFileSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import { pathToFileURL } from 'node:url';

const [chrome, dir, out, ...names] = process.argv.slice(2);
const profile = mkdtempSync(join(tmpdir(), 'approval-probe-'));
const proc = spawn(chrome, ['--headless=new', '--remote-debugging-port=0', `--user-data-dir=${profile}`, '--no-first-run', '--disable-gpu', 'about:blank']);
const wsURL = await new Promise((ok, fail) => {
  proc.stderr.on('data', (d) => { const m = /DevTools listening on (ws:\S+)/.exec(String(d)); if (m) ok(m[1]); });
  setTimeout(() => fail(new Error('chrome did not start')), 15000);
});
const browser = new WebSocket(wsURL);
await new Promise((ok) => browser.addEventListener('open', ok));
let seq = 0;
const waiting = new Map();
browser.addEventListener('message', (ev) => { const m = JSON.parse(ev.data); if (waiting.has(m.id)) { waiting.get(m.id)(m); waiting.delete(m.id); } });
const send = (method, params = {}, sessionId) => new Promise((ok) => { const id = ++seq; waiting.set(id, ok); browser.send(JSON.stringify({ id, method, params, sessionId })); });
const { result: { targetId } } = await send('Target.createTarget', { url: 'about:blank' });
const { result: { sessionId } } = await send('Target.attachToTarget', { targetId, flatten: true });
const results = [];
for (const name of names) {
  for (const width of [1440, 390, 320]) {
    await send('Emulation.setDeviceMetricsOverride', { width, height: 900, deviceScaleFactor: 1, mobile: false }, sessionId);
    await send('Page.navigate', { url: pathToFileURL(join(dir, `${name}.html`)).href }, sessionId);
    await new Promise((r) => setTimeout(r, 400));
    const probe = `(() => { const w = innerWidth; const over = [...document.querySelectorAll('body *')].filter(e => e.getBoundingClientRect().right > w + 0.5).map(e => e.tagName + (e.className ? '.' + e.className : '')).slice(0, 5); const split = [...document.querySelectorAll('summary')].filter(e => !e.children.length && e.firstChild).filter(e => { const r = document.createRange(); r.selectNodeContents(e); return r.getClientRects()[0].left - e.getBoundingClientRect().left < 20; }).map(e => e.textContent); return JSON.stringify({ scrollWidth: document.documentElement.scrollWidth, innerWidth: w, height: document.documentElement.scrollHeight, over, split }); })()`;
    const { result: { result: { value } } } = await send('Runtime.evaluate', { expression: probe }, sessionId);
    const v = JSON.parse(value);
    await send('Emulation.setDeviceMetricsOverride', { width, height: Math.min(v.height, 8000), deviceScaleFactor: 1, mobile: false }, sessionId);
    const shot = await send('Page.captureScreenshot', { format: 'png' }, sessionId);
    writeFileSync(join(out, `${name}-${width}.png`), Buffer.from(shot.result.data, 'base64'));
    results.push({ name, width, overflow: v.scrollWidth > v.innerWidth, ...v });
  }
}
writeFileSync(join(out, 'probe.json'), JSON.stringify(results, null, 2) + '\n');
console.log(results.filter((r) => r.overflow).length ? 'OVERFLOW' : 'NO_OVERFLOW', results.filter((r) => r.split.length).length ? 'SUMMARY_SPLIT' : 'SUMMARY_INLINE', results.length);
browser.close();
proc.kill();
