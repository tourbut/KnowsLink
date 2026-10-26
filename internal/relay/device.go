// Client-started connections bind Google consent to a public key; only the originating private key can collect credentials.
package relay

import (
	"crypto/ed25519"
	"encoding/base64"
	"net/http"
	"time"
)

func deviceBytes(token string, c *Connection) []byte {
	return []byte("KNOWSLINK-DEVICE\x00" + token + "\x00" + c.Client + "\x00" + c.Kid + "\x00" + c.Public)
}

func deviceProof(token, proof string, c *Connection) bool {
	key, e := base64.RawURLEncoding.DecodeString(c.Public)
	sig, e2 := base64.RawURLEncoding.DecodeString(proof)
	return e == nil && e2 == nil && len(key) == ed25519.PublicKeySize && ed25519.Verify(key, deviceBytes(token, c), sig)
}

// startDevice keeps a known host only as the default-name hint (A-04); any other value is simply not used.
func (st *State) startDevice(token, client, kid, public, proof, host string, now time.Time) (*Connection, error) {
	decoded, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil || len(decoded) != 32 || client != supportedClient || !kidPattern.MatchString(kid) || len(kid) > 128 {
		return nil, fault("invalid_schema")
	}
	c := &Connection{Device: true, Client: client, Mode: "register", Kid: kid, Public: public, State: "requested", Exp: now.Add(connectionTTL)}
	if hostNames[host] != "" {
		c.Host = host
	}
	if !deviceProof(token, proof, c) {
		return nil, fault("invalid_signature")
	}
	if st.Connections[hashToken(token)] != nil {
		return nil, fault("invalid_auth")
	}
	// Keep one-use tombstones for the existing 24h window while bounding anonymous storage.
	devices := 0
	for _, existing := range st.Connections {
		if existing.Device {
			devices++
		}
	}
	if devices >= 2000 {
		return nil, fault("capacity")
	}
	st.Connections[hashToken(token)] = c
	return c, nil
}

// approveDevice records the owner's consent and reserves the agent ID; the agent itself is made only when the client
// saves its key (A-01). A rate or capacity refusal ends the attempt so the polling client can report it.
func (st *State) approveDevice(id, owner, name string, now time.Time) error {
	c := st.Connections[id]
	if c == nil || !c.Device || c.Owner != owner || c.State != "prepared" || !now.Before(c.Exp) || st.Owners[owner] == nil || !st.Owners[owner].Active {
		return fault("invalid_auth")
	}
	if err := st.agentLimit(owner, now); err != nil {
		r := err.(agentRefused)
		c.State, c.Refused, c.Retry = "refused", r.Code, r.Retry
		return err
	}
	c.Agent, c.Name, c.State = newMemberAgentID(), name, "approved"
	return nil
}

func (s *Service) deviceAPI(start bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct{ Token, Client, Kid, Public, Proof, Host string }
		parseErr := decodeCommand(w, r, &input)
		value, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
			if ok, retry := s.requestHit(st, r, now, anonymousRate(s.clientIP(r))...); !ok {
				return retry, nil
			}
			var c *Connection
			e := parseErr
			if e == nil && start {
				c, e = st.startDevice(input.Token, input.Client, input.Kid, input.Public, input.Proof, input.Host, now)
			} else if e == nil {
				c = st.Connections[hashToken(input.Token)]
				switch {
				case c == nil || !c.Device || c.Client != input.Client || !deviceProof(input.Token, input.Proof, c):
					e = fault("invalid_auth")
				case c.State == "refused":
					e = agentRefused{c.Refused, c.Retry}
				case !c.open(now):
					e = fault("invalid_auth")
				}
			}
			if e != nil {
				return errorBody(e), nil
			}
			v := map[string]string{"state": c.State, "exp": c.Exp.Format(time.RFC3339)}
			if c.State == "approved" {
				if _, e = st.connection(input.Token, now); e != nil {
					return errorBody(e), nil
				}
				// The owner chose the name in the browser; the installer shows it beside the fingerprint (AS-05, F-07).
				v["owner"], v["agent"], v["name"] = c.Owner, c.Agent, c.Name
			}
			return v, nil
		})
		connectionResponse(w, value, err)
	}
}

// errorBody is a JSON refusal; an agent_rate refusal also names when a new agent fits again.
func errorBody(e error) map[string]string {
	v := map[string]string{"error": e.Error()}
	if r, ok := e.(agentRefused); ok && !r.Retry.IsZero() {
		v["retry_at"] = r.Retry.UTC().Add(time.Second - 1).Truncate(time.Second).Format(time.RFC3339)
	}
	return v
}

func connectionResponse(w http.ResponseWriter, value any, err error) {
	if retry, ok := value.(time.Time); ok {
		rateLimited(w, retry)
		return
	}
	if v, ok := value.(map[string]string); ok && v["error"] != "" {
		writeJSON(w, statusFor(fault(v["error"])), v)
		return
	}
	respond(w, value, err)
}

// deviceOutcome is what the member who started an attempt sees on a revisit once it is approved or ended (W-01).
func (st *State) deviceOutcome(c *Connection) map[string]any {
	titles := map[string]string{"consumed": "이 연결은 완료됐습니다", "approved": "연결을 승인했습니다", "cancelled": "이 연결은 취소됐습니다", "expired": "이 연결은 만료됐습니다"}
	if titles[c.State] == "" {
		return nil
	}
	name := c.Name
	if a := st.Agents[c.Agent]; a != nil {
		name = agentName(a) // the owner may have renamed it since
	}
	return map[string]any{"Title": titles[c.State], "Outcome": c.State, "Name": name}
}

func (s *Service) devicePage(w http.ResponseWriter, r *http.Request) {
	value, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
		id, session, retry, e := s.memberHit(w, st, r, now, memberRate)
		if !retry.IsZero() {
			return refusal{429, limited(retry), false, false}, nil
		}
		c := st.Connections[r.PathValue("id")]
		// Only the member who started the attempt learns its state or agent name; anyone else keeps the generic refusal.
		if c != nil && c.Device && c.Owner != "" && e == nil && st.Members[id].Owner == c.Owner {
			if c.State == "refused" {
				refusedErr := agentRefused{c.Refused, c.Retry}
				return refusal{statusFor(refusedErr), safeProblem(refusedErr), false, false}, nil
			}
			if v := st.deviceOutcome(c); v != nil {
				return v, nil
			}
		}
		if c == nil || !c.Device || !c.open(now) || !now.Before(c.Exp) {
			return refusal{401, safeProblem(fault("invalid_auth")), false, false}, nil
		}
		if c.Owner != "" && (e != nil || st.Members[id].Owner != c.Owner) {
			return refusal{403, "이 연결을 시작한 Google 계정으로 다시 시작하세요.", false, false}, nil
		}
		public, _ := base64.RawURLEncoding.DecodeString(c.Public)
		v := map[string]any{"Title": "내 클라이언트 연결", "ConnectionID": r.PathValue("id"), "Fingerprint": fingerprint(public), "Exp": clock(c.Exp), "State": c.State, "Google": s.Google != nil, "Here": "/connect/" + r.PathValue("id")}
		v["Login"] = c.Owner == "" || e != nil || now.Sub(session.Verified) >= reauthWindow
		if c.Owner != "" {
			v["Name"] = st.defaultName(c.Owner, c.Host, responseLang(w))
		}
		return v, nil
	})
	if err != nil {
		refused(w, "클라이언트 연결", refusal{503, "잠시 뒤 다시 시도하세요.", false, false})
		return
	}
	if v, ok := value.(refusal); ok {
		refused(w, "클라이언트 연결", v)
		return
	}
	render(w, 200, "device", value.(map[string]any))
}
