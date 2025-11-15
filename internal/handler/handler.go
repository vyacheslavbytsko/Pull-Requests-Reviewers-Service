package handler

import (
	"net/http"

	"github.com/vyacheslavbytsko/Pull-Requests-Reviewers-Service/internal/api"
	"github.com/vyacheslavbytsko/Pull-Requests-Reviewers-Service/internal/store"
)

type Handler struct {
	store *store.Store
}

func NewHandler() *Handler {
	return &Handler{
		store: store.NewStore(),
	}
}

func (h *Handler) PostPullRequestCreate(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) PostPullRequestMerge(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) PostPullRequestReassign(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
func (h *Handler) PostTeamAdd(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) GetTeamGet(w http.ResponseWriter, r *http.Request, params api.GetTeamGetParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) GetUsersGetReview(w http.ResponseWriter, r *http.Request, params api.GetUsersGetReviewParams) {
	w.Write([]byte("Hello World from PostPullRequestCreate!"))
}

func (h *Handler) PostUsersSetIsActive(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
