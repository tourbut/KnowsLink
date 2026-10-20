// Independent browser checks for the connect guide / resume UX (U8/U7/O-NM-03): real key events (CDP Input), 3-width layout matrix, focus walk, a11y names, captures.
// Usage: node browser.mjs <cdp-port> <label> <stage>   stage = layout | keys | after.  CAP=1 saves PNG captures. Synthetic A/B/C only; layout/keys are GET-only (+ /fixture/as).
// after = post-mutate.py read-only re-check (back/forward/reload do not create agents).
import fs from 'node:fs';
import path from 'node:path';
import {fileURLToPath} from 'node:url';
import {createHash} from 'node:crypto';

const out = path.dirname(fileURLToPath(import.meta.url));
const [port, label, stage] = [process.argv[2], process.argv[3], process.argv[4] || 'layout'];
const ORIGIN = 'http://127.0.0.1:18932';
const M = JSON.parse(fs.readFileSync(process.env.CONNECT_MANIFEST || 'C:/Users/shin/AppData/Local/KnowsLinkDevTools/connect-ux002-dev/fixture-manifest.json', 'utf8'));
const CN = M.connections, AG = M.agents;
const CAP = process.env.CAP === '1';
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
const WIDTHS_SAFE = null;
const WIDTHS = [[1440, 900], [390, 844], [320, 640]];
// click/Enter navigations bypass go(); wait out the shared per-member 429 window (40/min over /home/*) and reload to the same URL
const enterNavRL = async () => { await enterNav(); for (let i = 0; i < 12 && await rateLimited(); i++) { obs.memberWindow429Waits = (obs.memberWindow429Waits || 0) + 1; await sleep(8000); const w = new Promise(r => loadWaiters.push(r)); await js('location.reload()'); await Promise.race([w, sleep(8000)]); await sleep(150); } };
const shotIf = (n, full = true) => CAP ? shot(n, full) : null;
const openAll = `document.querySelectorAll('details').forEach(e=>e.open=true)`;
const axTree = async () => (await send('Accessibility.getFullAXTree')).nodes.filter(n => !n.ignored).map(n => ({role: n.role?.value, name: n.name?.value}));
const cardCount = () => js(`document.querySelectorAll('li[id^=agent-]').length`);
const pageFacts = () => js(`({title:document.title,h1:document.querySelectorAll('h1').length,main:!!document.querySelector('main'),skip:!!document.querySelector('a.skip[href="#main"]'),lang:document.documentElement.lang,cur:[...document.querySelectorAll('[aria-current]')].map(e=>e.getAttribute('aria-current')+':'+e.textContent.trim().slice(0,24))})`);
const A_CONN = ['waiting', 'prepared', 'approved', 'consumed'], B_CONN = ['expired', 'cancelled', 'device_lost'];

try {
  if (stage === 'layout') {
    const plan = [
      ['a', 'agents-a', '/home/agents'], ['a', 'agents-a-claude', '/home/agents?use=claude&start=first'], ['a', 'agents-a-other', '/home/agents?use=other&start=first'], ['a', 'agents-a-dots-continue', '/home/agents?use=dots&start=continue'],
      ...A_CONN.map(k => ['a', 'conn-' + k, '/home/connections/' + CN[k]]),
      ['a', 'friends-a', '/home/friends'], ['a', 'legacy-a', '/home/legacy', openAll],
      ['b', 'agents-b', '/home/agents'], ...B_CONN.map(k => ['b', 'conn-' + k, '/home/connections/' + CN[k]]), ['b', 'friends-b', '/home/friends'],
      ['c', 'agents-c', '/home/agents'], ['c', 'friends-c', '/home/friends'],
    ];
    let cur = '';
    for (const [who, name, p, prep] of plan) {
      if (cur !== who) { await as(who); cur = who; }
      await metrics(1440, 900); await go(ORIGIN + p); if (prep) await js(prep);
      const f = await pageFacts();
      check('page-facts-' + name, f.title.length > 0 && f.h1 === 1 && f.main && f.skip && f.lang === 'ko', f);
      if (name.startsWith('conn-')) { const ok = B_CONN.some(k => name === 'conn-' + k) ? f.cur.length === 0 : f.cur.length === 1 && f.cur[0].startsWith('step:'); check('aria-current-step-only-when-state-known-' + name, ok, f.cur); }
      for (const [w, h] of WIDTHS) {
        await metrics(w, h); await sleep(120);
        const l = await layout();
        check(`layout-${w}-${name}`, l.iw === w && l.over <= 0 && l.clipped.length === 0 && l.hidden.length === 0 && l.ell.length === 0 && l.nameOver.length === 0 && l.small.length === 0, l);
        await shotIf(`${slug(name)}-${w}`);
      }
    }
    await metrics(1440, 900);
    check('requests-same-origin-get-only', requests.every(q => q.origin === ORIGIN && q.method === 'GET'), requests.filter(q => q.origin !== ORIGIN || q.method !== 'GET'));
  }

  if (stage === 'keys') {
    await metrics(1440, 900); await as('a');
    // ---- selector by keyboard only: Tab into radio group, arrows pick, Tab to start group and button, Enter submits GET
    await go(ORIGIN + '/home/agents'); await js('document.activeElement.blur(); window.scrollTo(0,0)');
    const n0 = await cardCount();
    for (let i = 0; i < 30; i++) { await tab(); if (await js(`document.activeElement.name==='use'`)) break; }
    check('tab-reaches-use-radio-group', await js(`document.activeElement.name==='use'`));
    check('tab-enters-group-at-first-choice-when-none-selected', await js(`document.activeElement.value==='web'`), await js('document.activeElement.value'));
    for (let i = 0; i < 3; i++) await key('ArrowDown', 'ArrowDown', 40);
    check('arrow-down-x3-selects-claude', await js(`document.querySelector('input[name=use]:checked')?.value`) === 'claude', await js(`document.querySelector('input[name=use]:checked')?.value`));
    await key('ArrowDown', 'ArrowDown', 40); await key('ArrowDown', 'ArrowDown', 40); // other, then wrap -> web
    check('arrow-wraps-inside-group-no-escape', await js(`document.querySelector('input[name=use]:checked')?.value`) === 'web', await js(`document.querySelector('input[name=use]:checked')?.value`));
    await key('ArrowUp', 'ArrowUp', 38); check('arrow-up-wraps-to-other', await js(`document.querySelector('input[name=use]:checked')?.value`) === 'other');
    await key('ArrowUp', 'ArrowUp', 38); // claude
    await tab(); check('tab-leaves-use-group-to-start-group', await js(`document.activeElement.name==='start'`), await js('document.activeElement.name'));
    await key('ArrowDown', 'ArrowDown', 40); check('start-arrow-selects-continue', await js(`document.querySelector('input[name=start]:checked')?.value`) === 'continue');
    await key('ArrowUp', 'ArrowUp', 38);
    await tab(); check('tab-reaches-submit-button', await js(`document.activeElement.tagName==='BUTTON'&&document.activeElement.textContent.trim()==='안내 보기'`));
    await enterNavRL();
    const g = await js(`({q:location.search,h:location.hash,h3:document.querySelector('#use-guide h3')?.textContent,active:document.activeElement.id||document.activeElement.tagName,checked:document.querySelector('input[name=use]:checked')?.value,start:document.querySelector('input[name=start]:checked')?.value,cards:document.querySelectorAll('li[id^=agent-]').length})`);
    obs.focusAfterSelect = g.active;
    check('enter-submits-get-guide-shown-for-claude-first', g.q === '?use=claude&start=first' && g.h === '#use-guide' && /^Claude Code 작업 세션/.test(g.h3) && g.checked === 'claude' && g.start === 'first', g);
    check('selection-by-keyboard-created-no-agent', g.cards === n0, [g.cards, n0]);
    check('guide-region-scrolled-into-view', await js(`(()=>{const r=document.getElementById('use-guide').getBoundingClientRect();return r.top>=-1&&r.top<innerHeight})()`));
    await js('location.reload()'); await sleep(900);
    check('reload-guide-keeps-selection-and-card-count', await js(`document.querySelector('input[name=use]:checked')?.value`) === 'claude' && await cardCount() === n0, await cardCount());
    // ---- tab walk of guide page: names, focus, no trap, skip first
    let st = await tabStops(90);
    check('agents-tab-walk-names-focus-no-trap', st.at(-1).wrap === true && st.filter(s => !s.wrap).every(s => s.axName && focusOk(s)), st.filter(s => !s.wrap && !(s.axName && focusOk(s))));
    check('skip-link-first-stop-visible-focus', st[0].cls === 'skip' && st[0].fv, st[0]);
    obs.agentsTabNames = st.filter(s => !s.wrap).map(s => s.axName);
    const small = st.filter(s => !s.wrap && s.h > 0 && s.h < 43.5); check('agents-tab-stops-all-44px', small.length === 0, small.map(s => s.key + ':' + s.h));
    // ---- other -> web link
    await go(ORIGIN + '/home/agents?use=other&start=first#use-guide');
    await js(`document.querySelector('#use-guide a').focus()`); await enterNavRL();
    check('other-guide-link-goes-to-web-guide', await js(`location.search==='?use=web'&&!!document.querySelector('#use-guide h3')&&/^웹에서 바로 사용/.test(document.querySelector('#use-guide h3').textContent)`), await js('location.href'));
    // ---- connection pages: tab walk and exact return
    for (const k of A_CONN) {
      await go(ORIGIN + '/home/connections/' + CN[k]); st = await tabStops(30);
      check(`conn-${k}-tab-walk-names-focus-no-trap-44px`, st.at(-1).wrap === true && st.filter(s => !s.wrap).every(s => s.axName && focusOk(s) && s.h >= 43.5), st.filter(s => !s.wrap && !(s.axName && focusOk(s) && s.h >= 43.5)));
      obs['connTabNames-' + k] = st.filter(s => !s.wrap).map(s => s.axName);
    }
    for (const k of A_CONN) {
      const ag = AG[k === 'consumed' ? 'registered' : k];
      await go(ORIGIN + '/home/connections/' + CN[k]);
      await js(`[...document.querySelectorAll('a')].find(a=>/내 에이전트로 돌아가기/.test(a.textContent)).focus()`); await enterNavRL();
      const r = await js(`({p:location.pathname,h:location.hash,el:!!document.getElementById('agent-${ag}'),top:Math.round(document.getElementById('agent-${ag}')?.getBoundingClientRect().top)})`);
      check(`conn-${k}-return-link-lands-on-own-agent-card`, r.p === '/home/agents' && r.h === `#agent-${ag}` && r.el && r.top >= -2 && r.top < 400, r);
    }
    // resume card link -> connection page (waiting), then browser back keeps agents unchanged
    await go(ORIGIN + '/home/agents'); const n1 = await cardCount();
    await js(`[...document.querySelectorAll('a')].find(a=>a.textContent.trim()==='연결 상태 확인').focus()`); await enterNavRL();
    check('resume-link-opens-connection-page', await js(`location.pathname.startsWith('/home/connections/')&&document.querySelector('[aria-current=step]')!==null`), await js('location.pathname'));
    await js('history.back()'); await sleep(900); const n2 = await cardCount();
    await js('history.forward()'); await sleep(900); await js('history.back()'); await sleep(900);
    check('history-back-forward-do-not-change-agents', n1 === n2 && await cardCount() === n1 && await js(`location.pathname==='/home/agents'`), [n1, n2, await cardCount()]);
    // ---- friends A: accessible names of cancel buttons carry the typed tail
    await go(ORIGIN + '/home/friends'); st = await tabStops(80);
    check('friends-a-tab-walk-names-focus-no-trap-44px', st.at(-1).wrap === true && st.filter(s => !s.wrap).every(s => s.axName && focusOk(s) && s.h >= 43.5), st.filter(s => !s.wrap && !(s.axName && focusOk(s) && s.h >= 43.5)));
    const names = st.filter(s => !s.wrap).map(s => s.axName); obs.friendsATabNames = names;
    const cancels = names.filter(n => n && n.includes('친구 요청 취소'));
    check('friends-a-cancel-ax-names-carry-tail-or-marker', ['…zz99', '형식 다름', '이전 요청', ...Object.values(M.rerequest_codes).map(c => '…' + c.slice(-4))].every(t => cancels.some(n => n.includes(t))) && cancels.filter(n => /입력한 코드 …/.test(n)).length >= 5, cancels);
    check('friends-a-no-two-cancel-names-identical', new Set(cancels).size === cancels.length, cancels);
    // ---- AX of agents page: landmarks, radiogroups named by legend, headings
    await go(ORIGIN + '/home/agents?use=web&start=first'); const ax2 = await axTree();
    check('ax-agents-landmarks-radiogroups-headings', ax2.some(n => n.role === 'main') && ax2.some(n => n.role === 'banner') && ax2.some(n => n.role === 'radiogroup' && n.name === '사용할 곳') && ax2.some(n => n.role === 'radiogroup' && n.name === '연결 시작') && ax2.some(n => n.role === 'heading' && n.name === '내 에이전트와 연결 상태') && ax2.some(n => n.role === 'radio' && /^웹에서 바로 사용/.test(n.name || '')), ax2.filter(n => ['radiogroup', 'heading', 'radio'].includes(n.role)).map(n => n.role + ':' + n.name));
    // ---- B and C: tab walks
    await as('b'); await go(ORIGIN + '/home/agents'); st = await tabStops(140);
    check('agents-b-tab-walk-names-focus-no-trap', st.at(-1).wrap === true && st.filter(s => !s.wrap).every(s => s.axName && focusOk(s)), st.filter(s => !s.wrap && !(s.axName && focusOk(s))));
    await as('c'); await go(ORIGIN + '/home/agents'); st = await tabStops(60);
    check('agents-c-tab-walk-names-focus-no-trap', st.at(-1).wrap === true && st.filter(s => !s.wrap).every(s => s.axName && focusOk(s)), st.filter(s => !s.wrap && !(s.axName && focusOk(s))));
    check('requests-same-origin-get-only', requests.every(q => q.origin === ORIGIN && q.method === 'GET'), requests.filter(q => q.origin !== ORIGIN || q.method !== 'GET'));
  }

  if (stage === 'after') {
    // after mutate.py: A has 5 agents; GET-only browser checks of the new waiting agent, history and reload do not create or resubmit
    const MU = JSON.parse(fs.readFileSync(path.join(out, 'mutate-result.json'), 'utf8')); const loc = MU.location;
    await metrics(1440, 900); await as('a'); await go(ORIGIN + '/home/agents');
    check('after-agents-a-five-cards', await cardCount() === 5, await cardCount());
    // creation returns to the new agent's own card (design: 생성은 새 agent 위치로 복귀) -> /home/agents#agent-{id}
    check('after-create-location-is-agents-hash-of-new-card', loc === `/home/agents#agent-${MU.new_agent}`, loc);
    await go(ORIGIN + loc); await sleep(300);
    check('after-new-card-in-view-and-says-unused-no-key', await js(`(()=>{const e=document.getElementById('agent-${MU.new_agent}');const r=e.getBoundingClientRect();return r.top>=-2&&r.top<innerHeight&&/연결 상태: 미사용 · 등록된 키 없음/.test(e.innerText)&&!/키 등록됨/.test(e.innerText)&&!/온라인/.test(e.innerText)})()`), await js(`document.getElementById('agent-${MU.new_agent}').innerText.slice(0,200)`));
    for (const [w, h] of WIDTHS) { await metrics(w, h); await sleep(120); const l = await layout(); check(`layout-${w}-agents-a-after`, l.iw === w && l.over <= 0 && l.clipped.length === 0 && l.hidden.length === 0 && l.small.length === 0, l); await shotIf(`agents-a-after-${w}`); }
    await metrics(1440, 900);
    await js(`[...document.querySelectorAll('a')].find(a=>a.textContent.trim()==='연결 상태 확인' && a.closest('li').id==='agent-${MU.new_agent}')?.focus()`);
    const hasStatusLink = await js(`document.activeElement.closest('li')?.id==='agent-${MU.new_agent}'`);
    if (hasStatusLink) { await enterNavRL(); const f = await pageFacts(); check('after-new-agent-resume-link-opens-connection-page-step-1', f.cur.length === 1 && f.cur[0].startsWith('step:1') && f.h1 === 1 && f.main && f.skip && f.title.length > 0, f);
      for (const [w, h] of WIDTHS) { await metrics(w, h); await sleep(120); const l = await layout(); check(`layout-${w}-conn-new-a`, l.iw === w && l.over <= 0 && l.clipped.length === 0 && l.hidden.length === 0 && l.small.length === 0, l); await shotIf(`conn-new-a-${w}`); }
      await metrics(1440, 900); } else { obs.newCardHasNoStatusLink = true; }
    await go(ORIGIN + loc); await js('history.back()'); await sleep(900); await js('history.forward()'); await sleep(900); await js('location.reload()'); await sleep(1200);
    await go(ORIGIN + '/home/agents');
    check('after-back-forward-reload-keep-five-cards', await cardCount() === 5, await cardCount());
    check('after-no-post-from-browser', requests.every(q => q.origin === ORIGIN && q.method === 'GET'), requests.filter(q => q.method !== 'GET'));
  }
} finally { save(); ws.close(); }
const bad = checks.filter(c => !c.passed); console.log(`checks ${checks.length} failed ${bad.length}`); process.exit(bad.length ? 1 : 0);
