// Session-only agent and relationship screens reuse the owner state machine and never reveal another member's identity.
package relay

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type memberPair struct {
	*Pair
	Incoming bool
	// Own and Other name the member's live agent and the counterpart for a manual re-invite; Own is "" when none.
	Own, Other string
	// MineAgent and PeerAgent label Mine and Peer; a non-friend counterpart stays an unnamed identifier.
	MineAgent, PeerAgent humanAgent
}

// Deadline shows the invite deadline in the same KST form as connection deadlines.
func (p memberPair) Deadline() string {
	if p.Exp.IsZero() {
		return "없음"
	}
	return clock(p.Exp)
}

type memberAgent struct {
	ID, Status          string
	Name, Client, Short string // the same name-first label as the person pages
	Keys                []memberKey
	KeyFull             bool // the key record cap is reached; only a new agent can take a new key
}
type memberKey struct{ Kid, Fingerprint, Status string }

// ActiveKeys is the active key count the collapsed agent row shows.
func (a memberAgent) ActiveKeys() int {
	n := 0
	for _, k := range a.Keys {
		if k.Status == "활성" {
			n++
		}
	}
	return n
}

// Mine and Peer orient a relationship row from the member's side: only an incoming invite has the member as recipient.
func (p memberPair) Mine() string {
	if p.Incoming {
		return p.Recipient
	}
	return p.Inviter
}
func (p memberPair) Peer() string {
	if p.Incoming {
		return p.Inviter
	}
	return p.Recipient
}

func fingerprint(public []byte) string {
	sum := sha256.Sum256(public)
	return "SHA256:" + base64.RawURLEncoding.EncodeToString(sum[:])
}
func newMemberAgentID() string {
	return fmt.Sprintf("agent_%x", sha256.Sum256([]byte(randomToken())))[:28]
}
func ownAgents(st *State, owner string) []memberAgent {
	agents := []memberAgent{}
	for id, a := range st.Agents {
		if a.Owner != owner {
			continue
		}
		label := humanAgentView(id, a)
		v := memberAgent{ID: id, Status: "미연결", Name: label.Name, Client: label.Client, Short: label.Short, KeyFull: len(a.Keys) >= agentKeyRecords}
		if activeKeys(a) > 0 {
			v.Status = "연결 완료"
		}
		if a.Revoked {
			v.Status = "철회"
		}
		for kid, k := range a.Keys {
			state := "활성"
			if k.Revoked {
				state = "철회"
			}
			v.Keys = append(v.Keys, memberKey{kid, fingerprint(k.Public), state})
		}
		sort.Slice(v.Keys, func(i, j int) bool { return v.Keys[i].Kid < v.Keys[j].Kid })
		agents = append(agents, v)
	}
	sort.Slice(agents, func(i, j int) bool { return agents[i].ID < agents[j].ID })
	labels := make([]humanAgent, len(agents))
	for i, a := range agents {
		labels[i] = humanAgent{ID: a.ID, Short: a.Short}
	}
	for i, label := range distinguish(labels) {
		agents[i].Short = label.Short
	}
	return agents
}
func safeProblem(err error) string {
	switch err.Error() {
	case "capacity":
		return "운영 한도에 도달했습니다. 기존 키·관계는 철회할 수 있습니다."
	case "unsupported_client":
		return "미지원 클라이언트입니다. Node 22 로컬 클라이언트를 선택하세요."
	case "reauth_required":
		return "이 작업은 5분 안의 로그인 신원 재확인이 필요합니다. 아래 ‘로그인 다시 확인’으로 신원을 확인한 뒤 이 작업을 시작한 화면에서 다시 시도하세요."
	case "invalid_auth":
		return "연결 수단이 만료·취소·사용되었거나 유효하지 않습니다. 아직 새 키가 연결되지 않았습니다."
	case "agent_rate":
		r, _ := err.(agentRefused)
		return "오늘 새 에이전트를 많이 만들었습니다. " + clock(r.Retry) + " 뒤 다시 연결하세요. 기존 에이전트는 그대로 쓸 수 있습니다."
	case "service_capacity":
		return "지금은 서비스 용량 때문에 새 에이전트를 만들 수 없습니다. 기존 에이전트는 그대로 쓸 수 있습니다. 나중에 다시 연결하세요."
	default:
		return "처리할 수 없습니다. 대상과 현재 권한·상태를 확인하세요. 새 권한은 생성되지 않았습니다."
	}
}

const keyRecordsFull = "이 agent에는 새 키를 더 연결할 수 없습니다(키 기록 보호 상한). 키를 철회하거나 기다려도 이 agent에 새 키 공간은 생기지 않습니다. 새 agent를 만들어 따로 연결하세요. 활성 친구에는 새 agent도 연결됩니다. 기존 개별 연결만 있다면 사람 친구 요청을 새로 수락하세요. 기존 키와 관계는 철회하기 전까지 그대로입니다."

// connectLimit explains a connect capacity refusal; create reports whether a replacement agent can be made now.
func connectLimit(st *State, owner, agent, mode string, now time.Time) (problem string, create bool) {
	a := st.Agents[agent]
	switch {
	case a != nil && len(a.Keys) >= agentKeyRecords:
		if err := st.agentLimit(owner, now); err != nil {
			return keyRecordsFull + " 지금은 새 agent도 만들 수 없습니다. " + safeProblem(err), false
		}
		return keyRecordsFull, true
	case a != nil && mode == "register" && activeKeys(a) >= 3:
		return "활성 키는 agent당 3개입니다. 쓰지 않는 키를 철회하거나 회전 방식으로 연결하세요.", false
	}
	return safeProblem(fault("capacity")), false
}

// hostNames are default agent names (A-04): the hosts a client may name when it starts, plus remote Dots.
var hostNames = map[string]string{"claude": "Claude Code", "codex": "Codex", "grok": "Grok Bot", "dots": "Dots"}

// defaultName is the host name, numbered from 2 while owner's live agents or approved attempts hold it; an unknown
// host counts "내 에이전트 N" from 1. It never uses device, account or email details, since friends see the name.
func (st *State) defaultName(owner, host, lang string) string {
	taken := map[string]bool{}
	for _, a := range st.Agents {
		if a.Owner == owner && !a.Revoked {
			taken[a.Name] = true
		}
	}
	for _, c := range st.Connections {
		if c.Device && c.Owner == owner && c.State == "approved" {
			taken[c.Name] = true
		}
	}
	base, n := hostNames[host], 1
	name := base
	if base == "" {
		base = "내 에이전트"
		if lang == "en" {
			base = "My agent"
		}
		name = base + " 1"
	}
	for taken[name] {
		n++
		name = base + " " + strconv.Itoa(n)
	}
	return name
}

// ownAttempt is one client-started connection still waiting for its key on the agents page (M-02).
type ownAttempt struct{ ID, Client, Started, State string }

// ownAttempts lists the member's keyless client attempts until the sweep removes them a day after their deadline.
func (st *State) ownAttempts(owner string) []ownAttempt {
	out := []ownAttempt{}
	for id, c := range st.Connections {
		if !c.Device || c.Owner != owner || (c.State != "prepared" && c.State != "approved" && c.State != "expired") {
			continue
		}
		client := hostNames[c.Host]
		if client == "" {
			client = "Node 로컬 클라이언트"
		}
		out = append(out, ownAttempt{id, client, c.Exp.Add(-connectionTTL).In(kst).Format("15:04"), c.State})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

// inviteNotice tells the inviter whether the submit made a new pending invite or kept the current one unchanged.
func inviteNotice(before string) string {
	switch before {
	case "pending":
		return "invite-pending"
	case "active":
		return "invite-active"
	}
	return "invited"
}
func (s *Service) agentRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /v1/connect/start", s.deviceAPI(true))
	mux.HandleFunc("POST /v1/connect/poll", s.deviceAPI(false))
	mux.HandleFunc("GET /connect/{id}", s.devicePage)
	mux.HandleFunc("POST /home/device-confirm", s.memberAction("device-confirm"))
	mux.HandleFunc("POST /home/agents", s.memberAction("create"))
	mux.HandleFunc("POST /home/connect", s.memberAction("connect"))
	mux.HandleFunc("GET /home/connections/{id}", s.connectionPage)
	for _, path := range []string{"invite-deny", "confirm", "cancel", "attempt-remove", "key-revoke", "agent-revoke", "invites", "invite-decision", "unpair"} {
		mux.HandleFunc("POST /home/"+path, s.memberAction(path))
	}
	for _, path := range []string{"info", "prepare", "complete"} {
		mux.HandleFunc("POST /v1/connect/"+path, s.connectAPI(path))
	}
}
func (s *Service) memberAction(path string) http.HandlerFunc {
	denying := path == "invite-deny"
	if denying {
		path = "invite-decision"
	}
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 8192)
		parseErr := r.ParseForm()
		value, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
			clean := path == "cancel" || path == "attempt-remove" || path == "key-revoke" || path == "agent-revoke" || path == "unpair" || (path == "invite-decision" && r.FormValue("decision") == "deny")
			rates := memberRate
			if clean {
				rates = cleanupRate
			}
			id, session, retry, e := s.memberHit(w, st, r, now, rates)
			if !retry.IsZero() {
				return refusal{429, limited(retry), false, false}, nil
			}
			if e != nil {
				return refusal{401, "로그인 세션이 유효하지 않습니다. 다시 로그인하세요.", false, false}, nil
			}
			owner := st.Members[id].Owner
			// Return refusals as values: rejected work must commit its request budget.
			fail := func(e error) (any, error) {
				return refusal{statusFor(e), safeProblem(e), e.Error() == "reauth_required", false}, nil
			}
			if parseErr != nil || (denying && r.PostForm.Get("decision") != "deny") {
				return fail(fault("invalid_schema"))
			}
			if path == "create" || path == "connect" || path == "confirm" || path == "device-confirm" || path == "key-revoke" || path == "agent-revoke" {
				if now.Sub(session.Verified) >= reauthWindow {
					return fail(fault("reauth_required"))
				}
			}
			agent := r.FormValue("agent")
			switch path {
			case "device-confirm":
				if st.Members[id].Issuer != googleIssuer {
					return fail(fault("invalid_auth"))
				}
				name, valid := validName(r.FormValue("name"))
				if !valid {
					c := st.Connections[r.FormValue("connection")]
					if c == nil || !c.Device {
						return fail(fault("invalid_auth"))
					}
					public, _ := base64.RawURLEncoding.DecodeString(c.Public)
					problem, kept := nameRefusal(r.FormValue("name"))
					return map[string]any{"FormTemplate": "device", "FormStatus": 422, "Title": "내 클라이언트 연결", "ConnectionID": r.FormValue("connection"), "Name": kept, "State": c.State, "Fingerprint": fingerprint(public), "Exp": clock(c.Exp), "Problem": problem}, nil
				}
				if e := st.approveDevice(r.FormValue("connection"), owner, name, now); e != nil {
					return fail(e)
				}
				return "/connect/" + r.FormValue("connection"), nil
			case "create":
				name, valid := validName(r.FormValue("name"))
				if !valid {
					view, e := st.humanView(id, session, "agents", r, now)
					if e != nil {
						return fail(e)
					}
					problem, kept := nameRefusal(r.FormValue("name"))
					view["FormTemplate"], view["FormStatus"], view["Problem"], view["NewName"] = "human-agents", 422, problem, kept
					return view, nil
				}
				// A repeated submit of the same form (double click, resend) returns the agent it already made.
				attempt := r.FormValue("attempt")
				scope := agentCreateKey(owner, attempt)
				if len(attempt) < 16 || len(attempt) > 96 {
					scope = ""
				} else if made := st.Agents[st.Idempotency[scope]]; made != nil && made.Owner == owner && !made.Revoked {
					return "/home/agents#agent-" + st.Idempotency[scope], nil
				}
				if e := st.agentLimit(owner, now); e != nil {
					return fail(e)
				}
				agent = newMemberAgentID() // Avoid email-derived IDs and cross-member name collisions.
				st.addAgent(agent, &Agent{Name: name, Owner: owner, Keys: map[string]*Key{}}, now)
				if scope != "" {
					st.Idempotency[scope] = agent
				}
				return "/home/agents#agent-" + agent, nil
			case "connect":
				t, e := st.beginConnection(owner, agent, r.FormValue("client"), r.FormValue("mode"), now)
				if e != nil && e.Error() == "capacity" {
					problem, create := connectLimit(st, owner, agent, r.FormValue("mode"), now)
					return refusal{status: 409, problem: problem, create: create}, nil
				}
				if e != nil {
					return fail(e)
				}
				return map[string]any{"Title": "에이전트 연결", "Token": t, "ConnectionID": hashToken(t), "Agent": agent, "Name": agentName(st.Agents[agent]), "Short": shortID(agent), "Owner": owner, "Client": supportedClient, "Mode": r.FormValue("mode"), "State": "waiting", "Exp": clock(now.Add(connectionTTL))}, nil
			case "attempt-remove":
				// A client attempt has no key or agent yet, so removing it revokes nothing and keeps no record (A-02).
				c := st.Connections[r.FormValue("connection")]
				if c == nil || !c.Device || c.Owner != owner || c.State == "consumed" {
					return fail(fault("sender_not_allowed"))
				}
				delete(st.Connections, r.FormValue("connection"))
				return "/home/agents", nil
			case "confirm", "cancel":
				c := st.Connections[r.FormValue("connection")]
				if c == nil || c.Owner != owner || !now.Before(c.Exp) || c.State == "consumed" || c.State == "cancelled" {
					return fail(fault("invalid_auth"))
				}
				if path == "cancel" {
					c.State = "cancelled"
					if c.Agent != "" {
						return "/home/agents#agent-" + c.Agent, nil
					}
					return "/home/agents", nil
				}
				if c.Device || c.State != "prepared" {
					return fail(fault("sender_not_allowed"))
				}
				c.State = "approved"
				return "/home/connections/" + r.FormValue("connection"), nil
			case "agent-revoke":
				a := st.Agents[agent]
				if a == nil || a.Owner != owner {
					return fail(fault("sender_not_allowed"))
				}
				if !a.Revoked {
					a.Revoked = true
					a.Changed = now
				}
				a.Credential = ""
				revokeKeys(a, now)
				st.sweep(now)
				return "/home/agents", nil
			default:
				c := command{Agent: agent, Target: r.FormValue("target"), Kid: r.FormValue("kid"), Decision: r.FormValue("decision")}
				if path == "invite-decision" || path == "unpair" {
					generation, e := strconv.ParseInt(r.PostFormValue("generation"), 10, 64)
					p := st.Pairs[pairID(c.Agent, c.Target)]
					if e != nil || p == nil || p.Generation != generation {
						return fail(fault("sender_not_allowed"))
					}
				}
				before := ""
				if p := st.Pairs[pairID(c.Agent, c.Target)]; p != nil {
					before = p.State
				}
				_, e := st.operateAs(path, owner, "owner", c, now, false)
				if e != nil {
					return fail(e)
				}
				if path == "invites" {
					return "/home?n=" + inviteNotice(before), nil
				}
			}
			return "/home", nil
		})
		if err != nil {
			if err.Error() == "invalid_auth" {
				expired(w, r)
			} else {
				refusedTo(w, "내 에이전트", refusal{503, "일시적으로 처리할 수 없습니다. 잠시 뒤 다시 시도하세요.", false, false}, memberBack(path))
			}
			return
		}
		switch v := value.(type) {
		case refusal:
			refusedTo(w, "내 에이전트", v, memberBack(path))
		case string:
			http.Redirect(w, r, v, 303)
		default:
			view := value.(map[string]any)
			name, status := "connection", 200
			if template, ok := view["FormTemplate"].(string); ok {
				name = template
				status = view["FormStatus"].(int)
			}
			render(w, status, name, view)
		}
	}
}

// agentCreateKey starts with "#", which no agent ID (agentPattern) may start with, so it never equals a message
// idempotency key ("from/key"). It stays printable because Postgres JSONB rejects NUL.
func agentCreateKey(owner, attempt string) string { return "#agent-create/" + owner + "/" + attempt }

// memberBack is the internal list a failed member action returns to: legacy pair actions to the pair list, the rest
// to the agent list. No caller-supplied return URL is ever used.
func memberBack(path string) string {
	switch path {
	case "invites", "invite-decision", "unpair":
		return "/home/legacy#pairs"
	}
	return "/home/agents"
}
func (s *Service) connectionPage(w http.ResponseWriter, r *http.Request) {
	value, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
		id, _, retry, e := s.memberHit(w, st, r, now, memberRate)
		if !retry.IsZero() {
			return refusal{429, limited(retry), false, false}, nil
		}
		if e != nil {
			return nil, e
		}
		c := st.Connections[r.PathValue("id")]
		if c == nil || c.Owner != st.Members[id].Owner {
			return refusal{403, "연결을 조회할 수 없습니다.", false, false}, nil
		}
		public, _ := base64.RawURLEncoding.DecodeString(c.Public)
		var unsure *agentResume
		// An approved attempt past its deadline keeps its state until cleanup (A-02) but asks for the same check.
		if a := st.Agents[c.Agent]; c.Device && (c.State == "expired" || c.State == "approved" && !now.Before(c.Exp)) && a != nil && activeKeys(a) == 0 {
			unsure = &deviceUnsure // approved, then expired: the same check the agent card asks for (U8 승인 후 응답 유실)
		}
		state := c.State
		if unsure != nil {
			state = "expired" // past its deadline the page offers no completion command or cancel
		}
		return map[string]any{"Unsure": unsure, "Title": "에이전트 연결", "ConnectionID": r.PathValue("id"), "Agent": c.Agent, "Name": agentName(st.Agents[c.Agent]), "Short": shortID(c.Agent), "Client": c.Client, "Mode": c.Mode, "State": state, "Fingerprint": fingerprint(public), "Exp": clock(c.Exp)}, nil
	})
	if err != nil {
		if err.Error() == "invalid_auth" {
			expired(w, r)
		} else {
			refused(w, "연결 확인", refusal{503, "일시적으로 처리할 수 없습니다. 잠시 뒤 다시 시도하세요.", false, false})
		}
		return
	}
	if v, ok := value.(refusal); ok {
		refusedTo(w, "연결 확인", v, "/home/agents")
		return
	}
	render(w, 200, "connection", value.(map[string]any))
}
func (s *Service) connectAPI(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var c struct{ Token, Client, Kid, Public, Proof string }
		parseErr := decodeCommand(w, r, &c)
		value, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
			conn, e := st.connection(c.Token, now)
			rates := anonymousRate(s.clientIP(r))
			if e == nil {
				for id, m := range st.Members {
					if m.Owner == conn.Owner {
						rates = memberRate(id)
						break
					}
				}
			}
			if ok, retry := s.requestHit(st, r, now, rates...); !ok {
				return retry, nil
			}
			if parseErr != nil {
				return map[string]string{"error": parseErr.Error()}, nil
			}
			if e != nil {
				return map[string]string{"error": e.Error()}, nil
			}
			var result any
			switch path {
			case "info":
				if c.Client != conn.Client {
					e = fault("sender_not_allowed")
				} else {
					result = map[string]string{"owner": conn.Owner, "agent": conn.Agent, "client": conn.Client, "mode": conn.Mode, "state": conn.State, "exp": conn.Exp.Format(time.RFC3339)}
				}
			case "prepare":
				result, e = st.prepareConnection(c.Token, c.Client, c.Kid, c.Public, c.Proof, now)
			case "complete":
				result, e = st.completeConnection(c.Token, c.Client, c.Proof, now)
			}
			if e != nil {
				return map[string]string{"error": e.Error()}, nil
			}
			return result, nil
		})
		connectionResponse(w, value, err)
	}
}
