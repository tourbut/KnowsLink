// Migrate applies SQL-only goose migrations separately from relay startup.
package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/tourbut/KnowsLink/internal/config"
)

func main() {
	if err := run(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func run() error {
	databaseURL, err := config.DatabaseURL()
	if err != nil {
		return err
	}
	directory := os.Getenv("MIGRATIONS_DIR")
	if directory == "" {
		directory = "db/migrations"
	}
	entries, err := os.ReadDir(directory)
	if err != nil {
		return errors.New("migration directory is unavailable")
	}
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".go") {
			return errors.New("only SQL migrations are supported")
		}
	}
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return errors.New("database configuration failed")
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return errors.New("database ping failed")
	}
	provider, err := goose.NewProvider(goose.DialectPostgres, db, os.DirFS(directory), goose.WithDisableGlobalRegistry(true))
	if errors.Is(err, goose.ErrNoMigrations) {
		log.Print("no SQL migrations; no-op (no business schema exists)")
		return nil
	}
	if err != nil {
		return errors.New("SQL migration configuration failed")
	}
	if _, err := provider.Up(ctx); err != nil {
		return errors.New("SQL migration execution failed")
	}
	log.Print("SQL migrations completed")
	return nil
}
