package service

type PullRequestRepoInterface interface {
}

type PullRequestService struct {
	repo PullRequestRepoInterface
}

func NewPullRequestService(repo PullRequestRepoInterface) *PullRequestService {
	return &PullRequestService{repo: repo}
}
