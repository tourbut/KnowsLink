// Measure link heights, line spacing and horizontal overflow of the fixture pages at 1440/390/320 via CDP (V-01..V-03 recheck).
import { spawn } from 'node:child_process';
import fs from 'node:fs';
import path from 'node:path';
import { pathToFileURL } from 'node:url';
const [,, srcDir, outDir] = process.argv;
fs.mkdirSync(outDir, { recursive: true });
const chrome = spawn('C:/Program Files/Google/Chrome/Application/chrome.exe', ['--headless=new','--disable-gpu','--remote-debugging-port=9334','--user-data-dir='+path.join(outDir,'prof'),'--hide-scrollbars','about:blank'], {stdio:'ignore'});
const sleep = ms => new Promise(r => setTimeout(r, ms));
for (let i=0;i<50;i++){ try{ await (await fetch('http://127.0.0.1:9334/json/version')).json(); break;}catch{await sleep(200);} }
const probe = `(() => {
 const vis=[...document.querySelectorAll('a')].filter(a=>a.offsetParent!==null);
 const h=a=>Math.round(a.getBoundingClientRect().height);
 const inline=vis.filter(a=>a.closest('.policy p,.policy li,.auth p')&&!a.closest('footer'));
 const short=vis.filter(a=>h(a)<44&&!inline.includes(a)).map(a=>a.textContent.trim().slice(0,24)+':'+h(a));
 const cur=document.querySelector('footer a[aria-current]');
 const links=document.querySelector('p.links');
 const lk=links?[...links.querySelectorAll('a')].map(a=>a.getBoundingClientRect()):[];
 return {iw:innerWidth, hOverflow:document.documentElement.scrollWidth>innerWidth,
  inlineLinks:inline.map(a=>({t:a.textContent.trim().slice(0,20),display:getComputedStyle(a).display,lineH:Math.round(a.parentElement.getBoundingClientRect().height/Math.max(1,Math.round(a.parentElement.getBoundingClientRect().height/parseFloat(getComputedStyle(a.parentElement).lineHeight))))})),
  nonInlineBelow44:short, curWeight:cur?getComputedStyle(cur).fontWeight:null, curDeco:cur?getComputedStyle(cur).textDecorationLine:null,
  settingsGaps:lk.slice(1).map((r,i)=>Math.round(r.left-lk[i].right)), mainH:document.documentElement.scrollHeight};
})()`;
const out=[];
for (const p of ['policy-start-ko','policy-start-en','policy-privacy-ko','policy-privacy-en','policy-terms-ko','policy-terms-en','human-account','en-human-account']) for (const w of [1440,390,320]) {
  const t = await (await fetch('http://127.0.0.1:9334/json/new?about:blank',{method:'PUT'})).json();
  const ws = new WebSocket(t.webSocketDebuggerUrl); await new Promise(r=>ws.onopen=r);
  let id=0; const pend=new Map();
  ws.onmessage = e => { const m = JSON.parse(e.data); if(m.id&&pend.has(m.id)){pend.get(m.id)(m);pend.delete(m.id);} };
  const send=(method,params={})=>new Promise(r=>{const i=++id;pend.set(i,r);ws.send(JSON.stringify({id:i,method,params}));});
  await send('Page.enable');
  await send('Emulation.setDeviceMetricsOverride',{width:w,height:900,deviceScaleFactor:1,mobile:false});
  await send('Page.navigate',{url:pathToFileURL(path.join(srcDir,p+'.html')).href});
  await sleep(600);
  const m=(await send('Runtime.evaluate',{expression:probe,returnByValue:true})).result.result.value;
  const shot=await send('Page.captureScreenshot',{format:'png',captureBeyondViewport:true,clip:{x:0,y:0,width:w,height:m.mainH,scale:1}});
  fs.writeFileSync(path.join(outDir,`${p}-${w}.png`),Buffer.from(shot.result.data,'base64'));
  out.push({page:p,width:w,...m}); ws.close(); fetch('http://127.0.0.1:9334/json/close/'+t.id);
}
fs.writeFileSync(path.join(outDir,'measure.json'),JSON.stringify(out,null,1));
chrome.kill(); process.exit(0);
