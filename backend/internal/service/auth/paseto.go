package auth

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/o1egl/paseto"
)

var pasetoInstance = paseto.NewV2()

type TokenPayload struct {
	UserID     int       `json:"user_id"`
	Role       string    `json:"user_role"`
	Expiration time.Time `json:"expiration"`
}

func GeneratePasetoToken(userID int, userRole string) (string, error) {
	symmetricKey := []byte(os.Getenv("SYMMETRIC_KEY"))
	payload := TokenPayload{
		UserID:     userID,
		Role:       userRole,
		Expiration: time.Now().Add(24 * time.Hour),
	}
	
	encrypted, err := pasetoInstance.Encrypt(symmetricKey, payload, nil)
	return encrypted, err
}

func ValidatePasetoToken(tokenString string) (*TokenPayload, error) {
	symmetricKey := []byte(os.Getenv("SYMMETRIC_KEY"))
	var payload TokenPayload
	var footer string
	err := pasetoInstance.Decrypt(tokenString, symmetricKey, &payload, &footer)
	if err != nil {
		return nil, err
	}

	if time.Now().After(payload.Expiration) {
		return nil, fmt.Errorf("token expired")
	}

	log.Println(payload)
	return &payload, nil
}
