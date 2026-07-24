package database

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/riperaspberry/subscription-service/internal/config"
)

func New(cfg *config.Config) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	slog.Info("connecting to database",
		"host", cfg.DBHost,
		"port", cfg.DBPort,
		"database", cfg.DBName,
	)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		slog.Error("failed to create database pool", "error", err)
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		slog.Error("database ping failed", "error", err)
		pool.Close()
		return nil, fmt.Errorf("ping failed: %w", err)
	}

	slog.Info("database connection established")

	return pool, nil
}
