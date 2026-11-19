"""Minimal browser stand-in: manual cookie jar (relay cookies are Secure, loopback is http), no redirects followed."""
import http.client, json, re, time, urllib.parse, html
BASE = "127.0.0.1:18951"
class NoRedirect(Exception): pass
class Browser:
    def __init__(self): self.cookies = {}; self.lang = None
    def req(self, method, path, data=None, js=False, headers=None, cross=False):
        for _ in range(30):
            c = http.client.HTTPConnection(BASE, timeout=30)
            h = {"Sec-Fetch-Site": "cross-site" if cross else "same-origin", "Origin": "http://" + BASE}
            if self.cookies: h["Cookie"] = "; ".join(f"{k}={v}" for k, v in self.cookies.items())
            body = None
            if js: body = json.dumps(data); h["Content-Type"] = "application/json"
            elif data is not None: body = urllib.parse.urlencode(data); h["Content-Type"] = "application/x-www-form-urlencoded"
            h.update(headers or {})
            c.request(method, path, body, h); r = c.getresponse(); raw = r.read().decode("utf-8", "replace")
            for k, v in r.getheaders():
                if k.lower() == "set-cookie":
                    nv = v.split(";")[0]; n, _, val = nv.partition("=")
                    if "Max-Age=0" in v or "Max-Age=-1" in v or val == "": self.cookies.pop(n, None)
                    else: self.cookies[n] = val
            if r.status == 429 and not ("오늘 새 에이전트" in raw or "many new agents" in raw):
                time.sleep(min(float(r.getheader("Retry-After") or 5) + 0.5, 12)); continue
            return r.status, raw, dict(r.getheaders())
        raise RuntimeError("rate limited too long " + path)
    def get(self, p, **k): return self.req("GET", p, **k)
    def post(self, p, d=None, **k): return self.req("POST", p, d or {}, **k)
    def login(self, x): s, b, _ = self.req("POST", f"/fixture/as/{x}"); assert s == 200, (s, b); return self
    def text(self, raw): return html.unescape(re.sub(r"<[^>]+>", " ", re.sub(r"<(script|style)[\s\S]*?</\1>", "", raw)))
def fx(path, method="GET"):
    c = http.client.HTTPConnection(BASE, timeout=30); c.request(method, path); r = c.getresponse(); return r.status, r.read().decode()
def state(x): return json.loads(fx(f"/fixture/state?m={x}")[1])
