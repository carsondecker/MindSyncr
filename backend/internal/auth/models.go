package auth

import "github.com/google/uuid"

type RegisterRequest struct {
	Email           string `json:"email"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type RegisterResponse struct {
	Id       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
}
