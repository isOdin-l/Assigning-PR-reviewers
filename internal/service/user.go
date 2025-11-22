package service

type UserRepoInterface interface {
}

type UserService struct {
	repo UserRepoInterface
}

func NewUserService(repo UserRepoInterface) *UserService {
	return &UserService{repo: repo}
}
