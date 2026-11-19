// Captures fixture pages with allowed GETs (plus the fixture login POST) at 1440/390/320 and records layout probes.
// Usage: KL_A=<token> KL_B=<token> KL_C=<token> node capture.mjs <chrome.exe> <origin> <out-dir> <profile-dir> [acct:/path ...]
import { spawn } from 'node:child_process';
import { mkdirSync, writeFileSync } from 'node:fs';
import { join } from 'node:path';

const [chrome, origin, out, profile] = process.argv.slice(2);
mkdirSync(profile, { recursive: true });
// PAGES (JSON {acct: [path]}) and PROBE (file name) replace the default set for a supplementary run.
const pages = process.env.PAGES ? JSON.parse(process.env.PAGES) : {
  a: ['/home/agents', '/home/agents?use=web&start=first', '/home/agents?use=grok&start=first', '/home/agents?use=dots&start=first',
    '/home/agents?use=claude&start=first', '/home/agents?use=other&start=first', '/home/agents?use=grok&start=continue',
    '/home/friends', '/home/legacy', '/home/account'],
  b: ['/home/agents', '/home/agents?use=dots&start=continue'],
  c: ['/home/agents', '/home/friends'],
};
const extraConnections = process.argv.slice(6); // optional "acct:/path" entries
for (const e of extraConnections) { const [acct, path] = e.split(/:(.+)/); pages[acct].push(path); }
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
const evaluate = async (expression, sessionId) => (await send('Runtime.evaluate', { expression, awaitPromise: true, returnByValue: true }, sessionId)).result.result.value;
const sleep = (ms) => new Promise((r) => setTimeout(r, ms));
const results = [];
for (const [acct, list] of Object.entries(pages)) {
  const { result: { browserContextId } } = await send('Target.createBrowserContext');
  const { result: { targetId } } = await send('Target.createTarget', { url: 'about:blank', browserContextId });
  const { result: { sessionId } } = await send('Target.attachToTarget', { targetId, flatten: true });
  // The fixture session cookie is __Host- and Secure, which Chrome drops from an http fetch, so the login token
  // from POST /fixture/as/{acct} (taken outside the browser) is set directly.
  const login = (await send('Network.setCookie', { name: '__Host-kl_session', value: process.env['KL_' + acct.toUpperCase()], url: origin + '/', path: '/', secure: true, httpOnly: true, sameSite: 'Strict' }, sessionId)).result?.success;
  for (const path of list) {
    for (const width of [1440, 390, 320]) {
      await send('Emulation.setDeviceMetricsOverride', { width, height: 900, deviceScaleFactor: 1, mobile: false }, sessionId);
      const probe = `(() => {
        const w = innerWidth;
        const over = [...document.querySelectorAll('body *')].filter(e => e.getBoundingClientRect().right > w + 0.5).map(e => e.tagName + (e.className ? '.' + e.className : '')).slice(0, 5);
        const small = [...document.querySelectorAll('a[href], button, input:not([type=hidden]), select, summary, textarea')].filter(e => { const r = e.getBoundingClientRect(); return r.width > 0 && (r.height < 44 && !(e.tagName === 'A' && getComputedStyle(e).display === 'inline')); }).map(e => e.tagName + ':' + (e.textContent || e.value || e.name || '').trim().slice(0, 30) + ':' + Math.round(e.getBoundingClientRect().height)).slice(0, 12);
        const centered = [...document.querySelectorAll('main *')].filter(e => getComputedStyle(e).textAlign === 'center' && e.textContent.trim()).map(e => e.tagName).slice(0, 5);
        const text = document.body.innerText;
        return { title: document.title, scrollWidth: document.documentElement.scrollWidth, innerWidth: w, height: document.documentElement.scrollHeight, over, small, centered,
          online: /온라인/.test(text), onlineBadge: [...document.querySelectorAll('.badge, [class*=badge]')].some(e => /온라인/.test(e.textContent)),
          hidden: ['OpenClaw', 'Meta Muse', 'Claude Desktop', 'Cursor', 'Grok CLI'].filter(n => text.includes(n)),
          currentStep: [...document.querySelectorAll('[aria-current]')].map(e => e.getAttribute('aria-current') + ':' + e.textContent.trim().slice(0, 30)),
          text };
      })()`;
      // The shared fixture rate-limits each member (40/min) and tester uses the same accounts, so pace and back off.
      let v;
      for (let attempt = 0; ; attempt++) {
        await send('Page.navigate', { url: origin + path }, sessionId);
        await sleep(4000);
        v = await evaluate(probe, sessionId);
        if (!v.title.startsWith('요청 제한')) break;
        if (attempt === 6) throw new Error(`rate limited at ${acct} ${path} ${width}`);
        await sleep(30000);
      }
      await send('Emulation.setDeviceMetricsOverride', { width, height: Math.min(v.height, 8000), deviceScaleFactor: 1, mobile: false }, sessionId);
      const shot = await send('Page.captureScreenshot', { format: 'png' }, sessionId);
      const name = `${acct}${path.replace(/[^a-z0-9]+/gi, '-').replace(/-$/, '')}-${width}`;
      writeFileSync(join(out, `${name}.png`), Buffer.from(shot.result.data, 'base64'));
      if (width === 1440) writeFileSync(join(out, `${name}.txt`), v.text);
      delete v.text;
      results.push({ acct, path, width, login, overflow: v.scrollWidth > v.innerWidth, png: `${name}.png`, ...v });
    }
  }
}
writeFileSync(join(out, process.env.PROBE || 'probe.json'), JSON.stringify(results, null, 2) + '\n');
console.log(results.filter((r) => r.overflow).length ? 'OVERFLOW' : 'NO_OVERFLOW', results.length);
proc.kill();
process.exit(0);
