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
}

// Deadline shows the invite deadline in the same KST form as connection deadlines.
func (p memberPair) Deadline() string {
	if p.Exp.IsZero() {
		return "없음"
	}
	return clock(p.Exp)
}

type memberAgent struct {
	ID, Status string
	Keys       []memberKey
	KeyFull    bool // the key record cap is reached; only a new agent can take a new key
}
type memberKey struct{ Kid, Fingerprint, Status string }

func fingerprint(public []byte) string {
	sum := sha256.Sum256(public)
	return "SHA256:" + base64.RawURLEncoding.EncodeToString(sum[:])
}
func ownAgents(st *State, owner string) []memberAgent {
	agents := []memberAgent{}
	for id, a := range st.Agents {
		if a.Owner != owner {
			continue
		}
		v := memberAgent{ID: id, Status: "미연결", KeyFull: len(a.Keys) >= agentKeyRecords}
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
	return agents
}
func safeProblem(err error) string {
	switch err.Error() {
	case "capacity":
		return "운영 한도에 도달했습니다. 기존 키·관계는 철회할 수 있습니다."
	case "unsupported_client":
		return "미지원 클라이언트입니다. Node 22 로컬 클라이언트를 선택하세요."
	case "reauth_required":
		return "이 작업은 5분 안의 이메일 재확인이 필요합니다. 홈에서 이메일을 다시 확인하세요."
	case "invalid_auth":
		return "연결 수단이 만료·취소·사용되었거나 유효하지 않습니다. 아직 새 키가 연결되지 않았습니다."
	default:
		return "처리할 수 없습니다. 대상과 현재 권한·상태를 확인하세요. 새 권한은 생성되지 않았습니다."
	}
}

// Record and active-limit refusals name the actual limit and the next step that can work; none promises an exact
// cleanup time, payment, or that revoking a key or waiting frees key space on a live agent (C1).
var agentLimits = map[string]string{
	"active":  "활성 agent는 회원당 5개입니다. 새 agent를 만들려면 쓰지 않는 agent를 먼저 철회하세요. 철회한 agent의 키와 관계는 끝나며 되돌릴 수 없습니다.",
	"records": "철회한 agent 기록을 보존하는 중이라 지금은 새 agent를 만들 수 없습니다. 기존 agent와 키는 그대로입니다. 필요하면 쓰지 않는 agent를 철회하세요. 철회 기록은 최소 24시간 보존된 뒤 정리됩니다. 정리되어 홈 목록에서 사라진 뒤 다시 시도하세요.",
}

const keyRecordsFull = "이 agent에는 새 키를 더 연결할 수 없습니다(키 기록 보호 상한). 키를 철회하거나 기다려도 이 agent에 새 키 공간은 생기지 않습니다. 새 agent를 만들어 따로 연결하고, 각 상대와 새로 초대·수락하세요. 기존 키와 관계는 철회하기 전까지 그대로입니다."

// connectLimit explains a connect capacity refusal; create reports whether a replacement agent can be made now.
func connectLimit(st *State, owner, agent, mode string) (problem string, create bool) {
	a := st.Agents[agent]
	switch {
	case a != nil && len(a.Keys) >= agentKeyRecords:
		if limit := st.agentLimit(owner); limit != "" {
			return keyRecordsFull + " 지금은 새 agent도 만들 수 없습니다. " + agentProblem(limit), false
		}
		return keyRecordsFull, true
	case a != nil && mode == "register" && activeKeys(a) >= 3:
		return "활성 키는 agent당 3개입니다. 쓰지 않는 키를 철회하거나 회전 방식으로 연결하세요.", false
	}
	return safeProblem(fault("capacity")), false
}
func agentProblem(limit string) string {
	if m := agentLimits[limit]; m != "" {
		return m
	}
	return safeProblem(fault("capacity"))
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
	mux.HandleFunc("POST /home/agents", s.memberAction("create"))
	mux.HandleFunc("POST /home/connect", s.memberAction("connect"))
	mux.HandleFunc("GET /home/connections/{id}", s.connectionPage)
	for _, path := range []string{"confirm", "cancel", "key-revoke", "agent-revoke", "invites", "invite-decision", "unpair"} {
		mux.HandleFunc("POST /home/"+path, s.memberAction(path))
	}
	for _, path := range []string{"info", "prepare", "complete"} {
		mux.HandleFunc("POST /v1/connect/"+path, s.connectAPI(path))
	}
}
func (s *Service) memberAction(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 8192)
		parseErr := r.ParseForm()
		value, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
			clean := path == "cancel" || path == "key-revoke" || path == "agent-revoke" || path == "unpair" || (path == "invite-decision" && r.FormValue("decision") == "deny")
			rates := memberRate
			if clean {
				rates = cleanupRate
			}
			id, session, retry, e := s.memberHit(st, r, now, rates)
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
			if parseErr != nil {
				return fail(fault("invalid_schema"))
			}
			if path == "create" || path == "connect" || path == "confirm" || path == "key-revoke" || path == "agent-revoke" {
				if now.Sub(session.Verified) >= reauthWindow {
					return fail(fault("reauth_required"))
				}
			}
			agent := r.FormValue("agent")
			switch path {
			case "create":
				if limit := st.agentLimit(owner); limit != "" {
					return refusal{status: 409, problem: agentProblem(limit)}, nil
				}
				agent = fmt.Sprintf("agent_%x", sha256.Sum256([]byte(randomToken())))[:28] // Avoid email-derived IDs and cross-member name collisions.
				st.Agents[agent] = &Agent{Owner: owner, Keys: map[string]*Key{}}
				return "/home", nil
			case "connect":
				t, e := st.beginConnection(owner, agent, r.FormValue("client"), r.FormValue("mode"), now)
				if e != nil && e.Error() == "capacity" {
					problem, create := connectLimit(st, owner, agent, r.FormValue("mode"))
					return refusal{status: 409, problem: problem, create: create}, nil
				}
				if e != nil {
					return fail(e)
				}
				return map[string]any{"Title": "agent 연결 대기", "Token": t, "ConnectionID": hashToken(t), "Agent": agent, "Owner": owner, "Client": supportedClient, "Mode": r.FormValue("mode"), "State": "waiting", "Exp": clock(now.Add(connectionTTL))}, nil
			case "confirm", "cancel":
				c := st.Connections[r.FormValue("connection")]
				if c == nil || c.Owner != owner || !now.Before(c.Exp) || c.State == "consumed" || c.State == "cancelled" {
					return fail(fault("invalid_auth"))
				}
				if path == "cancel" {
					c.State = "cancelled"
				} else {
					if c.State != "prepared" {
						return fail(fault("sender_not_allowed"))
					}
					c.State = "approved"
				}
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
				refused(w, "내 agent", refusal{503, "일시적으로 처리할 수 없습니다. 잠시 뒤 다시 시도하세요.", false, false})
			}
			return
		}
		switch v := value.(type) {
		case refusal:
			refused(w, "내 agent", v)
		case string:
			http.Redirect(w, r, v, 303)
		default:
			render(w, 200, "connection", value.(map[string]any))
		}
	}
}
func (s *Service) connectionPage(w http.ResponseWriter, r *http.Request) {
	value, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
		id, _, retry, e := s.memberHit(st, r, now, memberRate)
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
		return map[string]any{"Title": "agent 연결 확인", "ConnectionID": r.PathValue("id"), "Agent": c.Agent, "Client": c.Client, "Mode": c.Mode, "State": c.State, "Fingerprint": fingerprint(public), "Exp": clock(c.Exp)}, nil
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
		refused(w, "연결 확인", v)
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
			if ok, retry := st.hit(now, rates...); !ok {
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
		if retry, ok := value.(time.Time); ok {
			rateLimited(w, retry)
			return
		}
		if v, ok := value.(map[string]string); ok && v["error"] != "" {
			respond(w, nil, fault(v["error"]))
			return
		}
		respond(w, value, err)
	}
}
