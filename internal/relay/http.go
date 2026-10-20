// HTTP API separates owner credentials from agent credentials and serves escaped approval forms.
package relay

import (
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"strings"
	"time"
)

type command struct {
	Agent    string `json:"agent,omitempty"`
	Kid      string `json:"kid,omitempty"`
	Public   string `json:"public,omitempty"`
	Proof    string `json:"proof,omitempty"`
	Target   string `json:"target,omitempty"`
	Decision string `json:"decision,omitempty"`
	ID       string `json:"id,omitempty"`
	Token    string `json:"token,omitempty"`
	Claim    string `json:"claim,omitempty"`
}

func bearer(r *http.Request) string {
	return strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
}
func ownerToken(r *http.Request) string {
	if _, password, ok := r.BasicAuth(); ok {
		return password
	}
	return bearer(r)
}
func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func statusFor(err error) int {
	switch err.Error() {
	case "invalid_auth":
		return 401
	case "sender_not_allowed", "human_invite_required", "disclosure_denied", "invalid_csrf":
		return 403
	case "idempotency_conflict", "id_collision", "duplicate_gate", "duplicate_result", "already_claimed", "invalid_lease", "invalid_gate", "key_exists":
		return 409
	case "rate_limited":
		return 429
	case "capacity":
		return 409
	case "unavailable", "abnormal_clock":
		return 503
	default:
		return 422
	}
}
func (s *Service) Handler() http.Handler {
	mux := http.NewServeMux()
	s.remoteRoutes(mux)
	mux.HandleFunc("GET /v1/registry", func(w http.ResponseWriter, r *http.Request) {
		if !RegistryValid() {
			respond(w, nil, fault("unavailable"))
			return
		}
		writeJSON(w, 200, map[string]any{"sha256": RegistrySHA256, "manifest": string(registry)})
	})
	mux.HandleFunc("POST /v1/owners", func(w http.ResponseWriter, r *http.Request) {
		// Local synthetic signup is not identity verification; public deployments leave it off so it cannot mint members or owners.
		if !s.SyntheticSignup {
			respond(w, nil, fault("sender_not_allowed"))
			return
		}
		result, err := s.transaction(r.Context(), func(st *State, _ time.Time) (any, error) {
			token := randomToken()
			id := "owner_" + randomToken()
			st.Owners[id] = &Owner{hashToken(token), true}
			return map[string]string{"owner": id, "credential": token}, nil
		})
		respond(w, result, err)
	})
	for _, path := range []string{"agents", "keys", "key-revoke", "invites", "invite-decision", "unpair", "owner-revoke", "pull", "persist", "ack", "claim", "gate-consume", "authorize"} {
		mux.HandleFunc("POST /v1/"+path, s.operation(path))
	}
	mux.HandleFunc("POST /v1/send", func(w http.ResponseWriter, r *http.Request) {
		raw, err := io.ReadAll(http.MaxBytesReader(w, r.Body, rawEnvelopeLimit))
		if err != nil {
			respond(w, nil, fault("invalid_json"))
			return
		}
		e, err := Parse(raw)
		if err != nil {
			respond(w, nil, err)
			return
		}
		result, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
			agent, err := st.principal(bearer(r), "agent")
			if err != nil {
				return nil, err
			}
			return st.ingest(agent, e, r.Header.Get("X-Execution-Claim"), now)
		})
		respond(w, result, err)
	})
	mux.HandleFunc("GET /v1/gates/{id}", func(w http.ResponseWriter, r *http.Request) {
		result, err := s.transaction(r.Context(), func(st *State, _ time.Time) (any, error) {
			agent, err := st.principal(bearer(r), "agent")
			if err != nil {
				return nil, err
			}
			g := st.Gates[r.PathValue("id")]
			if g == nil || g.To != agent {
				return nil, fault("sender_not_allowed")
			}
			v := map[string]any{"state": lookupVisible(g), "consumed": g.Consumed, "parent": g.Parent}
			if g.Lookup != nil && g.Lookup.State != "" {
				v["lookup"] = g.Lookup.State
			}
			return v, nil
		})
		respond(w, result, err)
	})
	mux.HandleFunc("GET /v1/contacts", func(w http.ResponseWriter, r *http.Request) {
		result, err := s.transaction(r.Context(), func(st *State, _ time.Time) (any, error) {
			agent, err := st.principal(bearer(r), "agent")
			if err != nil {
				return nil, err
			}
			contacts := []string{}
			for target := range st.Agents {
				if _, _, ok := st.authority(agent, target); ok {
					contacts = append(contacts, target)
				}
			}
			return contacts, nil
		})
		respond(w, result, err)
	})
	// One contact's public label for a work session's toast and send dialog (O-03): HP-05 agent name and client plus the owner's member name; no email, key or credential.
	mux.HandleFunc("GET /v1/contacts/{agent}", func(w http.ResponseWriter, r *http.Request) {
		result, err := s.transaction(r.Context(), func(st *State, _ time.Time) (any, error) {
			agent, err := st.principal(bearer(r), "agent")
			if err != nil {
				return nil, err
			}
			target := r.PathValue("agent")
			if _, _, ok := st.authority(agent, target); !ok {
				return nil, fault("sender_not_allowed")
			}
			a := st.Agents[target]
			v := humanAgentView(target, a)
			friend := ""
			if _, m := st.humanMember(a.Owner); m != nil {
				friend = m.Name
			}
			return map[string]string{"agent": target, "name": v.Name, "client": v.Client, "friend": friend}, nil
		})
		respond(w, result, err)
	})
	mux.HandleFunc("GET /v1/keys/{agent}/{kid}", func(w http.ResponseWriter, r *http.Request) {
		result, err := s.transaction(r.Context(), func(st *State, _ time.Time) (any, error) {
			agent, err := st.principal(bearer(r), "agent")
			if err != nil {
				return nil, err
			}
			target := r.PathValue("agent")
			_, _, allowed := st.authority(agent, target)
			if agent != target && !allowed {
				return nil, fault("sender_not_allowed")
			}
			a := st.Agents[target]
			if a == nil {
				return nil, fault("invalid_auth")
			}
			k := a.Keys[r.PathValue("kid")]
			if k == nil || k.Revoked {
				return nil, fault("invalid_auth")
			}
			return map[string]string{"public": base64.RawURLEncoding.EncodeToString(k.Public)}, nil
		})
		respond(w, result, err)
	})
	mux.HandleFunc("GET /v1/receipts/{id}", func(w http.ResponseWriter, r *http.Request) {
		result, err := s.transaction(r.Context(), func(st *State, _ time.Time) (any, error) {
			agent, err := st.principal(bearer(r), "agent")
			if err != nil {
				return nil, err
			}
			m := st.Messages[r.PathValue("id")]
			if m == nil || (m.Receipt.From != agent && m.Receipt.To != agent) || !st.current(m) {
				return nil, fault("sender_not_allowed")
			}
			v := map[string]any{"receipt": m.Receipt, "completion": m.Completion, "reply_to": m.Parent, "reply_id": m.ReplyID}
			if !m.SessionSeen.IsZero() {
				v["session_received_at"] = m.SessionSeen
			}
			return v, nil
		})
		respond(w, result, err)
	})
	mux.HandleFunc("GET /owner", s.ownerPage)
	mux.HandleFunc("GET /owner/gates/{id}", s.gatePage(basicOwner))
	mux.HandleFunc("POST /owner/gates/{id}", s.gateDecision(basicOwner))
	mux.HandleFunc("POST /owner/gates/{id}/deny", s.gateDecision(basicOwner))
	s.memberRoutes(mux)
	s.humanRoutes(mux)
	s.agentRoutes(mux)
	s.textRoutes(mux)
	mux.HandleFunc("GET /home/receipts", s.memberReceipt)
	// Stdlib Sec-Fetch-Site/Origin check rejects cross-origin browser writes; non-browser agent calls carry neither header.
	protected := crossOrigin.Handler(mux)
	ratedTest := s.testHandler(mux)
	bounded := s.boundedHTTP(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'unsafe-inline'; form-action 'self' https://accounts.google.com; frame-ancestors 'none'; base-uri 'none'")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		if strings.HasPrefix(r.URL.Path, "/v1/test/") {
			ratedTest.ServeHTTP(w, r)
			return
		}
		protected.ServeHTTP(w, r)
	}))
	root := http.NewServeMux()
	root.HandleFunc("GET /v1/text/events", s.textEvents)
	root.Handle("/", bounded)
	return root
}
func respond(w http.ResponseWriter, result any, err error) {
	if err != nil {
		writeJSON(w, statusFor(err), map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, 200, result)
}
func (s *Service) operation(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var c command
		raw, err := io.ReadAll(http.MaxBytesReader(w, r.Body, rawEnvelopeLimit))
		if err != nil {
			respond(w, nil, fault("invalid_json"))
			return
		}
		if len(raw) == 0 {
			raw = []byte("{}")
		}
		if err := Strict(raw, &c); err != nil {
			respond(w, nil, err)
			return
		}
		result, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) { return st.operate(path, bearer(r), c, now) })
		respond(w, result, err)
	}
}
func (st *State) operate(path, token string, c command, now time.Time) (any, error) {
	testOnly := strings.HasPrefix(path, "test-")
	if testOnly {
		path = strings.TrimPrefix(path, "test-")
	}
	kind := "agent"
	switch path {
	case "agents", "keys", "key-revoke", "invite-decision", "unpair", "owner-revoke":
		kind = "owner"
	}
	principal, err := st.principal(token, kind)
	if err != nil && path == "invites" {
		principal, err = st.principal(token, "owner")
		kind = "owner"
	}
	if err != nil {
		return nil, err
	}
	if testOnly {
		if !st.TestAgents[principal] {
			return nil, fault("sender_not_allowed")
		}
		if path != "pull" {
			m := st.Messages[c.ID]
			if m == nil || m.Receipt.Intent != "relay.test.message" {
				return nil, fault("sender_not_allowed")
			}
		}
	}
	if m := st.Messages[c.ID]; m != nil && m.Receipt.Intent == publicTextIntent {
		return nil, fault("sender_not_allowed")
	}
	return st.operateAs(path, principal, kind, c, now, testOnly)
}

func (st *State) operateAs(path, principal, kind string, c command, now time.Time, testOnly bool) (any, error) {
	switch path {
	case "agents", "keys":
		if !agentPattern.MatchString(c.Agent) || !kidPattern.MatchString(c.Kid) {
			return nil, fault("invalid_schema")
		}
		a := st.Agents[c.Agent]
		if path == "agents" && a != nil {
			return nil, fault("key_exists")
		}
		if path == "keys" && (a == nil || a.Owner != principal) {
			return nil, fault("sender_not_allowed")
		}
		public, e := base64.RawURLEncoding.DecodeString(c.Public)
		proof, e2 := base64.RawURLEncoding.DecodeString(c.Proof)
		if e != nil || e2 != nil || len(public) != 32 || !ed25519.Verify(public, popBytes(principal, c.Agent, c.Kid, c.Public), proof) {
			return nil, fault("invalid_signature")
		}
		credential := ""
		if a == nil {
			if !st.agentCapacity(principal) {
				return nil, fault("capacity")
			}
			credential = randomToken()
			a = &Agent{Owner: principal, Credential: hashToken(credential), Keys: map[string]*Key{}}
			st.Agents[c.Agent] = a
		}
		if a.Keys[c.Kid] != nil {
			return nil, fault("key_exists")
		}
		if len(a.Keys) >= agentKeyRecords {
			return nil, fault("capacity")
		}
		if path == "keys" {
			revokeKeys(a, now)
		}
		a.Keys[c.Kid] = &Key{Public: public, Changed: now}
		st.sweep(now)
		return map[string]string{"agent": c.Agent, "credential": credential}, nil
	case "key-revoke":
		a := st.Agents[c.Agent]
		if a == nil || a.Owner != principal || a.Keys[c.Kid] == nil {
			return nil, fault("sender_not_allowed")
		}
		if k := a.Keys[c.Kid]; !k.Revoked {
			k.Revoked = true
			k.Changed = now
		}
		st.sweep(now)
		return map[string]string{"state": "revoked"}, nil
	case "owner-revoke":
		st.Owners[principal].Active = false
		st.sweep(now)
		return map[string]string{"state": "revoked"}, nil
	case "invites":
		if st.HumanModel == 1 && !st.SyntheticSignup {
			return nil, fault("human_invite_required")
		}
		a, b := st.Agents[c.Agent], st.Agents[c.Target]
		if a == nil || b == nil || a.Revoked || b.Revoked || c.Agent == c.Target || (kind == "agent" && principal != c.Agent) || (kind == "owner" && a.Owner != principal) || !st.Owners[b.Owner].Active {
			return nil, fault("sender_not_allowed")
		}
		if _, m := st.humanMember(a.Owner); m != nil {
			return nil, fault("human_invite_required")
		}
		if _, m := st.humanMember(b.Owner); m != nil {
			return nil, fault("human_invite_required")
		}
		id := pairID(c.Agent, c.Target)
		p := st.Pairs[id]
		if p != nil && (p.State == "active" || p.State == "pending") {
			return p, nil
		}
		if !st.pairCapacity(a.Owner, b.Owner, false) {
			return nil, fault("capacity")
		}
		generation := int64(1)
		if p != nil {
			generation = p.Generation + 1
		}
		p = &Pair{A: c.Agent, B: c.Target, Inviter: c.Agent, Recipient: c.Target, State: "pending", Generation: generation, Exp: now.Add(24 * time.Hour)}
		st.Pairs[id] = p
		return p, nil
	case "invite-decision":
		p := st.Pairs[pairID(c.Agent, c.Target)]
		if p == nil || st.Agents[p.A] == nil || st.Agents[p.B] == nil || st.Agents[p.Recipient].Owner != principal || (c.Decision != "accept" && c.Decision != "deny") {
			return nil, fault("sender_not_allowed")
		}
		if c.Decision == "accept" && p.State != "active" {
			if _, m := st.humanMember(st.Agents[p.A].Owner); m != nil {
				return nil, fault("human_invite_required")
			}
			if _, m := st.humanMember(st.Agents[p.B].Owner); m != nil {
				return nil, fault("human_invite_required")
			}
		}
		if p.State == "active" && c.Decision == "accept" {
			return p, nil
		}
		if p.State != "pending" {
			return nil, fault("sender_not_allowed")
		}
		if c.Decision == "accept" {
			if !st.pairCapacity(st.Agents[p.A].Owner, st.Agents[p.B].Owner, true) {
				return nil, fault("capacity")
			}
			p.State = "active"
		} else {
			p.State = "denied"
		}
		return p, nil
	case "unpair":
		p := st.Pairs[pairID(c.Agent, c.Target)]
		if p == nil || st.Agents[p.A] == nil || st.Agents[p.B] == nil || (st.Agents[p.A].Owner != principal && st.Agents[p.B].Owner != principal) {
			return nil, fault("sender_not_allowed")
		}
		p.State = "revoked"
		st.sweep(now)
		return p, nil
	case "pull":
		return st.leaseMessageFor(principal, now, testOnly)
	case "persist", "ack":
		m, err := st.leased(principal, c.ID, c.Token, now)
		if err != nil {
			return nil, err
		}
		if path == "persist" {
			m.Persisted = true
			m.Inbox = append(json.RawMessage(nil), m.Envelope...)
			return map[string]string{"state": "persisted"}, nil
		}
		if !m.Persisted {
			return nil, fault("invalid_lease")
		}
		m.Receipt.State = "delivered"
		m.LeaseToken = ""
		if m.Receipt.Intent == "relay.result" || m.Receipt.Intent == publicTextIntent {
			if m.Receipt.Intent == publicTextIntent {
				m.Completion = "received"
				if parent := st.Messages[m.Parent]; parent != nil {
					parent.Completion = "reply_received"
				}
			}
			m.Envelope = nil
			m.Inbox = nil
		}
		return m.Receipt, nil
	case "claim":
		m := st.Messages[c.ID]
		if m == nil || m.Receipt.To != principal || m.Receipt.State != "delivered" || !m.Persisted || len(m.Inbox) == 0 || !st.current(m) || !now.Before(m.Receipt.Exp) || m.Receipt.Intent == "relay.result" || m.Receipt.Intent == publicTextIntent || m.Deliver != "agent" {
			return nil, fault("sender_not_allowed")
		}
		if m.Claimed {
			return nil, fault("already_claimed")
		}
		if !st.claimCapacity() {
			return nil, fault("capacity")
		}
		m.Claimed = true
		m.ClaimToken = randomToken()
		if m.Receipt.Intent == "relay.test.message" {
			m.Envelope = nil
			m.Inbox = nil
		}
		return map[string]string{"claim": m.ClaimToken, "policy": policy(m.Receipt.Intent)}, nil
	case "gate-consume":
		g := st.Gates[c.ID]
		if g == nil || g.To != principal || lookupVisible(g) != "approved" || g.Consumed || !now.Before(g.Exp) {
			return nil, fault("invalid_gate")
		}
		m, err := st.parentFor(principal, g.Parent, now)
		if err != nil {
			return nil, err
		}
		if m.ClaimToken != c.Claim || c.Claim == "" || m.Receipt.Digest != g.Digest || m.Generation != g.Generation {
			return nil, fault("invalid_gate")
		}
		g.Consumed = true
		// The disclosure is the only thing approval enables: the agent must echo it unchanged in its done result.
		v := map[string]any{"gate_passed": true, "executable": false, "disclosure": false}
		if g.Lookup != nil {
			v["lookup"] = g.Lookup.State
			if g.Lookup.State == "done" && len(g.Lookup.Disclosure) > 0 {
				v["disclosure"] = g.Lookup.Disclosure
			}
		}
		return v, nil
	case "authorize":
		m, err := st.parentFor(principal, c.ID, now)
		if err != nil {
			return nil, err
		}
		if c.Claim == "" || m.ClaimToken != c.Claim {
			return nil, fault("sender_not_allowed")
		}
		return map[string]any{"executable": false, "disclosure": false, "policy": policy(m.Receipt.Intent)}, nil
	}
	return nil, fault("invalid_schema")
}
func csrf(token, id string) string {
	mac := hmac.New(sha256.New, []byte(token))
	_, _ = mac.Write([]byte("gate:" + id))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

// The gate page shares the member visual tokens; the verified summary leads and the raw body stays in technical detail.

// ownerAuth returns the CSRF key, the owner resolver, and the page to return to.
type ownerAuth func(r *http.Request) (string, func(*State, time.Time) (string, error), string)

func basicOwner(r *http.Request) (string, func(*State, time.Time) (string, error), string) {
	token := ownerToken(r)
	return token, func(st *State, _ time.Time) (string, error) { return st.principal(token, "owner") }, "/owner"
}
func sessionOwner(r *http.Request) (string, func(*State, time.Time) (string, error), string) {
	token := readCookie(r, sessionCookie)
	return token, func(st *State, now time.Time) (string, error) {
		id, _, err := st.session(token, now)
		if err != nil {
			return "", err
		}
		return st.Members[id].Owner, nil
	}, "/home"
}

func (s *Service) gateOwner(w http.ResponseWriter, st *State, r *http.Request, now time.Time, auth ownerAuth, clean bool) (string, error) {
	_, resolve, back := auth(r)
	if back == "/home" {
		rates := memberRate
		if clean {
			rates = cleanupRate
		}
		id, _, retry, err := s.memberHit(w, st, r, now, rates)
		if !retry.IsZero() {
			return "", &budgetError{retry}
		}
		if err != nil {
			return "", err
		}
		return st.Members[id].Owner, nil
	}
	owner, err := resolve(st, now)
	rates := anonymousRate(s.clientIP(r))
	if err == nil {
		rates = memberRate(owner)
		if clean {
			rates = cleanupRate(owner)
		}
	}
	if ok, retry := s.requestHit(st, r, now, rates...); !ok {
		return "", &budgetError{retry}
	}
	return owner, err
}

type budgetError struct{ retry time.Time }

func (e *budgetError) Error() string { return "rate_limited" }

func verifiedGate(st *State, g *Gate, now time.Time) (*Envelope, bool) {
	m, h := st.Messages[g.Parent], st.Messages[g.ID]
	if m == nil || h == nil || len(m.Envelope) == 0 {
		return nil, false
	}
	e, err := Parse(m.Envelope)
	if err != nil || gateKinds[e.Intent] == "" {
		return nil, false
	}
	a := st.Agents[e.From]
	if a == nil || a.Keys[e.Sig.Kid] == nil || !e.Verify(a.Keys[e.Sig.Kid].Public) || e.Digest() != m.Receipt.Digest || e.ID != m.Receipt.ID || e.From != g.From || e.To != g.To || m.Generation != g.Generation || g.Digest != m.Receipt.Digest {
		return nil, false
	}
	return e, g.State == "pending" && st.current(m) && st.current(h) && now.Before(g.Exp) && now.Before(m.Receipt.Exp)
}
func (s *Service) gateView(w http.ResponseWriter, r *http.Request, auth ownerAuth) (any, error) {
	token, _, back := auth(r)
	return s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
		owner, err := s.gateOwner(w, st, r, now, auth, false)
		if err != nil {
			return nil, err
		}
		g := st.Gates[r.PathValue("id")]
		if g == nil || g.Owner != owner {
			return nil, fault("sender_not_allowed")
		}
		state, e, active := gateState(st, g, now)
		list := back
		if back == "/home" {
			list = "/home/legacy#gates"
		}
		view := map[string]any{"State": state, "ID": g.ID, "From": g.From, "To": g.To, "Parent": g.Parent, "Exp": clock(g.Exp), "Policy": g.Policy, "CSRF": csrf(token, g.ID), "Active": active, "Back": list, "Deny": back + "/gates/" + g.ID + "/deny", "FromAgent": st.gateAgent(owner, g.From), "ToAgent": st.gateAgent(owner, g.To)}
		if e != nil {
			body, _ := json.MarshalIndent(e.Body, "", "  ")
			summary, text := gateSummary(e)
			view["Summary"], view["Text"] = summary, text
			// A request whose content cannot be confirmed keeps no approve/deny form (3단위 '확인할 수 없습니다').
			view["Active"] = active && summary != nil
			view["Body"], view["Intent"], view["Scope"] = string(body), e.Intent, gateScopes[e.Intent]
			if e.Intent == "schedule.query" {
				// Only a connected member can answer; the fixed SAR-SCHEDULE-001 explanation replaces the U5 text then (SC-07).
				calendar := st.calendarState(owner)
				view["Calendar"] = calendar
				if calendar == "connected" {
					view["Scope"] = calendarGateScope
				} else {
					view["Active"] = false
				}
			}
			if e.Render != nil {
				view["Hint"] = e.Render["hint"]
			}
		}
		return view, nil
	})
}
func ownerError(w http.ResponseWriter, err error) {
	if err.Error() == "invalid_auth" {
		w.Header().Set("WWW-Authenticate", `Basic realm="KnowsLink owner", charset="UTF-8"`)
	}
	http.Error(w, err.Error(), statusFor(err))
}

// gateError sends members back to email login instead of the owner Basic prompt.
// Member refusals return to the approval list; a stale or repeated submission never reads as a fresh decision.
func gateError(w http.ResponseWriter, r *http.Request, back string, err error) {
	list := "/home"
	if back == "/home" {
		list = "/home/legacy#gates"
	}
	if e, ok := err.(*budgetError); ok {
		refusedTo(w, "요청 승인", refusal{status: 429, problem: limited(e.retry)}, list)
		return
	}
	if back == "/home" && err.Error() == "invalid_auth" {
		expired(w, r)
		return
	}
	if back == "/home" {
		problem := "이 승인 요청을 볼 수 없습니다. 승인 목록에서 내 요청을 다시 선택하세요."
		if r.Method == "POST" {
			problem = "이번 제출은 적용되지 않았습니다. 요청이 이미 처리됐거나 기한·권한·원문이 바뀌었을 수 있습니다. 승인 목록에서 현재 상태를 확인하세요."
		}
		refusedTo(w, "요청 승인", refusal{status: statusFor(err), problem: problem}, list)
		return
	}
	ownerError(w, err)
}
func (s *Service) gatePage(auth ownerAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		view, err := s.gateView(w, r, auth)
		if err != nil {
			_, _, back := auth(r)
			gateError(w, r, back, err)
			return
		}
		render(w, 200, "human-gate", view.(map[string]any))
	}
}
func (s *Service) gateDecision(auth ownerAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, _, back := auth(r)
		r.Body = http.MaxBytesReader(w, r.Body, 8*1024)
		parseErr := r.ParseForm()
		value, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
			owner, err := s.gateOwner(w, st, r, now, auth, r.PostForm.Get("decision") == "deny")
			if err != nil {
				return nil, err
			}
			if parseErr != nil || (strings.HasSuffix(r.URL.Path, "/deny") && r.PostForm.Get("decision") != "deny") {
				return nil, fault("invalid_schema")
			}
			if !hmac.Equal([]byte(r.PostForm.Get("csrf")), []byte(csrf(token, r.PathValue("id")))) {
				return nil, fault("invalid_csrf")
			}

			g := st.Gates[r.PathValue("id")]
			if g == nil || g.Owner != owner || g.State != "pending" || !now.Before(g.Exp) {
				return nil, fault("invalid_gate")
			}
			m := st.Messages[g.Parent]
			h := st.Messages[g.ID]
			if _, active := verifiedGate(st, g, now); !active || m == nil || h == nil {
				return nil, fault("invalid_gate")
			}
			decision := r.PostForm.Get("decision")
			if decision != "approve" && decision != "deny" {
				return nil, fault("invalid_schema")
			}
			var job *lookupJob
			if decision == "approve" {
				// An approval of a schedule.query is one lookup; a member without a connected calendar cannot approve (SC-02/07).
				e, _ := verifiedGate(st, g, now)
				if job, err = st.startLookup(g, e, now); err != nil {
					return nil, err
				}
			}
			h.Receipt.State = "delivered"
			h.Persisted = true
			h.Envelope = nil
			h.Inbox = nil
			g.State = "denied"
			if decision == "approve" {
				g.State = "approved"
			}
			return job, nil
		})
		if err != nil {
			gateError(w, r, back, err)
			return
		}
		if job, ok := value.(*lookupJob); ok && job != nil {
			s.runLookup(r.Context(), job)
		}
		// Both decisions land on the canonical gate page; /deny is POST-only and its GET would end in 405.
		http.Redirect(w, r, back+"/gates/"+r.PathValue("id"), http.StatusSeeOther)
	}
}

var dashboard = template.Must(template.New("owner").Parse(`<!doctype html><html lang="ko"><meta charset="utf-8"><title>KnowsLink Owner</title><h1>KnowsLink Owner 작업 화면</h1><p>가입·키·초대 관리는 인증된 /v1 API를 사용합니다. 원문은 요청 만료 또는 완료 후 지워집니다.</p><h2>등록 에이전트</h2><ul>{{range .Agents}}<li>{{.}}</li>{{end}}</ul><h2>승인 요청</h2><ul>{{range .Gates}}<li><a href="/owner/gates/{{.ID}}">{{.Parent}}</a> 상태: {{.State}}</li>{{end}}</ul></html>`))

func (s *Service) ownerPage(w http.ResponseWriter, r *http.Request) {
	view, err := s.transaction(r.Context(), func(st *State, _ time.Time) (any, error) {
		owner, err := st.principal(ownerToken(r), "owner")
		if err != nil {
			return nil, err
		}
		agents := []string{}
		gates := []*Gate{}
		for id, a := range st.Agents {
			if a.Owner == owner {
				agents = append(agents, id)
			}
		}
		for _, g := range st.Gates {
			if g.Owner == owner {
				gates = append(gates, g)
			}
		}
		return map[string]any{"Agents": agents, "Gates": gates}, nil
	})
	if err != nil {
		ownerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = dashboard.Execute(w, view)
}

func decodeCommand(w http.ResponseWriter, r *http.Request, target any) error {
	raw, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 8192))
	if err != nil {
		return fault("invalid_json")
	}
	return Strict(raw, target)
}
