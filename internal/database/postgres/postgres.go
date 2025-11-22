package postgres

import (
	"context"
	"log/slog"

	"github.com/isOdin-l/Assigning-PR-reviewers/configs"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresDB(ctx context.Context, cfg *configs.Config) (*pgxpool.Pool, error) {
	conectionString := cfg.DSN()
	conn, err := pgxpool.New(ctx, conectionString)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, err
	}

	slog.Info("Database connected")
	return conn, nil
}
