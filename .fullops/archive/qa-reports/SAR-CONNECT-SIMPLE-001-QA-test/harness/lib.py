import json, os, re, subprocess, sys, threading, time, urllib.parse, html
sys.path.insert(0, os.path.dirname(__file__))
from br import *
NODE = "D:/workspace/KnowsLink/.devtools/node-v22.22.2-win-x64/node.exe"
PLUGIN = "C:/Users/shin/orca/workspaces/KnowsLink/tester/adapters/dist/plugin.js"
URL = "http://" + BASE
OUT = os.path.dirname(os.path.abspath(__file__))
RESULTS = os.path.join(OUT, os.environ.get("QA_RESULTS", "results.jsonl"))
def rec(id, check, ok, detail=""):
    r = {"id": id, "check": check, "result": ok if ok in ("not_run", "INFO") else ("PASS" if ok else "FAIL"), "detail": str(detail)[:600]}
    with open(RESULTS, "a", encoding="utf-8") as f: f.write(json.dumps(r, ensure_ascii=False) + "\n")
    print(f"[{r['result']}] {id} {check} :: {r['detail'][:200]}", flush=True)
    return ok
def agents_active(x): return [a for a in (state(x)["agents"] or []) if not a["revoked"]]
def attempts(x): return state(x)["attempts"] or []
def field(raw, name):
    m = re.search(r'<input[^>]*name="%s"[^>]*value="([^"]*)"' % name, raw) or re.search(r'<input[^>]*value="([^"]*)"[^>]*name="%s"' % name, raw)
    return html.unescape(m.group(1)) if m else None
def conn_id(url): return url.rsplit("/", 1)[1]
def google_login(b, cid, acct):
    """What the person does in the browser: open the link, choose a Google account, return. Returns (status, body) of the callback."""
    s, raw, _ = b.get(f"/connect/{cid}")
    s, raw, h = b.post("/auth/google", {"connection": cid})
    assert s == 303, (s, raw[:200])
    q = urllib.parse.parse_qs(urllib.parse.urlparse(h["Location"]).query)
    # the account chooser page registers the PKCE challenge for the codes it offers
    pick = Browser(); pick.req("GET", urllib.parse.urlparse(h["Location"]).path + "?" + urllib.parse.urlparse(h["Location"]).query)
    code = f"{acct}.{q['nonce'][0]}"
    return b.get("/auth/google/callback?" + urllib.parse.urlencode({"code": code, "state": q["state"][0]}))
def consent(b, cid, name=None):
    """Consent page values and the explicit confirm. Returns (default_name, status, body)."""
    s, raw, _ = b.get(f"/connect/{cid}")
    default = field(raw, "name")
    s2, raw2, h2 = b.post("/home/device-confirm", {"connection": cid, "name": default if name is None else name})
    return default, s2, raw2, raw
def approve(acct, url, name=None, b=None):
    b = b or Browser().login(acct)
    cid = conn_id(url)
    cb = google_login(b, cid, acct)
    return b, cid, consent(b, cid, name)
def _guard(fn):
    def run(*a):
        try: fn(*a)
        except Exception:
            import traceback
            with open(os.path.join(OUT, "thread-errors.log"), "a", encoding="utf-8") as f: f.write(traceback.format_exc() + chr(10))
    return run
def cli(base, folder, *args, env=None, on_waiting=None, timeout=180):
    """Run the shipped plugin.js setup-connect; call on_waiting(url) once the URL event arrives. Returns (exit, events, stderr)."""
    e = dict(os.environ); e.pop("KNOWSLINK_MODE", None)
    e.update(env or {})
    p = subprocess.Popen([NODE, PLUGIN, "setup-connect", base, folder, *args], stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True, encoding="utf-8", env=e)
    events, fired = [], False
    t0 = time.time()
    for line in p.stdout:
        try: ev = json.loads(line)
        except Exception: ev = {"raw": line}
        events.append(ev)
        if ev.get("url") and on_waiting and not fired:
            fired = True; threading.Thread(target=_guard(on_waiting), args=(ev["url"],), daemon=True).start()
        if time.time() - t0 > timeout: p.kill(); break
    code = p.wait(); err = p.stderr.read()
    with open(os.path.join(OUT, "cli-output.log"), "a", encoding="utf-8") as f:
        f.write(json.dumps({"base": base, "folder": folder, "args": list(args), "exit": code, "events": events, "stderr": err}, ensure_ascii=False) + chr(10))
    return code, events, err
def friends(a, b_):
    """Make a and b friends through the product forms."""
    ba, bb = Browser().login(a), Browser().login(b_)
    code_b = state(b_)["friendCode"]
    s, raw, _ = ba.get("/home/friends")
    s, raw, h = ba.post("/home/friends/request", {"code": code_b, "consent": "yes"})
    s, raw, _ = bb.get("/home/friends")
    return ba, bb, raw
def hnode(*args, env=None):
    e = dict(os.environ); e.pop("KNOWSLINK_MODE", None); e.update(env or {})
    args = list(args)
    for n in range(6):
        if n and args[0] == "start": args[2] = args[2].split("-r")[0] + f"-r{n}"   # a 429 after mkdir leaves the folder behind
        r = subprocess.run([NODE, os.path.join(OUT, "h.mjs"), *args], capture_output=True, text=True, encoding="utf-8", env=e)
        try: v = json.loads(r.stdout.strip().splitlines()[-1])
        except Exception: v = {"error": r.stdout + r.stderr}
        if "429" in str(v.get("error", "")): time.sleep(20); continue
        return v
    return v
W = "D:/workspace/KnowsLink/.devtools/connect-simple-001-qa/work/run8"
def accept(bb):
    s, raw, _ = bb.get("/home/friends")
    f = re.search(r'name="friend" value="([^"]+)"><input type="hidden" name="generation" value="(\d+)"', raw)
    return bb.post("/home/friends/decision", {"friend": f.group(1), "generation": f.group(2), "consent": "yes", "decision": "accept"})
