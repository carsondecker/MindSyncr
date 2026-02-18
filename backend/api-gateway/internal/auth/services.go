package auth

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/carsondecker/MindSyncr/internal/db/sqlc"
	"github.com/carsondecker/MindSyncr/utils"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func (h *AuthHandler) registerService(ctx context.Context, email, username, password string) (UserWithRefresh, string, string, *utils.ServiceError) {
	passwordHash, err := hashPassword(password)
	if err != nil {
		return UserWithRefresh{}, "", "", &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrHashFail,
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
				return UserWithRefresh{}, "", "", &utils.ServiceError{
					StatusCode: http.StatusBadRequest,
					Code:       utils.ErrUserAlreadyExists,
					Message:    "this email is already in use",
				}
			}
		}
		return UserWithRefresh{}, "", "", &utils.ServiceError{
			StatusCode: http.StatusBadRequest,
			Code:       utils.ErrDbtxFail,
			Message:    err.Error(),
		}
	}

	jwtToken, err := utils.CreateJWT(row.ID, row.Email, row.Username, row.Role)
	if err != nil {
		return UserWithRefresh{}, "", "", &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrJwtFail,
			Message:    err.Error(),
		}
	}

	refreshToken, refreshRes, err := createRefreshToken(ctx, h.cfg.Queries, row.ID)
	if err != nil {
		return UserWithRefresh{}, "", "", &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrRefreshFail,
			Message:    err.Error(),
		}
	}

	res := UserWithRefresh{
		Id:              row.ID,
		Email:           row.Email,
		Username:        row.Username,
		Role:            row.Role,
		Status:          row.Status,
		IsEmailVerified: row.IsEmailVerified,
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
		RefreshToken:    refreshRes,
	}

	return res, jwtToken, refreshToken, nil
}

func (h *AuthHandler) loginService(ctx context.Context, email, password string) (UserWithRefresh, string, string, *utils.ServiceError) {
	row, err := h.cfg.Queries.GetUserForLogin(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return UserWithRefresh{}, "", "", &utils.ServiceError{
				StatusCode: http.StatusUnauthorized,
				Code:       utils.ErrInvalidCredentials,
				Message:    err.Error(),
			}
		}
		return UserWithRefresh{}, "", "", &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrDbtxFail,
			Message:    err.Error(),
		}
	}

	err = checkPassword(row.PasswordHash, password)
	if err != nil {
		return UserWithRefresh{}, "", "", &utils.ServiceError{
			StatusCode: http.StatusUnauthorized,
			Code:       utils.ErrInvalidCredentials,
			Message:    err.Error(),
		}
	}

	jwtToken, err := utils.CreateJWT(row.ID, row.Email, row.Username, row.Role)
	if err != nil {
		return UserWithRefresh{}, "", "", &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrJwtFail,
			Message:    err.Error(),
		}
	}

	refreshToken, refreshRes, err := createRefreshToken(ctx, h.cfg.Queries, row.ID)
	if err != nil {
		return UserWithRefresh{}, "", "", &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrRefreshFail,
			Message:    err.Error(),
		}
	}

	res := UserWithRefresh{
		Id:              row.ID,
		Email:           row.Email,
		Username:        row.Username,
		Role:            row.Role,
		Status:          row.Status,
		IsEmailVerified: row.IsEmailVerified,
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
		RefreshToken:    refreshRes,
	}

	return res, jwtToken, refreshToken, nil
}

func (h *AuthHandler) refreshService(ctx context.Context, token string) (string, string, RefreshTokenResponse, *utils.ServiceError) {
	tx, err := h.cfg.DB.BeginTx(ctx, nil)
	if err != nil {
		return "", "", RefreshTokenResponse{}, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrTxBeginFail,
			Message:    err.Error(),
		}
	}

	defer tx.Rollback()

	qtx := h.cfg.Queries.WithTx(tx)

	userId, sErr := isValidRefreshToken(ctx, qtx, token)
	if sErr != nil {
		return "", "", RefreshTokenResponse{}, sErr
	}

	err = qtx.RevokeUserTokens(ctx, userId)
	if err != nil {
		return "", "", RefreshTokenResponse{}, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrRefreshRevokeFail,
			Message:    err.Error(),
		}
	}

	refreshToken, res, err := createRefreshToken(ctx, qtx, userId)
	if err != nil {
		return "", "", RefreshTokenResponse{}, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrRefreshFail,
			Message:    err.Error(),
		}
	}

	jwtToken, sErr := createJWTById(ctx, qtx, userId)
	if sErr != nil {
		return "", "", RefreshTokenResponse{}, sErr
	}

	err = tx.Commit()
	if err != nil {
		return "", "", RefreshTokenResponse{}, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrTxCommitFail,
			Message:    err.Error(),
		}
	}

	return jwtToken, refreshToken, res, nil
}

func (h *AuthHandler) logoutService(ctx context.Context, userId uuid.UUID, token string) *utils.ServiceError {
	_, sErr := isValidRefreshToken(ctx, h.cfg.Queries, token)
	if sErr != nil {
		return sErr
	}

	err := h.cfg.Queries.RevokeUserTokens(ctx, userId)
	if err != nil {
		return &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrRefreshRevokeFail,
			Message:    err.Error(),
		}
	}

	return nil
}

func (h *AuthHandler) getUserService(ctx context.Context, userId uuid.UUID) (User, *utils.ServiceError) {
	row, err := h.cfg.Queries.GetUserById(ctx, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, &utils.ServiceError{
				StatusCode: http.StatusNotFound,
				Code:       utils.ErrUserNotFound,
				Message:    err.Error(),
			}
		}
		return User{}, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrDbtxFail,
			Message:    err.Error(),
		}
	}

	res := User{
		Id:              row.ID,
		Email:           row.Email,
		Username:        row.Username,
		Role:            row.Role,
		Status:          row.Status,
		IsEmailVerified: row.IsEmailVerified,
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
	}

	return res, nil
}
