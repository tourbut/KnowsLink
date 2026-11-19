// Independent browser checks for the receipt lookup form: real key events (CDP Input), viewport matrix, focus, a11y tree.
// Usage: node browser.mjs <cdp-port> <label> [short]   (synthetic fixture A only; read-only GETs except fixture session minting)
import fs from 'node:fs';
import path from 'node:path';
import {fileURLToPath} from 'node:url';
import {createHash} from 'node:crypto';

const out = path.dirname(fileURLToPath(import.meta.url));
const [port, label, mode] = [process.argv[2], process.argv[3], process.argv[4]];
const ORIGIN = 'http://127.0.0.1:18928';
const AB = '0199a3f2-4c10-7a11-8b22-334455667788', NONE = '0199a3f2-4c10-7a11-8b22-334455667799';
const version = await (await fetch(`http://127.0.0.1:${port}/json/version`)).json();
const target = (await (await fetch(`http://127.0.0.1:${port}/json/list`)).json()).find(t => t.type === 'page');
const ws = new WebSocket(target.webSocketDebuggerUrl);
await new Promise((res, rej) => { ws.onopen = res; ws.onerror = rej; });
let next = 0; const pending = new Map(), checks = [], captures = [], requests = [], loadWaiters = [];
const save = () => fs.writeFileSync(path.join(out, `browser-${label}.json`), JSON.stringify({browser: version.Browser, protocol: version['Protocol-Version'], origin: ORIGIN, mode: mode || 'full', checks, captures, requests}, null, 2));
ws.onmessage = ev => {
  const m = JSON.parse(ev.data);
  if (m.id) { const p = pending.get(m.id); pending.delete(m.id); m.error ? p.reject(new Error(m.error.message)) : p.resolve(m.result); }
  else if (m.method === 'Network.requestWillBeSent') { const u = new URL(m.params.request.url); requests.push({method: m.params.request.method, origin: u.origin, path: u.pathname}); }
  else if (m.method === 'Page.loadEventFired') { loadWaiters.splice(0).forEach(f => f()); }
};
const send = (method, params = {}) => new Promise((resolve, reject) => { const id = ++next; pending.set(id, {resolve, reject}); ws.send(JSON.stringify({id, method, params})); });
const js = async expr => { const r = await send('Runtime.evaluate', {expression: expr, returnByValue: true, awaitPromise: true}); if (r.exceptionDetails) throw new Error('eval failed: ' + expr.slice(0, 80)); return r.result.value; };
const sleep = ms => new Promise(r => setTimeout(r, ms));
function check(name, value, detail) { checks.push({name, passed: value === true, detail}); save(); if (value !== true) console.log('FAIL', name, JSON.stringify(detail ?? '')); else console.log('PASS', name); }
async function go(url) { const w = new Promise(r => loadWaiters.push(r)); await send('Page.navigate', {url}); await Promise.race([w, sleep(8000)]); await sleep(150); }
async function key(k, code, vk, mods = 0) { await send('Input.dispatchKeyEvent', {type: 'keyDown', key: k, code, windowsVirtualKeyCode: vk, modifiers: mods, text: k === 'Enter' ? '\r' : undefined}); await send('Input.dispatchKeyEvent', {type: 'keyUp', key: k, code, windowsVirtualKeyCode: vk, modifiers: mods}); await sleep(120); }
const tab = () => key('Tab', 'Tab', 9), shiftTab = () => key('Tab', 'Tab', 9, 8), enter = async () => { const w = new Promise(r => loadWaiters.push(r)); await key('Enter', 'Enter', 13); await Promise.race([w, sleep(1500)]); await sleep(150); };
async function selectAll() { await key('a', 'KeyA', 65, 2); }
const active = () => js(`(()=>{const e=document.activeElement;const cs=getComputedStyle(e);return {tag:e.tagName,id:e.id,cls:e.className,href:e.getAttribute('href'),text:(e.textContent||'').trim().slice(0,24),fv:e.matches(':focus-visible'),outline:cs.outlineStyle+' '+cs.outlineWidth+' '+cs.outlineColor,h:Math.round(e.getBoundingClientRect().height)}})()`);
async function shot(name) { const m = await js('({w:innerWidth,h:innerHeight,dpr:devicePixelRatio})'); const {data} = await send('Page.captureScreenshot', {format: 'png'}); const b = Buffer.from(data, 'base64'); fs.writeFileSync(path.join(out, `${label}-${name}.png`), b); captures.push({file: `${label}-${name}.png`, ...m, sha256: createHash('sha256').update(b).digest('hex')}); save(); }
async function metrics(w, h, dsf = 1) { await send('Emulation.setDeviceMetricsOverride', {width: w, height: h, deviceScaleFactor: dsf, mobile: false}); }
const layout = () => js(`({iw:innerWidth,dpr:devicePixelRatio,over:document.documentElement.scrollWidth-innerWidth,small:[...document.querySelectorAll('a,button,input,select,summary')].filter(e=>e.getClientRects().length&&e.getBoundingClientRect().height<44).map(e=>e.tagName+'#'+e.id+':'+Math.round(e.getBoundingClientRect().height)),clipped:[...document.querySelectorAll('main *')].filter(e=>e.getClientRects().length&&e.getBoundingClientRect().right>innerWidth+1).map(e=>e.tagName+'#'+e.id).slice(0,5)})`);

await send('Page.enable'); await send('Network.enable'); await send('Runtime.enable'); await send('Accessibility.enable');
// synthetic session (fixture-only POSTs, not a real login)
let sess = ''; const post = async p => (await fetch(ORIGIN + p, {method: 'POST', headers: {Origin: ORIGIN, ...(sess && {Cookie: '__Host-kl_session=' + sess})}}));
await post('/fixture/recover');
const r = await post('/fixture/as/a'); const token = /__Host-kl_session=([^;]+)/.exec(r.headers.get('set-cookie'))[1];
sess = token; await send('Network.setCookie', {name: '__Host-kl_session', value: token, url: ORIGIN + '/', secure: true, httpOnly: true, sameSite: 'Strict'});
await send('Network.deleteCookies', {name: '__Host-kl_receipt', url: ORIGIN + '/'}); // start without a recovery context
const URLQ = (a, i) => `${ORIGIN}/home/receipts?agent=${a}&id=${encodeURIComponent(i)}`;

try {
  await metrics(1440, 900);
  // --- first view: autofocus + tab order + label + skip
  await go(ORIGIN + '/home/receipts');
  let a = await active();
  check('first-view-autofocus-agent-select', a.id === 'receipt-agent', a);
  check('labels-bound', await js(`document.getElementById('receipt-agent').labels[0].textContent==='내 에이전트'&&document.getElementById('receipt-id').labels[0].textContent==='요청 ID 또는 답장 ID'&&document.querySelector('#receipt-form button').textContent==='상태 확인'`));
  const order = [];
  for (let i = 0; i < 4; i++) { await tab(); const x = await active(); order.push(x.id || x.href || x.text); if (i === 0) check('tab-moves-to-id-with-visible-focus', x.id === 'receipt-id' && x.fv && parseFloat(x.outline.split(' ')[1]) >= 2 && !x.outline.startsWith('none'), x); }
  check('tab-order-id-button-then-links', order[1] === '상태 확인' && order.slice(2).every(v => v.startsWith('/')), order);
  await shiftTab(); await shiftTab(); await shiftTab(); a = await active();
  check('shift-tab-returns-to-id', a.id === 'receipt-id', a);
  // skip link: first focusable in DOM order, reached by Shift+Tab from the autofocused field
  let hops = 0; do { await shiftTab(); a = await active(); hops++; } while (a.cls !== 'skip' && hops < 8);
  check('skip-link-reachable-and-visible-on-focus', a.cls === 'skip' && a.fv && await js(`document.activeElement.getBoundingClientRect().top>=0&&document.activeElement.getBoundingClientRect().height>=30`), {hops, a});
  await key('Enter', 'Enter', 13); await sleep(200); check('skip-link-targets-main', await js(`location.hash==='#main'`)); await tab(); a = await active();
  check('after-skip-next-tab-is-inside-main', await js(`!!document.activeElement.closest('main')`), a);

  // --- keyboard edit loop: bad -> alert -> fix -> success -> return link -> back/refresh
  await go(ORIGIN + '/home/receipts'); // autofocus is on the agent select: choose with ArrowDown, Tab to id, type
  await key('ArrowDown', 'ArrowDown', 40); await tab(); await send('Input.insertText', {text: 'bad-id'}); await enter();
  check('invalid-submit-keeps-url-input-select', await js(`location.search.includes('id=bad-id')&&document.getElementById('receipt-id').value==='bad-id'`));
  a = await active(); check('invalid-focus-on-field-with-alert', a.id === 'receipt-id' && await js(`!!document.querySelector('[role=alert]')&&document.getElementById('receipt-id').getAttribute('aria-invalid')==='true'`), a);
  check('error-link-reaches-input', await js(`(document.querySelector('[role=alert] a')||{getAttribute:()=>null}).getAttribute('href')==='#receipt-id'`));
  await selectAll(); await send('Input.insertText', {text: AB}); await key('Tab', 'Tab', 9); await shiftTab();
  await enter(); a = await active();
  check('fixed-input-resubmits-same-screen', await js(`location.pathname==='/home/receipts'&&document.getElementById('receipt-id').value==='${AB}'&&document.querySelector('#receipt-result dt').textContent==='전달'`), a);
  check('result-region-status-labelled', await js(`(()=>{const s=document.getElementById('receipt-result');return s.getAttribute('role')==='status'&&s.getAttribute('aria-labelledby')==='receipt-result-h'})()`));
  await shot('success-1440');
  // return link
  const ret = await js(`[...document.querySelectorAll('#receipt-result a')].find(e=>e.textContent.includes('조회 입력으로 돌아가기'))!==undefined`);
  check('return-to-input-link-present', ret === true);
  await js(`[...document.querySelectorAll('#receipt-result a')].find(e=>e.textContent.includes('조회 입력으로 돌아가기')).focus()`); await key('Enter', 'Enter', 13); await sleep(200);
  a = await active(); check('return-to-input-focuses-input-keeps-value', a.id === 'receipt-id' && await js(`document.getElementById('receipt-id').value==='${AB}'`), a);
  // refresh + back
  await js('location.reload()'); await sleep(800);
  check('refresh-keeps-selection-and-id', await js(`document.getElementById('receipt-id').value==='${AB}'&&document.getElementById('receipt-agent').value==='agent_a'`));
  await go(ORIGIN + '/home/receipts?agent=agent_a&id=' + NONE); await js('history.back()'); await sleep(900);
  check('back-restores-previous-query-input', await js(`location.search.includes('${AB}')&&document.getElementById('receipt-id').value==='${AB}'`));
  // bare entry (e.g. from nav) restores via recovery context
  await go(ORIGIN + '/home/receipts');
  check('bare-entry-restores-last-lookup', await js(`document.getElementById('receipt-agent').value==='agent_a'&&document.getElementById('receipt-id').value!==''`));
  // navigating away and returning via internal return link goes to a safe internal page
  check('return-links-are-internal', await js(`[...document.querySelectorAll('main a')].every(a=>new URL(a.href).origin===location.origin)`));

  // --- refusal state with keyboard retry
  await go(URLQ('agent_a', NONE));
  a = await active(); check('refusal-alert-focus-and-message', await js(`document.querySelector('[role=alert]').textContent.includes('이 요청의 상태를 확인할 수 없습니다. 내 에이전트와 요청 ID를 확인하세요.')`), a);
  check('refusal-keeps-select-and-id-editable', await js(`document.getElementById('receipt-agent').value==='agent_a'&&document.getElementById('receipt-id').value==='${NONE}'&&!document.getElementById('receipt-id').readOnly`));

  // --- accessibility tree (landmarks, heading, names)
  await go(URLQ('agent_a', AB));
  const ax = (await send('Accessibility.getFullAXTree')).nodes.filter(n => !n.ignored).map(n => ({role: n.role?.value, name: n.name?.value}));
  const has = (role, re) => ax.some(n => n.role === role && re.test(n.name || ''));
  check('ax-landmarks-heading-form', has('banner', /.*/) && has('main', /.*/) && has('heading', /^요청 상태 조회$/) && has('combobox', /^내 에이전트$/) && has('textbox', /^요청 ID 또는 답장 ID/) && has('button', /^상태 확인$/) && has('navigation', /현재 위치/), ax.filter(n => ['banner', 'main', 'heading', 'combobox', 'textbox', 'button', 'navigation'].includes(n.role)));
  check('ax-current-page-link', await js(`document.querySelector('nav[aria-label="현재 위치"] [aria-current="page"]').textContent==='요청 상태 조회'`));

  // --- contrast of rendered colors
  const lum = c => { const [r, g, b] = c.match(/\d+(\.\d+)?/g).slice(0, 3).map(Number).map(v => { v /= 255; return v <= .03928 ? v / 12.92 : ((v + .055) / 1.055) ** 2.4; }); return .2126 * r + .7152 * g + .0722 * b; };
  const cr = (a, b) => { const [x, y] = [lum(a), lum(b)].sort((p, q) => q - p); return (x + .05) / (y + .05); };
  const colors = await js(`(()=>{const g=(s,p)=>getComputedStyle(document.querySelector(s))[p];return {bg:g('body','backgroundColor'),text:g('main p','color'),muted:g('#receipt-help','color'),inputBorder:g('#receipt-id','borderTopColor'),btnBg:g('.primary','backgroundColor'),btnText:g('.primary','color'),link:g('main a','color')}})()`);
  await go(URLQ('agent_a', NONE));
  const prob = await js(`({bg:getComputedStyle(document.querySelector('.problem')).backgroundColor,fg:getComputedStyle(document.querySelector('.problem strong')||document.querySelector('.problem')).color,edge:getComputedStyle(document.querySelector('.problem')).borderLeftColor})`);
  const ratios = {text: cr(colors.text, colors.bg), muted: cr(colors.muted, colors.bg), link: cr(colors.link, colors.bg), button: cr(colors.btnText, colors.btnBg), inputBorder: cr(colors.inputBorder, colors.bg), problemText: cr(prob.fg === colors.text ? colors.text : prob.fg, prob.bg), problemEdge: cr(prob.edge, prob.bg)};
  check('contrast-text-4.5-nontext-3', ratios.text >= 4.5 && ratios.muted >= 4.5 && ratios.link >= 4.5 && ratios.button >= 4.5 && ratios.inputBorder >= 3 && ratios.problemEdge >= 3, ratios);

  // --- viewport matrix per state
  const states = {first: ORIGIN + '/home/receipts?agent=&id=x'.replace('?agent=&id=x', ''), success: URLQ('agent_a', AB), refusal: URLQ('agent_a', NONE), invalid: URLQ('agent_a', 'x'.repeat(100) + '!bad')};
  const dims = mode === 'short' ? [[390, 844, 1]] : [[1440, 900, 1], [390, 844, 1], [320, 640, 1], [720, 800, 2]];
  for (const [w, h, d] of dims) {
    await metrics(w, h, d);
    for (const [n, u] of Object.entries(states)) {
      await go(u); if (n === 'success') await js(`document.querySelectorAll('details').forEach(e=>e.open=true)`);
      const l = await layout();
      check(`layout-${w}x${d}-${n}`, l.iw === w && l.over <= 0 && l.small.length === 0 && l.clipped.length === 0, l);
      if (d === 1 && w < 1000 && ['success', 'refusal'].includes(n)) await shot(`${n}-${w}`);
    }
    if (mode !== 'short') { // keyboard focus visible and not clipped at narrow widths
      await go(states.refusal); await js('document.activeElement.blur()'); await tab(); await tab(); await tab(); const x = await active();
      check(`focus-visible-${w}x${d}`, x.fv && !x.outline.startsWith('none') && await js(`(()=>{const r=document.activeElement.getBoundingClientRect();return r.left>=0&&r.right<=innerWidth})()`), x);
    }
  }
  await metrics(1440, 900, 1);
  check('details-toggle-by-keyboard', await (async () => { await go(URLQ('agent_a', AB)); await js(`document.querySelector('details summary').focus()`); await key(' ', 'Space', 32); await sleep(150); const open = await js(`document.querySelector('details').open`); await key('Enter', 'Enter', 13); await sleep(150); const closed = await js(`!document.querySelector('details').open`); return open && closed; })());
  // closed-selector visibility of the stable distinguisher (U2: name + kind + stable value must be told apart)
  for (const [w, h, d] of mode === 'short' ? [] : [[1440, 900, 1], [390, 844, 1], [320, 640, 1]]) {
    await metrics(w, h, d); await go(URLQ('agent_a', AB));
    const sel = await js(`(()=>{const s=document.getElementById('receipt-agent');const t=s.options[s.selectedIndex].text;const cs=getComputedStyle(s);const c=document.createElement('canvas').getContext('2d');c.font=cs.fontStyle+' '+cs.fontWeight+' '+cs.fontSize+' '+cs.fontFamily;const avail=s.clientWidth-parseFloat(cs.paddingLeft)-parseFloat(cs.paddingRight)-24;return {text:t,textWidth:Math.round(c.measureText(t).width),avail:Math.round(avail)}})()`);
    check(`selector-closed-shows-full-name-kind-distinguisher-${w}`, sel.textWidth <= sel.avail, sel);
  }
  await metrics(1440, 900, 1);
  // 429 in a real browser at narrow widths (limit then recover; synthetic member only)
  if (mode !== 'short') {
    await post('/fixture/limit');
    for (const [w, h] of [[390, 844], [320, 640]]) { await metrics(w, h); await go(URLQ('agent_a', AB)); const l = await layout(); const st = await js(`({v:document.getElementById('receipt-id').value,a:(document.querySelector('[role=alert]')||{}).textContent})`); check(`layout-429-${w}`, l.over <= 0 && l.small.length === 0 && await js(`document.getElementById('receipt-id').value==='${AB}'&&!!document.querySelector('[role=alert]')&&document.querySelector('[role=alert]').textContent.includes('다시 시도하세요')`), {l, st}); }
    await post('/fixture/recover'); await metrics(1440, 900, 1);
  }
  check('requests-same-origin-get-only', requests.every(q => q.origin === ORIGIN && q.method === 'GET'), requests.filter(q => q.origin !== ORIGIN || q.method !== 'GET'));
} finally { save(); ws.close(); }
const bad = checks.filter(c => !c.passed); console.log(`checks ${checks.length} failed ${bad.length}`); process.exit(bad.length ? 1 : 0);
