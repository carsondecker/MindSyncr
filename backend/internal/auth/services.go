package auth

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/carsondecker/MindSyncr/internal/db/sqlc"
	"github.com/carsondecker/MindSyncr/internal/utils"
	"github.com/lib/pq"
)

func (h *AuthHandler) registerService(ctx context.Context, email, username, password string) (RegisterResponse, string, string, *utils.ServiceError) {
	passwordHash, err := hashPassword(password)
	if err != nil {
		return RegisterResponse{}, "", "", &utils.ServiceError{
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
				return RegisterResponse{}, "", "", &utils.ServiceError{
					StatusCode: http.StatusBadRequest,
					Code:       utils.ErrUserAlreadyExists,
					Message:    "this email is already in use",
				}
			}
		}
		return RegisterResponse{}, "", "", &utils.ServiceError{
			StatusCode: http.StatusBadRequest,
			Code:       utils.ErrDbtxFail,
			Message:    err.Error(),
		}
	}

	jwtToken, err := utils.CreateJWT(row.ID, row.Email, row.Username, row.Role)
	if err != nil {
		return RegisterResponse{}, "", "", &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrJwtFail,
			Message:    err.Error(),
		}
	}

	refreshToken, refreshRes, err := createRefreshToken(ctx, h.cfg.Queries, row.ID)
	if err != nil {
		return RegisterResponse{}, "", "", &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrRefreshFail,
			Message:    err.Error(),
		}
	}

	res := RegisterResponse{
		Id:           row.ID,
		Email:        row.Email,
		Username:     row.Username,
		CreatedAt:    row.CreatedAt,
		RefreshToken: refreshRes,
	}

	return res, jwtToken, refreshToken, nil
}

func (h *AuthHandler) loginService(ctx context.Context, email, password string) (LoginResponse, string, string, *utils.ServiceError) {
	row, err := h.cfg.Queries.GetUserForLogin(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return LoginResponse{}, "", "", &utils.ServiceError{
				StatusCode: http.StatusUnauthorized,
				Code:       utils.ErrInvalidCredentials,
				Message:    err.Error(),
			}
		}
		return LoginResponse{}, "", "", &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrDbtxFail,
			Message:    err.Error(),
		}
	}

	err = checkPassword(row.PasswordHash, password)
	if err != nil {
		return LoginResponse{}, "", "", &utils.ServiceError{
			StatusCode: http.StatusUnauthorized,
			Code:       utils.ErrInvalidCredentials,
			Message:    err.Error(),
		}
	}

	jwtToken, err := utils.CreateJWT(row.ID, row.Email, row.Username, row.Role)
	if err != nil {
		return LoginResponse{}, "", "", &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrJwtFail,
			Message:    err.Error(),
		}
	}

	refreshToken, refreshRes, err := createRefreshToken(ctx, h.cfg.Queries, row.ID)
	if err != nil {
		return LoginResponse{}, "", "", &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrRefreshFail,
			Message:    err.Error(),
		}
	}

	res := LoginResponse{
		Id:           row.ID,
		Email:        row.Email,
		Username:     row.Username,
		RefreshToken: refreshRes,
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

	isValid, userId, sErr := isValidRefreshToken(ctx, qtx, token)
	if sErr != nil {
		return "", "", RefreshTokenResponse{}, sErr
	}

	if !isValid {
		return "", "", RefreshTokenResponse{}, &utils.ServiceError{
			StatusCode: http.StatusUnauthorized,
			Code:       utils.ErrInvalidRefreshToken,
			Message:    "refresh token is invalid",
		}
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

	jwtToken, sErr := CreateJWTById(ctx, qtx, userId)
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

func (h *AuthHandler) logoutService(ctx context.Context, token string) *utils.ServiceError {
	isValid, userId, sErr := isValidRefreshToken(ctx, h.cfg.Queries, token)
	if sErr != nil {
		return sErr
	}

	if !isValid {
		return &utils.ServiceError{
			StatusCode: http.StatusUnauthorized,
			Code:       utils.ErrInvalidRefreshToken,
			Message:    "refresh token is invalid",
		}
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
