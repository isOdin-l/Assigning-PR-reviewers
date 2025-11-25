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

type UserServiceInterface interface {
	GetPRsWhereUserIsReviewer(ctx context.Context, userId string) (*api.ResponseGetPRsWhereUserIsReviewer, *models.ErrorResponse)
	PostUserSetIsActive(ctx context.Context, user *models.PostUserSetIsActive) (*api.ResponseSetUserActive, *models.ErrorResponse)
}

type UserHandler struct {
	service UserServiceInterface
}

func NewUserHandler(service UserServiceInterface) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetUsersGetReview(w http.ResponseWriter, r *http.Request) {
	var user api.GetUsersGetReview
	if err := chibind.DefaultBind(r, &user); err != nil {
		responser.RenderError(w, r, http.StatusInternalServerError, &models.ErrorResponse{Code: api.SERVERERROR, Message: err.Error()})
		slog.Error(fmt.Sprintf("Handler layer: %s", err.Error()))
		return
	}

	response, err := h.service.GetPRsWhereUserIsReviewer(r.Context(), user.UserId)
	if err.Code == api.NOTFOUND {
		responser.RenderError(w, r, http.StatusNotFound, err)
		slog.Info(fmt.Sprintf("Handler layer: %s", err.Message))
		return
	} else if err.Code == api.SERVERERROR {
		responser.RenderError(w, r, http.StatusInternalServerError, err)
		slog.Error(fmt.Sprintf("Hadnler layer: %s", err.Message))
	}

	responser.RenderResponse(w, r, http.StatusOK, *response)
}
func (h *UserHandler) PostUserSetIsActive(w http.ResponseWriter, r *http.Request) {
	var user api.PostUserSetIsActive
	if err := chibind.DefaultBind(r, &user); err != nil {
		responser.RenderError(w, r, http.StatusInternalServerError, &models.ErrorResponse{Code: api.SERVERERROR, Message: err.Error()})
		slog.Error(fmt.Sprintf("Handler layer: %s", err.Error()))
		return
	}

	response, err := h.service.PostUserSetIsActive(r.Context(), models.ConvertToPostUserSetIsActive(&user))
	if err.Code != api.SERVERERROR {
		responser.RenderError(w, r, http.StatusInternalServerError, err)
		slog.Error(fmt.Sprintf("Handler layer: %s", err.Message))
		return
	}

	responser.RenderResponse(w, r, http.StatusOK, *response)
}
