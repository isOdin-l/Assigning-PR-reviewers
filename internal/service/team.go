package service

import (
	"context"

	"github.com/isOdin-l/Assigning-PR-reviewers/internal/models"
)

type TeamRepoInterface interface {
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
	// Проверка на существование команды и создание команды
	err := s.repo.CreateTeam(ctx, team)
	if err != nil {
		return nil, err
	}

	return team, nil
}
func (s *TeamService) GetTeam(ctx context.Context, team *models.GetTeamParams) (*models.Team, error) {
	return s.repo.GetTeam(ctx, team.TeamName)
}
