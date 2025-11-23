package handler

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/isOdin-l/Assigning-PR-reviewers/gen/models"
	"github.com/isOdin-l/Assigning-PR-reviewers/tool/chibind"
)

type PullRequestServiceInterface interface {
	PostPullRequestCreate(pullRequestCreateInfo models.PostPullRequestCreateJSONBody) (models.PullRequest, models.ErrorResponse)
	PostPullRequestMerge(pullRequestInfo models.PostPullRequestMergeJSONBody) (models.PullRequest, models.ErrorResponse)
	PostPullRequestReassign(pullRequestReassign *models.PostPullRequestReassignJSONBody) (*models.ResponsePullRequestReassign, *models.ErrorResponse)
}

type PullRequestHandler struct {
	service PullRequestServiceInterface
}

func NewPullRequestHandler(service PullRequestServiceInterface) *PullRequestHandler {
	return &PullRequestHandler{service: service}
}

func (h *PullRequestHandler) PostPullRequestCreate(w http.ResponseWriter, r *http.Request) {
	var pullRequestCreateInfo models.PostPullRequestCreateJSONBody
	if err := chibind.DefaultBind(r, &pullRequestCreateInfo); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		slog.Error(err.Error())
		return
	}

	response, err := h.service.PostPullRequestCreate(pullRequestCreateInfo)

	switch err.Error.Code {
	case models.NOTFOUND:
		http.Error(w, err.Error.Message, http.StatusNotFound)
		slog.Info(err.Error.Message)
		return
	case models.PREXISTS:
		http.Error(w, err.Error.Message, http.StatusConflict)
		slog.Info(err.Error.Message)
		return
	}

	render.JSON(w, r, map[string]models.PullRequest{
		"pr": response,
	})

}
func (h *PullRequestHandler) PostPullRequestMerge(w http.ResponseWriter, r *http.Request) {
	var pullRequestInfo models.PostPullRequestMergeJSONBody
	if err := chibind.DefaultBind(r, &pullRequestInfo); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		slog.Error("bad request")
		return
	}

	response, err := h.service.PostPullRequestMerge(pullRequestInfo)
	if err.Error.Code == models.NOTFOUND {
		http.Error(w, "PR не найден", http.StatusNotFound)
		slog.Info(err.Error.Message)
		return
	}

	render.JSON(w, r, map[string]models.PullRequest{
		"pr": response,
	})

}
func (h *PullRequestHandler) PostPullRequestReassign(w http.ResponseWriter, r *http.Request) {
	var pullRequestReassign models.PostPullRequestReassignJSONBody
	if err := chibind.DefaultBind(r, &pullRequestReassign); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error(err.Error())
		return
	}

	response, err := h.service.PostPullRequestReassign(&pullRequestReassign)

	switch err.Error.Code {
	case models.NOTFOUND:
		http.Error(w, err.Error.Message, http.StatusNotFound)
		slog.Info(err.Error.Message)
		return
	case models.PRMERGED:
	case models.NOTASSIGNED:
	case models.NOCANDIDATE:
		http.Error(w, err.Error.Message, http.StatusConflict)
		slog.Info(err.Error.Message)
		return
	}

	render.JSON(w, r, map[string]any{
		"pr":          response.PullRequest,
		"replaced_by": response.NewReviewerId,
	})

}
