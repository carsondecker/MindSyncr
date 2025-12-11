package auth

import (
	"context"

	"github.com/carsondecker/MindSyncr/internal/db/sqlc"
	"github.com/carsondecker/MindSyncr/internal/utils"
	"github.com/lib/pq"
)

func (h *AuthHandler) registerService(ctx context.Context, email, username, password string) (sqlc.InsertUserRow, *utils.ServiceError) {
	passwordHash, err := hashPassword(password)
	if err != nil {
		return sqlc.InsertUserRow{}, &utils.ServiceError{
			StatusCode: 500,
			Code:       "HASH_FAIL",
			Message:    err.Error(),
		}
	}

	row, err := h.cfg.Queries.InsertUser(ctx, sqlc.InsertUserParams{
		Email:        email,
		Username:     username,
		PasswordHash: passwordHash,
	})

	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				return sqlc.InsertUserRow{}, &utils.ServiceError{
					StatusCode: 400,
					Code:       "USER_ALREADY_EXISTS",
					Message:    "this email is already in use",
				}
			}
		}
		return sqlc.InsertUserRow{}, &utils.ServiceError{
			StatusCode: 500,
			Code:       "DBTX_FAIL",
			Message:    err.Error(),
		}
	}

	return row, nil
}
