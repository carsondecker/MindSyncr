package auth

import "github.com/google/uuid"

type RegisterRequest struct {
	Email           string `json:"email" validation:"required,email"`
	Username        string `json:"username" validation:"required,min=8,max=36"`
	Password        string `json:"password" validation:"required,password"`
	ConfirmPassword string `json:"confirm_password" validation:"required,eqfield=Password"`
}

type RegisterResponse struct {
	Id       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
}
