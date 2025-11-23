package service

import (
	"github.com/isOdin-l/Assigning-PR-reviewers/internal/models"
)

type UserRepoInterface interface {
	GetPRsWhereUserIsReviewer(userId string) (*models.ResponsePRsWhereUserIsReviewer, error)
	PostUserSetIsActive(user *models.PostUserSetIsActive) (*models.ResponseUser, error)
}

type UserService struct {
	repo UserRepoInterface
}

func NewUserService(repo UserRepoInterface) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetPRsWhereUserIsReviewer(userId string) (*models.ResponsePRsWhereUserIsReviewer, error) {
	return s.repo.GetPRsWhereUserIsReviewer(userId)
}

func (s *UserService) PostUserSetIsActive(user *models.PostUserSetIsActive) (*models.ResponseUser, error) {
	return s.repo.PostUserSetIsActive(user)
}
