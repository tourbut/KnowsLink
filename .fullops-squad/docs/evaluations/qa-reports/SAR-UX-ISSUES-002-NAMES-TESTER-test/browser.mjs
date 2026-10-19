// Independent browser checks for the names UX (U2/U4/U7): real key events (CDP Input), 3-width layout matrix, focus walk, a11y names, captures.
// Usage: node browser.mjs <cdp-port> <label> <stage>   stage = layout | keys | after | recheck (5fa3333 fix recheck)
// Synthetic A/B/C only. layout/keys change no relationship; keys submits ONE friend request with a non-existent code (harmless no-op by contract).
import fs from 'node:fs';
import path from 'node:path';
import {fileURLToPath} from 'node:url';
import {createHash} from 'node:crypto';

const out = path.dirname(fileURLToPath(import.meta.url));
const [port, label, stage] = [process.argv[2], process.argv[3], process.argv[4] || 'layout'];
const ORIGIN = 'http://127.0.0.1:18930';
const M = JSON.parse(fs.readFileSync(process.env.NAMES_MANIFEST || 'C:/Users/shin/AppData/Local/KnowsLinkDevTools/names-ux002-dev/fixture-manifest.json', 'utf8'));
const AB = M.friend_ab, CA = M.request_ca;
const NONE = '0199a3f2-4c10-7a11-8b22-334455667799';
const version = await (await fetch(`http://127.0.0.1:${port}/json/version`)).json();
const target = (await (await fetch(`http://127.0.0.1:${port}/json/list`)).json()).find(t => t.type === 'page');
const ws = new WebSocket(target.webSocketDebuggerUrl);
await new Promise((res, rej) => { ws.onopen = res; ws.onerror = rej; });
let next = 0; const pending = new Map(), checks = [], captures = [], requests = [], loadWaiters = [], obs = {};
const resultFile = path.join(out, `browser-${label}-${stage}.json`);
const save = () => fs.writeFileSync(resultFile, JSON.stringify({browser: version.Browser, protocol: version['Protocol-Version'], origin: ORIGIN, stage, checks, observations: obs, captures, requests}, null, 2));
ws.onmessage = ev => {
  const m = JSON.parse(ev.data);
  if (m.id) { const p = pending.get(m.id); pending.delete(m.id); m.error ? p.reject(new Error(m.error.message)) : p.resolve(m.result); }
  else if (m.method === 'Network.requestWillBeSent') { const u = new URL(m.params.request.url); requests.push({method: m.params.request.method, origin: u.origin, path: u.pathname}); }
  else if (m.method === 'Page.loadEventFired') loadWaiters.splice(0).forEach(f => f());
};
const send = (method, params = {}) => new Promise((resolve, reject) => { const id = ++next; pending.set(id, {resolve, reject}); ws.send(JSON.stringify({id, method, params})); });
const js = async expr => { const r = await send('Runtime.evaluate', {expression: expr, returnByValue: true, awaitPromise: true}); if (r.exceptionDetails) throw new Error('eval failed: ' + expr.slice(0, 100) + ' :: ' + JSON.stringify(r.exceptionDetails.exception?.description || '').slice(0, 200)); return r.result.value; };
const sleep = ms => new Promise(r => setTimeout(r, ms));
function check(name, value, detail) { checks.push({name, passed: value === true, detail}); save(); console.log(value === true ? 'PASS' : 'FAIL', name, value === true ? '' : JSON.stringify(detail ?? '').slice(0, 400)); }
const rateLimited = () => js(`document.body.innerText.includes('요청이 많습니다')`);
async function go(url) {
  for (let i = 0; i < 12; i++) {
    const w = new Promise(r => loadWaiters.push(r)); await send('Page.navigate', {url}); await Promise.race([w, sleep(8000)]); await sleep(150);
    if (!(await rateLimited()) || url.includes('/fixture/')) return;
    obs.memberWindow429Waits = (obs.memberWindow429Waits || 0) + 1; await sleep(8000); // per-member request window shared by all /home pages
  }
}
async function key(k, code, vk, mods = 0, text) { await send('Input.dispatchKeyEvent', {type: 'keyDown', key: k, code, windowsVirtualKeyCode: vk, modifiers: mods, text: text ?? (k === 'Enter' ? '\r' : undefined)}); await send('Input.dispatchKeyEvent', {type: 'keyUp', key: k, code, windowsVirtualKeyCode: vk, modifiers: mods}); await sleep(100); }
const tab = () => key('Tab', 'Tab', 9), shiftTab = () => key('Tab', 'Tab', 9, 8), space = () => key(' ', 'Space', 32, 0, ' ');
const enterNav = async () => { const w = new Promise(r => loadWaiters.push(r)); await key('Enter', 'Enter', 13); await Promise.race([w, sleep(2500)]); await sleep(200); };
const selectAll = () => key('a', 'KeyA', 65, 2);
const metrics = (w, h) => send('Emulation.setDeviceMetricsOverride', {width: w, height: h, deviceScaleFactor: 1, mobile: false});
const slug = s => s.replace(/[^a-z0-9]+/gi, '-').replace(/^-|-$/g, '').slice(0, 40);
async function shot(name, full = true) {
  const m = await js('({w:innerWidth,h:innerHeight,dpr:devicePixelRatio})'); let params = {format: 'png'};
  if (full) { const lm = await send('Page.getLayoutMetrics'); params = {format: 'png', captureBeyondViewport: true, clip: {x: 0, y: 0, width: m.w, height: Math.min(Math.ceil(lm.cssContentSize.height), 6000), scale: 1}}; }
  const {data} = await send('Page.captureScreenshot', params); const b = Buffer.from(data, 'base64'); const f = `${label}-${name}.png`;
  fs.writeFileSync(path.join(out, f), b); captures.push({file: f, ...m, sha256: createHash('sha256').update(b).digest('hex')}); save();
}
const layout = () => js(`(()=>{
  const vis = e => e.getClientRects().length && getComputedStyle(e).visibility !== 'hidden';
  const hit = e => { const l = (e.type==='checkbox'||e.type==='radio') ? (e.closest('label')||e) : e; return l.getBoundingClientRect(); };
  const ctl = [...document.querySelectorAll('a,button,input:not([type=hidden]),select,summary,textarea')].filter(vis);
  const small = ctl.filter(e => hit(e).height < 43.5).map(e => e.tagName + '#' + (e.id || e.name || '') + ':' + Math.round(hit(e).height) + ':' + (e.textContent || e.value || '').trim().slice(0, 14));
  const clipped = [...document.querySelectorAll('body *')].filter(e => vis(e) && !e.closest('svg') && (e.getBoundingClientRect().right > innerWidth + 1 || e.getBoundingClientRect().left < -1) && getComputedStyle(e).position !== 'fixed').map(e => e.tagName + '#' + e.id + '.' + e.className).slice(0, 6);
  const hidden = [...document.querySelectorAll('main *')].filter(e => vis(e) && !/^(INPUT|TEXTAREA|SELECT)$/.test(e.tagName) && ['hidden','clip'].includes(getComputedStyle(e).overflowX) && e.scrollWidth > e.clientWidth + 1).map(e => e.tagName + '.' + e.className).slice(0, 6);
  const inputScroll = [...document.querySelectorAll('main input:not([type=checkbox]):not([type=radio]):not([type=hidden])')].filter(e => vis(e) && e.scrollWidth > e.clientWidth + 1).map(e => e.id + ':' + (e.value || '').length + 'ch');
  const ell = [...document.querySelectorAll('main *')].filter(e => vis(e) && getComputedStyle(e).textOverflow === 'ellipsis').map(e => e.tagName + '.' + e.className).slice(0, 6);
  const names = [...document.querySelectorAll('.name-preview,.wrap,main h1,main h2,.choice span,.consent span,summary h3')].filter(vis);
  const nameOver = names.filter(e => e.scrollWidth > e.clientWidth + 1 || e.getBoundingClientRect().right > innerWidth + 1).map(e => e.tagName + '.' + e.className + ':' + e.scrollWidth + '>' + e.clientWidth);
  const ta = document.getElementById('my-code'); const taClip = ta ? ta.scrollHeight > ta.clientHeight + 1 || ta.scrollWidth > ta.clientWidth + 1 : false;
  const selects = [...document.querySelectorAll('main select')].filter(vis).map(e => { const t = e.options[e.selectedIndex].text; const cs = getComputedStyle(e); const c = document.createElement('canvas').getContext('2d'); c.font = cs.fontStyle + ' ' + cs.fontWeight + ' ' + cs.fontSize + ' ' + cs.fontFamily; const avail = e.clientWidth - parseFloat(cs.paddingLeft) - parseFloat(cs.paddingRight) - 24; return {id: e.id, full: Math.round(c.measureText(t).width), avail: Math.round(avail)}; });
  return {iw: innerWidth, dpr: devicePixelRatio, inputScroll, selects, over: document.documentElement.scrollWidth - innerWidth, small, clipped, hidden, ell, nameOver, taClip, nameCount: names.length};
})()`);
const tabStops = async (cap = 70) => {
  await js('document.activeElement && document.activeElement.blur(); window.scrollTo(0,0)');
  const stops = []; const ax = (await send('Accessibility.getFullAXTree')).nodes; const byId = new Map(ax.filter(n => n.backendDOMNodeId).map(n => [n.backendDOMNodeId, n]));
  for (let i = 0; i < cap; i++) {
    await tab();
    const info = await js(`(()=>{const e=document.activeElement;if(e===document.body||e===document.documentElement)return {body:true};const cs=getComputedStyle(e);const r=e.getBoundingClientRect();const l=(e.type==='checkbox'||e.type==='radio')?(e.closest('label')||e):e;const lr=l.getBoundingClientRect();return {key:e.tagName+'|'+e.id+'|'+(e.getAttribute('href')||'')+'|'+(e.textContent||e.value||'').trim().slice(0,26),tag:e.tagName,id:e.id,cls:e.className,type:e.type||'',text:(e.textContent||e.value||'').trim().slice(0,26),fv:e.matches(':focus-visible'),outline:cs.outlineStyle+' '+cs.outlineWidth,shadow:cs.boxShadow!=='none',h:Math.round(lr.height),inx:r.left>=-1&&r.right<=innerWidth+1}})()`);
    if (info.body) { stops.push({wrap: true}); break; }
    if (stops.length > 1 && info.key === stops[0].key) { stops.push({wrap: true}); break; } // focus came back to the first stop: one full cycle, no trap (headless Edge cycles inside the page instead of leaving to <body>)
    const {objectId} = (await send('Runtime.evaluate', {expression: 'document.activeElement'})).result; const bid = (await send('DOM.describeNode', {objectId})).node.backendNodeId;
    const n = byId.get(bid); info.axName = n?.name?.value ?? null; info.axRole = n?.role?.value ?? null; stops.push(info);
  }
  return stops;
};
const focusOk = s => s.fv && ((!s.outline.startsWith('none') && parseFloat(s.outline.split(' ')[1]) >= 2) || s.shadow) && s.inx;

await send('Page.enable'); await send('Network.enable'); await send('Runtime.enable'); await send('Accessibility.enable'); await send('DOM.enable');
let sess = ''; const post = async p => fetch(ORIGIN + p, {method: 'POST', headers: {Origin: ORIGIN, ...(sess && {Cookie: '__Host-kl_session=' + sess})}});
async function as(who) { const r = await post('/fixture/as/' + who); sess = /__Host-kl_session=([^;]+)/.exec(r.headers.get('set-cookie'))[1]; await send('Network.setCookie', {name: '__Host-kl_session', value: sess, url: ORIGIN + '/', secure: true, httpOnly: true, sameSite: 'Strict'}); await send('Network.deleteCookies', {name: '__Host-kl_receipt', url: ORIGIN + '/'}); }
const URLQ = (a, i) => `${ORIGIN}/home/receipts?agent=${a}&id=${encodeURIComponent(i)}`;
const WIDTHS = [[1440, 900], [390, 844], [320, 640]];

try {
  if (stage === 'layout' || stage === 'after') {
    const plan = stage === 'layout' ? [
      ['c', 'friends-c', '/home/friends', true], ['a', 'friends-a', '/home/friends', true], ['b', 'friends-b', '/home/friends', false], ['a', 'account-a', '/home/account', true],
      ['a', 'agents-a', '/home/agents', true, `document.querySelectorAll('details').forEach(e=>e.open=true)`], ['a', 'legacy-a', '/home/legacy', true, `document.querySelectorAll('details').forEach(e=>e.open=true)`],
      ['a', 'receipts-a', '/home/receipts', true], ['a', 'receipts-refusal-a', URLQ('agent_a1', NONE).replace(ORIGIN, ''), true], ['a', 'receipts-invalid-a', URLQ('agent_a2', 'x'.repeat(60) + '!bad').replace(ORIGIN, ''), false],
      ['a', 'new-self-a', '/home/new?peer=self', true], ['a', 'new-peer-a', '/home/new?peer=' + AB, true], ['a', 'friend-ab-a', '/home/friends/' + AB, true, `document.querySelectorAll('details').forEach(e=>e.open=true)`],
      ['b', 'friend-ab-b', '/home/friends/' + AB, true, `document.querySelectorAll('details').forEach(e=>e.open=true)`], ['b', 'new-b', '/home/new', false],
      ['c', 'agents-c', '/home/agents', true], ['c', 'new-c', '/home/new', true],
    ] : [
      ['a', 'friends-a-after', '/home/friends', true], ['c', 'friends-c-after', '/home/friends', true], ['a', 'friend-ac-a-after', '/home/friends/' + CA, true], ['c', 'friend-ac-c-after', '/home/friends/' + CA, true],
      ['a', 'new-a-after', '/home/new', true], ['a', 'account-a-after', '/home/account', true], ['c', 'new-c-after', '/home/new', true],
    ];
    let cur = '';
    for (const [who, name, p, cap, prep] of plan) {
      if (cur !== who) { await as(who); cur = who; }
      await metrics(1440, 900); await go(ORIGIN + p); if (prep) await js(prep);
      if (await js(`document.title`) === '') { check('page-title-' + name, false, 'empty'); }
      for (const [w, h] of WIDTHS) {
        await metrics(w, h); await sleep(120);
        const l = await layout();
        if (l.selects.length) (obs.selects ||= {})[`${w}-${name}`] = l.selects; // closed <select> text width vs room (DEV judgement item 2); observation only
        if (l.inputScroll.length) (obs.inputScroll ||= {})[`${w}-${name}`] = l.inputScroll; // value scrolls inside the field (DEV: existing behavior); observation only
        check(`layout-${w}-${name}`, l.iw === w && l.over <= 0 && l.clipped.length === 0 && l.hidden.length === 0 && l.ell.length === 0 && l.nameOver.length === 0 && l.small.length === 0 && !l.taClip, l);
        if (cap) await shot(`${slug(name)}-${w}`);
      }
    }
    await metrics(1440, 900);
    check('requests-same-origin-get-only', requests.every(q => q.origin === ORIGIN && q.method === 'GET'), requests.filter(q => q.origin !== ORIGIN || q.method !== 'GET'));
  }

  if (stage === 'keys') {
    await metrics(1440, 900);
    // ---- A account: friend code is reachable, selectable, copyable, read-only
    await as('a'); await go(ORIGIN + '/home/account');
    let st = await tabStops();
    check('account-tab-walk-names-focus-no-trap', st.at(-1).wrap === true && st.filter(s => !s.wrap).every(s => s.axName && focusOk(s)), st.filter(s => !s.wrap && !(s.axName && focusOk(s))));
    obs.accountTabOrder = st.filter(s => !s.wrap).map(s => `${s.tag}#${s.id || s.text}`);
    check('skip-link-is-first-and-targets-main', st[0].cls === 'skip' && st[0].fv, st[0]);
    await js('document.activeElement.blur()'); await js('window.scrollTo(0,0)');
    for (let i = 0; i < 40; i++) { await tab(); if (await js(`document.activeElement.id==='my-code'`)) break; }
    check('friend-code-reachable-by-tab', await js(`document.activeElement.id==='my-code'`));
    const codeInfo = await js(`(()=>{const e=document.activeElement;return {ro:e.readOnly,label:e.labels[0].textContent,v:e.value,fv:e.matches(':focus-visible'),outline:getComputedStyle(e).outlineStyle}})()`);
    check('friend-code-is-readonly-with-visible-label-and-focus', codeInfo.ro && codeInfo.label === '내 친구 코드' && codeInfo.fv && codeInfo.outline !== 'none', codeInfo);
    await selectAll();
    const sel = await js(`(()=>{const e=document.activeElement;return {s:e.selectionStart,e:e.selectionEnd,len:e.value.length}})()`);
    check('ctrl-a-selects-whole-code', sel.s === 0 && sel.e === sel.len && sel.len === 43, sel);
    await send('Input.insertText', {text: 'zzz'}); check('typing-cannot-modify-readonly-code', await js(`document.activeElement.value.length===43&&!document.activeElement.value.includes('zzz')`));
    await send('Browser.grantPermissions', {permissions: ['clipboardReadWrite', 'clipboardSanitizedWrite'], origin: ORIGIN}).catch(() => {});
    await key('c', 'KeyC', 67, 2); await sleep(200);
    const clip = await js(`navigator.clipboard.readText().then(t=>({len:t.length,same:t===document.getElementById('my-code').value})).catch(e=>({err:String(e).slice(0,80)}))`);
    obs.clipboard = clip; check('ctrl-c-copies-whole-code', clip.same === true, clip);
    // ---- A friends: tab walk, accept/deny names, form reachability
    await go(ORIGIN + '/home/friends');
    st = await tabStops();
    check('friends-a-tab-walk-names-focus-no-trap', st.at(-1).wrap === true && st.filter(s => !s.wrap).every(s => s.axName && focusOk(s)), st.filter(s => !s.wrap && !(s.axName && focusOk(s))));
    const names = st.filter(s => !s.wrap).map(s => s.axName); obs.friendsATabNames = names;
    check('friends-a-repeated-action-names-carry-target-and-id', names.some(n => n === `친구 요청 수락: 합성 회원 c / ${CA}`) && names.some(n => n === `친구 요청 거절: 합성 회원 c / ${CA}`), names);
    check('friends-a-code-field-before-request-form', names.indexOf('내 친구 코드') >= 0 && names.indexOf('내 친구 코드') < names.indexOf('친구가 준 친구 코드'), names);
    // ---- back / forward does not mutate
    const before = await js(`document.querySelector('main').innerText`); await go(ORIGIN + '/home/account'); await js('history.back()'); await sleep(900);
    check('history-back-restores-friends-without-change', await js(`location.pathname==='/home/friends'&&document.querySelector('main').innerText`) === before);
    // ---- C: request flow by keyboard only
    await as('c'); await go(ORIGIN + '/home/friends');
    st = await tabStops();
    check('friends-c-tab-walk-names-focus-no-trap', st.at(-1).wrap === true && st.filter(s => !s.wrap).every(s => s.axName && focusOk(s)), st.filter(s => !s.wrap && !(s.axName && focusOk(s))));
    await js('document.activeElement.blur(); window.scrollTo(0,0)');
    for (let i = 0; i < 40; i++) { await tab(); if (await js(`document.activeElement.id==='friend-code'`)) break; }
    check('c-code-input-reachable', await js(`document.activeElement.id==='friend-code'`));
    // empty submit: native validation stops it, no POST leaves the browser
    const postsBefore = requests.filter(q => q.method === 'POST').length;
    await tab(); await tab(); // consent checkbox, then submit button
    check('c-tab-order-code-consent-button', await js(`document.activeElement.textContent==='친구 요청 보내기'`));
    await key('Enter', 'Enter', 13); await sleep(500);
    check('c-empty-submit-blocked-by-native-validation', requests.filter(q => q.method === 'POST').length === postsBefore && await js(`location.pathname==='/home/friends'&&document.getElementById('friend-code').matches(':invalid')`));
    // type a non-existent code, consent by Space, submit by Enter
    await js(`document.getElementById('friend-code').focus()`); await send('Input.insertText', {text: ' nonexistent-code-' + 'x'.repeat(20) + ' '}); await tab(); await space();
    check('c-consent-checked-by-space', await js(`document.getElementById('friend-consent').checked`));
    await tab(); await enterNav();
    const res = await js(`({p:location.pathname,alert:!!document.querySelector('[role=alert]'),status:[...document.querySelectorAll('[role=status]')].map(e=>e.textContent.trim().slice(0,160)),active:document.activeElement.tagName+'#'+document.activeElement.id,html:document.querySelector('main').innerText.slice(0,300)})`);
    obs.cRequestResult = res; check('c-request-result-announced-as-status-not-alert', res.alert === false && res.status.length >= 1 && res.status.some(t => t.includes('요청을 보냈습니다')), res);
    check('c-request-result-does-not-reveal-success-friend', !/친구가 되었|친구입니다|성공/.test(res.html), res.html);
    await shot('c-request-result-1440', false); await metrics(390, 844); await shot('c-request-result-390', false); await metrics(1440, 900);
    // ---- receipts radio selector via keyboard
    await as('a'); await go(ORIGIN + '/home/receipts');
    st = await tabStops(20);
    obs.receiptsTabOrder = st.filter(s => !s.wrap).map(s => `${s.tag}${s.type ? ':' + s.type : ''}#${s.id || s.text}`);
    check('receipts-tab-walk-names-focus', st.filter(s => !s.wrap).every(s => s.axName && focusOk(s)), st.filter(s => !s.wrap && !(s.axName && focusOk(s))));
    await go(ORIGIN + '/home/receipts'); await js('document.activeElement && document.activeElement.blur()');
    for (let i = 0; i < 30; i++) { await tab(); if (await js(`document.activeElement.type==='radio'`)) break; }
    check('receipts-tab-enters-radio-group', await js(`document.activeElement.type==='radio'`));
    await key('ArrowDown', 'ArrowDown', 40); const rv = await js(`document.querySelector('input[name=agent]:checked')?.value`);
    check('receipts-arrow-selects-next-agent-radio', rv === 'agent_a2' || rv === 'agent_a1', rv);
    await tab(); await send('Input.insertText', {text: NONE}); await enterNav();
    const rr = await js(`({q:location.search,alert:document.querySelector('[role=alert]')?.textContent.trim().slice(0,90),checked:document.querySelector('input[name=agent]:checked')?.value,id:document.getElementById('receipt-id').value})`);
    check('receipts-submit-keeps-radio-and-id-and-announces-refusal', rr.q.includes('agent=agent_a') && rr.checked === rv && rr.id === NONE && /확인할 수 없습니다/.test(rr.alert || ''), rr);
    // ---- accessibility tree: landmarks, headings, fieldset name
    const ax = (await send('Accessibility.getFullAXTree')).nodes.filter(n => !n.ignored).map(n => ({role: n.role?.value, name: n.name?.value}));
    const has = (role, re) => ax.some(n => n.role === role && re.test(n.name || ''));
    check('ax-receipts-landmarks-heading-radiogroup', has('banner', /.*/) && has('main', /.*/) && has('heading', /^요청 상태 조회$/) && has('radiogroup', /^내 에이전트$/) && has('radio', /긴한글.*Node 22 로컬 클라이언트 \/ agent_a1/) && has('textbox', /^요청 ID 또는 답장 ID/), ax.filter(n => ['radiogroup', 'radio', 'textbox', 'heading', 'main', 'banner'].includes(n.role)));
    await go(ORIGIN + '/home/friends'); const ax2 = (await send('Accessibility.getFullAXTree')).nodes.filter(n => !n.ignored).map(n => ({role: n.role?.value, name: n.name?.value}));
    check('ax-friends-landmarks-code-textbox', ax2.some(n => n.role === 'textbox' && n.name === '내 친구 코드') && ax2.some(n => n.role === 'heading' && n.name === '친구') && ax2.some(n => n.role === 'navigation' && n.name === '주요 메뉴') && ax2.some(n => n.role === 'link' && n.name === '친구'), ax2.filter(n => ['textbox', 'navigation', 'main'].includes(n.role)));
    // ---- contrast (rendered colors) on friends and account
    const lum = c => { const [r, g, b] = c.match(/\d+(\.\d+)?/g).slice(0, 3).map(Number).map(v => { v /= 255; return v <= .03928 ? v / 12.92 : ((v + .055) / 1.055) ** 2.4; }); return .2126 * r + .7152 * g + .0722 * b; };
    const cr = (a, b) => { const [x, y] = [lum(a), lum(b)].sort((p, q) => q - p); return (x + .05) / (y + .05); };
    const col = await js(`(()=>{const g=(s,p)=>getComputedStyle(document.querySelector(s))[p];return {bg:g('body','backgroundColor'),text:g('main p','color'),muted:g('#my-code-help','color'),codeBg:g('#my-code','backgroundColor'),codeText:g('#my-code','color'),codeBorder:g('#my-code','borderTopColor'),inputBorder:g('#friend-code','borderTopColor'),btnBg:g('.primary','backgroundColor'),btnText:g('.primary','color'),link:g('main a','color')}})()`);
    const ratios = {text: cr(col.text, col.bg), muted: cr(col.muted, col.bg), link: cr(col.link, col.bg), button: cr(col.btnText, col.btnBg), code: cr(col.codeText, col.codeBg), codeBorder: cr(col.codeBorder, col.bg), inputBorder: cr(col.inputBorder, col.bg)};
    obs.contrast = ratios; check('contrast-text-4.5-nontext-3', ratios.text >= 4.5 && ratios.muted >= 4.5 && ratios.link >= 4.5 && ratios.button >= 4.5 && ratios.code >= 4.5 && ratios.inputBorder >= 3 && ratios.codeBorder >= 3, ratios);
    // ---- skip link actually moves to main
    await go(ORIGIN + '/home/friends'); await js('document.activeElement.blur()'); await tab(); await key('Enter', 'Enter', 13); await sleep(200);
    check('skip-link-enter-targets-main', await js(`location.hash==='#main'`)); await tab();
    check('after-skip-next-tab-is-in-main', await js(`!!document.activeElement.closest('main')`));
    // ---- 200%-like check is NOT a real browser zoom: reported not_run in the report
  }

  if (stage === 'legend') {
    // f799250 D-NM-02 recheck: GET only. Shared radio legends at 3 widths: no horizontal overflow, legend wraps inside the viewport, every radio label stays visible.
    const NAME40 = '가'.repeat(40);
    const legends = `[...document.querySelectorAll('fieldset[role=radiogroup] > legend')].map(e=>{const r=e.getBoundingClientRect(),cs=getComputedStyle(e),lh=parseFloat(cs.lineHeight)||parseFloat(cs.fontSize)*1.4;return {id:e.parentElement.id,text:e.textContent,left:Math.round(r.left),right:Math.round(r.right),h:Math.round(r.height),lines:Math.round(r.height/lh),sw:e.scrollWidth,cw:e.clientWidth,ow:cs.overflowWrap,mw:cs.maxWidth}})`;
    const labels = `[...document.querySelectorAll('fieldset[role=radiogroup] label')].map(e=>{const r=e.getBoundingClientRect(),sp=e.querySelector('span'),sr=sp.getBoundingClientRect();return {t:sp.textContent,left:Math.round(r.left),right:Math.round(r.right),h:Math.round(r.height),spanRight:Math.round(sr.right),sw:sp.scrollWidth,cw:sp.clientWidth}})`;
    const plan = [['a', 'new-self-a', '/home/new?peer=self'], ['a', 'new-peer-a', '/home/new?peer=' + AB], ['b', 'new-peer-b', '/home/new?peer=' + AB], ['a', 'receipts-a', '/home/receipts'], ['a', 'receipts-refusal-a', `/home/receipts?agent=agent_a1&id=${encodeURIComponent(NONE)}`]];
    let cur = '';
    for (const [who, name, p] of plan) {
      if (cur !== who) { await as(who); cur = who; }
      await metrics(1440, 900); await go(ORIGIN + p);
      for (const [w, h] of WIDTHS) {
        await metrics(w, h); await sleep(150);
        const l = await layout(), lg = await js(legends), lb = await js(labels);
        (obs.legends ||= {})[`${w}-${name}`] = {legends: lg, labels: lb};
        check(`layout-${w}-${name}`, l.iw === w && l.over <= 0 && l.clipped.length === 0 && l.hidden.length === 0 && l.ell.length === 0 && l.nameOver.length === 0 && l.small.length === 0 && !l.taClip, l);
        check(`legend-inside-viewport-and-wrapping-${w}-${name}`, lg.length > 0 && lg.every(g => g.left >= 0 && g.right <= w && g.sw <= g.cw + 1 && g.ow === 'anywhere'), lg);
        check(`radio-labels-visible-${w}-${name}`, lb.every(x => x.left >= 0 && x.right <= w && x.spanRight <= w && x.sw <= x.cw + 1 && x.h >= 43.5), lb);
        await shot(`legend-${slug(name)}-${w}`);
      }
    }
    // D-NM-02 exact path: B -> A, A's name is 40 chars without spaces; the legend names it in full and wraps below 1440
    await as('b'); await go(ORIGIN + '/home/new?peer=' + AB);
    const ph = await js(`document.querySelector('fieldset#to > legend').textContent`);
    check('d-nm-02-path-legend-names-the-40-char-peer-in-full', ph === NAME40 + '의 수신 에이전트', ph);
    for (const [w, h] of WIDTHS.slice(1)) { await metrics(w, h); await sleep(150); const g = (await js(legends)).find(x => x.id === 'to'); check(`d-nm-02-legend-wraps-to-several-lines-${w}`, g.lines >= 2 && g.right <= w, g); }
    await metrics(1440, 900);
    check('requests-same-origin-get-only', requests.every(q => q.origin === ORIGIN && q.method === 'GET'), requests.filter(q => q.origin !== ORIGIN || q.method !== 'GET'));
  }

  if (stage === 'recheck') {
    // 5fa3333 recheck: only the screens the fix touched (GET), keyboard on the new-conversation radios, link navigation, one refused rename POST (changes nothing).
    const FORBID = '이름에 줄바꿈·제어 문자·방향 제어 문자는 쓸 수 없습니다. 보이는 글자로 1–40자를 입력하세요.';
    const rad = g => `[...document.querySelectorAll('input[name=${g}]')].map(e=>({v:e.value,c:e.checked,lab:e.closest('label').innerText.trim()}))`;
    const checkedOf = async g => (await js(rad(g))).filter(r => r.c).map(r => r.v);
    const openAll = `document.querySelectorAll('details').forEach(e=>e.open=true)`;
    const plan = [['a', 'new-self-a', '/home/new?peer=self'], ['a', 'new-peer-a', '/home/new?peer=' + AB], ['b', 'new-peer-b', '/home/new?peer=' + AB], ['a', 'agents-a', '/home/agents', openAll], ['a', 'friends-a', '/home/friends']];
    let cur = '';
    for (const [who, name, p, prep] of plan) {
      if (cur !== who) { await as(who); cur = who; }
      await metrics(1440, 900); await go(ORIGIN + p); if (prep) await js(prep);
      for (const [w, h] of WIDTHS) {
        await metrics(w, h); await sleep(120);
        const l = await layout();
        check(`layout-${w}-${name}`, l.iw === w && l.over <= 0 && l.clipped.length === 0 && l.hidden.length === 0 && l.ell.length === 0 && l.nameOver.length === 0 && l.small.length === 0 && !l.taClip, l);
        await shot(`${slug(name)}-${w}`);
      }
    }

    // ---- (1) new conversation: defaults, labels, AX tree, keyboard
    await as('a'); await metrics(390, 844); await go(ORIGIN + '/home/new?peer=self');
    obs.newSelfRadios = {from: await js(rad('from')), to: await js(rad('to'))};
    check('new-self-no-select-two-radiogroups', await js(`!document.querySelector('main select') && document.querySelectorAll('fieldset[role=radiogroup]').length===2`));
    check('new-self-defaults-from-a1-to-a2', JSON.stringify(await checkedOf('from')) === '["agent_a1"]' && JSON.stringify(await checkedOf('to')) === '["agent_a2"]', obs.newSelfRadios);
    let ax = (await send('Accessibility.getFullAXTree')).nodes.filter(n => !n.ignored).map(n => ({role: n.role?.value, name: n.name?.value}));
    obs.newSelfAx = ax.filter(n => ['radiogroup', 'radio', 'button'].includes(n.role));
    const grp = ax.filter(n => n.role === 'radiogroup').map(n => n.name), rn = ax.filter(n => n.role === 'radio').map(n => n.name);
    check('new-self-ax-groups-named-by-legend-radios-by-label', grp.includes('내 발신 에이전트') && grp.includes('내 수신 에이전트') && rn.length === 4 && rn.every(n => /^.+ \/ Node 22 로컬 클라이언트 \/ agent_a[12]$/.test(n)) && new Set(rn).size === 2, {grp, rn});
    await js('document.activeElement && document.activeElement.blur(); window.scrollTo(0,0)');
    for (let i = 0; i < 30; i++) { await tab(); if (await js(`document.activeElement.name==='from'`)) break; }
    check('new-self-tab-enters-from-group-on-checked-radio', await js(`document.activeElement.name==='from'&&document.activeElement.value==='agent_a1'&&document.activeElement.checked&&document.activeElement.matches(':focus-visible')`));
    await key('ArrowDown', 'ArrowDown', 40);
    check('new-self-arrow-down-selects-next-from', JSON.stringify(await checkedOf('from')) === '["agent_a2"]', await checkedOf('from'));
    for (let i = 0; i < 6; i++) { await tab(); if (await js(`document.activeElement.name==='to'`)) break; }
    check('new-self-tab-enters-to-group-on-checked-radio', await js(`document.activeElement.name==='to'&&document.activeElement.value==='agent_a2'&&document.activeElement.checked`));
    await key('ArrowUp', 'ArrowUp', 38);
    check('new-self-arrow-up-selects-previous-to', JSON.stringify(await checkedOf('to')) === '["agent_a1"]', await checkedOf('to'));
    await tab(); check('new-self-tab-reaches-open-button', await js(`document.activeElement.textContent.trim()==='대화 열기'`));
    await enterNav();
    const u1 = new URL(await js('location.href'));
    obs.newSelfOpen = {path: u1.pathname, query: u1.search}; // the fixture answers 403 for /home/conversation; the request itself is the evidence
    check('new-self-enter-requests-conversation-with-chosen-from-and-to', u1.pathname === '/home/conversation' && u1.searchParams.get('from') === 'agent_a2' && u1.searchParams.get('to') === 'agent_a1' && [...u1.searchParams.keys()].length === 2, obs.newSelfOpen);

    await go(ORIGIN + '/home/new?peer=' + AB);
    obs.newPeerRadios = {from: await js(rad('from')), to: await js(rad('to'))};
    check('new-peer-defaults-from-a1-to-b1', JSON.stringify(await checkedOf('from')) === '["agent_a1"]' && JSON.stringify(await checkedOf('to')) === '["agent_b1"]' && (await js(rad('to'))).length === 1, obs.newPeerRadios);
    ax = (await send('Accessibility.getFullAXTree')).nodes.filter(n => !n.ignored).map(n => ({role: n.role?.value, name: n.name?.value}));
    check('new-peer-ax-legend-names-peer', ax.some(n => n.role === 'radiogroup' && n.name === '합성 회원 b의 수신 에이전트') && ax.some(n => n.role === 'radio' && /^😀{40} \/ Node 22 로컬 클라이언트 \/ agent_b1$/u.test(n.name || '')), ax.filter(n => n.role === 'radiogroup' || n.role === 'radio'));
    await js('document.activeElement && document.activeElement.blur(); window.scrollTo(0,0)');
    for (let i = 0; i < 30; i++) { await tab(); if (await js(`document.activeElement.name==='from'`)) break; }
    await key('ArrowDown', 'ArrowDown', 40);
    for (let i = 0; i < 6; i++) { await tab(); if (await js(`document.activeElement.name==='to'`)) break; }
    check('new-peer-tab-enters-single-to-radio', await js(`document.activeElement.name==='to'&&document.activeElement.value==='agent_b1'`));
    await tab(); await enterNav();
    const u2 = new URL(await js('location.href'));
    check('new-peer-enter-requests-conversation-with-chosen-from-and-peer-agent', u2.pathname === '/home/conversation' && u2.searchParams.get('from') === 'agent_a2' && u2.searchParams.get('to') === 'agent_b1' && [...u2.searchParams.keys()].length === 2, u2.search);

    // ---- (2) my agents: the two links sit on separate lines, each reaches the right place
    for (const [w, h] of WIDTHS) {
      await metrics(w, h); await go(ORIGIN + '/home/agents'); await js(openAll); await sleep(120);
      const pairs = await js(`(()=>{const r=[...document.querySelectorAll('main a[href^="/home/receipts?agent="]')],l=[...document.querySelectorAll('main a[href="/home/legacy#agents"]')];return r.map((a,i)=>{const x=a.getBoundingClientRect(),y=l[i].getBoundingClientRect();return {agent:a.getAttribute('href'),sameParent:a.parentElement===l[i].parentElement,gap:Math.round(y.top-x.bottom),h1:Math.round(x.height),h2:Math.round(y.height),lineBreak:y.top>=x.bottom-1}})})()`);
      (obs.agentLinks ||= {})[w] = pairs;
      check(`agents-links-separate-lines-and-44px-${w}`, pairs.length === 2 && pairs.every(q => !q.sameParent && q.lineBreak && q.h1 >= 43.5 && q.h2 >= 43.5), pairs);
    }
    await metrics(390, 844);
    for (const ag of ['agent_a1', 'agent_a2']) {
      await go(ORIGIN + '/home/agents'); await js(openAll);
      await js(`document.querySelector('main a[href="/home/receipts?agent=${ag}#receipt-form"]').focus()`); await enterNav();
      const r = await js(`({p:location.pathname,q:location.search,h:location.hash,checked:document.querySelector('input[name=agent]:checked')?.value})`);
      check(`agents-receipt-link-opens-lookup-with-${ag}-selected`, r.p === '/home/receipts' && r.q === '?agent=' + ag && r.h === '#receipt-form' && r.checked === ag, r);
    }
    await go(ORIGIN + '/home/agents'); await js(openAll);
    await js(`document.querySelector('main a[href="/home/legacy#agents"]').focus()`); await enterNav();
    const lg = await js(`({p:location.pathname,h:location.hash,target:!!document.getElementById('agents'),t:document.title})`);
    check('agents-manage-link-opens-connection-page-at-agents', lg.p === '/home/legacy' && lg.h === '#agents' && lg.target, lg);

    // ---- (4) received request: refusal notice sits between accept and deny and describes the deny button
    await go(ORIGIN + '/home/friends');
    ax = (await send('Accessibility.getFullAXTree')).nodes.filter(n => !n.ignored);
    const deny = ax.find(n => n.role?.value === 'button' && (n.name?.value || '').startsWith('친구 요청 거절:')), accept = ax.find(n => n.role?.value === 'button' && (n.name?.value || '').startsWith('친구 요청 수락:'));
    obs.denyAx = {name: deny?.name?.value, description: deny?.description?.value, acceptDescription: accept?.description?.value ?? null};
    check('deny-button-accessible-description-is-the-refusal-notice', !!deny && deny.description?.value === '거절하면 요청한 사람의 기록에 거절로 표시됩니다. 알리지 않으려면 응답하지 않아도 기한이 지나면 요청이 끝납니다.' && !accept.description?.value, obs.denyAx);
    check('deny-notice-order-accept-notice-deny', await js(`(()=>{const f=document.querySelector('form[action="/home/friends/decision"]');const k=[...f.children].filter(e=>e.tagName!=='INPUT').map(e=>e.tagName+':'+(e.value||e.id||''));return JSON.stringify(k)})()`) .then(v => /^\["BUTTON:accept","P:deny-help-.+","BUTTON:deny"\]$/.test(v) ? true : v) === true);
    const st = await tabStops();
    const names = st.filter(s => !s.wrap).map(s => s.axName); obs.friendsTab = names;
    const ia = names.findIndex(n => (n || '').startsWith('친구 요청 수락:'));
    check('friends-tab-accept-then-deny-notice-not-a-stop', ia >= 0 && (names[ia + 1] || '').startsWith('친구 요청 거절:') && st.at(-1).wrap === true && st.filter(s => !s.wrap).every(s => s.axName && focusOk(s)), names);

    // ---- (3) refused rename in the browser: dedicated text, forbidden characters removed from the field, nothing saved
    await metrics(1440, 900); await go(ORIGIN + '/home/agents'); await js(openAll);
    const orig = await js(`document.getElementById('agent-name-agent_a2').value`);
    await js(`document.getElementById('agent-name-agent_a2').focus()`); await selectAll();
    await send('Input.insertText', {text: 'ab\u0007‮cd'});
    obs.typedBeforeSubmit = await js(`[...document.getElementById('agent-name-agent_a2').value].map(c=>c.codePointAt(0).toString(16))`);
    await enterNav();
    const rf = await js(`(()=>{const i=document.getElementById('agent-name-agent_a2');return {p:location.pathname,alert:[...document.querySelectorAll('[role=alert]')].map(e=>e.textContent.trim()),val:[...i.value].map(c=>c.codePointAt(0).toString(16)),inv:i.getAttribute('aria-invalid'),desc:i.getAttribute('aria-describedby'),title:document.title,text:document.body.innerText}})()`);
    obs.refusal = rf;
    check('rename-refusal-browser-dedicated-text-and-stripped-value', rf.p === '/home/agent-name' && rf.alert.some(t => t.includes(FORBID)) && rf.val.join(',') === '61,62,63,64' && rf.inv === 'true' && !/[ --‪-‮⁦-⁩]/.test(rf.text), {...rf, text: rf.text.length});
    for (const [w, h] of WIDTHS) {
      await metrics(w, h); await sleep(120); const l = await layout();
      check(`layout-${w}-rename-refusal-a`, l.iw === w && l.over <= 0 && l.clipped.length === 0 && l.hidden.length === 0 && l.nameOver.length === 0 && l.small.length === 0, l);
      await shot(`rename-refusal-a-${w}`);
    }
    await metrics(1440, 900); await go(ORIGIN + '/home/agents'); await js(openAll);
    check('rename-refusal-browser-saved-nothing', await js(`document.getElementById('agent-name-agent_a2').value`) === orig && orig.length === 39, orig);
    const posts = requests.filter(q => q.method === 'POST');
    check('requests-same-origin-and-exactly-one-post-the-refused-rename', requests.every(q => q.origin === ORIGIN) && posts.length === 1 && posts[0].path === '/home/agent-name', posts);
  }
} finally { save(); ws.close(); }
const bad = checks.filter(c => !c.passed); console.log(`checks ${checks.length} failed ${bad.length}`); process.exit(bad.length ? 1 : 0);
