package handler

import "net/http"

type PullRequestServiceInterface interface {
}

type PullRequestHandler struct {
	service PullRequestServiceInterface
}

func NewPullRequestHandler(service PullRequestServiceInterface) *PullRequestHandler {
	return &PullRequestHandler{service: service}
}

func (h *PullRequestHandler) PostPullRequestCreate(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
func (h *PullRequestHandler) PostPullRequestMerge(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
func (h *PullRequestHandler) PostPullRequestReassign(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
