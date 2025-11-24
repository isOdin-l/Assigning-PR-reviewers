package service

import (
	"context"
	"errors"

	"github.com/isOdin-l/Assigning-PR-reviewers/internal/models"
)

type TeamRepoInterface interface {
	IsTeamExist(ctx context.Context, teamName string) (int, error)
	GetTeam(ctx context.Context, teamName string) (*models.Team, error)
	CreateTeam(ctx context.Context, team *models.Team) error
}

type TeamService struct {
	repo TeamRepoInterface
}

func NewTeamService(repo TeamRepoInterface) *TeamService {
	return &TeamService{repo: repo}
}

func (s *TeamService) PostTeamAdd(ctx context.Context, team *models.Team) (*models.Team, error) {
	// проверка на то, что команда существует
	isExist, err := s.repo.IsTeamExist(ctx, team.TeamName)
	if err != nil {
		return nil, err
	} else if isExist == 1 {
		return nil, errors.New("team is exist") // TODO: выводить нормльную ошибку - потом поработаю с кастомными ошибками
	}

	// создание команды
	err = s.repo.CreateTeam(ctx, team)
	if err != nil {
		return nil, err
	}

	return team, nil
}
func (s *TeamService) GetTeam(ctx context.Context, team *models.GetTeamGetParams) (*models.Team, error) {
	return s.repo.GetTeam(ctx, team.TeamName)
}
