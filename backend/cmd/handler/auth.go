package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"gohelp/internal/models"
	"log"
	"net/http"
)

type Authorization interface {
	RegisterUser(ctx context.Context, user models.SignUp) error
	LoginUser(ctx context.Context, email, password string) (string, error)
}

// @Summary SignUp
// @Tags auth
// @Description create account
// @Accept  json
// @Produce  json
// @Param username query string true "your username"
// @Param password query string true "your password"
// @Param email query string true "your email"
// @Router /users/register [post]
func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	log.Println("signIn func running")
	input := models.SignUp{
		Username: r.URL.Query().Get("username"),
		Email:    r.URL.Query().Get("email"),
		Password: r.URL.Query().Get("password"),
	}
	fmt.Println(input)
	if err := validate.Struct(input); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}
	err := h.RegisterUser(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	log.Println("signUp func ended")
}
// @Summary SignIn
// @Tags auth
// @Description create account
// @Accept  json
// @Produce  json
// @Param email query string true "your email"
// @Param password query string true "your password"
// @Router /users/login [post]
func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	log.Println("signIn func running")
	credentials := models.LoginRequest{
		Email:    r.URL.Query().Get("email"),
		Password: r.URL.Query().Get("password"),
	}
	if err := validate.Struct(credentials); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}
	token, err := h.LoginUser(r.Context(), credentials.Email, credentials.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
	log.Println("signIn func ended")
}
