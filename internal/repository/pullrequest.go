package repository

import "github.com/jackc/pgx/v5/pgxpool"

type PullRequestRepo struct {
	db *pgxpool.Pool
}

func NewPullRequestRepo(db *pgxpool.Pool) *PullRequestRepo {
	return &PullRequestRepo{db: db}
}
