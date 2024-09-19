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

// @Summary Create New Discussion
// @Security BearerAuth
// @Tags discussions
// @Description You can post new discussion
// @Accept  json
// @Produce  json
// @Param title query string true "Title of discussion"
// @Param content query string true "Describe your problem here"
// @Router /discuss/discussions [post]
func (h *Handler) CreateDiscussion(w http.ResponseWriter, r *http.Request) {
	log.Println("CreateDisc func running")
	request := struct {
		Title   string `json:"title" validate:"required"`
		Content string `json:"content" validate:"required"`
	}{
		Title:   r.URL.Query().Get("title"),
		Content: r.URL.Query().Get("content"),
	}

	AuthorID := r.Context().Value(UserIDKey).(int)
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

// @Summary Comment discussion
// @Security BearerAuth
// @Tags discussions
// @Description You can comment a discussion
// @Accept  json
// @Produce  json
// @Param discussionID query string true "Id of discussion"
// @Param content query string true "Your comment"
// @Router /discuss/comments [post]
func (h *Handler) CreateComment(w http.ResponseWriter, r *http.Request) {
	log.Println("CreateCom func running")
	request := struct {
		DiscussionID string `json:"discussionID" validate:"required"`
		Content      string `json:"content" validate:"required"`
	}{
		DiscussionID: r.URL.Query().Get("discussionID"),
		Content:      r.URL.Query().Get("content"),
	}
	AuthorID := r.Context().Value(UserIDKey).(int)
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

// @Summary Get all discussions
// @Tags discussions
// @Description Get all discussions on site
// @Accept  json
// @Produce  json
// @Router /discussions [get]
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

// @Summary Get all discussions
// @Tags discussions
// @Description Get all discussions on site
// @Accept  json
// @Produce  json
// @Param discussionName query string true "Search term"
// @Router /search [get]
func (h *Handler) SearchDiscussionsByName(w http.ResponseWriter, r *http.Request) {
	log.Println("SearchDiscByName func running")
	request := struct {
		DiscussionName string `json:"discussionName" validate:"required"`
	}{
		DiscussionName: r.URL.Query().Get("discussionName"),
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

// @Summary Submit a vote
// @Security BearerAuth
// @Tags discussions
// @Description Submit a vote with either "like" or "dislike"
// @Accept  json
// @Produce  json
// @Param ElementId query string true "Id of discussion or comment"
// @Param vote query string true "The type of vote. Can be either 'like' or 'dislike'."
// @Router /discuss/vote [post]
func (h *Handler) Vote(w http.ResponseWriter, r *http.Request) {
	AuthorID := r.Context().Value(UserIDKey).(int)

	request := struct {
		ElementId string `json:"ElementId" validate:"required"`
		VoteType  string `json:"vote" validate:"required"`
	}{
		ElementId: r.URL.Query().Get("ElementId"),
		VoteType:  r.URL.Query().Get("vote"),
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

// @Summary Get full discussion
// @Security BearerAuth
// @Tags discussions
// @Description Get full display of discussion with comments
// @Accept  json
// @Produce  json
// @Param discussion_id query string true "Id of discussion"
// @Router /getdiscussion [get]
func (h *Handler) GetDiscussionWithComments(w http.ResponseWriter, r *http.Request) {

	request := struct {
		DiscussionId string `json:"discussion_id" validate:"required"`
	}{
		DiscussionId: r.URL.Query().Get("discussion_id"),
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
	w.WriteHeader(http.StatusOK)
}
