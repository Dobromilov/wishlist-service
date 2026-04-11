package storage

import (
	"context"
	"fmt"
	"log/slog"

	"wishlist-service/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	DB     *pgxpool.Pool
	logger *slog.Logger
}

func New(ctx context.Context, cfg config.Config, logger *slog.Logger) (*Storage, error) {
	pool, err := pgxpool.New(ctx, cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	logger.Info("connect DB")

	return &Storage{
		DB:     pool,
		logger: logger,
	}, nil
}

func (s *Storage) Close() {
	if s.DB != nil {
		s.DB.Close()
		s.logger.Info("DB connection closed")
	}
}
