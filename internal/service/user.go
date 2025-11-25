package service

import (
	"context"

	"github.com/isOdin-l/Assigning-PR-reviewers/internal/models"
	"github.com/isOdin-l/Assigning-PR-reviewers/pkg/api"
)

type UserRepoInterface interface {
	GetPRsWhereUserIsReviewer(ctx context.Context, userId string) (*models.PRsWhereUserIsReviewer, *models.ErrorResponse)
	PostUserSetIsActive(ctx context.Context, user *models.PostUserSetIsActive) (*models.User, *models.ErrorResponse)
}

type UserService struct {
	repo UserRepoInterface
}

func NewUserService(repo UserRepoInterface) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetPRsWhereUserIsReviewer(ctx context.Context, userId string) (*api.ResponseGetPRsWhereUserIsReviewer, *models.ErrorResponse) {
	response, err := s.repo.GetPRsWhereUserIsReviewer(ctx, userId)
	if err != nil {
		return nil, err
	}

	return models.ConvertToApiResponseGetPRsWhereUserIsReviewer(response), nil
}

func (s *UserService) PostUserSetIsActive(ctx context.Context, user *models.PostUserSetIsActive) (*api.ResponseSetUserActive, *models.ErrorResponse) {
	response, err := s.repo.PostUserSetIsActive(ctx, user)
	if err != nil {
		return nil, err
	}

	return models.ConvertToApiResponseSetUserActive(response), nil
}
