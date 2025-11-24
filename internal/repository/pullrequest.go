package repository

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/isOdin-l/Assigning-PR-reviewers/internal/database"
	"github.com/isOdin-l/Assigning-PR-reviewers/internal/models"
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

func (r *PullRequestRepo) IsAuthorExist(ctx context.Context, userId string) (int, error) {
	query, values, err := r.psql.
		Select("COUNT(1)").
		From(database.UsersTable).
		Where(sq.Eq{"id": userId}).ToSql()
	if err != nil {
		return -1, err
	}

	var isExist int
	err = r.db.QueryRow(ctx, query, values...).Scan(&isExist)
	if err != nil {
		return -1, err
	}
	return isExist, nil
}
func (r *PullRequestRepo) PullRequestCreate(ctx context.Context, pullRequest *models.PullRequestCreate) (*models.PullRequest, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	// Проверка на существование PR
	isExist, err := r.isPrWithAuthorExist(ctx, tx, pullRequest.PullRequestId, pullRequest.AuthorId)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	} else if isExist == 1 {
		tx.Rollback(ctx)
		return nil, errors.New("author/team already exist")
	}

	// получение свободных людей
	reviewersIds, err := r.getNotActiveUser(ctx, tx, pullRequest.AuthorId, createPrReviewersLimit)
	if err != pgx.ErrNoRows && err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	// Создание PR
	query, values, err := r.psql.
		Insert(database.PrTable).
		Columns(getInsertColumnsPr()...).
		Values(getInsertValuesPr(pullRequest)).ToSql()
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}
	_, err = tx.Exec(ctx, query, values...)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	for _, reviewer_id := range reviewersIds {
		// Соединение PR и ревьюеров
		if err = r.addReviewerToPr(ctx, tx, pullRequest.PullRequestId, reviewer_id); err != nil {
			tx.Rollback(ctx)
			return nil, err
		}
		_, err = tx.Exec(ctx, query, values...)
		if err != nil {
			tx.Rollback(ctx)
			return nil, err
		}

		// Обновление status для ревьюера на false
		if err = r.setReviewerStatus(ctx, tx, reviewer_id, false); err != nil {
			tx.Rollback(ctx)
			return nil, err
		}
	}

	return &models.PullRequest{
		PullRequestId:     pullRequest.PullRequestId,
		PullRequestName:   pullRequest.PullRequestName,
		AuthorId:          pullRequest.AuthorId,
		AssignedReviewers: reviewersIds,
		Status:            models.PullRequestStatusOPEN,
	}, tx.Commit(ctx)

}

func (r *PullRequestRepo) PullRequestMerge(ctx context.Context, pullRequest *models.PullRequestMerge) (*models.PullRequest, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	// Если не Merged, то делаем ему статус Merge
	query, values, err := r.psql.
		Update(database.PrTable).
		Set("status", models.PullRequestStatusMERGED).
		Where(sq.And{
			sq.Eq{"id": pullRequest.PullRequestId},
			sq.NotEq{"status": models.PullRequestStatusMERGED},
		}).ToSql()
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	tx.Exec(ctx, query, values...)

	// Получаем PullRequest
	response := &models.PullRequest{PullRequestId: pullRequest.PullRequestId}
	query, values, err = r.psql.
		Select(getSelectColumnsPrMerge()...).
		From(database.PrTable).
		Where(sq.Eq{"id": pullRequest.PullRequestId}).ToSql()
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}
	err = tx.QueryRow(ctx, query, values...).Scan(&response.AuthorId, response.PullRequestName, response.Status, response.MergedAt)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	// Получить AssignedReviewers
	response.AssignedReviewers, err = r.getReviewersIds(ctx, tx, pullRequest.PullRequestId)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	return response, tx.Commit(ctx)
}
func (r *PullRequestRepo) PullRequestReassign(ctx context.Context, pullRequest *models.PullRequestReassign) (*models.PullRequest, string, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, "", err
	}

	// Проверка на то, что PR не MERGED
	prStatus, err := r.getPullRequestStatus(ctx, tx, pullRequest.PullRequestId)
	if err != nil {
		tx.Rollback(ctx)
		return nil, "", err
	} else if string(prStatus) == string(models.PullRequestStatusMERGED) {
		tx.Rollback(ctx)
		return nil, "", errors.New("cannot reassign on merged PR")
	}

	// Проверка на то, есть ли такой ревьюер в этом PR
	isExist, err := r.isReviewerExistInPr(ctx, tx, pullRequest.PullRequestId, pullRequest.OldUserId)
	if err != nil {
		tx.Rollback(ctx)
		return nil, "", err
	} else if isExist == 0 {
		tx.Rollback(ctx)
		return nil, "", errors.New("reviewer is not assigned to this PR")
	}

	// Получаем неактивного участника команды; TODO: сделать проверку на то, что мы получили 0 пользователей
	reviewerId, err := r.getNotActiveUser(ctx, tx, pullRequest.OldUserId, findNewReviewer)
	if err != nil {
		tx.Rollback(ctx)
		return nil, "", err
	}

	// Удаляем прошлого ревьюера из ревьюеров
	if err = r.deleteReviewer(ctx, tx, pullRequest.OldUserId); err != nil {
		tx.Rollback(ctx)
		return nil, "", err
	}

	// Помечаем статус этого прошлого ревьюера как is_active = true
	if err = r.setReviewerStatus(ctx, tx, pullRequest.OldUserId, true); err != nil {
		tx.Rollback(ctx)
		return nil, "", err
	}

	// Добавляем нового ревьюера в таблчику к пулреквесту
	if err = r.addReviewerToPr(ctx, tx, pullRequest.PullRequestId, reviewerId[0]); err != nil {
		tx.Rollback(ctx)
		return nil, "", err
	}
	// Меняем статус этого ревьюера на is_active = false
	if err = r.setReviewerStatus(ctx, tx, reviewerId[0], false); err != nil {
		tx.Rollback(ctx)
		return nil, "", err
	}

	// Собираем PullRequest
	response, err := r.getPullRequest(ctx, tx, pullRequest.PullRequestId)
	if err != nil {
		tx.Rollback(ctx)
		return nil, "", err
	}
	response.AssignedReviewers, err = r.getReviewersIds(ctx, tx, pullRequest.PullRequestId)
	if err != nil {
		tx.Rollback(ctx)
		return nil, "", err
	}

	// Уходим в закат
	return response, reviewerId[0], nil
}

func (r *PullRequestRepo) isPrWithAuthorExist(ctx context.Context, tx pgx.Tx, prId string, authorId string) (int, error) {
	query, values, err := r.psql.
		Select("COUNT(1)").
		From(database.PrTable).
		Where(sq.Eq{"id": prId, "author_id": authorId}).ToSql()
	if err != nil {
		return -1, nil
	}

	var isExist int
	err = tx.QueryRow(ctx, query, values...).Scan(&isExist)
	if err != nil {
		return -1, err
	}
	return isExist, nil
}

func (r *PullRequestRepo) addReviewerToPr(ctx context.Context, tx pgx.Tx, prId string, reviewerId string) error {
	query, values, err := r.psql.
		Insert(database.PrReviewsTable).
		Columns(getInsertColumnsPrReviewer()...).
		Values(prId, reviewerId).ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, query, values...)
	return err
}

func (r *PullRequestRepo) setReviewerStatus(ctx context.Context, tx pgx.Tx, reviewerId string, status bool) error {
	query, values, err := r.psql.
		Update(database.UsersTable).
		Set("is_active", status).
		Where(sq.Eq{"id": reviewerId}).ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, query, values...)
	return err
}

func (r *PullRequestRepo) deleteReviewer(ctx context.Context, tx pgx.Tx, oldUserId string) error {
	query, values, err := r.psql.
		Delete(database.PrReviewsTable).
		Where(sq.Eq{"reviewer_id": oldUserId}).ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, query, values...)
	return err
}

func (r *PullRequestRepo) getReviewersIds(ctx context.Context, tx pgx.Tx, pullRequestId string) ([]string, error) {
	query, values, err := r.psql.
		Select("reviewer_id").
		From(database.PrReviewsTable).
		Where(sq.Eq{"pr_id": pullRequestId}).ToSql()
	if err != nil {
		return []string{}, err
	}

	rows, err := tx.Query(ctx, query, values...)
	if err != nil {
		return []string{}, err
	}

	reviewersIds, err := pgx.CollectRows(rows, pgx.RowToStructByName[string])
	if err != nil {
		return []string{}, err
	}

	return reviewersIds, nil
}

func (r *PullRequestRepo) getNotActiveUser(ctx context.Context, tx pgx.Tx, userId string, limitNewReviewers uint64) ([]string, error) {
	// Получаем team_name
	subquery := r.psql.
		Select("team_name").
		From(database.UsersTable).
		Where(sq.Eq{"id": userId})

	// Находим количество участников команды от 0 до limitNewReviewers, у которых is_active = true, которые не userId, которые в одной team_name с userId
	query, values, err := r.psql.
		Select("id").
		From(database.UsersTable).
		Where("team_name = ($)", subquery).
		Where(sq.Eq{"is_active": true}, sq.NotEq{"id": userId}).
		Limit(limitNewReviewers).ToSql()
	if err != nil {
		return []string{}, err
	}

	rows, err := tx.Query(ctx, query, values...)
	if err != nil {
		return []string{}, err
	}
	reviewers, err := pgx.CollectRows(rows, pgx.RowToStructByName[string])
	if err != nil {
		return []string{}, err
	}

	return reviewers, nil

}

func (r *PullRequestRepo) getPullRequestStatus(ctx context.Context, tx pgx.Tx, pullRequestId string) (string, error) {
	query, values, err := r.psql.
		Select("status").
		From(database.PrTable).
		Where(sq.Eq{"id": pullRequestId}).ToSql()
	if err != nil {
		return "", err
	}

	var prStatus string
	if err = tx.QueryRow(ctx, query, values...).Scan(&prStatus); err != nil {
		return "", err
	}

	return prStatus, nil
}

func (r *PullRequestRepo) isReviewerExistInPr(ctx context.Context, tx pgx.Tx, prId string, reviewerId string) (int, error) {
	query, values, err := r.psql.
		Select("COUNT(1)").
		From(database.PrReviewsTable).
		Where(sq.Eq{"pr_id": prId, "reviewer_id": reviewerId}).ToSql()
	if err != nil {
		return -1, err
	}
	var isExist int
	err = tx.QueryRow(ctx, query, values...).Scan(&isExist)

	return isExist, err
}

func (r *PullRequestRepo) getPullRequest(ctx context.Context, tx pgx.Tx, pullRequestId string) (*models.PullRequest, error) {
	// Получаем id, name, author_id, status для PR
	query, values, err := r.psql.
		Select(getSelectColumnsPr()...).
		From(database.PrTable).
		Where(sq.Eq{"id": pullRequestId}).ToSql()
	if err != nil {
		return nil, err
	}

	var response models.PullRequest
	if err := tx.QueryRow(ctx, query, values...).Scan(&response.PullRequestId, &response.AuthorId, &response.PullRequestName, &response.Status); err != nil {
		return nil, err
	}

	return &response, nil
}
