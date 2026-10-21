// Headless Chrome check for the synthetic calendar fixture. Reads rendered text and horizontal overflow.
// Usage: node browser.mjs <origin> <stage>   stage = early | reconnect | off
import {spawn} from 'node:child_process';
import fs from 'node:fs';
import net from 'node:net';
import path from 'node:path';
import {fileURLToPath} from 'node:url';

const out = path.dirname(fileURLToPath(import.meta.url));
const manifest = JSON.parse(fs.readFileSync('C:/Users/shin/AppData/Local/KnowsLinkDevTools/schedule001-qa/browser-fixture.json', 'utf8'));
const [origin, stage] = [process.argv[2], process.argv[3]];
if (!origin.startsWith('http://127.0.0.1:')) throw new Error('loopback origin required');
const chrome = 'C:/Program Files/Google/Chrome/Application/chrome.exe';
const port = await new Promise((resolve, reject) => {
  const probe = net.createServer();
  probe.listen(0, '127.0.0.1', () => { const chosen = probe.address().port; probe.close(() => resolve(chosen)); });
  probe.on('error', reject);
});
const profile = path.join(process.env.TEMP || '/tmp', `schedqa-chrome-${stage}-${Date.now()}`);
const child = spawn(chrome, [`--headless=new`, `--disable-gpu`, `--remote-debugging-port=${port}`, `--user-data-dir=${profile}`, '--no-first-run', '--disable-extensions', 'about:blank'], {stdio: 'ignore'});
const sleep = ms => new Promise(r => setTimeout(r, ms));
let version;
for (let i = 0; i < 40; i++) {
  try { version = await (await fetch(`http://127.0.0.1:${port}/json/version`)).json(); break; } catch { await sleep(150); }
}
if (!version) throw new Error('chrome debug port did not open');
const target = (await (await fetch(`http://127.0.0.1:${port}/json/list`)).json()).find(t => t.type === 'page');
const ws = new WebSocket(target.webSocketDebuggerUrl);
await new Promise((res, rej) => { ws.onopen = res; ws.onerror = rej; });
let next = 0;
const pending = new Map();
const loadWaiters = [];
ws.onmessage = ev => {
  const message = JSON.parse(ev.data);
  if (message.method === 'Page.loadEventFired') loadWaiters.splice(0).forEach(fn => fn());
  if (!message.id) return;
  const waiter = pending.get(message.id);
  pending.delete(message.id);
  message.error ? waiter.reject(new Error(message.error.message)) : waiter.resolve(message.result);
};
const send = (method, params = {}) => new Promise((resolve, reject) => { const id = ++next; pending.set(id, {resolve, reject}); ws.send(JSON.stringify({id, method, params})); });
const js = async expr => {
  const result = await send('Runtime.evaluate', {expression: expr, returnByValue: true, awaitPromise: true});
  if (result.exceptionDetails) throw new Error('eval failed');
  return result.result.value;
};
await send('Network.enable');
await send('Page.enable');
await send('Emulation.setDeviceMetricsOverride', {width: 1440, height: 900, deviceScaleFactor: 1, mobile: false});

// Page CSP is default-src 'none', so sign-in is a Node request. Cookie values stay in memory.
async function sign(who) {
  const response = await fetch(origin + '/fixture/as/' + who, {method: 'POST', headers: {Origin: origin}});
  if (response.status !== 204) throw new Error('sign-in status ' + response.status);
  const headers = typeof response.headers.getSetCookie === 'function' ? response.headers.getSetCookie() : [];
  if (headers.length === 0) throw new Error('sign-in returned no cookies');
  let sawSession = false;
  const pairs = [];
  for (const header of headers) {
    const name = header.split('=', 1)[0];
    const value = header.slice(name.length + 1).split(';', 1)[0];
    const hostCookie = name.startsWith('__Host-');
    const stored = await send('Network.setCookie', {
      name, value, url: origin + '/', path: '/', httpOnly: true,
      secure: hostCookie, sameSite: hostCookie ? 'Strict' : 'Lax',
    });
    if (!stored.success) throw new Error('browser refused a sign-in cookie');
    pairs.push(name + '=' + value);
    if (name === '__Host-kl_session') sawSession = true;
  }
  if (!sawSession) throw new Error('sign-in missing session cookie');
  const probe = await fetch(origin + '/home/account', {redirect: 'manual', headers: {Origin: origin, Cookie: pairs.join('; ')}});
  if (probe.status !== 200) throw new Error('session probe ' + probe.status);
}
async function open(url, lang) {
  if (lang) await send('Network.setCookie', {name: '__Host-kl_lang', value: lang, url: origin + '/', secure: true, httpOnly: true, path: '/', sameSite: 'Strict'});
  else await send('Network.deleteCookies', {name: '__Host-kl_lang', url: origin + '/'});
  let text = '';
  for (let i = 0; i < 8; i++) {
    const done = new Promise(resolve => loadWaiters.push(resolve));
    await send('Page.navigate', {url});
    await Promise.race([done, sleep(8000)]);
    await sleep(200);
    text = await js(`(() => { document.querySelectorAll('details').forEach(e => { e.open = true; }); return document.body ? document.body.innerText : ''; })()`);
    if (!text.includes('요청이 많습니다')) return text;
    await sleep(8000);
  }
  return text;
}
async function overflow(width) {
  await send('Emulation.setDeviceMetricsOverride', {width, height: 900, deviceScaleFactor: 1, mobile: width < 800});
  await sleep(150);
  return js('({width: window.innerWidth, extra: document.documentElement.scrollWidth - window.innerWidth})');
}

const gates = manifest.gates;
const checks = [];
function check(name, ok, detail) { checks.push({name, passed: !!ok, detail: ok ? '' : String(detail).slice(0, 240)}); console.log(ok ? 'PASS' : 'FAIL', name); }

try {
  if (stage === 'early') {
    await sign('b');
    let text = await open(origin + '/home/account');
    check('browser-b-connected-ko', text.includes('연결됨 — 바쁜 시간만 확인합니다') && text.includes('아직 상대에게 전달되지 않은 조회 요청과 결과'), text);
    text = await open(origin + '/home/account', 'en');
    check('browser-b-connected-en', text.includes('Connected — only busy times are checked') && text.includes('Results already delivered cannot be taken back.'), text);
    text = await open(origin + '/home/gates/' + gates.connected_pending_2h_30);
    check('browser-gate-ko', text.includes('승인하면 이 요청의 기간에 대해 바쁜 시간대만 한 번 공유합니다') && text.includes('이 요청 승인') && !text.includes('가능 시간 공개 정책이 없어'), text);
    text = await open(origin + '/home/gates/' + gates.connected_pending_2h_30, 'en');
    check('browser-gate-en', text.includes('Approving shares only your busy time slots') && text.includes('Approve this request'), text);
    for (const width of [1440, 390, 320]) {
      const box = await overflow(width);
      check('browser-gate-overflow-' + width, box.extra <= 0 && box.width === width, JSON.stringify(box));
    }
    await sign('c');
    text = await open(origin + '/home/gates/' + gates.unconnected_pending);
    check('browser-v1-ko', text.includes('이 요청은 승인할 수 없습니다') && text.includes('캘린더가 연결되지 않아') && !text.includes('이 요청 승인'), text);
    text = await open(origin + '/home/gates/' + gates.unconnected_pending, 'en');
    check('browser-v1-en', text.includes('This request cannot be approved') && !text.includes('Approve this request'), text);
    for (const width of [1440, 390, 320]) {
      const box = await overflow(width);
      check('browser-unconnected-overflow-' + width, box.extra <= 0 && box.width === width, JSON.stringify(box));
    }
  } else if (stage === 'reconnect') {
    await sign('b');
    const text = await open(origin + '/home/account');
    check('browser-v3-ko', text.includes('다시 연결 필요') && text.includes('아직 상대에게 전달되지 않은 조회 요청과 결과') && text.includes('캘린더 다시 연결'), text);
    const en = await open(origin + '/home/account', 'en');
    check('browser-v3-en', en.includes('Reconnection needed') && en.includes('Disconnecting discards the token'), en);
  } else if (stage === 'off') {
    await sign('b');
    const text = await open(origin + '/home/account?n=calendar-off');
    check('browser-off-ko', text.includes('캘린더 연결을 해제했습니다') && text.includes('아직 전달되지 않은 조회 요청과 결과') && !text.includes('연결됨 — 바쁜 시간만 확인합니다'), text);
    const en = await open(origin + '/home/account?n=calendar-off', 'en');
    check('browser-off-en', en.includes('Calendar disconnected.') && en.includes('not yet delivered ended as declined'), en);
  } else throw new Error('unknown stage');
} finally {
  const file = path.join(out, `browser-${stage}.json`);
  fs.writeFileSync(file, JSON.stringify({browser: version.Browser, stage, origin, checks, fail: checks.filter(c => !c.passed).length}, null, 2));
  ws.close();
  child.kill();
}
console.log('FAILS=' + checks.filter(c => !c.passed).length);
process.exit(checks.some(c => !c.passed) ? 1 : 0);
