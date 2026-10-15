// Receipt input fixtures exercise account isolation, bounded restoration and escaped native form recovery.
package relay

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func receiptFixture() (*State, time.Time) {
	st := newState()
	st.Members["member-a"] = &Member{Owner: "owner-a"}
	st.Members["member-b"] = &Member{Owner: "owner-b"}
	st.Agents["agent_a"] = &Agent{Owner: "owner-a", Name: strings.Repeat("긴이름", 13), Client: "Node 22 로컬 클라이언트"}
	st.Agents["agent_a2"] = &Agent{Owner: "owner-a", Name: st.Agents["agent_a"].Name, Client: "Node 22 로컬 클라이언트"}
	st.Agents["agent_b"] = &Agent{Owner: "owner-b", Name: "다른 회원의 비공개 이름"}
	return st, time.Date(2026, 10, 15, 1, 0, 0, 0, time.UTC)
}

func draftRequest(t *testing.T, d receiptDraft, now time.Time) *http.Request {
	t.Helper()
	w := httptest.NewRecorder()
	saveReceiptDraft(w, d, now)
	r := httptest.NewRequest("GET", "/home/receipts", nil)
	for _, c := range w.Result().Cookies() {
		if !c.HttpOnly || !c.Secure || c.SameSite != http.SameSiteStrictMode || c.Path != "/" {
			t.Fatal("weak draft cookie")
		}
		r.AddCookie(c)
	}
	return r
}

func TestReceiptDraftIsolationAndBounds(t *testing.T) {
	st, now := receiptFixture()
	d := receiptDraft{Member: "member-a", Agent: "agent_a", ID: firstID, Exp: now.Add(sessionIdle + codeTTL)}
	r := draftRequest(t, d, now)
	for _, c := range []struct {
		name, member, path, agent, id string
	}{
		{"same-member-relogin", "member-a", "/home/receipts", "agent_a", firstID},
		{"other-member", "member-b", "/home/receipts", "", ""},
		{"other-member-old-url", "member-b", "/home/receipts?agent=agent_a&id=" + firstID, "", ""},
		{"edit", "member-a", "/home/receipts?agent=agent_a2&id=" + secondID, "agent_a2", secondID},
		{"back", "member-a", "/home/receipts?agent=agent_a&id=" + firstID, "agent_a", firstID},
		{"foreign-agent", "member-a", "/home/receipts?agent=agent_b&id=" + firstID, "", ""},
		{"oversized", "member-a", "/home/receipts?agent=agent_a&id=" + strings.Repeat("x", 129), "", ""},
		{"control", "member-a", "/home/receipts?agent=agent_a&id=%0A", "", ""},
	} {
		t.Run(c.name, func(t *testing.T) {
			r.URL, _ = url.Parse(c.path)
			v, got := st.receiptView(c.member, r, now)
			if got.Agent != c.agent || got.ID != c.id {
				t.Fatalf("unexpected restored input: %#v", got)
			}
			for _, a := range v["Own"].([]humanAgent) {
				if st.Agents[a.ID].Owner != st.Members[c.member].Owner {
					t.Fatal("foreign agent in selector")
				}
			}
		})
	}
	r.URL, _ = url.Parse("/home/receipts")
	if got := readReceiptDraft(r, d.Exp); got.Member != "" {
		t.Fatal("expired draft restored")
	}
	r = httptest.NewRequest("GET", "/home/receipts", nil)
	r.AddCookie(&http.Cookie{Name: receiptCookie, Value: strings.Repeat("x", 1025)})
	if got := readReceiptDraft(r, now); got.Member != "" {
		t.Fatal("oversized cookie restored")
	}
}

func TestReceiptFormRecoveryFixtures(t *testing.T) {
	st, now := receiptFixture()
	for _, name := range []string{"first", "invalid", "unavailable", "limited", "success", "revoked"} {
		t.Run(name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/home/receipts?agent=agent_a&id="+firstID, nil)
			v, _ := st.receiptView("member-a", r, now)
			switch name {
			case "first":
				v["ID"] = ""
			case "invalid":
				v["ID"], v["Problem"] = `"><script>bad()</script>`, "요청 ID 형식이 올바르지 않습니다."
			case "unavailable":
				v["Problem"], v["Focus"] = receiptUnavailable, "receipt-problem"
			case "limited":
				v["Problem"], v["Focus"] = limited(now.Add(time.Minute)), "receipt-problem"
			case "success":
				v["Transport"], v["Completion"], v["ReplyState"], v["Focus"] = transportNames["delivered"], "수신 클라이언트가 보관했습니다.", "아직 관련 답장이 없습니다.", "receipt-result"
				v["From"], v["To"], v["Exp"], v["TTL"] = "agent_a", "agent_b", clock(now), 180
			case "revoked":
				st.Agents["agent_a"].Revoked = true
				v, _ = st.receiptView("member-a", r, now)
				if v["Agent"] != "agent_a" || v["Focus"] != "receipt-agent" {
					t.Fatal("revoked selection silently replaced")
				}
			}
			body := renderFixture(t, uiFixture{name, "receipt", v})
			if strings.Contains(body, "<script>") || strings.Contains(body, "다른 회원의 비공개 이름") {
				t.Fatal("untrusted HTML or foreign name")
			}
			if len(formSignatures(body)) != 2 || !strings.Contains(body, `name="agent"`) || !strings.Contains(body, `name="id" value="`) {
				t.Fatal("editable form lost")
			}
			if strings.Count(body, " autofocus") != 1 {
				t.Fatal("focus target not unique")
			}
			if dir := os.Getenv("KNOWSLINK_UI_FIXTURE_DIR"); dir != "" {
				if err := os.WriteFile(filepath.Join(dir, "receipt-"+name+".html"), []byte(body), 0600); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}
