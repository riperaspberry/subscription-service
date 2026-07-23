package database

import (
	"context"
	"fmt"

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
	pool, err := pgxpool.New(
		context.Background(),
		dsn,
	)

	if err != nil {
		return nil, err
	}

	err = pool.Ping(context.Background())

if err != nil {
	return nil, fmt.Errorf("ping failed: %w", err)
}

	return pool, nil
}