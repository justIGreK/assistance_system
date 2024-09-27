package pkg

import (
	"net/http"
	"os"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func InitOAuth() {
	goth.UseProviders(
		google.New(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"), "http://localhost:8080/auth/google/callback"),
	)
}

// StartGoogleOAuth initiates Google OAuth flow
func StartGoogleOAuth(w http.ResponseWriter, r *http.Request) {
	if _, err := gothic.GetProviderName(r); err != nil {
		provider := r.URL.Path[len("/auth/"):]
		if provider == "" {
			http.Error(w, "Provider not selected", http.StatusBadRequest)
			return
		}
		r.URL.RawQuery = "provider=" + provider
	}
	gothic.BeginAuthHandler(w, r)
}

// CompleteGoogleOAuth completes the Google OAuth flow
func CompleteGoogleOAuth(w http.ResponseWriter, r *http.Request) (goth.User, error) {
	user, err := gothic.CompleteUserAuth(w, r)
	return user, err
}
