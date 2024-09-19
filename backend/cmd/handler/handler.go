package handler

import (
	"gohelp/internal/service/auth"
	"gohelp/internal/service/forum"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Authorization
	Forum
}

func NewHandler(user *auth.UserService, forum *forum.ForumService) *Handler {
	return &Handler{Authorization: user, Forum: forum}
}

func (h *Handler) InitRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/discussions", h.GetDiscussionsWithCountOfComments)
	r.Get("/search", h.SearchDiscussionsByName)
	r.Get("/getdiscussion", h.GetDiscussionWithComments)
	r.Route("/users", func(r chi.Router) {
		r.Post("/register", h.SignUp)
		r.Post("/login", h.SignIn)
	})
	r.Route("/discuss", func(r chi.Router) {
		r.Use(h.AuthMiddleware)
		r.Post("/discussions", h.CreateDiscussion)
		r.Post("/comments", h.CreateComment)
		r.Post("/vote", h.Vote)
	})

	return r
}
