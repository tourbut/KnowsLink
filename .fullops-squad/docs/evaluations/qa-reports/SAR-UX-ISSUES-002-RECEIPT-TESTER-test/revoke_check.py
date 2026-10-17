"""Post-revoke state check. The fixture revoked agent_a during the first full run (irreversible), so a fresh A session
gets a forged recovery context {Agent: agent_a, ID: AB} instead of replaying the pre-revoke lookup."""
import base64, json, re, sys
sys.path.insert(0, "."); from client import *

AB = "0199a3f2-4c10-7a11-8b22-334455667788"
res = []
def chk(name, fn):
    try: fn(); res.append({"check": name, "result": "pass"})
    except Exception as e: res.append({"check": name, "result": "fail", "detail": f"{type(e).__name__}: {e}"})

b = Browser(); b.post("/fixture/recover"); b.post("/fixture/as/a"); b.get("/home/receipts")
ck = json.loads(base64.urlsafe_b64decode(b.cookies["__Host-kl_receipt"] + "=="))
ck.update(Agent="agent_a", ID=AB)
b.cookies["__Host-kl_receipt"] = base64.urlsafe_b64encode(json.dumps(ck, separators=(",", ":")).encode()).decode().rstrip("=")
opts = lambda x: re.findall(r'<option value="([^"]*)"', x)

def recovery():
    s, h, x = b.get("/home/receipts")
    assert s == 403 and opts(x) == ["", "agent_a2", "agent_a"], (s, opts(x))
    assert re.search(r'<option value="agent_a" selected disabled>[^<]*/ 철회</option>', x) and '<option value="agent_a2" selected' not in x
    assert "선택한 내 에이전트가 철회되었습니다. 사용할 내 에이전트를 다시 선택하세요." in x and 'aria-invalid="true"' in x
    assert re.search(r'name="id" value="%s"' % AB, x)
chk("revoked_agent_recovery_context_is_not_silently_replaced", recovery)

def url():
    s, h, x = b.get(q(agent="agent_a", id=AB))
    assert s == 403 and '<option value="agent_a2" selected' not in x and "철회되었습니다" in x
    s, h, x = b.get(q(agent="agent_a2", id=AB)); assert s == 403 and "이 요청의 상태를 확인할 수 없습니다." in x
chk("revoked_agent_in_url_and_remaining_agent_still_refused_for_nonparticipant", url)

def other():
    c = Browser(); c.post("/fixture/as/b"); c.cookies["__Host-kl_receipt"] = b.cookies["__Host-kl_receipt"]
    s, h, x = c.get("/home/receipts")
    assert s == 200 and opts(x) == ["", "agent_b"] and "agent_a" not in x and AB not in x
chk("account_B_never_sees_revoked_A_context", other)

json.dump(res, open("revoke-result.json", "w", encoding="utf-8"), ensure_ascii=False, indent=1)
for r in res: print(r["result"].upper(), r["check"], r.get("detail", ""))
sys.exit(0 if all(r["result"] == "pass" for r in res) else 1)
