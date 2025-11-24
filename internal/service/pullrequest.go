package service

import (
	"context"
	"errors"

	"github.com/isOdin-l/Assigning-PR-reviewers/internal/models"
	"github.com/isOdin-l/Assigning-PR-reviewers/pkg/api"
)

type PullRequestRepoInterface interface {
	IsAuthorExist(ctx context.Context, userId string) (int, error)
	PullRequestCreate(ctx context.Context, pullRequest *models.PullRequestCreate) (*models.PullRequest, error)
	PullRequestMerge(ctx context.Context, pullRequest *models.PullRequestMerge) (*models.PullRequest, error)
	PullRequestReassign(ctx context.Context, pullRequest *models.PullRequestReassign) (*models.PullRequest, string, error)
}

type PullRequestService struct {
	repo PullRequestRepoInterface
}

func NewPullRequestService(repo PullRequestRepoInterface) *PullRequestService {
	return &PullRequestService{repo: repo}
}

func (s *PullRequestService) PullRequestCreate(ctx context.Context, pullRequest *models.PullRequestCreate) (*api.ResponsePullRequestCreate, error) {
	// Проверка на существование команды и автора
	isExists, err := s.repo.IsAuthorExist(ctx, pullRequest.AuthorId)
	if err != nil {
		return nil, err
	} else if isExists == 0 {
		return nil, errors.New("author/team not found")
	}

	// Проверка на существование PR, Создание PR, Назначение ревьюеров
	response, err := s.repo.PullRequestCreate(ctx, pullRequest)
	if err != nil {
		return nil, err
	}

	return models.ConvertToApiPullRequestCreate(*response), nil
}
func (s *PullRequestService) PullRequestMerge(ctx context.Context, pullRequest *models.PullRequestMerge) (*api.ResponsePullRequestMerge, error) {
	// Помечаем PR как Merged: проверка на merged. Если merged, то просто выдаём инфу о PR. Если не Merged, то меняем на MERGED и выводим инфу о PR
	response, err := s.repo.PullRequestMerge(ctx, pullRequest)
	if err != nil {
		return nil, err
	}

	return models.ConvertToApiPullRequestMerge(response), nil
}
func (s *PullRequestService) PullRequestReassign(ctx context.Context, pullRequest *models.PullRequestReassign) (*api.ResponsePullRequestReassign, error) {
	// Проверка на то, что PR и пользователь существуют
	// Проверка на MERGED, на то, что пользователь назначен ревьюером, на доступность кандидатов
	// Обработка 409 code
	responseRp, replacedByUserId, err := s.repo.PullRequestReassign(ctx, pullRequest)
	if err != nil {
		return nil, err
	}

	return models.ConvertToApiPullRequestReassign(responseRp, replacedByUserId), nil
}
