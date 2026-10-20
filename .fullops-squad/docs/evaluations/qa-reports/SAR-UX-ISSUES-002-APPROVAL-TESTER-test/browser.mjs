// Independent browser checks for the approval summary UX (U5/U7): real key events (CDP Input), 3-width layout matrix, focus walk, a11y names, captures.
// Usage: node browser.mjs <cdp-port> <label> <stage>   stage = layout | keys | flow | after
// Synthetic A/B/C only. layout/keys are GET-only. flow does ONE real deny on pending_decide_2 (last step); after only reads.
import fs from 'node:fs';
import path from 'node:path';
import {fileURLToPath} from 'node:url';
import {createHash} from 'node:crypto';

const out = path.dirname(fileURLToPath(import.meta.url));
const [port, label, stage] = [process.argv[2], process.argv[3], process.argv[4] || 'layout'];
const ORIGIN = 'http://127.0.0.1:18931';
const M = JSON.parse(fs.readFileSync('C:/Users/shin/AppData/Local/KnowsLinkDevTools/approval-ux002-dev/fixture-manifest.json', 'utf8'));
const G = M.gates;
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
const WIDTHS = [[1440, 900], [390, 844], [320, 640]];

const CAP = label.startsWith('chrome'); // designer captures come from Chrome only
const shotIf = (n) => CAP ? shot(n) : null;
const openAll = `document.querySelectorAll('details').forEach(e=>e.open=true)`;
// text/pre containers that could clip: dd, code, message-text and pre must neither overflow the page nor hide text
const textLayout = () => js(`(()=>{
  const vis = e => e.getClientRects().length && getComputedStyle(e).visibility !== 'hidden';
  const t = [...document.querySelectorAll('main dd, main code, main .message-text, main pre, main p, main li, main h1, main h2, main h3')].filter(vis);
  const over = t.filter(e => e.getBoundingClientRect().right > innerWidth + 1 || e.getBoundingClientRect().left < -1).map(e => e.tagName + '.' + e.className);
  const clip = t.filter(e => ['hidden','clip'].includes(getComputedStyle(e).overflowX) && e.scrollWidth > e.clientWidth + 1).map(e => e.tagName + '.' + e.className);
  const pre = [...document.querySelectorAll('main pre')].filter(vis).map(e => ({sw: e.scrollWidth, cw: e.clientWidth, scrolls: e.scrollWidth > e.clientWidth + 1, tabindex: e.getAttribute('tabindex'), ox: getComputedStyle(e).overflowX, ws: getComputedStyle(e).whiteSpace, wrap: getComputedStyle(e).overflowWrap}));
  return {over, clip, pre};
})()`);
const WIDTHS2 = [[1440, 900], [390, 844], [320, 640]];
const gatePath = k => '/home/gates/' + G[k];

try {
  await send('Page.enable'); await send('Network.enable'); await send('Runtime.enable'); await send('Accessibility.enable'); await send('DOM.enable');
  if (stage === 'layout' || stage === 'after') {
    const plan = stage === 'layout' ? [
      ['c', 'list-c-empty', '/home/legacy#gates', true],
      ['a', 'list-a', '/home/legacy#gates', true],
      ['a', 'home-a', '/home', true],
      ['a', 'detail-query-hint', gatePath('pending_query_hint'), true],
      ['a', 'detail-query-hint-open', gatePath('pending_query_hint'), true, openAll],
      ['a', 'detail-commit-long', gatePath('pending_commit_long'), true],
      ['a', 'detail-commit-long-open', gatePath('pending_commit_long'), true, openAll],
      ['a', 'detail-decide-1', gatePath('pending_decide_1'), true],
      ['a', 'detail-decide-2', gatePath('pending_decide_2'), false],
      ['a', 'detail-expired', gatePath('expired'), true],
      ['a', 'detail-processed', gatePath('processed'), true],
      ['a', 'detail-original-missing', gatePath('original_missing'), true],
      ['b', 'detail-foreign-b', gatePath('pending_decide_1'), true],
      ['a', 'friends-a', '/home/friends', true],
    ] : [
      ['a', 'after-list-a', '/home/legacy#gates', true],
      ['a', 'after-detail-decide-1-approved', gatePath('pending_decide_1'), true],
      ['a', 'after-detail-decide-2-denied', gatePath('pending_decide_2'), true],
      ['a', 'after-home-a', '/home', true],
    ];
    let cur = '';
    for (const [who, name, p, cap, prep] of plan) {
      if (cur !== who) { await as(who); cur = who; }
      await metrics(1440, 900); await go(ORIGIN + p); if (prep) await js(prep);
      check('title-h1-main-' + name, (await js('document.title')).length > 0 && await js(`document.querySelectorAll('h1').length===1 && document.querySelectorAll('main').length===1`));
      for (const [w, h] of WIDTHS2) {
        await metrics(w, h); await sleep(120);
        const l = await layout(), tl = await textLayout();
        if (tl.pre.some(x => x.scrolls)) (obs.preScroll ||= {})[`${w}-${name}`] = tl.pre;
        check(`layout-${w}-${name}`, l.iw === w && l.over <= 0 && l.clipped.length === 0 && l.hidden.length === 0 && l.ell.length === 0 && l.small.length === 0 && tl.over.length === 0 && tl.clip.length === 0, {...l, tl});
        if (cap) await shotIf(`${slug(name)}-${w}`);
      }
    }
    await metrics(1440, 900);
    // reflow equivalent of browser 200% at 1280 CSS px (a 640px viewport); NOT a real browser zoom
    if (stage === 'layout') {
      await as('a');
      for (const [name, p, prep] of [['list-a', '/home/legacy#gates'], ['detail-commit-long-open', gatePath('pending_commit_long'), openAll], ['detail-query-hint-open', gatePath('pending_query_hint'), openAll]]) {
        await metrics(1440, 900); await go(ORIGIN + p); if (prep) await js(prep);
        await metrics(640, 450); await sleep(120); const l = await layout(), tl = await textLayout();
        check(`reflow-640-equiv-200pct-${name}`, l.iw === 640 && l.over <= 0 && l.clipped.length === 0 && l.hidden.length === 0 && tl.over.length === 0 && tl.clip.length === 0, {...l, tl});
      }
      await metrics(1440, 900);
    }
    check('requests-same-origin-get-only', requests.every(q => q.origin === ORIGIN && q.method === 'GET'), requests.filter(q => q.origin !== ORIGIN || q.method !== 'GET'));
  }

  if (stage === 'keys') {
    await metrics(1440, 900);
    await as('a');
    // ---- detail page (pending, hint): tab walk, names, focus, 44px, no trap, skip link
    await go(ORIGIN + gatePath('pending_query_hint'));
    let st = await tabStops();
    obs.detailTabNames = st.filter(s => !s.wrap).map(s => s.axName); obs.detailTabOrder = st.filter(s => !s.wrap).map(s => `${s.tag}:${s.text}`);
    check('detail-tab-first-stop-is-skip-link', st[0].cls === 'skip' && st[0].fv, st[0]);
    check('detail-tab-order-header-nav-hint-approve-deny-technical-return', JSON.stringify(st.filter(s => !s.wrap).map(s => s.text)) === JSON.stringify(['본문으로 건너뛰기', 'KnowsLink', '내 계정', '승인 목록', '참고 문구 — 승인 근거 아님', '이 요청 승인', '거절', '기술 상세', '승인 목록으로 돌아가기']), obs.detailTabOrder);
    check('detail-tab-walk-names-visible-focus-44px-no-trap', st.at(-1).wrap === true && st.filter(s => !s.wrap).every(s => s.axName && focusOk(s) && s.h >= 43.5), st.filter(s => !s.wrap && !(s.axName && focusOk(s) && s.h >= 43.5)));
    // Space/Enter on summaries toggle without moving focus
    await js('window.scrollTo(0,0)');
    await js(`document.querySelector('details summary').focus()`);
    await space(); const o1 = await js(`document.querySelector('details').open`); const f1 = await js(`document.activeElement.tagName`);
    await space(); const o2 = await js(`document.querySelector('details').open`);
    await key('Enter', 'Enter', 13); const o3 = await js(`document.querySelector('details').open`);
    check('hint-details-toggle-by-space-and-enter-focus-stays', o1 === true && o2 === false && o3 === true && f1 === 'SUMMARY', {o1, o2, o3, f1});
    await js(`document.querySelector('button[value=approve]').focus()`); await shiftTab();
    check('shift-tab-goes-back-from-approve-to-hint-summary', await js(`document.activeElement.tagName==='SUMMARY'&&document.activeElement.textContent.includes('참고 문구')`));
    const ax = (await send('Accessibility.getFullAXTree')).nodes;
    const roles = ax.filter(n => !n.ignored).map(n => [n.role?.value, n.name?.value]);
    obs.landmarks = roles.filter(r => ['main', 'navigation', 'heading', 'banner', 'alert', 'status', 'form'].includes(r[0]));
    check('ax-main-one-h1-nav-named-heading-order', roles.filter(r => r[0] === 'main').length === 1 && roles.filter(r => r[0] === 'heading' && r[1] === '요청 승인').length === 1 && roles.some(r => r[0] === 'navigation' && r[1] === '현재 위치') && JSON.stringify(roles.filter(r => r[0] === 'heading').map(r => r[1]).slice(0, 5)) === JSON.stringify(['요청 승인', '누가 누구에게', '요청한 내용', '승인하면 허용되는 범위', '기한·상태']), obs.landmarks);
    check('ax-current-page-marked-only-on-request-approval', await js(`[...document.querySelectorAll('[aria-current]')].map(e=>e.getAttribute('aria-current')+':'+e.textContent.trim()).join('|')`) === 'page:요청 승인');
    await go(ORIGIN + gatePath('pending_query_hint')); await tab(); await key('Enter', 'Enter', 13); await sleep(200);
    check('skip-link-moves-to-main', await js(`location.hash==='#main'`), await js('location.href'));

    // ---- original_missing: alert, only safe return
    await go(ORIGIN + gatePath('original_missing'));
    const om = await js(`({alerts:[...document.querySelectorAll('[role=alert]')].map(e=>e.textContent.trim()),statuses:[...document.querySelectorAll('[role=status]')].length,buttons:document.querySelectorAll('main button,main form,main input').length,links:[...document.querySelectorAll('main a')].map(e=>e.getAttribute('href')+'|'+e.textContent.trim())})`);
    check('original-missing-alert-and-no-controls-and-safe-return', om.alerts.length === 1 && om.alerts[0].includes('승인할 내용을 확인할 수 없습니다') && om.buttons === 0 && om.links.includes('/home/legacy#gates|승인 목록으로 돌아가기'), om);
    st = await tabStops(); check('original-missing-tab-walk-no-trap', st.at(-1).wrap === true && st.filter(s => !s.wrap).every(s => s.axName && focusOk(s) && s.h >= 43.5), st);
    for (const k of ['expired', 'processed']) {
      await go(ORIGIN + gatePath(k));
      const r = await js(`({status:[...document.querySelectorAll('[role=status]')].map(e=>e.textContent.trim()),ctl:document.querySelectorAll('main button,main form,main input').length})`);
      check(`${k}-status-note-and-no-controls`, r.status.length === 1 && r.status[0].includes('승인하거나 거절할 수 없습니다') && r.ctl === 0, r);
    }

    // ---- list: repeated link names, Enter opens detail, return link comes back to the list section
    await go(ORIGIN + '/home/legacy#gates');
    st = await tabStops(80);
    const reviews = st.filter(s => (s.text || '') === '요청 검토');
    obs.listReviewNames = reviews.map(s => s.axName);
    check('list-review-links-7-names-with-type-requester-expiry', reviews.length === 7 && reviews.every(s => /^요청 검토: .+ \/ .+ \/ 기한 \d{4}-/.test(s.axName)), obs.listReviewNames);
    obs.listReviewNamesDistinct = new Set(reviews.map(s => s.axName)).size;
    check('list-tab-walk-visible-focus-44px-no-trap', st.at(-1).wrap === true && st.filter(s => !s.wrap).every(s => s.axName && focusOk(s) && s.h >= 43.5), st.filter(s => !s.wrap && !(s.axName && focusOk(s) && s.h >= 43.5)));
    await js(`document.querySelector('#gates a').focus()`); await enterNav();
    check('enter-on-review-opens-detail-with-title-and-h1', await js(`location.pathname.startsWith('/home/gates/')&&document.title==='요청 승인 · KnowsLink'&&document.querySelector('h1').textContent==='요청 승인'`), await js('location.href'));
    await js(`[...document.querySelectorAll('main a')].find(a=>a.textContent.includes('승인 목록으로 돌아가기')).focus()`); await enterNav();
    check('return-link-goes-to-approval-list-section', await js(`location.pathname==='/home/legacy'&&location.hash==='#gates'&&document.querySelector(':target')?.id==='gates'`), await js('location.href'));
    await go(ORIGIN + gatePath('pending_decide_1')); await js(`document.querySelector('nav[aria-label="현재 위치"] a').focus()`); await enterNav();
    check('breadcrumb-link-goes-to-approval-list-section', await js(`location.pathname==='/home/legacy'&&location.hash==='#gates'`), await js('location.href'));

    // ---- foreign member: alert + return
    await as('b'); await go(ORIGIN + gatePath('pending_decide_1'));
    st = await tabStops(); check('foreign-403-tab-walk-and-alert', st.at(-1).wrap === true && await js(`document.querySelectorAll('[role=alert]').length===1&&document.title==='요청 승인 · KnowsLink'`) && st.filter(s => !s.wrap).every(s => s.axName && focusOk(s) && s.h >= 43.5), st);

    // ---- O-NM-08 friends: accept -> deny, deny description = help text
    await as('a'); await go(ORIGIN + '/home/friends');
    st = await tabStops(); const fn = st.filter(s => !s.wrap).map(s => s.axName);
    const ia = fn.findIndex(n => (n || '').startsWith('친구 요청 수락:'));
    check('onm08-tab-accept-then-deny-help-not-a-stop', ia >= 0 && (fn[ia + 1] || '').startsWith('친구 요청 거절:'), fn);
    const dd = await js(`(()=>{const b=document.querySelector('button[value=deny]');const h=document.getElementById(b.getAttribute('aria-describedby'));return {help:h.textContent.trim(),isPrev:b.previousElementSibling===h,mt:getComputedStyle(h).marginTop,mb:getComputedStyle(h).marginBottom,acceptBottom:document.querySelector('button[value=accept]').getBoundingClientRect().bottom,helpTop:h.getBoundingClientRect().top,helpBottom:h.getBoundingClientRect().bottom,denyTop:b.getBoundingClientRect().top}})()`);
    obs.onm08 = dd;
    check('onm08-help-between-accept-and-deny-with-spacing', dd.isPrev && dd.helpTop > dd.acceptBottom && dd.denyTop >= dd.helpBottom && parseFloat(dd.mt) >= 12, dd);
    const axd = (await send('Accessibility.getFullAXTree')).nodes.find(n => n.role?.value === 'button' && (n.name?.value || '').startsWith('친구 요청 거절'));
    check('onm08-deny-ax-description-equals-help-text', axd && (axd.description?.value || '').includes('거절하면 요청한 사람의 기록에 거절로 표시됩니다'), axd && axd.description);
    check('requests-same-origin-get-only', requests.every(q => q.origin === ORIGIN && q.method === 'GET'), requests.filter(q => q.origin !== ORIGIN || q.method !== 'GET'));
  }

  if (stage === 'flow') {
    // ONE real deny on pending_decide_2 with real keys, then back-button behaviour
    await metrics(1440, 900); await as('a'); await go(ORIGIN + '/home/legacy#gates');
    const before = await js(`document.querySelector('#gates').textContent.replace(/\\s+/g,' ')`);
    await go(ORIGIN + gatePath('pending_decide_2'));
    check('flow-decide2-pending-with-buttons-before', await js(`!!document.querySelector('button[value=deny]')&&!!document.querySelector('button[value=approve]')`));
    await js(`document.querySelector('button[value=deny]').focus()`);
    const w = new Promise(r => loadWaiters.push(r)); await key('Enter', 'Enter', 13); await Promise.race([w, sleep(3000)]); await sleep(400);
    const r1 = await js(`({path:location.pathname,title:document.title,h1:document.querySelector('h1')?.textContent,badge:document.querySelector('#state-h')?.parentElement.textContent.replace(/\\s+/g,' '),ctl:document.querySelectorAll('main button,main form').length,status:[...document.querySelectorAll('[role=status]')].map(e=>e.textContent.trim()),alert:[...document.querySelectorAll('[role=alert]')].map(e=>e.textContent.trim()),text:document.body.innerText.slice(0,300)})`);
    obs.afterDeny = r1;
    const posts = requests.filter(q => q.method === 'POST' && !q.path.startsWith('/fixture/'));
    check('flow-deny-posts-once-to-the-deny-url-and-lands-on-detail', posts.length === 1 && posts[0].path === `/home/gates/${G.pending_decide_2}/deny` && r1.path === gatePath('pending_decide_2'), {posts, r1});
    check('flow-after-deny-shows-denied-no-controls-and-status-note', /거절/.test(r1.badge) && r1.ctl === 0 && r1.status.length === 1, r1);
    await shotIf('flow-after-deny-1440');
    const hist = await send('Page.getNavigationHistory');
    obs.history = hist.entries.map(e => e.url.replace(ORIGIN, ''));
    const prev = hist.entries[hist.currentIndex - 1];
    const w2 = new Promise(r => loadWaiters.push(r)); await send('Page.navigateToHistoryEntry', {entryId: prev.id}); await Promise.race([w2, sleep(3000)]); await sleep(400);
    const r2 = await js(`({path:location.pathname,hasForm:!!document.querySelector('main form'),approve:!!document.querySelector('button[value=approve]'),state:document.querySelector('#state-h')?.parentElement.textContent.replace(/\\s+/g,' ')})`);
    obs.afterBack = r2;
    check('flow-back-after-deny-shows-fresh-closed-state-not-stale-form', r2.path === gatePath('pending_decide_2') && r2.hasForm === false && r2.approve === false && /거절/.test(r2.state || ''), r2);
    await go(ORIGIN + '/home/legacy#gates');
    const after = await js(`document.querySelector('#gates').textContent.replace(/\\s+/g,' ')`);
    check('flow-list-shows-denied-row-and-differs-from-before', after !== before && /거절/.test(after), {before: before.length, after: after.length});
    // note: the page CSP (default-src 'none') blocks fetch(), so the pending count after deny is read over HTTP in mutate.py 2 (the first run died here after all checks were saved)
  }

  if (stage === 'onm04') {
    // O-NM-04: fixture blocks /home/receipts (403), so the anonymous-429 page is rendered from the pinned SHA's own template (go test, KNOWSLINK_UI_FIXTURE_DIR) and opened from disk. Template render, not a live 429.
    const url = 'file:///' + path.join(out, 'render', 'receipt-anonymous-limited.html').split(path.sep).join('/');
    await metrics(1440, 900); await send('Page.navigate', {url}); await sleep(600);
    const r = await js(`({title:document.title,h1:document.querySelectorAll('h1').length,alerts:[...document.querySelectorAll('[role=alert]')].map(e=>e.textContent.trim()),radios:document.querySelectorAll('input,select,textarea,fieldset,legend,form').length,emptyHref:document.querySelectorAll('a[href=""]').length,links:[...document.querySelectorAll('main a')].map(e=>e.getAttribute('href')+'|'+e.textContent.trim()),restore:document.body.innerText.includes('복원'),retry:/\\d{4}-\\d\\d-\\d\\d \\d\\d:\\d\\d:\\d\\d KST 이후 다시 시도/.test(document.body.innerText)})`);
    obs.render = r;
    check('onm04-no-selector-form-or-legend', r.radios === 0, r);
    check('onm04-retry-time-and-real-login-link-no-empty-href', r.retry && r.links.includes('/|로그인 화면으로 이동') && r.emptyHref === 0, r);
    check('onm04-no-restore-note-for-anonymous', r.restore === false, r);
    check('onm04-title-one-h1-one-alert', r.title === '요청 상태 조회 · KnowsLink' && r.h1 === 1 && r.alerts.length === 1, r);
    const st = await tabStops(); obs.onm04Tab = st.map(s => s.axName || (s.wrap ? 'wrap' : s.text));
    check('onm04-tab-walk-names-focus-44px', st.at(-1).wrap === true && st.filter(s => !s.wrap).every(s => s.axName && focusOk(s) && s.h >= 43.5), st);
    for (const [w, h] of WIDTHS2) { await metrics(w, h); await sleep(120); const l = await layout(), tl = await textLayout(); check(`layout-${w}-onm04-anonymous-429`, l.iw === w && l.over <= 0 && l.clipped.length === 0 && l.small.length === 0 && tl.over.length === 0, {...l, tl}); await shotIf(`onm04-anonymous-429-${w}`); }
    await send('Page.navigate', {url: url.replace('receipt-anonymous-limited', 'receipt-limited')}); await sleep(600);
    obs.memberLimited = await js(`({restore:document.body.innerText.includes('복원'),emptyHref:document.querySelectorAll('a[href=""]').length,link:[...document.querySelectorAll('[role=alert] a')].map(e=>e.getAttribute('href')+'|'+e.textContent)})`);
  }

  if (stage === 'contrast') {
    // WCAG 2.x contrast of every visible text node vs its effective background (solid colours only) on the approval pages
    const calc = `(()=>{const parse=c=>{const m=c.match(/rgba?\\(([^)]+)\\)/);if(!m)return null;const p=m[1].split(/[ ,\\/]+/).map(Number);return {r:p[0],g:p[1],b:p[2],a:p[3]===undefined?1:p[3]}};
      const lum=c=>{const f=v=>{v/=255;return v<=0.03928?v/12.92:Math.pow((v+0.055)/1.055,2.4)};return 0.2126*f(c.r)+0.7152*f(c.g)+0.0722*f(c.b)};
      const bg=e=>{for(;e;e=e.parentElement){const c=parse(getComputedStyle(e).backgroundColor);if(c&&c.a>0.99)return c}return {r:255,g:255,b:255,a:1}};
      const out=[];const w=document.createTreeWalker(document.body,NodeFilter.SHOW_TEXT);let n;
      while(n=w.nextNode()){if(!n.textContent.trim())continue;const e=n.parentElement;if(!e.getClientRects().length||['STYLE','SCRIPT'].includes(e.tagName))continue;const cs=getComputedStyle(e);const f=parse(cs.color),b=bg(e);const l1=lum(f),l2=lum(b);const r=(Math.max(l1,l2)+0.05)/(Math.min(l1,l2)+0.05);const px=parseFloat(cs.fontSize);const big=px>=24||(px>=18.66&&parseInt(cs.fontWeight)>=700);out.push({t:n.textContent.trim().slice(0,24),r:Math.round(r*100)/100,need:big?3:4.5,tag:e.tagName+'.'+e.className})}
      return out})()`;
    const pages = [['a', '/home/legacy#gates'], ['a', gatePath('pending_query_hint')], ['a', gatePath('pending_commit_long')], ['a', gatePath('expired')], ['a', gatePath('processed')], ['a', gatePath('original_missing')], ['b', gatePath('pending_decide_1')], ['a', '/home/friends']];
    let cur = '';
    for (const [who, p] of pages) {
      if (cur !== who) { await as(who); cur = who; }
      await go(ORIGIN + p); await js(openAll);
      const rows = await js(calc); const low = rows.filter(x => x.r < x.need);
      obs['contrast ' + p.slice(0, 22) + ' ' + who] = {texts: rows.length, min: Math.min(...rows.map(x => x.r))};
      check('contrast-text-4.5-' + slug(p) + '-' + who, rows.length > 0 && low.length === 0, low);
    }
  }

  if (stage === 'p409') {
    // the 409 body of a duplicate click (HTTP response of mutate.py 1) opened from disk for layout and capture
    const url = 'file:///' + path.join(out, 'render', 'gate-409-refused.html').split(path.sep).join('/');
    await metrics(1440, 900); await send('Page.navigate', {url}); await sleep(600);
    check('p409-title-h1-alert-return', await js(`document.title==='요청 승인 · KnowsLink'&&document.querySelectorAll('h1').length===1&&document.querySelectorAll('[role=alert]').length===1&&!!document.querySelector('a[href="/home/legacy#gates"]')`));
    for (const [w, h] of WIDTHS2) { await metrics(w, h); await sleep(120); const l = await layout(), tl = await textLayout(); check(`layout-${w}-gate-409-refused`, l.iw === w && l.over <= 0 && l.clipped.length === 0 && l.small.length === 0 && tl.over.length === 0, {...l, tl}); await shotIf(`gate-409-refused-${w}`); }
  }
} finally { save(); ws.close(); }
const bad = checks.filter(c => !c.passed); console.log(`checks ${checks.length} failed ${bad.length}`); process.exit(bad.length ? 1 : 0);
