package handler

import "net/http"

type TeamServiceInterface interface {
}

type TeamHandler struct {
	service TeamServiceInterface
}

func NewTeamHandler(service TeamServiceInterface) *TeamHandler {
	return &TeamHandler{service: service}
}

func (h *TeamHandler) PostTeamAdd(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
func (h *TeamHandler) GetTeamGet(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
