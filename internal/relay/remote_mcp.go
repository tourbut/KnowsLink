// Remote MCP tools use OAuth installation authority and the existing signed text, lease, persist and ACK rules.
package relay

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"time"
)

type remoteRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

func remoteResult(w http.ResponseWriter, id json.RawMessage, v any, err error) {
	out := map[string]any{"jsonrpc": "2.0", "id": id}
	if err != nil {
		out["error"] = map[string]any{"code": -32602, "message": "Request refused", "data": map[string]string{"reason": err.Error()}}
	} else {
		out["result"] = v
	}
	writeJSON(w, 200, out)
}

func (s *Service) remoteMCP(w http.ResponseWriter, r *http.Request) {
	o := s.remoteOrigin()
	var in remoteRequest
	// A 4096-byte text can need the whole raw envelope budget once a host JSON-escapes it.
	raw, err := io.ReadAll(http.MaxBytesReader(w, r.Body, rawEnvelopeLimit))
	if err != nil || Strict(raw, &in) != nil || in.JSONRPC != "2.0" || len(in.Method) > 128 {
		writeJSON(w, 400, map[string]string{"error": "invalid_request"})
		return
	}
	modern, ok := remoteEra(w, r, in)
	if !ok {
		return
	}
	var install string
	_, err = s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
		var e error
		install, _, e = st.remotePrincipal(bearer(r), o+"/mcp", now)
		return nil, e
	})
	if err != nil && statusFor(err) == 503 {
		respond(w, nil, err)
		return
	}
	if o == "" || err != nil {
		challenge := fmt.Sprintf(`Bearer resource_metadata="%s/.well-known/oauth-protected-resource", scope="%s", error="invalid_token", error_description="Connect your Google account"`, o, remoteScope)
		w.Header().Set("WWW-Authenticate", challenge)
		writeJSON(w, 401, map[string]string{"error": "invalid_token"})
		return
	}
	if len(in.ID) == 0 {
		if in.Method == "notifications/initialized" {
			w.WriteHeader(202)
			return
		}
		writeJSON(w, 400, map[string]string{"error": "invalid_request"})
		return
	}
	// A modern client tells an absent method from a refusal by 404 -32601; a legacy client keeps the refusal reason.
	if modern && !slices.Contains(remoteModernMethods, in.Method) {
		remoteError(w, 404, in.ID, -32601, "Method not found", nil)
		return
	}
	if in.Method == "events/subscribe" || in.Method == "events/unsubscribe" {
		s.remoteSubscription(w, r, install, in, modern)
		return
	}
	v, err := s.transaction(r.Context(), func(st *State, now time.Time) (any, error) {
		id, i, e := st.remotePrincipal(bearer(r), o+"/mcp", now)
		if e != nil || id != install {
			return nil, fault("invalid_auth")
		}
		switch in.Method {
		case "initialize":
			var p struct {
				Protocol     string          `json:"protocolVersion"`
				Capabilities json.RawMessage `json:"capabilities"`
				Client       json.RawMessage `json:"clientInfo"`
			}
			if Strict(in.Params, &p) != nil {
				return nil, fault("invalid_schema")
			}
			// The spec has the server answer the requested revision when it serves it, so an older client is not forced up.
			version := "2025-11-25"
			if slices.Contains(remoteVersions, p.Protocol) {
				version = p.Protocol
			}
			return map[string]any{"protocolVersion": version, "capabilities": map[string]any{"tools": map[string]any{}}, "serverInfo": remoteServerInfo}, nil
		case "server/discover":
			return map[string]any{"resultType": "complete", "supportedVersions": []string{remoteModern}, "capabilities": map[string]any{"tools": map[string]any{}, "events": map[string]any{}}, "serverInfo": remoteServerInfo, "instructions": remoteInstructions}, nil
		case "ping":
			return map[string]any{}, nil
		case "tools/list":
			return map[string]any{"tools": remoteTools()}, nil
		case "events/list":
			return map[string]any{"events": []any{remoteEventDefinition()}}, nil
		case "tools/call":
			var p struct {
				Name      string          `json:"name"`
				Arguments json.RawMessage `json:"arguments"`
				Meta      json.RawMessage `json:"_meta,omitempty"`
			}
			if Strict(in.Params, &p) != nil {
				return nil, fault("invalid_schema")
			}
			value, e := st.remoteTool(i, p.Name, p.Arguments, now)
			if e != nil {
				return nil, e
			}
			raw, _ := json.Marshal(value)
			return map[string]any{"content": []any{map[string]string{"type": "text", "text": string(raw)}}, "structuredContent": value, "isError": false}, nil
		default:
			return nil, fault("unsupported_method")
		}
	})
	remoteReply(w, in, modern, v, err)
}

func remoteUUID(now time.Time) string {
	b := make([]byte, 16)
	if _, e := rand.Read(b); e != nil {
		panic(e)
	}
	millis := now.UnixMilli()
	for index := 5; index >= 0; index-- {
		b[index] = byte(millis)
		millis >>= 8
	}
	b[6] = (b[6] & 15) | 112
	b[8] = (b[8] & 63) | 128
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[:4], b[4:6], b[6:8], b[8:10], b[10:])
}
