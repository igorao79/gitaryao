package api

import (
	"encoding/json"
	"net/http"
	"regexp"

	"gitserv/internal/auth"
	"gitserv/internal/gitops"
	"gitserv/internal/models"
)

var validNameRe = regexp.MustCompile(`^[a-zA-Z0-9_\-\.]+$`)

type createRepoRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsPrivate   bool   `json:"is_private"`
}

type errorResponse struct {
	Error string `json:"error"`
}

// CreateRepo handles POST /api/repos (authenticated)
func (s *Server) CreateRepo(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetClaimsFromContext(r.Context())
	if claims == nil {
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "Authentication required"})
		return
	}

	var req createRepoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Invalid JSON"})
		return
	}

	if req.Name == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Name is required"})
		return
	}

	if !validNameRe.MatchString(req.Name) {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "Invalid repo name"})
		return
	}

	owner := claims.Username

	if gitops.RepoExists(s.ReposDir, owner, req.Name) {
		writeJSON(w, http.StatusConflict, errorResponse{Error: "Repository already exists"})
		return
	}

	if err := gitops.InitBareRepo(s.ReposDir, owner, req.Name); err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "Failed to create repository"})
		return
	}

	repo, err := s.Repos.Create(claims.UserID, req.Name, req.Description, req.IsPrivate)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "Failed to save repository"})
		return
	}

	writeJSON(w, http.StatusCreated, repo)
}

// ListMyRepos handles GET /api/repos (authenticated — user's repos)
func (s *Server) ListMyRepos(w http.ResponseWriter, r *http.Request) {
	claims := auth.GetClaimsFromContext(r.Context())
	if claims == nil {
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "Authentication required"})
		return
	}

	repos, err := s.Repos.ListByOwnerID(claims.UserID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "Failed to list repos"})
		return
	}

	if repos == nil {
		repos = []models.Repository{}
	}

	writeJSON(w, http.StatusOK, repos)
}

// ListPublicRepos handles GET /api/repos/public (no auth needed)
func (s *Server) ListPublicRepos(w http.ResponseWriter, r *http.Request) {
	repos, err := s.Repos.ListPublic()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "Failed to list repos"})
		return
	}

	if repos == nil {
		repos = []models.Repository{}
	}

	writeJSON(w, http.StatusOK, repos)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
