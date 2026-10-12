// Human friendships authorize current endpoints; historical text copies never grant transport or execution authority.
package relay

import (
	"crypto/hmac"
	"fmt"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

type Friend struct {
	ID, A, B, Inviter, State       string
	NameA, NameB                   string
	Generation                     int64
	Exp                            time.Time
	BlockA, BlockB, SuppressLegacy bool
}

type ConversationCopy struct {
	OwnerA, OwnerB, PersonA, PersonB, AgentA, AgentB string
	Text, Authorship                                 string
}

func validName(raw string) (string, bool) {
	name := strings.TrimSpace(raw)
	if !utf8.ValidString(name) || utf8.RuneCountInString(name) < 1 || utf8.RuneCountInString(name) > 40 {
		return name, false
	}
	for _, r := range name {
		if unicode.IsControl(r) || unicode.Is(unicode.Bidi_Control, r) {
			return name, false
		}
	}
	return name, true
}

func agentName(a *Agent) string {
	if a == nil || a.Name == "" {
		return "이름 없는 에이전트"
	}
	return a.Name
}

func (st *State) humanMember(owner string) (string, *Member) {
	for id, m := range st.Members {
		if m.Owner == owner && m.Active && (m.Issuer == emailIssuer || m.Issuer == googleIssuer) && st.Owners[owner] != nil && st.Owners[owner].Active {
			return id, m
		}
	}
	return "", nil
}

func (st *State) migrateHuman(now time.Time) error {
	if st.HumanModel > 1 || st.HumanModel < 0 {
		return fault("unavailable")
	}
	if st.Friends == nil {
		st.Friends = map[string]*Friend{}
	}
	if st.HumanModel == 1 {
		return nil
	}
	for _, p := range st.Pairs {
		if p.State == "pending" {
			p.State = "expired"
		}
	}
	for _, m := range st.Members {
		// Existing members must explicitly opt into all of their own endpoints.
		m.SelfConnect = false
		if m.FriendCode == "" {
			m.FriendCode = randomToken()
		}
	}
	st.HumanModel = 1
	return nil
}

// authority returns a disjoint source plus generation, so accepting a friendship cannot revive an old pair's gate.
func (st *State) authority(from, to string) (string, int64, bool) {
	a, b := st.Agents[from], st.Agents[to]
	if a == nil || b == nil || a.Revoked || b.Revoked || from == to {
		return "", 0, false
	}
	oa, ob := st.Owners[a.Owner], st.Owners[b.Owner]
	if oa == nil || ob == nil || !oa.Active || !ob.Active {
		return "", 0, false
	}
	for _, m := range st.Members {
		if (m.Owner == a.Owner || m.Owner == b.Owner) && (m.Issuer == emailIssuer || m.Issuer == googleIssuer) && !m.Active {
			return "", 0, false
		}
	}
	ia, ma := st.humanMember(a.Owner)
	ib, mb := st.humanMember(b.Owner)
	if ma != nil && mb != nil {
		if ia == ib && ma.SelfConnect {
			return "self:" + ia, 1, true
		}
		if f := st.Friends[pairID(ia, ib)]; f != nil {
			if f.BlockA || f.BlockB {
				return "", 0, false
			}
			if f.State == "active" {
				return "friend:" + f.ID, f.Generation, true
			}
			if f.SuppressLegacy {
				return "", 0, false
			}
		}
	}
	p := st.Pairs[pairID(from, to)]
	if p == nil || p.State != "active" {
		return "", 0, false
	}
	return "", p.Generation, true
}

func (st *State) sweepFriends(now time.Time) {
	for _, f := range st.Friends {
		if f.State == "pending" && !now.Before(f.Exp) {
			f.State = "expired"
		}
	}
}

// requestFriend gives every code outcome the same public response; no target information escapes before acceptance.
func (st *State) requestFriend(id, code string, now time.Time) {
	m := st.Members[id]
	if m == nil || m.Name == "" {
		return
	}
	_, eligible := st.humanMember(m.Owner)
	if eligible == nil {
		return
	}
	target := ""
	for other, peer := range st.Members {
		if other != id && peer.Name != "" && peer.FriendCode != "" && hmac.Equal([]byte(peer.FriendCode), []byte(code)) {
			if _, p := st.humanMember(peer.Owner); p != nil {
				target = other
			}
		}
	}
	// An indistinguishable outgoing placeholder prevents the request list becoming a code-existence oracle.
	placeholder := "undisclosed:" + hashToken(code)
	if target == "" {
		target = placeholder
	}
	key := pairID(id, target)
	old := st.Friends[key]
	if old != nil && (old.BlockA || old.BlockB) {
		target = placeholder
		key = pairID(id, target)
		old = st.Friends[key]
	}
	if old != nil && (old.State == "pending" || old.State == "active") {
		return
	}
	active, pending, mine, theirs := 0, 0, 0, 0
	for _, f := range st.Friends {
		if f.State == "active" {
			active++
		}
		if f.State == "pending" {
			pending++
			if f.Inviter == id {
				mine++
			}
			if (f.A == target || f.B == target) && f.Inviter != target {
				theirs++
			}
		}
	}
	if theirs >= 10 {
		target = placeholder
		key = pairID(id, target)
		old = st.Friends[key]
		if old != nil && old.State == "pending" {
			return
		}
	}
	if active >= 400 || pending >= 200 || mine >= 10 || (old == nil && len(st.Friends) >= 20000) {
		return
	}
	if old == nil {
		old = &Friend{ID: randomToken()[:22], A: id, B: target}
		st.Friends[key] = old
	}
	old.Generation++
	old.Inviter, old.State, old.Exp = id, "pending", now.Add(24*time.Hour)
	if a := st.Members[old.A]; a != nil {
		old.NameA = a.Name
	}
	if b := st.Members[old.B]; b != nil {
		old.NameB = b.Name
	}
}

func (st *State) friendByID(id, publicID string) *Friend {
	for _, f := range st.Friends {
		if f.ID == publicID && (f.A == id || f.B == id) {
			return f
		}
	}
	return nil
}

func (st *State) decideFriend(id, publicID, decision string, generation int64, now time.Time) error {
	f := st.friendByID(id, publicID)
	if f == nil || f.Generation != generation {
		return fault("sender_not_allowed")
	}
	switch decision {
	case "accept":
		for _, endpoint := range []string{f.A, f.B} {
			m := st.Members[endpoint]
			if m == nil || m.Name == "" {
				return fault("sender_not_allowed")
			}
			if _, eligible := st.humanMember(m.Owner); eligible == nil {
				return fault("sender_not_allowed")
			}
		}
		if f.Inviter == id || f.BlockA || f.BlockB || (f.State != "pending" && f.State != "active") || (f.State == "pending" && !now.Before(f.Exp)) {
			return fault("sender_not_allowed")
		}
		if f.State == "active" {
			return nil
		}
		active, ca, cb := 0, 0, 0
		for _, other := range st.Friends {
			if other.State == "active" {
				active++
				if other.A == f.A || other.B == f.A {
					ca++
				}
				if other.A == f.B || other.B == f.B {
					cb++
				}
			}
		}
		if active >= 400 || ca >= 20 || cb >= 20 {
			return fault("capacity")
		}
		f.State = "active"
	case "deny", "cancel":
		if f.State != "pending" || (decision == "cancel") != (f.Inviter == id) {
			return fault("sender_not_allowed")
		}
		f.State = map[string]string{"deny": "denied", "cancel": "cancelled"}[decision]
	case "remove":
		if f.State != "active" {
			return fault("sender_not_allowed")
		}
		f.State, f.SuppressLegacy = "removed", true
	case "block":
		f.State, f.SuppressLegacy = "removed", true
		if id == f.A {
			f.BlockA = true
		} else {
			f.BlockB = true
		}
	case "unblock":
		if id == f.A {
			if !f.BlockA {
				return fault("sender_not_allowed")
			}
			f.BlockA = false
		} else {
			if !f.BlockB {
				return fault("sender_not_allowed")
			}
			f.BlockB = false
		}
	default:
		return fault("invalid_schema")
	}
	st.sweep(now)
	return nil
}

func (st *State) copyConversation(m *Message, body, authorship string) {
	a, b := st.Agents[m.Receipt.From], st.Agents[m.Receipt.To]
	_, ma := st.humanMember(a.Owner)
	_, mb := st.humanMember(b.Owner)
	if ma == nil || mb == nil {
		return
	}
	m.Conversation = &ConversationCopy{OwnerA: a.Owner, OwnerB: b.Owner, PersonA: ma.Name, PersonB: mb.Name, AgentA: agentName(a), AgentB: agentName(b), Text: body, Authorship: authorship}
}

func shortID(id string) string {
	if len(id) > 8 {
		return fmt.Sprintf("%x", []byte(hashToken(id))[:4])
	}
	return id
}
