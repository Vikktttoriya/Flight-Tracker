package database

import (
	"context"
	"fmt"
	"time"

	"github.com/Vikktttoriya/flight-tracker/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func NewPostgres(cfg config.DBConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
	)

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		zap.L().Error("Failed to parse database configuration", zap.Error(err))
		return nil, fmt.Errorf("parse pg config: %w", err)
	}

	poolCfg.MaxConns = 10
	poolCfg.MinConns = 2
	poolCfg.MaxConnIdleTime = 10 * time.Minute
	poolCfg.MaxConnLifetime = 30 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		zap.L().Error("Failed to create database connection pool", zap.Error(err))
		return nil, fmt.Errorf("create pg pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		zap.L().Error("Database connection test failed", zap.Error(err))
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return pool, nil
}
