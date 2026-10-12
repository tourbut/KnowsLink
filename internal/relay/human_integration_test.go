//go:build integration

// Actual synthetic OTP, browser forms, signed APIs and Postgres prove the complete person-to-person conversation flow.
package relay

import (
	"context"
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestHumanPostgres(t *testing.T) {
	if os.Getenv("TEST_SYNTHETIC_DATABASE") != "1" || os.Getenv("TEST_DATABASE_URL") == "" {
		t.Fatal("isolated database required")
	}
	pool, e := pgxpool.New(context.Background(), os.Getenv("TEST_DATABASE_URL"))
	if e != nil {
		t.Fatal(e)
	}
	defer pool.Close()
	s, mail := identityService(t, pool)
	h := s.Handler()
	browsers := map[string]*browser{}
	members := map[string]string{}
	for i, person := range []string{"a", "b", "c"} {
		b := &browser{t, h, fmt.Sprintf("192.0.2.%d", i+11), map[string]string{}}
		browsers[person] = b
		email := person + "@humanpeer.invalid"
		expect(t, b.do("POST", "/auth/start", url.Values{"email": {email}}), 303)
		expect(t, b.do("POST", "/auth/verify", url.Values{"code": {mail.code(email)}}), 303)
		expect(t, b.do("GET", "/home", nil), 303, "/home/name")
		expect(t, b.do("POST", "/home/profile", url.Values{"name": {"합성 회원 " + person}}), 303)
		expect(t, b.do("GET", "/home", nil), 200, "대화", "첫 에이전트를 추가하세요")
		s.mutateState(t, func(st *State) {
			id, _, e := st.session(b.cookies[sessionCookie], time.Now())
			if e != nil {
				t.Fatal(e)
			}
			members[person] = id
		})
	}
	a, b, c := browsers["a"], browsers["b"], browsers["c"]
	agents, tokens := map[string]string{}, map[string]string{}
	keys := map[string]ed25519.PrivateKey{}
	create := func(person, label string) string {
		t.Helper()
		client := browsers[person]
		name := "합성 도우미 " + label
		expect(t, client.do("POST", "/home/agents", url.Values{"name": {name}}), 303, "/home/agents")
		agent := ""
		s.mutateState(t, func(st *State) {
			for id, a := range st.Agents {
				if a.Owner == st.Members[members[person]].Owner && a.Name == name {
					agent = id
				}
			}
		})
		if agent == "" {
			t.Fatal("created endpoint missing")
		}
		token, id := browserGrant(t, client, agent, "register")
		key, proof := clientPrepare(t, h, token, "key1")
		expect(t, client.do("POST", "/home/confirm", url.Values{"connection": {id}}), 303)
		result := publicCall(t, h, "/v1/connect/complete", "", map[string]string{"token": token, "client": supportedClient, "proof": proof}, 200)
		tokens[label] = result["credential"].(string)
		keys[label] = key
		agents[label] = agent
		return agent
	}
	aa, ab, cc := create("a", "a"), create("b", "b"), create("c", "c")
	a2 := create("a", "a2")
	expect(t, a.do("POST", "/home/conversation/send", url.Values{"from": {aa}, "to": {a2}, "text": {"합성 자기 연결"}, "attempt": {"human-self-send-0001"}}), 303)
	selfLease := publicCall(t, h, "/v1/text/pull", tokens["a2"], map[string]string{}, 200)
	selfDelivery := map[string]string{"id": selfLease["envelope"].(map[string]any)["id"].(string), "token": selfLease["lease_token"].(string)}
	publicCall(t, h, "/v1/text/persist", tokens["a2"], selfDelivery, 200)
	publicCall(t, h, "/v1/text/ack", tokens["a2"], selfDelivery, 200)
	var code string
	var friend *Friend
	s.mutateState(t, func(st *State) { code = st.Members[members["b"]].FriendCode; st.Rates = map[string][]time.Time{} })
	form := url.Values{"from": {aa}, "to": {ab}, "text": {"합성 테스트"}, "attempt": {"human-before-consent"}}
	expect(t, a.do("POST", "/home/conversation/send", form), 403)
	expect(t, a.do("POST", "/home/friends/request", url.Values{"code": {code}, "consent": {"yes"}}), 303, "n=requested")
	s.mutateState(t, func(st *State) { friend = st.Friends[pairID(members["a"], members["b"])] })
	expect(t, a.do("GET", "/home/friends/"+friend.ID, nil), 403)
	expect(t, b.do("GET", "/home/friends", nil), 200, "합성 회원 a", "서로의 현재")
	decision := url.Values{"friend": {friend.ID}, "generation": {fmt.Sprint(friend.Generation)}, "decision": {"accept"}, "consent": {"yes"}}
	expect(t, c.do("POST", "/home/friends/decision", decision), 403)
	expect(t, a.do("POST", "/home/friends/decision", decision), 403)
	expect(t, b.do("POST", "/home/friends/decision", decision, "Origin", "https://foreign.invalid"), 403)
	expect(t, b.do("POST", "/home/friends/decision", decision), 303)
	view := a.do("GET", "/home/friends/"+friend.ID, nil)
	expect(t, view, 200, "합성 회원 b", "합성 도우미 b")
	for _, forbidden := range []string{"@humanpeer.invalid", members["b"], "SHA256:", tokens["b"]} {
		if strings.Contains(view.Body.String(), forbidden) {
			t.Fatal("friend view disclosed private information")
		}
	}
	newAgent := create("b", "b2")
	expect(t, a.do("GET", "/home/friends/"+friend.ID, nil), 200, "합성 도우미 b2")
	publicCall(t, h, "/v1/invites", tokens["a"], map[string]string{"agent": aa, "target": newAgent}, 403)
	keyRequest := httptest.NewRequest("GET", "/v1/keys/"+cc+"/key1", nil)
	keyRequest.Header.Set("Authorization", "Bearer "+tokens["a"])
	keyResponse := httptest.NewRecorder()
	h.ServeHTTP(keyResponse, keyRequest)
	expect(t, keyResponse, 403)
	for i := 0; i < 6; i++ {
		person, from, to, token := "a", aa, ab, tokens["b"]
		if i%2 == 1 {
			person, from, to, token = "b", ab, aa, tokens["a"]
		}
		attempt := fmt.Sprintf("human-continuous-%04d", i)
		expect(t, browsers[person].do("POST", "/home/conversation/send", url.Values{"from": {from}, "to": {to}, "text": {fmt.Sprintf("합성 왕복 %d", i)}, "attempt": {attempt}}), 303)
		lease := publicCall(t, h, "/v1/text/pull", token, map[string]string{}, 200)
		envelope := lease["envelope"].(map[string]any)
		if envelope["authorship"] != "owner" || envelope["v"] != "knowslink.owner-text.v1" || envelope["sig"] != nil {
			t.Fatal("owner falsely signed")
		}
		payload := map[string]string{"id": envelope["id"].(string), "token": lease["lease_token"].(string)}
		publicCall(t, h, "/v1/text/ack", token, payload, 409)
		publicCall(t, h, "/v1/text/persist", token, payload, 200)
		publicCall(t, h, "/v1/text/ack", token, payload, 200)
		if i == 0 {
			s.mutateState(t, func(st *State) { msg := st.Messages[payload["id"]]; msg.Receipt.Exp = time.Now().Add(-time.Second) })
		}
	}
	s.mutateState(t, func(st *State) { st.Rates = map[string][]time.Time{} })
	chat := a.do("GET", conversationURL(aa, ab), nil)
	expect(t, chat, 200, "내가 작성", "상대가 작성", "전달됨", "24시간")
	if strings.Contains(chat.Body.String(), "읽음") {
		t.Fatal("invented read receipt")
	}
	signed := signedText(remoteUUID(time.Now()), ab, aa, "human-signed-path-01", "", "합성 agent 작성", time.Now().Add(180*time.Second), keys["b"])
	raw, _ := json.Marshal(signed)
	publicCall(t, h, "/v1/text/send", tokens["b"], json.RawMessage(raw), 200)
	expect(t, a.do("GET", conversationURL(aa, ab), nil), 200, "에이전트가 작성")
	form = url.Values{"from": {aa}, "to": {ab}, "text": {"합성 멱등"}, "attempt": {"human-idempotent-01"}}
	expect(t, a.do("POST", "/home/conversation/send", form), 303)
	expect(t, a.do("POST", "/home/conversation/send", form), 303)
	form.Set("text", "changed")
	expect(t, a.do("POST", "/home/conversation/send", form), 409)
	form.Set("attempt", "human-forged-sender")
	form.Set("from", ab)
	form.Set("to", aa)
	expect(t, a.do("POST", "/home/conversation/send", form), 403)
	form.Set("from", aa)
	form.Set("to", ab)
	form.Set("text", strings.Repeat("한", 1366))
	expect(t, a.do("POST", "/home/conversation/send", form), 422)
	expect(t, c.do("GET", conversationURL(aa, ab), nil), 403)
	lease := publicCall(t, h, "/v1/text/pull", tokens["a"], map[string]string{}, 200)
	pending := lease["envelope"].(map[string]any)["id"].(string)
	decision.Set("decision", "remove")
	expect(t, a.do("POST", "/home/friends/decision", decision), 303)
	publicCall(t, h, "/v1/text/persist", tokens["a"], map[string]string{"id": pending, "token": lease["lease_token"].(string)}, 409)
	expect(t, a.do("GET", conversationURL(aa, ab), nil), 200, "지금은 메시지를 보낼 수 없습니다", "합성 왕복")
	form.Set("text", "after removal")
	expect(t, a.do("POST", "/home/conversation/send", form), 403)
	// The schema trigger rejects the old writer and preserves already committed deny state.
	if _, e = pool.Exec(context.Background(), `UPDATE relay_state SET data=data-'HumanModel'-'Friends'`); e == nil {
		t.Fatal("old writer downgrade accepted")
	}
	s.mutateState(t, func(st *State) {
		if st.Friends[pairID(members["a"], members["b"])].State != "removed" {
			t.Fatal("deny state lost")
		}
		for _, msg := range st.Messages {
			msg.Receipt.Accepted = time.Now().Add(-24 * time.Hour)
		}
		st.Rates = map[string][]time.Time{}
	})
	expect(t, a.do("GET", conversationURL(aa, ab), nil), 403)
	guest := &browser{t, h, "192.0.2.99", map[string]string{}}
	expect(t, guest.do("GET", conversationURL(aa, ab), nil), 303)
	for _, path := range []string{"/home/friends", "/home/account", "/home/agents"} {
		r := httptest.NewRequest(http.MethodGet, path, nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Header().Get("Content-Security-Policy") == "" || w.Header().Get("Cache-Control") != "no-store" {
			t.Fatal("security headers missing")
		}
	}
}

// A connection-private SQL fixture exercises an old JSON reader path without replacing another test's singleton.
func TestHumanLegacyPostgres(t *testing.T) {
	if os.Getenv("TEST_SYNTHETIC_DATABASE") != "1" || os.Getenv("TEST_DATABASE_URL") == "" {
		t.Fatal("isolated database required")
	}
	cfg, err := pgxpool.ParseConfig(os.Getenv("TEST_DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}
	cfg.MaxConns = 1
	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	ctx := context.Background()
	for _, sql := range []string{
		`CREATE TEMP TABLE relay_state (LIKE public.relay_state INCLUDING ALL)`,
		`CREATE TRIGGER human_writer_guard BEFORE UPDATE OF data ON relay_state FOR EACH ROW EXECUTE FUNCTION public.guard_human_writer()`,
	} {
		if _, err := pool.Exec(ctx, sql); err != nil {
			t.Fatal(err)
		}
	}
	legacy, ids, _, now := humanFixture(t)
	legacy.HumanModel = 0
	for _, m := range legacy.Members {
		m.Name, m.FriendCode = "", ""
	}
	for _, a := range legacy.Agents {
		a.Name = ""
	}
	legacy.Pairs[pairID("agent_a1", "agent_b1")] = &Pair{A: "agent_a1", B: "agent_b1", State: "active", Generation: 7}
	legacy.Pairs[pairID("agent_a2", "agent_b2")] = &Pair{A: "agent_a2", B: "agent_b2", State: "pending", Generation: 3, Exp: now.Add(time.Hour)}
	raw, err := json.Marshal(legacy)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `INSERT INTO relay_state(singleton,data) VALUES(true,$1)`, raw); err != nil {
		t.Fatal(err)
	}
	s := &Service{Pool: pool}
	_, err = s.transaction(ctx, func(st *State, at time.Time) (any, error) {
		if st.HumanModel != 1 || st.Members[ids["a"]].SelfConnect || st.Members[ids["a"]].FriendCode == "" || st.Agents["agent_a1"].Keys["key1"].Credential != legacy.Agents["agent_a1"].Keys["key1"].Credential {
			t.Fatal("SQL migration changed key/identity or widened consent")
		}
		if _, gen, ok := st.authority("agent_a1", "agent_b1"); !ok || gen != 7 {
			t.Fatal("SQL exact legacy authority lost")
		}
		if _, _, ok := st.authority("agent_a2", "agent_b2"); ok || st.Pairs[pairID("agent_a2", "agent_b2")].State != "expired" {
			t.Fatal("SQL pending became authority")
		}
		if _, _, ok := st.authority("agent_a1", "agent_a2"); ok {
			t.Fatal("SQL legacy self widened")
		}
		if _, err := st.operateAs("invites", st.Members[ids["a"]].Owner, "owner", command{Agent: "agent_a1", Target: "agent_b2"}, at, false); err == nil {
			t.Fatal("SQL legacy public grant accepted")
		}
		return nil, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = pool.Exec(ctx, `UPDATE relay_state SET data='{}'::jsonb`); err == nil {
		t.Fatal("SQL old writer discarded migration")
	}
	var version int
	if err = pool.QueryRow(ctx, `SELECT (data->>'HumanModel')::integer FROM relay_state`).Scan(&version); err != nil || version != 1 {
		t.Fatal("SQL migration not durable", err)
	}
}
