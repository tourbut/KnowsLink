// Header and version checks of the 2026-07-28 /mcp request path, without a database.
package relay

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/oauth2"
)

// SAR-MCP-UPGRADE-001: 2026-07-28 header checks run before any DB work; a request that passes them reaches the (dead) DB as 503.
func TestRemoteModernHeaderChecks(t *testing.T) {
	pool, err := pgxpool.New(context.Background(), "postgres://x@127.0.0.1:1/x?connect_timeout=1") // nothing listens
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	s := &Service{Pool: pool, Google: &GoogleLogin{config: oauth2.Config{RedirectURL: "https://relay.example/auth/google/callback"}}}
	modern := `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"knowslink_status","arguments":{},"_meta":{"io.modelcontextprotocol/protocolVersion":"%s","io.modelcontextprotocol/clientCapabilities":{}}}}`
	for _, c := range []struct {
		name, body, version, method, tool string
		code                              int
		want                              string
	}{
		{"legacy without headers", `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"x","version":"1"}}}`, "", "", "", 503, ""},
		{"legacy header", `{"jsonrpc":"2.0","id":1,"method":"tools/list"}`, "2025-11-25", "tools/list", "", 503, ""},
		{"modern matching", fmt.Sprintf(modern, "2026-07-28"), "2026-07-28", "tools/call", "knowslink_status", 503, ""},
		{"modern base64 name", fmt.Sprintf(modern, "2026-07-28"), "2026-07-28", "tools/call", "=?base64?" + base64.StdEncoding.EncodeToString([]byte("knowslink_status")) + "?=", 503, ""},
		{"version header differs from _meta", fmt.Sprintf(modern, "2026-07-28"), "2025-11-25", "tools/call", "knowslink_status", 400, `"code":-32020`},
		{"method header differs", fmt.Sprintf(modern, "2026-07-28"), "2026-07-28", "tools/list", "knowslink_status", 400, `"code":-32020`},
		{"tool header differs", fmt.Sprintf(modern, "2026-07-28"), "2026-07-28", "tools/call", "knowslink_send", 400, `"code":-32020`},
		{"legacy revision in _meta", fmt.Sprintf(modern, "2025-11-25"), "2025-11-25", "tools/call", "knowslink_status", 400, `"code":-32022`},
		// Mcp-Name mirrors params.name only on named methods; resources/read's uri passes the check (then 404 -32601 after auth).
		{"unnamed method ignores Mcp-Name", `{"jsonrpc":"2.0","id":1,"method":"resources/read","params":{"uri":"file:///x","_meta":{"io.modelcontextprotocol/protocolVersion":"2026-07-28","io.modelcontextprotocol/clientCapabilities":{}}}}`, "2026-07-28", "resources/read", "file:///x", 503, ""},
		// SAR-MCP-SPEC-CHECK-001: a modern request missing a required header or clientCapabilities is refused.
		{"modern without Mcp-Method", fmt.Sprintf(modern, "2026-07-28"), "2026-07-28", "", "knowslink_status", 400, `"code":-32020`},
		{"modern tools/call without Mcp-Name", fmt.Sprintf(modern, "2026-07-28"), "2026-07-28", "tools/call", "", 400, `"code":-32020`},
		{"modern without clientCapabilities", `{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{"_meta":{"io.modelcontextprotocol/protocolVersion":"2026-07-28"}}}`, "2026-07-28", "tools/list", "", 400, `"code":-32602`},
		{"unknown header revision", `{"jsonrpc":"2.0","id":1,"method":"tools/list"}`, "1900-01-01", "", "", 400, `"requested":"1900-01-01"`},
	} {
		r := httptest.NewRequest("POST", "/mcp", strings.NewReader(c.body))
		r.Header.Set("Authorization", "Bearer "+strings.Repeat("A", 43))
		for k, v := range map[string]string{"MCP-Protocol-Version": c.version, "Mcp-Method": c.method, "Mcp-Name": c.tool} {
			if v != "" {
				r.Header.Set(k, v)
			}
		}
		w := httptest.NewRecorder()
		s.remoteMCP(w, r)
		if w.Code != c.code || !strings.Contains(w.Body.String(), c.want) {
			t.Error(c.name, w.Code, w.Body.String())
		}
	}
}

// SAR-MCP-SPEC-CHECK-001: a 2026-07-28 DCR request carries application_type; validation passes it to the (dead) DB as 503.
func TestRemoteRegisterApplicationType(t *testing.T) {
	pool, err := pgxpool.New(context.Background(), "postgres://x@127.0.0.1:1/x?connect_timeout=1")
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	s := &Service{Pool: pool, Google: &GoogleLogin{config: oauth2.Config{RedirectURL: "https://relay.example/auth/google/callback"}}}
	for body, code := range map[string]int{
		`{"client_name":"x","redirect_uris":["http://127.0.0.1:3118/callback"],"application_type":"native"}`: 503,
		`{"client_name":"x","redirect_uris":["https://chatgpt.com/callback"],"application_type":"web"}`:      503,
	} {
		w := httptest.NewRecorder()
		s.remoteRegister(w, httptest.NewRequest("POST", "/remote/oauth/register", strings.NewReader(body)))
		if w.Code != code {
			t.Error(body, w.Code, w.Body.String())
		}
	}
}
