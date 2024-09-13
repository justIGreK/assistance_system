package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/o1egl/paseto"
)

var pasetoInstance = paseto.NewV2()

func GeneratePasetoToken(userID int) (string, error) {
	symmetricKey  := []byte(os.Getenv("SYMMETRIC_KEY"))
	token := paseto.JSONToken{
		Expiration: time.Now().Add(24 * time.Hour),
		Subject:    fmt.Sprintf("%d", userID),
	}
	encrypted, err := pasetoInstance.Encrypt(symmetricKey, token, nil)
	return encrypted, err
}

func ValidatePasetoToken(tokenString string) (int, error) {
	symmetricKey  := []byte(os.Getenv("SYMMETRIC_KEY"))
	var token paseto.JSONToken
	var footer string
	err := pasetoInstance.Decrypt(tokenString, symmetricKey, &token, &footer)
	if err != nil {
		return 0, err
	}

	if time.Now().After(token.Expiration) {
		return 0, fmt.Errorf("token expired")
	}

	var userID int
	fmt.Sscanf(token.Subject, "%d", &userID)
	return userID, nil
}
