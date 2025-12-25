package auth

import (
	"time"

	"github.com/google/uuid"
)

type RegisterRequest struct {
	Email           string `json:"email" validate:"required,email"`
	Username        string `json:"username" validate:"required,min=8,max=36"`
	Password        string `json:"password" validate:"required,password"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserWithRefresh struct {
	Id              uuid.UUID            `json:"id"`
	Email           string               `json:"email"`
	Username        string               `json:"username"`
	Role            string               `json:"role"`
	Status          string               `json:"status"`
	IsEmailVerified bool                 `json:"is_email_verified"`
	CreatedAt       time.Time            `json:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at"`
	RefreshToken    RefreshTokenResponse `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	Id              uuid.UUID `json:"id"`
	Email           string    `json:"email"`
	Username        string    `json:"username"`
	Role            string    `json:"role"`
	Status          string    `json:"status"`
	IsEmailVerified bool      `json:"is_email_verified"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
