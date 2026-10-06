//go:build integration

// HTTP admission under slow bodies and rejected cleanup requests, against the shared Postgres state.
package relay

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// Review H1/M1: body reception holds no admission slot, and only a verified principal ending its own record spends cleanup.
func TestHTTPSlowBodyAndCleanupAdmission(t *testing.T) {
	pool := messagePool(t)
	f := setup(t, pool)
	f.addAgent("agent_d")
	f.addAgent("agent_e")
	f.mutate(func(st *State) { st.Rates = map[string][]time.Time{} })
	release := make(chan struct{})
	var once sync.Once
	unblock := func() { once.Do(func() { close(release) }) }
	defer unblock()
	var block atomic.Bool
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if block.Load() {
			<-release
		}
		w.WriteHeader(200)
	})
	server := httptest.NewServer(f.s.boundedHTTP(next))
	defer server.Close()
	revoke := `{"agent":"agent_a","kid":"key1"}`
	// Slow senders announce 100 bytes and send 1; before the fix each held a slot until the 10s read deadline.
	var conns []net.Conn
	slow := func(path, auth string, n int) {
		for range n {
			c, err := net.Dial("tcp", strings.TrimPrefix(server.URL, "http://"))
			if err != nil {
				t.Fatal(err)
			}
			conns = append(conns, c)
			fmt.Fprintf(c, "POST %s HTTP/1.1\r\nHost: x\r\nContent-Type: application/json\r\n%sContent-Length: 100\r\n\r\n{", path, auth)
		}
	}
	probe := func(path, token, body string) int {
		r, _ := http.NewRequest("POST", server.URL+path, strings.NewReader(body))
		r.Header.Set("Authorization", "Bearer "+token)
		resp, err := http.DefaultClient.Do(r)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		return resp.StatusCode
	}
	slow("/v1/ack", "", 8)
	slow("/v1/key-revoke", "Authorization: Bearer "+f.tokens["agent_a_owner"]+"\r\n", 8)
	slow("/auth/start", "", 24)
	time.Sleep(300 * time.Millisecond)
	if code := probe("/v1/key-revoke", f.tokens["agent_a_owner"], revoke); code != 200 {
		t.Fatalf("cleanup behind slow bodies got %d", code)
	}
	if code := probe("/v1/registry", f.tokens["agent_a"], "{}"); code != 200 {
		t.Fatalf("new work behind slow bodies got %d", code)
	}
	f.mutate(func(st *State) {
		if len(st.Rates["http:ip:127.0.0.1"]) != 0 || len(st.Rates["cleanup:member:"+f.owners["agent_a"]]) != 1 {
			t.Fatalf("slow bodies spent rate before admission: %v", st.Rates)
		}
	})
	inFlight := func(want int) (clean int) {
		t.Helper()
		deadline := time.Now().Add(5 * time.Second)
		for {
			n := 0
			clean = 0
			f.mutate(func(st *State) {
				n = len(st.HTTP)
				for _, a := range st.HTTP {
					if a.Clean {
						clean++
					}
				}
			})
			if n == want {
				return clean
			}
			if time.Now().After(deadline) {
				t.Fatalf("admitted %d want %d", n, want)
			}
			time.Sleep(20 * time.Millisecond)
		}
	}
	// A slow body that finally arrives malformed still reaches admission and spends its budget once.
	limited := 0
	for _, c := range conns {
		fmt.Fprint(c, strings.Repeat("x", 99))
		resp, err := http.ReadResponse(bufio.NewReader(c), nil)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		c.Close()
		if resp.StatusCode == 429 {
			limited++
		}
	}
	inFlight(0)
	f.mutate(func(st *State) {
		if limited != 2 || len(st.Rates["http:ip:127.0.0.1"]) != 31 {
			t.Fatalf("32 anonymous slow bodies: %d limited, ip hits %d", limited, len(st.Rates["http:ip:127.0.0.1"]))
		}
		st.Rates = map[string][]time.Time{}
	})
	// Requests that cannot end their own record hold the new local channel and shared new records only.
	block.Store(true)
	done := make(chan int, 32)
	send := func(path, auth, body string, header ...string) {
		go func() {
			r := httptest.NewRequest("POST", path, strings.NewReader(body))
			if auth != "" {
				r.Header.Set("Authorization", "Bearer "+auth)
			}
			for i := 0; i+1 < len(header); i += 2 {
				r.Header.Set(header[i], header[i+1])
			}
			w := httptest.NewRecorder()
			f.s.boundedHTTP(next).ServeHTTP(w, r)
			done <- w.Code
		}()
	}
	send("/v1/ack", "", `{"id":"x","token":"y"}`)                                             // anonymous
	send("/v1/key-revoke", "invented-token", revoke)                                          // wrong credential
	send("/v1/key-revoke", f.tokens["agent_b_owner"], revoke)                                 // another owner's key
	send("/v1/key-revoke", f.tokens["agent_a_owner"], `{"agent":`)                            // malformed body
	send("/v1/key-revoke", f.tokens["agent_a_owner"], revoke, "Sec-Fetch-Site", "cross-site") // cross-origin
	send("/v1/ack", f.tokens["agent_b"], `{"id":"none","token":"none"}`)                      // no lease
	send("/owner/gates/none/deny", f.tokens["agent_b_owner"], "csrf=bad&decision=deny", "Content-Type", "application/x-www-form-urlencoded")
	if clean := inFlight(7); clean != 0 {
		t.Fatalf("%d rejected cleanup requests took cleanup admission", clean)
	}
	own := func(id string) string { return `{"agent":"` + id + `","kid":"key1"}` }
	for _, id := range []string{"agent_a", "agent_b", "agent_c", "agent_d"} {
		send("/v1/key-revoke", f.tokens[id+"_owner"], own(id))
	}
	if clean := inFlight(11); clean != 4 {
		t.Fatalf("own cleanup admitted %d of 4", clean)
	}
	// A fifth owner meets the full local cleanup channel and the first owner's repeat meets its own held cleanup slot.
	// Both are refused before DB admission, spending no rate.
	for _, id := range []string{"agent_e", "agent_a"} {
		r := httptest.NewRequest("POST", "/v1/key-revoke", strings.NewReader(own(id)))
		r.Header.Set("Authorization", "Bearer "+f.tokens[id+"_owner"])
		w := httptest.NewRecorder()
		f.s.boundedHTTP(next).ServeHTTP(w, r)
		if w.Code != 429 || !strings.Contains(w.Body.String(), "capacity") {
			t.Fatalf("%s extra cleanup got %d", id, w.Code)
		}
	}
	f.mutate(func(st *State) {
		if len(st.Rates["http:member:"+f.owners["agent_b"]]) != 2 || len(st.Rates["cleanup:member:"+f.owners["agent_e"]]) != 0 {
			t.Fatalf("cleanup budget spent by rejected requests: %v", st.Rates)
		}
		for _, id := range []string{"agent_a", "agent_b", "agent_c", "agent_d"} {
			if len(st.Rates["cleanup:member:"+f.owners[id]]) != 1 {
				t.Fatalf("own cleanup budget %s: %v", id, st.Rates)
			}
		}
	})
	unblock()
	for range 11 {
		if code := <-done; code != 200 {
			t.Fatalf("admitted request got %d", code)
		}
	}
	inFlight(0)
}
