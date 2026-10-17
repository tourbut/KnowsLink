"""Minimal cookie-aware client for the receipt fixture (loopback only, synthetic A/B/C)."""
import http.client, re, urllib.parse
ORIGIN = "http://127.0.0.1:18928"
HOST, PORT = "127.0.0.1", 18928
STYLE = re.compile(r"<style>.*?</style>", re.S)

class Browser:
    def __init__(self):
        self.cookies = {}
    def req(self, method, path):
        c = http.client.HTTPConnection(HOST, PORT, timeout=15)
        h = {"Origin": ORIGIN}
        if self.cookies:
            h["Cookie"] = "; ".join(f"{k}={v}" for k, v in self.cookies.items())
        c.request(method, path, headers=h)
        r = c.getresponse(); body = r.read().decode("utf-8", "replace")
        for k, v in r.getheaders():
            if k.lower() == "set-cookie":
                name, _, rest = v.partition("=")
                val = rest.split(";")[0]
                if "Max-Age=0" in v or val == "":
                    self.cookies.pop(name, None)
                else:
                    self.cookies[name] = val
        c.close()
        return r.status, dict((k.lower(), v) for k, v in r.getheaders()), body
    def get(self, path): return self.req("GET", path)
    def post(self, path): return self.req("POST", path)

def strip_style(b): return STYLE.sub("", b)
def q(**kw): return "/home/receipts?" + urllib.parse.urlencode(kw)
