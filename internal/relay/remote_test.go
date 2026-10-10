// Remote protocol checks cover installation isolation, PKCE, private Inbox durability, callback security and delivery leases.
package relay

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/netip"
	"net/url"
	"strings"
	"testing"
	"time"
)

func remoteFixture(t *testing.T, st *State, owner string, now time.Time) (string, *remoteInstall, string, url.Values) {
	t.Helper()
	r := st.remote()
	r.Clients["client"] = &remoteClient{Name: "host", Redirects: []string{"https://chatgpt.com/callback"}, Exp: now.Add(remoteLifetime)}
	verifier := strings.Repeat("A", 43)
	sum := sha256.Sum256([]byte(verifier))
	a := &remoteAttempt{Client: "client", Redirect: "https://chatgpt.com/callback", Challenge: base64.RawURLEncoding.EncodeToString(sum[:]), Resource: "https://relay.example/mcp", Owner: owner, Exp: now.Add(codeTTL)}
	code, e := st.approveRemote(a, now)
	if e != nil {
		t.Fatal(e)
	}
	form := url.Values{"grant_type": {"authorization_code"}, "client_id": {"client"}, "code": {code}, "redirect_uri": {a.Redirect}, "resource": {a.Resource}, "code_verifier": {verifier}}
	v, e := st.exchangeRemote(form, a.Resource, now)
	if e != nil {
		t.Fatal(e)
	}
	token := v.(map[string]any)["access_token"].(string)
	form.Set("refresh_token", v.(map[string]any)["refresh_token"].(string))
	id, i, e := st.remotePrincipal(token, a.Resource, now)
	if e != nil {
		t.Fatal(e)
	}
	return id, i, token, form
}

func TestRemoteOAuthIsolationAndRefresh(t *testing.T) {
	st, _, now := textFixture()
	original := st.Agents["agent_a"].Keys["key1"]
	id, a, token, form := remoteFixture(t, st, "agent_a", now)
	_, b, other, _ := remoteFixture(t, st, "agent_a", now)
	if a.Agent == b.Agent || a.Kid == b.Kid || bytes.Equal(a.Private, b.Private) || st.Agents["agent_a"].Keys["key1"] != original {
		t.Fatal("installations share keys or changed local key")
	}
	if _, e := st.principal(token, "agent"); e == nil {
		t.Fatal("remote bearer escaped to local API")
	}
	if _, _, e := st.remotePrincipal(token, "https://foreign.example/mcp", now); e == nil {
		t.Fatal("wrong audience")
	}
	if _, _, e := st.remotePrincipal(token, form.Get("resource"), now.Add(time.Hour)); e == nil {
		t.Fatal("expired access token")
	}
	if _, e := st.exchangeRemote(form, form.Get("resource"), now); e == nil {
		t.Fatal("code replay")
	}
	var refresh string
	// Mint another grant with known response to test rotation and client binding.
	r := st.remote()
	verifier := strings.Repeat("B", 43)
	sum := sha256.Sum256([]byte(verifier))
	code := randomToken()
	r.Codes[hashToken(code)] = &remoteCode{Install: id, Client: "client", Redirect: form.Get("redirect_uri"), Challenge: base64.RawURLEncoding.EncodeToString(sum[:]), Resource: form.Get("resource"), Exp: now.Add(time.Minute)}
	form.Set("code", code)
	form.Set("code_verifier", "wrong")
	if _, e := st.exchangeRemote(form, form.Get("resource"), now); e == nil {
		t.Fatal("PKCE bypass")
	}
	form.Set("code_verifier", verifier)
	form.Set("redirect_uri", "https://evil.example/callback")
	if _, e := st.exchangeRemote(form, form.Get("resource"), now); e == nil {
		t.Fatal("redirect binding bypass")
	}
	form.Set("redirect_uri", "https://chatgpt.com/callback")
	v, e := st.exchangeRemote(form, form.Get("resource"), now)
	if e != nil {
		t.Fatal(e)
	}
	refresh = v.(map[string]any)["refresh_token"].(string)
	rf := url.Values{"grant_type": {"refresh_token"}, "client_id": {"wrong"}, "refresh_token": {refresh}, "resource": {form.Get("resource")}}
	if _, e := st.exchangeRemote(rf, form.Get("resource"), now); e == nil {
		t.Fatal("refresh client binding bypass")
	}
	rf.Set("client_id", "client")
	rotated, e := st.exchangeRemote(rf, form.Get("resource"), now)
	if e != nil {
		t.Fatal(e)
	}
	access := rotated.(map[string]any)["access_token"].(string)
	st.Agents[a.Agent].Keys[a.Kid].Revoked = true
	st.sweepRemote(now)
	if _, _, e := st.remotePrincipal(access, form.Get("resource"), now); e == nil || st.Remote.Installs[id] != nil {
		t.Fatal("key revocation kept private key or authority")
	}
	if _, _, e := st.remotePrincipal(other, form.Get("resource"), now); e != nil {
		t.Fatal("revocation affected different installation")
	}
}

func TestRemoteRefreshReplayRevokesOnlyVerifiedFamily(t *testing.T) {
	st, _, now := textFixture()
	id, _, _, form := remoteFixture(t, st, "agent_a", now)
	_, _, other, otherForm := remoteFixture(t, st, "agent_b", now)
	rf := url.Values{"grant_type": {"refresh_token"}, "client_id": {"client"}, "resource": {form.Get("resource")}, "refresh_token": {form.Get("refresh_token")}}
	v, err := st.exchangeRemote(rf, form.Get("resource"), now)
	if err != nil {
		t.Fatal(err)
	}
	access := v.(map[string]any)["access_token"].(string)
	old := rf.Get("refresh_token")
	rf.Set("refresh_token", old[:len(old)-1]+"!")
	if _, err = st.exchangeRemote(rf, form.Get("resource"), now); err == nil || st.Remote.Installs[id] == nil {
		t.Fatal("forged family token revoked authority")
	}
	rf.Set("refresh_token", old)
	rf.Set("resource", "https://foreign.example/mcp")
	if _, err = st.exchangeRemote(rf, form.Get("resource"), now); err == nil || st.Remote.Installs[id] == nil {
		t.Fatal("wrong audience replay revoked authority")
	}
	rf.Set("resource", form.Get("resource"))
	refusal, err := st.exchangeRemote(rf, form.Get("resource"), now)
	if _, ok := refusal.(remoteGrantRefusal); !ok || err != nil {
		t.Fatal("replay refusal must commit family revocation", err)
	}
	if _, _, err = st.remotePrincipal(access, form.Get("resource"), now); err == nil || st.Remote.Installs[id] != nil {
		t.Fatal("old refresh replay left active access/key")
	}
	if _, _, err = st.remotePrincipal(other, otherForm.Get("resource"), now); err != nil {
		t.Fatal("replay revoked other installation", err)
	}
}

func TestRemoteGoogleBindingAndSignedInbox(t *testing.T) {
	st, keys, now := textFixture()
	pending := &remoteAttempt{Exp: now.Add(codeTTL)}
	st.remote().Attempts["attempt"] = pending
	login := st.finishGoogleLogin(&googleAttempt{Remote: "attempt", Exp: now.Add(codeTTL)}, "remote-subject", "remote@example.test", now)
	if login.Problem != "" || pending.Owner == "" {
		t.Fatal("Google did not bind remote request")
	}
	if v := st.finishGoogleLogin(&googleAttempt{Remote: "attempt", Exp: now.Add(codeTTL)}, "foreign-subject", "foreign@example.test", now); v.Problem == "" {
		t.Fatal("foreign Google identity replaced binding")
	}
	_, i, _, _ := remoteFixture(t, st, "agent_b", now)
	st.Pairs[pairID("agent_a", i.Agent)] = &Pair{A: "agent_a", B: i.Agent, State: "active", Generation: 1}
	e := signedText(firstID, "agent_a", i.Agent, "remote-inbox-key-01", "", "private text", now.Add(180*time.Second), keys["agent_a"])
	if _, err := st.ingestPublicText("agent_a", e, now); err != nil {
		t.Fatal(err)
	}
	if _, err := st.remoteTool(i, "knowslink_receive", nil, now); err != nil {
		t.Fatal(err)
	}
	if len(i.Inbox) != 1 || st.Messages[firstID].Receipt.State != "delivered" || st.Messages[firstID].Envelope != nil {
		t.Fatal("missing durable private Inbox before relay ACK")
	}
	raw, _ := json.Marshal(st)
	restarted := newState()
	if err := json.Unmarshal(raw, restarted); err != nil {
		t.Fatal(err)
	}
	var restored *remoteInstall
	for _, candidate := range restarted.Remote.Installs {
		restored = candidate
	}
	v, err := restarted.remoteTool(restored, "knowslink_receive", nil, now)
	if err != nil || !bytes.Contains(mustJSON(v), []byte("private text")) {
		t.Fatal("response loss/restart lost Inbox", err)
	}
	send := json.RawMessage(`{"to":"agent_a","text":"explicit reply","idempotency_key":"remote-reply-key-01","reply_to":"` + firstID + `"}`)
	v, err = restarted.remoteTool(restored, "knowslink_send", send, now)
	if err != nil {
		t.Fatal(err)
	}
	first := v.(map[string]any)["receipt"].(Receipt)
	v, err = restarted.remoteTool(restored, "knowslink_send", send, now)
	if err != nil || v.(map[string]any)["receipt"].(Receipt).ID != first.ID {
		t.Fatal("write idempotency", err)
	}
	if _, err := restarted.remoteTool(restored, "knowslink_ack", json.RawMessage(`{"id":"`+firstID+`"}`), now); err != nil || len(restored.Inbox) != 0 {
		t.Fatal("explicit Inbox ACK", err)
	}
	if _, err := st.remoteTool(i, "knowslink_send", json.RawMessage(`{"to":"agent_c","text":"no relation","idempotency_key":"remote-unpaired-01"}`), now); err == nil {
		t.Fatal("automatic pair authorization")
	}
	i.Inbox[firstID] = json.RawMessage(`{"text":"expired"}`)
	st.sweepRemote(now.Add(180 * time.Second))
	if len(i.Inbox) != 0 {
		t.Fatal("Inbox TTL")
	}
}

func mustJSON(v any) []byte { raw, _ := json.Marshal(v); return raw }

func TestRemoteDeliveryClaimsRetryAndRevocation(t *testing.T) {
	st, keys, now := textFixture()
	id, i, _, _ := remoteFixture(t, st, "agent_b", now)
	st.Pairs[pairID("agent_a", i.Agent)] = &Pair{A: "agent_a", B: i.Agent, State: "active", Generation: 1}
	e := signedText(firstID, "agent_a", i.Agent, "remote-event-key-01", "", "never webhook this", now.Add(180*time.Second), keys["agent_a"])
	if _, err := st.ingestPublicText("agent_a", e, now); err != nil {
		t.Fatal(err)
	}
	st.Remote.Subscriptions["sub"] = &remoteSubscription{Install: id, URL: "https://callback.example/event", Secret: "whsec_" + base64.StdEncoding.EncodeToString(make([]byte, 32)), Exp: now.Add(time.Hour)}
	work := st.claimRemote(now)
	if work.Delivery == nil {
		t.Fatal("pending catch-up")
	}
	job := *work.Delivery
	if work = st.claimRemote(now); work.Delivery != nil || work.Next.IsZero() {
		t.Fatal("double cross-process claim")
	}
	st.finishRemote(job, 503, false, now)
	if st.claimRemote(now).Delivery != nil {
		t.Fatal("retry ignored backoff")
	}
	second := st.claimRemote(now.Add(time.Second))
	if second.Delivery == nil || second.Delivery.Subscription.EventID != job.Subscription.EventID {
		t.Fatal("retry changed event ID")
	}
	st.finishRemote(*second.Delivery, 200, false, now.Add(time.Second))
	if st.claimRemote(now.Add(2*time.Second)).Delivery != nil {
		t.Fatal("duplicate acknowledged event")
	}
	// A crashed sender cannot reset its retry budget after every lease expiry.
	sub := st.Remote.Subscriptions["sub"]
	sub.Seen = ""
	for attempt := 0; attempt < 3; attempt++ {
		if st.claimRemote(now.Add(time.Duration(attempt)*16*time.Second)).Delivery == nil {
			t.Fatal("missing crash-recovery claim")
		}
	}
	if st.claimRemote(now.Add(49*time.Second)).Delivery != nil {
		t.Fatal("crash recovery bypassed three-attempt bound")
	}
	for _, status := range []int{410, 413} {
		sub.Seen = ""
		st.Remote.Subscriptions["sub"] = sub
		work := st.claimRemote(now.Add(time.Minute))
		if work.Delivery == nil {
			t.Fatal("terminal response fixture")
		}
		st.finishRemote(*work.Delivery, status, false, now.Add(time.Minute))
		if st.Remote.Subscriptions["sub"] != nil {
			t.Fatal("terminal delivery response retried")
		}
	}
	st.Agents[i.Agent].Revoked = true
	st.sweepRemote(now)
	if len(st.Remote.Subscriptions) != 0 || st.Remote.Installs[id] != nil {
		t.Fatal("revocation retained secret/key")
	}
}

type remoteRoundTrip func(*http.Request) (*http.Response, error)

func (f remoteRoundTrip) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestRemoteWebhookSignatureAndSSRF(t *testing.T) {
	for _, raw := range []string{"127.0.0.1", "::1", "::ffff:127.0.0.1", "10.0.0.1", "100.64.0.1", "169.254.169.254", "192.0.2.1", "198.18.0.1", "224.0.0.1", "fd00::1", "fe80::1", "64:ff9b::7f00:1", "2001:db8::1", "2002:7f00:1::"} {
		if remotePublicIP(netip.MustParseAddr(raw)) {
			t.Fatal("SSRF IP allowed", raw)
		}
	}
	for _, raw := range []string{"8.8.8.8", "1.1.1.1", "2606:4700:4700::1111"} {
		if !remotePublicIP(netip.MustParseAddr(raw)) {
			t.Fatal("public IP blocked", raw)
		}
	}
	for _, raw := range []string{"http://example.com", "https://user:pass@example.com", "https://example.com/#x", "https://example.com:444/path"} {
		if remoteURL(raw) {
			t.Fatal("unsafe URL allowed", raw)
		}
	}
	client := remoteHTTPClient()
	if client.CheckRedirect(&http.Request{}, nil) == nil {
		t.Fatal("redirect enabled")
	}
	_, _, err := remotePost(context.Background(), client, "event", "subscription", "https://127.0.0.1/", "whsec_"+base64.StdEncoding.EncodeToString(make([]byte, 32)), "", []byte(`{}`))
	if err == nil {
		t.Fatal("localhost callback reached")
	}
	secret := "whsec_" + base64.StdEncoding.EncodeToString(bytes.Repeat([]byte{7}, 32))
	body := []byte(`{"data":{"arrival":true}}`)
	client = &http.Client{Transport: remoteRoundTrip(func(r *http.Request) (*http.Response, error) {
		raw, _ := io.ReadAll(r.Body)
		key, _ := remoteSecret(secret)
		mac := hmac.New(sha256.New, key)
		_, _ = mac.Write([]byte("evt." + r.Header.Get("webhook-timestamp") + "."))
		_, _ = mac.Write(raw)
		want := "v1," + base64.StdEncoding.EncodeToString(mac.Sum(nil))
		if r.Header.Get("webhook-signature") != want || r.Header.Get("X-MCP-Subscription-Id") != "sub" || r.Header.Get("webhook-id") != "evt" || !bytes.Equal(raw, body) {
			t.Fatal("Standard Webhooks signature or exact body")
		}
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`{}`)), Header: make(http.Header)}, nil
	})}
	if status, _, err := remotePost(context.Background(), client, "evt", "sub", "https://callback.example/", secret, "", body); err != nil || status != 200 {
		t.Fatal("signed delivery", err)
	}
	if _, err := remoteSecret("whsec_" + base64.StdEncoding.EncodeToString([]byte("short"))); err == nil {
		t.Fatal("short signing key")
	}
	for _, mode := range []string{"valid", "wrong", "redirect", "large"} {
		client := &http.Client{Transport: remoteRoundTrip(func(r *http.Request) (*http.Response, error) {
			var payload struct{ Type, Challenge string }
			raw, _ := io.ReadAll(r.Body)
			_ = json.Unmarshal(raw, &payload)
			if payload.Type != "verification" || len(payload.Challenge) != 43 || r.Header.Get("webhook-signature") == "" {
				t.Fatal("unsigned or predictable verification challenge")
			}
			status := 200
			if mode == "wrong" {
				payload.Challenge = "wrong"
			}
			if mode == "redirect" {
				status = 302
			}
			answer := string(mustJSON(map[string]string{"challenge": payload.Challenge}))
			if mode == "large" {
				answer = strings.Repeat("x", 4097)
			}
			return &http.Response{StatusCode: status, Body: io.NopCloser(strings.NewReader(answer)), Header: make(http.Header)}, nil
		})}
		err := verifyRemoteCallback(context.Background(), client, "sub", "https://callback.example/", secret)
		if (mode == "valid") != (err == nil) {
			t.Fatal("callback challenge response", mode, err)
		}
	}
}

func TestRemoteAdmissionAndEventVersions(t *testing.T) {
	st, _, now := textFixture()
	id, i, token, form := remoteFixture(t, st, "agent_a", now)
	before := st.eventVersion(now)
	if st.remoteAdmission(token, now) != i.Agent || st.liveCredentials()[hashToken(token)].kind != "remote" {
		t.Fatal("remote admission identity")
	}
	if st.remoteAdmission(token, now.Add(time.Hour)) != "" {
		t.Fatal("admission kept expired token")
	}
	st.remote().Subscriptions["sub"] = &remoteSubscription{Install: id, Secret: "test-secret", Exp: now.Add(time.Hour)}
	st.remote().Subscriptions["sub"].Seen = "bookkeeping"
	if before != st.eventVersion(now) {
		t.Fatal("subscription bookkeeping caused DB notification loop")
	}
	st.Agents[i.Agent].Keys[i.Kid].Revoked = true
	if st.remoteAdmission(token, now) != "" || before == st.eventVersion(now) {
		t.Fatal("current key revocation did not change notification/admission")
	}
	st.sweepRemote(now)
	if _, _, e := st.remotePrincipal(token, form.Get("resource"), now); e == nil {
		t.Fatal("revoked OAuth token survived sweep")
	}
}

func TestRemoteSubscriptionIdentityAndExpiry(t *testing.T) {
	st, _, now := textFixture()
	id, _, _, _ := remoteFixture(t, st, "agent_a", now)
	var input remoteSubscriptionInput
	input.Name = remoteEvent
	input.Delivery.URL = "https://callback.example/events"
	input.Arguments = json.RawMessage(`{}`)
	a := remoteSubscriptionID(id, input)
	b := remoteSubscriptionID("other", input)
	if a == b || a != remoteSubscriptionID(id, input) {
		t.Fatal("subscription identity is not deterministic and installation scoped")
	}
	st.Remote.Subscriptions[a] = &remoteSubscription{Install: id, URL: input.Delivery.URL, Secret: "test-signing-secret", OldSecret: "test-old-secret", Rotated: now, Exp: now.Add(time.Hour)}
	raw, _ := json.Marshal(st)
	restored := newState()
	_ = json.Unmarshal(raw, restored)
	if restored.Remote.Subscriptions[a] == nil {
		t.Fatal("subscription lost across restart")
	}
	restored.sweepRemote(now.Add(time.Minute))
	if restored.Remote.Subscriptions[a].OldSecret != "" {
		t.Fatal("old rotation secret retained")
	}
	restored.sweepRemote(now.Add(time.Hour))
	if restored.Remote.Subscriptions[a] != nil {
		t.Fatal("expired subscription retained")
	}
	for _, raw := range []string{"null", `{"extra":true}`, `[]`} {
		if remoteEmpty(json.RawMessage(raw)) {
			t.Fatal("invalid subscription filter accepted")
		}
	}
	if !remoteEmpty(json.RawMessage(`{ }`)) {
		t.Fatal("empty filter rejected")
	}
}
