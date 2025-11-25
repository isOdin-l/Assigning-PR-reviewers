package handler

import (
	"context"
	"net/http"

	"github.com/isOdin-l/Assigning-PR-reviewers/internal/models"
	"github.com/isOdin-l/Assigning-PR-reviewers/pkg/api"
	"github.com/isOdin-l/Assigning-PR-reviewers/tool/chibind"
	"github.com/isOdin-l/Assigning-PR-reviewers/tool/responser"
)

type TeamServiceInterface interface {
	PostTeamAdd(ctx context.Context, team *models.Team) (*api.ResponseTeam, *models.ErrorResponse)
	GetTeam(ctx context.Context, team *models.GetTeamParams) (*api.ResponseTeam, *models.ErrorResponse)
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
		responser.HandleError(w, r, &models.ErrorResponse{Code: api.SERVERERROR, Message: err.Error()})
		return
	}

	response, err := h.service.PostTeamAdd(r.Context(), models.ConvertToTeam(&team))
	if err != nil {
		responser.HandleError(w, r, err)
		return
	}

	responser.RenderResponse(w, r, http.StatusCreated, *response)
}
func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	var team api.GetTeamParams
	if err := chibind.DefaultBind(r, &team); err != nil {
		responser.HandleError(w, r, &models.ErrorResponse{Code: api.SERVERERROR, Message: err.Error()})
		return
	}

	response, err := h.service.GetTeam(r.Context(), models.ConvertToGetTeamParams(&team))
	if err != nil {
		responser.HandleError(w, r, err)
		return
	}

	responser.RenderResponse(w, r, http.StatusOK, *response)
}
