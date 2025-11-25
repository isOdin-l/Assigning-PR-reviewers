package handler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/isOdin-l/Assigning-PR-reviewers/internal/models"
	"github.com/isOdin-l/Assigning-PR-reviewers/pkg/api"
	"github.com/isOdin-l/Assigning-PR-reviewers/tool/chibind"
	"github.com/isOdin-l/Assigning-PR-reviewers/tool/responser"
)

type PullRequestServiceInterface interface {
	PullRequestCreate(ctx context.Context, pullRequest *models.PullRequestCreate) (*api.ResponsePullRequestCreate, *models.ErrorResponse)
	PullRequestMerge(ctx context.Context, pullRequest *models.PullRequestMerge) (*api.ResponsePullRequestMerge, *models.ErrorResponse)
	PullRequestReassign(ctx context.Context, pullRequest *models.PullRequestReassign) (*api.ResponsePullRequestReassign, *models.ErrorResponse)
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
		responser.RenderError(w, r, http.StatusInternalServerError, &models.ErrorResponse{Code: api.SERVERERROR, Message: err.Error()})
		slog.Error(fmt.Sprintf("Handler layer: %s", err.Error()))
		return
	}

	response, err := h.service.PullRequestCreate(r.Context(), models.ConvertToPullRequestCreate(&pullRequest))
	switch err.Code {
	case api.NOTFOUND:
		responser.RenderError(w, r, http.StatusNotFound, err)
		slog.Info(fmt.Sprintf("Handler layer: %s", err.Message))
		return
	case api.PREXISTS:
		responser.RenderError(w, r, http.StatusConflict, err)
		slog.Info(fmt.Sprintf("Handler layer: %s", err.Message))
		return
	case api.SERVERERROR:
		responser.RenderError(w, r, http.StatusInternalServerError, err)
		slog.Error(fmt.Sprintf("Handler layer: %s", err.Message))
		return
	}

	responser.RenderResponse(w, r, http.StatusCreated, *response)
}
func (h *PullRequestHandler) PullRequestMerge(w http.ResponseWriter, r *http.Request) {
	var pullRequest api.PullRequestMerge
	if err := chibind.DefaultBind(r, &pullRequest); err != nil {
		responser.RenderError(w, r, http.StatusInternalServerError, &models.ErrorResponse{Code: api.SERVERERROR, Message: err.Error()})
		slog.Error(fmt.Sprintf("Handler layer: %s", err.Error()))
		return
	}

	response, err := h.service.PullRequestMerge(r.Context(), models.ConvertToPullRequestMerge(&pullRequest))
	if err.Code != api.NOTFOUND {
		responser.RenderError(w, r, http.StatusNotFound, err)
		slog.Info(fmt.Sprintf("Handler layer: %s", err.Message))
		return
	}

	responser.RenderResponse(w, r, http.StatusOK, *response)
}
func (h *PullRequestHandler) PullRequestReassign(w http.ResponseWriter, r *http.Request) {
	var pullRequest api.PullRequestReassign
	if err := chibind.DefaultBind(r, &pullRequest); err != nil {
		responser.RenderError(w, r, http.StatusInternalServerError, &models.ErrorResponse{Code: api.SERVERERROR, Message: err.Error()})
		slog.Error(fmt.Sprintf("Handler layer: %s", err.Error()))
		return
	}

	response, err := h.service.PullRequestReassign(r.Context(), models.ConvertToPullRequestReassign(&pullRequest))
	switch err.Code {
	case api.NOTFOUND:
		responser.RenderError(w, r, http.StatusNotFound, err)
		slog.Info(fmt.Sprintf("Handler layer: %s", err.Message))
		return
	case api.PRMERGED:
	case api.NOTASSIGNED:
	case api.NOCANDIDATE:
		responser.RenderError(w, r, http.StatusConflict, err)
		slog.Info(fmt.Sprintf("Handler layer: %s", err.Message))
		return
	case api.SERVERERROR:
		responser.RenderError(w, r, http.StatusInternalServerError, err)
		slog.Error(fmt.Sprintf("Handler layer: %s", err.Message))
		return
	}

	responser.RenderResponse(w, r, http.StatusOK, *response)
}
