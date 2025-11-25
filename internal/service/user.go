package service

import (
	"context"

	"github.com/isOdin-l/Assigning-PR-reviewers/internal/models"
	"github.com/isOdin-l/Assigning-PR-reviewers/pkg/api"
)

type UserRepoInterface interface {
	GetPRsWhereUserIsReviewer(ctx context.Context, userId string) (*models.PRsWhereUserIsReviewer, error)
	PostUserSetIsActive(ctx context.Context, user *models.PostUserSetIsActive) (*models.User, error)
}

type UserService struct {
	repo UserRepoInterface
}

func NewUserService(repo UserRepoInterface) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetPRsWhereUserIsReviewer(ctx context.Context, userId string) (*api.ResponseGetPRsWhereUserIsReviewer, error) {
	response, err := s.repo.GetPRsWhereUserIsReviewer(ctx, userId)
	if err != nil {
		return nil, err
	}

	return models.ConvertToApiResponseGetPRsWhereUserIsReviewer(response), nil
}

func (s *UserService) PostUserSetIsActive(ctx context.Context, user *models.PostUserSetIsActive) (*api.ResponseSetUserActive, error) {
	response, err := s.repo.PostUserSetIsActive(ctx, user)
	if err != nil {
		return nil, err
	}

	return models.ConvertToApiResponseSetUserActive(response), nil
}
