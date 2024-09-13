package models

type User struct {
    ID           int    `json:"id"`
    Username     string `json:"username" binding:"required,alphanum,min=6"`
    Email        string `json:"email" binding:"required,min=6"`
    Password     string `json:"password,omitempty" binding:"required,min=6"`
    PasswordHash string `json:"-"`
}
type SignUp struct {
    Username     string `json:"username" binding:"required,alphanum,min=6"`
    Email        string `json:"email" binding:"required,min=6"`
    Password     string `json:"password,omitempty" binding:"required,min=6"`
}
type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}
