// Measure only task-owned synthetic HTML in a dedicated Chrome profile; never submit forms.
import fs from 'node:fs';
import path from 'node:path';
import {fileURLToPath} from 'node:url';
import {createHash} from 'node:crypto';

const out = path.dirname(fileURLToPath(import.meta.url));
const origin = 'http://127.0.0.1:18896';
const version = await (await fetch('http://127.0.0.1:18897/json/version')).json();
const target = (await (await fetch('http://127.0.0.1:18897/json/list')).json()).find(t => t.type === 'page');
const ws = new WebSocket(target.webSocketDebuggerUrl);
await new Promise((resolve, reject) => { ws.onopen = resolve; ws.onerror = reject; });
let next = 0;
const pending = new Map(), checks = [], captures = [], requests = [];
const report = {source: 'f2b8cb622dc397034423582978e77560e0355771', browser: version.Browser, target: target.id, origin, checks, captures, requests};
const save = () => fs.writeFileSync(path.join(out, 'results.json'), JSON.stringify(report, null, 2));
ws.onmessage = event => {
  const m = JSON.parse(event.data);
  if (m.id) { const p = pending.get(m.id); pending.delete(m.id); m.error ? p.reject(new Error(m.error.message)) : p.resolve(m.result); }
  if (m.method === 'Network.requestWillBeSent') requests.push({url: m.params.request.url, method: m.params.request.method});
};
const send = (method, params = {}) => new Promise((resolve, reject) => { const id = ++next; pending.set(id, {resolve, reject}); ws.send(JSON.stringify({id, method, params})); });
async function js(expression) {
  const r = await send('Runtime.evaluate', {expression, returnByValue: true, awaitPromise: true});
  if (r.exceptionDetails) throw new Error(JSON.stringify(r.exceptionDetails));
  return r.result.value;
}
function check(name, passed, value) { checks.push({name, passed, value}); save(); if (!passed) throw new Error(name); }
async function go(file) {
  const url = `${origin}/${file}.html`;
  await send('Page.navigate', {url});
  for (let n = 0; n < 100; n++) {
    if (await js(`location.href===${JSON.stringify(url)} && document.readyState==='complete'`)) return;
    await new Promise(resolve => setTimeout(resolve, 50));
  }
  throw new Error('navigation timeout');
}
function luminance(rgb) {
  const v = rgb.match(/[\d.]+/g).slice(0, 3).map(Number).map(x => x / 255).map(x => x <= 0.04045 ? x / 12.92 : ((x + 0.055) / 1.055) ** 2.4);
  return v[0] * 0.2126 + v[1] * 0.7152 + v[2] * 0.0722;
}
async function color(selector, name) {
  const value = await js(`(()=>{const e=document.querySelector(${JSON.stringify(selector)}),s=getComputedStyle(e);return {foreground:s.color,background:s.backgroundColor,height:e.getBoundingClientRect().height,outline:s.outlineWidth,outlineStyle:s.outlineStyle,outlineOffset:s.outlineOffset,outlineColor:s.outlineColor,font:s.fontSize}})()`);
  const a = luminance(value.foreground), b = luminance(value.background);
  value.ratio = (Math.max(a, b) + 0.05) / (Math.min(a, b) + 0.05);
  check(name, value.ratio >= 4.5 && value.height >= 44, value);
  return value;
}
async function shot(name) {
  const {data} = await send('Page.captureScreenshot', {format: 'png'});
  const bytes = Buffer.from(data, 'base64');
  fs.writeFileSync(path.join(out, name + '.png'), bytes);
  captures.push({file: name + '.png', sha256: createHash('sha256').update(bytes).digest('hex')}); save();
}
try {
  check('contrast-formula-reference', Math.abs((luminance('rgb(255,255,255)') + 0.05) / (luminance('rgb(0,90,98)') + 0.05) - 7.958060) < 0.000001);
  await send('Page.enable'); await send('Network.enable'); await send('DOM.enable'); await send('CSS.enable');
  for (const width of [1440, 390, 320]) {
    await send('Emulation.setDeviceMetricsOverride', {width, height: width === 1440 ? 900 : 844, deviceScaleFactor: 1, mobile: false});
    await go('start-google');
    check('viewport-' + width, await js(`innerWidth===${width}`));
    await color('.hero .primary', 'public-contrast-' + width);
    const spacing = await js("(()=>{const e=document.querySelector('.examples');return {gap:e.nextElementSibling.getBoundingClientRect().top-e.getBoundingClientRect().bottom,margin:getComputedStyle(e).marginBottom}})()");
    check('example-spacing-' + width, Math.abs(spacing.gap - 24) < 0.02, spacing);
    await js("document.querySelectorAll('details').forEach(e=>e.open=true)");
    check('public-open-overflow-' + width, await js('document.documentElement.scrollWidth<=innerWidth'));
    await js("document.querySelector('.examples').scrollIntoView()"); await shot('examples-' + width);
    await go('human-home-empty');
    await color('.conversation-list .primary', 'member-contrast-' + width);
    check('member-overflow-' + width, await js('document.documentElement.scrollWidth<=innerWidth'));
    check('ordinary-link-color-' + width, await js("getComputedStyle(document.querySelector('.conversation-list a:not(.button)')).color==='rgb(24, 43, 49)'"));
    await shot('member-empty-' + width);
    const {root} = await send('DOM.getDocument');
    const {nodeId} = await send('DOM.querySelector', {nodeId: root.nodeId, selector: '.conversation-list .primary'});
    for (const state of ['hover', 'focus-visible']) {
      await send('CSS.forcePseudoState', {nodeId, forcedPseudoClasses: [state]});
      const value = await color('.conversation-list .primary', `member-${state}-${width}`);
      if (state === 'focus-visible') check('focus-outline-' + width, value.outline === '3px' && value.outlineStyle === 'solid' && value.outlineOffset === '2px' && value.outlineColor === 'rgb(0, 90, 98)', value);
    }
    await send('CSS.forcePseudoState', {nodeId, forcedPseudoClasses: []});
    await go('human-home');
    check('populated-list-color-' + width, await js("[...document.querySelectorAll('.conversation-list a:not(.button)')].every(e=>getComputedStyle(e).color==='rgb(24, 43, 49)')"));
    if (width === 1440) await shot('member-populated-1440');
    await go('human-account');
    check('member-account-overflow-' + width, await js('document.documentElement.scrollWidth<=innerWidth'));
  }
  check('only-fixture-get-requests', requests.every(r => r.url.startsWith(origin + '/') && r.method === 'GET'), requests);
  console.log(`PASS ${checks.length} checks`);
} finally { save(); ws.close(); }
