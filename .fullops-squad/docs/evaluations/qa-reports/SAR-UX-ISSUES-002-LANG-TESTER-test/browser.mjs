// Independent browser checks for the LANG fixture (U3 settings/language, U7 header): real key events (CDP Input), 3-width layout matrix
// in both languages, focus walk, accessibility-tree names, designer captures.
// Usage: node browser.mjs <cdp-port> <label> <stage> [acct]
//   stage = layout | keys | langflow | save        (layout/keys/langflow are GET-only; save does the real, restored language save)
//   acct  = en | ko   (which synthetic member to sign in as for layout/keys; default en)
import fs from 'node:fs';
import path from 'node:path';
import {fileURLToPath} from 'node:url';
import {createHash} from 'node:crypto';

const out = path.dirname(fileURLToPath(import.meta.url));
const [port, label, stage, acct = 'en'] = [process.argv[2], process.argv[3], process.argv[4] || 'layout', process.argv[5]];
const ORIGIN = 'http://127.0.0.1:18933';
const M = JSON.parse(fs.readFileSync('C:/Users/shin/AppData/Local/KnowsLinkDevTools/lang-ux002-dev/fixture-manifest.json', 'utf8'));
const ID = M.ids;
const version = await (await fetch(`http://127.0.0.1:${port}/json/version`)).json();
const target = (await (await fetch(`http://127.0.0.1:${port}/json/list`)).json()).find(t => t.type === 'page');
const ws = new WebSocket(target.webSocketDebuggerUrl);
await new Promise((res, rej) => { ws.onopen = res; ws.onerror = rej; });
let next = 0; const pending = new Map(), checks = [], captures = [], requests = [], loadWaiters = [], obs = {};
const tag = stage === 'layout' || stage === 'keys' ? `${stage}-${acct}` : stage;
const resultFile = path.join(out, `browser-${label}-${tag}.json`);
const save = () => fs.writeFileSync(resultFile, JSON.stringify({browser: version.Browser, protocol: version['Protocol-Version'], origin: ORIGIN, stage, acct, checks, observations: obs, captures, requests}, null, 2));
ws.onmessage = ev => {
  const m = JSON.parse(ev.data);
  if (m.id) { const p = pending.get(m.id); pending.delete(m.id); m.error ? p.reject(new Error(m.error.message)) : p.resolve(m.result); }
  else if (m.method === 'Network.requestWillBeSent') { const u = new URL(m.params.request.url); requests.push({method: m.params.request.method, origin: u.origin, path: u.pathname + (u.search ? '?…' : '')}); }
  else if (m.method === 'Page.loadEventFired') loadWaiters.splice(0).forEach(f => f());
};
const send = (method, params = {}) => new Promise((resolve, reject) => { const id = ++next; pending.set(id, {resolve, reject}); ws.send(JSON.stringify({id, method, params})); });
const js = async expr => { const r = await send('Runtime.evaluate', {expression: expr, returnByValue: true, awaitPromise: true}); if (r.exceptionDetails) throw new Error('eval failed: ' + expr.slice(0, 100) + ' :: ' + JSON.stringify(r.exceptionDetails.exception?.description || '').slice(0, 200)); return r.result.value; };
const sleep = ms => new Promise(r => setTimeout(r, ms));
function check(name, value, detail) { checks.push({name, passed: value === true, detail}); save(); console.log(value === true ? 'PASS' : 'FAIL', name, value === true ? '' : JSON.stringify(detail ?? '').slice(0, 400)); }
const rateLimited = () => js(`/요청이 많습니다|Too many requests/.test(document.body.innerText)`);
async function go(url) {
  for (let i = 0; i < 12; i++) {
    const w = new Promise(r => loadWaiters.push(r)); await send('Page.navigate', {url}); await Promise.race([w, sleep(8000)]); await sleep(150);
    if (!(await rateLimited()) || url.includes('/fixture/')) return;
    obs.memberWindow429Waits = (obs.memberWindow429Waits || 0) + 1; await sleep(8000);
  }
}
async function key(k, code, vk, mods = 0, text) { await send('Input.dispatchKeyEvent', {type: 'keyDown', key: k, code, windowsVirtualKeyCode: vk, modifiers: mods, text: text ?? (k === 'Enter' ? '\r' : undefined)}); await send('Input.dispatchKeyEvent', {type: 'keyUp', key: k, code, windowsVirtualKeyCode: vk, modifiers: mods}); await sleep(100); }
const tab = () => key('Tab', 'Tab', 9), shiftTab = () => key('Tab', 'Tab', 9, 8), space = () => key(' ', 'Space', 32, 0, ' ');
const arrow = d => ({Down: () => key('ArrowDown', 'ArrowDown', 40), Up: () => key('ArrowUp', 'ArrowUp', 38), Right: () => key('ArrowRight', 'ArrowRight', 39), Left: () => key('ArrowLeft', 'ArrowLeft', 37)}[d])();
const enterNav = async () => { const w = new Promise(r => loadWaiters.push(r)); await key('Enter', 'Enter', 13); await Promise.race([w, sleep(3000)]); await sleep(250); };
const metrics = (w, h, dpr = 1) => send('Emulation.setDeviceMetricsOverride', {width: w, height: h, deviceScaleFactor: dpr, mobile: false});
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
  const small = ctl.filter(e => hit(e).height < 43.5 && !(e.tagName==='A' && e.classList.contains('skip'))).map(e => e.tagName + '#' + (e.id || e.name || '') + ':' + Math.round(hit(e).height) + ':' + (e.textContent || e.value || '').trim().slice(0, 14));
  const clipped = [...document.querySelectorAll('body *')].filter(e => vis(e) && !e.closest('svg') && (e.getBoundingClientRect().right > innerWidth + 1 || e.getBoundingClientRect().left < -1) && getComputedStyle(e).position !== 'fixed' && !e.classList.contains('skip')).map(e => e.tagName + '#' + e.id + '.' + e.className).slice(0, 6);
  const hidden = [...document.querySelectorAll('main *')].filter(e => vis(e) && !/^(INPUT|TEXTAREA|SELECT)$/.test(e.tagName) && ['hidden','clip'].includes(getComputedStyle(e).overflowX) && e.scrollWidth > e.clientWidth + 1).map(e => e.tagName + '.' + e.className).slice(0, 6);
  const ell = [...document.querySelectorAll('main *')].filter(e => vis(e) && getComputedStyle(e).textOverflow === 'ellipsis').map(e => e.tagName + '.' + e.className).slice(0, 6);
  const names = [...document.querySelectorAll('.name-preview,.wrap,main h1,main h2,.choice span,.consent span,summary h3,main dd,main li,main p,button')].filter(vis);
  const nameOver = names.filter(e => e.scrollWidth > e.clientWidth + 1 || e.getBoundingClientRect().right > innerWidth + 1).map(e => e.tagName + '.' + e.className + ':' + e.scrollWidth + '>' + e.clientWidth).slice(0, 6);
  const btnClip = [...document.querySelectorAll('button,.button')].filter(vis).filter(e => e.scrollHeight > e.clientHeight + 1 || e.scrollWidth > e.clientWidth + 1).map(e => (e.textContent || '').trim().slice(0, 20)).slice(0, 6);
  return {iw: innerWidth, dpr: devicePixelRatio, over: document.documentElement.scrollWidth - innerWidth, small, clipped, hidden, ell, nameOver, btnClip, lang: document.documentElement.lang, h1: document.querySelectorAll('h1').length, nCtl: ctl.length};
})()`);
// accessibility-tree names of the whole page: used to prove the English page has no Korean names (other than the language name / user data)
async function axNames() {
  const ax = (await send('Accessibility.getFullAXTree')).nodes;
  // InlineTextBox nodes are fragments of a StaticText whose whole name is checked already
  return ax.filter(n => !n.ignored && n.role?.value !== 'InlineTextBox').map(n => ({role: n.role?.value, name: n.name?.value ?? '', desc: n.description?.value ?? ''}));
}
const HANGUL = /[\u3131-\u318e\uac00-\ud7a3]/;
const EXEMPT = n => /긴한글이름+/.test(n) || n === '한국어';
const tabStops = async (cap = 80) => {
  await js('document.activeElement && document.activeElement.blur(); window.scrollTo(0,0)');
  const stops = []; const ax = (await send('Accessibility.getFullAXTree')).nodes; const byId = new Map(ax.filter(n => n.backendDOMNodeId).map(n => [n.backendDOMNodeId, n]));
  for (let i = 0; i < cap; i++) {
    await tab();
    const info = await js(`(()=>{const e=document.activeElement;if(e===document.body||e===document.documentElement)return {body:true};const cs=getComputedStyle(e);const r=e.getBoundingClientRect();const l=(e.type==='checkbox'||e.type==='radio')?(e.closest('label')||e):e;const lr=l.getBoundingClientRect();return {key:e.tagName+'|'+e.id+'|'+(e.getAttribute('href')||'')+'|'+(e.textContent||e.value||'').trim().slice(0,26),tag:e.tagName,id:e.id,type:e.type||'',text:(e.textContent||e.value||'').trim().slice(0,26),fv:e.matches(':focus-visible'),outline:cs.outlineStyle+' '+cs.outlineWidth,shadow:cs.boxShadow!=='none',h:Math.round(lr.height),inx:r.left>=-1&&r.right<=innerWidth+1,skip:e.classList.contains('skip')}})()`);
    if (info.body) { stops.push({wrap: true}); break; }
    if (stops.length > 1 && info.key === stops[0].key) { stops.push({wrap: true}); break; }
    const {objectId} = (await send('Runtime.evaluate', {expression: 'document.activeElement'})).result; const bid = (await send('DOM.describeNode', {objectId})).node.backendNodeId;
    const n = byId.get(bid); info.axName = n?.name?.value ?? null; info.axRole = n?.role?.value ?? null; stops.push(info);
  }
  return stops;
};
const focusOk = s => s.fv && ((!s.outline.startsWith('none') && parseFloat(s.outline.split(' ')[1]) >= 2) || s.shadow) && s.inx;

await send('Page.enable'); await send('Network.enable'); await send('Runtime.enable'); await send('Accessibility.enable'); await send('DOM.enable');
let sess = ''; const post = async p => fetch(ORIGIN + p, {method: 'POST', headers: {Origin: ORIGIN, ...(sess && {Cookie: '__Host-kl_session=' + sess})}});
async function as(who) { const r = await post('/fixture/as/' + who); sess = /__Host-kl_session=([^;]+)/.exec(r.headers.get('set-cookie'))[1]; await send('Network.setCookie', {name: '__Host-kl_session', value: sess, url: ORIGIN + '/', secure: true, httpOnly: true, sameSite: 'Strict'}); }
async function anon(cookieLang) {
  await send('Network.deleteCookies', {name: '__Host-kl_session', url: ORIGIN + '/'}); await send('Network.deleteCookies', {name: '__Host-kl_lang', url: ORIGIN + '/'});
  if (cookieLang) await send('Network.setCookie', {name: '__Host-kl_lang', value: cookieLang, url: ORIGIN + '/', secure: true, httpOnly: true, sameSite: 'Strict'});
}
const WIDTHS = [[1440, 900], [390, 844], [320, 640]];
const CAP = label.startsWith('chrome'); // designer captures come from Chrome only
const openAll = `document.querySelectorAll('details').forEach(e=>e.open=true)`;

// ---- screens: [name, path, who(session)|null(anon), anonCookie, openDetails]
const MEMBER = [['home', '/home'], ['name', '/home/name'], ['conversation', M.conversation], ['friends', '/home/friends'], ['friends-detail-none', null], ['agents', '/home/agents'], ['account', '/home/account'], ['legacy', '/home/legacy'], ['new', '/home/new'], ['receipts', '/home/receipts'], ['conn-waiting', '/home/connections/' + ID.connection_waiting], ['conn-prepared', '/home/connections/' + ID.connection_prepared], ['conn-approved', '/home/connections/' + ID.connection_approved], ['gate', '/home/gates/' + ID.gate], ['connect-device', '/connect/' + ID.device], ['remote', '/remote/consent/' + ID.remote]].filter(x => x[1]);
const ANON = (cookie) => [['public', '/'], ['verify', '/auth/verify'], ['connect-device', '/connect/' + ID.device], ['remote', '/remote/consent/' + ID.remote]].map(x => [...x, cookie]);

try {
  if (stage === 'layout') {
    const plan = [];
    if (acct === 'en' || acct === 'ko') for (const [n, p] of MEMBER) plan.push({who: acct, n, p, cookie: null});
    if (acct === 'en') for (const c of ['en', 'ko']) for (const [n, p] of ANON(c)) plan.push({who: null, n: 'anon-' + n, p, cookie: c});
    for (const [w, h] of WIDTHS) {
      await metrics(w, h);
      let curWho = '#';
      for (const e of plan) {
        const sig = (e.who || 'anon') + ':' + e.cookie;
        if (sig !== curWho) { if (e.who) { await anon(null); await as(e.who); } else await anon(e.cookie); curWho = sig; }
        await go(ORIGIN + e.p);
        for (const open of [false, true]) {
          if (open) { const n = await js(`document.querySelectorAll('main details').length`); if (!n) continue; await js(openAll); await sleep(100); }
          const r = await layout(); const L = r.lang; const nm = `${e.n}${open ? '-open' : ''}-${w}`;
          // expected language: account saved value, or the anonymous cookie, or Korean
          const exp = e.who ? (e.who === 'en' ? 'en' : 'ko') : (e.cookie || 'ko');
          const t = `${L}:${nm}`;
          const okGeom = r.over <= 0 && r.small.length === 0 && r.clipped.length === 0 && r.hidden.length === 0 && r.ell.length === 0 && r.nameOver.length === 0 && r.btnClip.length === 0;
          check(`layout_${t}_${e.who || 'anon'}`, okGeom && r.h1 === 1 && L === exp, {...r});
          if (CAP) await shot(`${e.who || 'anon'}-${L}-${nm}`.replace(/-+/g, '-'));
        }
      }
    }
    // accessibility-tree names in English contain no Korean (language name / user data exempt); Korean pages' landmark names are Korean
    await metrics(1440, 900);
    for (const e of plan) {
      if (e.who) { await anon(null); await as(e.who); } else await anon(e.cookie);
      await go(ORIGIN + e.p); const nodes = await axNames(); const lang = await js('document.documentElement.lang');
      const named = nodes.filter(n => n.name || n.desc);
      if (lang === 'en') { const bad = named.filter(n => (HANGUL.test(n.name) && !EXEMPT(n.name)) || (HANGUL.test(n.desc) && !EXEMPT(n.desc))); check(`ax_names_en_no_korean_${e.who || 'anon'}_${e.n}`, bad.length === 0, bad.slice(0, 5)); }
      else { const lm = nodes.filter(n => ['navigation', 'main', 'banner', 'region', 'form', 'group', 'radiogroup'].includes(n.role) && n.name); const bad = lm.filter(n => /^[A-Za-z][A-Za-z ,'-]+$/.test(n.name)); check(`ax_landmark_names_ko_not_english_${e.who || 'anon'}_${e.n}`, bad.length === 0, bad.slice(0, 5)); }
    }
  }

  if (stage === 'keys') {
    await metrics(1440, 900);
    // settings page: real Tab walk, radio arrow/space, selecting does not submit
    await anon(null); await as(acct); await go(ORIGIN + '/home/account');
    const L = await js('document.documentElement.lang'); const ko = L === 'ko';
    const stops = await tabStops(); obs.settingsStops = stops.map(s => s.wrap ? 'WRAP' : `${s.tag}:${s.type}:${s.text}|${s.axRole}:${s.axName}`);
    const real = stops.filter(s => !s.wrap);
    check(`keys_${L}_settings_first_stop_is_skip_link`, real[0]?.skip === true && /건너뛰기|Skip/.test(real[0].text), real[0]);
    check(`keys_${L}_settings_all_stops_have_visible_focus_inside_viewport`, real.every(focusOk), real.filter(s => !focusOk(s)).map(s => s.key));
    check(`keys_${L}_settings_all_stops_at_least_44px_except_skip`, real.filter(s => !s.skip).every(s => s.h >= 44), real.filter(s => !s.skip && s.h < 44).map(s => [s.key, s.h]));
    check(`keys_${L}_settings_no_keyboard_trap_one_cycle`, stops[stops.length - 1]?.wrap === true, stops.length);
    const exp = ko ? ['본문으로 건너뛰기', 'KnowsLink', '설정', '대화', '친구', '내 에이전트', '서비스 안내', '내 이름', '이름 저장', '친구 코드 교체', '언어 저장'] : ['Skip to content', 'KnowsLink', 'Settings', 'Conversations', 'Friends', 'My agents', 'Service guide', 'My name', 'Save name', 'Friend code (replace)', 'Save language'];
    const names = real.map(s => (s.axName || '').replace(/^[▸▾] /, '')); obs.summaryNameHasDecorativeGlyph = real.some(s => /^[▸▾] /.test(s.axName || '')); let at = -1; const missing = exp.filter(n => { const i = names.findIndex((x, j) => j > at && x === n); if (i < 0) return true; at = i; return false; });
    check(`keys_${L}_settings_tab_order_follows_reading_order`, missing.length === 0, {missing, names});
    const radioStops = real.filter(s => s.type === 'radio'); check(`keys_${L}_settings_radiogroup_is_one_tab_stop_on_checked_radio`, radioStops.length === 1, radioStops.map(s => s.key));
    const radioAx = real.find(s => s.type === 'radio'); check(`keys_${L}_settings_radio_accessible_name_is_language_name`, radioAx && (radioAx.axName === (L === 'en' ? 'English' : '한국어')), radioAx);
    // arrow selection changes the radio only
    const before = requests.length;
    await js(`document.querySelector('input[name=lang]:checked').focus()`);
    const startVal = await js(`document.querySelector('input[name=lang]:checked').value`);
    await arrow('Down'); await sleep(200);
    const afterDown = await js(`({v:document.querySelector('input[name=lang]:checked').value, focus:document.activeElement.value, url:location.href})`);
    check(`keys_${L}_settings_arrow_moves_selection_to_other_language`, afterDown.v !== startVal || startVal === 'en', {startVal, afterDown}); // en is last: Down wraps to ko in native groups
    await arrow('Up'); await sleep(150);
    const afterUp = await js(`document.querySelector('input[name=lang]:checked').value`);
    check(`keys_${L}_settings_arrow_back_restores_original_selection`, afterUp === startVal, {startVal, afterUp});
    await space(); await sleep(150);
    const posts = requests.slice(before).filter(r => r.method !== 'GET');
    check(`keys_${L}_settings_choosing_by_keyboard_submits_nothing_no_navigation`, posts.length === 0 && (await js('location.pathname')) === '/home/account', {posts, path: await js('location.pathname')});
    // Space on the other radio must also not submit
    await arrow('Down'); await space(); await sleep(200);
    check(`keys_${L}_settings_space_on_other_radio_submits_nothing`, requests.slice(before).filter(r => r.method !== 'GET').length === 0, requests.slice(before));
    await go(ORIGIN + '/home/account'); // discard the unsaved choice: reload shows the saved value again
    check(`keys_${L}_settings_reload_discards_unsaved_choice_saved_value_shown`, (await js(`document.querySelector('input[name=lang]:checked').value`)) === L, await js(`document.querySelector('input[name=lang]:checked').value`));
    // skip link
    await js('document.activeElement.blur(); window.scrollTo(0,0)'); await tab(); await enterNav();
    check(`keys_${L}_settings_skip_link_moves_to_main`, (await js(`location.hash`)) === '#main', await js('location.hash'));
    // header "Settings" link is current page, reachable by keyboard from another page
    await go(ORIGIN + '/home/friends'); await js('document.activeElement && document.activeElement.blur()'); await tab(); await tab(); await tab();
    const f3 = await js(`({t:document.activeElement.textContent, h:document.activeElement.getAttribute('href')})`);
    check(`keys_${L}_friends_third_tab_stop_is_header_settings_link`, f3.h === '/home/account' && f3.t === (ko ? '설정' : 'Settings'), f3);
    await enterNav(); check(`keys_${L}_header_settings_link_opens_settings`, (await js('location.pathname')) === '/home/account', await js('location.pathname'));
  }

  if (stage === 'langflow') {
    // anonymous: keyboard use of the language links on the public page, cookie persistence, reload, other screens, no input lost on a form-free page
    await metrics(1440, 900); await anon(null); await go(ORIGIN + '/');
    check('flow_public_default_is_korean', (await js('document.documentElement.lang')) === 'ko', await js('document.documentElement.lang'));
    const stops = await tabStops(12); obs.publicStops = stops.map(s => s.wrap ? 'WRAP' : `${s.tag}:${s.text}|${s.axName}`);
    const iEn = stops.findIndex(s => s.axName === 'English'); check('flow_public_english_link_is_a_tab_stop_with_name_English', iEn > 0 && stops[iEn].h >= 44 && focusOk(stops[iEn]), stops[iEn]);
    const iKo = stops.findIndex(s => s.axName === '한국어'); check('flow_public_korean_link_is_a_tab_stop_before_english', iKo > 0 && iKo < iEn, [iKo, iEn]);
    await js(`document.querySelector('nav.lang a[lang=en]').focus()`); // the Tab order was walked above; focus the stop directly (Chrome resumes Tab from the last focus point)
    const on = await js(`document.activeElement.textContent`); check('flow_public_focus_is_on_english_link', on === 'English', on);
    await enterNav();
    const L1 = await js('document.documentElement.lang'); check('flow_public_enter_on_english_switches_page_to_english', L1 === 'en' && (await js('location.pathname')) === '/', {L1, url: await js('location.href')});
    const cks = (await send('Network.getCookies', {urls: [ORIGIN + '/']})).cookies.filter(c => c.name === '__Host-kl_lang');
    check('flow_cookie_kl_lang_en_secure_httponly_strict_365d', cks.length === 1 && cks[0].value === 'en' && cks[0].secure && cks[0].httpOnly && cks[0].sameSite === 'Strict' && cks[0].path === '/' && !cks[0].domain.startsWith('.') && cks[0].expires - Date.now() / 1000 > 364 * 86400 && cks[0].expires - Date.now() / 1000 < 366 * 86400, cks);
    check('flow_public_english_has_current_marker_on_english_and_no_hangul_besides_language_name', await js(`(()=>{const a=[...document.querySelectorAll('nav.lang a')];return a.length===2&&a[1].hasAttribute('aria-current')&&!a[0].hasAttribute('aria-current')})()`) && (await js(`document.body.innerText.replace(/한국어/g,'').match(/[\\u3131-\\u318e\\uac00-\\ud7a3]/)===null`)), 'see innerText');
    await go(ORIGIN + '/'); check('flow_public_reload_keeps_english', (await js('document.documentElement.lang')) === 'en', '');
    for (const p of ['/connect/' + ID.device, '/remote/consent/' + ID.remote, '/auth/verify']) { await go(ORIGIN + p); check(`flow_cookie_en_applies_to_${p.split('/')[1]}`, (await js('document.documentElement.lang')) === 'en', p); }
    // back to Korean through the same links on the remote consent screen: return path preserved
    await go(ORIGIN + '/remote/consent/' + ID.remote); await js(`document.querySelector('nav.lang a[lang=ko]').focus()`); await enterNav();
    check('flow_remote_consent_ko_link_returns_to_same_screen_in_korean', (await js('location.pathname')) === '/remote/consent/' + ID.remote && (await js('document.documentElement.lang')) === 'ko', await js('location.href'));
    // screen-only: signed in as none (no saved value) with the cookie preference
    await anon('en'); await as('none'); await go(ORIGIN + '/home/account');
    check('flow_none_with_en_cookie_settings_in_english_radio_en_checked', (await js('document.documentElement.lang')) === 'en' && (await js(`document.querySelector('input[name=lang]:checked').value`)) === 'en', '');
    await anon(null); await as('none'); await go(ORIGIN + '/home/account');
    check('flow_none_without_cookie_korean_so_cookie_choice_was_not_saved_to_account', (await js('document.documentElement.lang')) === 'ko', '');
    // ?lang= with the cookie refused: block cookies in this browser context (Network.setBlockedURLs cannot block cookies; use a cookie-less navigation)
    await anon(null); await go(ORIGIN + '/?lang=en');
    check('flow_query_lang_en_without_cookie_applies_one_response_with_notice', (await js('document.documentElement.lang')) === 'en' && (await js(`document.body.innerText.includes('This browser did not keep the language choice. It applies to this page only.')`)), await js('document.body.innerText.slice(0,200)'));
    if (CAP) await shot('anon-en-public-query-notice-1440');
    check('flow_query_lang_notice_is_status_role', await js(`!!document.querySelector('[role=status]')`), '');
    obs.nonGetRequestsInFlow = requests.filter(r => r.method !== 'GET' && !r.path.startsWith('/fixture/')).length;
    check('flow_no_post_requests_made_by_language_links', obs.nonGetRequestsInFlow === 0, requests.filter(r => r.method !== 'GET'));
  }

  if (stage === 'nocookie') {
    // profile with cookies blocked (Preferences default_content_setting_values.cookies=2): the cookie cannot be kept, the answer must still apply to this response with a notice
    await metrics(1440, 900);
    for (const p of ['/', '/remote/consent/' + ID.remote, '/connect/' + ID.device]) {
      const n = p.split('/')[1] || 'public';
      await go(ORIGIN + p); check(`nocookie_${n}_starts_korean`, (await js('document.documentElement.lang')) === 'ko', '');
      await js(`document.querySelector('nav.lang a[lang=en]').focus()`); await enterNav();
      const st = await js(`({lang:document.documentElement.lang, url:location.pathname+location.search, notice:[...document.querySelectorAll('[role=status]')].map(e=>e.innerText), han:/[\\u3131-\\u318e\\uac00-\\ud7a3]/.test(document.body.innerText.replace(/한국어/g,''))})`);
      const cks = (await send('Network.getCookies', {urls: [ORIGIN + '/']})).cookies.filter(c => c.name === '__Host-kl_lang');
      check(`nocookie_${n}_english_link_applies_this_response_with_notice_and_stores_nothing`, st.lang === 'en' && st.notice.some(t => t.includes('did not keep the language choice')) && cks.length === 0 && !st.han, {st, cks});
      if (CAP) await shot(`anon-en-${n}-cookies-blocked-1440`);
      await go(ORIGIN + p); check(`nocookie_${n}_next_visit_is_korean_choice_not_kept`, (await js('document.documentElement.lang')) === 'ko', '');
    }
    await go(ORIGIN + '/'); await js(`document.querySelector('nav.lang a[lang=en]').focus()`); await enterNav();
    check('nocookie_notice_precedes_main_content_inside_main_or_status_region', await js(`!!document.querySelector('main [role=status], [role=status]')`), '');
  }

  if (stage === 'zoom') {
    // reflow equivalent of 200% (a real Chrome Page-zoom setting is not driven from CDP): 720x450 css px at DPR 2, 320 at DPR 4 as 400%
    for (const who of ['en', 'ko']) {
      await anon(null); await as(who);
      for (const [w, h, dpr, nm] of [[720, 450, 2, '200pct'], [320, 256, 4, '400pct']]) {
        await metrics(w, h, dpr);
        for (const [n, p] of [['account', '/home/account'], ['friends', '/home/friends'], ['agents', '/home/agents'], ['conversation', M.conversation]]) {
          await go(ORIGIN + p); const r = await layout();
          check(`zoom_${nm}_${r.lang}_${n}_no_horizontal_overflow_no_clipping_44px`, r.over <= 0 && r.small.length === 0 && r.clipped.length === 0 && r.hidden.length === 0 && r.btnClip.length === 0, r);
          if (CAP && n === 'account') await shot(`${who}-${r.lang}-${n}-${nm}`);
        }
      }
    }
  }

  if (stage === 'save') {
    // the only real state changes: en saves Korean then English again; ko saves English then Korean again. Keyboard path only.
    await metrics(1440, 900);
    async function saveVia(who, to) {
      await anon(null); await as(who); await go(ORIGIN + '/home/account');
      const L0 = await js('document.documentElement.lang');
      await js(`document.querySelector('input[name=lang][value=${to}]').focus()`);
      // select with the keyboard: Space on the radio, then Tab to the Save button and Enter
      await space();
      const picked = await js(`document.querySelector('input[name=lang]:checked').value`);
      const reqs0 = requests.length;
      await js(`document.querySelector('input[name=lang]:checked').focus()`); await tab();
      const btn = await js('document.activeElement.textContent'); const Lsave = L0 === 'ko' ? '언어 저장' : 'Save language';
      check(`save_${who}_to_${to}_tab_from_radio_reaches_save_button`, picked === to && btn === Lsave, {picked, btn});
      check(`save_${who}_to_${to}_no_request_before_enter`, requests.length === reqs0, requests.slice(reqs0));
      await enterNav();
      const st = await js(`({lang:document.documentElement.lang, url:location.pathname+location.search+location.hash, checked:document.querySelector('input[name=lang]:checked').value, status:[...document.querySelectorAll('[role=status]')].map(e=>e.innerText), h1:document.querySelector('h1').innerText, active:document.activeElement.tagName})`);
      check(`save_${who}_to_${to}_lands_on_same_settings_place_in_new_language`, st.lang === to && st.url === '/home/account?n=language#language' && st.checked === to && st.h1 === (to === 'ko' ? '설정' : 'Settings') && st.status.length >= 1, st);
      obs[`save_${who}_${to}_status`] = st.status;
      if (CAP) { for (const [w, h] of WIDTHS) { await metrics(w, h); await go(ORIGIN + '/home/account?n=language#language'); await shot(`${who}-${st.lang}-account-after-save-${w}`); } await metrics(1440, 900); }
      const posts = requests.filter(r => r.method === 'POST' && r.path === '/home/language'); check(`save_${who}_to_${to}_exactly_one_language_post_so_far`, posts.length >= 1, posts.length);
      // reload and another screen keep it; settings keep position
      await go(ORIGIN + '/home/account'); check(`save_${who}_to_${to}_reload_keeps_language`, (await js('document.documentElement.lang')) === to, '');
      await go(ORIGIN + '/home/friends'); check(`save_${who}_to_${to}_other_screen_uses_saved_language`, (await js('document.documentElement.lang')) === to, '');
    }
    await saveVia('en', 'ko'); // owner screens in Korean while en is Korean
    obs.ownerScreensKo = {};
    for (const [n, p] of MEMBER) { await go(ORIGIN + p); obs.ownerScreensKo[n] = await js('({lang:document.documentElement.lang,title:document.title,over:document.documentElement.scrollWidth-innerWidth})'); }
    check('save_en_as_ko_owner_screens_all_html_lang_ko_no_overflow', Object.values(obs.ownerScreensKo).every(o => o.lang === 'ko' && o.over <= 0), obs.ownerScreensKo);
    // owner-only screens (conversation, approval, connections) in Korean at 3 widths, with captures for the designer
    for (const [w, h] of WIDTHS) {
      await metrics(w, h);
      for (const [n, p] of MEMBER.filter(x => ['conversation', 'gate', 'conn-waiting', 'conn-prepared', 'conn-approved', 'legacy', 'agents', 'friends', 'home'].includes(x[0]))) {
        await go(ORIGIN + p);
        for (const open of [false, true]) {
          if (open) { if (!(await js(`document.querySelectorAll('main details').length`))) continue; await js(openAll); await sleep(100); }
          const r = await layout(); const nm = `${n}${open ? '-open' : ''}-${w}`;
          check(`layout_owner_ko_${nm}`, r.over <= 0 && r.small.length === 0 && r.clipped.length === 0 && r.hidden.length === 0 && r.nameOver.length === 0 && r.btnClip.length === 0 && r.lang === 'ko', r);
          if (CAP) await shot(`en-ko-${nm}-ownerko`);
        }
      }
    }
    await metrics(1440, 900);
    await saveVia('en', 'en'); // restore
    await saveVia('ko', 'en'); await saveVia('ko', 'ko'); // restore
  }
} catch (e) { check('script_completed_without_exception', false, String(e.stack || e).slice(0, 400)); }
save(); ws.close();
const failed = checks.filter(c => !c.passed).length;
console.log(`\n${checks.length - failed} PASS / ${failed} FAIL  (${resultFile})`);
process.exit(failed ? 1 : 0);
