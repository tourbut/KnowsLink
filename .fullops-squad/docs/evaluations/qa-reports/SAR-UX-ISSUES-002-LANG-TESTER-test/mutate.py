"""State-changing HTTP checks (run LAST; the designer shares the fixture). Only POST /home/language on the en and ko accounts, each restored.
Run: python mutate.py -> mutate-result.json (exit 1 on FAIL). none account is never saved: a saved value cannot be unset."""
import json, re, sys
sys.path.insert(0, ".")
from client import *
from checks_lib import chk, results, lang_of, hangul_residue, kst, main_text

def counts(b):
    g = lambda p: b.get(p)[2]
    f, a, lg, c, g1 = g("/home/friends"), g("/home/agents"), g("/home/legacy"), g(M["conversation"]), g("/home/gates/" + ID["gate"])
    sess = lambda t: re.sub(r"(Session|세션).{0,80}?KST", "", re.sub(r"<[^>]+>", " ", t))
    return {"friend_li": len(re.findall(r"<li[ >]", main(f))), "agent_details": len(re.findall(r"<details", main(a))), "legacy_li": len(re.findall(r"<li[ >]", main(lg))),
            "messages": len(re.findall(r'class="message-meta"', c)), "gate_forms": len(re.findall(r"<form", g1)), "kst": kst(sess(lg) + c + g1 + f + a)}

def invariants(b):
    """Language-independent parts of the owner screens: ids, hrefs, form actions, hidden values, deadline strings, signed original JSON."""
    out = {}
    for k, p in {"gate": "/home/gates/" + ID["gate"], "conversation": M["conversation"], "legacy": "/home/legacy", "agents": "/home/agents", "conn-waiting": "/home/connections/" + ID["connection_waiting"], "conn-prepared": "/home/connections/" + ID["connection_prepared"], "conn-approved": "/home/connections/" + ID["connection_approved"], "friends": "/home/friends"}.items():
        body = b.get(p)[2]; m = re.sub(r"(Session|세션).{0,160}?KST", "", body, flags=re.S) if k == "legacy" else body
        out[k] = {"hrefs": re.findall(r'(?:href|action)="([^"]*)"', m), "hidden": [h for h in re.findall(r'<input type="hidden" name="(\w+)" value="([^"]*)"', m) if h[0] not in ("csrf", "attempt")],  # csrf/attempt tokens are per session by design
 "ids": sorted(set(re.findall(r"agent_\w+|0199a3f2-[0-9a-f-]{23}|[0-9a-f]{8}(?=<)", m))),
                  "kst": kst(m), "json": [re.sub(r"&lt;[^&]*?&gt;", "<PLACEHOLDER>", c) for c in re.findall(r"<pre[^>]*>(.*?)</pre>", m, re.S)], "code": [re.sub(r"&lt;[^&]*?&gt;", "<PLACEHOLDER>", c) for c in re.findall(r"<code[^>]*>(.*?)</code>", m, re.S)]}  # a command keeps its words; its <placeholder> is prose
    return out

en = as_("en"); ko = as_("ko")
before_en = counts(en); inv_en = invariants(en)
chk("pre_state_en_saved_english_ko_saved_korean", lang_of(en.get("/home/account")[2]) == "en" and lang_of(ko.get("/home/account")[2]) == "ko", "")

# ---- rejected values: previous saved value and selection kept, failure is told, English refusal in English
for form, label in [({"lang": "xx"}, "xx"), ({"lang": ""}, "empty"), ({}, "missing"), ({"lang": "EN"}, "uppercase"), ({"lang": "ko\u0000"}, "nul"), ({"lang": "ko", "lang2": "en"}, "extra_field_ok_ko")]:
    b = as_("en")
    s, h, body = b.post("/home/language", form)
    if label == "extra_field_ok_ko":   # a valid value with an unknown extra field is a legitimate save: restore below
        chk("lang_post_valid_with_extra_field_saves_ko_then_restored", s == 303, (s, h.get("location")))
        b2 = as_("en"); b2.post("/home/language", {"lang": "en"}); continue
    m = main(body)
    chk(f"reject_{label}_422_alert_in_saved_language", s == 422 and lang_of(body) == "en" and 'role="alert"' in m and hangul_residue(body.replace("Choose 한국어 or English.", "")) == [], (s, lang_of(body), main_text(m)[:200]))  # b11d5b4 V-LANG-02: designer-specified sentence shows the visible label '한국어'
    chk(f"reject_{label}_keeps_previous_selection_en_checked", [r[0] for r in re.findall(r'<input type="radio" name="lang" value="(\w+)"( checked)?>', m) if r[1]] == ["en"], re.findall(r'<input type="radio"[^>]*>', m))
    chk(f"reject_{label}_does_not_change_saved_value_other_session", lang_of(as_("en").get("/home/account")[2]) == "en", "")
    chk(f"reject_{label}_message_says_failure_not_saved_and_save_button_still_there", ("not" in main_text(m).lower() or "unsupported" in main_text(m).lower()) and "Save language" in m, main_text(m)[:300])
s, h, body = as_("ko").post("/home/language", {"lang": "xx"})
chk("reject_xx_for_ko_member_shows_korean_alert_ko_checked", s == 422 and lang_of(body) == "ko" and 'role="alert"' in body and re.findall(r'value="(\w+)" checked', body) == ["ko"], (s, lang_of(body)))

# ---- successful save en -> ko, observed over several sessions, then restored
b = as_("en"); other = as_("en")
open_other = other.get("/home/account")[2]            # a second open "tab" before the change
s, h, body = b.post("/home/language", {"lang": "ko"})
chk("save_ko_303_to_same_settings_place", s == 303 and h.get("location") == "/home/account?n=language#language", (s, h.get("location")))
chk("save_ko_response_sets_no_language_cookie", "kl_lang" not in h.get("set-cookie", ""), h.get("set-cookie"))
s, h, body = b.get(h["location"])
m = main(body)
chk("save_ko_landing_is_korean_with_status_and_ko_checked", lang_of(body) == "ko" and h.get("content-language") == "ko" and 'role="status"' in m and re.findall(r'value="(\w+)" checked', m) == ["ko"], (lang_of(body), re.findall(r'role="status"[^>]*>.*?</', m)))
print("   INFO save status text:", re.findall(r'role="status"[^>]*>(.*?)</', m))
chk("save_ko_landing_title_h1_settings_in_korean", re.search(r"<h1>(.*?)</h1>", m).group(1) == "설정" and "설정" in re.search(r"<title>(.*?)</title>", body).group(1), "")
chk("save_ko_open_other_session_not_forced_reload_but_next_request_is_ko", lang_of(open_other) == "en" and lang_of(other.get("/home/friends")[2]) == "ko", "")
chk("save_ko_reload_keeps", lang_of(b.get("/home/account")[2]) == "ko" and "role=\"status\"" not in b.get("/home/account")[2], "")
chk("save_ko_new_login_session_same_account_keeps", lang_of(as_("en").get("/home/account")[2]) == "ko", "")
fresh = Browser(); fresh.post("/fixture/as/en"); fresh.post("/fixture/as/en")
chk("save_ko_new_browser_without_cookies_gets_ko_after_login", lang_of(fresh.get("/home")[2]) == "ko", "")
chk("save_ko_anonymous_browser_and_other_account_unaffected", lang_of(Browser().get("/")[2]) == "ko" and lang_of(as_("none").get("/home/account")[2]) == "ko" and lang_of(as_("ko").get("/home/account")[2]) == "ko", "")
# owner screens in Korean: same ids/hrefs/deadlines/originals as in English; no residual English UI words
bk = as_("en")
inv_ko = invariants(bk)
for k in inv_en:
    chk(f"invariant_owner_{k}_hrefs_hidden_ids_kst_json_identical_en_vs_ko", inv_en[k] == inv_ko[k], {x: (inv_en[k][x], inv_ko[k][x]) for x in inv_en[k] if inv_en[k][x] != inv_ko[k][x]})
for k, p in {"gate": "/home/gates/" + ID["gate"], "conversation": M["conversation"], "legacy": "/home/legacy", "agents": "/home/agents", "conn-waiting": "/home/connections/" + ID["connection_waiting"], "conn-prepared": "/home/connections/" + ID["connection_prepared"], "conn-approved": "/home/connections/" + ID["connection_approved"], "friends": "/home/friends", "new": "/home/new", "receipts": "/home/receipts", "name": "/home/name"}.items():
    s, h, body = bk.get(p)
    chk(f"owner_{k}_in_korean_lang_header_title", s == 200 and lang_of(body) == "ko" and h.get("content-language") == "ko", (s, lang_of(body)))
    if k in ("gate", "legacy", "conversation"):
        pre = re.sub(r"<(pre|code)[^>]*>.*?</\1>", "", body, flags=re.S)
        words = sorted(set(w for w in re.findall(r"\b[a-z]{4,}\b", main_text(main(pre))) if w not in {"with", "long", "name", "helper", "registered", "approved", "prepared", "waiting", "local", "client", "synthetic", "member", "node", "agent", "kind", "query", "schedule", "fixture", "message", "bytes", "example", "idle", "test", "text", "wake"}))  # user data / protocol and technical terms kept as in the source
        chk(f"owner_{k}_korean_page_has_no_english_ui_words", words == [], words[:10])
after_ko = counts(bk)
chk("state_counts_and_kst_unchanged_after_language_save", {x: v for x, v in before_en.items()} == after_ko, (before_en, after_ko))
# anonymous save blocked still; saving does not touch the other account
chk("ko_account_still_korean_en_save_did_not_copy", lang_of(as_("ko").get("/home/account")[2]) == "ko", "")
# restore en
s, h, body = as_("en").post("/home/language", {"lang": "en"})
chk("restore_en_303", s == 303, s)
chk("restored_en_is_english_again", lang_of(as_("en").get("/home/account")[2]) == "en", "")
# idempotent save of the current value
b = as_("en"); s, h, body = b.post("/home/language", {"lang": "en"}); chk("save_same_value_again_is_ok_and_stays_en", s == 303 and lang_of(b.get(h["location"])[2]) == "en", s)

# ---- ko account: English then restored
b = as_("ko"); s, h, body = b.post("/home/language", {"lang": "en"})
page = b.get(h["location"])[2]
chk("ko_account_saves_english_lands_in_english_settings_same_place", s == 303 and h.get("location") == "/home/account?n=language#language" and lang_of(page) == "en" and hangul_residue(page) == [] and 'role="status"' in page, (s, lang_of(page), hangul_residue(page)))
print("   INFO ko->en status text:", re.findall(r'role="status"[^>]*>(.*?)</', page))
chk("ko_account_in_english_does_not_affect_en_anonymous_or_none", lang_of(Browser().get("/")[2]) == "ko" and lang_of(as_("none").get("/home/account")[2]) == "ko", "")
ko_en_counts = counts(as_("ko"))
chk("ko_account_english_conversation_still_ok_and_no_hangul_except_names", hangul_residue(as_("ko").get(M["conversation"])[2]) == [], hangul_residue(as_("ko").get(M["conversation"])[2]))
b = as_("ko"); s, h, body = b.post("/home/language", {"lang": "ko"})
chk("restore_ko", s == 303 and lang_of(as_("ko").get("/home/account")[2]) == "ko", s)
chk("final_state_en_english_ko_korean_none_unsaved", lang_of(as_("en").get("/home/account")[2]) == "en" and lang_of(as_("ko").get("/home/account")[2]) == "ko" and lang_of(as_("none").get("/home/account")[2]) == "ko", "")

json.dump({"origin": ORIGIN, "waited429": Browser.waited429, "results": results}, open("mutate-result.json", "w", encoding="utf-8"), ensure_ascii=False, indent=1)
f = [r for r in results if r["result"] == "FAIL"]
print(f"\n{len(results)-len(f)} PASS / {len(f)} FAIL (429 waits {Browser.waited429})")
sys.exit(1 if f else 0)
