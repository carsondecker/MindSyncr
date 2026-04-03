package auth

import (
	"database/sql"
	"net/http"
	"testing"
	"time"

	"github.com/carsondecker/MindSyncr/utils"
	"github.com/google/uuid"
)

var mockDBError = &utils.ServiceError{
	StatusCode: 500,
	Code:       utils.ErrDbtxFail,
	Message:    "",
}

type MockRefreshToken struct {
	userId    uuid.UUID
	tokenHash string
	isRevoked bool
}

type MockAuthRepository struct {
	users         []InternalUser
	refreshTokens []MockRefreshToken
	returnError   bool
}

func (r *MockAuthRepository) BeginTx() (*sql.Tx, *utils.ServiceError) {
	if r.returnError {
		return nil, mockDBError
	}

	return &sql.Tx{}, nil
}

func (r *MockAuthRepository) GetTxRepo(tx *sql.Tx) AuthRepository {
	return r
}

func (r *MockAuthRepository) Register(email, username, passwordHash string) (User, *utils.ServiceError) {
	if r.returnError {
		return User{}, mockDBError
	}

	uuid, _ := uuid.NewUUID()
	internalUser := InternalUser{
		User: User{
			Id:       uuid,
			Email:    email,
			Username: username,
		},
		PasswordHash: passwordHash,
	}

	r.users = append(r.users, internalUser)

	return internalUser.User, nil
}

func (r *MockAuthRepository) InsertRefreshToken(userId uuid.UUID, tokenHash string, expiresAt time.Time) (RefreshTokenResponse, *utils.ServiceError) {
	if r.returnError {
		return RefreshTokenResponse{}, mockDBError
	}

	return RefreshTokenResponse{}, nil
}

func (r *MockAuthRepository) GetInternalUser(email string) (InternalUser, *utils.ServiceError) {
	if r.returnError {
		return InternalUser{}, mockDBError
	}

	for _, internalUser := range r.users {
		if internalUser.Email == email {
			return InternalUser{}, nil
		}
	}

	return InternalUser{}, &utils.ServiceError{
		StatusCode: http.StatusUnauthorized,
		Code:       utils.ErrInvalidCredentials,
		Message:    "",
	}
}

func (r *MockAuthRepository) CheckValidRefreshToken(tokenHash string) (uuid.UUID, *utils.ServiceError) {
	if r.returnError {
		return uuid.Nil, mockDBError
	}

	for _, refreshToken := range r.refreshTokens {
		if refreshToken.tokenHash == tokenHash {
			return refreshToken.userId, nil
		}
	}

	return uuid.Nil, &utils.ServiceError{
		StatusCode: http.StatusUnauthorized,
		Code:       utils.ErrInvalidRefreshToken,
		Message:    "",
	}
}

func (r *MockAuthRepository) RevokeUserTokens(userId uuid.UUID) *utils.ServiceError {
	if r.returnError {
		return mockDBError
	}

	for _, refreshToken := range r.refreshTokens {
		if refreshToken.userId == userId {
			refreshToken.isRevoked = true
		}
	}
	return nil
}

func (r *MockAuthRepository) GetUserById(userId uuid.UUID) (User, *utils.ServiceError) {
	if r.returnError {
		return User{}, mockDBError
	}

	for _, internalUser := range r.users {
		if internalUser.Id == userId {
			return internalUser.User, nil
		}
	}

	return User{}, &utils.ServiceError{
		StatusCode: http.StatusNotFound,
		Code:       utils.ErrUserNotFound,
		Message:    "",
	}
}

func TestRegister(t *testing.T) {
	tcs := []struct {
		name string
		test func(*AuthService, *testing.T)
	}{
		{
			name: "success",
			test: func(t *testing.T) {

			},
		},
	}

	for _, tc := range tcs {

	}
}
