// Approval summaries read only the verified parent envelope; names are reading aids and never approval authority.
package relay

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// gateKinds lists the parent intents a person can judge; any other intent is shown as unverifiable.
var gateKinds = map[string]string{"schedule.query": "일정 가능 시간 확인", "schedule.commit": "일정 확정"}

// gateScopes restates policy() for people: approval passes the gate only and enables nothing else. A member with a
// connected calendar sees calendarGateScope for schedule.query instead (SAR-SCHEDULE-001 화면 문구 원칙).
var gateScopes = map[string]string{
	"schedule.query":  "가능 시간 공개 정책이 없어 승인해도 상대에게 가능 시간이 공개되지 않습니다.",
	"schedule.commit": "승인해도 일정이 확정되거나 실행되지 않습니다. 이 서비스에는 일정 실행 기능이 없습니다.",
}

const calendarGateScope = "승인하면 이 요청의 기간에 대해 바쁜 시간대만 한 번 공유합니다. 일정을 만들거나 바꾸지 않습니다."

type gateField struct{ Label, Value string }

type gateRow struct {
	ID, State, Kind, Exp string
	From, To             humanAgent
}

// gateState is what a person can still act on: a pending gate past its deadline or binding is no longer pending.
func gateState(st *State, g *Gate, now time.Time) (string, *Envelope, bool) {
	e, active := verifiedGate(st, g, now)
	switch {
	case g.State != "pending":
		return g.State, e, false
	case !now.Before(g.Exp):
		return "expired", e, false
	case e == nil:
		return "unavailable", nil, false
	case !active:
		return "revoked", e, false
	}
	return "pending", e, true
}

// gateSummary maps the verified body fields to labels; free text stays separate and is never part of the scope.
// Parse already guarantees these shapes; any mismatch (a new intent or schema) yields no summary, so the page shows
// the "cannot confirm" state instead of panicking.
func gateSummary(e *Envelope) (fields []gateField, text string) {
	ok := true
	span := func(key string) string {
		v, isMap := e.Body[key].(map[string]any)
		ok = ok && isMap
		start, _ := instant(v["start"])
		end, _ := instant(v["end"])
		return clock(start) + " – " + clock(end)
	}
	fields = []gateField{{"요청 유형", gateKinds[e.Intent]}}
	switch e.Intent {
	case "schedule.query":
		minutes, isNumber := e.Body["granularity_min"].(float64)
		ok = ok && isNumber
		fields = append(fields, gateField{"확인할 기간", span("window")}, gateField{"시간 단위", fmt.Sprintf("%d분", int64(minutes))}, gateField{"공개되는 내용", "바쁜 시간대만"})
	case "schedule.commit":
		zone, isZone := e.Body["timezone"].(string)
		commitment, isText := e.Body["commitment"].(string)
		ok = ok && isZone && isText
		fields = append(fields, gateField{"일정", span("slot")}, gateField{"요청에 적힌 시간대", zone})
		text = commitment
	default:
		ok = false
	}
	if !ok {
		return nil, ""
	}
	return fields, text
}

// gateAgent falls back to the already visible ID's short form when the name is not permitted.
func (st *State) gateAgent(owner, id string) humanAgent {
	a := st.visibleAgent(owner, id)
	if a.Short == "" {
		a.Short = shortID(id)
	}
	return a
}

func (st *State) gateRows(owner string, now time.Time) []gateRow {
	rows := []gateRow{}
	for _, g := range st.Gates {
		if g.Owner != owner {
			continue
		}
		// Decided gates keep the server-bound policy label; only pending ones pay for signature checks.
		state, kind := g.State, "확인할 수 없는 요청"
		if g.State == "pending" {
			var e *Envelope
			if state, e, _ = gateState(st, g, now); e != nil {
				kind = gateKinds[e.Intent]
			}
		} else if intent, _, _ := strings.Cut(g.Policy, ":"); gateKinds[intent] != "" {
			kind = gateKinds[intent]
		}
		rows = append(rows, gateRow{ID: g.ID, State: state, Kind: kind, Exp: clock(g.Exp), From: st.gateAgent(owner, g.From), To: st.gateAgent(owner, g.To)})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].ID < rows[j].ID })
	return rows
}
