// Owner pre-approval (D057): between two active self agents of one owner, only fixed orca.*.v1 text kinds pass without a
// per-send approval. Grants are made and ended only from the owner's web session; any miss is approval_required.
package relay

import (
	"encoding/json"
	"sort"
	"strings"
	"time"

	jcs "github.com/cyberphone/json-canonicalization/go/src/webpki.org/jsoncanonicalizer"
)

const (
	grantBodyLimit = 2048
	grantHourly    = 30
	grantTTL       = 30 * 24 * time.Hour
	grantsPerOwner = 30
)

var grantKinds = []string{"orca.task.v1", "orca.decision.v1", "orca.answer.v1"}

// Grant covers one (from, to, kind). FromKeys/ToKeys hold each agent's active keys at creation (kids are per agent, so
// they stay apart): a later revoked one ends the grant, and keys added later are not covered, so nothing revives a
// grant (the generation rule).
type Grant struct {
	ID, Owner, From, To, Kind string
	FromKeys, ToKeys          []string
	Created, Exp              time.Time
	Sent                      []time.Time `json:",omitempty"` // accepted sends in the last hour
}

func grantKind(kind string) bool {
	for _, k := range grantKinds {
		if k == kind {
			return true
		}
	}
	return false
}

func activeKids(a *Agent) []string {
	kids := []string{}
	for kid, k := range a.Keys {
		if !k.Revoked {
			kids = append(kids, kid)
		}
	}
	sort.Strings(kids)
	return kids
}

// grantLive reports whether a stored grant still has the same owner, self authority and every snapshotted key.
func (st *State) grantLive(g *Grant, now time.Time) bool {
	a, b := st.Agents[g.From], st.Agents[g.To]
	if !now.Before(g.Exp) || a == nil || b == nil || a.Owner != g.Owner || b.Owner != g.Owner {
		return false
	}
	if scope, _, ok := st.authority(g.From, g.To); !ok || !strings.HasPrefix(scope, "self:") {
		return false
	}
	for agent, kids := range map[*Agent][]string{a: g.FromKeys, b: g.ToKeys} {
		for _, kid := range kids {
			if k := agent.Keys[kid]; k == nil || k.Revoked {
				return false
			}
		}
	}
	return true
}

func (st *State) sweepGrants(now time.Time) {
	for id, g := range st.Grants {
		if !st.grantLive(g, now) {
			delete(st.Grants, id)
			continue
		}
		kept := g.Sent[:0]
		for _, t := range g.Sent {
			if now.Sub(t) < time.Hour {
				kept = append(kept, t)
			}
		}
		g.Sent = kept
	}
}

// createGrant replaces the same (from, to, kind) with a fresh 30-day grant (renewal is a new owner consent); the hourly
// count carries over so renewing never resets the rate.
func (st *State) createGrant(owner, from, to, kind string, now time.Time) error {
	a, b := st.Agents[from], st.Agents[to]
	if !grantKind(kind) || a == nil || b == nil || a.Owner != owner || b.Owner != owner {
		return fault("sender_not_allowed")
	}
	if scope, _, ok := st.authority(from, to); !ok || !strings.HasPrefix(scope, "self:") {
		return fault("sender_not_allowed")
	}
	fromKeys, toKeys := activeKids(a), activeKids(b)
	if len(fromKeys) == 0 || len(toKeys) == 0 {
		return fault("sender_not_allowed")
	}
	var sent []time.Time
	count := 0
	for id, g := range st.Grants {
		if g.Owner != owner {
			continue
		}
		if g.From == from && g.To == to && g.Kind == kind {
			sent = g.Sent
			delete(st.Grants, id)
			continue
		}
		count++
	}
	if count >= grantsPerOwner {
		return fault("capacity")
	}
	id := randomToken()
	st.Grants[id] = &Grant{ID: id, Owner: owner, From: from, To: to, Kind: kind, FromKeys: fromKeys, ToKeys: toKeys, Created: now, Exp: now.Add(grantTTL), Sent: sent}
	return nil
}

func (st *State) revokeGrant(owner, id string) error {
	g := st.Grants[id]
	if g == nil || g.Owner != owner {
		return fault("sender_not_allowed")
	}
	delete(st.Grants, id)
	return nil
}

// grantFor returns the grant that admits this signed text, or nil; it never says which check failed.
func (st *State) grantFor(e *PublicText, scope string, now time.Time) *Grant {
	if !grantKind(e.Kind) || !strings.HasPrefix(scope, "self:") || !validGrantBody(e.Kind, e.Text) {
		return nil
	}
	owner := st.Agents[e.From].Owner
	for _, g := range st.Grants {
		if g.Owner != owner || g.From != e.From || g.To != e.To || g.Kind != e.Kind || !st.grantLive(g, now) || len(g.Sent) >= grantHourly {
			continue
		}
		for _, kid := range g.FromKeys {
			if kid == e.Sig.Kid {
				return g
			}
		}
	}
	return nil
}

// validGrantBody accepts only the canonical JSON of one fixed body shape within 2 KB.
func validGrantBody(kind, text string) bool {
	if len(text) > grantBodyLimit {
		return false
	}
	canonical, err := jcs.Transform([]byte(text))
	if err != nil || string(canonical) != text {
		return false
	}
	var body map[string]any
	if json.Unmarshal([]byte(text), &body) != nil {
		return false
	}
	filled := func(key string, max int) bool {
		s, ok := body[key].(string)
		return ok && strings.TrimSpace(s) != "" && len([]rune(s)) <= max
	}
	optional := func(key string, check func() bool) bool {
		_, present := body[key]
		return !present || check()
	}
	item := func() bool {
		n, ok := body["item"].(float64)
		return ok && n >= 1 && n <= 2147483647 && n == float64(int64(n))
	}
	switch kind {
	case "orca.task.v1":
		return object(body, []string{"text"}, "project") && filled("text", grantBodyLimit) && optional("project", func() bool { return filled("project", 200) })
	case "orca.decision.v1":
		options := func() bool {
			list, ok := body["options"].([]any)
			if !ok || len(list) == 0 || len(list) > 9 {
				return false
			}
			for _, o := range list {
				s, ok := o.(string)
				if !ok || strings.TrimSpace(s) == "" || len([]rune(s)) > 200 {
					return false
				}
			}
			return true
		}
		return object(body, []string{"bot", "question"}, "item", "options") && filled("bot", 80) && filled("question", grantBodyLimit) && optional("item", item) && optional("options", options)
	case "orca.answer.v1":
		return object(body, []string{"answer", "by"}, "item") && filled("answer", grantBodyLimit) && filled("by", 80) && optional("item", item)
	}
	return false
}

// humanGrant is the owner's view row: names and kind only, never a body.
type humanGrant struct {
	ID, From, To, Kind, Exp string
	Sent                    int
}

func (st *State) ownGrants(owner string, now time.Time, lang string) []humanGrant {
	rows := []humanGrant{}
	for _, g := range st.Grants {
		if g.Owner == owner {
			rows = append(rows, humanGrant{g.ID, agentName(st.Agents[g.From]), agentName(st.Agents[g.To]), g.Kind, humanClock(g.Exp, now, lang), len(g.Sent)})
		}
	}
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].From+rows[i].To+rows[i].Kind < rows[j].From+rows[j].To+rows[j].Kind
	})
	return rows
}
