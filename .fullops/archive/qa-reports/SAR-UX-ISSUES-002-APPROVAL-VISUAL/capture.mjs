// Captures the approval fixture with allowed GETs only (plus same-origin POST /fixture/as/{a|c} for a session) at 1440/390/320.
// No approve/deny POST. Usage: node capture.mjs <chrome.exe> <origin> <manifest.json> <out-dir>
import { spawn } from 'node:child_process';
import { mkdtempSync, readFileSync, writeFileSync, mkdirSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join } from 'node:path';

const [chrome, origin, manifestPath, out] = process.argv.slice(2);
const gates = JSON.parse(readFileSync(manifestPath, 'utf8')).gates;
mkdirSync(out, { recursive: true });
const profile = mkdtempSync(join(tmpdir(), 'approval-visual-'));
const proc = spawn(chrome, ['--headless=new', '--remote-debugging-port=0', `--user-data-dir=${profile}`, '--no-first-run', '--disable-gpu', 'about:blank']);
process.on('exit', () => proc.kill());
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
const evaluate = async (expression) => (await send('Runtime.evaluate', { expression, awaitPromise: true, returnByValue: true }, sessionId)).result.result.value;
// Paced so the shared fixture's rate limit is not tripped for the tester.
const go = async (path) => { await send('Page.navigate', { url: origin + path }, sessionId); await new Promise((r) => setTimeout(r, 2500)); };
// The session cookie is __Host-/Secure, which headless Chrome drops on http loopback, so send it as a header instead.
const login = async (who) => {
  const r = await fetch(`${origin}/fixture/as/${who}`, { method: 'POST', headers: { Origin: origin } });
  const cookie = (r.headers.get('set-cookie') || '').split(';')[0];
  if (r.status !== 204 || !cookie) throw new Error(`login ${who} failed: ${r.status}`);
  await send('Network.enable', {}, sessionId);
  await send('Network.setExtraHTTPHeaders', { headers: { Cookie: cookie } }, sessionId);
  return r.status;
};

// Measures page overflow, small targets, non-left text, shapes to avoid, colors in use, headings and main text.
const probe = `(() => {
  const w = innerWidth, vis = (e) => { const r = e.getBoundingClientRect(); return r.width > 0 && r.height > 0; };
  const over = [...document.querySelectorAll('body *')].filter(e => vis(e) && e.getBoundingClientRect().right > w + 0.5).map(e => e.tagName + (e.className ? '.' + e.className : '')).slice(0, 5);
  const small = [...document.querySelectorAll('a[href],button,summary,input:not([type=hidden]),select,textarea')].filter(vis)
    .filter(e => { const r = e.getBoundingClientRect(); return r.height < 44 || (r.width < 44 && e.tagName !== 'A'); })
    .map(e => ({ tag: e.tagName, text: (e.innerText || e.value || '').trim().slice(0, 40), h: Math.round(e.getBoundingClientRect().height), w: Math.round(e.getBoundingClientRect().width) }));
  const notLeft = [...document.querySelectorAll('main *')].filter(e => vis(e) && e.innerText && ['center', 'right', 'end'].includes(getComputedStyle(e).textAlign)).map(e => e.tagName + ':' + e.innerText.trim().slice(0, 30)).slice(0, 5);
  const shapes = [...document.querySelectorAll('body *')].filter(vis).flatMap(e => { const s = getComputedStyle(e), r = e.getBoundingClientRect(), out = [];
    if (s.backgroundImage !== 'none') out.push('gradient/image:' + e.tagName);
    if (s.boxShadow !== 'none') out.push('shadow:' + e.tagName);
    if (parseFloat(s.borderTopLeftRadius) > 4) out.push('radius' + s.borderTopLeftRadius + ':' + e.tagName + ':' + (e.innerText || '').trim().slice(0, 20));
    return out; }).slice(0, 10);
  const colors = new Set(); for (const e of document.querySelectorAll('body, body *')) { if (!vis(e) && e !== document.body) continue; const s = getComputedStyle(e); colors.add(s.color); if (s.backgroundColor !== 'rgba(0, 0, 0, 0)') colors.add('bg ' + s.backgroundColor); }
  const headings = [...document.querySelectorAll('h1,h2,h3,h4')].filter(vis).map(h => h.tagName + ' ' + h.innerText.trim());
  const details = [...document.querySelectorAll('details')].map(d => (d.querySelector('summary')?.innerText.trim() || '') + (d.open ? ' [open]' : ' [closed]'));
  const current = [...document.querySelectorAll('[aria-current]')].map(e => e.getAttribute('aria-current') + ':' + e.innerText.trim());
  const roles = [...document.querySelectorAll('[role=alert],[role=status]')].map(e => e.getAttribute('role') + ':' + e.innerText.trim().slice(0, 80));
  const anchors = [...document.querySelectorAll('main a[href]')].map(a => a.innerText.trim().slice(0, 40) + ' -> ' + a.getAttribute('href'));
  const buttons = [...document.querySelectorAll('main button')].map(b => b.innerText.trim());
  return { title: document.title, scrollWidth: document.documentElement.scrollWidth, innerWidth: w, height: document.documentElement.scrollHeight, over, small, notLeft, shapes, colors: [...colors], headings, details, current, roles, anchors, buttons, main: (document.querySelector('main') || document.body).innerText };
})()`;

const pages = [
  ['a', 'home-a', '/home'],
  ['a', 'legacy-a', '/home/legacy'],
  ['a', 'friends-a', '/home/friends'],
  ...Object.entries(gates).map(([name, id]) => ['a', `gate-${name}`, `/home/gates/${id}`]),
  ['c', 'legacy-c-empty', '/home/legacy'],
];
const results = [];
let who = '';
for (const [account, name, path] of pages) {
  if (who !== account) { results.push({ login: account, status: await login(account) }); who = account; }
  for (const width of [1440, 390, 320]) {
    await send('Emulation.setDeviceMetricsOverride', { width, height: 900, deviceScaleFactor: 1, mobile: false }, sessionId);
    await go(path);
    const v = await evaluate(probe);
    // Stop instead of hammering the shared fixture when the session is missing or rate limited.
    if (/서비스 안내와 시작하기|요청 제한/.test(v.title)) throw new Error(`${name} ${width}: unexpected page "${v.title}"`);
    await send('Emulation.setDeviceMetricsOverride', { width, height: Math.min(v.height, 8000), deviceScaleFactor: 1, mobile: false }, sessionId);
    const shot = await send('Page.captureScreenshot', { format: 'png' }, sessionId);
    writeFileSync(join(out, `${name}-${width}.png`), Buffer.from(shot.result.data, 'base64'));
    // The approval list sits deep in a long legacy page, so keep a readable crop of it too.
    const clip = await evaluate(`(() => { const r = document.querySelector('section#gates')?.getBoundingClientRect(); return r && { x: 0, y: r.top + scrollY - 8, width: innerWidth, height: r.height + 16, scale: 1 }; })()`);
    if (clip) writeFileSync(join(out, `${name}-${width}-gates.png`), Buffer.from((await send('Page.captureScreenshot', { format: 'png', clip }, sessionId)).result.data, 'base64'));
    results.push({ name, path, width, overflow: v.scrollWidth > v.innerWidth, ...v, main: width === 1440 ? v.main : undefined });
  }
}
writeFileSync(join(out, 'capture.json'), JSON.stringify(results, null, 2) + '\n');
console.log(results.filter((r) => r.overflow).length ? 'OVERFLOW' : 'NO_OVERFLOW', results.filter((r) => r.name).length);
browser.close();
proc.kill();
