// Relay serves the relay.v1 API, member and owner web pages, and install files.
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
	_ "time/tzdata" // the runtime image has no zoneinfo; schedule.commit validates IANA zones with LoadLocation

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
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
	redisURL, err := config.RedisURL()
	if err != nil {
		return err
	}
	address, err := config.ListenAddress()
	if err != nil {
		return err
	}
	poolConfig, err := relay.PoolConfig(databaseURL)
	if err != nil {
		return errors.New("database configuration failed")
	}
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
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
	redisOptions, err := redis.ParseURL(redisURL)
	if err != nil {
		return errors.New("redis configuration failed")
	}
	redisOptions.DisableIdentity = true
	// A single self-hosted Redis sends no maintenance notices; skip the handshake probe.
	redisOptions.MaintNotificationsConfig = &maintnotifications.Config{Mode: maintnotifications.ModeDisabled}
	rdb := redis.NewClient(redisOptions)
	defer func() { _ = rdb.Close() }()
	pingCtx, cancel = context.WithTimeout(ctx, 5*time.Second)
	err = rdb.Ping(pingCtx).Err()
	cancel()
	if err != nil {
		return errors.New("redis ping failed")
	}
	mux := http.NewServeMux()
	// Ready means both stores answer: without Redis every relay transaction is unavailable.
	mux.Handle("/healthz", healthHandler(func(ctx context.Context) error {
		if err := pool.Ping(ctx); err != nil {
			return err
		}
		return rdb.Ping(ctx).Err()
	}))
	testAgents, err := relay.TestAgentAllowlist(os.Getenv("KNOWSLINK_TEST_AGENTS"))
	if err != nil {
		return err
	}
	service := &relay.Service{Pool: pool, Redis: rdb, Stopping: ctx, TestAgents: testAgents, SyntheticSignup: os.Getenv("KNOWSLINK_SYNTHETIC_SIGNUP") == "1", InstallDir: os.Getenv("KNOWSLINK_INSTALL_DIR")}
	googleCtx, googleCancel := context.WithTimeout(ctx, 6*time.Second)
	service.Google, err = relay.ConfigureGoogle(googleCtx, os.Getenv("KNOWSLINK_GOOGLE_CLIENT_ID"), os.Getenv("KNOWSLINK_GOOGLE_CLIENT_SECRET"), os.Getenv("KNOWSLINK_GOOGLE_REDIRECT_URL"))
	googleCancel()
	if err != nil {
		return err
	}
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
	go service.RemoteEvents(ctx)
	mux.Handle("/", service.Handler())
	server := &http.Server{
		Addr: address, Handler: mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
	result := make(chan error, 1)
	go func() { result <- server.ListenAndServe() }()
	log.Printf("relay starting; vendor effects and disclosure disabled; email login configured=%t; google login configured=%t; synthetic signup=%t", service.Mail != nil, service.Google != nil, service.SyntheticSignup)
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
