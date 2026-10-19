// Designer recheck of DEV fix 5fa3333 on the names fixture: captures + layout facts for the 4 fixed screens.
// Usage: node probe.mjs <cdp-port>. Allowed GETs only, plus ONE POST /home/agent-name with a control character (refused, not saved).
import fs from 'node:fs';
import path from 'node:path';
import {fileURLToPath} from 'node:url';
import {createHash} from 'node:crypto';

const out = path.dirname(fileURLToPath(import.meta.url));
const port = process.argv[2];
const ORIGIN = 'http://127.0.0.1:18930';
const M = JSON.parse(fs.readFileSync('C:/Users/shin/AppData/Local/KnowsLinkDevTools/names-ux002-dev/fixture-manifest.json', 'utf8'));
const target = (await (await fetch(`http://127.0.0.1:${port}/json/list`)).json()).find(t => t.type === 'page');
const version = await (await fetch(`http://127.0.0.1:${port}/json/version`)).json();
const ws = new WebSocket(target.webSocketDebuggerUrl);
await new Promise((res, rej) => { ws.onopen = res; ws.onerror = rej; });
let next = 0; const pending = new Map(), loadWaiters = [], requests = [], result = {browser: version.Browser, origin: ORIGIN, screens: {}, captures: []};
ws.onmessage = ev => {
  const m = JSON.parse(ev.data);
  if (m.id) { const p = pending.get(m.id); pending.delete(m.id); m.error ? p.reject(new Error(m.error.message)) : p.resolve(m.result); }
  else if (m.method === 'Network.requestWillBeSent' && m.params.request.url.startsWith(ORIGIN)) requests.push(m.params.request.method + ' ' + new URL(m.params.request.url).pathname);
  else if (m.method === 'Network.responseReceived' && m.params.response.url.startsWith(ORIGIN)) requests.push(m.params.response.status + ' ' + new URL(m.params.response.url).pathname);
  else if (m.method === 'Page.loadEventFired') loadWaiters.splice(0).forEach(f => f());
};
const send = (method, params = {}) => new Promise((resolve, reject) => { const id = ++next; pending.set(id, {resolve, reject}); ws.send(JSON.stringify({id, method, params})); });
const js = async expr => { const r = await send('Runtime.evaluate', {expression: expr, returnByValue: true, awaitPromise: true}); if (r.exceptionDetails) throw new Error(JSON.stringify(r.exceptionDetails).slice(0, 300)); return r.result.value; };
const sleep = ms => new Promise(r => setTimeout(r, ms));
const save = () => fs.writeFileSync(path.join(out, 'probe.json'), JSON.stringify({...result, requests}, null, 2));
const loaded = () => { const w = new Promise(r => loadWaiters.push(r)); return Promise.race([w, sleep(8000)]).then(() => sleep(200)); };
async function go(url) {
  for (let i = 0; i < 12; i++) {
    const w = loaded(); await send('Page.navigate', {url}); await w;
    if (!(await js(`document.body.innerText.includes('요청이 많습니다')`))) return;
    await sleep(8000); // per-member /home request window
  }
  throw new Error('rate limited: ' + url);
}
async function key(k, code, vk) { await send('Input.dispatchKeyEvent', {type: 'keyDown', key: k, code, windowsVirtualKeyCode: vk}); await send('Input.dispatchKeyEvent', {type: 'keyUp', key: k, code, windowsVirtualKeyCode: vk}); await sleep(100); }
async function shot(name) {
  const w = await js('innerWidth'); const lm = await send('Page.getLayoutMetrics');
  const {data} = await send('Page.captureScreenshot', {format: 'png', captureBeyondViewport: true, clip: {x: 0, y: 0, width: w, height: Math.min(Math.ceil(lm.cssContentSize.height), 6000), scale: 1}});
  const b = Buffer.from(data, 'base64'); const f = `designer-${name}.png`;
  fs.writeFileSync(path.join(out, f), b); result.captures.push({file: f, sha256: createHash('sha256').update(b).digest('hex')});
}
// Page-wide facts: horizontal overflow and controls under 44px (radio/checkbox measured by their label).
const common = `(()=>{const vis=e=>e.getClientRects().length&&getComputedStyle(e).visibility!=='hidden';
  const hit=e=>((e.type==='radio'||e.type==='checkbox')?(e.closest('label')||e):e).getBoundingClientRect();
  const small=[...document.querySelectorAll('main a,main button,main input:not([type=hidden]),main select,main summary,main textarea')].filter(vis).filter(e=>hit(e).height<43.5).map(e=>e.tagName+':'+Math.round(hit(e).height)+':'+(e.textContent||e.value||'').trim().slice(0,12));
  return {over:document.documentElement.scrollWidth-innerWidth,small,selects:document.querySelectorAll('main select').length};})()`;
const newFacts = `(()=>[...document.querySelectorAll('main fieldset')].map(f=>({legend:f.querySelector('legend')?.innerText,
  radios:[...f.querySelectorAll('input[type=radio]')].map(r=>{const l=r.closest('label');const s=l.querySelector('span')||l;const lr=l.getBoundingClientRect();
    return {value:r.value,checked:r.checked,label:s.innerText.replace(/\\n/g,' '),h:Math.round(lr.height),right:Math.round(lr.right),cut:s.scrollWidth>s.clientWidth+1}})})))()`;
const agentFacts = `(()=>{document.querySelectorAll('main details').forEach(d=>d.open=true);
  const pick=t=>[...document.querySelectorAll('main a')].filter(a=>a.innerText.trim()===t).map(a=>{const r=a.getBoundingClientRect();return {top:Math.round(r.top),bottom:Math.round(r.bottom),h:Math.round(r.height),parent:a.parentElement.tagName}});
  return {receipts:pick('요청 상태 조회'),manage:pick('연결·키 관리')};})()`;
const friendFacts = `(()=>[...document.querySelectorAll('main button')].filter(b=>b.innerText.trim()==='거절').map(b=>{const id=b.getAttribute('aria-describedby');const d=id&&document.getElementById(id);
  return {describedby:id,text:d?.innerText,beforeButton:!!d&&!!(d.compareDocumentPosition(b)&Node.DOCUMENT_POSITION_FOLLOWING),color:d&&getComputedStyle(d).color,h:Math.round(b.getBoundingClientRect().height)}}))()`;

await send('Page.enable'); await send('Network.enable'); await send('Runtime.enable');
let sess = '';
async function as(who) {
  const r = await fetch(ORIGIN + '/fixture/as/' + who, {method: 'POST', headers: {Origin: ORIGIN}});
  sess = /__Host-kl_session=([^;]+)/.exec(r.headers.get('set-cookie'))[1];
  await send('Network.setCookie', {name: '__Host-kl_session', value: sess, url: ORIGIN + '/', secure: true, httpOnly: true, sameSite: 'Strict'});
}
const WIDTHS = [[1440, 900], [390, 844], [320, 640]];
const plan = [
  ['new-self', `/home/new?peer=self`, newFacts],
  ['new-peer', `/home/new?peer=${M.friend_ab}`, newFacts],
  ['agents', `/home/agents`, agentFacts],
  ['friends', `/home/friends`, friendFacts],
];
try {
  await as('a');
  for (const [w, h] of WIDTHS) {
    await send('Emulation.setDeviceMetricsOverride', {width: w, height: h, deviceScaleFactor: 1, mobile: false});
    for (const [name, p, facts] of plan) {
      await go(ORIGIN + p);
      const f = await js(facts); // agentFacts opens details before the capture
      result.screens[`${name}-${w}`] = {...(await js(common)), facts: f};
      await shot(`${name}-${w}`); save();
    }
  }
  // Keyboard on the self-conversation radios at 320: focus the checked sender, ArrowDown moves the check inside the group only.
  await go(ORIGIN + '/home/new?peer=self');
  await js(`document.querySelector('main fieldset input[type=radio]:checked').focus()`);
  const before = await js(newFacts); await key('ArrowDown', 'ArrowDown', 40); const after = await js(newFacts);
  result.keyboard = {before: before.map(g => g.radios.map(r => r.checked)), after: after.map(g => g.radios.map(r => r.checked)), focused: await js(`document.activeElement.value`)};
  // One refused rename: control character in the name. The relay must answer 422 and keep the stored name.
  await send('Emulation.setDeviceMetricsOverride', {width: 390, height: 844, deviceScaleFactor: 1, mobile: false});
  await go(ORIGIN + '/home/agents');
  const id = await js(`(()=>{const i=document.querySelector('main form input[name=name]');return i.id})()`);
  const stored = await js(`document.getElementById(${JSON.stringify(id)}).value`);
  const w = loaded();
  await js(`(()=>{const i=document.getElementById(${JSON.stringify(id)});i.value='bad\\u0007name';i.form.requestSubmit();})()`);
  await w; await sleep(300);
  result.refusal = {input: id, stored, ...(await js(`(()=>{const i=document.getElementById(${JSON.stringify(id)});const a=[...document.querySelectorAll('[role=alert]')].map(e=>e.innerText.trim());
    return {value:i?.value,invalid:i?.getAttribute('aria-invalid'),alerts:a,title:document.title}})()`)), ...(await js(common))};
  await shot('refusal-390');
  await go(ORIGIN + '/home/agents');
  result.refusal.storedAfter = await js(`document.getElementById(${JSON.stringify(id)}).value`);
  save(); console.log('done');
} catch (e) { result.error = String(e); save(); console.error(e); process.exitCode = 1; }
ws.close();
