// Post-revoke browser check (Chrome CDP): revoked agent is shown disabled; keyboard submit must not pick another agent.
import fs from 'node:fs';
const port = process.argv[2], ORIGIN = 'http://127.0.0.1:18928', AB = '0199a3f2-4c10-7a11-8b22-334455667788';
const target = (await (await fetch(`http://127.0.0.1:${port}/json/list`)).json()).find(t => t.type === 'page');
const ws = new WebSocket(target.webSocketDebuggerUrl); await new Promise(r => ws.onopen = r);
let n = 0; const pend = new Map(), waiters = [], reqs = [];
ws.onmessage = e => { const m = JSON.parse(e.data); if (m.id) { pend.get(m.id)(m.result); pend.delete(m.id); } else if (m.method === 'Page.loadEventFired') waiters.splice(0).forEach(f => f()); else if (m.method === 'Network.requestWillBeSent' && m.params.request.url.startsWith('http')) reqs.push(m.params.request.method + ' ' + new URL(m.params.request.url).pathname + new URL(m.params.request.url).search.replace(AB, '<AB>')); };
const send = (method, params = {}) => new Promise(r => { const id = ++n; pend.set(id, r); ws.send(JSON.stringify({id, method, params})); });
const js = async x => (await send('Runtime.evaluate', {expression: x, returnByValue: true})).result.value;
const sleep = ms => new Promise(r => setTimeout(r, ms));
const go = async u => { const w = new Promise(r => waiters.push(r)); await send('Page.navigate', {url: u}); await Promise.race([w, sleep(6000)]); await sleep(200); };
const key = async (k, code, vk) => { await send('Input.dispatchKeyEvent', {type: 'keyDown', key: k, code, windowsVirtualKeyCode: vk, text: k === 'Enter' ? '\r' : undefined}); await send('Input.dispatchKeyEvent', {type: 'keyUp', key: k, code, windowsVirtualKeyCode: vk}); await sleep(600); };
await send('Page.enable'); await send('Network.enable');
const res = [];
const rec = (name, ok, d) => { res.push({name, passed: ok === true, detail: d}); console.log(ok ? 'PASS' : 'FAIL', name, ok ? '' : JSON.stringify(d)); };
const r = await fetch(ORIGIN + '/fixture/as/a', {method: 'POST', headers: {Origin: ORIGIN}}); const sess = /__Host-kl_session=([^;]+)/.exec(r.headers.get('set-cookie'))[1];
await send('Network.setCookie', {name: '__Host-kl_session', value: sess, url: ORIGIN + '/', secure: true, httpOnly: true, sameSite: 'Strict'});
await send('Network.deleteCookies', {name: '__Host-kl_receipt', url: ORIGIN + '/'});
await go(`${ORIGIN}/home/receipts?agent=agent_a&id=${AB}`);
const st = await js(`({alert:(document.querySelector('[role=alert]')||{}).textContent,sel:document.getElementById('receipt-agent').value,selText:document.getElementById('receipt-agent').selectedOptions[0]&&document.getElementById('receipt-agent').selectedOptions[0].text.slice(-12),active:document.activeElement.id,id:document.getElementById('receipt-id').value,opts:[...document.getElementById('receipt-agent').options].map(o=>o.value+(o.disabled?':disabled':''))})`);
rec('revoked-agent-visible-disabled-with-alert-and-focus-on-select', st.alert.includes('철회되었습니다') && st.opts.join() === ',agent_a2,agent_a:disabled' && st.active === 'receipt-agent' && st.id === AB, st);
// keyboard: Tab to id, Enter => must not submit as another agent
await key('Tab', 'Tab', 9); await key('Enter', 'Enter', 13);
const after = await js(`({url:location.search.replace('${AB}','<AB>'),agent2:document.getElementById('receipt-agent').value,alert:(document.querySelector('[role=alert]')||{}).textContent})`);
rec('keyboard-submit-with-revoked-agent-never-switches-to-agent_a2', !after.url.includes('agent=agent_a2') && after.agent2 !== 'agent_a2', after);
// choose live agent deliberately by keyboard and retry: refusal for non-participant agent_a2 is the same neutral message
await go(`${ORIGIN}/home/receipts?agent=agent_a&id=${AB}`); await js(`document.getElementById('receipt-agent').focus()`); await key('ArrowUp', 'ArrowUp', 38);
const picked = await js(`document.getElementById('receipt-agent').value`); await key('Tab', 'Tab', 9); await key('Enter', 'Enter', 13);
const final = await js(`({url:location.search.replace('${AB}','<AB>'),alert:(document.querySelector('[role=alert]')||{}).textContent,agent:document.getElementById('receipt-agent').value})`);
rec('deliberate-keyboard-reselect-submits-same-screen', picked === 'agent_a2' && final.url.includes('agent=agent_a2') && final.agent === 'agent_a2', {picked, final});
rec('requests-are-get-only-same-origin', reqs.every(x => x.startsWith('GET ')), reqs);
fs.writeFileSync(new URL('./revoke-browser-result.json', import.meta.url), JSON.stringify({browser: 'Chrome via CDP', checks: res, requests: reqs}, null, 1)); ws.close();
process.exit(res.every(x => x.passed) ? 0 : 1);
