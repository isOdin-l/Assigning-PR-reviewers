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

type TeamRepo struct {
	db   *pgxpool.Pool
	psql sq.StatementBuilderType
}

func NewTeamRepo(db *pgxpool.Pool) *TeamRepo {
	return &TeamRepo{db: db, psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar)}
}

func (r *TeamRepo) CreateTeam(ctx context.Context, team *models.Team) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	isExist, err := r.isTeamExist(ctx, tx, team.TeamName)
	if err != nil {
		return err
	} else if isExist != 0 {
		return errors.New("team already exists")
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
		case 0: // Если пользователя нет, то создаём запись о нём
			if err := r.createTeamMemberTx(ctx, tx, team.TeamName, v); err != nil {
				tx.Rollback(ctx)
				return err
			}
		case 1: // Если пользователесь существует, то обновляем его команду TODO: возможно нужно обновлять и другие его данные
			if err := r.updateTeamMemberTx(ctx, tx, team.TeamName, v.UserId); err != nil {
				tx.Rollback(ctx)
				return err
			}
		}
	}
	return tx.Commit(ctx)
}

func (r *TeamRepo) GetTeam(ctx context.Context, teamName string) (*models.Team, error) {
	query, values, err := r.psql.
		Select(getSelectColumnsTeam()...).
		From(database.UsersTable).
		Where(sq.Eq{"team_name": teamName}).ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Query(ctx, query, values...)
	if err != nil {
		return nil, err
	}

	teamMembers, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.TeamMember])
	if err != nil {
		return nil, err
	}

	return &models.Team{TeamName: teamName, Members: teamMembers}, nil
}

func (r *TeamRepo) isTeamExist(ctx context.Context, tx pgx.Tx, teamName string) (int, error) {
	query, values, err := r.psql.
		Select("COUNT(1)").
		From(database.UsersTable).
		Where(sq.Eq{"team_name": teamName}).ToSql()
	if err != nil {
		return -1, err
	}
	var isExist int
	err = tx.QueryRow(ctx, query, values...).Scan(&isExist)

	return isExist, err
}

func (r *TeamRepo) isTeamMemberExistTx(ctx context.Context, tx pgx.Tx, userId string) (int, error) {
	query, values, err := r.psql.
		Select("COUNT(1)").
		From(database.UsersTable).
		Where(sq.Eq{"id": userId}).ToSql()
	if err != nil {
		return -1, nil
	}

	var isExist int
	err = tx.QueryRow(ctx, query, values...).Scan(&isExist)
	return isExist, err
}

func (r *TeamRepo) createTeamMemberTx(ctx context.Context, tx pgx.Tx, teamName string, member models.TeamMember) error {
	query, values, err := r.psql.
		Insert(database.UsersTable).
		Columns(getInsertColumnsTeamMember()...).
		Values(member.UserId, teamName, member.Username, member.IsActive).ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, query, values...)
	return err
}

func (r *TeamRepo) updateTeamMemberTx(ctx context.Context, tx pgx.Tx, teamName string, userId string) error {
	query, values, err := r.psql.
		Update(database.UsersTable).
		Set("team_name", teamName).
		Where(sq.Eq{"id": userId}).ToSql()

	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, query, values...)
	return err
}
