package auth

import (
	"context"
	"errors"
	"fmt"
	"gohelp/internal/models"
	"gohelp/internal/storage/postgresql"
	"gohelp/util"

	"github.com/markbates/goth"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	CreateUser(ctx context.Context, user models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserById(ctx context.Context, userID int) (*models.User, error)
	ChangeBanStatus(ctx context.Context, userID int, status bool) error
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

	if user.Banned {
		return "", errors.New("your account is banned")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := GeneratePasetoToken(user.ID, user.Role)
	if err != nil {
		return "", fmt.Errorf("error during generating token: %v", err)
	}

	return token, nil
}

func (s *UserService) UsersActions(ctx context.Context, userID int, action string) (*models.User, error) {
	user, err := s.GetUserById(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error during getting user by id: %v", err)
	}
	var status bool
	if action == "ban" {
		status = true
	} else {
		status = false
	}

	err = s.ChangeBanStatus(ctx, userID, status)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GoogleAuth(ctx context.Context, googleUser goth.User) (string, error){
	user, err := s.GetUserByEmail(ctx, googleUser.Email)
	if user.Banned{
		return "", errors.New("your account is banned")
	}
	if err !=  nil{
		newUser := models.User{
			Username:     util.GenerateNickname(),
			Email:        googleUser.Email,
			Password:     " ",
			PasswordHash: " ",
		}
		s.CreateUser(ctx, newUser)
		user, err = s.GetUserByEmail(ctx, googleUser.Email)
		if err != nil{
			return "", err
		}
	}
	token, err:= GeneratePasetoToken(user.ID, user.Role)
	if err != nil {
		return "", fmt.Errorf("error during generating token: %v", err)
	}
	return token, nil

}
