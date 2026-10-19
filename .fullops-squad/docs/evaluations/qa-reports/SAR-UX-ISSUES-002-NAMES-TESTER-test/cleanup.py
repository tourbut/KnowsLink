"""Fixture hygiene after mutate.py/followup.py: C cancels every placeholder 'sent' row except the one real request to B (so the designer captures are readable).
Each cancel is the product's normal 요청 취소 action; the result is also a check (a placeholder row can be cancelled and disappears)."""
import html as H, json, os, re, sys
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__))); from client import *
B, C = as_("b"), as_("c")
rows = lambda: re.findall(r'<li><h3>보낸 친구 요청</h3>.*?name="friend" value="([^"]+)"><input type="hidden" name="generation" value="(\d+)">', main(C.get("/home/friends")[2]), re.S)
recv_b = re.findall(r'<li><h3>(.*?)의 친구 요청</h3>.*?name="friend" value="([^"]+)"><input type="hidden" name="generation" value="(\d+)"', main(B.get("/home/friends")[2]), re.S)
assert len(recv_b) == 1 and recv_b[0][0].startswith("합성 회원 c"), recv_b
keep = recv_b[0][1]; before = rows(); assert keep in [r[0] for r in before]
for fid, gen in before:
    if fid == keep: continue
    s, h, b = C.post("/home/friends/decision", {"friend": fid, "generation": gen, "decision": "cancel"}); assert s in (200, 303), (s, fid)
after = rows(); res = {"before": len(before), "after": len(after), "kept_is_b_request": [r[0] for r in after] == [keep], "b_still_has_pending": len(re.findall(r"의 친구 요청</h3>", main(B.get("/home/friends")[2]))) == 1, "member_window_429_retries": Browser.waited429}
json.dump(res, open(os.path.join(os.path.dirname(os.path.abspath(__file__)), "cleanup-result.json"), "w"), indent=1); print(res)
sys.exit(0 if res["kept_is_b_request"] and res["b_still_has_pending"] else 1)
