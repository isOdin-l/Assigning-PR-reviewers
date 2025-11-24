package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/isOdin-l/Assigning-PR-reviewers/internal/models"
	"github.com/isOdin-l/Assigning-PR-reviewers/pkg/api"
	"github.com/isOdin-l/Assigning-PR-reviewers/tool/chibind"
)

type PullRequestServiceInterface interface {
	PullRequestCreate(ctx context.Context, pullRequest *models.PullRequestCreate) (*api.ResponsePullRequestCreate, error)
	PullRequestMerge(ctx context.Context, pullRequest *models.PullRequestMerge) (*api.ResponsePullRequestMerge, error)
	PullRequestReassign(ctx context.Context, pullRequest *models.PullRequestReassign) (*api.ResponsePullRequestReassign, error)
}

type PullRequestHandler struct {
	service PullRequestServiceInterface
}

func NewPullRequestHandler(service PullRequestServiceInterface) *PullRequestHandler {
	return &PullRequestHandler{service: service}
}

func (h *PullRequestHandler) PullRequestCreate(w http.ResponseWriter, r *http.Request) {
	var pullRequest api.PullRequestCreate
	if err := chibind.DefaultBind(r, &pullRequest); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		slog.Error(err.Error())
		return
	}

	response, err := h.service.PullRequestCreate(r.Context(), models.ConvertToPullRequestCreate(&pullRequest))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error(err.Error())
		return
	}
	// switch err.Error.Code {
	// case models.NOTFOUND:
	// 	http.Error(w, err.Error.Message, http.StatusNotFound)
	// 	slog.Info(err.Error.Message)
	// 	return
	// case models.PREXISTS:
	// 	http.Error(w, err.Error.Message, http.StatusConflict)
	// 	slog.Info(err.Error.Message)
	// 	return
	// }

	render.JSON(w, r, *response)

}
func (h *PullRequestHandler) PullRequestMerge(w http.ResponseWriter, r *http.Request) {
	var pullRequest api.PullRequestMerge
	if err := chibind.DefaultBind(r, &pullRequest); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		slog.Error("bad request")
		return
	}

	response, err := h.service.PullRequestMerge(r.Context(), models.ConvertToPullRequestMerge(&pullRequest))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error(err.Error())
		return
	}
	// if err.Error.Code == models.NOTFOUND {
	// 	http.Error(w, "PR не найден", http.StatusNotFound)
	// 	slog.Info(err.Error.Message)
	// 	return
	// }

	render.JSON(w, r, *response)

}
func (h *PullRequestHandler) PullRequestReassign(w http.ResponseWriter, r *http.Request) {
	var pullRequest api.PullRequestReassign
	if err := chibind.DefaultBind(r, &pullRequest); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error(err.Error())
		return
	}

	response, err := h.service.PullRequestReassign(r.Context(), models.ConvertToPullRequestReassign(&pullRequest))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error(err.Error())
		return
	}
	// switch err.Error.Code {
	// case models.NOTFOUND:
	// 	http.Error(w, err.Error.Message, http.StatusNotFound)
	// 	slog.Info(err.Error.Message)
	// 	return
	// case models.PRMERGED:
	// case models.NOTASSIGNED:
	// case models.NOCANDIDATE:
	// 	http.Error(w, err.Error.Message, http.StatusConflict)
	// 	slog.Info(err.Error.Message)
	// 	return
	// }

	render.JSON(w, r, *response)
}
