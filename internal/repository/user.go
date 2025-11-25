package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/isOdin-l/Assigning-PR-reviewers/internal/database"
	"github.com/isOdin-l/Assigning-PR-reviewers/internal/models"
	"github.com/isOdin-l/Assigning-PR-reviewers/pkg/api"
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

func (r *UserRepo) GetPRsWhereUserIsReviewer(ctx context.Context, userId string) (*models.PRsWhereUserIsReviewer, *models.ErrorResponse) {
	query, args, err := r.psql.
		Select(getColumnsUserIsReviewer()...).
		From(database.PrTable).
		InnerJoin(getJoinUserIsReviewer()).
		Where(sq.Eq{"reviewer_id": userId}).ToSql()

	if err != nil {
		return nil, models.InternalError(err)
	}

	rows, err := r.db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, models.InternalError(err)
	}

	pullRequestFromRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.PullRequestShort])
	if err != nil {
		return nil, models.InternalError(err)
	}
	return &models.PRsWhereUserIsReviewer{User_id: userId, PullRequests: pullRequestFromRows}, nil
}

func (r *UserRepo) PostUserSetIsActive(ctx context.Context, user *models.PostUserSetIsActive) (*models.User, *models.ErrorResponse) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, models.InternalError(err)
	}

	// проверяем - есть ли такой пользователь
	if err := r.isUserExist(ctx, tx, user.UserId); err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	// Обновляем значение is_active
	query, values, err := r.psql.
		Update(database.UsersTable).
		Set("is_active", user.IsActive).
		Where(sq.Eq{"id": user.UserId}).
		Suffix(getUpdateUserIsActiveSuffix()).ToSql()
	if err != nil {
		tx.Rollback(ctx)
		return nil, models.InternalError(err)
	}

	userResponse := models.ConvertToUser(user)
	err = r.db.QueryRow(ctx, query, values...).Scan(&userResponse.Username, &userResponse.TeamName)
	if err != nil {
		tx.Rollback(ctx)
		return nil, models.InternalError(err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, models.InternalError(err)
	}

	return userResponse, nil
}

func (r *UserRepo) isUserExist(ctx context.Context, tx pgx.Tx, userId string) *models.ErrorResponse {
	query, values, err := r.psql.
		Select("COUNT(1)").
		From(database.UsersTable).
		Where(sq.Eq{"id": userId}).ToSql()
	if err != nil {
		return models.InternalError(err)
	}

	var isExist int
	err = tx.QueryRow(ctx, query, values...).Scan(&isExist)
	if err != nil {
		return models.InternalError(err)
	} else if isExist != 1 {
		return models.ErrCreate(api.NOTFOUND, "пользователь не найден")
	}

	return nil
}
