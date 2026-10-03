// Migration source checks reject missing directories and Go code before DB access.
package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRejectInvalidMigrationSources(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://localhost:1/knowslink")
	directory := t.TempDir()
	t.Setenv("MIGRATIONS_DIR", filepath.Join(directory, "missing"))
	if err := run(); err == nil || err.Error() != "migration directory is unavailable" {
		t.Fatalf("expected directory error, got %v", err)
	}
	if err := os.WriteFile(filepath.Join(directory, "00001_probe.go"), []byte("package migration\n"), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("MIGRATIONS_DIR", directory)
	if err := run(); err == nil || err.Error() != "only SQL migrations are supported" {
		t.Fatalf("expected SQL-only error, got %v", err)
	}
}
