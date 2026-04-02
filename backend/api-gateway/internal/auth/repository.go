package auth

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/carsondecker/MindSyncr/internal/db/sqlc"
	"github.com/carsondecker/MindSyncr/utils"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type AuthRepository interface {
	Register(email, username, passwordHash string) (User, *utils.ServiceError)
	InsertRefreshToken(userId uuid.UUID, tokenHash string, expiresAt time.Time) (RefreshTokenResponse, *utils.ServiceError)
	GetInternalUser(email string) (InternalUser, *utils.ServiceError)
}

type PostgresAuthRepository struct {
	db      *sql.DB
	queries *sqlc.Queries
}

func NewPostgresAuthRepository(db *sql.DB, queries *sqlc.Queries) *PostgresAuthRepository {
	return &PostgresAuthRepository{
		db,
		queries,
	}
}

func (r *PostgresAuthRepository) Register(email, username, passwordHash string) (User, *utils.ServiceError) {
	row, err := r.queries.InsertUser(context.Background(), sqlc.InsertUserParams{
		Email:        email,
		Username:     username,
		PasswordHash: passwordHash,
	})

	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				return User{}, &utils.ServiceError{
					StatusCode: http.StatusBadRequest,
					Code:       utils.ErrUserAlreadyExists,
					Message:    "this email is already in use",
				}
			}
		}
		return User{}, &utils.ServiceError{
			StatusCode: http.StatusBadRequest,
			Code:       utils.ErrDbtxFail,
			Message:    err.Error(),
		}
	}

	return User{
		Id:              row.ID,
		Email:           row.Email,
		Username:        row.Username,
		Role:            row.Role,
		Status:          row.Status,
		IsEmailVerified: row.IsEmailVerified,
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
	}, nil
}

func (r *PostgresAuthRepository) InsertRefreshToken(userId uuid.UUID, tokenHash string, expiresAt time.Time) (RefreshTokenResponse, *utils.ServiceError) {
	row, err := r.queries.InsertRefreshToken(context.Background(), sqlc.InsertRefreshTokenParams{
		UserID:    userId,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return RefreshTokenResponse{}, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrRefreshFail,
			Message:    fmt.Sprintf("failed to create refresh token: %w", err),
		}
	}

	return RefreshTokenResponse{
		ExpiresAt: row.ExpiresAt,
		CreatedAt: row.CreatedAt,
	}, nil
}

func (r *PostgresAuthRepository) GetInternalUser(email string) (InternalUser, *utils.ServiceError) {
	row, err := r.queries.GetUserForLogin(context.Background(), email)
	if err != nil {
		if err == sql.ErrNoRows {
			return InternalUser{}, &utils.ServiceError{
				StatusCode: http.StatusUnauthorized,
				Code:       utils.ErrInvalidCredentials,
				Message:    err.Error(),
			}
		}
		return InternalUser{}, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrDbtxFail,
			Message:    err.Error(),
		}
	}

	return InternalUser{
		User: User{
			Id:              row.ID,
			Email:           row.Email,
			Username:        row.Username,
			Role:            row.Role,
			Status:          row.Status,
			IsEmailVerified: row.IsEmailVerified,
			CreatedAt:       row.CreatedAt,
			UpdatedAt:       row.UpdatedAt,
		},
		PasswordHash: row.PasswordHash,
	}, nil
}
