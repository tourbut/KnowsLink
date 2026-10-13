// Session-only person, friendship and conversation pages share existing origin, rate and authorization boundaries.
package relay

import (
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"
)

type humanAgent struct{ ID, Name, Client, Status, Short string }
type humanFriend struct {
	ID, Name, State, Exp string
	Generation           int64
	Incoming             bool
}
type humanConversation struct{ ID, Name, FromName, ToName, State, URL string }
type humanMessage struct{ ID, Author, Agent, Short, Text, When, State, Authorship string }

func humanAgentView(id string, a *Agent) humanAgent {
	client := a.Client
	if client == "" {
		client = "Node 22 로컬 클라이언트"
	}
	status := "연결 필요"
	if activeKeys(a) > 0 {
		status = "등록됨"
	}
	if a.Revoked {
		status = "철회"
	}
	return humanAgent{ID: id, Name: agentName(a), Client: client, Status: status, Short: shortID(id)}
}

func (st *State) humanAgents(owner string) []humanAgent {
	agents := []humanAgent{}
	for id, a := range st.Agents {
		if a.Owner == owner && !a.Revoked {
			agents = append(agents, humanAgentView(id, a))
		}
	}
	sort.Slice(agents, func(i, j int) bool { return agents[i].ID < agents[j].ID })
	return agents
}

func conversationURL(from, to string) string {
	return "/home/conversation?from=" + url.QueryEscape(from) + "&to=" + url.QueryEscape(to)
}

func (s *Service) humanRoutes(mux *http.ServeMux) {
	for _, path := range []string{"name", "friends", "agents", "account", "new", "conversation"} {
		mux.HandleFunc("GET /home/"+path, s.humanPage(path))
	}
	mux.HandleFunc("GET /home/friends/{id}", s.humanPage("friend"))
	mux.HandleFunc("GET /home/legacy", s.legacyHomePage)
	for _, path := range []string{"profile", "agent-name", "friends/request", "friends/decision", "friend-code", "self-connect", "conversation/send"} {
		mux.HandleFunc("POST /home/"+path, s.humanAction(path))
	}
}

func (st *State) humanView(id string, session *Session, path string, r *http.Request, now time.Time) (map[string]any, error) {
	m := st.Members[id]
	v := map[string]any{"Title": "대화", "Section": "conversations", "Name": m.Name, "Masked": maskEmail(m.Email), "Code": m.FriendCode, "SelfConnected": m.SelfConnect, "Recent": now.Sub(session.Verified) < reauthWindow, "Attempt": randomToken(), "From": humanAgent{}, "To": humanAgent{}}
	if r.URL.Query().Get("n") == "requested" {
		v["Notice"] = "요청을 보냈습니다. 가능한 요청은 상대가 확인할 수 있습니다"
	}
	own := st.humanAgents(m.Owner)
	v["Own"] = own
	friends, requests := []humanFriend{}, []humanFriend{}
	for _, f := range st.Friends {
		if f.A != id && f.B != id {
			continue
		}
		peer := f.A
		if peer == id {
			peer = f.B
		}
		name := "요청 대기"
		incoming := f.Inviter != id
		if incoming {
			name = f.NameA
			if peer == f.B {
				name = f.NameB
			}
		}
		if f.State == "active" {
			if p := st.Members[peer]; p != nil {
				if _, eligible := st.humanMember(p.Owner); eligible != nil {
					name = p.Name
				}
			}
		}
		state := f.State
		blockedByMe := f.A == id && f.BlockA || f.B == id && f.BlockB
		if blockedByMe {
			state = "blocked"
		}
		if !blockedByMe && (f.BlockA || f.BlockB) {
			state = "removed"
		}
		entry := humanFriend{ID: f.ID, Name: name, State: state, Generation: f.Generation, Incoming: incoming, Exp: clock(f.Exp)}
		if state == "active" {
			friends = append(friends, entry)
		} else {
			requests = append(requests, entry)
		}
	}
	sort.Slice(friends, func(i, j int) bool { return friends[i].ID < friends[j].ID })
	sort.Slice(requests, func(i, j int) bool { return requests[i].ID < requests[j].ID })
	v["Friends"], v["Requests"] = friends, requests
	gates := []*Gate{}
	for _, g := range st.Gates {
		if g.Owner == m.Owner && g.State == "pending" {
			gates = append(gates, g)
		}
	}
	v["PendingGates"], v["Gates"] = len(gates), gates
	for _, p := range st.Pairs {
		a, b := st.Agents[p.A], st.Agents[p.B]
		if a != nil && a.Owner == m.Owner || b != nil && b.Owner == m.Owner {
			v["Legacy"] = true
		}
	}
	conversations := []humanConversation{}
	seen := map[string]bool{}
	messages := []*Message{}
	for _, msg := range st.Messages {
		c := msg.Conversation
		if c != nil && (c.OwnerA == m.Owner || c.OwnerB == m.Owner) && now.Before(msg.Receipt.Accepted.Add(24*time.Hour)) {
			messages = append(messages, msg)
		}
	}
	sort.Slice(messages, func(i, j int) bool { return messages[i].Receipt.Accepted.After(messages[j].Receipt.Accepted) })
	for _, msg := range messages {
		key := pairID(msg.Receipt.From, msg.Receipt.To)
		if seen[key] {
			continue
		}
		seen[key] = true
		c := msg.Conversation
		from, to, person, fa, ta := msg.Receipt.From, msg.Receipt.To, c.PersonB, c.AgentA, c.AgentB
		if c.OwnerA != m.Owner {
			from, to, person, fa, ta = to, from, c.PersonA, c.AgentB, c.AgentA
		}
		if c.OwnerA == c.OwnerB {
			person = "내 에이전트"
		}
		conversations = append(conversations, humanConversation{ID: key, Name: person, FromName: fa, ToName: ta, State: humanTransport(msg.Receipt.State), URL: conversationURL(from, to)})
	}
	v["Conversations"] = conversations
	switch path {
	case "name":
		v["Title"] = "내 이름"
	case "account":
		v["Title"], v["Section"] = "내 계정", "account"
	case "agents":
		v["Title"], v["Section"] = "내 에이전트", "agents"
	case "friends":
		v["Title"], v["Section"] = "친구", "friends"
	case "friend", "new":
		v["Title"] = "새 대화"
		peerID := r.URL.Query().Get("peer")
		if path == "friend" {
			peerID = r.PathValue("id")
			v["Title"], v["Section"] = "친구", "friends"
		}
		if peerID == "self" {
			if m.SelfConnect {
				v["Agents"] = own
			}
		} else if peerID != "" {
			f := st.friendByID(id, peerID)
			if f == nil || f.State != "active" || f.BlockA || f.BlockB {
				return nil, fault("sender_not_allowed")
			}
			peer := f.A
			if peer == id {
				peer = f.B
			}
			p := st.Members[peer]
			if p == nil {
				return nil, fault("sender_not_allowed")
			}
			if _, eligible := st.humanMember(p.Owner); eligible == nil {
				return nil, fault("sender_not_allowed")
			}
			v["Peer"] = humanFriend{ID: f.ID, Name: p.Name, State: f.State, Generation: f.Generation}
			v["Agents"] = st.humanAgents(p.Owner)
		}
	case "conversation":
		from, to := r.URL.Query().Get("from"), r.URL.Query().Get("to")
		history := st.historical(m.Owner, from, to, now)
		a, b := st.Agents[from], st.Agents[to]
		_, _, allowed := st.authority(from, to)
		canSend := allowed && a != nil && a.Owner == m.Owner && b != nil && activeKeys(b) > 0
		if !canSend && len(history) == 0 {
			return nil, fault("sender_not_allowed")
		}
		if canSend {
			v["From"], v["To"] = humanAgentView(from, a), humanAgentView(to, b)
			_, peer := st.humanMember(b.Owner)
			if peer != nil && b.Owner != m.Owner {
				v["Peer"] = humanFriend{Name: peer.Name}
			}
		} else {
			c := history[0].Conversation
			fa, ta, person := c.AgentA, c.AgentB, c.PersonB
			if history[0].Receipt.From != from {
				fa, ta, person = c.AgentB, c.AgentA, c.PersonA
			}
			v["From"], v["To"] = humanAgent{ID: from, Name: fa, Short: shortID(from)}, humanAgent{ID: to, Name: ta, Short: shortID(to)}
			v["Peer"] = humanFriend{Name: person}
		}
		sort.Slice(history, func(i, j int) bool { return history[i].Receipt.Accepted.Before(history[j].Receipt.Accepted) })
		rows := []humanMessage{}
		for _, msg := range history {
			c := msg.Conversation
			auth := c.Authorship
			if auth == "owner" && c.OwnerA != m.Owner {
				auth = "상대가 작성"
			}
			rows = append(rows, humanMessage{ID: msg.Receipt.ID, Author: c.PersonA, Agent: c.AgentA, Short: shortID(msg.Receipt.From), Text: c.Text, When: humanClock(msg.Receipt.Accepted, now), State: humanTransport(msg.Receipt.State), Authorship: auth})
		}
		v["Title"], v["Messages"], v["CanSend"], v["ConversationID"] = "대화", rows, canSend, pairID(from, to)
		if attempt := r.URL.Query().Get("attempt"); len(attempt) >= 16 && len(attempt) <= 96 {
			v["Attempt"] = attempt
			if message := st.Messages[st.Idempotency[from+"/web:"+attempt]]; message != nil && message.Conversation != nil && message.Conversation.OwnerA == m.Owner && message.Receipt.To == to {
				v["Notice"] = "이 발송은 이미 접수됐습니다. 대화에서 전달 상태를 확인하세요."
				v["Text"] = message.Conversation.Text
			}
		}
	}
	return v, nil
}

func humanTransport(state string) string {
	switch state {
	case "queued", "leased":
		return "pending"
	case "delivered":
		return "delivered"
	case "failed:expired":
		return "expired"
	}
	return "failed"
}
func humanClock(t, now time.Time) string {
	t, now = t.In(kst), now.In(kst)
	prefix := t.Format("1월 2일 ")
	if t.Year() == now.Year() && t.YearDay() == now.YearDay() {
		prefix = "오늘 "
	} else if t.Year() != now.Year() {
		prefix = t.Format("2006년 1월 2일 ")
	}
	hour := t.Hour()
	part := "오전 "
	if hour >= 12 {
		part = "오후 "
	}
	hour %= 12
	if hour == 0 {
		hour = 12
	}
	return prefix + part + strconv.Itoa(hour) + t.Format(":04")
}

func (s *Service) humanPage(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		v, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
			id, session, retry, e := s.memberHit(st, r, now, memberRate)
			if !retry.IsZero() {
				return refusal{status: 429, problem: limited(retry)}, nil
			}
			if e != nil {
				return nil, e
			}
			if st.Members[id].Name == "" && path != "name" {
				return "/home/name", nil
			}
			return st.humanView(id, session, path, r, now)
		})
		if err != nil {
			if err.Error() == "invalid_auth" {
				expired(w, r)
			} else {
				refused(w, "대화", refusal{status: statusFor(err), problem: safeProblem(err)})
			}
			return
		}
		switch result := v.(type) {
		case refusal:
			refused(w, "대화", result)
		case string:
			http.Redirect(w, r, result, 303)
		default:
			render(w, 200, "human-"+path, result.(map[string]any))
		}
	}
}

func (s *Service) humanAction(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 16384)
		parseErr := r.ParseForm()
		v, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
			rates := memberRate
			decision := r.FormValue("decision")
			if path == "friends/decision" && decision != "accept" {
				rates = cleanupRate
			}
			id, session, retry, e := s.memberHit(st, r, now, rates)
			if !retry.IsZero() {
				return refusal{status: 429, problem: limited(retry)}, nil
			}
			if e != nil {
				return nil, e
			}
			if parseErr != nil {
				return refusal{status: 422, problem: "입력을 확인하세요."}, nil
			}
			m := st.Members[id]
			fail := func(e error) (any, error) { return refusal{status: statusFor(e), problem: safeProblem(e)}, nil }
			switch path {
			case "profile":
				name, ok := validName(r.FormValue("name"))
				if !ok {
					return map[string]any{"Title": "내 이름", "Name": r.FormValue("name"), "Problem": "이름은 1–40자로 입력하세요"}, nil
				}
				m.Name = name
				return "/home", nil
			case "agent-name":
				a := st.Agents[r.FormValue("agent")]
				name, ok := validName(r.FormValue("name"))
				if a == nil || a.Owner != m.Owner || a.Revoked {
					return fail(fault("sender_not_allowed"))
				}
				if !ok {
					view, e := st.humanView(id, session, "agents", r, now)
					if e != nil {
						return fail(e)
					}
					view["FormTemplate"], view["FormStatus"], view["Problem"], view["EditID"], view["EditName"] = "human-agents", 422, "이름은 1–40자로 입력하세요", r.FormValue("agent"), r.FormValue("name")
					return view, nil
				}
				a.Name = name
				return "/home/agents", nil
			case "friend-code":
				m.FriendCode = randomToken()
				return "/home/account", nil
			case "self-connect":
				if r.FormValue("consent") != "yes" {
					return fail(fault("invalid_schema"))
				}
				if _, eligible := st.humanMember(m.Owner); eligible == nil {
					return fail(fault("sender_not_allowed"))
				}
				m.SelfConnect = true
				return "/home/new?peer=self", nil
			case "friends/request":
				if r.FormValue("consent") != "yes" {
					return fail(fault("invalid_schema"))
				}
				st.requestFriend(id, r.FormValue("code"), now)
				return "/home/friends?n=requested", nil
			case "friends/decision":
				g, e := strconv.ParseInt(r.FormValue("generation"), 10, 64)
				if e != nil {
					return fail(fault("invalid_schema"))
				}
				if decision == "accept" && r.FormValue("consent") != "yes" {
					return fail(fault("invalid_schema"))
				}
				if e = st.decideFriend(id, r.FormValue("friend"), decision, g, now); e != nil {
					return fail(e)
				}
				return "/home/friends", nil
			case "conversation/send":
				from, to := r.FormValue("from"), r.FormValue("to")
				_, e := st.sendOwnerText(id, from, to, r.FormValue("text"), r.FormValue("attempt"), now)
				if e != nil {
					q := url.Values{"from": {from}, "to": {to}}
					copyR := r.Clone(r.Context())
					copyR.URL = &url.URL{RawQuery: q.Encode()}
					view, ve := st.humanView(id, session, "conversation", copyR, now)
					if ve != nil {
						return fail(e)
					}
					view["Problem"], view["Text"], view["Attempt"], view["FormStatus"] = safeProblem(e), r.FormValue("text"), r.FormValue("attempt"), statusFor(e)
					return view, nil
				}
				return conversationURL(from, to), nil
			}
			return fail(fault("invalid_schema"))
		})
		if err != nil {
			if err.Error() == "invalid_auth" {
				expired(w, r)
			} else {
				refused(w, "대화", refusal{status: 503, problem: safeProblem(err)})
			}
			return
		}
		switch result := v.(type) {
		case refusal:
			refused(w, "대화", result)
		case string:
			http.Redirect(w, r, result, 303)
		case map[string]any:
			template := "human-name"
			status := 422
			if name, ok := result["FormTemplate"].(string); ok {
				template = name
				status = result["FormStatus"].(int)
			}
			if path == "conversation/send" {
				template = "human-conversation"
				status = result["FormStatus"].(int)
			}
			render(w, status, template, result)
		}
	}
}
