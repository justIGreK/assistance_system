package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"gohelp/internal/models"
	"gohelp/util"
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
	UpdateDiscussion(ctx context.Context, discussionID, content string, authorID int) (*models.Discussion, error)
	UpdateComment(ctx context.Context, commentID, content string, authorID int) (*models.Comment, error)
	DeleteDiscussion(ctx context.Context, commentID string, authorID int) error 
	DeleteComment(ctx context.Context, commentID string, authorID int) error 
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
		Title   string `json:"title" validate:"required,max=35"`
		Content string `json:"content" validate:"required,max=70"`
	}{
		Title:   r.URL.Query().Get("title"),
		Content: r.URL.Query().Get("content"),
	}

	AuthorID := r.Context().Value(UserIDKey).(int)
	if err := validate.Struct(request); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}
	if err := util.ValidateTitle(request.Title); err != nil {
		http.Error(w, "Invalid title: "+err.Error(), http.StatusBadRequest)
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
		Content      string `json:"content" validate:"required,max=70"`
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
		DiscussionName string `json:"discussionName" validate:"required,max=35"`
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
// @Param vote query string true "The type of vote. Can be either 'like' or 'dislike'." Enums(like, dislike)
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
}

// @Summary Update discussion
// @Security BearerAuth
// @Tags discussions
// @Accept  json
// @Produce  json
// @Param discussion_id query string true "Id of discussion"
// @Param content query string true "New content field"
// @Router /discuss/discussions/edit [put]
func (h *Handler) UpdateDiscussion(w http.ResponseWriter, r *http.Request) {
	AuthorID := r.Context().Value(UserIDKey).(int)
	request := struct {
		DiscussionID string `json:"discussion_id" validate:"required"`
		Content      string `json:"content" validate:"max=70"`
	}{
		DiscussionID: r.URL.Query().Get("discussion_id"),
		Content:      r.URL.Query().Get("content"),
	}
	if err := validate.Struct(request); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}
	discussion, err := h.Forum.UpdateDiscussion(r.Context(), request.DiscussionID, request.Content, AuthorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"Updated discussion": discussion,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// @Summary Update comment
// @Security BearerAuth
// @Tags discussions
// @Accept  json
// @Produce  json
// @Param comment_id query string true "Id of comment"
// @Param content query string true "New content field"
// @Router /discuss/comments/edit [put]
func (h *Handler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	AuthorID := r.Context().Value(UserIDKey).(int)
	request := struct {
		CommentID string `json:"comment_id" validate:"required"`
		Content   string `json:"content" validate:"max=70"`
	}{
		CommentID: r.URL.Query().Get("comment_id"),
		Content:   r.URL.Query().Get("content"),
	}
	if request.Content == "" {
		http.Error(w, "Content field cant be empty", http.StatusBadRequest)
		return
	}
	if err := validate.Struct(request); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}
	comment, err := h.Forum.UpdateComment(r.Context(), request.CommentID, request.Content, AuthorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"Updated comment": comment,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// // @Summary Update discussion
// // @Security BearerAuth
// // @Tags discussions
// // @Accept  json
// // @Produce  json
// // @Param discussion_id query string true "Id of discussion"
// // @Router /discuss/discussions/delete [delete]
// func (h *Handler) DeleteDiscussion(w http.ResponseWriter, r *http.Request) {
// 	AuthorID := r.Context().Value(UserIDKey).(int)
// 	request := struct {
// 		DiscussionID string `json:"discussion_id" validate:"required"`
// 	}{
// 		DiscussionID: r.URL.Query().Get("discussion_id"),
// 	}
// 	if err := validate.Struct(request); err != nil {
// 		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	err := h.Forum.DeleteDiscussion(r.Context(), request.DiscussionID, AuthorID)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// }

// @Summary Delete comment
// @Security BearerAuth
// @Tags discussions
// @Accept  json
// @Produce  json
// @Param comment_id query string true "Id of comment"
// @Router /discuss/comments/delete [delete]
func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	AuthorID := r.Context().Value(UserIDKey).(int)
	request := struct {
		CommentID string `json:"comment_id" validate:"required"`
	}{
		CommentID: r.URL.Query().Get("comment_id"),
	}
	if err := validate.Struct(request); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}
	err := h.Forum.DeleteComment(r.Context(), request.CommentID, AuthorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
