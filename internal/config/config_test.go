// Configuration checks cover missing values, safe errors, and listen boundaries.
package config

import (
	"strings"
	"testing"
)

func TestDatabaseURL(t *testing.T) {
	for _, value := range []string{"", " ", "https://localhost/db", "postgres://localhost", "postgres://example-secret@%invalid/db"} {
		t.Setenv("DATABASE_URL", value)
		_, err := DatabaseURL()
		if err == nil || strings.Contains(err.Error(), "example-secret") {
			t.Fatalf("invalid setting must fail without exposing secret: %v", err)
		}
	}
	t.Setenv("DATABASE_URL", "postgres://localhost/knowslink")
	if _, err := DatabaseURL(); err != nil {
		t.Fatal(err)
	}
}

func TestListenAddress(t *testing.T) {
	t.Setenv("RELAY_ADDR", "")
	if value, err := ListenAddress(); value != "127.0.0.1:8080" || err != nil {
		t.Fatalf("unexpected default: %s, %v", value, err)
	}
	for _, value := range []string{":8080", "127.0.0.1:0", "127.0.0.1:65536", "127.0.0.1:http"} {
		t.Setenv("RELAY_ADDR", value)
		if _, err := ListenAddress(); err == nil {
			t.Fatal("invalid listen setting accepted")
		}
	}
	t.Setenv("RELAY_ADDR", "[::1]:8080")
	if _, err := ListenAddress(); err != nil {
		t.Fatal(err)
	}
}
