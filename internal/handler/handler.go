package handler

type ServiceInterface interface {
	PullRequestServiceInterface
	TeamServiceInterface
	UserServiceInterface
}

type Handler struct {
	PullRequestHandler
	TeamHandler
	UserHandler
}

func NewHandler(service ServiceInterface) *Handler {
	return &Handler{
		PullRequestHandler: *NewPullRequestHandler(service),
		TeamHandler:        *NewTeamHandler(service),
		UserHandler:        *NewUserHandler(service),
	}
}
