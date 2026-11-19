# CS-AC-08 follow-up: E still holds nameless agents seeded by earlier partial runs; name them all and the empty section must show no title (M-01).
from lib import *
E = Browser().login("e")
s, raw, _ = E.get("/home/agents")
ids = re.findall(r'<form method="post" action="/home/agent-name"><input type="hidden" name="agent" value="(agent_[^"]+)"><label for="unnamed-name', raw)
for i in ids: E.post("/home/agent-name", {"agent": i, "name": "정리 " + i[-4:]})
s, raw2, _ = E.get("/home/agents")
rec("CS-AC-08", f"after naming all {len(ids)} remaining nameless agents (inline form, one by one) the '이름을 정할 에이전트' section has no title at all (M-01)", len(ids) >= 1 and "이름을 정할 에이전트" not in raw2 and 'id="unnamed"' not in raw2, f"named={len(ids)}")
