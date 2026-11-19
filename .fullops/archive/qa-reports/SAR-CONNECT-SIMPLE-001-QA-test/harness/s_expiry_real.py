# Real 10-minute expiry: nobody opens the link; the shipped CLI must end with reason expired (KO) and no agent.
from lib import *
W = "D:/workspace/KnowsLink/.devtools/connect-simple-001-qa/work/expiry-real"
t0 = time.time()
code, ev, err = cli(URL, W, "--host", "codex", timeout=900)
print(json.dumps({"exit": code, "events": ev, "stderr": err, "seconds": round(time.time() - t0)}, ensure_ascii=False))
last = ev[-1] if ev else {}
rec("CS-AC-11", "expired (real 10-min wait, KO): reason code, text, next action", code == 2 and last.get("reason") == "expired" and "다시 연결" in last.get("message", ""), json.dumps(last, ensure_ascii=False))
