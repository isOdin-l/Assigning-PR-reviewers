package service

type TeamRepoInterface interface {
}

type TeamService struct {
	repo TeamRepoInterface
}

func NewTeamService(repo TeamRepoInterface) *TeamService {
	return &TeamService{repo: repo}
}
