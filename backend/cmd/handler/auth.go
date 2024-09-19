package handler

import (
	"context"
	"encoding/json"
	"gohelp/internal/models"
	"log"
	"net/http"
)

type Authorization interface {
	RegisterUser(ctx context.Context, user models.SignUp) error
	LoginUser(ctx context.Context, email, password string) (string, error)
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var u models.SignUp
	log.Println("signUp func running")
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		log.Println(r.Body)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	if err := validate.Struct(u); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}
	err := h.RegisterUser(r.Context(), u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	log.Println("signUp func ended")
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	log.Println("signIn func running")
	var credentials models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
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
