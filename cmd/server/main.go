package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"wishlist-service/internal/config"
	"wishlist-service/internal/handler"
	"wishlist-service/internal/service"
	"wishlist-service/internal/storage"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	cfg := config.Load()
	logger.Info("starting server", "port", cfg.AppPort)

	ctx := context.Background()

	db, err := storage.New(ctx, cfg, logger)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := storage.RunMigrations(cfg.DSN(), "migrations", logger); err != nil {
		logger.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	userRepo := storage.NewUserRepository(db.DB)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := handler.NewAuthHandler(authService, logger)

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	mux.HandleFunc("POST /auth/register", authHandler.Register)
	mux.HandleFunc("POST /auth/login", authHandler.Login)

	server := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: mux,
	}

	go func() {
		logger.Info("server listening", "addr", cfg.AppPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	logger.Info("server stopped gracefully")
}
