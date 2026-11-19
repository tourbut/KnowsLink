"""Known-limit probe: the shared admission 429 page follows the browser cookie, not the member's saved language (DEV report). GET only."""
import sys, re, json
sys.path.insert(0, ".")
from client import *
b = as_("ko"); b.cookies["__Host-kl_lang"] = "en"
hit = None
for i in range(60):
    s, h, body = b.req("GET", "/home/account", patient=False)
    if s == 429: hit = (i, s, re.search(r'<html lang="(\w+)"', body).group(1), h.get("content-language"), re.search(r"<title>(.*?)</title>", body).group(1), text(main(body))[:160]); break
print("first 429:", hit)
json.dump({"first_429": hit, "saved_language": "ko", "cookie": "en"}, open("rate429-result.json", "w", encoding="utf-8"), ensure_ascii=False, indent=1)
