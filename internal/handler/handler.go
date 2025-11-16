package handler

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/vyacheslavbytsko/Pull-Requests-Reviewers-Service/internal/api"
	"github.com/vyacheslavbytsko/Pull-Requests-Reviewers-Service/internal/db"
)

type Handler struct {
	db   *db.DB
	rand *rand.Rand
}

func NewHandler(db *db.DB) *Handler {
	return &Handler{
		db:   db,
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
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
	ctx := r.Context()
	var body api.PostPullRequestCreateJSONBody

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, api.INVALIDREQUEST, "invalid request body", http.StatusBadRequest)
		return
	}

	if body.PullRequestId == "" {
		writeError(w, api.INVALIDREQUEST, "pull_request_id is required", http.StatusBadRequest)
		return
	}

	if body.PullRequestName == "" {
		writeError(w, api.INVALIDREQUEST, "pull_request_name is required", http.StatusBadRequest)
		return
	}

	if body.AuthorId == "" {
		writeError(w, api.INVALIDREQUEST, "author_id is required", http.StatusBadRequest)
		return
	}

	var exists bool
	err := h.db.Pool.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM prs WHERE pull_request_id=$1)",
		body.PullRequestId,
	).Scan(&exists)

	if err != nil {
		writeError(w, api.INTERNALERROR, "database error", http.StatusInternalServerError)
		return
	}

	if exists {
		writeError(w, api.PREXISTS, "PR id already exists", http.StatusConflict)
		return
	}

	err = h.db.Pool.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM users WHERE user_id=$1)",
		body.AuthorId,
	).Scan(&exists)

	if err != nil {
		writeError(w, api.INTERNALERROR, "database error", http.StatusInternalServerError)
		return
	}

	if !exists {
		writeError(w, api.NOTFOUND, "author not found", http.StatusNotFound)
		return
	}

	var teamName string
	err = h.db.Pool.QueryRow(ctx, "SELECT team_name FROM users WHERE user_id=$1", body.AuthorId).Scan(&teamName)

	if err != nil {
		writeError(w, api.INTERNALERROR, "failed to get author's team", http.StatusInternalServerError)
		return
	}

	// Get rows where all conditions are met
	rows, err := h.db.Pool.Query(ctx, `SELECT user_id FROM users WHERE team_name=$1 AND is_active=TRUE AND user_id<>$2`, teamName, body.AuthorId)

	if err != nil {
		writeError(w, api.INTERNALERROR, "failed to get team members", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	reviewerCandidates := make([]string, 0, 2)

	for rows.Next() {
		var userId string
		if err := rows.Scan(&userId); err != nil {
			writeError(w, api.INTERNALERROR, "failed to scan reviewer", http.StatusInternalServerError)
			return
		}
		reviewerCandidates = append(reviewerCandidates, userId)
	}

	// We shuffle candidates
	if len(reviewerCandidates) > 1 {
		h.rand.Shuffle(len(reviewerCandidates), func(i, j int) {
			reviewerCandidates[i], reviewerCandidates[j] = reviewerCandidates[j], reviewerCandidates[i]
		})
	}

	assignedReviewers := reviewerCandidates
	if len(assignedReviewers) > 2 {
		assignedReviewers = assignedReviewers[:2]
	}

	var createdAt *time.Time
	err = h.db.Pool.QueryRow(ctx, `
			INSERT INTO prs(pull_request_id, pull_request_name, author_id, status, assigned_reviewers)
			VALUES($1, $2, $3, $4, $5)
			RETURNING created_at
		`, body.PullRequestId, body.PullRequestName, body.AuthorId, "OPEN", assignedReviewers).Scan(&createdAt)

	if err != nil {
		writeError(w, api.INTERNALERROR, "failed to create PR", http.StatusInternalServerError)
		return
	}

	pr := api.PullRequest{
		PullRequestId:     body.PullRequestId,
		PullRequestName:   body.PullRequestName,
		AuthorId:          body.AuthorId,
		Status:            api.PullRequestStatusOPEN,
		AssignedReviewers: assignedReviewers,
		CreatedAt:         createdAt,
		MergedAt:          nil,
	}

	writeJSON(w, http.StatusCreated, map[string]api.PullRequest{"pr": pr})
}

func (h *Handler) PostPullRequestMerge(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) PostPullRequestReassign(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) PostTeamAdd(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
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

	var exists bool
	err := h.db.Pool.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM teams WHERE team_name=$1)",
		team.TeamName,
	).Scan(&exists)

	if err != nil {
		writeError(w, api.INTERNALERROR, "database error", http.StatusInternalServerError)
		return
	}

	if exists {
		writeError(w, api.TEAMEXISTS, "team_name already exists", http.StatusBadRequest)
		return
	}

	_, err = h.db.Pool.Exec(ctx, "INSERT INTO teams(team_name) VALUES($1)", team.TeamName)

	if err != nil {
		writeError(w, api.INTERNALERROR, "failed to create team", http.StatusInternalServerError)
		return
	}

	for _, m := range team.Members {
		_, err := h.db.Pool.Exec(ctx, `
			INSERT INTO users(user_id, username, team_name, is_active)
			VALUES($1, $2, $3, $4)
			ON CONFLICT (user_id) DO UPDATE
			SET username = EXCLUDED.username,
			    team_name = EXCLUDED.team_name,
			    is_active = EXCLUDED.is_active
		`, m.UserId, m.Username, team.TeamName, m.IsActive)
		if err != nil {
			writeError(w, api.INTERNALERROR, fmt.Sprintf("failed to insert user %s", m.UserId), http.StatusInternalServerError)
			return
		}
	}

	writeJSON(w, http.StatusCreated, team)
}

func (h *Handler) GetTeamGet(w http.ResponseWriter, r *http.Request, params api.GetTeamGetParams) {
	ctx := r.Context()
	teamName := params.TeamName

	var exists bool
	err := h.db.Pool.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM teams WHERE team_name=$1)",
		teamName,
	).Scan(&exists)

	if err != nil {
		writeError(w, api.INTERNALERROR, "database error", http.StatusInternalServerError)
		return
	}

	if !exists {
		writeError(w, api.NOTFOUND, "team not found", http.StatusNotFound)
		return
	}

	rows, err := h.db.Pool.Query(ctx, `SELECT user_id, username, is_active FROM users WHERE team_name=$1`, teamName)

	if err != nil {
		writeError(w, api.INTERNALERROR, "failed to fetch users", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var members []api.TeamMember
	for rows.Next() {
		var m api.TeamMember
		if err := rows.Scan(&m.UserId, &m.Username, &m.IsActive); err != nil {
			writeError(w, api.INTERNALERROR, "failed to fetch member", http.StatusInternalServerError)
			return
		}
		members = append(members, m)
	}

	team := api.Team{
		TeamName: teamName,
		Members:  members,
	}

	writeJSON(w, http.StatusOK, team)
}

func (h *Handler) GetUsersGetReview(w http.ResponseWriter, r *http.Request, params api.GetUsersGetReviewParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) PostUsersSetIsActive(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var body api.PostUsersSetIsActiveJSONBody

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, api.INVALIDREQUEST, "invalid request body", http.StatusBadRequest)
		return
	}

	if body.UserId == "" {
		writeError(w, api.INVALIDREQUEST, "user_id is required", http.StatusBadRequest)
		return
	}

	// TODO: check if is_active is in body and not just "false" by default

	cmdTag, err := h.db.Pool.Exec(ctx, `
		UPDATE users 
		SET is_active=$1
		WHERE user_id=$2
	`, body.IsActive, body.UserId)
	if err != nil {
		writeError(w, api.INTERNALERROR, "failed to update user", http.StatusInternalServerError)
		return
	}

	if cmdTag.RowsAffected() == 0 {
		writeError(w, api.NOTFOUND, "user not found", http.StatusNotFound)
		return
	}

	var user api.User
	err = h.db.Pool.QueryRow(ctx, `
		SELECT user_id, username, team_name, is_active 
		FROM users 
		WHERE user_id=$1
	`, body.UserId).Scan(&user.UserId, &user.Username, &user.TeamName, &user.IsActive)
	if err != nil {
		writeError(w, api.INTERNALERROR, "failed to fetch user", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, user)
}
