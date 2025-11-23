package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/isOdin-l/Assigning-PR-reviewers/gen/models"
	"github.com/isOdin-l/Assigning-PR-reviewers/tool/chibind"
)

type UserServiceInterface interface {
	GetUserIsReviewer(userInfo models.GetUsersGetReview) (models.ResponseUsersGetReview, error)
	PostUserSetIsActive(userInfo models.PostUserSetIsActiveJSONBody) (models.User, models.ErrorResponse)
}

type UserHandler struct {
	service UserServiceInterface
}

func NewUserHandler(service UserServiceInterface) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetUsersGetReview(w http.ResponseWriter, r *http.Request) {
	var user models.GetUsersGetReview
	if err := chibind.DefaultBind(r, &user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		slog.Error(fmt.Sprintf("Error in request data: %s", err.Error()))
		return
	}

	response, err := h.service.GetUserIsReviewer(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error(fmt.Sprintf("Error in server process: %s", err.Error()))
		return
	}

	render.JSON(w, r, map[string]any{
		"user_id":       response.User_id,
		"pull_requests": response.PullRequests,
	})
}
func (h *UserHandler) PostUserSetIsActive(w http.ResponseWriter, r *http.Request) {
	var user models.PostUserSetIsActiveJSONBody
	if err := chibind.DefaultBind(r, &user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		slog.Error(fmt.Sprintf("Error in request data: %s", err.Error()))
		return
	}

	response, err := h.service.PostUserSetIsActive(user)
	if err.Error.Code == models.NOTFOUND {
		http.Error(w, err.Error.Message, http.StatusNotFound)
		return
	}

	render.JSON(w, r, map[string]models.User{
		"user": response,
	})
}
