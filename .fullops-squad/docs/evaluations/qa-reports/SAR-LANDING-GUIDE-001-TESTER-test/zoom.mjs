// Use only the task-owned Chrome profile's real Page zoom setting; never substitute CSS zoom or pinch.
import fs from 'node:fs';
import path from 'node:path';
import {fileURLToPath} from 'node:url';
const out=path.dirname(fileURLToPath(import.meta.url));
const target=(await (await fetch('http://127.0.0.1:18891/json/list')).json()).find(t=>t.type==='page');
const ws=new WebSocket(target.webSocketDebuggerUrl);await new Promise((resolve,reject)=>{ws.onopen=resolve;ws.onerror=reject;});
let id=0;const pending=new Map();
ws.onmessage=e=>{const m=JSON.parse(e.data);if(m.id){const p=pending.get(m.id);pending.delete(m.id);m.error?p.reject(m.error):p.resolve(m.result);}};
const send=(method,params={})=>new Promise((resolve,reject)=>{pending.set(++id,{resolve,reject});ws.send(JSON.stringify({id,method,params}));});
const js=async expression=>{const r=await send('Runtime.evaluate',{expression,returnByValue:true});if(r.exceptionDetails)throw new Error('evaluation failed');return r.result.value;};
const find="(()=>{const roots=[document];for(let i=0;i<roots.length;i++){for(const e of roots[i].querySelectorAll('*')){if(e.shadowRoot)roots.push(e.shadowRoot);if(e.tagName==='SELECT'&&[...e.options].some(o=>o.textContent.trim()==='200%'))return e;}}return null;})()";
try{
  if(process.argv[2]==='inspect'){
    await send('Emulation.clearDeviceMetricsOverride');
    await send('Page.navigate',{url:'chrome://settings/appearance'});
    for(let i=0;i<100;i++){if(await js(`!!${find}`))break;await new Promise(r=>setTimeout(r,100));}
    const state=await js(`(()=>{const e=${find};return e?{id:e.id,index:e.selectedIndex,options:[...e.options].map(o=>({text:o.textContent.trim(),value:o.value}))}:null;})()`);
    fs.writeFileSync(path.join(out,'zoom-settings-before.json'),JSON.stringify(state,null,2));console.log(JSON.stringify(state));
  }else if(process.argv[2]==='expanded'){
    await send('Page.navigate',{url:'http://127.0.0.1:18782/auth/verify'});
    for(let i=0;i<50;i++){if(await js("!!document.querySelector('#start input')"))break;await new Promise(r=>setTimeout(r,100));}
    await js("document.querySelectorAll('details').forEach(e=>e.open=true);document.getElementById('start').scrollIntoView()");
    const result=await js("({width:innerWidth,overflow:document.documentElement.scrollWidth,dpr:devicePixelRatio,scale:visualViewport.scale,cssZoom:getComputedStyle(document.body).zoom,details:[...document.querySelectorAll('details')].every(e=>e.open),form:!!document.querySelector('#start input'),formWidth:document.querySelector('#start input').getBoundingClientRect().width,footer:!!document.querySelector('footer')})");
    result.passed=result.dpr===2&&result.scale===1&&result.cssZoom==='1'&&result.overflow<=result.width&&result.details&&result.form&&result.formWidth<result.width&&result.footer;
    fs.writeFileSync(path.join(out,'zoom-expanded-results.json'),JSON.stringify(result,null,2));
    console.log(JSON.stringify(result));if(!result.passed)process.exitCode=1;
  }else{
    const before=await js(`(()=>{const e=${find};return {index:e.selectedIndex,target:[...e.options].findIndex(o=>o.textContent.trim()==='200%')};})()`);
    if(before.target<0)throw new Error('200% option absent');
    await js(`${find}.focus()`);
    const key=before.index<before.target?'ArrowDown':'ArrowUp';const count=Math.abs(before.target-before.index);
    for(let i=0;i<count;i++){await send('Input.dispatchKeyEvent',{type:'keyDown',key,code:key,windowsVirtualKeyCode:key==='ArrowDown'?40:38});await send('Input.dispatchKeyEvent',{type:'keyUp',key,code:key,windowsVirtualKeyCode:key==='ArrowDown'?40:38});}
    const selected=await js(`${find}.selectedOptions[0].textContent.trim()`);
    await send('Page.navigate',{url:'http://127.0.0.1:18782/'});
    for(let i=0;i<50;i++){if(await js("!!document.querySelector('.hero')"))break;await new Promise(r=>setTimeout(r,100));}
    const result=await js('({width:innerWidth,height:innerHeight,dpr:devicePixelRatio,scale:visualViewport.scale,cssZoom:getComputedStyle(document.body).zoom,overflow:document.documentElement.scrollWidth})');
    result.setting=selected;result.passed=selected==='200%'&&result.cssZoom==='1'&&result.scale===1&&result.dpr===2&&result.overflow<=result.width;
    fs.writeFileSync(path.join(out,'zoom-results.json'),JSON.stringify(result,null,2));
    const {data}=await send('Page.captureScreenshot',{format:'png'});fs.writeFileSync(path.join(out,'browser-zoom-200.png'),Buffer.from(data,'base64'));
    console.log(JSON.stringify(result));if(!result.passed)process.exitCode=1;
  }
}finally{ws.close();}
