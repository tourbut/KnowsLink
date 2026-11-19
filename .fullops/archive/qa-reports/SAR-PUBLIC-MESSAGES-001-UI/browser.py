"""Private Chrome/CDP capture helper; credentials stay outside persisted UI evidence."""

import base64
import json
import subprocess
import tempfile
import time
from pathlib import Path

import websocket


class Browser:
    def __init__(self):
        self.profile = tempfile.TemporaryDirectory(prefix="sar-messages-ui-chrome-")
        self.chrome = subprocess.Popen(
            ["/usr/bin/google-chrome", "--headless", "--no-sandbox", "--disable-gpu",
             "--remote-debugging-port=0", "--remote-allow-origins=*",
             "--user-data-dir=" + self.profile.name, "about:blank"],
            stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL,
        )
        active = Path(self.profile.name) / "DevToolsActivePort"
        for _ in range(100):
            if active.exists():
                break
            time.sleep(.05)
        lines = active.read_text().splitlines()
        self.ws = websocket.create_connection("ws://127.0.0.1:" + lines[0] + lines[1], timeout=15)
        self.counter = 0
        self.responses = []
        self.sessions = []

    def call(self, method, params=None, session=None):
        self.counter += 1
        request = {"id": self.counter, "method": method, "params": params or {}}
        if session:
            request["sessionId"] = session
        self.ws.send(json.dumps(request))
        while True:
            result = json.loads(self.ws.recv())
            if result.get('method') == 'Network.responseReceived':
                response = result['params']['response']
                self.responses.append({'url': response['url'].split('?')[0], 'status': response['status']})
            if result.get("id") == self.counter:
                if "error" in result:
                    raise RuntimeError(result["error"])
                return result.get("result", {})

    def new(self):
        context = self.call("Target.createBrowserContext")['browserContextId']
        target = self.call("Target.createTarget", {"url": "about:blank", "browserContextId": context})['targetId']
        session = self.call("Target.attachToTarget", {"targetId": target, "flatten": True})['sessionId']
        self.call("Page.enable", session=session)
        self.call("Network.enable", session=session)
        self.sessions.append((context, session))
        return session

    def evaluate(self, session, expression):
        result = self.call("Runtime.evaluate", {"expression": expression, "returnByValue": True,
                                                "awaitPromise": True}, session)
        if "exceptionDetails" in result:
            raise RuntimeError("browser expression failed")
        return result['result'].get('value')

    def navigate(self, session, url):
        self.call("Page.navigate", {"url": url}, session)
        time.sleep(.15)
        for _ in range(100):
            if self.evaluate(session, "document.readyState") == "complete":
                break
            time.sleep(.03)

    def viewport(self, session, width):
        self.call("Emulation.setDeviceMetricsOverride", {"width": width, "height": 844 if width == 390 else 900,
                  "deviceScaleFactor": 1, "mobile": width == 390}, session)

    def capture(self, session, directory, name, width):
        self.viewport(session, width)
        # Hide identity and private inputs before both the text snapshot and the screenshot.
        self.evaluate(session, """(() => {
          for(const input of document.querySelectorAll('input')) {
            if(['email','code','csrf','token'].includes(input.name) || /^[A-Za-z0-9_-]{43}$/.test(input.value)) {
              input.style.visibility='hidden';
            }
          }
          for(const node of document.querySelectorAll('p,dd')) {
            if(node.textContent.includes('@')) { node.textContent='[비공개 신원 가림]'; }
          }
        })()""")
        metrics = self.evaluate(session, """({viewport:innerWidth,scroll:document.documentElement.scrollWidth,
          scripts:document.querySelectorAll('script').length,injected:window.uiInjected===true,
          approve:!!document.querySelector('button[value="approve"]'),
          text:document.body.innerText})""")
        size = self.call("Page.getLayoutMetrics", session=session)['cssContentSize']
        data = self.call("Page.captureScreenshot", {"format": "png", "captureBeyondViewport": True,
             "clip": {"x": 0, "y": 0, "width": width, "height": size['height'], "scale": 1}}, session)
        Path(directory, name + ".png").write_bytes(base64.b64decode(data['data']))
        Path(directory, name + ".txt").write_text(metrics.pop('text') + "\n")
        return metrics

    def close(self):
        for context, _ in self.sessions:
            self.call("Target.disposeBrowserContext", {"browserContextId": context})
        self.ws.close()
        self.chrome.terminate()
        self.chrome.wait(timeout=10)
        self.profile.cleanup()
