package repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/isOdin-l/Assigning-PR-reviewers/internal/database"
	"github.com/isOdin-l/Assigning-PR-reviewers/internal/models"
	"github.com/isOdin-l/Assigning-PR-reviewers/pkg/api"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TeamRepo struct {
	db   *pgxpool.Pool
	psql sq.StatementBuilderType
}

func NewTeamRepo(db *pgxpool.Pool) *TeamRepo {
	return &TeamRepo{db: db, psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar)}
}

func (r *TeamRepo) CreateTeam(ctx context.Context, team *models.Team) *models.ErrorResponse {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return models.InternalError(err)
	}

	// Проверяем на то, что команды нет
	isExist, er := r.isTeamExist(ctx, tx, team.TeamName)
	if er != nil {
		tx.Rollback(ctx)
		return er
	} else if isExist != 0 {
		tx.Rollback(ctx)
		return models.ErrCreate(api.TEAMEXISTS, fmt.Sprintf("%s already exists", team.TeamName))
	}

	// Перебор всех мемберов
	for _, v := range team.Members {
		// Проверка на существование
		isExist, err := r.isTeamMemberExistTx(ctx, tx, v.UserId)
		if err != nil {
			tx.Rollback(ctx)
			return err
		}

		switch isExist {
		case false: // Если пользователя нет, то создаём запись о нём
			if err := r.createTeamMemberTx(ctx, tx, team.TeamName, v); err != nil {
				tx.Rollback(ctx)
				return err
			}
		case true: // Если пользователесь существует, то обновляем его команду TODO: возможно нужно обновлять и другие его данные
			if err := r.updateTeamMemberTx(ctx, tx, team.TeamName, v.UserId); err != nil {
				tx.Rollback(ctx)
				return err
			}
		}
	}
	if err = tx.Commit(ctx); err != nil {
		tx.Rollback(ctx)
		return models.InternalError(err)
	}

	return nil
}

func (r *TeamRepo) GetTeam(ctx context.Context, teamName string) (*models.Team, *models.ErrorResponse) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, models.InternalError(err)
	}

	isExist, er := r.isTeamExist(ctx, tx, teamName)
	if er != nil {
		tx.Rollback(ctx)
		return nil, models.InternalError(err)
	} else if isExist == 0 {
		tx.Rollback(ctx)
		return nil, models.ErrCreate(api.NOTFOUND, fmt.Sprintf("%s not found", teamName))
	}

	query, values, err := r.psql.
		Select(getSelectColumnsTeam()...).
		From(database.UsersTable).
		Where(sq.Eq{"team_name": teamName}).ToSql()
	if err != nil {
		tx.Rollback(ctx)
		return nil, models.InternalError(err)
	}
	rows, err := r.db.Query(ctx, query, values...)
	if err != nil {
		tx.Rollback(ctx)
		return nil, models.InternalError(err)
	}

	teamMembers, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.TeamMember])
	if err != nil {
		tx.Rollback(ctx)
		return nil, models.InternalError(err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, models.InternalError(err)
	}

	return &models.Team{TeamName: teamName, Members: teamMembers}, nil
}

func (r *TeamRepo) isTeamExist(ctx context.Context, tx pgx.Tx, teamName string) (int, *models.ErrorResponse) {
	query, values, err := r.psql.
		Select("COUNT(1)").
		From(database.UsersTable).
		Where(sq.Eq{"team_name": teamName}).ToSql()
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

func (r *TeamRepo) isTeamMemberExistTx(ctx context.Context, tx pgx.Tx, userId string) (bool, *models.ErrorResponse) {
	query, values, err := r.psql.
		Select("COUNT(1)").
		From(database.UsersTable).
		Where(sq.Eq{"id": userId}).ToSql()
	if err != nil {
		return false, models.InternalError(err)
	}

	var isExist int
	err = tx.QueryRow(ctx, query, values...).Scan(&isExist)
	if err != nil {
		return false, models.InternalError(err)
	}

	if isExist == 0 {
		return false, nil
	}
	return true, nil
}

func (r *TeamRepo) createTeamMemberTx(ctx context.Context, tx pgx.Tx, teamName string, member models.TeamMember) *models.ErrorResponse {
	query, values, err := r.psql.
		Insert(database.UsersTable).
		Columns(getInsertColumnsTeamMember()...).
		Values(member.UserId, teamName, member.Username, member.IsActive).ToSql()
	if err != nil {
		return models.InternalError(err)
	}

	if _, err = tx.Exec(ctx, query, values...); err != nil {
		return models.InternalError(err)
	}

	return nil
}

func (r *TeamRepo) updateTeamMemberTx(ctx context.Context, tx pgx.Tx, teamName string, userId string) *models.ErrorResponse {
	query, values, err := r.psql.
		Update(database.UsersTable).
		Set("team_name", teamName).
		Where(sq.Eq{"id": userId}).ToSql()

	if err != nil {
		return models.InternalError(err)
	}
	if _, err = tx.Exec(ctx, query, values...); err != nil {
		return models.InternalError(err)
	}

	return nil
}
