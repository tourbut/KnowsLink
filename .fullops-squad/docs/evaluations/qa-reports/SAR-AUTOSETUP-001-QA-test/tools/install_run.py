"""SAR-AUTOSETUP-001-QA: run the on-screen setup command in an isolated box and play the browser person (A=claude, B=codex).
usage: python install_run.py <box> <app> <lang> <member> <answers comma list> [--cancel-first]"""
import json, os, re, subprocess, sys, threading, time, urllib.request, urllib.parse
from pathlib import Path
sys.path.insert(0, os.path.dirname(__file__))
from qa import get, RELAY
D = Path(__file__).parent
box, app, lang, member, answers = sys.argv[1:6]
root = D / "work" / box
for s in ("local", "claude", "codex", "tmp"): (root / s).mkdir(parents=True, exist_ok=True)
env = dict(os.environ, LOCALAPPDATA=str(root / "local"), CLAUDE_CONFIG_DIR=str(root / "claude"), CODEX_HOME=str(root / "codex"), TEMP=str(root / "tmp"), TMP=str(root / "tmp"))
shown = get(f"/home/agents?use={app}&lang={lang}", member)[1]
import html
cmd = html.unescape(re.search(r'<textarea id="setup-command"[^>]*>(.*?)</textarea>', shown, re.S).group(1))
cmd = cmd.replace("https://link.knowslog.com", RELAY)  # the only edit: production URL -> isolated relay
print("COMMAND:", cmd, flush=True)
events = []
def browser(cid):
    time.sleep(1.5)
    t0 = time.time()
    post = lambda path, form: urllib.request.urlopen(urllib.request.Request(RELAY + path, data=urllib.parse.urlencode(form).encode(), method="POST")).read()
    post("/fixture/google", {"connection": cid, "member": member})  # stand-in for the Google OAuth callback
    st, page = get(f"/connect/{cid}", member)
    fp = re.search(r'aria-label="공개키 지문">([^<]+)<', page)
    ev = {"status": st, "page_fingerprint": fp.group(1) if fp else None, "has_name_form": "device-agent-name" in page}
    name = f"QA {app} 에이전트"
    # real browser submit: the page's own form (session cookie of the member)
    import http.client
    tok = json.loads(urllib.request.urlopen(RELAY + "/fixture/login?member=" + member).read())["session"]
    rq = urllib.request.Request(RELAY + "/home/device-confirm", data=urllib.parse.urlencode({"connection": cid, "name": name}).encode(), method="POST")
    rq.add_header("Cookie", "__Host-kl_session=" + tok)
    class NR(urllib.request.HTTPRedirectHandler):
        def redirect_request(self, *a, **k): return None
    try: r = urllib.request.build_opener(NR).open(rq); code = r.status
    except urllib.error.HTTPError as e: code = e.code
    ev["confirm_status"] = code; ev["name"] = name; events.append(ev)
command = f"[Console]::OutputEncoding=[Text.Encoding]::UTF8; {cmd}; exit $LASTEXITCODE"
started = time.time()
p = subprocess.Popen(["powershell.exe", "-NoProfile", "-Command", command], stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.STDOUT, env=env)
p.stdin.write((answers.replace(",", "\n") + "\n").encode()); p.stdin.close()
lines = []
for raw in p.stdout:
    line = raw.decode("utf-8", "replace").rstrip()
    lines.append(f"[{time.time()-started:6.1f}s] {line}")
    m = re.search(r"/connect/([A-Za-z0-9_-]{43})", line)
    if m: threading.Thread(target=browser, args=(m.group(1),), daemon=True).start()
code = p.wait(); secs = time.time() - started
time.sleep(2)
out = "\n".join(lines)
(D / "work" / f"{box}-{time.strftime('%H%M%S')}.log").write_text(f"exit={code} seconds={secs:.1f}\nbrowser={json.dumps(events, ensure_ascii=False)}\n{out}\n", encoding="utf-8")
print(out); print(f"EXIT={code} SECONDS={secs:.1f}"); print("BROWSER", json.dumps(events, ensure_ascii=False))
