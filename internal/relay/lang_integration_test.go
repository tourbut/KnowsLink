//go:build integration

// U3-A: saving a language, failing to, refreshing, signing in again, switching accounts, the anonymous choice, a
// damaged locale, and that choosing a language never performs or loses anything else.
package relay

import (
	"fmt"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestLanguagePreferencePostgres(t *testing.T) {
	s, _ := identityService(t, messagePool(t))
	h := s.Handler()
	newBrowser := func(ip string) *browser { return &browser{t, h, ip, map[string]string{}} }
	// signIn mints a verified email session directly (the mail code path is covered by the identity tests and its
	// 60-second resend spacing would refuse the repeated sign-ins this check needs) and names the member once.
	signIn := func(b *browser, email string) string {
		t.Helper()
		id := ""
		s.mutateState(t, func(st *State) {
			v := st.signIn(emailIssuer, email, email, "", "", time.Now())
			if v.Session == "" {
				t.Fatal("sign-in refused", v.Problem)
			}
			b.cookies[sessionCookie] = v.Session
			id = st.Sessions[hashToken(v.Session)].Member
			if st.Members[id].Name == "" {
				st.Members[id].Name = "Synthetic " + email[:6]
			}
		})
		return id
	}
	locale := func(id string) string {
		v := ""
		s.mutateState(t, func(st *State) { v = st.Members[id].Locale })
		return v
	}
	snapshot := func() string {
		out := ""
		s.mutateState(t, func(st *State) {
			out = fmt.Sprint(len(st.Friends), len(st.Agents), len(st.Gates), len(st.Messages), len(st.Connections), len(st.Pairs))
		})
		return out
	}
	english := func(body string) {
		t.Helper()
		if !strings.Contains(body, `<html lang="en">`) {
			t.Fatal("page not English")
		}
		if left := koreanLeft(body); len(left) > 0 {
			t.Fatalf("Korean left in English page: %q", left)
		}
	}
	a := newBrowser("192.0.2.31")
	idA := signIn(a, "lang-a@example.test")
	before := snapshot()

	// Default Korean settings page in the U3 order, nothing checked but the current language.
	w := a.do("GET", "/home/account", nil)
	expect(t, w, 200, `<html lang="ko">`, "<h1>설정</h1>", `value="ko" checked`, `action="/home/language"`)
	if w.Header().Get("Content-Language") != "ko" || strings.Contains(w.Body.String(), `value="en" checked`) {
		t.Fatal("default language", w.Header())
	}
	body := w.Body.String()
	if !(strings.Index(body, `id="profile-h"`) < strings.Index(body, `id="code-h"`) && strings.Index(body, `id="code-h"`) < strings.Index(body, `id="language-h"`) && strings.Index(body, `id="language-h"`) < strings.Index(body, `id="session-h"`)) {
		t.Fatal("settings order must be name, friend code, language, sign-in and security")
	}
	if strings.Contains(body, "내 계정") || !strings.Contains(a.do("GET", "/home", nil).Body.String(), `href="/home/account"`) {
		t.Fatal("header must name the settings entry")
	}

	// Save English: the same settings place comes back in English with a status notice and the choice checked.
	expect(t, a.do("POST", "/home/language", url.Values{"lang": {"en"}}), 303, "/home/account?n=language#language")
	w = a.do("GET", "/home/account?n=language", nil)
	expect(t, w, 200, `value="en" checked`, `role="status"`, `id="language"`)
	english(w.Body.String())
	if w.Header().Get("Content-Language") != "en" || locale(idA) != "en" || a.cookies[langCookie] != "" {
		t.Fatal("member save must set the member locale only", locale(idA), a.cookies)
	}
	// The saved value decides the public page too, so it offers the settings place instead of browser links (V-LANG-01).
	w = a.do("GET", "/", nil)
	expect(t, w, 200, `<html lang="en">`, `<a href="/home/account#language">Change in Settings</a>`)
	if strings.Contains(w.Body.String(), "/lang/") || w.Header().Get(langSavedHeader) != "" {
		t.Fatal("saved member must not get browser language links or the marker", w.Header())
	}
	for _, path := range []string{"/home", "/home/friends", "/home/agents", "/home/new", "/home/legacy", "/home/receipts"} {
		w = a.do("GET", path, nil)
		expect(t, w, 200)
		english(w.Body.String())
	}
	// Another browser of the same account sees English too; refusals follow the member's language.
	a2 := newBrowser("192.0.2.32")
	signIn(a2, "lang-a@example.test")
	w = a2.do("GET", "/home", nil)
	expect(t, w, 200)
	english(w.Body.String())
	w = a2.do("POST", "/home/agents", url.Values{"name": {"‎"}})
	expect(t, w, 422)
	english(w.Body.String())

	// A redirect of a member with a saved value renders nothing, yet must not carry the internal markers either.
	for _, w := range []*httptest.ResponseRecorder{a.do("POST", "/home/language", url.Values{"lang": {"en"}}), a.do("POST", "/auth/logout", nil)} {
		expect(t, w, 303)
		if sent := w.Result().Header; sent.Get(langSavedHeader) != "" || sent.Get(langUnsavedHeader) != "" {
			t.Fatal("redirect carried an internal language marker", sent)
		}
	}
	// Signed out: the public page follows the browser, which holds no choice, so Korean.
	expect(t, a.do("GET", "/", nil), 200, `<html lang="ko">`)
	// Account switch in the same browser: B has no saved value and inherits nothing from A.
	idB := signIn(a, "lang-b@example.test")
	expect(t, a.do("GET", "/home", nil), 200, `<html lang="ko">`)
	if locale(idB) != "" {
		t.Fatal("B must not inherit A's language")
	}

	// Anonymous explicit choice: one cookie, a return to the same page, English for the next visit, and the first
	// sign-in of a member without a saved value applies it to the screen without saving it.
	expect(t, a.do("POST", "/auth/logout", nil), 303)
	w = a.do("GET", "/lang/en?back=%2F", nil)
	expect(t, w, 303, "/?lang=en")
	if a.cookies[langCookie] != "en" {
		t.Fatal("anonymous choice must be stored in the language cookie")
	}
	w = a.do("GET", "/", nil)
	expect(t, w, 200, `<html lang="en">`, `lang="en" aria-current="true"`)
	if strings.Contains(w.Body.String(), `role="status"`) {
		t.Fatal("a stored cookie shows no unsaved notice")
	}
	english(w.Body.String())
	signIn(a, "lang-b@example.test")
	w = a.do("GET", "/home/account", nil)
	expect(t, w, 200, `value="en" checked`, `id="language-unsaved"`, "not saved to your account yet")
	english(w.Body.String())
	expect(t, a.do("GET", "/", nil), 200, `/lang/ko?back=`) // no saved value: the browser links still work
	if locale(idB) != "" {
		t.Fatal("screen application must not save the browser choice to the account")
	}
	// B saves Korean explicitly: the member value wins over the browser cookie; after sign-out the cookie applies again.
	expect(t, a.do("POST", "/home/language", url.Values{"lang": {"ko"}}), 303)
	expect(t, a.do("GET", "/home/friends", nil), 200, `<html lang="ko">`)
	if locale(idB) != "ko" || a.cookies[langCookie] != "en" {
		t.Fatal("member save must leave the browser cookie alone")
	}
	if strings.Contains(a.do("GET", "/home/account", nil).Body.String(), "language-unsaved") {
		t.Fatal("saved value must not be called unsaved")
	}
	// The member request limit refuses at admission, before any member lookup: the browser cookie (en) decides the
	// language even though B saved Korean (U3 LANG-VISUAL).
	s.mutateState(t, func(st *State) {
		for range 40 {
			st.Rates["http:member:"+idB] = append(st.Rates["http:member:"+idB], time.Now())
		}
	})
	w = a.do("GET", "/home/friends", nil)
	expect(t, w, 429, `<html lang="en">`)
	english(w.Body.String())
	if w.Header().Get("Content-Language") != "en" {
		t.Fatal("refusal language header", w.Header())
	}
	s.mutateState(t, func(st *State) { delete(st.Rates, "http:member:"+idB) })
	expect(t, a.do("POST", "/auth/logout", nil), 303)
	expect(t, a.do("GET", "/", nil), 200, `<html lang="en">`)

	// Failure: an unsupported value keeps the saved choice, shows the error and the saved radio, nothing else changes.
	signIn(a, "lang-b@example.test")
	w = a.do("POST", "/home/language", url.Values{"lang": {"fr"}})
	expect(t, w, 422, `role="alert"`, `value="ko" checked`, `<html lang="ko">`, "<h1>설정</h1>")
	if locale(idB) != "ko" {
		t.Fatal("invalid value must not be saved")
	}
	// A damaged stored value falls back to Korean without exposing it; the value is distinctive so random page
	// tokens cannot match it by chance.
	s.mutateState(t, func(st *State) { st.Members[idB].Locale = "damaged-locale" })
	w = a.do("GET", "/home/account", nil)
	expect(t, w, 200, `<html lang="en">`) // browser cookie en is the next valid source
	if strings.Contains(w.Body.String(), "damaged-locale") {
		t.Fatal("damaged locale leaked")
	}
	delete(a.cookies, langCookie)
	expect(t, a.do("GET", "/home/account", nil), 200, `<html lang="ko">`, `value="ko" checked`)

	// Anonymous save is refused as an expired session and changes no member.
	anon := newBrowser("192.0.2.33")
	expect(t, anon.do("POST", "/home/language", url.Values{"lang": {"en"}}), 303, "/?n=expired")
	// ?lang= without a cookie applies to that response only and says so; it never becomes the stored choice.
	w = anon.do("GET", "/?lang=en", nil)
	expect(t, w, 200, `<html lang="en">`, `role="status"`, "only")
	if anon.cookies[langCookie] != "" || w.Header().Get(langUnsavedHeader) != "" {
		t.Fatal("hint must not persist or leak the marker", anon.cookies, w.Header())
	}
	expect(t, anon.do("GET", "/", nil), 200, `<html lang="ko">`)
	expect(t, anon.do("GET", "/lang/en?back=https%3A%2F%2Fevil.invalid%2F", nil), 303, "/?lang=en")

	// Language flows create or change no friend, agent, gate, message, connection or legacy pair, and the same
	// deadline reads identically in both languages.
	if after := snapshot(); after != before {
		t.Fatal("language actions changed other state", before, after)
	}
	var codeB string
	s.mutateState(t, func(st *State) { codeB = st.Members[idB].FriendCode })
	signIn(a2, "lang-a@example.test")
	expect(t, a2.do("POST", "/home/friends/request", url.Values{"code": {codeB}, "consent": {"yes"}}), 303)
	var deadline string
	s.mutateState(t, func(st *State) { deadline = clock(st.Friends[pairID(idA, idB)].Exp) })
	w = a2.do("GET", "/home/friends", nil)
	expect(t, w, 200, `<html lang="en">`, deadline)
	english(w.Body.String())
	expect(t, a.do("GET", "/home/friends", nil), 200, `<html lang="ko">`, deadline)
}
