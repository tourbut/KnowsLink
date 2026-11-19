# CS-AC-10: retry in a folder with leftover attempt files, a pre-made empty folder and a fresh folder, on C: and D:; no cleanup by the user; moved folder not deleted.
from lib import *
import glob
ROOTS = {"D:": W + "/folders", "C:": "C:/Users/shin/orca/workspaces/KnowsLink/tester/.fullops-squad/local/SAR-CONNECT-SIMPLE-001-QA/work-c7"}
def acl(path):
    ps = "[Console]::OutputEncoding=[Text.Encoding]::UTF8; $a=Get-Acl -LiteralPath '" + path + "'; $a.Owner; $a.AreAccessRulesProtected; ($a.Access | ForEach-Object { $_.IdentityReference.Value + '|' + $_.FileSystemRights + '|' + $_.AccessControlType }) -join ';'"
    return subprocess.run(["pwsh", "-NoProfile", "-Command", ps], capture_output=True, text=True, encoding="utf-8", errors="replace").stdout.strip().splitlines()
me = os.environ.get("USERNAME", "")
def acl_ok(path):
    lines = acl(path)
    owner, prot, rules = lines[0], lines[1], lines[2] if len(lines) > 2 else ""
    ids = {r.split("|")[0].lower() for r in rules.split(";") if r}
    allowed = lambda i: me.lower() in i or i in ("nt authority\\system", "builtin\\administrators")
    return me.lower() in owner.lower() and all(allowed(i) for i in ids), f"owner={owner} protected={prot} rules={rules}"
def ok(url): approve("d", url)
def cancel(url):
    b = Browser().login("d"); cid = conn_id(url); google_login(b, cid, "d"); b.post("/home/cancel", {"connection": cid})
for drive, root in ROOTS.items():
    os.makedirs(root, exist_ok=True)
    # F1: leftover attempt files from a cancelled attempt
    f1 = f"{root}/leftover"
    code, ev, err = cli(URL, f1, "--host", "claude", on_waiting=cancel)
    leftover = sorted(os.listdir(f1))
    code, ev, err = cli(URL, f1, "--host", "claude", on_waiting=ok)
    last = ev[-1]
    moved = glob.glob(f1 + ".unused-*")
    kept = sorted(os.listdir(moved[0])) if moved else []
    rec("CS-AC-10", f"{drive} leftover attempt files: retry connects with no manual cleanup", code == 0 and last.get("state") == "connected" and os.path.exists(f1 + "/agent.json"), json.dumps(last, ensure_ascii=False))
    rec("CS-AC-10", f"{drive} leftover files were moved to .unused-<time>, not deleted", bool(moved) and "login.json" in kept and "private.pem" in kept, f"leftover_before={leftover} moved={[os.path.basename(m) for m in moved]} kept={kept}")
    ok1, d1 = acl_ok(f1); rec("CS-AC-10", f"{drive} new private folder ACL: owner current user, only user/SYSTEM/Administrators", ok1, d1[:400])
    # F2: waiting attempt in folder resumes without a new start
    f2 = f"{root}/resume"
    r = hnode("start", URL, f2); n0 = json.loads(fx("/fixture/hits")[1])["starts"]
    code, ev, err = cli(URL, f2, "--host", "claude", on_waiting=ok)
    n1 = json.loads(fx("/fixture/hits")[1])["starts"]
    rec("CS-AC-10", f"{drive} waiting attempt left in the folder is resumed (same link, no new start)", code == 0 and ev[0].get("resumed") == "true" and ev[0].get("url") == r["url"] and n1 == n0, f"first_event={ev[0]} starts {n0}->{n1}")
    # F3: pre-made empty folder
    f3 = f"{root}/empty"; os.makedirs(f3, exist_ok=True)
    code, ev, err = cli(URL, f3, "--host", "claude", on_waiting=ok)
    rec("CS-AC-10", f"{drive} pre-made empty folder is used as is (connected, nothing moved)", code == 0 and ev[-1].get("state") == "connected" and not glob.glob(f3 + ".unused-*") and os.path.exists(f3 + "/agent.json"), json.dumps(ev[-1], ensure_ascii=False))
    ok3, d3 = acl_ok(f3); rec("CS-AC-10", f"{drive} pre-made empty folder ACL after connect: owner current user, only user/SYSTEM/Administrators", ok3, d3[:400])
    # F4: fresh folder
    f4 = f"{root}/fresh"
    code, ev, err = cli(URL, f4, "--host", "claude", on_waiting=ok)
    rec("CS-AC-10", f"{drive} brand-new folder connects", code == 0 and ev[-1].get("state") == "connected", json.dumps(ev[-1], ensure_ascii=False))
    # F5: connected folder is reused as the connected agent (R-01/AS-10)
    code, ev, err = cli(URL, f4, "--host", "claude")
    rec("CS-AC-10", f"{drive} folder with agent.json is reused as connected (no new attempt)", code == 0 and ev[-1].get("existing") == "true", json.dumps(ev[-1], ensure_ascii=False))
    # F6: unusable path (parent is a file) -> private_folder (KO/EN), nothing deleted
    blocker = f"{root}/blocker-file"; open(blocker, "w").write("x")
    code, ev, err = cli(URL, blocker + "/child", "--host", "claude")
    last = ev[-1]
    rec("CS-AC-10", f"{drive} path whose parent is a file -> reason private_folder with the path and next action (KO)", code == 2 and last.get("reason") == "private_folder" and "child" in last.get("message", "") and "기본 위치" in last.get("message", ""), json.dumps(last, ensure_ascii=False))
    code, ev, err = cli(URL, blocker + "/child", "--host", "claude", "--lang", "en")
    last = ev[-1]
    rec("CS-AC-11", f"{drive} private_folder (EN): code, text, next action", code == 2 and last.get("reason") == "private_folder" and "default location" in last.get("message", ""), json.dumps(last, ensure_ascii=False))
