package auth

import (
	"context"
	"errors"
	"fmt"
	"gohelp/internal/models"
	"gohelp/internal/storage/postgresql"

	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	CreateUser(ctx context.Context, user models.User) error
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
}

type UserService struct {
	UserRepo
}

func NewUserService(userRepo *postgresql.UserRepository) *UserService {
	return &UserService{UserRepo: userRepo}
}

func (s *UserService) RegisterUser(ctx context.Context, user models.SignUp) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	newUser := models.User{
		Username:     user.Username,
		Email:        user.Email,
		Password:     user.Password,
		PasswordHash: string(hashedPassword),
	}
	err = s.CreateUser(ctx, newUser)
	if err != nil {
		return fmt.Errorf("error during creating user:%v", err)
	}
	return nil
}

func (s *UserService) LoginUser(ctx context.Context, email, password string) (string, error) {
	user, err := s.GetUserByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := GeneratePasetoToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("error during generating token: %v", err)
	}

	return token, nil
}
