package handler

import (
	"context"
	"gohelp/internal/service/auth"
	"net/http"
	"strings"
)

type contextKey string
const( 
	UserIDKey contextKey = "user_id"
	UserRoleKey contextKey = "user_role"
)

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		payload, err := auth.ValidatePasetoToken(token)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserRoleKey, payload.Role)
		ctx = context.WithValue(ctx, UserIDKey, payload.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
