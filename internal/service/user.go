package service

import (
	"context"

	"github.com/isOdin-l/Assigning-PR-reviewers/internal/models"
)

type UserRepoInterface interface {
	GetPRsWhereUserIsReviewer(ctx context.Context, userId string) (*models.ResponsePRsWhereUserIsReviewer, error)
	PostUserSetIsActive(ctx context.Context, user *models.PostUserSetIsActive) (*models.User, error)
}

type UserService struct {
	repo UserRepoInterface
}

func NewUserService(repo UserRepoInterface) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetPRsWhereUserIsReviewer(ctx context.Context, userId string) (*models.ResponsePRsWhereUserIsReviewer, error) {
	return s.repo.GetPRsWhereUserIsReviewer(ctx, userId)
}

func (s *UserService) PostUserSetIsActive(ctx context.Context, user *models.PostUserSetIsActive) (*models.User, error) {
	return s.repo.PostUserSetIsActive(ctx, user)
}
