package db

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yaninyzwitty/tron-payment-gateway/packages/shared/config"
)

func DbConnect(ctx context.Context, cfg *config.Config, cocroachDBPass string) (*pgxpool.Pool, error) {
	userInfo := url.UserPassword(cfg.DatabaseConfig.User, cocroachDBPass)

	dbURL := url.URL{
		Scheme:   "postgres",
		User:     userInfo,
		Host:     fmt.Sprintf("%s:%d", cfg.DatabaseConfig.Host, cfg.DatabaseConfig.Port),
		Path:     cfg.DatabaseConfig.Database,
		RawQuery: "sslmode=verify-full",
	}

	dsn := dbURL.String()

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgx pool config: %w", err)
	}

	// Set pool parameters
	poolCfg.MaxConns = int32(cfg.DatabaseConfig.MaxConnections)
	poolCfg.MinConns = 2
	poolCfg.MaxConnLifetime = time.Hour

	// Initialize pool using the parsed config
	dbpool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}
	if err := dbpool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	return dbpool, nil
	// see if it works
}
