package service

import (
	"context"

	"github.com/isOdin-l/Assigning-PR-reviewers/internal/models"
	"github.com/isOdin-l/Assigning-PR-reviewers/pkg/api"
)

type TeamRepoInterface interface {
	GetTeam(ctx context.Context, teamName string) (*models.Team, *models.ErrorResponse)
	CreateTeam(ctx context.Context, team *models.Team) *models.ErrorResponse
}

type TeamService struct {
	repo TeamRepoInterface
}

func NewTeamService(repo TeamRepoInterface) *TeamService {
	return &TeamService{repo: repo}
}

func (s *TeamService) PostTeamAdd(ctx context.Context, team *models.Team) (*api.ResponseTeam, *models.ErrorResponse) {
	// Проверка на существование команды и создание команды
	err := s.repo.CreateTeam(ctx, team)
	if err != nil {
		return nil, err
	}

	return models.ConvertToApiResponseTeam(team), nil
}
func (s *TeamService) GetTeam(ctx context.Context, team *models.GetTeamParams) (*api.ResponseTeam, *models.ErrorResponse) {
	response, err := s.repo.GetTeam(ctx, team.TeamName)
	if err != nil {
		return nil, err
	}
	return models.ConvertToApiResponseTeam(response), nil
}
