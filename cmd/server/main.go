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

	wishlistRepo := storage.NewWishlistRepository(db.DB)
	itemRepo := storage.NewWishlistItemRepository(db.DB)
	wishlistService := service.NewWishlistService(wishlistRepo, itemRepo)
	wishlistHandler := handler.NewWishlistHandler(wishlistService, logger)
	itemHandler := handler.NewWishlistItemHandler(wishlistService, logger)
	publicHandler := handler.NewPublicHandler(wishlistService, logger)

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	mux.HandleFunc("POST /auth/register", authHandler.Register)
	mux.HandleFunc("POST /auth/login", authHandler.Login)

	mux.Handle("POST /api/wishlists", handler.AuthMiddleware(authService, http.HandlerFunc(wishlistHandler.Create)))
	mux.Handle("GET /api/wishlists", handler.AuthMiddleware(authService, http.HandlerFunc(wishlistHandler.GetMy)))
	mux.Handle("PUT /api/wishlists/", handler.AuthMiddleware(authService, http.HandlerFunc(wishlistHandler.Update)))
	mux.Handle("DELETE /api/wishlists/", handler.AuthMiddleware(authService, http.HandlerFunc(wishlistHandler.Delete)))
	mux.Handle("POST /api/wishlists/{id}/items", handler.AuthMiddleware(authService, http.HandlerFunc(itemHandler.Create)))
	mux.Handle("GET /api/wishlists/{id}/items", handler.AuthMiddleware(authService, http.HandlerFunc(itemHandler.GetByWishlist)))
	mux.Handle("DELETE /api/wishlists/{wishlistId}/items/{itemId}", handler.AuthMiddleware(authService, http.HandlerFunc(itemHandler.Delete)))

	mux.HandleFunc("GET /public/{token}", publicHandler.GetWishlist)
	mux.HandleFunc("POST /public/{token}/reserve", publicHandler.ReserveItem)

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
