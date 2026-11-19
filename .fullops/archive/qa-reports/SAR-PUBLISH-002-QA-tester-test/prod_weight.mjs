// Read-only: measure aria-current weight on prod /privacy.
import { spawn } from 'node:child_process'; import fs from 'node:fs'; import path from 'node:path'; import os from 'node:os';
const prof = fs.mkdtempSync(path.join(os.tmpdir(), 'pub002l-'));
const chrome = spawn('C:/Program Files/Google/Chrome/Application/chrome.exe', ['--headless=new','--disable-gpu','--remote-debugging-port=9347','--user-data-dir='+prof,'about:blank'], {stdio:'ignore'});
const sleep = ms => new Promise(r => setTimeout(r, ms));
for (let i=0;i<50;i++){ try{ await (await fetch('http://127.0.0.1:9347/json/version')).json(); break;}catch{await sleep(200);} }
const t = await (await fetch('http://127.0.0.1:9347/json/new?about:blank',{method:'PUT'})).json();
const ws = new WebSocket(t.webSocketDebuggerUrl); await new Promise(r=>ws.onopen=r);
let id=0; const pend=new Map(); ws.onmessage=e=>{const m=JSON.parse(e.data); if(m.id&&pend.has(m.id)){pend.get(m.id)(m);pend.delete(m.id);}};
const send=(method,params={})=>new Promise(r=>{const i=++id;pend.set(i,r);ws.send(JSON.stringify({id:i,method,params}));});
await send("Page.navigate",{url:"https://link.knowslog.com/privacy"}); await sleep(3500);
const r = await send('Runtime.evaluate',{expression:`JSON.stringify([...document.querySelectorAll('footer a, a[aria-current]')].map(a=>[a.textContent.trim(),a.getAttribute('aria-current'),getComputedStyle(a).fontWeight,getComputedStyle(a).textDecorationLine]))`,returnByValue:true});
console.log(r.result.result.value); chrome.kill(); process.exit(0);
