// Independent CDP layout and synthetic login checks using installed Chrome and Node WebSocket.
import fs from 'node:fs';
import path from 'node:path';
import {fileURLToPath} from 'node:url';
import {createHash} from 'node:crypto';

const out = path.dirname(fileURLToPath(import.meta.url));
const version = await (await fetch('http://127.0.0.1:18891/json/version')).json();
const target = (await (await fetch('http://127.0.0.1:18891/json/list')).json()).find(t=>t.type==='page');
const ws = new WebSocket(target.webSocketDebuggerUrl);
await new Promise((resolve,reject)=>{ws.onopen=resolve;ws.onerror=reject;});
let next=0, actions=0;
const started=Date.now(), pending=new Map(), checks=[], captures=[], requests=[];
const save=()=>fs.writeFileSync(path.join(out,'browser-results.json'),JSON.stringify({browser:version.Browser,protocol:version['Protocol-Version'],target:target.id,source:'4795d619513fde4fbcd4a55b94ec91b056ebbaf0',actions,seconds:(Date.now()-started)/1000,checks,captures,requests},null,2));
ws.onmessage=event=>{const m=JSON.parse(event.data);if(m.id){const p=pending.get(m.id);pending.delete(m.id);m.error?p.reject(new Error(m.error.message)):p.resolve(m.result);}else if(m.method==='Network.requestWillBeSent'){const u=new URL(m.params.request.url);requests.push({method:m.params.request.method,origin:u.origin,path:u.pathname});}};
function send(method,params={}){return new Promise((resolve,reject)=>{const id=++next;pending.set(id,{resolve,reject});ws.send(JSON.stringify({id,method,params}));});}
async function js(expression){const r=await send('Runtime.evaluate',{expression,returnByValue:true,awaitPromise:true});if(r.exceptionDetails)throw new Error('Browser evaluation failed');return r.result.value;}
function check(name,value){checks.push({name,passed:value===true,value});save();if(value!==true)throw new Error(name);}
function action(){if(++actions>60||Date.now()-started>1200000)throw new Error('scenario budget exceeded');}
async function go(url){action();const done=new Promise(resolve=>{const previous=ws.onmessage;ws.onmessage=e=>{previous(e);if(JSON.parse(e.data).method==='Page.loadEventFired'){ws.onmessage=previous;resolve();}};});await send('Page.navigate',{url});await done;}
async function shot(name){const metrics=await js('({width:innerWidth,height:innerHeight,dpr:devicePixelRatio,cssZoom:getComputedStyle(document.body).zoom})');const {data}=await send('Page.captureScreenshot',{format:'png'});const bytes=Buffer.from(data,'base64');fs.writeFileSync(path.join(out,name+'.png'),bytes);captures.push({file:name+'.png',...metrics,browserZoom:1,sha256:createHash('sha256').update(bytes).digest('hex')});save();}
async function key(key,code,virtual){action();await send('Input.dispatchKeyEvent',{type:'keyDown',key,code,windowsVirtualKeyCode:virtual,text:key==='Enter'?'\r':undefined});await send('Input.dispatchKeyEvent',{type:'keyUp',key,code,windowsVirtualKeyCode:virtual});}
await send('Page.enable');await send('Network.enable');
try {
  for(const width of [1440,390,320]){
    await send('Emulation.setDeviceMetricsOverride',{width,height:width===1440?900:844,deviceScaleFactor:1,mobile:false});
    await go('http://127.0.0.1:18782/');
    check('viewport-'+width,await js(`innerWidth===${width}`));
    check('hero-in-viewport-'+width,await js("document.querySelector('.hero .primary').getBoundingClientRect().bottom<=innerHeight"));
    await shot('hero-'+width);
    await js("document.querySelectorAll('details').forEach(e=>e.open=true)");
    check('open-details-overflow-'+width,await js('document.documentElement.scrollWidth<=innerWidth'));
    check('targets-44px-'+width,await js("[...document.querySelectorAll('a,button,input,summary')].filter(e=>e.getClientRects().length).every(e=>e.getBoundingClientRect().height>=44)"));
    check('url-string-'+width,await js("[...document.querySelectorAll('code')].some(e=>e.textContent==='https://link.knowslog.com/mcp')"));
    if(width===320){await js("document.getElementById('install').scrollIntoView()");await shot('install-open-320');}
    if(width===390){await js("document.getElementById('examples').scrollIntoView()");await shot('examples-390');}
  }
  await send('Emulation.setEmulatedMedia',{features:[{name:'prefers-reduced-motion',value:'reduce'}]});
  check('reduced-motion-no-animation',await js("matchMedia('(prefers-reduced-motion:reduce)').matches && [...document.querySelectorAll('body *')].every(e=>getComputedStyle(e).animationName==='none'&&getComputedStyle(e).transitionDuration==='0s')"));
  await go('http://127.0.0.1:18784/');await shot('none-320');
  await go('http://127.0.0.1:18781/auth/google/callback?error=access_denied');
  check('google-error-label',await js("!!document.querySelector('[role=alert]') && !!document.querySelector('a[href=\"#start\"]')"));
  await shot('google-error-320');
  await go('http://127.0.0.1:18783/');
  await js("document.querySelector('input[name=email]').focus()");
  await send('Input.insertText',{text:'landing@example.test'});action();
  await key('Enter','Enter',13);
  for(let i=0;i<50 && await js('location.pathname')!=='/auth/verify';i++)await new Promise(r=>setTimeout(r,100));
  check('browser-email-enter',await js("location.pathname==='/auth/verify'"));
  const code=await (await fetch('http://127.0.0.1:18783/fixture/code')).text();
  if(!/^\d{6}$/.test(code))throw new Error('synthetic code absent');
  await js("document.querySelector('input[name=code]').focus()");
  await send('Input.insertText',{text:code});action();await key('Enter','Enter',13);
  for(let i=0;i<50 && await js('location.pathname')==='/auth/verify';i++)await new Promise(r=>setTimeout(r,100));
  check('browser-verified-home',await js("location.pathname==='/home'"));
  await send('Emulation.setDeviceMetricsOverride',{width:1440,height:900,deviceScaleFactor:1,mobile:false});
  await shot('member-home-1440');
  check('member-css-isolation',await js("!document.body.classList.contains('landing') && getComputedStyle(document.querySelector('h1')).fontSize!=='40px'"));
  await go('http://127.0.0.1:18783/home/account');await shot('member-account-1440');
  await js("[...document.querySelectorAll('a')].find(e=>e.textContent==='서비스 안내').click()");action();
  for(let i=0;i<50 && await js('location.pathname')!=='/';i++)await new Promise(r=>setTimeout(r,100));
  check('account-to-public-authenticated',await js("location.pathname==='/' && document.body.innerText.includes('로그인되어 있습니다.') && !document.body.innerText.includes('비공개 QA 이름') && !document.body.innerText.includes('landing@example.test')"));
  check('no-external-browser-requests',requests.every(r=>['http://127.0.0.1:18781','http://127.0.0.1:18782','http://127.0.0.1:18783','http://127.0.0.1:18784'].includes(r.origin)));
  console.log(`Browser checks passed: ${checks.length}; actions: ${actions}`);
} finally {save();ws.close();}
