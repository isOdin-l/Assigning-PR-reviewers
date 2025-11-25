package handler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/isOdin-l/Assigning-PR-reviewers/internal/models"
	"github.com/isOdin-l/Assigning-PR-reviewers/pkg/api"
	"github.com/isOdin-l/Assigning-PR-reviewers/tool/chibind"
)

type TeamServiceInterface interface {
	PostTeamAdd(ctx context.Context, team *models.Team) (*api.ResponseTeam, error)
	GetTeam(ctx context.Context, team *models.GetTeamParams) (*api.ResponseTeam, error)
}

type TeamHandler struct {
	service TeamServiceInterface
}

func NewTeamHandler(service TeamServiceInterface) *TeamHandler {
	return &TeamHandler{service: service}
}

func (h *TeamHandler) PostTeamAdd(w http.ResponseWriter, r *http.Request) {
	var team api.Team
	if err := chibind.DefaultBind(r, &team); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		slog.Error(fmt.Sprintf("Error while parsing data: %v", err.Error()))
		return
	}

	response, err := h.service.PostTeamAdd(r.Context(), models.ConvertToTeam(&team))
	// if err.Error.Code == models.TEAMEXISTS {
	// 	http.Error(w, "Команда уже существует", http.StatusBadRequest)
	// 	slog.Info(err.Error.Message)
	// 	return
	// }
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error(err.Error())
		return
	}

	render.JSON(w, r, *response)

}
func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	var team api.GetTeamParams
	if err := chibind.DefaultBind(r, &team); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		slog.Error(fmt.Sprintf("Erorr while parsing data: %v", err.Error()))
		return
	}

	response, err := h.service.GetTeam(r.Context(), models.ConvertToGetTeamParams(&team))
	// if err.Error.Code == models.NOTFOUND {
	// 	http.Error(w, err.Error.Message, http.StatusNotFound)
	// 	slog.Info(err.Error.Message)
	// 	return
	// }
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error(err.Error())
		return
	}

	render.JSON(w, r, *response)
}
