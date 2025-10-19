package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yaninyzwitty/tron-payment-gateway/packages/shared/config"
)

func DbConnect(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	url := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?pool_max_conns=%d",
		cfg.DatabaseConfig.User, cfg.DatabaseConfig.Password, cfg.DatabaseConfig.Host, cfg.DatabaseConfig.Database, cfg.DatabaseConfig.MaxConnections)

	dbpool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to create new pool %v", err)
	}
	return dbpool, nil
}
