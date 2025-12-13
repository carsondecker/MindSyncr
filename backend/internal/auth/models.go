package auth

import "github.com/google/uuid"

type RegisterRequest struct {
	Email           string `json:"email" validate:"required,email"`
	Username        string `json:"username" validate:"required,min=8,max=36"`
	Password        string `json:"password" validate:"required,password"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

type RegisterResponse struct {
	Id       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Id       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
}
