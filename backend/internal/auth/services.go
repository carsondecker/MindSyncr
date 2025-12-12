package auth

import (
	"context"

	"github.com/carsondecker/MindSyncr/internal/db/sqlc"
	"github.com/carsondecker/MindSyncr/internal/utils"
	"github.com/lib/pq"
)

func (h *AuthHandler) registerService(ctx context.Context, email, username, password string) (RegisterResponse, string, *utils.ServiceError) {
	passwordHash, err := hashPassword(password)
	if err != nil {
		return RegisterResponse{}, "", &utils.ServiceError{
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
				return RegisterResponse{}, "", &utils.ServiceError{
					StatusCode: 400,
					Code:       "USER_ALREADY_EXISTS",
					Message:    "this email is already in use",
				}
			}
		}
		return RegisterResponse{}, "", &utils.ServiceError{
			StatusCode: 500,
			Code:       "DBTX_FAIL",
			Message:    err.Error(),
		}
	}

	jwtToken, err := utils.CreateJWT(row.ID, row.Email, row.Username)

	if err != nil {
		return RegisterResponse{}, "", &utils.ServiceError{
			StatusCode: 500,
			Code:       "JWT_FAIL",
			Message:    err.Error(),
		}
	}

	res := RegisterResponse{
		Id:       row.ID,
		Email:    row.Email,
		Username: row.Username,
	}

	return res, jwtToken, nil
}
