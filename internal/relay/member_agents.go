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
}

type memberAgent struct {
	ID, Status string
	Keys       []memberKey
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
		v := memberAgent{ID: id, Status: "미연결"}
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
		if err := r.ParseForm(); err != nil {
			render(w, 422, "head", map[string]any{"Title": "내 agent", "Problem": "입력 크기나 형식이 올바르지 않습니다."})
			return
		}
		token := readCookie(r, sessionCookie)
		value, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
			id, session, e := st.session(token, now)
			if e != nil {
				return nil, e
			}
			owner := st.Members[id].Owner
			clean := path == "cancel" || path == "key-revoke" || path == "agent-revoke" || path == "unpair" || (path == "invite-decision" && r.FormValue("decision") == "deny")
			rates := memberRate(id)
			if clean {
				rates = cleanupRate(id)
			}
			if ok, retry := st.hit(now, rates...); !ok {
				return refusal{429, limited(retry)}, nil
			}
			// Return refusals as values: rejected work must commit its request budget.
			fail := func(e error) (any, error) { return refusal{statusFor(e), safeProblem(e)}, nil }
			if path == "create" || path == "connect" || path == "confirm" || path == "key-revoke" || path == "agent-revoke" {
				if now.Sub(session.Verified) >= reauthWindow {
					return fail(fault("reauth_required"))
				}
			}
			agent := r.FormValue("agent")
			switch path {
			case "create":
				if !st.agentCapacity(owner) {
					return fail(fault("capacity"))
				}
				agent = fmt.Sprintf("agent_%x", sha256.Sum256([]byte(randomToken())))[:28] // Avoid email-derived IDs and cross-member name collisions.
				st.Agents[agent] = &Agent{Owner: owner, Keys: map[string]*Key{}}
				return "/home", nil
			case "connect":
				t, e := st.beginConnection(owner, agent, r.FormValue("client"), r.FormValue("mode"), now)
				if e != nil {
					return fail(e)
				}
				return map[string]any{"Title": "agent 연결 대기", "Token": t, "ConnectionID": hashToken(t), "Agent": agent, "Owner": owner, "Client": supportedClient, "Mode": r.FormValue("mode"), "Exp": clock(now.Add(connectionTTL))}, nil
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
				a.Revoked = true
				a.Credential = ""
				for _, k := range a.Keys {
					k.Revoked = true
					k.Changed = now
				}
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
				_, e := st.operateAs(path, owner, "owner", c, now, false)
				if e != nil {
					return fail(e)
				}
			}
			return "/home", nil
		})
		if err != nil {
			if err.Error() == "invalid_auth" {
				expired(w, r)
			} else {
				render(w, 503, "head", map[string]any{"Title": "내 agent", "Problem": "일시적으로 처리할 수 없습니다."})
			}
			return
		}
		switch v := value.(type) {
		case refusal:
			render(w, v.status, "head", map[string]any{"Title": "내 agent", "Problem": v.problem})
		case string:
			http.Redirect(w, r, v, 303)
		default:
			render(w, 200, "connection", value.(map[string]any))
		}
	}
}
func (s *Service) connectionPage(w http.ResponseWriter, r *http.Request) {
	value, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
		id, _, e := st.session(readCookie(r, sessionCookie), now)
		if e != nil {
			return nil, e
		}
		if ok, retry := st.hit(now, memberRate(id)...); !ok {
			return refusal{429, limited(retry)}, nil
		}
		c := st.Connections[r.PathValue("id")]
		if c == nil || c.Owner != st.Members[id].Owner {
			return refusal{403, "연결을 조회할 수 없습니다."}, nil
		}
		public, _ := base64.RawURLEncoding.DecodeString(c.Public)
		return map[string]any{"Title": "agent 연결 확인", "ConnectionID": r.PathValue("id"), "Agent": c.Agent, "Client": c.Client, "Mode": c.Mode, "State": c.State, "Fingerprint": fingerprint(public), "Exp": clock(c.Exp)}, nil
	})
	if err != nil {
		if err.Error() == "invalid_auth" {
			expired(w, r)
		} else {
			render(w, 503, "head", map[string]any{"Title": "연결 확인", "Problem": "일시적으로 처리할 수 없습니다. 잠시 뒤 다시 시도하세요."})
		}
		return
	}
	if v, ok := value.(refusal); ok {
		render(w, v.status, "head", map[string]any{"Title": "연결 확인", "Problem": v.problem})
		return
	}
	render(w, 200, "connection", value.(map[string]any))
}
func (s *Service) connectAPI(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var c struct{ Token, Client, Kid, Public, Proof string }
		if err := decodeCommand(w, r, &c); err != nil {
			respond(w, nil, err)
			return
		}
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
				return map[string]string{"error": "rate_limited", "retry_at": retry.Format(time.RFC3339)}, nil
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
		if v, ok := value.(map[string]string); ok && v["error"] != "" {
			if v["retry_at"] != "" {
				w.Header().Set("Retry-After", fmt.Sprint(60))
			}
			respond(w, nil, fault(v["error"]))
			return
		}
		respond(w, value, err)
	}
}
