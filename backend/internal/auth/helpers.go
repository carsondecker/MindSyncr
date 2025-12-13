package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/carsondecker/MindSyncr/internal/db/sqlc"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(hash), err
}

func checkPassword(hashedPassword string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

func (h *AuthHandler) createRefreshToken(ctx context.Context, id uuid.UUID) (string, RefreshTokenResponse, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", RefreshTokenResponse{}, fmt.Errorf("failed to generate random token: %w", err)
	}

	token := base64.URLEncoding.EncodeToString(b)
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	row, err := h.cfg.Queries.InsertRefreshToken(ctx, sqlc.InsertRefreshTokenParams{
		UserID:    id,
		Token:     token,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return "", RefreshTokenResponse{}, fmt.Errorf("failed to create refresh token: %w", err)
	}

	return row.Token, RefreshTokenResponse{
		ExpiresAt: row.ExpiresAt,
		CreatedAt: row.CreatedAt,
	}, nil
}
