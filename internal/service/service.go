package service

type RepoInterface interface {
	PullRequestRepoInterface
	TeamRepoInterface
	UserRepoInterface
}

type Service struct {
	PullRequestService
	TeamService
	UserService
}

func NewService(repo RepoInterface) *Service {
	return &Service{
		PullRequestService: *NewPullRequestService(repo),
		TeamService:        *NewTeamService(repo),
		UserService:        *NewUserService(repo),
	}
}
