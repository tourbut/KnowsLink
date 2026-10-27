# CS-AC-07 (no host): default name is '내 에이전트 <번호>' (re-run of the last step of s_attempts, which lost its approver thread to a 429 wait).
from lib import *
box = {}
def fourth(url):
    b, cid, c = approve("c", url); box["d"] = c[0]
code, ev, err = cli(URL, W + "/at-nohost2", on_waiting=fourth)
rec("CS-AC-07", "no host: default 내 에이전트 <번호> on the consent screen and saved with it", (box.get("d") or "").startswith("내 에이전트 ") and ev[-1].get("name") == box.get("d"), f"default={box.get('d')} event={json.dumps(ev[-1], ensure_ascii=False)}")
