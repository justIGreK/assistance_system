package models

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username" validate:"required,min=6"`
	Email        string `json:"email" validate:"required,min=6"`
	Password     string `json:"password,omitempty" validate:"required,min=6"`
	PasswordHash string `json:"-"`
	Role         string `json:"role"`
	Banned       bool   `json:"banned"`
}
type SignUp struct {
	Username string `json:"username" validate:"required,min=6,max=15"`
	Email    string `json:"email" validate:"required,min=6,max=20"`
	Password string `json:"password,omitempty" validate:"required,min=6,max=30"`
}
type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

const (
	CustomerRole       string = "customer"
	AdministrationRole string = "admin"
)
