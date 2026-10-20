"""Sent-pending cap (mine >= 10) counts held (blocked/receive-limit) rows like visible ones. Run once AFTER rerequest.py (A has 7 pending sent rows,
3 of them held from D/E/unknown path). Sends 4 synthetic non-existent codes: 3 must appear (7+3=10), the 4th must be dropped silently.
If held rows were not counted (5 counted), the 4th would appear. Then cancels the 3 added rows. Output cap-result.json."""
import json, re, sys
sys.path.insert(0, ".")
from client import as_, main, text
results = []
def chk(name, ok, detail=""):
    results.append({"name": name, "result": "PASS" if ok else "FAIL", "detail": detail if not ok else ""}); print("PASS" if ok else "FAIL", name, "" if ok else str(detail)[:400])
ROW = re.compile(r"<(?:li|article|section)[^>]*>(?:(?!</(?:li|article|section)>).)*?(?:보낸 친구 요청|친구 요청 기록)(?:(?!</(?:li|article|section)>).)*?</(?:li|article|section)>", re.S)
a = as_("a")
page = lambda: a.get("/home/friends")[2]
pend = lambda: [r for r in ROW.findall(main(page())) if "합성 회원 C의 친구 요청" not in r and "요청 대기" in text(r)]
chk("start_7_pending_sent_rows", len(pend()) == 7, len(pend()))
resp = []
for n in range(1, 5):
    s, h, d = a.post("/home/friends/request", {"code": f"qa-cap-fake-000{n}", "consent": "yes"}); resp.append((s, h.get("location")))
    print("sent", n, "pending now", len(pend()))
chk("four_posts_same_public_response", len(set(resp)) == 1, resp)
fin = pend()
tails = [f"000{n}" for n in range(1, 5)]
shown = [t for t in tails if any("…" + t in text(r) for r in fin)]
chk("pending_sent_capped_at_10_held_rows_counted", len(fin) == 10 and shown == tails[:3], (len(fin), shown))
for t in shown:  # cleanup: cancel the added rows
    pg = page()
    for f in re.findall(r'<form method="post" action="/home/friends/decision">.*?</form>', pg, re.S):
        if "…" + t in f:
            a.post("/home/friends/decision", {"friend": re.search(r'name="friend" value="([^"]+)"', f).group(1), "generation": re.search(r'name="generation" value="([^"]+)"', f).group(1), "decision": "cancel"})
chk("cleanup_back_to_7_pending", len(pend()) == 7, len(pend()))
json.dump({"results": results}, open("cap-result.json", "w", encoding="utf-8"), ensure_ascii=False, indent=1)
sys.exit(1 if any(r["result"] == "FAIL" for r in results) else 0)
