// Test transport exposes only paired allowlisted agents; owner and business operations stay outside the machine path.
package relay

import (
	"io"
	"net/http"
	"strings"
	"time"
)

func TestAgentAllowlist(value string) (map[string]bool, error) {
	agents := map[string]bool{}
	if value == "" {
		return agents, nil
	}
	parts := strings.Split(value, ",")
	if len(parts) != 2 {
		return nil, fault("KNOWSLINK_TEST_AGENTS requires two distinct agent IDs")
	}
	for _, agent := range parts {
		if !agentPattern.MatchString(agent) || agents[agent] {
			return nil, fault("KNOWSLINK_TEST_AGENTS requires two distinct agent IDs")
		}
		agents[agent] = true
	}
	return agents, nil
}

func (s *Service) testHandler(normal http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := s.transaction(r.Context(), func(st *State, _ time.Time) (any, error) {
			agent, err := st.principal(bearer(r), "agent")
			if err != nil {
				return nil, err
			}
			if !s.TestAgents[agent] {
				return nil, fault("sender_not_allowed")
			}
			return nil, nil
		})
		if err != nil {
			respond(w, nil, err)
			return
		}
		suffix := strings.TrimPrefix(r.URL.Path, "/v1/test")
		allowed := r.Method == "GET" && (suffix == "/registry" || strings.HasPrefix(suffix, "/keys/") || strings.HasPrefix(suffix, "/receipts/"))
		if r.Method == "POST" {
			switch suffix {
			case "/send", "/pull", "/persist", "/ack", "/claim":
				allowed = true
			}
		}
		if !allowed {
			respond(w, nil, fault("sender_not_allowed"))
			return
		}
		if suffix == "/send" {
			s.testSend(w, r)
			return
		}
		if r.Method == "POST" {
			s.operation("test-"+strings.TrimPrefix(suffix, "/"))(w, r)
			return
		}
		if strings.HasPrefix(suffix, "/keys/") {
			parts := strings.Split(suffix, "/")
			if len(parts) != 4 || !s.TestAgents[parts[2]] {
				respond(w, nil, fault("sender_not_allowed"))
				return
			}
		}
		if strings.HasPrefix(suffix, "/receipts/") {
			_, err := s.transaction(r.Context(), func(st *State, _ time.Time) (any, error) {
				m := st.Messages[strings.TrimPrefix(suffix, "/receipts/")]
				if m == nil || m.Receipt.Intent != "relay.test.message" {
					return nil, fault("sender_not_allowed")
				}
				return nil, nil
			})
			if err != nil {
				respond(w, nil, err)
				return
			}
		}
		// Work on a request copy. The outer machine path cannot reach owner/signup/rotation/pairing.
		copy := r.Clone(r.Context())
		copy.URL.Path = "/v1" + suffix
		copy.URL.RawPath = ""
		normal.ServeHTTP(w, copy)
	})
}

func (s *Service) testSend(w http.ResponseWriter, r *http.Request) {
	// Parse signed fields once with the same strict parser as the normal transport.
	raw, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 64*1024))
	if err != nil {
		respond(w, nil, fault("invalid_json"))
		return
	}
	parsed, err := Parse(raw)
	if err != nil {
		respond(w, nil, err)
		return
	}
	if parsed.Intent != "relay.test.message" {
		respond(w, nil, fault("sender_not_allowed"))
		return
	}
	value, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
		agent, err := st.principal(bearer(r), "agent")
		if err != nil {
			return nil, err
		}
		return st.ingest(agent, parsed, "", now)
	})
	respond(w, value, err)
}
