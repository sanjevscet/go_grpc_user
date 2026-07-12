package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool(ctx context.Context, databaseUrl string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, databaseUrl)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}
	return pool, nil
}
