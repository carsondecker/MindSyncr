package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/carsondecker/MindSyncr/internal/db/sqlc"
	"github.com/carsondecker/MindSyncr/internal/utils"
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

func createRefreshToken(ctx context.Context, q *sqlc.Queries, userId uuid.UUID) (string, RefreshTokenResponse, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", RefreshTokenResponse{}, fmt.Errorf("failed to generate random token: %w", err)
	}

	token := base64.URLEncoding.EncodeToString(b)
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	hashBytes := sha256.Sum256([]byte(token))
	tokenHash := base64.RawURLEncoding.EncodeToString(hashBytes[:])

	row, err := q.InsertRefreshToken(ctx, sqlc.InsertRefreshTokenParams{
		UserID:    userId,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return "", RefreshTokenResponse{}, fmt.Errorf("failed to create refresh token: %w", err)
	}

	return token, RefreshTokenResponse{
		ExpiresAt: row.ExpiresAt,
		CreatedAt: row.CreatedAt,
	}, nil
}

func isValidRefreshToken(ctx context.Context, q *sqlc.Queries, token string) (bool, uuid.UUID, *utils.ServiceError) {
	hashBytes := sha256.Sum256([]byte(token))
	tokenHash := base64.RawURLEncoding.EncodeToString(hashBytes[:])

	userId, err := q.CheckValidRefreshToken(ctx, tokenHash)
	if err != nil {
		return false, uuid.Nil, &utils.ServiceError{
			StatusCode: http.StatusUnauthorized,
			Code:       utils.ErrInvalidRefreshToken,
			Message:    err.Error(),
		}
	}

	if userId == uuid.Nil {
		return false, uuid.Nil, &utils.ServiceError{
			StatusCode: http.StatusUnauthorized,
			Code:       utils.ErrInvalidRefreshToken,
			Message:    "could not get user id from refresh token entry",
		}
	}

	return true, userId, nil
}

func CreateJWTById(ctx context.Context, q *sqlc.Queries, userId uuid.UUID) (string, *utils.ServiceError) {
	row, err := q.GetUserById(ctx, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", &utils.ServiceError{
				StatusCode: http.StatusNotFound,
				Code:       utils.ErrUserNotFound,
				Message:    err.Error(),
			}
		}
		return "", &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrDbtxFail,
			Message:    err.Error(),
		}
	}

	jwtToken, err := utils.CreateJWT(row.ID, row.Email, row.Username, row.Role)
	if err != nil {
		return "", &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrJwtFail,
			Message:    err.Error(),
		}
	}

	return jwtToken, nil
}
