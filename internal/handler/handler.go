package handler

import (
	"encoding/json"
	"fmt"
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

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, code api.ErrorResponseErrorCode, msg string, status int) {
	resp := api.ErrorResponse{}
	resp.Error.Code = code
	resp.Error.Message = msg
	writeJSON(w, status, resp)
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
	var team api.Team

	if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
		writeError(w, api.INVALIDREQUEST, "invalid request body", http.StatusBadRequest)
		return
	}

	if team.TeamName == "" {
		writeError(w, api.INVALIDREQUEST, "team_name is required", http.StatusBadRequest)
		return
	}

	if len(team.Members) == 0 {
		writeError(w, api.INVALIDREQUEST, "members cannot be empty", http.StatusBadRequest)
		return
	}

	for i, m := range team.Members {
		if m.UserId == "" || m.Username == "" {
			writeError(w, api.INVALIDREQUEST, fmt.Sprintf("member at index %d is invalid", i), http.StatusBadRequest)
			return
		}
	}

	h.store.Mu.Lock()
	defer h.store.Mu.Unlock()

	if _, exists := h.store.Teams[team.TeamName]; exists {
		writeError(w, api.TEAMEXISTS, "team_name already exists", http.StatusBadRequest)
		return
	}

	h.store.Teams[team.TeamName] = &api.Team{
		TeamName: team.TeamName,
		Members:  team.Members,
	}

	for _, m := range team.Members {
		if user, exists := h.store.Users[m.UserId]; exists {
			user.UserId = m.UserId
			user.Username = m.Username
			user.IsActive = m.IsActive
			user.TeamName = team.TeamName
		} else {
			h.store.Users[m.UserId] = &api.User{
				UserId:   m.UserId,
				Username: m.Username,
				IsActive: m.IsActive,
				TeamName: team.TeamName,
			}
		}
	}

	writeJSON(w, http.StatusCreated, team)
}

func (h *Handler) GetTeamGet(w http.ResponseWriter, r *http.Request, params api.GetTeamGetParams) {
	teamName := params.TeamName
	team, exists := h.store.Teams[teamName]

	if !exists {
		writeError(w, "NOT_FOUND", "team not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusCreated, team)
}

func (h *Handler) GetUsersGetReview(w http.ResponseWriter, r *http.Request, params api.GetUsersGetReviewParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) PostUsersSetIsActive(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
