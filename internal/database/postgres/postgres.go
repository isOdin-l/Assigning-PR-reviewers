package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/isOdin-l/Assigning-PR-reviewers/configs"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresDB(ctx context.Context, cfg *configs.Config) (*pgxpool.Pool, error) {
	// TODO: переместить подключение в cfg
	conectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.DbUsername, cfg.DbPassword, cfg.DdHost, cfg.DbPort, cfg.DbName)
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
