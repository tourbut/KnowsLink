// Measure rendered colors, including the representative authenticated member CTA.
import fs from 'node:fs';
import path from 'node:path';
import {fileURLToPath} from 'node:url';
const target=(await (await fetch('http://127.0.0.1:18891/json/list')).json()).find(t=>t.type==='page');
const ws=new WebSocket(target.webSocketDebuggerUrl);
await new Promise((resolve,reject)=>{ws.onopen=resolve;ws.onerror=reject;});
let id=0;const pending=new Map();
ws.onmessage=e=>{const m=JSON.parse(e.data);if(m.id){const p=pending.get(m.id);pending.delete(m.id);m.error?p.reject(m.error):p.resolve(m.result);}};
const send=(method,params={})=>new Promise((resolve,reject)=>{pending.set(++id,{resolve,reject});ws.send(JSON.stringify({id,method,params}));});
const evaluate=async expression=>(await send('Runtime.evaluate',{expression,returnByValue:true})).result.value;
function luminance(rgb){const v=rgb.match(/[\d.]+/g).slice(0,3).map(Number).map(x=>x/255).map(x=>x<=0.04045?x/12.92:((x+0.055)/1.055)**2.4);return v[0]*0.2126+v[1]*0.7152+v[2]*0.0722;}
const results=[];
try{
  for(const [url,selector] of [['http://127.0.0.1:18782/','.hero .primary'],['http://127.0.0.1:18783/home','.conversation-list .primary']]){
    await send('Page.navigate',{url});
    for(let i=0;i<50;i++){if(await evaluate(`location.href===${JSON.stringify(url)}&&!!document.querySelector(${JSON.stringify(selector)})`))break;await new Promise(r=>setTimeout(r,100));}
    const value=await evaluate(`(()=>{const e=document.querySelector(${JSON.stringify(selector)});const s=getComputedStyle(e);return {foreground:s.color,background:s.backgroundColor,font:s.fontSize,weight:s.fontWeight,outline:s.outlineWidth}})()`);
    const a=luminance(value.foreground),b=luminance(value.background),ratio=(Math.max(a,b)+0.05)/(Math.min(a,b)+0.05);
    results.push({url,selector,...value,ratio,minimum:4.5,passed:ratio>=4.5});
  }
  fs.writeFileSync(path.join(path.dirname(fileURLToPath(import.meta.url)),'contrast-results.json'),JSON.stringify(results,null,2));
  console.log(JSON.stringify(results));
  if(results.some(r=>!r.passed))process.exitCode=1;
}finally{ws.close();}
