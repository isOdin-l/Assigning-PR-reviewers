package responser

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/isOdin-l/Assigning-PR-reviewers/internal/models"
	"github.com/isOdin-l/Assigning-PR-reviewers/pkg/api"
)

func RenderResponse(w http.ResponseWriter, r *http.Request, status int, response any) {
	render.Status(r, status)
	render.JSON(w, r, response)
}

func renderError(w http.ResponseWriter, r *http.Request, status int, err *models.ErrorResponse) {
	render.Status(r, status)
	render.JSON(w, r, models.ConvertToApiErrorResponse(err))
}

func HandleError(w http.ResponseWriter, r *http.Request, err *models.ErrorResponse) {
	switch err.Code {
	case api.TEAMEXISTS:
		slog.Info(fmt.Sprintf("Handler layer: %s", err.Message))
		renderError(w, r, http.StatusBadRequest, err)
	case api.NOTFOUND:
		slog.Info(fmt.Sprintf("Handler layer: %s", err.Message))
		renderError(w, r, http.StatusNotFound, err)
	case api.PRMERGED:
		slog.Info(fmt.Sprintf("Handler layer: %s", err.Message))
		renderError(w, r, http.StatusConflict, err)
	case api.NOTASSIGNED:
		slog.Info(fmt.Sprintf("Handler layer: %s", err.Message))
		renderError(w, r, http.StatusConflict, err)
	case api.NOCANDIDATE:
		slog.Info(fmt.Sprintf("Handler layer: %s", err.Message))
		renderError(w, r, http.StatusConflict, err)
	case api.PREXISTS:
		slog.Info(fmt.Sprintf("Handler layer: %s", err.Message))
		renderError(w, r, http.StatusConflict, err)
	case api.SERVERERROR:
		slog.Error(fmt.Sprintf("Handler layer: %s", err.Message))
		renderError(w, r, http.StatusInternalServerError, err)
	}
}
