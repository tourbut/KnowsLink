// Retest 71dcac7: observe the closed receipt selector (observation only, no pass/fail change). Usage: node selector-observe.mjs <cdp-port> <label>
import fs from 'node:fs'; import path from 'node:path'; import {fileURLToPath} from 'node:url'; import {createHash} from 'node:crypto';
const out = path.dirname(fileURLToPath(import.meta.url));
const [port, label] = [process.argv[2], process.argv[3]];
const ORIGIN = 'http://127.0.0.1:18928', AB = '0199a3f2-4c10-7a11-8b22-334455667788';
const target = (await (await fetch(`http://127.0.0.1:${port}/json/list`)).json()).find(t => t.type === 'page');
const ws = new WebSocket(target.webSocketDebuggerUrl); await new Promise((res, rej) => { ws.onopen = res; ws.onerror = rej; });
let next = 0; const pending = new Map(), waiters = [], rows = [], captures = [];
ws.onmessage = ev => { const m = JSON.parse(ev.data); if (m.id) { const p = pending.get(m.id); pending.delete(m.id); m.error ? p.reject(new Error(m.error.message)) : p.resolve(m.result); } else if (m.method === 'Page.loadEventFired') waiters.splice(0).forEach(f => f()); };
const send = (method, params = {}) => new Promise((resolve, reject) => { const id = ++next; pending.set(id, {resolve, reject}); ws.send(JSON.stringify({id, method, params})); });
const js = async e => (await send('Runtime.evaluate', {expression: e, returnByValue: true})).result.value;
const sleep = ms => new Promise(r => setTimeout(r, ms));
const go = async u => { const w = new Promise(r => waiters.push(r)); await send('Page.navigate', {url: u}); await Promise.race([w, sleep(8000)]); await sleep(200); };
await send('Page.enable'); await send('Network.enable');
const post = async (p, s) => fetch(ORIGIN + p, {method: 'POST', headers: {Origin: ORIGIN, ...(s && {Cookie: '__Host-kl_session=' + s})}});
await post('/fixture/recover'); const r = await post('/fixture/as/a'); const token = /__Host-kl_session=([^;]+)/.exec(r.headers.get('set-cookie'))[1];
await send('Network.setCookie', {name: '__Host-kl_session', value: token, url: ORIGIN + '/', secure: true, httpOnly: true, sameSite: 'Strict'});
await send('Network.deleteCookies', {name: '__Host-kl_receipt', url: ORIGIN + '/'});
const measure = `(()=>{const s=document.getElementById('receipt-agent');const o=s.options[s.selectedIndex];const t=o.text;const cs=getComputedStyle(s);const c=document.createElement('canvas').getContext('2d');c.font=cs.fontStyle+' '+cs.fontWeight+' '+cs.fontSize+' '+cs.fontFamily;
const avail=s.clientWidth-parseFloat(cs.paddingLeft)-parseFloat(cs.paddingRight)-24;const w=x=>c.measureText(x).width;
let n=0;while(n<t.length&&w(t.slice(0,n+1))<=avail)n++;const visible=t.slice(0,n);
const parts=t.split(' / ');const dis=parts[0],name=parts[1],kind=parts[2];
const nameChars=Math.max(0,Math.min(name.length,visible.length-(dis.length+3)));
return {value:s.value,text:t,textWidth:Math.round(w(t)),avail:Math.round(avail),selectWidth:Math.round(s.getBoundingClientRect().width),visible,distinguisherFits:visible.startsWith(dis),nameCharsVisible:nameChars,nameChars:name.length,kindVisibleChars:Math.max(0,visible.length-(dis.length+3+name.length+3)),kindChars:kind.length,kindFullyVisible:n>=t.length,font:c.font}})()`;
for (const [w, h] of [[1440, 900], [390, 844], [320, 640]]) {
  await send('Emulation.setDeviceMetricsOverride', {width: w, height: h, deviceScaleFactor: 1, mobile: false});
  for (const ag of ['agent_a', 'agent_a2']) {
    await go(`${ORIGIN}/home/receipts?agent=${ag}&id=${AB}`);
    const m = await js(measure); rows.push({width: w, agent: ag, ...m});
    if (w === 1440 || true) { const {data} = await send('Page.captureScreenshot', {format: 'png', clip: await js(`(()=>{const r=document.getElementById('receipt-agent').getBoundingClientRect();return {x:Math.max(0,r.left+scrollX-16),y:Math.max(0,r.top+scrollY-40),width:Math.min(innerWidth,r.width+32),height:r.height+64,scale:2}})()`)}); const b = Buffer.from(data, 'base64'); const f = `${label}-selector-closed-${ag}-${w}.png`; fs.writeFileSync(path.join(out, f), b); captures.push({file: f, width: w, agent: ag, sha256: createHash('sha256').update(b).digest('hex')}); }
  }
}
// one full-viewport capture per width for designer context (agent_a)
for (const [w, h] of [[1440, 900], [390, 844], [320, 640]]) { await send('Emulation.setDeviceMetricsOverride', {width: w, height: h, deviceScaleFactor: 1, mobile: false}); await go(`${ORIGIN}/home/receipts?agent=agent_a&id=${AB}`); const {data} = await send('Page.captureScreenshot', {format: 'png'}); const b = Buffer.from(data, 'base64'); const f = `${label}-viewport-agent_a-${w}.png`; fs.writeFileSync(path.join(out, f), b); captures.push({file: f, width: w, agent: 'agent_a', full: true, sha256: createHash('sha256').update(b).digest('hex')}); }
fs.writeFileSync(path.join(out, `selector-observe-${label}.json`), JSON.stringify({origin: ORIGIN, browser: (await (await fetch(`http://127.0.0.1:${port}/json/version`)).json()).Browser, rows, captures}, null, 2));
for (const x of rows) console.log(x.width, x.agent, `tw=${x.textWidth} avail=${x.avail}`, `dis=${x.distinguisherFits} name=${x.nameCharsVisible}/${x.nameChars} kind=${x.kindVisibleChars}/${x.kindChars}`, JSON.stringify(x.visible));
ws.close();
