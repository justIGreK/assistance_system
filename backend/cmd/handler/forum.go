package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"gohelp/internal/models"
	"log"
	"net/http"

	"github.com/go-playground/validator"
)

type Forum interface {
	CreateDiscussion(ctx context.Context, title, content string, AuthorID int) (string, error)
	CreateComment(ctx context.Context, discussionID, content string, AuthorID int) (string, error)
	GetDiscussionWithComments(ctx context.Context, discussionID string) (*models.Discussion, []models.Comment, error)
	GetAllDiscussionsWithCountOfComments(ctx context.Context) ([]models.DiscussionWithCount, error)
	SearchDiscussionsByName(ctx context.Context, searchTerm string) ([]models.Discussion, error)
	Vote(ctx context.Context, userID int, discussionID, voteType string) error
}

var validate = validator.New()

func (h *Handler) CreateDiscussion(w http.ResponseWriter, r *http.Request) {
	log.Println("CreateDisc func running")
	var request struct {
		Title   string `json:"title" validate:"required"`
		Content string `json:"content" validate:"required"`
	}

	AuthorID := r.Context().Value(UserIDKey).(int)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if err := validate.Struct(request); err != nil {
        http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
        return
    }

	id, err := h.Forum.CreateDiscussion(r.Context(), request.Title, request.Content, AuthorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id})
	log.Println("CreateDisc func ended")
}

func (h *Handler) CreateComment(w http.ResponseWriter, r *http.Request) {
	log.Println("CreateCom func running")
	var request struct {
		DiscussionID string `json:"discussion_id" validate:"required"`
		Content      string `json:"content" validate:"required"`
	}
	AuthorID := r.Context().Value(UserIDKey).(int)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if err := validate.Struct(request); err != nil {
        http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
        return
    }

	id, err := h.Forum.CreateComment(r.Context(), request.DiscussionID, request.Content, AuthorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id})
	log.Println("CreateCom func ended")
}

func (h *Handler) GetDiscussionsWithCountOfComments(w http.ResponseWriter, r *http.Request) {
	log.Println("GetDiscWithCom func running")
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
	log.Println("GetDiscWithCom func ended")
}

func (h *Handler) SearchDiscussionsByName(w http.ResponseWriter, r *http.Request) {
	log.Println("SearchDiscByName func running")
	var request struct {
		DiscussionName string `json:"discussionName" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if err := validate.Struct(request); err != nil {
        http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
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
	log.Println("SearchDiscByName func ended")
}

func (h *Handler) Vote(w http.ResponseWriter, r *http.Request) {
	AuthorID := r.Context().Value(UserIDKey).(int)

	var request struct {
		ElementId string `json:"element_id" validate:"required"`
		VoteType     string `json:"vote_type" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if err := validate.Struct(request); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}
	err := h.Forum.Vote(r.Context(), AuthorID, request.ElementId, request.VoteType)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to vote: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetDiscussionWithComments(w http.ResponseWriter, r *http.Request) {
	var request struct {
		DiscussionId string `json:"discussion_id" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if err := validate.Struct(request); err != nil {
        http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
        return
    }
	fmt.Println(request.DiscussionId)
	discussion, comments, err := h.Forum.GetDiscussionWithComments(r.Context(), request.DiscussionId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"discussion": discussion,
		"comments":   comments,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
