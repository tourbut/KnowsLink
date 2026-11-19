import hashlib, json, os, re, sys, urllib.request, urllib.parse, urllib.error
RELAY = "http://127.0.0.1:18441"
def get(path, member=None, form=None):
    cookies = {}
    if member:
        tok = json.loads(urllib.request.urlopen(RELAY + "/fixture/login?member=" + member).read())["session"]
        cookies = {"__Host-kl_session": tok}
    req = urllib.request.Request(RELAY + path, data=urllib.parse.urlencode(form).encode() if form else None)
    if cookies: req.add_header("Cookie", "; ".join(f"{k}={v}" for k, v in cookies.items()))
    class NR(urllib.request.HTTPRedirectHandler):
        def redirect_request(self, *a, **k): return None
    try: r = urllib.request.build_opener(NR).open(req)
    except urllib.error.HTTPError as e: r = e
    return r.status, r.read().decode("utf-8", "replace")
if __name__ == "__main__":
    for use in sys.argv[1:]:
        st, body = get("/home/agents?use=" + use, "a")
        print(use, st, len(body))
        open(os.path.join(os.path.dirname(__file__), "work", f"screen-{use}.html"), "w", encoding="utf-8").write(body)
