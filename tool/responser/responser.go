package responser

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/isOdin-l/Assigning-PR-reviewers/internal/models"
)

// func Rend(w http.ResponseWriter, r *http.Request, err *models.ErrorResponse, response any){
// 	switch err.Code {
// 	case api.NOTFOUND:

// 	}
// }

func RenderResponse(w http.ResponseWriter, r *http.Request, status int, response any) {
	render.Status(r, status)
	render.JSON(w, r, response)
}

func RenderError(w http.ResponseWriter, r *http.Request, status int, err *models.ErrorResponse) {
	render.Status(r, status)
	render.JSON(w, r, models.ConvertToApiErrorResponse(err))
}
