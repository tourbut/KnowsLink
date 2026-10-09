# Measure exact Go fixture templates in a private headless Chrome profile; retain only public layout evidence.
import base64,json,subprocess,tempfile,time
from pathlib import Path
import websocket
profile=tempfile.TemporaryDirectory(prefix='knowslink-messages-chrome-')
chrome=subprocess.Popen(['/usr/bin/google-chrome','--headless','--no-sandbox','--disable-gpu','--remote-debugging-port=0','--remote-allow-origins=*','--user-data-dir='+profile.name,'about:blank'],stdout=subprocess.DEVNULL,stderr=subprocess.DEVNULL)
ws=None
try:
 active=Path(profile.name)/'DevToolsActivePort'
 for _ in range(100):
  if active.exists():break
  time.sleep(.05)
 lines=active.read_text().splitlines();ws=websocket.create_connection('ws://127.0.0.1:'+lines[0]+lines[1],timeout=10)
 counter=0
 def call(method,params=None,session=None):
  global counter
  counter+=1;request={'id':counter,'method':method,'params':params or {}}
  if session:request['sessionId']=session
  ws.send(json.dumps(request))
  while True:
   result=json.loads(ws.recv())
   if result.get('id')==counter:
    if 'error' in result:raise RuntimeError(result['error'])
    return result.get('result',{})
 target=call('Target.createTarget',{'url':'about:blank'})['targetId']
 session=call('Target.attachToTarget',{'targetId':target,'flatten':True})['sessionId']
 call('Page.enable',session=session)
 call('Emulation.setDeviceMetricsOverride',{'width':390,'height':844,'deviceScaleFactor':1,'mobile':True},session)
 results=[]
 for path in sorted(Path('/tmp/knowslink-messages-ui').glob('*.html')):
  call('Page.navigate',{'url':path.as_uri()},session)
  for _ in range(100):
   r=call('Runtime.evaluate',{'expression':'document.readyState','returnByValue':True},session)
   if r['result'].get('value')=='complete':break
   time.sleep(.02)
  width=call('Runtime.evaluate',{'expression':'JSON.stringify({viewport:innerWidth,scroll:document.documentElement.scrollWidth,body:document.body.scrollWidth})','returnByValue':True},session)
  metrics=json.loads(width['result']['value']);metrics['page']=path.name;results.append(metrics)
  if metrics['viewport']!=390 or metrics['scroll']>390:raise RuntimeError(metrics)
  if path.name=='receipt-failed-expired.html':
   capture=call('Page.captureScreenshot',{'format':'png','captureBeyondViewport':False},session)
   Path('/tmp/knowslink-messages-ui/receipt-mobile.png').write_bytes(base64.b64decode(capture['data']))
 Path('/tmp/knowslink-messages-ui/mobile-width.json').write_text(json.dumps({'viewport':'390x844','metrics':results,'visualAcceptance':'designer pending'},indent=2)+'\n')
 print(json.dumps(results))
finally:
 if ws:ws.close()
 chrome.terminate();chrome.wait(timeout=10);profile.cleanup()
