//go:build integration

// Real Postgres LISTEN/NOTIFY and HTTP streams prove cross-process delivery, catch-up, isolation and cleanup.
package relay

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestTextEventsHTTP(t *testing.T) {
	pool := messagePool(t)
	f := setup(t, pool)
	f.mutate(func(st *State) {
		for id, owner := range f.owners {
			st.Members[id] = &Member{Owner: owner, Active: true}
		}
	})
	other := &Service{Pool: pool}
	server := httptest.NewUnstartedServer(other.Handler())
	server.Config.WriteTimeout = 10 * time.Second
	server.Start()
	defer server.Close()
	open := func(token string) (*http.Response, context.CancelFunc) {
		t.Helper()
		ctx, cancel := context.WithCancel(context.Background())
		req, _ := http.NewRequestWithContext(ctx, "GET", server.URL+"/v1/text/events", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		res, err := (&http.Client{Timeout: 25 * time.Second}).Do(req)
		if err != nil {
			cancel()
			t.Fatal(err)
		}
		t.Cleanup(func() { cancel(); _ = res.Body.Close() })
		return res, cancel
	}
	bad, badCancel := open("wrong")
	if bad.StatusCode != 401 {
		t.Fatal("wrong key admitted")
	}
	badCancel()
	res, cancel := open(f.tokens["agent_b"])
	if res.StatusCode != 200 || res.Header.Get("Content-Type") != "text/event-stream" {
		t.Fatal("stream admission")
	}
	frames := make(chan string, 16)
	go func() {
		defer close(frames)
		scan := bufio.NewScanner(res.Body)
		frame := ""
		for scan.Scan() {
			if scan.Text() == "" {
				frames <- frame
				frame = ""
			} else {
				frame += scan.Text() + "\n"
			}
		}
	}()
	next := func(want string) {
		t.Helper()
		select {
		case frame := <-frames:
			if frame != want {
				t.Fatalf("unexpected metadata frame %q", frame)
			}
		case <-time.After(20 * time.Second):
			t.Fatal("missing SSE frame")
		}
	}
	next("event: ready\ndata: {}\n")
	// Heartbeat must survive the server's normal 10s write timeout and must not touch the singleton row.
	var before, after int64
	if err := pool.QueryRow(context.Background(), "SELECT epoch FROM relay_state").Scan(&before); err != nil {
		t.Fatal(err)
	}
	next(": heartbeat\n")
	if err := pool.QueryRow(context.Background(), "SELECT epoch FROM relay_state").Scan(&after); err != nil || before != after {
		t.Fatal("idle stream performed DB transactions", err)
	}
	e := signedText(firstID, "agent_a", "agent_b", "events-http-send-01", "", "private body never in SSE", time.Now().Add(time.Minute), f.private["agent_a"])
	raw, _ := json.Marshal(e)
	if w := f.request("POST", "/v1/text/send", f.tokens["agent_a"], raw, ""); w.Code != 200 {
		t.Fatal("text send", w.Code)
	}
	next("event: arrival\ndata: {}\n")
	// Another agent receives ready only, even while B has pending text.
	foreign, foreignCancel := open(f.tokens["agent_c"])
	reader := bufio.NewReader(foreign.Body)
	for _, want := range []string{"event: ready\n", "data: {}\n", "\n"} {
		if line, err := reader.ReadString('\n'); err != nil || line != want {
			t.Fatal("foreign ready", err)
		}
	}
	foreignCancel()
	cancel()
	// Reconnect catches pending data without relying on notification history.
	reconnected, reconnectCancel := open(f.tokens["agent_b"])
	data := make([]byte, len("event: ready\ndata: {}\n\nevent: arrival\ndata: {}\n\n"))
	if _, err := io.ReadFull(reconnected.Body, data); err != nil || string(data) != "event: ready\ndata: {}\n\nevent: arrival\ndata: {}\n\n" {
		t.Fatal("pending catch-up", err)
	}
	// Existing API still leases/persists/ACKs the original signed text.
	lease := f.call("POST", "/v1/text/pull", f.tokens["agent_b"], map[string]any{}, 200)
	body := map[string]any{"id": firstID, "token": lease["lease_token"]}
	for _, action := range []string{"persist", "ack"} {
		f.call("POST", "/v1/text/"+action, f.tokens["agent_b"], body, 200)
	}
	f.call("POST", "/v1/key-revoke", f.tokens["agent_b_owner"], map[string]any{"agent": "agent_b", "kid": "key1"}, 200)
	remaining, err := io.ReadAll(reconnected.Body)
	if err != nil || strings.Contains(string(remaining), "arrival") {
		t.Fatal("revoked stream did not close cleanly", err)
	}
	reconnectCancel()
	denied, deniedCancel := open(f.tokens["agent_b"])
	if denied.StatusCode != 401 {
		t.Fatal("revoked reconnect admitted")
	}
	deniedCancel()
	deadline := time.Now().Add(3 * time.Second)
	for {
		count := -1
		f.mutate(func(st *State) { count = len(st.Streams) })
		if count == 0 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("stream admissions leaked")
		}
		time.Sleep(10 * time.Millisecond)
	}
}
