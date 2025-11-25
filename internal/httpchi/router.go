package httpchi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type HandlerInterface interface {
	PullRequestCreate(w http.ResponseWriter, r *http.Request)
	PullRequestMerge(w http.ResponseWriter, r *http.Request)
	PullRequestReassign(w http.ResponseWriter, r *http.Request)
	PostTeamAdd(w http.ResponseWriter, r *http.Request)
	GetTeam(w http.ResponseWriter, r *http.Request)
	GetUsersGetReview(w http.ResponseWriter, r *http.Request)
	PostUserSetIsActive(w http.ResponseWriter, r *http.Request)
}

func NewRouter(h HandlerInterface) chi.Router {
	r := chi.NewRouter()

	// if options.ErrorHandlerFunc == nil {
	// 	options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
	// 		http.Error(w, err.Error(), http.StatusBadRequest)
	// 	}
	// }
	// wrapper := ServerInterfaceWrapper{
	// 	Handler:            si,
	// 	HandlerMiddlewares: options.Middlewares,
	// 	ErrorHandlerFunc:   options.ErrorHandlerFunc,
	// }

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v0", func(r chi.Router) {
		r.Route("/pullRequest", func(r chi.Router) {
			r.Post("/create", h.PullRequestCreate)
			r.Post("/merge", h.PullRequestMerge)
			r.Post("/reassign", h.PullRequestReassign)
		})

		r.Route("/team", func(r chi.Router) {
			r.Post("/add", h.PostTeamAdd)
			r.Get("/get/{team_name}", h.GetTeam)
		})

		r.Route("/users", func(r chi.Router) {
			r.Post("/setIsActive", h.PostUserSetIsActive)
			r.Get("/getReview/{user_id}", h.GetUsersGetReview)
		})
	})
	return r
}
