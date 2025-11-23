package postgres

import (
	"context"
	"log/slog"

	"github.com/isOdin-l/Assigning-PR-reviewers/configs"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	*pgxpool.Pool
}

func NewPostgresDB(ctx context.Context, cfg *configs.Config) (*pgxpool.Pool, error) {
	conn, err := pgxpool.New(ctx, cfg.DSN())
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, err
	}

	slog.Info("Database connected")
	return conn, nil
}
