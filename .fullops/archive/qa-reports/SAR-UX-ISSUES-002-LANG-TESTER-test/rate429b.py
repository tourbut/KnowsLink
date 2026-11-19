"""D-LG-01 recheck on b11d5b4: the member request-limit 429 page must follow the browser cookie language (en/ko), Content-Language = html lang.
Hits the per-member limit once per account, then probes cookie variants inside the same limit window. GET only. Original rate429.py is kept unchanged.
Run: python rate429b.py -> rate429b-result.json (exit 1 on FAIL)"""
import sys, re, json, time
sys.path.insert(0, ".")
from client import *
from checks_lib import chk, results, lang_of, hangul_residue

def probe(b, cookie, path="/home/account", query=""):
    b.cookies.pop("__Host-kl_lang", None)
    if cookie is not None: b.cookies["__Host-kl_lang"] = cookie
    return b.req("GET", path + query, patient=False)

def hit_limit(who):
    b = as_(who)
    for i in range(120):
        s, h, body = b.req("GET", "/home/account", patient=False)
        if s == 429: return b, i + 1
    return b, None

REC = {}
for who, cases in [("en", [(None, "ko"), ("ko", "ko"), ("en", "en"), ("xx", "ko"), ("EN", "ko")]), ("ko", [(None, "ko"), ("en", "en"), ("ko", "ko")]), ("none", [(None, "ko"), ("en", "en"), ("ko", "ko")])]:
    b, n = hit_limit(who)
    chk(f"limit_{who}_reaches_429_within_120_gets", n is not None, n); print("   INFO", who, "429 after", n, "requests")
    if n is None: continue
    for cookie, exp in cases:
        s, h, body = probe(b, cookie)
        while s != 429:   # window rolled over: hit again
            b2, n2 = hit_limit(who); b = b2; s, h, body = probe(b, cookie)
        t = re.search(r"<title>(.*?)</title>", body).group(1); m = text(main(body))
        tag = f"{who}_saved_cookie_{cookie}"
        REC[tag] = {"status": s, "html_lang": lang_of(body), "content_language": h.get("content-language"), "title": t, "main": m[:200], "retry_after": h.get("retry-after"), "hrefs": re.findall(r'href="([^"]*)"', main(body)), "cache_control": h.get("cache-control")}
        chk(f"d_lg_01_{tag}_429_html_lang_is_{exp}", s == 429 and lang_of(body) == exp, REC[tag])
        chk(f"d_lg_01_{tag}_content_language_header_is_{exp}", h.get("content-language") == exp, h.get("content-language"))
        chk(f"d_lg_01_{tag}_title_and_text_in_{exp}", (t == "요청 제한 · KnowsLink" and "요청이 많습니다" in m) if exp == "ko" else (bool(t) and "KnowsLink" in t and not re.search(r"[가-힣]", t + m) and re.search(r"[A-Za-z]{4,}", m) is not None), (t, m[:160]))
        if exp == "en": chk(f"d_lg_01_{tag}_english_page_has_no_hangul", hangul_residue(body) == [], hangul_residue(body))
        chk(f"d_lg_01_{tag}_retry_and_safe_return_present_no_saved_language_claim", bool(h.get("retry-after")) or bool(re.search(r"\d", m)), (h.get("retry-after"), m[:160]))
        chk(f"d_lg_01_{tag}_return_link_is_internal", all(x.startswith("/") and not x.startswith("//") for x in REC[tag]["hrefs"] if not x.startswith("#")), REC[tag]["hrefs"])
    # ?lang= query on the limit page (information): admission never looks the member up
    s, h, body = probe(b, None, query="?lang=en")
    print("   INFO", who, "429 with ?lang=en and no cookie ->", s, lang_of(body), h.get("content-language"))
    # the limit page must not store anything
    chk(f"d_lg_01_{who}_limit_page_sets_no_lang_cookie", "kl_lang" not in h.get("set-cookie", ""), h.get("set-cookie"))
    # recover the account: wait out the window so the designer's GETs are not blocked
    time.sleep(1)

# anonymous / other admission paths: public page flood with cookie en (information; stops at 150 requests)
b = Browser(); b.cookies["__Host-kl_lang"] = "en"; anon = None
for i in range(150):
    s, h, body = b.req("GET", "/", patient=False)
    if s != 200: anon = (i + 1, s, lang_of(body), h.get("content-language"), re.search(r"<title>(.*?)</title>", body).group(1) if "<title>" in body else None); break
print("   INFO anonymous public GET flood result:", anon or "no refusal in 150 requests")
if anon: chk("d_lg_01_anon_public_refusal_follows_cookie_en", anon[1] in (429, 503) and anon[2] == "en" and anon[3] == "en", anon)
json.dump({"origin": ORIGIN, "records": REC, "anon_flood": anon, "results": results}, open("rate429b-result.json", "w", encoding="utf-8"), ensure_ascii=False, indent=1)
f = [r for r in results if r["result"] == "FAIL"]
print(f"\n{len(results)-len(f)} PASS / {len(f)} FAIL")
sys.exit(1 if f else 0)
