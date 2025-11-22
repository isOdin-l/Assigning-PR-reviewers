package httpchi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type HandlerInterface interface {
	PostPullRequestCreate(w http.ResponseWriter, r *http.Request)
	PostPullRequestMerge(w http.ResponseWriter, r *http.Request)
	PostPullRequestReassign(w http.ResponseWriter, r *http.Request)
	PostTeamAdd(w http.ResponseWriter, r *http.Request)
	GetTeamGet(w http.ResponseWriter, r *http.Request)
	GetUsersGetReview(w http.ResponseWriter, r *http.Request)
	PostUsersSetIsActive(w http.ResponseWriter, r *http.Request)
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
			r.Post("/create", h.PostPullRequestCreate)
			r.Post("/merge", h.PostPullRequestMerge)
			r.Post("/reassign", h.PostPullRequestReassign)
		})

		r.Route("/team", func(r chi.Router) {
			r.Post("/add", h.PostTeamAdd)
			r.Get("/get", h.GetTeamGet)
		})

		r.Route("/users", func(r chi.Router) {
			r.Post("/setIsActive", h.PostUsersSetIsActive)
			r.Get("/getReview", h.GetUsersGetReview)
		})
	})
	return r
}
