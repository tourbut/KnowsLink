"""5-minute re-auth window: a session older than 5 minutes must not create an agent (POST /home/agents refused, nothing created).
Run (takes ~5.5 min): python -I -X utf8 stale.py   -> stale-result.json. Synthetic A only. One refused POST; no state change expected."""
import json, re, sys, time
sys.path.insert(0, ".")
from client import as_, main, text

def attempt_of(m):
    return re.search(r'<form method="post" action="/home/agents"><input type="hidden" name="attempt" value="([^"]+)"', m).group(1)

a = as_("a"); t0 = time.time()
s, h, d = a.get("/home/agents"); m = main(d)
cards_before = len(re.findall(r'<li id="agent-[^"]+"', m)); att = attempt_of(m)
while time.time() - t0 < 310: time.sleep(5)
s2, h2, d2 = a.post("/home/agents", {"attempt": att, "name": "stale window probe"})
s3, h3, d3 = a.get("/home/agents")
fresh = as_("a"); s4, h4, d4 = fresh.get("/home/agents")
out = {"age_s": round(time.time() - t0), "cards_before": cards_before, "post_status": s2, "post_location": h2.get("location"), "post_text": text(main(d2))[:300],
       "stale_get_status": s3, "stale_get_cards": len(re.findall(r'<li id="agent-[^"]+"', main(d3))), "fresh_cards": len(re.findall(r'<li id="agent-[^"]+"', main(d4))),
       "probe_name_present_anywhere": "stale window probe" in d3 + d4}
out["o_ct_04_wording"] = {"names_reauth_button_below": "아래 ‘로그인 다시 확인’" in out["post_text"], "retry_where_started": "이 작업을 시작한 화면에서 다시 시도하세요" in out["post_text"], "no_home_confirm_wording": "홈에서" not in out["post_text"],
    "has_reauth_button": "로그인 다시 확인" in d2 and "<button" in d2 or 'href="' in d2 and "로그인 다시 확인" in d2, "has_return_link": "내 에이전트로 돌아가기" in text(main(d2)) and 'href="/home/agents"' in d2}
out["pass"] = s2 not in (200, 201, 303) or ("재확인" in out["post_text"] or "다시 로그인" in out["post_text"]) and out["fresh_cards"] == cards_before and not out["probe_name_present_anywhere"]
json.dump(out, open("stale-result.json", "w", encoding="utf-8"), ensure_ascii=False, indent=1); print(json.dumps(out, ensure_ascii=False))
sys.exit(0 if out["fresh_cards"] == cards_before and not out["probe_name_present_anywhere"] and all(out["o_ct_04_wording"].values()) else 1)
