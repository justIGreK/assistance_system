package handler

import (
	"gohelp/internal/service/auth"
	"gohelp/internal/service/forum"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "gohelp/docs"
)

type Handler struct {
	Users
	Forum
}

func NewHandler(user *auth.UserService, forum *forum.ForumService) *Handler {
	return &Handler{Users: user, Forum: forum}
}

func (h *Handler) InitRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Get("/discussions", h.GetDiscussionsWithCountOfComments)
	r.Get("/search", h.SearchDiscussionsByName)
	r.Get("/getdiscussion", h.GetDiscussionWithComments)
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.SignUp)
		r.Post("/login", h.SignIn)
	})
	r.Route("/users", func(r chi.Router) {
		r.Use(h.AuthMiddleware)
		r.Put("/actions", h.UsersActions)
	})
	r.Route("/discuss", func(r chi.Router) {
		r.Use(h.AuthMiddleware)
		r.Post("/discussions", h.CreateDiscussion)
		r.Post("/comments", h.CreateComment)
		r.Post("/vote", h.Vote)
		r.Put("/discussions/edit", h.UpdateDiscussion)
		r.Put("/comments/edit", h.UpdateComment)
		r.Delete("/discussions/delete", h.DeleteDiscussion)
		r.Delete("/comments/delete", h.DeleteComment)
	})

	return r
}
