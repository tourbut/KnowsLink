// Receipt drafts are bounded browser input, never authority; restoration requires the same authenticated member.
package relay

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

const receiptCookie = "__Host-kl_receipt"

type receiptDraft struct {
	Member, Agent, ID string
	Exp               time.Time
}

// Admission can refuse before the route runs; only its authenticated DB view may be rendered with input.
func receiptFailure(w http.ResponseWriter, r *http.Request, status int, problem string, view map[string]any, draft receiptDraft, now time.Time) bool {
	if r.Method != "GET" || r.URL.Path != "/home/receipts" {
		return false
	}
	if view == nil {
		view = map[string]any{"Title": "요청 상태 조회", "Unavailable": true}
		problem += " 저장된 입력은 로그인 확인이 가능해지면 복원됩니다."
	}
	view["Problem"], view["Focus"] = problem, "receipt-problem"
	delete(view, "Transport")
	if draft.Member != "" {
		saveReceiptDraft(w, draft, now)
	}
	render(w, status, "receipt", view)
	return true
}

func safeReceiptInput(value string) bool {
	return len(value) <= 128 && utf8.ValidString(value) && !strings.ContainsFunc(value, func(r rune) bool {
		return unicode.IsControl(r) || unicode.Is(unicode.Cf, r)
	})
}

func readReceiptDraft(r *http.Request, now time.Time) receiptDraft {
	raw := readCookie(r, receiptCookie)
	if len(raw) > 1024 {
		return receiptDraft{}
	}
	b, err := base64.RawURLEncoding.DecodeString(raw)
	var d receiptDraft
	if err != nil || json.Unmarshal(b, &d) != nil || !now.Before(d.Exp) || d.Exp.After(now.Add(sessionIdle+codeTTL)) ||
		!safeReceiptInput(d.Member) || !safeReceiptInput(d.Agent) || !safeReceiptInput(d.ID) {
		return receiptDraft{}
	}
	return d
}

func saveReceiptDraft(w http.ResponseWriter, d receiptDraft, now time.Time) {
	b, _ := json.Marshal(d)
	setCookie(w, receiptCookie, base64.RawURLEncoding.EncodeToString(b), int(d.Exp.Sub(now).Seconds()))
}

func (st *State) receiptView(member string, r *http.Request, now time.Time) (map[string]any, receiptDraft) {
	d := readReceiptDraft(r, now)
	if d.Member != member {
		d = receiptDraft{}
	}
	d.Member, d.Exp = member, now.Add(sessionIdle+codeTTL)
	q := r.URL.Query()
	if q.Has("agent") || q.Has("id") {
		d.Agent, d.ID = q.Get("agent"), q.Get("id")
	}
	v := map[string]any{"Title": "요청 상태 조회", "Member": member, "Status": 200, "Focus": "receipt-id"}
	if !safeReceiptInput(d.Agent) {
		d.Agent, d.ID = "", ""
		v["Status"], v["Problem"] = 422, "입력이 너무 길거나 허용되지 않는 문자가 있습니다. 내 에이전트를 다시 선택하세요."
	}
	if !safeReceiptInput(d.ID) {
		d.ID = ""
		v["Status"], v["Problem"] = 422, "입력이 너무 길거나 허용되지 않는 문자가 있습니다. 요청 ID를 다시 입력하세요."
	}
	if d.Agent == "" {
		v["Focus"] = "receipt-agent"
	}
	owner := st.Members[member].Owner
	a := st.Agents[d.Agent]
	if d.Agent != "" && (a == nil || a.Owner != owner) {
		// An old URL in another account must not echo the previous account's input.
		d.Agent, d.ID = "", ""
		v["Status"], v["Problem"], v["Focus"] = 403, receiptUnavailable, "receipt-agent"
	}
	own := st.humanAgents(owner)
	if a != nil && a.Owner == owner && a.Revoked {
		own = append(own, humanAgentView(d.Agent, a))
		v["Status"], v["Problem"], v["Focus"] = 403, "선택한 내 에이전트가 철회되었습니다. 사용할 내 에이전트를 다시 선택하세요.", "receipt-agent"
	}
	v["Own"], v["Agent"], v["ID"] = distinguish(own), d.Agent, d.ID
	return v, d
}
