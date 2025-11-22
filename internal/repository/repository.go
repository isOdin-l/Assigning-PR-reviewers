package repository

import "github.com/jackc/pgx/v5/pgxpool"

type Repository struct {
	PullRequestRepo
	TeamRepo
	UserRepo
}

// TODO: Change pgxpool.Pool to db interface
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		PullRequestRepo: *NewPullRequestRepo(db),
		TeamRepo:        *NewTeamRepo(db),
		UserRepo:        *NewUserRepo(db),
	}
}
