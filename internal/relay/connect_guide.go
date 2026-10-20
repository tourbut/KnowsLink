// Member connection guide (U8): host choices shown before anything is issued, and per-agent resume rows computed only
// from stored connection, key and remote-install state. Display only: it never marks a host installed or online.
package relay

import "time"

type useChoice struct{ ID, Label, Note string }

// useChoices follow the U8 host table order and wording; hosts outside the target scope are never named.
var useChoices = []useChoice{
	{"web", "웹에서 바로 사용", "설치 필요 없음"},
	{"grok", "Grok Bot (내 컴퓨터)", "연결 안내 있음 · 기존 계정 사용 확인"},
	{"dots", "Dots (원격)", "연결 안내 있음 · 기존 계정 사용 확인"},
	{"claude", "Claude Code 작업 세션 (내 컴퓨터)", "연결 안내 있음 · 시험 환경 확인 (Claude Code 2.1.302, Windows)"},
	{"other", "지원 확인이 필요한 호스트", "Codex, Hermes Agent, 그 밖의 호스트"},
}

// agentResume is one agent's observed connection stage and the single next step that can still work.
type agentResume struct{ State, Next, Link, LinkText string }

const keyManage = "/home/legacy#agents"

func (st *State) agentResume(id string, now time.Time) agentResume {
	a := st.Agents[id]
	// The newest attempt wins; tombstones stay 24h after expiry, so an ended attempt is still explainable.
	var last *Connection
	lastID := ""
	for cid, c := range st.Connections {
		if c.Agent == id && c.Owner == a.Owner && (last == nil || c.Exp.After(last.Exp) || c.Exp.Equal(last.Exp) && cid > lastID) {
			last, lastID = c, cid
		}
	}
	if last != nil && now.Before(last.Exp) {
		page := "/home/connections/" + lastID
		if last.Device {
			page = "/connect/" + lastID
		}
		switch last.State {
		case "waiting":
			return agentResume{"키 등록 전 · 내 컴퓨터의 준비 대기", "기존 시도의 상태를 확인하세요. 연결 수단은 다시 표시하지 않습니다.", page, "연결 상태 확인"}
		case "prepared":
			return agentResume{"키 등록 전 · 지문 확인 필요", "내 컴퓨터에 표시된 지문과 같은지 확인하세요. 다르면 취소하세요.", page, "지문 확인"}
		case "approved":
			return agentResume{"브라우저 동의 완료 · 클라이언트 저장 확인 필요", "시작한 클라이언트에서 마무리와 저장 결과를 확인하세요. 승인만으로는 연결되지 않습니다.", page, "연결 상태 확인"}
		}
	}
	if install, exchanged, returning := st.remoteResume(id, now); install {
		switch {
		case exchanged:
			return agentResume{"키 등록됨 · 원격 앱의 토큰 교환 확인", "앱에서 연결 상태를 확인하세요. 이 화면은 온라인 여부를 확인하지 않습니다.", "", ""}
		case returning:
			return agentResume{"동의 완료 · 앱 복귀 확인 필요", "연결을 시작한 앱에서 토큰 교환 결과를 확인하세요. 앱으로 이동한 것만으로 연결되지 않습니다.", "", ""}
		}
		return agentResume{"확인 필요 · 앱의 토큰 교환을 확인하지 못함", "시작한 앱에서 연결 상태를 먼저 확인하세요. 연결되지 않았다면 앱에서 연결을 다시 시작하세요. 이 에이전트는 자동으로 철회되지 않습니다.", "", ""}
	}
	if activeKeys(a) > 0 {
		return agentResume{"키 등록됨", "앱에서 연결 상태를 확인하세요. 클라이언트가 자격을 저장하지 못했다면 기존 키를 덮어쓰지 말고 그 키를 철회한 뒤 새 키 추가로 다시 연결하세요.", keyManage, "연결·키 관리"}
	}
	if last != nil && last.Device && (last.State == "approved" || last.State == "expired") {
		return agentResume{"확인 필요 · 승인 뒤 클라이언트 저장을 확인하지 못함", "시작한 클라이언트의 저장 결과를 먼저 확인하세요. 쓰지 않는 에이전트라면 직접 철회할 수 있습니다. 자동 철회는 없습니다.", keyManage, "연결·키 관리"}
	}
	if last != nil && last.State != "consumed" {
		ended := "이전 시도 만료"
		if last.State == "cancelled" {
			ended = "이전 시도 취소"
		}
		return agentResume{ended + " · 키 없음", "이전 연결 수단은 다시 쓸 수 없습니다. 연결·키 관리에서 새로 시작하세요.", keyManage, "연결·키 관리"}
	}
	if last != nil {
		return agentResume{"활성 키 없음", "이 에이전트로 다시 연결하려면 연결·키 관리에서 새 키 추가로 시작하세요. 쓰지 않는 에이전트라면 직접 철회할 수 있습니다.", keyManage, "연결·키 관리"}
	}
	return agentResume{"미사용 · 등록된 키 없음", "이 에이전트로 연결하려면 연결·키 관리에서 시작하세요. 쓰지 않는 에이전트라면 직접 철회할 수 있습니다. 자동 철회는 없습니다.", keyManage, "연결·키 관리"}
}

// remoteResume reports whether the agent belongs to a remote install, and whether the app exchanged a token or still
// holds an unexchanged one-minute code.
func (st *State) remoteResume(agent string, now time.Time) (install, exchanged, returning bool) {
	if st.Remote == nil {
		return false, false, false
	}
	for id, in := range st.Remote.Installs {
		if in.Agent != agent {
			continue
		}
		install = true
		for _, t := range st.Remote.Tokens {
			exchanged = exchanged || (t.Install == id && now.Before(t.Exp))
		}
		for _, c := range st.Remote.Codes {
			returning = returning || (c.Install == id && now.Before(c.Exp))
		}
	}
	return install, exchanged, returning
}
