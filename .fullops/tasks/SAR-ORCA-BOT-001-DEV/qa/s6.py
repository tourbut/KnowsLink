# S6 loopback round trip for SAR-ORCA-BOT-001-DEV: real relay handlers + Node adapter CLI + fake orca (argv log only).
import json, os, subprocess, time, urllib.parse, http.client
D = r"D:\workspace\KnowsLink\.devtools\SAR-ORCA-BOT-001-DEV"
PLUGIN = r"C:\Users\shin\orca\workspaces\KnowsLink\dev\adapters\dist\plugin.js"
NODE = r"D:\workspace\KnowsLink\.devtools\node-v22.22.2-win-x64\node.exe"
m = json.load(open(os.path.join(D, "manifest.json")))
origin, A = m["origin"], m["agents"]
GROK, ORCA = os.path.join(D, "agents", "grok"), os.path.join(D, "agents", "orca")
LOG = os.path.join(D, "fake-orca-%d.log" % int(time.time()))
FAIL = LOG + ".fail"
env = dict(os.environ, KNOWSLINK_ORCA_BIN=os.path.join(D, "fakeorca", "orca.exe"), FAKE_ORCA_LOG=LOG, KNOWSLINK_MODE="public-node")
out = []
def say(k, v): out.append({"step": k, "result": v}); print(k, "=>", v, flush=True)
def cli(*args):
    p = subprocess.run([NODE, PLUGIN, "orca", *args], capture_output=True, text=True, encoding="utf-8", env=env, timeout=60)
    return p.returncode, p.stdout.strip(), p.stderr.strip()
def post(path, form, cookie=True, bearer=None):
    u = urllib.parse.urlparse(origin)
    c = http.client.HTTPConnection(u.hostname, u.port)
    h = {"Content-Type": "application/x-www-form-urlencoded", "Origin": origin, "Sec-Fetch-Site": "same-origin"}
    if cookie: h["Cookie"] = f'{m["session_cookie"]}={m["session"]}'
    if bearer: h["Authorization"] = "Bearer " + bearer
    c.request("POST", path, urllib.parse.urlencode(form), h)
    r = c.getresponse(); r.read(); return r.status, r.getheader("Location")
def body(name, obj):
    p = os.path.join(D, name); open(p, "w", encoding="utf-8").write(json.dumps(obj, ensure_ascii=False)); return p
def lines():
    return [json.loads(l) for l in open(LOG, encoding="utf-8").read().splitlines()] if os.path.exists(LOG) else []
def wait_log(n, timeout=60):
    end = time.time() + timeout
    while time.time() < end:
        if len(lines()) >= n: return lines()
        time.sleep(1)
    return lines()
TASK_TEXT = "S6 비밀 본문: 쿠어에게 jevpet 확인 부탁"
task = body("task.json", {"text": TASK_TEXT, "project": "jevpet"})
say("S1 send task without grant", cli("send", GROK, "--to", A["orca"], "--kind", "orca.task.v1", "--file", task))
grok_cred = json.load(open(os.path.join(GROK, "agent.json")))["credential"]
say("S2 grant-create with agent bearer + confirmed, no session", post("/home/grant-create", {"from": A["grok"], "to": A["orca"], "kind": "orca.task.v1", "consent": "yes", "confirmed": "true"}, cookie=False, bearer=grok_cred))
for frm, to, kind in [("grok", "orca", "orca.task.v1"), ("orca", "grok", "orca.decision.v1"), ("grok", "orca", "orca.answer.v1")]:
    say(f"web grant {frm}->{to} {kind}", post("/home/grant-create", {"from": A[frm], "to": A[to], "kind": kind, "consent": "yes"}))
watch = subprocess.Popen([NODE, PLUGIN, "orca", "watch", ORCA, "--bot", "s6-fixture-bot"], stdout=subprocess.PIPE, stderr=subprocess.STDOUT, text=True, encoding="utf-8", env=env)
say("send task with grant", cli("send", GROK, "--to", A["orca"], "--kind", "orca.task.v1", "--file", task))
mail = wait_log(1); say("doorbell 1 argv", mail)
say("S3 doorbell has no body/sender", bool(mail) and all(TASK_TEXT not in " ".join(a) and A["grok"] not in " ".join(a) for a in mail))
say("S5 orca receive task", cli("receive", ORCA))
dec = body("decision.json", {"bot": "s6-fixture-bot", "item": 3, "question": "어느 브랜치로 머지할까요?", "options": ["main", "dev"]})
say("orca send decision", cli("send", ORCA, "--to", A["grok"], "--kind", "orca.decision.v1", "--file", dec))
say("grok receive decision", cli("receive", GROK))
ans = body("answer.json", {"item": 3, "answer": "1", "by": "owner"})
say("grok send answer", cli("send", GROK, "--to", A["orca"], "--kind", "orca.answer.v1", "--file", ans))
mail = wait_log(2); say("S4 doorbell 2 argv (answer --item)", mail[-1] if len(mail) >= 2 else None)
say("orca receive answer", cli("receive", ORCA))
open(FAIL, "w").write("1")
say("send task 2 (mail fails)", cli("send", GROK, "--to", A["orca"], "--kind", "orca.task.v1", "--file", task, "--key", "s6-task-two-key-0001"))
fails = wait_log(3); say("failed doorbell attempt", fails[-1] if len(fails) >= 3 else None)
open(FAIL, "w").write("0")
retry = wait_log(4, 90); say("S3 retried doorbell on later poll", retry[-1] if len(retry) >= 4 else None)
time.sleep(30)
say("S3 lines after 30 s more (no re-ring after success)", len(lines()))
say("orca receive task 2", cli("receive", ORCA))
say("S1 orca->grok task (no grant that direction)", cli("send", ORCA, "--to", A["grok"], "--kind", "orca.task.v1", "--file", task))
watch.terminate(); wout = watch.communicate(timeout=10)[0]
say("S5 watch stdout has no body", TASK_TEXT not in wout and "어느 브랜치" not in wout)
open(os.path.join(D, "watch.out"), "w", encoding="utf-8").write(wout)
json.dump(out, open(os.path.join(D, "s6-result.json"), "w", encoding="utf-8"), ensure_ascii=False, indent=1)
