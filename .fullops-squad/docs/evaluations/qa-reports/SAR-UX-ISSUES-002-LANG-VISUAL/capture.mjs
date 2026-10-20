// Captures the language fixture with allowed GETs (plus the fixture login POST, taken outside the browser) at 1440/390/320.
// Usage: KL_EN=<token> KL_KO=<token> KL_NONE=<token> node capture.mjs <chrome.exe> <origin> <out-dir> <profile-dir> <manifest.json>
import { spawn } from 'node:child_process';
import { mkdirSync, readFileSync, writeFileSync } from 'node:fs';
import { join } from 'node:path';

const [chrome, origin, out, profile, manifestPath] = process.argv.slice(2);
mkdirSync(profile, { recursive: true });
mkdirSync(out, { recursive: true });
const { ids, conversation } = JSON.parse(readFileSync(manifestPath, 'utf8'));
const conn = (s) => `/home/connections/${ids['connection_' + s]}`;
// Each context: account session (or anonymous) and the __Host-kl_lang cookie a GET /lang/{code} would set.
const contexts = process.env.CONTEXTS ? JSON.parse(process.env.CONTEXTS) : [
  { name: 'en', acct: 'en', paths: ['/home/account', '/home/account?n=language', '/home', '/home/name', '/home/friends', '/home/agents', '/home/agents?use=claude&start=first',
    conversation, '/home/new', `/home/gates/${ids.gate}`, '/home/receipts', `/home/receipts?id=${ids.message}`, '/home/legacy',
    conn('waiting'), conn('prepared'), conn('approved'), `/remote/consent/${ids.remote}`, '/'] },
  { name: 'ko', acct: 'ko', paths: ['/home/account', '/home/friends', '/home/agents'] },
  { name: 'none', acct: 'none', paths: ['/home/account'] },
  { name: 'none-cookie-en', acct: 'none', lang: 'en', paths: ['/home/account'] },
  { name: 'en-cookie-ko', acct: 'en', lang: 'ko', paths: ['/home/account'] },
  { name: 'anon-ko', paths: ['/', '/auth/verify', `/connect/${ids.device}`] },
  { name: 'anon-en', lang: 'en', paths: ['/', '/auth/verify', `/connect/${ids.device}`] },
];
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
const probe = `(() => {
  const w = innerWidth;
  const over = [...document.querySelectorAll('body *')].filter(e => e.getBoundingClientRect().right > w + 0.5).map(e => e.tagName + (e.className ? '.' + e.className : '')).slice(0, 5);
  const clipped = [...document.querySelectorAll('main *')].filter(e => e.scrollWidth > e.clientWidth + 1 && getComputedStyle(e).overflowX !== 'visible').map(e => e.tagName + ':' + e.textContent.trim().slice(0, 30)).slice(0, 5);
  const small = [...document.querySelectorAll('a[href], button, input:not([type=hidden]), select, summary, textarea')].filter(e => { const r = e.getBoundingClientRect(); return r.width > 0 && (r.height < 44 && !(e.tagName === 'A' && getComputedStyle(e).display === 'inline')); }).map(e => e.tagName + ':' + (e.textContent || e.value || e.name || '').trim().slice(0, 30) + ':' + Math.round(e.getBoundingClientRect().height)).slice(0, 12);
  const centered = [...document.querySelectorAll('main *')].filter(e => getComputedStyle(e).textAlign === 'center' && e.textContent.trim()).map(e => e.tagName).slice(0, 5);
  const text = document.body.innerText;
  const hangul = (text.replace(/한국어/g, '').match(/[\\uac00-\\ud7a3]+/g) || []).slice(0, 15);
  const langLinks = [...document.querySelectorAll('a[href^="/lang/"]')].map(e => { const r = e.getBoundingClientRect(); return { text: e.textContent.trim(), href: e.getAttribute('href'), current: e.getAttribute('aria-current'), hreflang: e.getAttribute('hreflang'), lang: e.getAttribute('lang'), top: Math.round(r.top + scrollY), left: Math.round(r.left), h: Math.round(r.height), inHeader: !!e.closest('header') }; });
  const radios = [...document.querySelectorAll('input[type=radio]')].map(e => e.name + '=' + e.value + (e.checked ? '*' : ''));
  const selects = document.querySelectorAll('select').length;
  const scripts = document.querySelectorAll('script').length;
  const nav = [...document.querySelectorAll('header a')].map(e => e.textContent.trim() + (e.getAttribute('aria-current') ? '[' + e.getAttribute('aria-current') + ']' : ''));
  const h = [...document.querySelectorAll('h1,h2')].map(e => e.tagName + ':' + e.textContent.trim().slice(0, 50));
  return { title: document.title, htmlLang: document.documentElement.lang, scrollWidth: document.documentElement.scrollWidth, innerWidth: w, height: document.documentElement.scrollHeight,
    over, clipped, small, centered, hangul, langLinks, radios, selects, scripts, nav, h, text };
})()`;
const results = [];
for (const c of contexts) {
  const { result: { browserContextId } } = await send('Target.createBrowserContext');
  const { result: { targetId } } = await send('Target.createTarget', { url: 'about:blank', browserContextId });
  const { result: { sessionId } } = await send('Target.attachToTarget', { targetId, flatten: true });
  // __Host- Secure cookies are dropped by Chrome over http, so the fixture session and the language cookie are set directly.
  const cookie = (name, value) => send('Network.setCookie', { name, value, url: origin + '/', path: '/', secure: true, httpOnly: true, sameSite: 'Strict' }, sessionId);
  if (c.acct) await cookie('__Host-kl_session', process.env['KL_' + c.acct.toUpperCase()]);
  if (c.lang) await cookie('__Host-kl_lang', c.lang);
  for (const path of c.paths) {
    for (const width of [1440, 390, 320]) {
      await send('Emulation.setDeviceMetricsOverride', { width, height: 900, deviceScaleFactor: 1, mobile: false }, sessionId);
      let v;
      // The fixture rate-limits each member and tester shares the accounts, so pace and back off.
      for (let attempt = 0; ; attempt++) {
        await send('Page.navigate', { url: origin + path }, sessionId);
        await sleep(4000);
        v = await evaluate(probe, sessionId);
        if (!/^(요청 제한|Too many|Rate)/i.test(v.title)) break;
        if (attempt === 6) throw new Error(`rate limited at ${c.name} ${path} ${width}`);
        await sleep(30000);
      }
      await send('Emulation.setDeviceMetricsOverride', { width, height: Math.min(v.height, 8000), deviceScaleFactor: 1, mobile: false }, sessionId);
      const shot = await send('Page.captureScreenshot', { format: 'png' }, sessionId);
      const name = `${c.name}${path.replace(/[^a-z0-9]+/gi, '-').replace(/-$/, '').slice(0, 60)}-${width}`;
      writeFileSync(join(out, `${name}.png`), Buffer.from(shot.result.data, 'base64'));
      if (width === 1440) writeFileSync(join(out, `${name}.txt`), v.text);
      delete v.text;
      results.push({ ctx: c.name, path, width, overflow: v.scrollWidth > v.innerWidth, png: `${name}.png`, ...v });
    }
  }
}
writeFileSync(join(out, process.env.PROBE || 'probe.json'), JSON.stringify(results, null, 2) + '\n');
console.log(results.filter((r) => r.overflow).length ? 'OVERFLOW' : 'NO_OVERFLOW', results.length);
proc.kill();
process.exit(0);
