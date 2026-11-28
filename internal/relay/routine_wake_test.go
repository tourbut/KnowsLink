// Routine wake checks cover URL and SSRF refusal, rate limits, result handling, sealing and the metadata-only body.
package relay

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"
	"time"
)

const routineTestURL = "https://api2.cursor.sh/automations/webhook/abc_DEF-123"
const routineTestKey = "routine-secret-key-0123456789"
const routineThirdID = "0199a3f2-4c10-7a11-8b22-3344556677bb"

func routineFixture(t *testing.T) (*State, *Sealer, string, func(id string, at time.Time), time.Time) {
	t.Helper()
	st, keys, now := textFixture()
	_, install, _, _ := remoteFixture(t, st, "agent_b", now)
	st.Pairs[pairID("agent_a", install.Agent)] = &Pair{A: "agent_a", B: install.Agent, State: "active", Generation: 1}
	sealer, err := NewSealer(base64.StdEncoding.EncodeToString(make([]byte, 32)))
	if err != nil {
		t.Fatal(err)
	}
	arrive := func(id string, at time.Time) {
		t.Helper()
		text := signedText(id, "agent_a", install.Agent, "routine-"+id, "", "secret body from sender", at.Add(180*time.Second), keys["agent_a"])
		if _, err := st.ingestPublicText("agent_a", text, at); err != nil {
			t.Fatal(err)
		}
	}
	return st, sealer, install.Agent, arrive, now
}

func TestRoutineRegistrationRefusesBadInputAndSealsKey(t *testing.T) {
	st, sealer, agent, _, now := routineFixture(t)
	for _, bad := range []string{
		"http://api2.cursor.sh/automations/webhook/abc",
		"https://api2.cursor.sh.evil.example/automations/webhook/abc",
		"https://evil.example/automations/webhook/abc",
		"https://api2.cursor.sh/automations/webhook/",
		"https://api2.cursor.sh/automations/webhook/abc/../x",
		"https://api2.cursor.sh/automations/webhook/abc?x=1",
		"https://user@api2.cursor.sh/automations/webhook/abc",
		"https://api2.cursor.sh:8443/automations/webhook/abc",
		"https://api2.cursor.sh/automations/webhook/" + strings.Repeat("a", 129),
	} {
		if err := st.setRoutine(sealer, "agent_b", agent, bad, routineTestKey, now); err == nil || err.Error() != "routine_url" {
			t.Fatalf("accepted %q: %v", bad, err)
		}
	}
	for _, bad := range []string{"", "has space", "line\nbreak", strings.Repeat("k", 513)} {
		if err := st.setRoutine(sealer, "agent_b", agent, routineTestURL, bad, now); err == nil || err.Error() != "routine_key" {
			t.Fatalf("accepted key %q: %v", bad, err)
		}
	}
	if err := st.setRoutine(nil, "agent_b", agent, routineTestURL, routineTestKey, now); err == nil || err.Error() != "routine_unavailable" {
		t.Fatal("registered without a seal key", err)
	}
	if err := st.setRoutine(sealer, "agent_a", agent, routineTestURL, routineTestKey, now); err == nil || err.Error() != "sender_not_allowed" {
		t.Fatal("another owner registered a routine", err)
	}
	if err := st.setRoutine(sealer, "agent_b", "agent_b", routineTestURL, routineTestKey, now); err == nil || err.Error() != "sender_not_allowed" {
		t.Fatal("a local agent got a routine", err)
	}
	if len(st.Remote.Routines) != 0 {
		t.Fatal("refused registration stored a record")
	}
	if err := st.setRoutine(sealer, "agent_b", agent, " "+routineTestURL+" ", " "+routineTestKey+" ", now); err != nil {
		t.Fatal(err)
	}
	w := st.Remote.Routines[agent]
	if plain, err := sealer.Open(w.Key, agent); err != nil || plain != routineTestKey || w.URL != routineTestURL {
		t.Fatal("sealed key did not round-trip", plain, err)
	}
	if _, err := sealer.Open(w.Key, "other_agent"); err == nil {
		t.Fatal("sealed key opened under another agent")
	}
	raw, _ := json.Marshal(st)
	views := st.routineViews(st.humanAgents("agent_b"), now, "ko")
	shown, _ := json.Marshal(views)
	if strings.Contains(string(raw), routineTestKey) || strings.Contains(string(shown), routineTestKey) || strings.Contains(string(shown), "abc_DEF") {
		t.Fatal("key or full URL leaked into state or view")
	}
	if len(views) != 1 || views[0].Tail != "-123" || !views[0].Registered || views[0].Broken {
		t.Fatalf("view: %+v", views)
	}
	if err := st.removeRoutine("agent_a", agent); err == nil {
		t.Fatal("another owner removed the routine")
	}
	if err := st.removeRoutine("agent_b", agent); err != nil || st.Remote.Routines[agent] != nil {
		t.Fatal("owner could not remove the routine", err)
	}
}

func TestRoutineSpacingHourlyCapAndNoRepeat(t *testing.T) {
	st, sealer, agent, arrive, now := routineFixture(t)
	if err := st.setRoutine(sealer, "agent_b", agent, routineTestURL, routineTestKey, now); err != nil {
		t.Fatal(err)
	}
	if job, _ := st.claimRoutine(now); job != nil {
		t.Fatal("empty queue rang")
	}
	arrive(firstID, now)
	job, _ := st.claimRoutine(now)
	if job == nil || job.Pending != 1 {
		t.Fatal("first arrival did not ring", job)
	}
	st.finishRoutine(*job, 200, now)
	if again, _ := st.claimRoutine(now.Add(time.Minute)); again != nil {
		t.Fatal("the same waiting set rang twice")
	}
	later := now.Add(30 * time.Second)
	arrive(secondID, later)
	job, next := st.claimRoutine(later)
	if job != nil || !next.Equal(now.Add(routineSpacing)) {
		t.Fatalf("a new arrival inside 2 minutes rang: next=%v", next)
	}
	job, _ = st.claimRoutine(now.Add(routineSpacing))
	if job == nil || job.Pending != 2 {
		t.Fatal("the coalesced arrival did not ring after 2 minutes", job)
	}
	st.finishRoutine(*job, 200, now.Add(routineSpacing))
	// A set that only shrank (a message was received) does not ring again.
	st.Messages[secondID].Receipt.State = "delivered"
	if again, _ := st.claimRoutine(now.Add(10 * time.Minute)); again != nil {
		t.Fatal("a shrinking set rang")
	}

	w := st.Remote.Routines[agent]
	w.Window, w.Count = now, routineHourly
	arrive(routineThirdID, now.Add(10*time.Minute))
	job, next = st.claimRoutine(now.Add(10 * time.Minute))
	if job != nil || !next.Equal(now.Add(time.Hour)) {
		t.Fatalf("the 21st POST in an hour was sent: next=%v", next)
	}
	arrive(gateID, now.Add(time.Hour))
	job, _ = st.claimRoutine(now.Add(time.Hour))
	if job == nil || w.Count != 1 || !w.Window.Equal(now.Add(time.Hour)) {
		t.Fatal("the next hourly window did not reopen", w.Count)
	}
}

func TestRoutineResultHandling(t *testing.T) {
	for _, c := range []struct {
		name    string
		status  []int
		posts   int
		broken  bool
		pausedH bool
	}{
		{"started", []int{200}, 1, false, false},
		{"5xx retries once", []int{503, 503}, 2, false, false},
		{"connection failure retries once", []int{0, 200}, 2, false, false},
		{"unknown answer is not retried", []int{routineUnknown}, 1, false, false},
		{"other 4xx is not retried", []int{400}, 1, false, false},
		{"429 waits for the hourly window", []int{429}, 1, false, true},
		{"401", []int{401}, 1, true, false},
		{"403", []int{403}, 1, true, false},
		{"404", []int{404}, 1, true, false},
		{"410", []int{410}, 1, true, false},
		{"lost seal key", []int{routineSealLost}, 1, true, false},
	} {
		t.Run(c.name, func(t *testing.T) {
			st, sealer, agent, arrive, now := routineFixture(t)
			if err := st.setRoutine(sealer, "agent_b", agent, routineTestURL, routineTestKey, now); err != nil {
				t.Fatal(err)
			}
			arrive(firstID, now)
			at, posts := now, 0
			for _, status := range c.status {
				job, next := st.claimRoutine(at)
				if job == nil {
					t.Fatalf("post %d not claimed (next %v)", posts+1, next)
				}
				posts++
				st.finishRoutine(*job, status, at)
				at = at.Add(routineRetry)
			}
			w := st.Remote.Routines[agent]
			if job, _ := st.claimRoutine(at.Add(routineSpacing)); job != nil {
				t.Fatal("an extra POST was claimed")
			}
			if posts != c.posts || w.Broken != c.broken {
				t.Fatalf("posts=%d broken=%v", posts, w.Broken)
			}
			if c.status[len(c.status)-1] == 200 && w.LastOK.IsZero() {
				t.Fatal("success time missing")
			}
			if c.pausedH {
				arrive(secondID, now.Add(time.Hour))
				if job, _ := st.claimRoutine(now.Add(time.Hour)); job == nil {
					t.Fatal("429 did not resume in the next window")
				}
			}
			if c.broken {
				arrive(secondID, now.Add(time.Hour))
				if job, _ := st.claimRoutine(now.Add(time.Hour)); job != nil {
					t.Fatal("a broken registration rang")
				}
				if err := st.setRoutine(sealer, "agent_b", agent, routineTestURL, routineTestKey, now.Add(time.Hour)); err != nil {
					t.Fatal(err)
				}
				if job, _ := st.claimRoutine(now.Add(time.Hour)); job == nil {
					t.Fatal("re-registration did not resume")
				}
			}
		})
	}
}

func TestRoutinePostHeadersBodyAndOutcomes(t *testing.T) {
	var got *http.Request
	var body []byte
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/slow" {
			time.Sleep(300 * time.Millisecond)
		}
		got, body = r, must(io.ReadAll(r.Body))
		w.Write([]byte(`{"runUuid":"run"}`))
	}))
	defer server.Close()
	client := server.Client()
	if status := routinePost(context.Background(), client, server.URL+"/hook", routineTestKey, 2); status != 200 {
		t.Fatal(status)
	}
	if got.Method != "POST" || got.Header.Get("Authorization") != "Bearer "+routineTestKey || got.Header.Get("X-Automation-Key") != routineTestKey || got.Header.Get("Content-Type") != "application/json" {
		t.Fatal("headers", got.Header)
	}
	var fields map[string]any
	if err := json.Unmarshal(body, &fields); err != nil {
		t.Fatal(err)
	}
	keys := []string{}
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	if strings.Join(keys, ",") != "event,next,pending" || fields["event"] != "knowslink_text_arrived" || fields["pending"] != 2.0 || fields["next"] != routineNext {
		t.Fatalf("body %s", body)
	}
	if strings.Contains(string(body), "secret body") || strings.Contains(string(body), "agent_a") {
		t.Fatal("body leaked message text or sender")
	}
	client.Timeout = 100 * time.Millisecond
	if status := routinePost(context.Background(), client, server.URL+"/slow", routineTestKey, 1); status != routineUnknown {
		t.Fatal("a sent request without an answer must not count as a connection failure", status)
	}
	closed := server.URL
	server.Close()
	if status := routinePost(context.Background(), client, closed+"/hook", routineTestKey, 1); status != 0 {
		t.Fatal("refused connection", status)
	}
	// The production client refuses loopback, private and non-443 targets before anything is sent.
	for _, target := range []string{"https://127.0.0.1/automations/webhook/x", "https://10.0.0.1/automations/webhook/x", "https://[::1]/automations/webhook/x", "https://169.254.169.254/automations/webhook/x"} {
		if status := routinePost(context.Background(), remoteHTTPClient(), target, routineTestKey, 1); status != 0 {
			t.Fatal("SSRF target reached", target, status)
		}
	}
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
