// Read-only production QA via headless Chrome + CDP: language links, Tab order, 640px reflow.
// Only GETs/link navigation/Tab keys/scroll on https://link.knowslog.com; no form submit, no login.
import { spawn } from 'node:child_process';
import fs from 'node:fs';
import path from 'node:path';
import os from 'node:os';
const ORIGIN = 'https://link.knowslog.com';
const out = process.argv[2];
fs.mkdirSync(path.join(out, 'raw'), { recursive: true });
const prof = fs.mkdtempSync(path.join(os.tmpdir(), 'pub002-'));
const chrome = spawn('C:/Program Files/Google/Chrome/Application/chrome.exe', ['--headless=new', '--disable-gpu', '--remote-debugging-port=9344', '--user-data-dir=' + prof, '--hide-scrollbars', 'about:blank'], { stdio: 'ignore' });
const sleep = ms => new Promise(r => setTimeout(r, ms));
for (let i = 0; i < 50; i++) { try { await (await fetch('http://127.0.0.1:9344/json/version')).json(); break; } catch { await sleep(200); } }
async function session() {
  const t = await (await fetch('http://127.0.0.1:9344/json/new?about:blank', { method: 'PUT' })).json();
  const ws = new WebSocket(t.webSocketDebuggerUrl); await new Promise(r => ws.onopen = r);
  let id = 0; const pend = new Map();
  ws.onmessage = e => { const m = JSON.parse(e.data); if (m.id && pend.has(m.id)) { pend.get(m.id)(m); pend.delete(m.id); } };
  const send = (method, params = {}) => new Promise(r => { const i = ++id; pend.set(i, r); ws.send(JSON.stringify({ id: i, method, params })); });
  const ev = async expr => { const r = await send('Runtime.evaluate', { expression: expr, returnByValue: true }); if (r.result.exceptionDetails) throw new Error(JSON.stringify(r.result.exceptionDetails).slice(0, 300)); return r.result.result.value; };
  return { send, ev, close: () => { ws.close(); fetch('http://127.0.0.1:9344/json/close/' + t.id); } };
}
async function go(s, url) { await s.send('Page.navigate', { url }); await sleep(1200); }
const pages = ['/', '/privacy', '/terms'];
const res = { lang: [], tab: [], reflow: [] };
const state = () => `({url:location.pathname+location.search, lang:document.documentElement.lang, title:document.title, h1:document.querySelector('h1')?.textContent.trim(), cur:[...document.querySelectorAll('[aria-current]')].map(a=>a.textContent.trim()+'|'+a.getAttribute('href')+'|'+a.getAttribute('aria-current')), toggles:[...document.querySelectorAll('a[href^="/lang/"]')].map(a=>a.getAttribute('href'))})`;

// 1. language links (fresh session per page: cookie jar cleared)
for (const p of pages) {
  const s = await session();
  await s.send('Network.enable'); await s.send('Network.clearBrowserCookies');
  await go(s, ORIGIN + p);
  const a0 = await s.ev(state());
  const steps = [{ step: 'initial', ...a0 }];
  for (const target of ['en', 'ko']) {
    const href = await s.ev(`document.querySelector('a[href^="/lang/${target}"]')?.getAttribute('href')`);
    await s.send('Page.navigate', { url: ORIGIN + href }); await sleep(1500);
    steps.push({ step: 'click ' + href, ...(await s.ev(state())) });
  }
  res.lang.push({ page: p, steps });
  s.close();
}

// 2. Tab traversal, KO and EN, with focus-style probe (fresh jar each time)
const focusProbe = `(() => { const e=document.activeElement; if(!e||e===document.body) return {body:true}; const cs=getComputedStyle(e), r=e.getBoundingClientRect(); return {tag:e.tagName, text:e.textContent.trim().slice(0,40), href:e.getAttribute('href'), cur:e.getAttribute('aria-current'), outline:cs.outlineStyle+' '+cs.outlineWidth+' '+cs.outlineColor, offset:cs.outlineOffset, boxShadow:cs.boxShadow, fv:e.matches(':focus-visible'), top:Math.round(r.top+scrollY), left:Math.round(r.left), inView:r.top>=0&&r.bottom<=innerHeight}; })()`;
const linkOrder = `[...document.querySelectorAll('a[href],button,input,select,textarea')].filter(e=>e.offsetParent!==null||getComputedStyle(e).position==='fixed').map(e=>{const r=e.getBoundingClientRect();return {text:(e.textContent||e.value||'').trim().slice(0,40),href:e.getAttribute('href'),top:Math.round(r.top+scrollY),left:Math.round(r.left)}})`;
for (const lang of ['ko', 'en']) for (const p of pages) {
  const s = await session();
  await s.send('Network.enable'); await s.send('Network.clearBrowserCookies');
  await s.send('Emulation.setDeviceMetricsOverride', { width: 1280, height: 900, deviceScaleFactor: 1, mobile: false });
  await s.send('Page.navigate', { url: `${ORIGIN}/lang/${lang}?back=${encodeURIComponent(p)}` }); await sleep(1500);
  const dom = await s.ev(linkOrder);
  const seq = [];
  for (let i = 0; i < 60; i++) {
    await s.send('Input.dispatchKeyEvent', { type: 'rawKeyDown', key: 'Tab', code: 'Tab', windowsVirtualKeyCode: 9 });
    await s.send('Input.dispatchKeyEvent', { type: 'keyUp', key: 'Tab', code: 'Tab', windowsVirtualKeyCode: 9 });
    await sleep(60);
    const f = await s.ev(focusProbe);
    if (f.body || (seq.length && seq[0].href === f.href && seq[0].text === f.text && seq.length > 1)) { seq.push(f); break; }
    seq.push(f);
  }
  const shot = await s.send('Page.captureScreenshot', { format: 'png', captureBeyondViewport: true });
  fs.writeFileSync(path.join(out, 'raw', `tab-${lang}-${p.replace(/\//g, '_') || 'root'}.png`), Buffer.from(shot.result.data, 'base64'));
  res.tab.push({ lang, page: p, url: await s.ev('location.pathname'), htmlLang: await s.ev('document.documentElement.lang'), dom, seq });
  s.close();
}

// 3. 640px (= 200% of 1280) reflow + inline link line spacing
const reflowProbe = `(() => { const d=document,de=d.documentElement;
 const wide=[...d.querySelectorAll('body *')].filter(e=>e.getBoundingClientRect().right>innerWidth+1 && !e.closest('table,.table-wrap,[style*="overflow"]')).slice(0,5).map(e=>e.tagName+'.'+e.className+':'+Math.round(e.getBoundingClientRect().right));
 const tabs=[...d.querySelectorAll('table')].map(t=>({cw:t.parentElement.clientWidth,sw:t.scrollWidth,ov:getComputedStyle(t.parentElement).overflowX}));
 const inl=[...d.querySelectorAll('p a, li a')].filter(a=>a.offsetParent).map(a=>{const cs=getComputedStyle(a);const rects=[...a.getClientRects()];return {text:a.textContent.trim().slice(0,30),display:cs.display,minH:cs.minHeight,lineH:Math.round(parseFloat(cs.lineHeight)||0),maxRectH:Math.round(Math.max(...rects.map(r=>r.height))),parent:Math.round(parseFloat(getComputedStyle(a.parentElement).lineHeight)||0)}});
 const fixedOverlap=[...d.querySelectorAll('body *')].filter(e=>getComputedStyle(e).position==='fixed'&&e.offsetParent!==undefined).length;
 const clipped=[...d.querySelectorAll('h1,h2,p,li,a,td,th')].filter(e=>e.scrollWidth>e.clientWidth+1 && getComputedStyle(e).display!=='inline' && !e.closest('table,.table-wrap')).slice(0,5).map(e=>e.tagName+':'+e.textContent.trim().slice(0,20));
 return {iw:innerWidth,sw:de.scrollWidth,hOverflow:de.scrollWidth>innerWidth,wide,tabs,inl,clipped,fixedOverlap,h:de.scrollHeight}; })()`;
for (const lang of ['ko', 'en']) for (const p of pages) for (const w of [640, 320]) {
  const s = await session();
  await s.send('Network.enable'); await s.send('Network.clearBrowserCookies');
  await s.send('Emulation.setDeviceMetricsOverride', { width: w, height: 800, deviceScaleFactor: w === 640 ? 2 : 1, mobile: false });
  await s.send('Page.navigate', { url: `${ORIGIN}/lang/${lang}?back=${encodeURIComponent(p)}` }); await sleep(1500);
  const m = await s.ev(reflowProbe);
  const shot = await s.send('Page.captureScreenshot', { format: 'png', captureBeyondViewport: true, clip: { x: 0, y: 0, width: w, height: m.h, scale: 1 } });
  fs.writeFileSync(path.join(out, 'raw', `reflow-${lang}-${p.replace(/\//g, '_') || 'root'}-${w}.png`), Buffer.from(shot.result.data, 'base64'));
  res.reflow.push({ lang, page: p, width: w, ...m });
  s.close();
}
fs.writeFileSync(path.join(out, 'prod-result.json'), JSON.stringify(res, null, 1));
chrome.kill();
process.exit(0);
