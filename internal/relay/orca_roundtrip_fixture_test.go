//go:build integration && orcafixture

// Opt-in loopback S6 fixture (SAR-ORCA-BOT-001-DEV): one synthetic email member with two self agents (Grok and Orca
// roles) whose key folders are written for the Node adapter, the real relay handlers, and the member session token in
// the manifest so the grant is made through the real web route. No real account, mail or Orca bot is involved.
package relay

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestOrcaRoundtripFixture(t *testing.T) {
	dir, manifest := os.Getenv("ORCA_FIXTURE_DIR"), os.Getenv("ORCA_FIXTURE_MANIFEST")
	if os.Getenv("TEST_SYNTHETIC_DATABASE") != "1" || dir == "" || manifest == "" {
		t.Fatal("explicit isolated fixture configuration required")
	}
	const life = 2 * time.Hour
	s, _ := identityService(t, messagePool(t))
	server := httptest.NewServer(s.Handler())
	defer server.Close()
	email := "orca-fixture@example.test"
	var session string
	agents := map[string]string{"grok": "agent_fx_grok", "orca": "agent_fx_orca"}
	s.mutateState(t, func(st *State) {
		now := time.Now()
		v := st.signIn(emailIssuer, email, email, "", "", now)
		session = v.Session
		m := st.Members[st.Sessions[hashToken(v.Session)].Member]
		m.Name, m.SelfConnect = "합성 주인", true
		for role, id := range agents {
			public, private, _ := ed25519.GenerateKey(rand.Reader)
			credential := randomToken()
			st.Agents[id] = &Agent{Name: "합성 " + role, Client: supportedClient, Owner: m.Owner, Keys: map[string]*Key{"fixture-key": {Public: public, Changed: now, Credential: hashToken(credential)}}, Created: now}
			// The role folder already exists with a private ACL; files written inside inherit it.
			der, _ := x509.MarshalPKCS8PrivateKey(private)
			conf, _ := json.Marshal(map[string]string{"relay": server.URL + "/", "agent": id, "kid": "fixture-key", "credential": credential, "name": "합성 " + role})
			if err := os.WriteFile(filepath.Join(dir, role, "private.pem"), pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der}), 0600); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(dir, role, "agent.json"), conf, 0600); err != nil {
				t.Fatal(err)
			}
		}
	})
	data, _ := json.Marshal(map[string]any{"origin": server.URL, "agents": agents, "session_cookie": sessionCookie, "session": session})
	if err := os.WriteFile(manifest, data, 0600); err != nil {
		t.Fatal(err)
	}
	t.Log("orca fixture ready; synthetic account only")
	<-time.After(life)
}
