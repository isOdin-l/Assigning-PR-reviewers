package handler

import "net/http"

type UserServiceInterface interface {
}

type UserHandler struct {
	service UserServiceInterface
}

func NewUserHandler(service UserServiceInterface) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetUsersGetReview(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
func (h *UserHandler) PostUsersSetIsActive(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
