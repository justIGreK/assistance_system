package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"gohelp/internal/models"
	"gohelp/pkg"
	"log"
	"net/http"
	"strconv"

	"github.com/markbates/goth"
)

type Users interface {
	RegisterUser(ctx context.Context, user models.SignUp) error
	LoginUser(ctx context.Context, email, password string) (string, error)
	UsersActions(ctx context.Context, userID int, action string) (*models.User, error)
	GoogleAuth(ctx context.Context, user goth.User) (string, error)
}

// @Summary SignUp
// @Tags users
// @Description create account
// @Accept  json
// @Produce  json
// @Param username query string true "your username"
// @Param password query string true "your password"
// @Param email query string true "your email"
// @Router /auth/register [post]
func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	log.Println("signUP func running")
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

const (
	googleMethod  = "google"
	classicMethod = "classic"
)

// @Summary SignIn
// @Tags users
// @Description create account
// @Accept  json
// @Produce  json
// @Param auth_method query string true "Authorization method" Enums(classic, google)
// @Param email query string false "your email"
// @Param password query string false "your password"
// @Router /auth/login [post]
func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	if authMethod := r.URL.Query().Get("auth_method"); authMethod != classicMethod {
		response := fmt.Sprintf("For this action, please follow this link: http://localhost:8080/auth/%v", authMethod)
		json.NewEncoder(w).Encode(response)
		return
	}
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

func (h *Handler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	pkg.StartGoogleOAuth(w, r)
}

func (h *Handler) GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	user, err := pkg.CompleteGoogleOAuth(w, r)
	if err != nil {
		http.Error(w, "Failed to authenticate with Google", http.StatusUnauthorized)
		return
	}
	// fmt.Println("Email: " + user.Email +
	// 	"\n Name: " + user.Name +
	// 	"\n FirstName: " + user.FirstName +
	// 	"\n LastName: " + user.LastName +
	// 	"\n Description: " + user.Description +
	// 	"\n UserID: " + user.UserID +
	// 	"\n AvatarURL: " + user.AvatarURL +
	// 	"\n Location: " + user.Location +
	// 	"\n AccessToken: " + user.AccessToken)
	token, err := h.GoogleAuth(r.Context(), user)
	if err != nil{
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(token)
}

// @Summary Change status of user
// @Security BearerAuth
// @Tags users
// @Accept  json
// @Produce  json
// @Param user_id query string true "Id of User"
// @Param action query string true "The type of action. Can be either 'ban' or 'unban'." Enums(ban, unban)
// @Router /users/actions [put]
func (h *Handler) UsersActions(w http.ResponseWriter, r *http.Request) {
	UserRole := r.Context().Value(UserRoleKey).(string)
	if UserRole != models.AdministrationRole {
		http.Error(w, "You dont have permisions to do this", http.StatusUnauthorized)
		return
	}
	strID := r.URL.Query().Get("user_id")
	id, err := strconv.Atoi(strID)
	if err != nil {
		http.Error(w, "Invalid 'id' parameter", http.StatusBadRequest)
		return
	}
	request := struct {
		UserID int    `json:"user_id" validate:"required"`
		Action string `json:"action" validate:"required"`
	}{
		UserID: id,
		Action: r.URL.Query().Get("action"),
	}
	if err := validate.Struct(request); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}
	_, err = h.Users.UsersActions(r.Context(), request.UserID, request.Action)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if request.Action == "ban" {
		err = h.Forum.DeleteFullHistory(r.Context(), request.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("operation is completed")
}
