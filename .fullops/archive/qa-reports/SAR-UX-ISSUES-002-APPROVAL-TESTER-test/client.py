"""Minimal cookie-aware client for the names fixture (loopback only, synthetic A/B/C)."""
import http.client, re, time, urllib.parse
ORIGIN = "http://127.0.0.1:18931"
HOST, PORT = "127.0.0.1", 18931
STYLE = re.compile(r"<style>.*?</style>", re.S)

class Browser:
    def __init__(self):
        self.cookies = {}
    waited429 = 0  # class-wide count of member-window 429 retries, reported as evidence
    def req(self, method, path, form=None, origin=ORIGIN, patient=True):
        for _ in range(12):
            r = self._req(method, path, form, origin)
            # every /home page shares a per-member request window; wait it out unless the caller is testing the 429 itself
            if r[0] != 429 or not patient or path.startswith("/fixture") or "요청이 많습니다" not in r[2]: return r
            Browser.waited429 += 1; time.sleep(8)
        return r
    def _req(self, method, path, form, origin):
        c = http.client.HTTPConnection(HOST, PORT, timeout=20)
        h = {"Origin": origin}
        if self.cookies:
            h["Cookie"] = "; ".join(f"{k}={v}" for k, v in self.cookies.items())
        body = None
        if form is not None:
            body = urllib.parse.urlencode(form, encoding="utf-8"); h["Content-Type"] = "application/x-www-form-urlencoded"
        c.request(method, path, body=body, headers=h)
        r = c.getresponse(); data = r.read().decode("utf-8", "replace")
        for k, v in r.getheaders():
            if k.lower() == "set-cookie":
                name, _, rest = v.partition("="); val = rest.split(";")[0]
                if "Max-Age=0" in v or val == "": self.cookies.pop(name, None)
                else: self.cookies[name] = val
        c.close()
        return r.status, dict((k.lower(), v) for k, v in r.getheaders()), data
    def get(self, path): return self.req("GET", path)
    def post(self, path, form=None, **kw): return self.req("POST", path, form if form is not None else None, **kw)

def strip(b): return re.sub(r"<svg.*?</svg>", "<svg/>", STYLE.sub("", b), flags=re.S)
def main(b):
    m = re.search(r"<main.*</main>", strip(b), re.S); return m.group(0) if m else strip(b)
def as_(who):
    br = Browser(); s = br.post(f"/fixture/as/{who}")[0]; assert s == 204, s; return br
def text(b): return re.sub(r"\s+", " ", re.sub(r"<[^>]+>", " ", b)).strip()
