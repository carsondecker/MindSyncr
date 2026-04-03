package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/carsondecker/MindSyncr/utils"

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

func (s *AuthService) createRefreshToken(userId uuid.UUID) (string, RefreshTokenResponse, *utils.ServiceError) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", RefreshTokenResponse{}, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrRefreshFail,
			Message:    fmt.Sprintf("failed to generate random token: %s", err.Error()),
		}
	}

	token := base64.URLEncoding.EncodeToString(b)
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	hashBytes := sha256.Sum256([]byte(token))
	tokenHash := base64.RawURLEncoding.EncodeToString(hashBytes[:])

	refreshTokenRes, sErr := s.repo.InsertRefreshToken(userId, tokenHash, expiresAt)
	if sErr != nil {
		return "", RefreshTokenResponse{}, sErr
	}

	return token, refreshTokenRes, nil
}

func (s *AuthService) isValidRefreshToken(token string) (uuid.UUID, *utils.ServiceError) {
	hashBytes := sha256.Sum256([]byte(token))
	tokenHash := base64.RawURLEncoding.EncodeToString(hashBytes[:])

	userId, sErr := s.repo.CheckValidRefreshToken(tokenHash)
	if sErr != nil {
		return uuid.Nil, sErr
	}

	if userId == uuid.Nil {
		return uuid.Nil, &utils.ServiceError{
			StatusCode: http.StatusUnauthorized,
			Code:       utils.ErrInvalidRefreshToken,
			Message:    "could not get user id from refresh token entry",
		}
	}

	return userId, nil
}

func (s *AuthService) createJWTById(userId uuid.UUID) (string, *utils.ServiceError) {
	user, sErr := s.repo.GetUserById(userId)
	if sErr != nil {
		return "", sErr
	}

	jwtToken, err := utils.CreateJWT(user.Id, user.Email, user.Username, user.Role)
	if err != nil {
		return "", &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrJwtFail,
			Message:    err.Error(),
		}
	}

	return jwtToken, nil
}

func getRefreshToken(r *http.Request) (string, *utils.ServiceError) {
	refreshTokenCookie, err := r.Cookie("refresh_token")
	if err != nil || refreshTokenCookie.Value == "" {
		return "", &utils.ServiceError{
			StatusCode: http.StatusBadRequest,
			Code:       utils.ErrBadRequest,
			Message:    "no refresh token cookie provided",
		}
	}

	return refreshTokenCookie.Value, nil
}
