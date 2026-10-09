// Client onboarding checks Google binding, explicit approval, independent keys, expiry, replay and private-key possession.
package relay

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func deviceFixture(t *testing.T, st *State, now time.Time) (string, *Connection, ed25519.PrivateKey) {
	t.Helper()
	public, private, _ := ed25519.GenerateKey(rand.Reader)
	token := randomToken()
	c := &Connection{Client: supportedClient, Kid: "device_key", Public: base64.RawURLEncoding.EncodeToString(public)}
	proof := base64.RawURLEncoding.EncodeToString(ed25519.Sign(private, deviceBytes(token, c)))
	result, err := st.startDevice(token, c.Client, c.Kid, c.Public, proof, now)
	if err != nil {
		t.Fatal(err)
	}
	return token, result, private
}

func TestDeviceGoogleConnection(t *testing.T) {
	st, now := connectionFixture()
	foreignMember := st.signIn(googleIssuer, "existing-foreign", "foreign@example.test", "", "", now)
	if foreignMember.Problem != "" {
		t.Fatal("foreign member fixture")
	}
	var agents, credentials []string
	for range 2 {
		token, c, private := deviceFixture(t, st, now)
		id := hashToken(token)
		attempt := &googleAttempt{Connection: id, Exp: now.Add(connectionTTL)}
		if err := st.approveDevice(id, "owner", now); err == nil {
			t.Fatal("approval before Google identity")
		}
		login := st.finishGoogleLogin(attempt, "same-subject", "same@example.test", now)
		if login.Problem != "" || c.State != "prepared" || c.Owner == "" || c.Agent != "" {
			t.Fatal("Google must bind identity without granting agent access")
		}
		owner := c.Owner
		members := len(st.Members)
		if foreign := st.finishGoogleLogin(attempt, "foreign-subject", "foreign@example.test", now); foreign.Problem == "" || len(st.Members) != members || c.Owner != owner {
			t.Fatal("other Google account claimed connection")
		}
		sessions := len(st.Sessions)
		if foreign := st.finishGoogleLogin(attempt, "existing-foreign", "foreign@example.test", now); foreign.Problem == "" || len(st.Members) != members || len(st.Sessions) != sessions || c.Owner != owner || c.State != "prepared" {
			t.Fatal("existing foreign member replaced binding or gained a session")
		}
		if err := st.approveDevice(id, "other-owner", now); err == nil {
			t.Fatal("foreign approval")
		}
		if _, err := st.completeConnection(token, supportedClient, "", now); err == nil {
			t.Fatal("credential before consent")
		}
		if err := st.approveDevice(id, owner, now); err != nil {
			t.Fatal(err)
		}
		if err := st.approveDevice(id, owner, now); err == nil {
			t.Fatal("approval replay created an agent")
		}
		_, attacker, _ := ed25519.GenerateKey(rand.Reader)
		badProof := base64.RawURLEncoding.EncodeToString(ed25519.Sign(attacker, connectionBytes(token, c)))
		if _, err := st.completeConnection(token, supportedClient, badProof, now); err == nil {
			t.Fatal("foreign client collected credential")
		}
		proof := base64.RawURLEncoding.EncodeToString(ed25519.Sign(private, connectionBytes(token, c)))
		if _, err := st.completeConnection(token, "other-client", proof, now); err == nil {
			t.Fatal("client substitution")
		}
		value, err := st.completeConnection(token, supportedClient, proof, now)
		if err != nil {
			t.Fatal(err)
		}
		credential := value.(map[string]string)["credential"]
		if agent, err := st.principal(credential, "agent"); err != nil || agent != c.Agent {
			t.Fatal("wrong credential scope")
		}
		if _, err := st.completeConnection(token, supportedClient, proof, now); err == nil {
			t.Fatal("completion replay")
		}
		if st.finishGoogleLogin(attempt, "same-subject", "same@example.test", now).Problem == "" {
			t.Fatal("login replay reopened consumed request")
		}
		agents, credentials = append(agents, c.Agent), append(credentials, credential)
	}
	if agents[0] == agents[1] || credentials[0] == credentials[1] || len(st.Members) != 2 {
		t.Fatal("clients must share member but have independent agents and credentials")
	}
	if len(st.Pairs) != 0 {
		t.Fatal("onboarding silently paired agents")
	}
}

func TestDeviceCapacityRetention(t *testing.T) {
	st, now := connectionFixture()
	token, c, private := deviceFixture(t, st, now)
	for i := range 1998 {
		state := []string{"requested", "prepared", "consumed", "cancelled"}[i%4]
		st.Connections[randomToken()] = &Connection{Device: true, State: state, Exp: c.Exp}
	}
	// Member-started connections do not consume the anonymous device budget.
	st.Connections["member"] = &Connection{State: "waiting", Exp: c.Exp}
	newToken := randomToken()
	proof := base64.RawURLEncoding.EncodeToString(ed25519.Sign(private, deviceBytes(newToken, c)))
	if _, err := st.startDevice(newToken, c.Client, c.Kid, c.Public, proof, now); err != nil {
		t.Fatal("member connections consumed device capacity or 2000th request was rejected")
	}
	newToken = randomToken()
	proof = base64.RawURLEncoding.EncodeToString(ed25519.Sign(private, deviceBytes(newToken, c)))
	retained := c.Exp.Add(24*time.Hour - time.Nanosecond)
	st.sweepConnections(retained)
	if _, err := st.startDevice(newToken, c.Client, c.Kid, c.Public, proof, retained); err == nil || len(st.Connections) != 2001 || st.Connections[hashToken(token)] == nil {
		t.Fatal("capacity or 24h tombstone retention changed")
	}
	cleared := c.Exp.Add(24 * time.Hour)
	st.sweepConnections(cleared)
	if _, err := st.startDevice(newToken, c.Client, c.Kid, c.Public, proof, cleared); err != nil || len(st.Connections) != 1 {
		t.Fatal("anonymous capacity did not recover after retention")
	}
}

func TestDeviceFailureBoundaries(t *testing.T) {
	st, now := connectionFixture()
	token, c, private := deviceFixture(t, st, now)
	proof := base64.RawURLEncoding.EncodeToString(ed25519.Sign(private, deviceBytes(token, c)))
	if !deviceProof(token, proof, c) || deviceProof("other-token", proof, c) {
		t.Fatal("poll proof not bound to token")
	}
	copy := *c
	copy.Client = "other-client"
	if deviceProof(token, proof, &copy) {
		t.Fatal("poll proof not bound to client")
	}
	if _, err := st.startDevice(token, c.Client, c.Kid, c.Public, proof, now); err == nil {
		t.Fatal("start replay")
	}
	if _, err := st.startDevice(randomToken(), c.Client, c.Kid, c.Public, proof, now); err == nil {
		t.Fatal("start proof substitution")
	}
	attempt := &googleAttempt{Connection: hashToken(token), Exp: now.Add(connectionTTL)}
	if st.finishGoogleLogin(attempt, "subject", "member@example.test", c.Exp).Problem == "" {
		t.Fatal("expired Google binding")
	}
	st.finishGoogleLogin(attempt, "subject", "member@example.test", now)
	if st.approveDevice(hashToken(token), c.Owner, c.Exp) == nil {
		t.Fatal("expired consent")
	}
	if st.approveDevice(hashToken(token), c.Owner, now) != nil {
		t.Fatal("approval")
	}
	completion := base64.RawURLEncoding.EncodeToString(ed25519.Sign(private, connectionBytes(token, c)))
	if _, err := st.completeConnection(token, supportedClient, completion, c.Exp); err == nil {
		t.Fatal("expired completion")
	}
	st.Owners[c.Owner].Active = false
	if _, err := st.completeConnection(token, supportedClient, completion, now); err == nil {
		t.Fatal("inactive owner")
	}
	for len(st.Connections) < 2000 {
		st.Connections[randomToken()] = &Connection{Device: true}
	}
	token = randomToken()
	proof = base64.RawURLEncoding.EncodeToString(ed25519.Sign(private, deviceBytes(token, c)))
	if _, err := st.startDevice(token, c.Client, c.Kid, c.Public, proof, now); err == nil {
		t.Fatal("unbounded anonymous storage")
	}
}

func TestDeviceConsentTemplate(t *testing.T) {
	for _, state := range []string{"login", "prepared", "approved"} {
		w := httptest.NewRecorder()
		render(w, 200, "device", map[string]any{"Title": "connect", "ConnectionID": "public-id", "Fingerprint": "SHA256:public", "State": state, "Login": state == "login", "Google": true})
		body := w.Body.String()
		if !strings.Contains(body, "SHA256:public") || strings.Contains(body, "<no value>") {
			t.Fatal("missing fingerprint")
		}
		if state == "prepared" && !strings.Contains(body, "/home/device-confirm") {
			t.Fatal("missing explicit consent")
		}
		if state == "login" && !strings.Contains(body, "/auth/google") {
			t.Fatal("missing Google login")
		}
		if state == "approved" && strings.Contains(body, "<form") {
			t.Fatal("approved connection offers repeated consent")
		}
	}
}
