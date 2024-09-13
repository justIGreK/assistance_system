package handler

import (
	"context"
	"encoding/json"
	"gohelp/internal/models"
	"net/http"
)

type Forum interface {
	CreateDiscussion(ctx context.Context, title, content string, AuthorID int) (string, error)
	CreateComment(ctx context.Context, discussionID, content string, AuthorID int) (string, error)
	GetDiscussionWithComments(ctx context.Context, discussionID string) (*models.Discussion, []models.Comment, error)
	GetAllDiscussionsWithCountOfComments(ctx context.Context) ([]models.DiscussionWithCount, error)
	SearchDiscussionsByName(ctx context.Context, searchTerm string) ([]models.Discussion, error) 
}

func (h *Handler) CreateDiscussion(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}

	AuthorID := r.Context().Value(UserIDKey).(int)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	id, err := h.Forum.CreateDiscussion(r.Context(), request.Title, request.Content, AuthorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func (h *Handler) CreateComment(w http.ResponseWriter, r *http.Request) {
	var request struct {
		DiscussionID string `json:"discussion_id" binding:"required"`
		Content      string `json:"content" binding:"required"`
	}
	AuthorID := r.Context().Value(UserIDKey).(int)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	id, err := h.Forum.CreateComment(r.Context(), request.DiscussionID, request.Content, AuthorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func (h *Handler) GetDiscussionsWithComments(w http.ResponseWriter, r *http.Request) {
	discussion, err := h.Forum.GetAllDiscussionsWithCountOfComments(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"discussion": discussion,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) SearchDiscussionsByName(w http.ResponseWriter, r *http.Request){
	var request struct {
		DiscussionName string `json:"discussionName" binding:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	discussions, err := h.Forum.SearchDiscussionsByName(r.Context(), request.DiscussionName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"discussions": discussions,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}