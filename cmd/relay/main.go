// Relay serves local synthetic relay.v1 transport and owner human-gate UI.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tourbut/KnowsLink/internal/config"
	"github.com/tourbut/KnowsLink/internal/relay"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := run(ctx); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	databaseURL, err := config.DatabaseURL()
	if err != nil {
		return err
	}
	address, err := config.ListenAddress()
	if err != nil {
		return err
	}
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return errors.New("database configuration failed")
	}
	defer pool.Close()
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	err = pool.Ping(pingCtx)
	cancel()
	if err != nil {
		return errors.New("database ping failed")
	}
	mux := http.NewServeMux()
	mux.Handle("/healthz", healthHandler(pool.Ping))
	testAgents, err := relay.TestAgentAllowlist(os.Getenv("KNOWSLINK_TEST_AGENTS"))
	if err != nil {
		return err
	}
	service := &relay.Service{Pool: pool, TestAgents: testAgents, SyntheticSignup: os.Getenv("KNOWSLINK_SYNTHETIC_SIGNUP") == "1"}
	if smtpURL := os.Getenv("KNOWSLINK_SMTP_URL"); smtpURL != "" {
		if service.Mail, err = relay.SMTPMailer(smtpURL, os.Getenv("KNOWSLINK_MAIL_FROM")); err != nil {
			return err
		}
	}
	// Trust the edge client-IP header only behind the Tunnel; otherwise rate keys use the socket address.
	switch header := os.Getenv("KNOWSLINK_CLIENT_IP_HEADER"); header {
	case "", "CF-Connecting-IP":
		service.ClientIPHeader = header
	default:
		return errors.New("KNOWSLINK_CLIENT_IP_HEADER must be empty or CF-Connecting-IP")
	}
	go service.Cleanup(ctx)
	mux.Handle("/", service.Handler())
	server := &http.Server{
		Addr: address, Handler: mux,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
	result := make(chan error, 1)
	go func() { result <- server.ListenAndServe() }()
	log.Printf("relay starting; vendor effects and disclosure disabled; email login configured=%t; synthetic signup=%t", service.Mail != nil, service.SyntheticSignup)
	select {
	case err := <-result:
		if !errors.Is(err, http.ErrServerClosed) {
			return errors.New("relay listener failed")
		}
		return nil
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			return errors.New("relay shutdown failed")
		}
		return nil
	}
}

func healthHandler(ping func(context.Context) error) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := ping(ctx); err != nil {
			http.Error(w, "not ready", http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte("ok\n"))
	})
	return mux
}
