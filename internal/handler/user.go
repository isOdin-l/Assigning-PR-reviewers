package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/isOdin-l/Assigning-PR-reviewers/internal/models"
	"github.com/isOdin-l/Assigning-PR-reviewers/pkg/api"
	"github.com/isOdin-l/Assigning-PR-reviewers/tool/chibind"
)

type UserServiceInterface interface {
	GetPRsWhereUserIsReviewer(userId string) (models.ResponsePRsWhereUserIsReviewer, error)
	PostUserSetIsActive(user *models.PostUserSetIsActive) (*models.ResponseUser, *api.ErrorResponse)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		slog.Error(fmt.Sprintf("Error in request data: %s", err.Error()))
		return
	}

	response, err := h.service.GetPRsWhereUserIsReviewer(user.UserId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error(fmt.Sprintf("Error in server process: %s", err.Error()))
		return
	}

	render.JSON(w, r, response)
}
func (h *UserHandler) PostUserSetIsActive(w http.ResponseWriter, r *http.Request) {
	var user api.PostUserSetIsActive
	if err := chibind.DefaultBind(r, &user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		slog.Error(fmt.Sprintf("Error in request data: %s", err.Error()))
		return
	}

	response, err := h.service.PostUserSetIsActive(models.ConvertToPostUserSetIsActive(user))
	if err != nil && err.Error.Code == api.NOTFOUND {
		http.Error(w, err.Error.Message, http.StatusNotFound)
		return
	}

	render.JSON(w, r, response)
}
