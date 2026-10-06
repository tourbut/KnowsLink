// A member looks up one owned receipt by ID; text bodies and other members' identities never appear in this view.
package relay

import (
	"net/http"
	"time"
)

var transportNames = map[string]string{
	"queued":              "queued — 대기 중. 상대 수신이 확인되지 않았습니다.",
	"leased":              "leased — 수신 확인 중. 아직 ACK가 없습니다.",
	"delivered":           "delivered — 수신 클라이언트의 저장·ACK가 확인되었습니다.",
	"failed:expired":      "failed:expired — 기한 만료. 상대가 오프라인이거나 수동 pull하지 않았을 수 있습니다. 연결·관계를 확인하고 새 key로 명시 송신하세요.",
	"failed:revoked":      "failed:revoked — 키·agent·관계 권한이 끝났습니다. 현재 키와 새 관계 수락을 확인하세요. 옛 요청은 복구되지 않습니다.",
	"failed:max_attempts": "failed:max_attempts — 수신 확인 3회가 끝났습니다. 클라이언트 상태를 확인한 뒤 새 key로 명시 송신하세요.",
}

func (s *Service) memberReceipt(w http.ResponseWriter, r *http.Request) {
	v, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
		member, _, retry, err := s.memberHit(st, r, now, memberRate)
		if !retry.IsZero() {
			return refusal{status: 429, problem: limited(retry)}, nil
		}
		if err != nil {
			return nil, err
		}
		owner := st.Members[member].Owner
		m := st.Messages[r.URL.Query().Get("id")]
		agent := r.URL.Query().Get("agent")
		a := st.Agents[agent]
		if a == nil || a.Owner != owner || m == nil || (m.Receipt.From != agent && m.Receipt.To != agent) {
			return refusal{status: 403, problem: "이 agent의 요청을 조회할 수 없습니다. 자기 agent와 요청 ID를 확인하세요."}, nil
		}
		completion := m.Completion
		if completion == "" {
			completion = "아직 처리 결과가 없습니다."
		}
		replyState := "아직 관련 답장이 없습니다."
		if reply := st.Messages[m.ReplyID]; reply != nil {
			replyState = transportNames[reply.Receipt.State]
		}
		if !now.Before(m.Receipt.Exp) && m.ReplyID == "" {
			replyState = "답장 기한 만료. 새 연결 확인 요청을 명시적으로 보내세요."
		}
		return map[string]any{"Title": "연결 확인·receipt", "Agent": agent, "ID": m.Receipt.ID, "From": m.Receipt.From, "To": m.Receipt.To, "Parent": m.Parent, "Reply": m.ReplyID, "Transport": transportNames[m.Receipt.State], "Completion": completion, "ReplyState": replyState, "Exp": clock(m.Receipt.Exp), "TTL": m.Receipt.Exp.Sub(m.Receipt.Accepted).Seconds()}, nil
	})
	if err != nil {
		if err.Error() == "invalid_auth" {
			expired(w, r)
		} else {
			refused(w, "receipt", refusal{status: 503, problem: "상태를 확인할 수 없습니다. 잠시 뒤 수동으로 다시 조회하세요."})
		}
		return
	}
	if refusal, ok := v.(refusal); ok {
		refused(w, "receipt", refusal)
		return
	}
	render(w, 200, "receipt", v.(map[string]any))
}
