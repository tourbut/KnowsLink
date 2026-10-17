"""Observation only: anonymous (no session) GET /home/receipts volume, to see whether the anonymous IP limit (review low 2) is reachable via fixture."""
import sys, json, re; sys.path.insert(0, "."); from client import *
b = Browser(); seen = {}; first429 = None; body429 = None
for i in range(1, 121):
    s, h, x = b.get("/home/receipts")
    seen[s] = seen.get(s, 0) + 1
    if s == 429 and first429 is None:
        first429 = i; body429 = x
        break
out = {"requests": i, "status_counts": seen, "first_429_at": first429}
if body429:
    m = re.search(r"<main.*</main>", strip_style(body429), re.S).group(0)
    out["has_no_agents_text"] = "조회할 내 에이전트가 없습니다" in m
    out["has_agents_link"] = "/home/agents" in m
    out["alert"] = re.findall(r'role="alert"[^>]*>([^<]*)', m)
    out["main_text"] = re.sub(r"\s+", " ", re.sub(r"<[^>]+>", " ", m)).strip()[:400]
else:
    s, h, x = b.get("/home/receipts"); out["last_status"] = s; out["last_location"] = h.get("location")
json.dump(out, open("anon429-result.json", "w", encoding="utf-8"), ensure_ascii=False, indent=1); print(json.dumps(out, ensure_ascii=False, indent=1))
