"""Recheck of DEV b11d5b4 (product 665895a): saved-language link swap (V-LANG-01/F3), README hreflang (F2), 'your' wording (V-LANG-02),
unsaved note (V-LANG-03). GET + /lang/* (cookie only) + one rejected POST /home/language. New checks only; no expectation of checks.py is weakened.
Run: python recheck.py -> recheck-result.json (exit 1 on FAIL)."""
import json, re, sys, urllib.parse
sys.path.insert(0, ".")
from client import *
from dump import SCREENS
from checks_lib import chk, results, lang_of, hangul_residue

SETTINGS = {"en": ("English", "Change in Settings"), "ko": ("한국어", "설정에서 변경")}
PAGES = {"public": "/", "device": "/connect/" + ID["device"], "remote": "/remote/consent/" + ID["remote"]}
SPAN = lambda L: re.compile(r'<span class="lang">%s<span aria-hidden="true">·</span><a href="/home/account#language">%s</a></span>' % SETTINGS[L])
leaks = []

# ---------------------------------------------------------------- 1 saved member: link swap on the three screens, either cookie
for who in ["en", "ko"]:
    other = "ko" if who == "en" else "en"
    for cookie in [None, other, who]:
        for k, p in PAGES.items():
            b = as_(who)
            if cookie: b.cookies["__Host-kl_lang"] = cookie
            s, h, body = b.get(p)
            leaks += [(who, k, x) for x in h if x.startswith("x-knowslink")]
            tag = f"{who}_{k}_cookie_{cookie}"
            chk(f"swap_{tag}_replaces_language_links_with_settings_link", s == 200 and lang_of(body) == who and bool(SPAN(who).search(body)), (s, lang_of(body)))
            chk(f"swap_{tag}_no_lang_route_no_lang_nav", "/lang/" not in body and 'class="lang"' in body and '<nav class="lang"' not in body, re.findall(r'/lang/[^"]*', body))
            chk(f"swap_{tag}_no_internal_marker_in_page_or_headers", "X-KnowsLink-Lang" not in body and not [x for x in h if x.startswith("x-knowslink")], list(h))
    # the link target works: settings, language section, saved value checked
    b = as_(who); s, h, body = b.get("/home/account#language")
    chk(f"swap_{who}_settings_target_opens_language_section_with_saved_value", s == 200 and 'id="language"' in body and re.findall(r'value="(\w+)" checked', main(body)) == [who], s)
    # saved member can not change language by a stale /lang route: cookie only changes, saved value still wins
    b = as_(who); b.get(f"/lang/{other}?back=/"); s, h, body = b.get("/")
    chk(f"swap_{who}_manual_lang_route_does_not_override_saved_and_page_keeps_settings_link", lang_of(body) == who and bool(SPAN(who).search(body)), lang_of(body))
print("   INFO x-knowslink headers on GET 200 member pages:", sorted(set(leaks)) or "none")

# refusal 403 screens and /auth/verify for a saved member (what link do they show?)
for who in ["ko"]:
    b = as_(who)
    for k in ["conn-waiting", "gate"]:
        s, h, body = b.get(SCREENS[k])
        nav = re.search(r'<nav class="lang".*?</nav>|<span class="lang">.*?</span></span>', body, re.S)
        print("   INFO", who, k, s, lang_of(body), nav.group(0)[:200] if nav else "no lang control")
        chk(f"refusal_{k}_{who}_is_403_in_saved_language", s == 403 and lang_of(body) == who, (s, lang_of(body)))
for who in ["en", "ko"]:
    s, h, body = as_(who).get("/auth/verify")
    nav = re.search(r'<nav class="lang".*?</nav>|<span class="lang">.*?</span></span>', body, re.S)
    print("   INFO", who, "/auth/verify", s, lang_of(body), nav.group(0)[:140] if nav else "no lang control")

# ---------------------------------------------------------------- 2 unsaved member (none) and anonymous: /lang/ links remain and work
for k, p in PAGES.items():
    for cookie, exp in [(None, "ko"), ("en", "en")]:
        b = as_("none")
        if cookie: b.cookies["__Host-kl_lang"] = cookie
        s, h, body = b.get(p)
        nav = re.search(r'<nav class="lang".*?</nav>', body, re.S)
        chk(f"unsaved_none_{k}_cookie_{cookie}_keeps_lang_links_not_settings_link", s == 200 and lang_of(body) == exp and bool(nav) and nav.group(0).count("/lang/") == 2 and "/home/account#language" not in body, (s, nav and nav.group(0)[:100]))
    b = as_("none")
    s, h, body = b.get(p); nav = re.search(r'<nav class="lang".*?</nav>', body, re.S).group(0)
    en_href = re.search(r'<a href="([^"]+)" lang="en"', nav).group(1).replace("&amp;", "&")
    s, h, _ = b.get(en_href)
    chk(f"unsaved_none_{k}_english_link_sets_cookie_and_returns_same_screen", s == 303 and b.cookies.get("__Host-kl_lang") == "en" and h.get("location", "").split("?")[0] == p, (s, h.get("location"), b.cookies))
    s, h, body = b.get(h["location"])
    chk(f"unsaved_none_{k}_after_english_link_page_is_english_and_current_marked", lang_of(body) == "en" and 'lang="en" aria-current' in body, lang_of(body))
    ko_href = re.search(r'<a href="([^"]+)" lang="ko"', re.search(r'<nav class="lang".*?</nav>', body, re.S).group(0)).group(1).replace("&amp;", "&")
    s, h, _ = b.get(ko_href); s2, h2, body2 = b.get(h["location"])
    chk(f"unsaved_none_{k}_korean_link_returns_to_korean", lang_of(body2) == "ko", lang_of(body2))
    # anonymous: same two links work (already in checks.py; here the click path per screen)
    b = Browser(); s, h, body = b.get(p); nav = re.search(r'<nav class="lang".*?</nav>', body, re.S).group(0)
    s, h, _ = b.get(re.search(r'<a href="([^"]+)" lang="en"', nav).group(1).replace("&amp;", "&")); s2, h2, body2 = b.get(h["location"])
    chk(f"anon_{k}_english_link_works", s == 303 and lang_of(body2) == "en" and '<nav class="lang"' in body2, (s, lang_of(body2)))

# ---------------------------------------------------------------- 3 README links: hreflang=ko plus language note
def guide_links(body): return re.findall(r'<a [^>]*href="https://github\.com/[^"]*README[^"]*"[^>]*>.*?</a>', body)
en_links = []
def collect(L, bro, label, path):
    s, h, body = bro.get(path); ls = guide_links(body)
    for l in ls:
        tag = re.sub(r"[^a-z0-9]+", "_", (label + path).lower()).strip("_")
        text_ = re.sub(r"<[^>]+>", "", l)
        if L == "en":
            chk(f"readme_en_{tag}_{len(en_links)}_hreflang_ko_and_in_korean_note", 'hreflang="ko"' in l and text_.endswith("(in Korean)"), l[:200])
            en_links.append(l)
        else:
            chk(f"readme_ko_{tag}_{len(ls)}_hreflang_ko_no_english_note", 'hreflang="ko"' in l and "(in Korean)" not in text_ and "Korean" not in text_, l[:200])
    return len(ls), body
total = 0
for L in ["en", "ko"]:
    for who in ["anon", "en", "ko"]:
        if (L == "en") != (who in ("anon", "en")) and who != "anon": continue
        for q in ["/", "/home/agents?use=claude", "/home/agents?use=claude&start=first", "/home/agents?use=claude&start=continue", "/home/agents?use=grok", "/home/agents?use=grok&start=first", "/home/agents?use=grok&start=continue", "/home/agents?use=dots", "/home/agents?use=dots&start=first", "/home/agents?use=dots&start=continue", "/home/agents?use=other", "/home/agents?use=web", "/home/legacy", "/home/new"]:
            if who == "anon":
                if q != "/": continue
                bro = Browser()
                if L == "en": bro.cookies["__Host-kl_lang"] = "en"
            else: bro = as_(who)
            n, _ = collect(L, bro, f"{who}_{L}_", q); total += n
chk("readme_links_found_in_both_languages_at_least_10_each", len([x for x in en_links]) >= 10, len(en_links))
print("   INFO README links checked (en):", len(en_links), "total:", total)
# non-README GitHub links in English screens: list for information
b = as_("en"); infos = set()
for k, p in SCREENS.items():
    s, h, body = b.get(p)
    for l in re.findall(r'<a [^>]*href="https?://[^"]*"[^>]*>.*?</a>', body):
        if "README" not in l: infos.add(re.sub(r"\s+", " ", l)[:160])
print("   INFO other external links on en screens:", sorted(infos) or "none")

# ---------------------------------------------------------------- 4 'your' wording (V-LANG-02): 'my' only inside the three proper names
ALL = {}
b = as_("en")
for k, p in SCREENS.items(): ALL[k] = b.get(p)[2]
for q in ["?use=claude&start=first", "?use=claude&start=continue", "?use=grok&start=first", "?use=grok&start=continue", "?use=dots&start=first", "?use=dots&start=continue", "?use=other", "?use=web", "?use=claude&start=existing"]:
    ALL["agents" + q] = b.get("/home/agents" + q)[2]
for q in ["/home/receipts?agent=agent_b", "/home/receipts?agent=agent_fx_waiting", "/home/new?to=agent_a"]:
    ALL[q] = b.get(q)[2]
ALL["anon-public"] = (lambda x: (x.cookies.update({"__Host-kl_lang": "en"}), x.get("/"))[1][2])(Browser())
ALL["conn-approved-detail"] = b.get("/home/connections/" + ID["connection_approved"])[2]
def vis(body):
    t = re.sub(r"<(script|style)[^>]*>.*?</\1>", " ", body, flags=re.S)
    attrs = " ".join(re.findall(r'(?:aria-label|title|alt|placeholder)="([^"]*)"', t))
    return re.sub(r"\s+", " ", re.sub(r"<[^>]+>", " ", t) + " " + attrs)
PROPER = re.compile(r"My agents|My name|My friend code", re.I)
for k, body in ALL.items():
    t = PROPER.sub("", vis(body)); bad = re.findall(r".{25}\bmy\b.{25}", t, re.I)
    chk(f"wording_{re.sub(r'[^a-z0-9]+', '_', k.lower())}_no_first_person_my_outside_proper_names", not bad, bad[:3])
allt = " ".join(vis(x) for x in ALL.values())
chk("wording_old_first_person_phrases_absent", not re.search(r"on my computer|on my account|of my agent|my sending agent|Select my agent|Sign in and choose my name|Connect my (client|agent)|Manage my agents|Start a conversation with my agent|My own agents|to my account", allt, re.I), re.findall(r".{20}(?:on my computer|my sending agent|Select my agent|Connect my|Manage my).{20}", allt, re.I)[:3])
SAMPLES = ["your computer", "Connect your client", "Manage your agents", "Your own agents can talk", "Sign in and choose your name", "Connect your agent", "your sending agent", "Start a conversation with your agent", "Selecting a language changes nothing until you save. You stay signed in."]
found = {s: (s in allt) for s in SAMPLES}
print("   INFO designer sample phrases found on reachable en screens:", found)
chk("wording_designer_sample_phrases_present_at_least_6_of_9", sum(found.values()) >= 6, found)
chk("wording_designer_sample_phrases_core_present", all(found[s] for s in ["your computer", "Connect your client", "Selecting a language changes nothing until you save. You stay signed in."]), found)
chk("wording_en_settings_help_exact_and_ko_help_unchanged", "Selecting a language changes nothing until you save. You stay signed in." in ALL["account"] and "고르기만으로는 바뀌지 않으며 다시 로그인하지 않습니다" in as_("ko").get("/home/account")[2], "")
s, h, body = as_("en").post("/home/language", {"lang": "xx"})
chk("wording_unsupported_language_error_uses_visible_label_422", s == 422 and re.search(r"Choose (<span[^>]*>)?한국어(</span>)? or English", body) is not None, text(main(body))[:250])
chk("wording_en_screens_still_no_hangul_outside_names_and_label_after_wording_change", all(hangul_residue(b2) == [] for k, b2 in ALL.items() if k != "verify"), {k: hangul_residue(b2) for k, b2 in ALL.items() if k != "verify" and hangul_residue(b2)})  # /auth/verify is cookie-only (O-LG-01), Korean for a cookieless member

# ---------------------------------------------------------------- 5 unsaved note (V-LANG-03)
NOTE = {"ko": "이 브라우저의 선택을 표시합니다. 아직 계정에 저장되지 않았습니다.", "en": ["This browser’s choice is shown. It is not saved to your account yet.", "This browser's choice is shown. It is not saved to your account yet."]}
for cookie, L in [(None, "ko"), ("ko", "ko"), ("en", "en")]:
    b = as_("none")
    if cookie: b.cookies["__Host-kl_lang"] = cookie
    s, h, body = b.get("/home/account"); m = main(body)
    note = re.search(r'<p id="language-unsaved"[^>]*>(.*?)</p>', m)
    want = NOTE[L] if isinstance(NOTE[L], list) else [NOTE[L]]
    chk(f"unsaved_note_none_cookie_{cookie}_shown_in_{L}_and_linked_to_radiogroup", s == 200 and lang_of(body) == L and bool(note) and any(w in note.group(1) for w in want) and 'aria-describedby="language-unsaved language-help"' in m, (note and note.group(1)))
s, h, body = as_("none").get("/home/account?lang=en")
chk("unsaved_note_none_query_lang_en_note_in_english_plus_cookie_notice", lang_of(body) == "en" and 'id="language-unsaved"' in body, lang_of(body))
for who in ["en", "ko"]:
    for cookie in [None, "en", "ko"]:
        b = as_(who)
        if cookie: b.cookies["__Host-kl_lang"] = cookie
        body = b.get("/home/account")[2]
        chk(f"unsaved_note_absent_for_saved_{who}_cookie_{cookie}", 'language-unsaved' not in body and 'aria-describedby="language-help"' in body, "")
body = as_("none").get("/home/account")[2]
chk("unsaved_none_settings_has_no_settings_self_link_swap", "/home/account#language" not in body or body.count("/home/account#language") == 0, body.count("/home/account#language"))

json.dump({"origin": ORIGIN, "waited429": Browser.waited429, "results": results}, open("recheck-result.json", "w", encoding="utf-8"), ensure_ascii=False, indent=1)
f = [r for r in results if r["result"] == "FAIL"]
print(f"\n{len(results)-len(f)} PASS / {len(f)} FAIL (429 waits {Browser.waited429})")
sys.exit(1 if f else 0)
