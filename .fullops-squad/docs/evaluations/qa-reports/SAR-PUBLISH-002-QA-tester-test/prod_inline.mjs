// Read-only: measure line pitch of paragraphs that contain sentence-embedded links (V-01 regression) at 640px and 320px.
import { spawn } from 'node:child_process';
import fs from 'node:fs'; import path from 'node:path'; import os from 'node:os';
const ORIGIN = 'https://link.knowslog.com';
const prof = fs.mkdtempSync(path.join(os.tmpdir(), 'pub002i-'));
const chrome = spawn('C:/Program Files/Google/Chrome/Application/chrome.exe', ['--headless=new','--disable-gpu','--remote-debugging-port=9345','--user-data-dir='+prof,'about:blank'], {stdio:'ignore'});
const sleep = ms => new Promise(r => setTimeout(r, ms));
for (let i=0;i<50;i++){ try{ await (await fetch('http://127.0.0.1:9345/json/version')).json(); break;}catch{await sleep(200);} }
const probe = `(() => {
 const res=[];
 for (const a of document.querySelectorAll('main a, section a, p a, li a')) {
  if(!a.offsetParent) continue;
  const par=a.parentElement; if(!['P','LI','SPAN'].includes(par.tagName)) continue;
  const others=[...par.childNodes].filter(n=>n!==a&&n.nodeType===3&&n.textContent.trim().length>3).length;
  if(!others) continue; // not embedded in a sentence
  const cs=getComputedStyle(a), pr=document.createRange(); pr.selectNodeContents(par);
  const tops=[...new Set([...pr.getClientRects()].map(r=>Math.round(r.top)))].sort((x,y)=>x-y);
  const pitch=tops.slice(1).map((t,i)=>t-tops[i]).filter(d=>d>4);
  res.push({text:a.textContent.trim().slice(0,32), parent:par.tagName+'.'+par.className, display:cs.display, minH:cs.minHeight, bodyLine:getComputedStyle(par).lineHeight, lines:tops.length, pitch});
 }
 return res; })()`;
const out=[];
for (const lang of ['ko','en']) for (const p of ['/','/privacy','/terms']) for (const w of [1280,640,320]) {
  const t = await (await fetch('http://127.0.0.1:9345/json/new?about:blank',{method:'PUT'})).json();
  const ws = new WebSocket(t.webSocketDebuggerUrl); await new Promise(r=>ws.onopen=r);
  let id=0; const pend=new Map(); ws.onmessage=e=>{const m=JSON.parse(e.data); if(m.id&&pend.has(m.id)){pend.get(m.id)(m);pend.delete(m.id);}};
  const send=(method,params={})=>new Promise(r=>{const i=++id;pend.set(i,r);ws.send(JSON.stringify({id:i,method,params}));});
  await send('Network.enable'); await send('Network.clearBrowserCookies');
  await send('Emulation.setDeviceMetricsOverride',{width:w,height:800,deviceScaleFactor:1,mobile:false});
  await send('Page.navigate',{url:`${ORIGIN}/lang/${lang}?back=${encodeURIComponent(p)}`}); await sleep(1500);
  const r = await send('Runtime.evaluate',{expression:probe,returnByValue:true});
  out.push({lang,page:p,width:w,links:r.result.result.value}); ws.close(); fetch('http://127.0.0.1:9345/json/close/'+t.id);
}
fs.writeFileSync('prod-inline.json',JSON.stringify(out,null,1)); chrome.kill(); process.exit(0);
