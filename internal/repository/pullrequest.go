package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/isOdin-l/Assigning-PR-reviewers/internal/database"
	"github.com/isOdin-l/Assigning-PR-reviewers/internal/models"
	"github.com/isOdin-l/Assigning-PR-reviewers/pkg/api"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const createPrReviewersLimit = 2
const findNewReviewer = 1

type PullRequestRepo struct {
	db   *pgxpool.Pool
	psql sq.StatementBuilderType
}

func NewPullRequestRepo(db *pgxpool.Pool) *PullRequestRepo {
	return &PullRequestRepo{db: db, psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar)}
}

func (r *PullRequestRepo) PullRequestCreate(ctx context.Context, pullRequest *models.PullRequestCreate) (*models.PullRequest, *models.ErrorResponse) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, models.InternalError(err)
	}

	// Проверка на существование автора
	if isExist, err := r.isUserExist(ctx, tx, pullRequest.AuthorId); err != nil {
		tx.Rollback(ctx)
		return nil, err
	} else if isExist == 0 {
		tx.Rollback(ctx)
		return nil, models.ErrCreate(api.NOTFOUND, "Автор/команда не найдены")
	}

	// Проверка на существование PR
	if isExist, err := r.isPrExist(ctx, tx, pullRequest.PullRequestId); err != nil {
		tx.Rollback(ctx)
		return nil, err
	} else if isExist != 0 {
		tx.Rollback(ctx)
		return nil, models.ErrCreate(api.PREXISTS, "PR id already exists")
	}

	// получение свободных людей
	reviewersIds, er := r.getNotActiveUserCreate(ctx, tx, pullRequest.AuthorId, createPrReviewersLimit)
	if er != nil {
		tx.Rollback(ctx)
		return nil, er
	}

	// Создание PR
	query, values, err := r.psql.
		Insert(database.PrTable).
		Columns(getInsertColumnsPr()...).
		Values(pullRequest.PullRequestId, pullRequest.AuthorId, pullRequest.PullRequestName).ToSql()
	if err != nil {
		tx.Rollback(ctx)
		return nil, models.InternalError(err)
	}
	_, err = tx.Exec(ctx, query, values...)
	if err != nil {
		tx.Rollback(ctx)
		return nil, models.InternalError(err)
	}

	for _, reviewer := range *reviewersIds {
		// Соединение PR и ревьюеров
		if err := r.addReviewerToPr(ctx, tx, pullRequest.PullRequestId, reviewer); err != nil {
			tx.Rollback(ctx)
			return nil, err
		}
		// Обновление status для ревьюера на false
		if err := r.setReviewerStatus(ctx, tx, reviewer, false); err != nil {
			tx.Rollback(ctx)
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, models.InternalError(err)
	}

	return &models.PullRequest{
		PullRequestId:     pullRequest.PullRequestId,
		PullRequestName:   pullRequest.PullRequestName,
		AuthorId:          pullRequest.AuthorId,
		AssignedReviewers: *reviewersIds,
		Status:            models.PullRequestStatusOPEN,
	}, nil

}
func (r *PullRequestRepo) PullRequestMerge(ctx context.Context, pullRequest *models.PullRequestMerge) (*models.PullRequest, *models.ErrorResponse) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, models.InternalError(err)
	}
	// Проверяем pr на существование
	isExist, er := r.isPrExist(ctx, tx, pullRequest.PullRequestId)
	if er != nil {
		tx.Rollback(ctx)
		return nil, er
	} else if isExist == 0 {
		tx.Rollback(ctx)
		return nil, models.ErrCreate(api.NOTFOUND, "PR not found")
	}

	// Если не Merged, то делаем ему статус Merge
	query, values, err := r.psql.
		Update(database.PrTable).
		Set("status", models.PullRequestStatusMERGED).
		Set("merged_at", time.Now().UTC()).
		Where(sq.And{
			sq.Eq{"id": pullRequest.PullRequestId},
			sq.NotEq{"status": models.PullRequestStatusMERGED},
		}).ToSql()
	if err != nil {
		tx.Rollback(ctx)
		return nil, models.InternalError(err) // TODO: ошибку скорее со стороны пользователя
	}

	if _, err = tx.Exec(ctx, query, values...); err != nil {
		return nil, models.InternalError(err)
	}

	// Получаем PullRequest
	response := &models.PullRequest{PullRequestId: pullRequest.PullRequestId}
	query, values, err = r.psql.
		Select(getSelectColumnsPrMerge()...).
		From(database.PrTable).
		Where(sq.Eq{"id": pullRequest.PullRequestId}).ToSql()
	if err != nil {
		tx.Rollback(ctx)
		return nil, models.InternalError(err)
	}
	if err = tx.QueryRow(ctx, query, values...).Scan(&response.AuthorId, &response.PullRequestName, &response.Status, &response.MergedAt); err != nil {
		tx.Rollback(ctx)
		return nil, models.InternalError(err)
	}

	// Получить AssignedReviewers
	response.AssignedReviewers, er = r.getReviewersIds(ctx, tx, pullRequest.PullRequestId)
	if er != nil {
		tx.Rollback(ctx)
		return nil, models.InternalError(err)
	}

	if err = tx.Commit(ctx); err != nil {
		tx.Rollback(ctx)
		return nil, models.InternalError(err)
	}

	return response, nil
}
func (r *PullRequestRepo) PullRequestReassign(ctx context.Context, pullRequest *models.PullRequestReassign) (*models.PullRequest, string, *models.ErrorResponse) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, "", models.InternalError(err)
	}

	// Проверка на существование ревьюера и pr (отдельно)
	if isExist, err := r.isUserExist(ctx, tx, pullRequest.OldUserId); err != nil {
		tx.Rollback(ctx)
		return nil, "", err
	} else if isExist == 0 {
		tx.Rollback(ctx)
		return nil, "", models.ErrCreate(api.NOTFOUND, "PR или пользователь не найден")
	}
	if isExist, err := r.isPrExist(ctx, tx, pullRequest.PullRequestId); err != nil {
		tx.Rollback(ctx)
		return nil, "", err
	} else if isExist == 0 {
		tx.Rollback(ctx)
		return nil, "", models.ErrCreate(api.NOTFOUND, "PR или пользователь не найден")
	}

	// Проверка на то, что PR не MERGED
	if prStatus, err := r.getPullRequestStatus(ctx, tx, pullRequest.PullRequestId); err != nil {
		tx.Rollback(ctx)
		return nil, "", err
	} else if prStatus == "MERGED" {
		tx.Rollback(ctx)
		return nil, "", models.ErrCreate(api.PRMERGED, "cannot reassign on merged PR")
	}

	// Проверка на то, является ли этот пользователь ревьюером этому PR
	if isExist, err := r.isReviewerExistInPr(ctx, tx, pullRequest.PullRequestId, pullRequest.OldUserId); err != nil {
		tx.Rollback(ctx)
		return nil, "", err
	} else if isExist == 0 {
		tx.Rollback(ctx)
		return nil, "", models.ErrCreate(api.NOTASSIGNED, "reviewer is not assigned to this PR")
	}

	// Получаем неактивного участника команды
	reviewerId, er := r.getNotActiveUserReassign(ctx, tx, pullRequest.OldUserId, pullRequest.PullRequestId, findNewReviewer)
	if er != nil {
		tx.Rollback(ctx)
		return nil, "", er
	} else if len(*reviewerId) == 0 {
		tx.Rollback(ctx)
		return nil, "", models.ErrCreate(api.NOCANDIDATE, "no active replacement candidate in team")
	}

	// Логика смены одного ревьюера на другого:
	// Удаляем прошлого ревьюера из ревьюеров
	if err := r.deleteReviewer(ctx, tx, pullRequest.OldUserId); err != nil {
		tx.Rollback(ctx)
		return nil, "", err
	}
	// Помечаем статус этого прошлого ревьюера как is_active = true
	if err := r.setReviewerStatus(ctx, tx, pullRequest.OldUserId, true); err != nil {
		tx.Rollback(ctx)
		return nil, "", err
	}
	// Добавляем нового ревьюера в таблчику к пулреквесту
	if err := r.addReviewerToPr(ctx, tx, pullRequest.PullRequestId, (*reviewerId)[0]); err != nil {
		tx.Rollback(ctx)
		return nil, "", err
	}
	// Меняем статус этого ревьюера на is_active = false
	if err := r.setReviewerStatus(ctx, tx, (*reviewerId)[0], false); err != nil {
		tx.Rollback(ctx)
		return nil, "", err
	}

	// Собираем PullRequest
	response, er := r.getPullRequest(ctx, tx, pullRequest.PullRequestId)
	if er != nil {
		tx.Rollback(ctx)
		return nil, "", er
	}
	response.AssignedReviewers, er = r.getReviewersIds(ctx, tx, pullRequest.PullRequestId)
	if er != nil {
		tx.Rollback(ctx)
		return nil, "", er
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, "", models.InternalError(err)
	}
	// Уходим в закат
	return response, (*reviewerId)[0], nil
}

func (r *PullRequestRepo) isPrExist(ctx context.Context, tx pgx.Tx, prId string) (int, *models.ErrorResponse) {
	query, values, err := r.psql.
		Select("COUNT(1)").
		From(database.PrTable).
		Where(sq.Eq{"id": prId}).ToSql()
	if err != nil {
		return -1, models.InternalError(err)
	}

	var isExist int
	err = tx.QueryRow(ctx, query, values...).Scan(&isExist)
	if err != nil {
		return -1, models.InternalError(err)
	}
	return isExist, nil
}

func (r *PullRequestRepo) isUserExist(ctx context.Context, tx pgx.Tx, userId string) (int, *models.ErrorResponse) {
	query, values, err := r.psql.
		Select("COUNT(1)").
		From(database.UsersTable).
		Where(sq.Eq{"id": userId}).ToSql()
	if err != nil {
		return -1, models.InternalError(err)
	}

	var isExist int
	if err := tx.QueryRow(ctx, query, values...).Scan(&isExist); err != nil {
		return -1, models.InternalError(err)
	}

	return isExist, nil
}

func (r *PullRequestRepo) isReviewerExistInPr(ctx context.Context, tx pgx.Tx, prId string, reviewerId string) (int, *models.ErrorResponse) {
	query, values, err := r.psql.
		Select("COUNT(1)").
		From(database.PrReviewsTable).
		Where(sq.Eq{"pr_id": prId}).
		Where(sq.Eq{"reviewer_id": reviewerId}).ToSql()
	if err != nil {
		return -1, models.InternalError(err)
	}
	var isExist int
	if err := tx.QueryRow(ctx, query, values...).Scan(&isExist); err != nil {
		return -1, models.InternalError(err)
	}

	return isExist, nil
}

func (r *PullRequestRepo) addReviewerToPr(ctx context.Context, tx pgx.Tx, prId string, reviewerId string) *models.ErrorResponse {
	query, values, err := r.psql.
		Insert(database.PrReviewsTable).
		Columns(getInsertColumnsPrReviewer()...).
		Values(prId, reviewerId).ToSql()
	if err != nil {
		return models.InternalError(err)
	}

	if _, err := tx.Exec(ctx, query, values...); err != nil {
		return models.InternalError(err)
	}
	return nil
}

func (r *PullRequestRepo) setReviewerStatus(ctx context.Context, tx pgx.Tx, reviewerId string, status bool) *models.ErrorResponse {
	query, values, err := r.psql.
		Update(database.UsersTable).
		Set("is_active", status).
		Where(sq.Eq{"id": reviewerId}).ToSql()
	if err != nil {
		return models.InternalError(err)
	}

	if _, err = tx.Exec(ctx, query, values...); err != nil {
		return models.InternalError(err)
	}

	return nil
}

func (r *PullRequestRepo) deleteReviewer(ctx context.Context, tx pgx.Tx, oldUserId string) *models.ErrorResponse {
	query, values, err := r.psql.
		Delete(database.PrReviewsTable).
		Where(sq.Eq{"reviewer_id": oldUserId}).ToSql()
	if err != nil {
		return models.InternalError(err)
	}

	if _, err = tx.Exec(ctx, query, values...); err != nil {
		return models.InternalError(err)
	}

	return nil
}

func (r *PullRequestRepo) getReviewersIds(ctx context.Context, tx pgx.Tx, pullRequestId string) ([]string, *models.ErrorResponse) {
	query, values, err := r.psql.
		Select("reviewer_id").
		From(database.PrReviewsTable).
		Where(sq.Eq{"pr_id": pullRequestId}).ToSql()
	if err != nil {
		return nil, models.InternalError(err)
	}

	rows, err := tx.Query(ctx, query, values...)
	if err != nil {
		return nil, models.InternalError(err)
	}

	reviewersIds, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.ReviewerId])
	if err != nil {
		return nil, models.InternalError(err)
	}

	return *models.ConvertToStringReviewer(&reviewersIds), nil
}

func (r *PullRequestRepo) getNotActiveUserCreate(ctx context.Context, tx pgx.Tx, authorId string, limitNewReviewers uint64) (*[]string, *models.ErrorResponse) {
	// Получаем team_name
	subquery, subvalues, err := r.psql.
		Select("team_name").
		From(database.UsersTable).
		Where(sq.Eq{"id": authorId}).ToSql()
	if err != nil {
		return nil, models.InternalError(err)
	}

	var team_name string
	if err = tx.QueryRow(ctx, subquery, subvalues...).Scan(&team_name); err != nil {
		return nil, models.InternalError(err)
	}

	// Находим количество участников команды от 0 до limitNewReviewers, у которых is_active = true, которые не userId, которые в одной team_name с userId
	query, values, err := r.psql.
		Select("id").
		From(database.UsersTable).
		Where(sq.Eq{"is_active": true}).
		Where(sq.Eq{"team_name": team_name}).
		Where(sq.NotEq{"id": authorId}).
		Limit(limitNewReviewers).ToSql()
	if err != nil {
		return nil, models.InternalError(err)
	}

	rows, err := tx.Query(ctx, query, values...)
	if err != nil {
		return nil, models.InternalError(err)
	}
	reviewers, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.UserId])
	if err != nil {
		return nil, models.InternalError(err)
	}

	return models.ConvertToStringUser(&reviewers), nil
}

func (r *PullRequestRepo) getNotActiveUserReassign(ctx context.Context, tx pgx.Tx, userId string, prId string, limitNewReviewers uint64) (*[]string, *models.ErrorResponse) {
	// Получаем автора
	query, values, err := r.psql.
		Select("author_id").
		From(database.PrTable).
		Where(sq.Eq{"id": prId}).ToSql()
	if err != nil {
		return nil, models.InternalError(err)
	}

	var author_id string
	if err = tx.QueryRow(ctx, query, values...).Scan(&author_id); err != nil {
		return nil, models.InternalError(err)
	}

	// Получаем team_name
	subquery, subvalues, err := r.psql.
		Select("team_name").
		From(database.UsersTable).
		Where(sq.Eq{"id": userId}).ToSql()
	if err != nil {
		return nil, models.InternalError(err)
	}

	var team_name string
	if err = tx.QueryRow(ctx, subquery, subvalues...).Scan(&team_name); err != nil {
		return nil, models.InternalError(err)
	}

	// Находим количество участников команды от 0 до limitNewReviewers, у которых is_active = true, которые не userId, author_id, которые в одной team_name с userId
	query, values, err = r.psql.
		Select("id").
		From(database.UsersTable).
		Where(sq.Eq{"is_active": true}).
		Where(sq.Eq{"team_name": team_name}).
		Where(sq.NotEq{"id": userId}).
		Where(sq.NotEq{"id": author_id}).
		Limit(limitNewReviewers).ToSql()
	if err != nil {
		return nil, models.InternalError(err)
	}

	rows, err := tx.Query(ctx, query, values...)
	if err != nil {
		return nil, models.InternalError(err)
	}
	reviewers, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.UserId])
	if err != nil {
		return nil, models.InternalError(err)
	}

	return models.ConvertToStringUser(&reviewers), nil
}

func (r *PullRequestRepo) getPullRequestStatus(ctx context.Context, tx pgx.Tx, pullRequestId string) (string, *models.ErrorResponse) {
	query, values, err := r.psql.
		Select("status").
		From(database.PrTable).
		Where(sq.Eq{"id": pullRequestId}).ToSql()
	if err != nil {
		return "", models.InternalError(err)
	}

	var prStatus string
	if err = tx.QueryRow(ctx, query, values...).Scan(&prStatus); err != nil {
		return "", models.InternalError(err)
	}

	return prStatus, nil
}

func (r *PullRequestRepo) getPullRequest(ctx context.Context, tx pgx.Tx, pullRequestId string) (*models.PullRequest, *models.ErrorResponse) {
	// Получаем id, name, author_id, status для PR
	query, values, err := r.psql.
		Select(getSelectColumnsPr()...).
		From(database.PrTable).
		Where(sq.Eq{"id": pullRequestId}).ToSql()
	if err != nil {
		return nil, models.InternalError(err)
	}

	var response models.PullRequest
	if err := tx.QueryRow(ctx, query, values...).Scan(&response.PullRequestId, &response.AuthorId, &response.PullRequestName, &response.Status); err != nil {
		return nil, models.InternalError(err)
	}

	return &response, nil
}
