package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/isOdin-l/Assigning-PR-reviewers/gen/models"
	"github.com/isOdin-l/Assigning-PR-reviewers/tool/chibind"
)

type TeamServiceInterface interface {
	GetTeamGet(userInfo models.GetTeamGetParams) (models.Team, models.ErrorResponse)
	PostTeamAdd(teamInfo models.Team) (models.Team, models.ErrorResponse)
}

type TeamHandler struct {
	service TeamServiceInterface
}

func NewTeamHandler(service TeamServiceInterface) *TeamHandler {
	return &TeamHandler{service: service}
}

func (h *TeamHandler) PostTeamAdd(w http.ResponseWriter, r *http.Request) {
	var teamInfo models.Team
	if err := chibind.DefaultBind(r, &teamInfo); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		slog.Error(fmt.Sprintf("Error while parsing data: %v", err.Error()))
		return
	}

	team, err := h.service.PostTeamAdd(teamInfo)
	if err.Error.Code == models.TEAMEXISTS {
		http.Error(w, "Команда уже существует", http.StatusBadRequest)
		slog.Info(err.Error.Message)
		return
	}

	render.JSON(w, r, map[string]any{
		"team_name": team.TeamName,
		"members":   team.Members,
	})

}
func (h *TeamHandler) GetTeamGet(w http.ResponseWriter, r *http.Request) {
	var teamInfo models.GetTeamGetParams
	if err := chibind.DefaultBind(r, &teamInfo); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		slog.Error(fmt.Sprintf("Erorr while parsing data: %v", err.Error()))
		return
	}

	team, err := h.service.GetTeamGet(teamInfo)
	if err.Error.Code == models.NOTFOUND {
		http.Error(w, err.Error.Message, http.StatusNotFound)
		slog.Info(err.Error.Message)
		return
	}

	render.JSON(w, r, map[string]any{
		"team_name": team.TeamName,
		"members":   team.Members,
	})
}
