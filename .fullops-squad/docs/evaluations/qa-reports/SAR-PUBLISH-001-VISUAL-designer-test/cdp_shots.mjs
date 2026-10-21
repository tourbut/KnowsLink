// Capture + measure fixture pages at true CSS viewport widths via CDP.
import { spawn } from 'node:child_process';
import fs from 'node:fs';
import path from 'node:path';
import { pathToFileURL } from 'node:url';
const [,, srcDir, outDir] = process.argv;
fs.mkdirSync(outDir, { recursive: true });
const chrome = spawn('C:/Program Files/Google/Chrome/Application/chrome.exe', ['--headless=new','--disable-gpu','--remote-debugging-port=9333','--user-data-dir='+path.join(outDir,'..','prof'),'--hide-scrollbars','about:blank'], {stdio:'ignore'});
const sleep = ms => new Promise(r => setTimeout(r, ms));
let ver; for (let i=0;i<50;i++){ try{ ver = await (await fetch('http://127.0.0.1:9333/json/version')).json(); break;}catch{await sleep(200);} }
const pages = ['policy-privacy-ko','policy-privacy-en','policy-terms-ko','policy-terms-en','policy-start-ko','policy-start-en','human-account','en-human-account'];
const widths = [1440,390,320];
const results = [];
async function session(){
  const t = await (await fetch('http://127.0.0.1:9333/json/new?about:blank',{method:'PUT'})).json();
  const ws = new WebSocket(t.webSocketDebuggerUrl); await new Promise(r=>ws.onopen=r);
  let id=0; const pend=new Map();
  ws.onmessage = e => { const m = JSON.parse(e.data); if(m.id&&pend.has(m.id)){pend.get(m.id)(m);pend.delete(m.id);} };
  const send=(method,params={})=>new Promise(r=>{const i=++id;pend.set(i,r);ws.send(JSON.stringify({id:i,method,params}));});
  return {send, close:()=>{ws.close(); fetch('http://127.0.0.1:9333/json/close/'+t.id);}};
}
const probe = `(() => {
 const d=document, de=d.documentElement, cs=getComputedStyle, b=d.body;
 const vis=[...d.querySelectorAll('a')].filter(a=>a.offsetParent!==null);
 const foot=d.querySelector('footer');
 const wide=[...d.querySelectorAll('body *')].filter(e=>e.getBoundingClientRect().right>innerWidth+1 && !e.closest('table,.table-wrap,[style*="overflow"]')).slice(0,5).map(e=>e.tagName+'.'+e.className);
 const tabs=[...d.querySelectorAll('table')].map(t=>({cw:(t.parentElement.clientWidth),sw:t.scrollWidth, pov:cs(t.parentElement).overflowX}));
 const mono=[...d.querySelectorAll('body *')].filter(e=>/mono|courier|consolas/i.test(cs(e).fontFamily)).length;
 const grad=[...d.querySelectorAll('*')].filter(e=>/gradient/.test(cs(e).backgroundImage)).length;
 const pill=[...d.querySelectorAll('a,button')].filter(e=>e.offsetParent&&parseFloat(cs(e).borderRadius)>=20 && e.getBoundingClientRect().height<80).map(e=>e.textContent.trim().slice(0,20));
 const em=(d.body.innerText.match(/\\p{Extended_Pictographic}/gu)||[]).length;
 const ko=(d.body.innerText.match(/[\u3131-\u318E\uAC00-\uD7A3]+/g)||[]);
 return {title:d.title, lang:de.lang, h1:d.querySelectorAll('h1').length, current:[...d.querySelectorAll('[aria-current]')].map(a=>a.textContent.trim()+'|'+a.getAttribute('href')), iw:innerWidth, sw:de.scrollWidth, bodySw:b.scrollWidth, hOverflow:de.scrollWidth>innerWidth, wide, tabs, mono, grad, pill, emoji:em, bg:cs(b).backgroundColor, font:cs(b).fontFamily, color:cs(b).color, footerLinks: foot?[...foot.querySelectorAll('a')].map(a=>a.textContent.trim()+'|'+a.getAttribute('href')):null, koText:ko.slice(0,12), hasDetails:d.querySelectorAll('details').length, scripts:d.querySelectorAll('script').length, imgs:d.querySelectorAll('img').length, h1text:(d.querySelector('h1')||{}).textContent, mainH: de.scrollHeight};
})()`;
for (const p of pages) for (const w of widths) {
  const s = await session();
  await s.send('Page.enable');
  await s.send('Emulation.setDeviceMetricsOverride',{width:w,height:900,deviceScaleFactor:1,mobile:false});
  await s.send('Page.navigate',{url:pathToFileURL(path.join(srcDir,p+'.html')).href});
  await sleep(700);
  const r = await s.send('Runtime.evaluate',{expression:probe,returnByValue:true});
  if(!r.result.result.value){console.log(JSON.stringify(r.result).slice(0,500));process.exit(1);} const m = r.result.result.value;
  const shot = await s.send('Page.captureScreenshot',{format:'png',captureBeyondViewport:true,clip:{x:0,y:0,width:w,height:m.mainH,scale:1}});
  fs.writeFileSync(path.join(outDir,`${p}-${w}.png`),Buffer.from(shot.result.data,'base64'));
  results.push({page:p,width:w,...m});
  s.close();
}
fs.writeFileSync(path.join(outDir,'..','measure.json'),JSON.stringify(results,null,1));
chrome.kill();
process.exit(0);
