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

func (st *State) startDevice(token, client, kid, public, proof string, now time.Time) (*Connection, error) {
	decoded, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil || len(decoded) != 32 || client != supportedClient || !kidPattern.MatchString(kid) || len(kid) > 128 {
		return nil, fault("invalid_schema")
	}
	c := &Connection{Device: true, Client: client, Mode: "register", Kid: kid, Public: public, State: "requested", Exp: now.Add(connectionTTL)}
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

func (st *State) approveDevice(id, owner string, now time.Time) error {
	c := st.Connections[id]
	if c == nil || !c.Device || c.Owner != owner || c.State != "prepared" || !now.Before(c.Exp) || st.Owners[owner] == nil || !st.Owners[owner].Active {
		return fault("invalid_auth")
	}
	if !st.agentCapacity(owner) {
		return fault("capacity")
	}
	agent := newMemberAgentID()
	st.Agents[agent] = &Agent{Owner: owner, Keys: map[string]*Key{}}
	c.Agent, c.State = agent, "approved"
	return nil
}

func (s *Service) deviceAPI(start bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct{ Token, Client, Kid, Public, Proof string }
		parseErr := decodeCommand(w, r, &input)
		value, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
			if ok, retry := s.requestHit(st, r, now, anonymousRate(s.clientIP(r))...); !ok {
				return retry, nil
			}
			var c *Connection
			e := parseErr
			if e == nil && start {
				c, e = st.startDevice(input.Token, input.Client, input.Kid, input.Public, input.Proof, now)
			} else if e == nil {
				c = st.Connections[hashToken(input.Token)]
				if c == nil || !c.Device || c.Client != input.Client || !now.Before(c.Exp) || c.State == "cancelled" || c.State == "consumed" || !deviceProof(input.Token, input.Proof, c) {
					e = fault("invalid_auth")
				}
			}
			if e != nil {
				return map[string]string{"error": e.Error()}, nil
			}
			v := map[string]string{"state": c.State, "exp": c.Exp.Format(time.RFC3339)}
			if c.State == "approved" {
				if _, e = st.connection(input.Token, now); e != nil {
					return map[string]string{"error": e.Error()}, nil
				}
				v["owner"], v["agent"] = c.Owner, c.Agent
			}
			return v, nil
		})
		if retry, ok := value.(time.Time); ok {
			rateLimited(w, retry)
			return
		}
		if v, ok := value.(map[string]string); ok && v["error"] != "" {
			respond(w, nil, fault(v["error"]))
			return
		}
		respond(w, value, err)
	}
}

func (s *Service) devicePage(w http.ResponseWriter, r *http.Request) {
	value, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
		id, session, retry, e := s.memberHit(st, r, now, memberRate)
		if !retry.IsZero() {
			return refusal{429, limited(retry), false, false}, nil
		}
		c := st.Connections[r.PathValue("id")]
		if c == nil || !c.Device || !now.Before(c.Exp) || c.State == "cancelled" || c.State == "consumed" {
			return refusal{401, safeProblem(fault("invalid_auth")), false, false}, nil
		}
		if c.Owner != "" && (e != nil || st.Members[id].Owner != c.Owner) {
			return refusal{403, "이 연결을 시작한 Google 계정으로 다시 시작하세요.", false, false}, nil
		}
		public, _ := base64.RawURLEncoding.DecodeString(c.Public)
		v := map[string]any{"Title": "내 클라이언트 연결", "ConnectionID": r.PathValue("id"), "Fingerprint": fingerprint(public), "Exp": clock(c.Exp), "State": c.State, "Google": s.Google != nil}
		v["Login"] = c.Owner == "" || e != nil || now.Sub(session.Verified) >= reauthWindow
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
