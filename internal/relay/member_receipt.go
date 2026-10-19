// A member looks up one owned receipt by ID; text bodies and other members' identities never appear in this view.
package relay

import (
	"net/http"
	"strings"
	"time"
)

const receiptUnavailable = "이 요청의 상태를 확인할 수 없습니다. 내 에이전트와 요청 ID를 확인하세요."

var transportNames = map[string]string{
	"queued":              "대기 중. 상대 수신이 확인되지 않았습니다.",
	"leased":              "수신 확인 중. 아직 ACK가 없습니다.",
	"delivered":           "전달됨. 수신 클라이언트의 저장·ACK가 확인되었습니다.",
	"failed:expired":      "전달 기한 만료. 상대가 오프라인이거나 아직 받지 않았을 수 있습니다. 현재 연결·친구를 확인하고 필요하면 새 메시지로 보내세요.",
	"failed:revoked":      "전달 권한이 끝났습니다. 현재 연결과 친구 수락을 확인하세요. 옛 요청은 복구되지 않습니다.",
	"failed:max_attempts": "수신 확인 3회가 끝났습니다. 클라이언트 상태를 확인한 뒤 필요하면 새 메시지로 보내세요.",
}

func (s *Service) memberReceipt(w http.ResponseWriter, r *http.Request) {
	var view map[string]any
	var draft receiptDraft
	var checkedAt time.Time
	_, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
		member, _, retry, err := s.memberHit(st, r, now, memberRate)
		if err != nil {
			if !retry.IsZero() {
				view = map[string]any{"Title": "요청 상태 조회", "Status": 429, "Problem": limited(retry), "Focus": "receipt-problem", "Unavailable": true}
				return nil, nil
			}
			return nil, err
		}
		view, draft = st.receiptView(member, r, now)
		checkedAt = now
		if !retry.IsZero() {
			view["Status"], view["Problem"], view["Focus"] = 429, limited(retry), "receipt-problem"
			return nil, nil
		}
		if view["Problem"] != nil || (!r.URL.Query().Has("id") && draft.ID == "") {
			return nil, nil
		}
		if draft.Agent == "" {
			view["Status"], view["Problem"], view["Focus"] = 422, "내 에이전트를 선택하세요.", "receipt-agent"
			return nil, nil
		}
		if !uuidPattern.MatchString(strings.TrimSpace(draft.ID)) {
			view["Status"], view["Problem"] = 422, "요청 ID 형식이 올바르지 않습니다. 받은 요청 ID 또는 답장 ID를 확인하세요."
			return nil, nil
		}
		m := st.Messages[strings.TrimSpace(draft.ID)]
		if m == nil || (m.Receipt.From != draft.Agent && m.Receipt.To != draft.Agent) {
			view["Status"], view["Problem"], view["Focus"] = 403, receiptUnavailable, "receipt-problem"
			return nil, nil
		}
		completion := "처리 결과가 있습니다. 기술 상세에서 확인하세요."
		switch m.Completion {
		case "":
			completion = "아직 처리 결과가 없습니다."
		case "received":
			completion = "수신 클라이언트가 보관했습니다. 읽음·업무 실행 완료는 아닙니다."
		case "reply_received":
			completion = "관련 답장을 수신 클라이언트가 보관했습니다."
		case "failed:expired":
			completion = "처리 결과를 받기 전에 기한이 끝났습니다."
		}
		replyState := "아직 관련 답장이 없습니다."
		if reply := st.Messages[m.ReplyID]; reply != nil {
			replyState = transportNames[reply.Receipt.State]
		}
		if !now.Before(m.Receipt.Exp) && m.ReplyID == "" {
			replyState = "답장 기한 만료. 새 연결 확인 요청을 명시적으로 보내세요."
		}
		for key, value := range map[string]any{"From": m.Receipt.From, "To": m.Receipt.To, "Parent": m.Parent, "Reply": m.ReplyID, "Transport": transportNames[m.Receipt.State], "State": m.Receipt.State, "Completion": completion, "RawCompletion": m.Completion, "ReplyState": replyState, "Exp": clock(m.Receipt.Exp), "TTL": m.Receipt.Exp.Sub(m.Receipt.Accepted).Seconds(), "Focus": "receipt-result"} {
			view[key] = value
		}
		return nil, nil
	})
	if err != nil {
		if err.Error() == "invalid_auth" {
			expired(w, r)
			return
		}
		if view == nil {
			// Without an authenticated DB read, retain the browser draft but reveal no previous account input.
			view = map[string]any{"Title": "요청 상태 조회", "Unavailable": true}
		}
		delete(view, "Transport")
		view["Status"], view["Problem"], view["Focus"] = 503, "상태를 확인할 수 없습니다. 잠시 뒤 같은 화면에서 다시 조회하세요. 저장된 입력은 로그인 확인이 가능해지면 복원됩니다.", "receipt-problem"
	}
	if draft.Member != "" {
		saveReceiptDraft(w, draft, checkedAt)
	}
	render(w, view["Status"].(int), "receipt", view)
}
