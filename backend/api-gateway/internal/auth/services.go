package auth

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/carsondecker/MindSyncr/utils"
	"github.com/google/uuid"
)

func (s *AuthService) registerService(ctx context.Context, email, username, password string) (UserWithRefresh, string, string, *utils.ServiceError) {
	passwordHash, err := hashPassword(password)
	if err != nil {
		return UserWithRefresh{}, "", "", &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrHashFail,
			Message:    err.Error(),
		}
	}

	user, sErr := s.repo.Register(email, username, passwordHash)

	if sErr != nil {
		return UserWithRefresh{}, "", "", sErr
	}

	jwtToken, err := utils.CreateJWT(user.Id, user.Email, user.Username, user.Role)
	if err != nil {
		return UserWithRefresh{}, "", "", &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrJwtFail,
			Message:    err.Error(),
		}
	}

	refreshToken, refreshRes, sErr := s.createRefreshToken(user.Id)
	if sErr != nil {
		return UserWithRefresh{}, "", "", sErr
	}

	res := UserWithRefresh{
		User:         user,
		RefreshToken: refreshRes,
	}

	return res, jwtToken, refreshToken, nil
}

func (s *AuthService) loginService(ctx context.Context, email, password string) (UserWithRefresh, string, string, *utils.ServiceError) {
	internalUser, sErr := s.repo.GetInternalUser(email)
	if sErr != nil {
		return UserWithRefresh{}, "", "", sErr
	}

	err := checkPassword(internalUser.PasswordHash, password)
	if err != nil {
		return UserWithRefresh{}, "", "", &utils.ServiceError{
			StatusCode: http.StatusUnauthorized,
			Code:       utils.ErrInvalidCredentials,
			Message:    err.Error(),
		}
	}

	jwtToken, err := utils.CreateJWT(internalUser.Id, internalUser.Email, internalUser.Username, internalUser.Role)
	if err != nil {
		return UserWithRefresh{}, "", "", &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrJwtFail,
			Message:    err.Error(),
		}
	}

	refreshToken, refreshRes, sErr := s.createRefreshToken(internalUser.Id)
	if sErr != nil {
		return UserWithRefresh{}, "", "", sErr
	}

	res := UserWithRefresh{
		User:         internalUser.User,
		RefreshToken: refreshRes,
	}

	return res, jwtToken, refreshToken, nil
}

func (h *AuthService) refreshService(ctx context.Context, token string) (string, string, RefreshTokenResponse, *utils.ServiceError) {
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

func (h *AuthService) logoutService(ctx context.Context, userId uuid.UUID, token string) *utils.ServiceError {
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

func (h *AuthService) getUserService(ctx context.Context, userId uuid.UUID) (User, *utils.ServiceError) {
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
