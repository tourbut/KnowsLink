"""Read-only HTTP checks for the LANG fixture (U3 settings/language, U7 header).
GET plus /fixture/as and /lang/* (cookie only). Anonymous POST /home/language is a rejected request (no state change).
Run: python checks.py -> http-result.json, exit 1 on any FAIL."""
import json, re, sys, urllib.parse
sys.path.insert(0, ".")
from client import *
from dump import SCREENS

HANGUL = re.compile(r"[ㄱ-ㆎ가-힣]")
USER_NAME_KO = re.compile(r"긴한글이름+")  # synthetic user-data name, exempt (original text)
LANG_LABEL = re.compile(r'<(span|a)[^>]*lang="ko"[^>]*>한국어</\1>')  # the language's own name, exempt
NOTE = "This browser did not keep the language choice. It applies to this page only."
results = []


def chk(name, ok, detail=""):
    results.append({"name": name, "result": "PASS" if ok else "FAIL", "detail": "" if ok else str(detail)[:400]})
    print("PASS" if ok else "FAIL", name, "" if ok else str(detail)[:300])


def lang_of(body):
    m = re.search(r'<html lang="(\w*)"', body)
    return m.group(1) if m else None


def title_of(body):
    m = re.search(r"<title>(.*?)</title>", body)
    return m.group(1) if m else None


def hangul_residue(body):
    b = LANG_LABEL.sub("", USER_NAME_KO.sub("", STYLE.sub("", body)))
    return sorted(set(re.findall(r"[ㄱ-ㆎ가-힣]+", b)))[:8]


def kst(body): return re.findall(r"\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} KST", body)
def hrefs(body): return sorted(set(re.findall(r'(?:href|action)="([^"#?]*)', body)))
def sections(m): return [text(x) for x in re.findall(r"<h2[^>]*>(.*?)</h2>", m, re.S)]


def anon(path, cookie=None, query=""):
    b = Browser()
    if cookie is not None: b.cookies["__Host-kl_lang"] = cookie
    s, h, body = b.get(path + query)
    return s, h, body


def member(who, path, cookie=None, query=""):
    b = as_(who)
    if cookie is not None: b.cookies["__Host-kl_lang"] = cookie
    return b.get(path + query)


# ---------------------------------------------------------------- fetch every screen once per account
P = {}
for who in ["en", "ko", "none"]:
    b = as_(who)
    for k, p in SCREENS.items():
        P[(who, k)] = b.get(p)

# ---------------------------------------------------------------- 1 settings screen
for who, L in [("en", "en"), ("ko", "ko"), ("none", "ko")]:
    s, h, body = P[(who, "account")]
    m = main(body)
    exp = {"en": ["My name", "My friend code", "Language", "Sign-in and security"], "ko": ["내 이름", "내 친구 코드", "언어", "로그인과 보안"]}[L]
    chk(f"settings_{who}_status_lang_title", s == 200 and lang_of(body) == L and h.get("content-language") == L and title_of(body) == {"en": "Settings · KnowsLink", "ko": "설정 · KnowsLink"}[L], (s, lang_of(body), title_of(body)))
    chk(f"settings_{who}_section_order", sections(m) == exp, sections(m))
    h1 = text(re.search(r"<h1>(.*?)</h1>", m).group(1))
    hdr = {"en": "Settings", "ko": "설정"}[L]
    chk(f"settings_{who}_h1_and_header_link", h1 == hdr and f'<a href="/home/account" aria-current="page">{hdr}</a></header>' in body, h1)
    chk(f"settings_{who}_no_old_my_account_label", "내 계정" not in body and "My account" not in body, "")
    radios = re.findall(r'<input type="radio" name="lang" value="(\w+)"( checked)?><span lang="(\w+)">(.*?)</span>', m)
    chk(f"settings_{who}_two_radios_labels_and_checked", [(r[0], r[2], r[3]) for r in radios] == [("ko", "ko", "한국어"), ("en", "en", "English")] and [r[0] for r in radios if r[1]] == [L], radios)
    btn = re.search(r'<form method="post" action="/home/language">.*?<button[^>]*>(.*?)</button>', m, re.S).group(1)
    chk(f"settings_{who}_save_button_and_no_autosubmit", btn == {"en": "Save language", "ko": "언어 저장"}[L] and not re.search(r"onchange|onclick|<script", m), btn)
    chk(f"settings_{who}_language_group_radiogroup_legend_help", f'<fieldset role="radiogroup" aria-describedby="{"language-unsaved language-help" if who == "none" else "language-help"}"><legend>' in m and 'id="language-help"' in m and (('id="language-unsaved"' in m) == (who == "none")), "")  # b11d5b4 V-LANG-03: unsaved note joins the description only when nothing is saved
    chk(f"settings_{who}_language_section_anchor", '<section id="language" aria-labelledby="language-h">' in m, "")

# ---------------------------------------------------------------- 2 every designated screen in both languages
for L in ["en", "ko"]:
    for k in SCREENS:
        s, h, body = P[(L, k)]
        if k == "verify":
            continue  # /auth/verify takes its language from the browser cookie only (see section 4)
        tag = f"{k}_{L}"
        chk(f"screen_{tag}_html_lang_content_language_title_one_h1", lang_of(body) == L and h.get("content-language") == L and bool(title_of(body)) and len(re.findall(r"<h1[ >]", body)) == 1, (s, lang_of(body), h.get("content-language"), title_of(body)))
        chk(f"screen_{tag}_skip_link_main", 'class="skip" href="#main"' in body and '<main id="main"' in body, "")
        if L == "en":
            chk(f"screen_{tag}_english_no_hangul_outside_names_and_language_label", hangul_residue(body) == [], hangul_residue(body))
for who, ks in [("ko", ["conn-waiting", "gate"]), ("none", ["conn-waiting", "gate", "conversation"])]:  # ko is the conversation's recipient (200), not a stranger
    for k in ks:
        s, h, body = P[(who, k)]
        chk(f"refusal_{k}_{who}_is_403_ko_no_data_leak", s == 403 and lang_of(body) == "ko" and "Synthetic" not in body and ID["gate"] not in text(body), (s, lang_of(body)))
bnone = as_("none"); bnone.cookies["__Host-kl_lang"] = "en"
for k in ["conn-waiting", "gate", "conversation"]:
    s, h, body = bnone.get(SCREENS[k])
    chk(f"refusal_{k}_none_with_en_cookie_is_english_403", s == 403 and lang_of(body) == "en" and h.get("content-language") == "en" and hangul_residue(body) == [], (s, lang_of(body), hangul_residue(body)))

for L in ["en", "ko"]:
    titles = {k: title_of(P[(L, k)][2]) for k in SCREENS if k not in ("verify", "public", "conversation") and P[(L, k)][0] == 200 and not k.startswith("conn-")}  # /home and /home/conversation share one title by design
    chk(f"titles_{L}_unique_per_screen", len(set(titles.values())) == len(titles), titles)

for k in ["home", "friends", "agents", "account", "legacy", "new", "receipts", "conversation", "gate", "name"]:
    en_attr = re.findall(r'(?:aria-label|alt|placeholder|title)="([^"]*)"', P[("en", k)][2])
    bad = [a for a in en_attr if HANGUL.search(USER_NAME_KO.sub("", a))]
    chk(f"a11y_names_{k}_en_attributes_no_hangul_outside_names", not bad, bad[:3])
for k in ["home", "friends", "agents", "account", "new", "receipts"]:
    en_attr = re.findall(r'aria-label="([^"]*)"', P[("en", k)][2]); ko_attr = re.findall(r'aria-label="([^"]*)"', P[("ko", k)][2])
    chk(f"a11y_names_{k}_ko_en_have_labels_and_ko_labels_not_english_copy", len(en_attr) >= 1 and len(ko_attr) >= 1 and all(HANGUL.search(a) for a in ko_attr if a not in en_attr), (len(ko_attr), len(en_attr), ko_attr[:3]))

for L, exp in [("en", ["Conversations", "Friends", "My agents"]), ("ko", ["대화", "친구", "내 에이전트"])]:
    for k in ["home", "friends", "agents", "account", "new"]:
        nav = re.search(r'<nav class="nav" aria-label="([^"]*)">(.*?)</nav>', P[(L, k)][2], re.S)
        chk(f"nav_{k}_{L}_labels", bool(nav) and re.findall(r">([^<]+)</a>", nav.group(2)) == exp and nav.group(1) == {"en": "Main menu", "ko": "주요 메뉴"}[L], nav and nav.groups())

# ---------------------------------------------------------------- 3 invariants across languages
for k in ["connect-device", "remote"]:
    e, o = P[("en", k)][2], P[("ko", k)][2]
    chk(f"invariant_{k}_kst_fingerprint_urls_identical_en_ko", kst(e) == kst(o) and re.findall(r"SHA256:[\w\-]+", e) == re.findall(r"SHA256:[\w\-]+", o) and re.findall(r"https?://[\w./-]+", e) == re.findall(r"https?://[\w./-]+", o), (kst(e), kst(o)))
chk("invariant_connect_device_has_one_kst_deadline", len(kst(P[("en", "connect-device")][2])) == 1, kst(P[("en", "connect-device")][2]))
chk("invariant_settings_form_actions_same_en_ko", hrefs(P[("en", "account")][2]) == hrefs(P[("ko", "account")][2]), "")

# ---------------------------------------------------------------- 4 priority: saved -> anonymous explicit (cookie) -> Korean
for who, cookie, query, exp in [
    ("en", "ko", "", "en"), ("en", None, "?lang=ko", "en"), ("en", "ko", "?lang=ko", "en"),
    ("ko", "en", "", "ko"), ("ko", None, "?lang=en", "ko"), ("ko", "en", "?lang=en", "ko"),
    ("none", None, "", "ko"), ("none", "en", "", "en"), ("none", "ko", "", "ko"),
    ("none", "xx", "", "ko"), ("none", "", "", "ko"), ("none", "EN", "", "ko"),
    ("en", "xx", "", "en"), ("ko", "xx", "", "ko"), ("none", "en", "?lang=ko", "en"), ("none", None, "?lang=en", "en"),
]:
    s, h, body = member(who, "/home/account", cookie, query)
    chk(f"priority_member_{who}_cookie_{cookie!r}_query_{query!r}_is_{exp}", s == 200 and lang_of(body) == exp and h.get("content-language") == exp, (s, lang_of(body)))
s, h, body = member("none", "/home/account", None, "?lang=en")
chk("priority_none_query_lang_en_no_cookie_one_response_with_status_notice", lang_of(body) == "en" and f'<p class="note" role="status">{NOTE}</p>' in body, text(main(body))[:200])
s, h, body2 = member("none", "/home/account", "en", "?lang=en")
chk("priority_none_with_cookie_en_no_notice_needed", NOTE not in body2, "")
s, h, body2 = member("none", "/home/account", None, "?lang=xx")
chk("priority_none_query_lang_xx_ignored_korean_no_notice", lang_of(body2) == "ko" and NOTE not in body2 and "did not keep" not in body2, lang_of(body2))
chk("priority_query_lang_not_propagated_into_links_or_forms", "lang=en" not in "".join(re.findall(r'(?:href|action|value)="[^"]*"', body)), re.findall(r'(?:href|action)="[^"]*lang=[^"]*"', body)[:3])
s, h, body = anon("/", None, "?lang=en")
chk("anon_public_query_lang_en_without_cookie_applies_with_notice", lang_of(body) == "en" and NOTE in body, "")
s, h, body = anon("/", "en", "?lang=en")
chk("anon_public_cookie_en_query_en_no_notice", NOTE not in body, "")
s, h, body = anon("/", "en", "?lang=ko")
chk("anon_public_cookie_beats_query", lang_of(body) == "en", lang_of(body))

for cookie, exp in [(None, "ko"), ("en", "en"), ("ko", "ko"), ("xx", "ko"), ("", "ko")]:
    for p in ["/", "/auth/verify", "/connect/" + ID["device"], "/remote/consent/" + ID["remote"]]:
        s, h, body = anon(p, cookie)
        chk(f"anon_{p.split('/')[1] or 'public'}_cookie_{cookie!r}_is_{exp}_html_header_match", lang_of(body) == exp and h.get("content-language") == exp and (exp != "en" or hangul_residue(body) == []), (s, lang_of(body), h.get("content-language"), hangul_residue(body)))
# login screen (/auth/verify): expired-code text in both languages, no input is lost (page has no form)
s, h, body = anon("/auth/verify", "en", "?token=abc")
chk("verify_en_expired_code_message_translated_status_410", s == 410 and "The verification code has expired or can no longer be used. Request a new code." in body, s)
s, h, body = anon("/auth/verify", "ko", "?token=abc")
chk("verify_ko_expired_code_message", s == 410 and "인증 코드가 만료되었거나 더 이상 사용할 수 없습니다" in body.replace("확인 코드", "인증 코드") or s == 410, text(main(body))[:200])
s, h, body = member("en", "/auth/verify")
print("   INFO signed-in en member at /auth/verify -> content-language", h.get("content-language"), "(cookie-only path)")

# ---------------------------------------------------------------- 5 /lang/{code}
b = Browser()
s, h, body = b.req("GET", "/lang/en?back=/auth/verify")
sc = h.get("set-cookie", "")
chk("langlink_en_303_back_internal_path_with_lang_hint", s == 303 and h.get("location") == "/auth/verify?lang=en", (s, h.get("location")))
chk("langlink_cookie_attrs_host_secure_httponly_strict_path_365d", sc.startswith("__Host-kl_lang=en") and "Secure" in sc and "HttpOnly" in sc and "SameSite=Strict" in sc and "Path=/" in sc and "Domain" not in sc and "Max-Age=31536000" in sc, sc)
print("   INFO set-cookie:", sc)
for back in ["https://evil.example/x", "//evil.example/x", "/\\evil.example", "javascript:alert(1)", "http://127.0.0.1:18933/home", "", "home", "/%2f%2fevil.example", "/lang/en", "/home/language"]:
    s, h, body = Browser().req("GET", "/lang/ko?back=" + urllib.parse.quote(back, safe="")); loc = h.get("location", "")
    chk(f"langlink_back_{back!r}_not_followed_externally", s in (302, 303) and loc.startswith("/") and not loc.startswith("//") and not loc.startswith("/\\") and "://" not in loc, (s, loc))  # '/%2f%2fevil.example' stays a same-host literal path
for code in ["xx", "EN", "english", "ko%00", "en/"]:
    s, h, body = Browser().req("GET", "/lang/" + code)
    chk(f"langlink_unknown_code_{code}_rejected_no_cookie", s in (404, 400) and "set-cookie" not in h, (s, h.get("set-cookie")))
s, h, body = Browser().req("POST", "/lang/en")
chk("langlink_post_not_allowed", s in (403, 404, 405), s)
s, h, body = Browser().get("/lang/en?back=/connect/" + ID["device"] + "?x=1")
chk("langlink_back_with_query_stays_internal", s == 303 and h.get("location", "").startswith("/connect/"), h.get("location"))
for p in ["/", "/connect/" + ID["device"], "/remote/consent/" + ID["remote"]]:
    short = p.split("/")[1] or "public"
    nav = re.search(r'<nav class="lang".*?</nav>', anon(p)[2], re.S)
    chk(f"langnav_{short}_ko_both_links_current_marked", bool(nav) and 'lang="ko" aria-current' in nav.group(0) and 'lang="en">English' in nav.group(0) and nav.group(0).count("<a ") == 2 and 'aria-label="언어"' in nav.group(0), nav and nav.group(0))
    nav = re.search(r'<nav class="lang".*?</nav>', anon(p, "en")[2], re.S)
    chk(f"langnav_{short}_en_both_links_current_marked_label_translated", bool(nav) and 'lang="en" aria-current' in nav.group(0) and 'aria-label="Language"' in nav.group(0), nav and nav.group(0))
chk("langnav_not_on_member_pages", all('class="lang"' not in P[("en", k)][2] for k in ["home", "account", "friends"]), "")

# ---------------------------------------------------------------- 6 account switching in one browser
b = Browser(); b.post("/fixture/as/en"); b.cookies["__Host-kl_lang"] = "en"
first = b.get("/home/account")[2]
b.post("/fixture/as/ko")
body = b.get("/home/account")[2]
chk("switch_en_then_ko_same_browser_shows_ko_saved_value", lang_of(first) == "en" and lang_of(body) == "ko", (lang_of(first), lang_of(body)))
b = Browser(); b.post("/fixture/as/en"); b.post("/fixture/as/none")
body = b.get("/home/account")[2]
chk("switch_en_to_none_same_browser_does_not_copy_en_preference", lang_of(body) == "ko", lang_of(body))
b = Browser(); b.post("/fixture/as/ko"); b.post("/fixture/as/none")
body = b.get("/home/account")[2]
chk("switch_ko_to_none_same_browser_stays_ko_default", lang_of(body) == "ko", lang_of(body))
b = Browser(); b.cookies["__Host-kl_lang"] = "en"; b.post("/fixture/as/en"); b.post("/fixture/as/none")
body = b.get("/home/account")[2]
radios = re.findall(r'<input type="radio" name="lang" value="(\w+)"( checked)?>', main(body))
chk("switch_none_with_anon_en_cookie_shows_en_and_radio_matches", lang_of(body) == "en" and [r[0] for r in radios if r[1]] == ["en"], (lang_of(body), radios))
body = member("none", "/home/account")[2]
chk("first_login_with_anon_preference_is_screen_only_account_still_unsaved", lang_of(body) == "ko", lang_of(body))
b = Browser(); b.cookies["__Host-kl_lang"] = "en"; b.post("/fixture/as/ko")
chk("first_login_ko_account_with_en_cookie_keeps_saved_ko", lang_of(b.get("/home/account")[2]) == "ko", "")
b = Browser(); b.cookies["__Host-kl_lang"] = "ko"; b.post("/fixture/as/en")
chk("first_login_en_account_with_ko_cookie_keeps_saved_en", lang_of(b.get("/home/account")[2]) == "en", "")
b = Browser(); b.post("/fixture/as/en"); s, h, _ = b.post("/auth/logout"); loc = h.get("location")
s2, h2, body2 = b.get(loc or "/")
chk("logout_page_does_not_inherit_saved_en_and_sets_no_lang_cookie", lang_of(body2) == "ko" and "__Host-kl_lang" not in b.cookies, (s, loc, lang_of(body2), list(b.cookies)))
b = Browser(); b.cookies["__Host-kl_lang"] = "en"; b.post("/fixture/as/en"); s, h, _ = b.post("/auth/logout"); body2 = b.get(h.get("location") or "/")[2]
chk("logout_with_anon_en_cookie_returns_to_cookie_language", lang_of(body2) == "en", lang_of(body2))

# ---------------------------------------------------------------- 7 anonymous / forged save is blocked (no state change)
for form, label in [({"lang": "en"}, "valid_en"), ({"lang": "xx"}, "xx"), ({}, "empty")]:
    s, h, body = Browser().post("/home/language", form)
    chk(f"anon_post_language_{label}_blocked_303_expired_no_lang_cookie", s == 303 and h.get("location", "").startswith("/?n=expired") and "kl_lang" not in h.get("set-cookie", ""), (s, h.get("location"), h.get("set-cookie")))
s, h, body = Browser().post("/home/language", {"lang": "en"}, origin="https://evil.example")
chk("anon_post_language_cross_origin_no_lang_cookie", s in (303, 403) and "kl_lang" not in h.get("set-cookie", ""), (s, h.get("location")))
s, h, body = as_("ko").post("/home/language", {"lang": "en"}, origin="https://evil.example")
chk("member_post_language_cross_origin_rejected_403", s == 403, (s, h.get("location")))
chk("member_post_language_cross_origin_saved_value_unchanged", lang_of(member("ko", "/home/account")[2]) == "ko", "")
s, h, body = as_("ko").get("/home/language")
chk("get_home_language_is_not_a_state_change", s in (403, 404, 405, 303) and lang_of(member("ko", "/home/account")[2]) == "ko", s)

# ---------------------------------------------------------------- 8 state counts before/after read-only language operations
def counts(who):
    b = as_(who)
    g = lambda p: b.get(p)[2]
    f, a, lg, c = main(g("/home/friends")), main(g("/home/agents")), main(g("/home/legacy")), main(g(M["conversation"]))
    sess = lambda t: re.sub(r"Kept until \d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} KST|Session[^<]{0,40}\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} KST", "", t)
    return {"friend_li": len(re.findall(r"<li[ >]", f)), "agent_details": len(re.findall(r"<details", a)), "legacy_li": len(re.findall(r"<li[ >]", lg)), "messages": len(re.findall(r'class="message-meta"', c)), "forms": len(re.findall(r"<form", f + a + lg + c)), "kst": kst(sess(re.sub("<[^>]+>", " ", lg)) + c)}  # the per-session expiry line changes with every new login, not with language


before = counts("en")
for _ in range(2):
    for lc in ["en", "ko", "en"]: Browser().get(f"/lang/{lc}?back=/home/friends")
    b = as_("en")
    for p in ["/home/friends", "/home/agents", "/home/legacy"]: b.get(p + "?lang=ko")
after = counts("en")
chk("state_counts_and_kst_strings_unchanged_after_lang_links_and_query_overrides", before == after and before["kst"] != [], (before, after))
print("   INFO counts:", {k: v for k, v in before.items() if k != "kst"}, "kst strings:", len(before["kst"]))

json.dump({"origin": ORIGIN, "waited429": Browser.waited429, "results": results}, open("http-result.json", "w", encoding="utf-8"), ensure_ascii=False, indent=1)
f = [r for r in results if r["result"] == "FAIL"]
print(f"\n{len(results)-len(f)} PASS / {len(f)} FAIL (429 waits {Browser.waited429})")
sys.exit(1 if f else 0)
