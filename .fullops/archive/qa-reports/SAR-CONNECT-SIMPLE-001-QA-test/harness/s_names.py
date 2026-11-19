# CS-AC-07 (web add path default + numbering) and CS-AC-08 (nameless existing agent: name form visible, no fold) on member E.
from lib import *
E = Browser().login("e")
def default_web():
    s, raw, _ = E.get("/home/agents")
    m = re.search(r'<input id="new-agent-name"[^>]*value="([^"]*)"', raw)
    return html.unescape(m.group(1)) if m else None, field(raw, "attempt")
d1, att = default_web()
s, raw, h = E.post("/home/agents", {"attempt": att, "name": d1})
d2, att = default_web()
rec("CS-AC-07", "web add path: default name filled in; after creating it the next default takes the next number", bool(d1) and d1 != d2 and d1.startswith("내 에이전트 ") and d2.startswith("내 에이전트 "), f"first={d1} second={d2}")
rec("CS-AC-07", "Dots remote path default name 'Dots': fixture cannot start the remote OAuth consent; DEV evidence only", "not_run", "no remote OAuth client in fixture; see DEV report TestDefaultNames + remote tests")
aid = fx("/fixture/nameless?m=e", "POST")[1]
s, raw, _ = E.get("/home/agents")
body = raw
i = body.find("이름을 정할 에이전트")
sec = body[i:i + 8000] if i >= 0 else ""
form = re.search(r"<form[^>]*>(?:(?!</form>)[\s\S])*?" + aid + r"[\s\S]*?</form>", sec)
folded = "<details" in sec.split("</form>")[0] if form else True
rec("CS-AC-08", "nameless agent: '이름을 정할 에이전트' section is first, with the name input and save button not inside a <details>", i >= 0 and bool(form) and "이름 저장" in sec and not folded, (form.group(0)[:500] if form else sec[:300]))
fa = re.search(r'action="([^"]+)"', form.group(0)).group(1) if form else ""
rec("CS-AC-08", "section order: 이름 필요 section comes before the agent list (M-01)", i >= 0 and i < body.find("내 에이전트와 연결 상태"), f"pos name-section={i} list={body.find('내 에이전트와 연결 상태')}")
fields = dict(re.findall(r'<input type="hidden" name="([^"]+)" value="([^"]*)"', form.group(0))) if form else {}
fields["name"] = "바로 정한 이름"
s, raw2, h = E.post(fa, fields)
named = [a for a in agents_active("e") if a["id"] == aid]
rec("CS-AC-08", "saving the name from the visible form names the agent without opening anything else", s in (200, 303) and named and named[0]["name"] == "바로 정한 이름", f"status={s} action={fa} agent={named}")
s, raw3, _ = E.get("/home/agents")
rec("CS-AC-08", "after naming, the '이름을 정할 에이전트' section is gone (empty section shows no title)", "이름을 정할 에이전트" not in raw3 and "바로 정한 이름" in raw3, "")
