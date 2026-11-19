"""Irreversible creation checks on the last free agent slot of A (run once, after every read-only check and after stale.py).
Run: python -I -X utf8 mutate.py -> mutate-result.json. Synthetic A/B only. Exactly one agent is created; refusals create nothing."""
import json, re, sys, threading
sys.path.insert(0, ".")
from client import Browser, as_, main, text

results = []
def chk(name, ok, detail=""):
    results.append({"name": name, "result": "PASS" if ok else "FAIL", "detail": detail if not ok else ""})
    print("PASS" if ok else "FAIL", name, "" if ok else str(detail)[:400])

CARD = re.compile(r'<li id="agent-([^"]+)"')
def agents(b):
    s, h, d = b.get("/home/agents"); m = main(d); return s, m, CARD.findall(m)

a = as_("a")
s, m, ids0 = agents(a)
att = re.search(r'name="attempt" value="([^"]+)"', m).group(1)
legacy0 = text(a.get("/home/legacy")[2])
reg_card0 = text(re.search(r'<li id="agent-agent_fx_registered">.*?</li>', m, re.S).group(0))
chk("pre_four_agents_one_free_slot", len(ids0) == 4 and "활성 4/5" in legacy0, (ids0, re.findall(r"활성 \d/\d", legacy0)))
# double click: two simultaneous submits of the same form (same attempt)
res = [None, None]
def post(i):
    br = Browser(); br.cookies = dict(a.cookies)
    res[i] = br.req("POST", "/home/agents", {"attempt": att, "name": "Connect QA 신규 에이전트"}, patient=False)
th = [threading.Thread(target=post, args=(i,)) for i in range(2)]; [t.start() for t in th]; [t.join() for t in th]
locs = [r[1].get("location") for r in res]
chk("double_click_both_submits_succeed_to_same_location", all(r[0] == 303 for r in res) and locs[0] == locs[1] and locs[0], [(r[0], r[1].get("location")) for r in res])
s, m, ids1 = agents(a)
new = [i for i in ids1 if i not in ids0]
chk("double_click_created_exactly_one_agent", len(ids1) == 5 and len(new) == 1, (ids0, ids1))
# refresh / back-resubmit of the same form later
r3 = a.post("/home/agents", {"attempt": att, "name": "Connect QA 신규 에이전트"}, patient=False)
s, m, ids2 = agents(a)
chk("later_resubmit_same_attempt_same_location_no_second_agent", r3[0] == 303 and r3[1].get("location") == locs[0] and len(ids2) == 5, (r3[0], r3[1].get("location"), len(ids2)))
# a different attempt past the 5-agent limit is refused, nothing created
r4 = a.post("/home/agents", {"attempt": "second-attempt-token-qa-0000000000000000000", "name": "Connect QA 여섯째"}, patient=False)
s, m, ids3 = agents(a)
chk("sixth_agent_with_new_attempt_refused_nothing_created", r4[0] != 303 and len(ids3) == 5 and "여섯째" not in m, (r4[0], r4[1].get("location"), text(r4[2])[:160], len(ids3)))
# existing key untouched, new card is waiting-not-registered
reg_card1 = text(re.search(r'<li id="agent-agent_fx_registered">.*?</li>', m, re.S).group(0))
chk("existing_registered_key_card_unchanged", reg_card0 == reg_card1, (reg_card0[:200], reg_card1[:200]))
newcard = text(re.search(r'<li id="agent-%s">.*?</li>' % re.escape(new[0]), m, re.S).group(0)) if new else ""
chk("new_agent_card_not_claimed_connected", new and "Connect QA 신규 에이전트" in newcard and "✓ 등록됨" not in newcard and "키 등록됨" not in newcard, newcard[:300])
bb = as_("b"); sb, mb, idb = agents(bb)
chk("new_agent_not_visible_to_other_member", not new or (new[0] not in mb and "Connect QA" not in mb), idb)
json.dump({"origin": "http://127.0.0.1:18932", "location": locs[0], "new_agent": new[0] if new else None, "attempt_prefix": att[:6], "results": results}, open("mutate-result.json", "w", encoding="utf-8"), ensure_ascii=False, indent=1)
sys.exit(1 if any(r["result"] == "FAIL" for r in results) else 0)
