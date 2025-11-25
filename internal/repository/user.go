package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/isOdin-l/Assigning-PR-reviewers/internal/database"
	"github.com/isOdin-l/Assigning-PR-reviewers/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db   *pgxpool.Pool
	psql sq.StatementBuilderType
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db, psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar)}
}

func (r *UserRepo) GetPRsWhereUserIsReviewer(ctx context.Context, userId string) (*models.PRsWhereUserIsReviewer, error) {
	query, args, err := r.psql.
		Select(getColumnsUserIsReviewer()...).
		From(database.PrTable).
		InnerJoin(getJoinUserIsReviewer()).
		Where(sq.Eq{"reviewer_id": userId}).ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}

	pullRequestFromRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.PullRequestShort])
	if err != nil {
		return nil, err
	} // проверить - мб ошибка будет при пустом pullrequest, тогда надо задавать pullrequest как пустой массив
	return &models.PRsWhereUserIsReviewer{User_id: userId, PullRequests: pullRequestFromRows}, nil
}

func (r *UserRepo) PostUserSetIsActive(ctx context.Context, user *models.PostUserSetIsActive) (*models.User, error) { // error -> api.ErrorResponse
	query, values, err := r.psql.
		Update(database.UsersTable).
		Set("is_active", user.IsActive).
		Where(sq.Eq{"id": user.UserId}).
		Suffix(getUpdateUserIsActiveSuffix()).ToSql()
	if err != nil {
		return nil, err
	}

	userResponse := models.ConvertToUser(user)
	r.db.QueryRow(ctx, query, values...).Scan(&userResponse.Username, &userResponse.TeamName)

	return userResponse, nil
}
