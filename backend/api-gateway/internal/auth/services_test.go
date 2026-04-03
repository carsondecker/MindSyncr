package auth

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/carsondecker/MindSyncr/internal/db/sqlc"
	"github.com/carsondecker/MindSyncr/internal/sutils"
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

	for _, internalUser := range r.users {
		if internalUser.Email == email {
			return User{}, &utils.ServiceError{
				StatusCode: http.StatusBadRequest,
				Code:       utils.ErrUserAlreadyExists,
				Message:    "this email is already in use",
			}
		}
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
			return internalUser, nil
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

	for i := range r.refreshTokens {
		if r.refreshTokens[i].userId == userId {
			r.refreshTokens[i].isRevoked = true
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

// TODO - need to use http tests for most of these things
func TestRegister(t *testing.T) {
	tcs := []struct {
		name            string
		email           string
		username        string
		password        string
		confirmPassword string
		success         bool
		setup           func(*MockAuthRepository)
	}{
		{
			name:     "success",
			email:    "johndoe@example.com",
			username: "jonathandoe",
			password: "Iamsoreal123!",
			success:  true,
			setup:    nil,
		},
		{
			name:            "fail - duplicate email",
			email:           "johndoe@example.com",
			username:        "jonathandoe",
			password:        "Iamsoreal123!",
			confirmPassword: "Iamsoreal123!",
			success:         false,
			setup: func(r *MockAuthRepository) {
				r.Register("johndoe@example.com", "testtesttest", "Password123!")
			},
		},
		{
			name:            "fail - duplicate email, case insensitive",
			email:           "JohnDoe@example.com",
			username:        "jonathandoe",
			password:        "Iamsoreal123!",
			confirmPassword: "Iamsoreal123!",
			success:         false,
			setup: func(r *MockAuthRepository) {
				r.Register("johndoe@example.com", "testtesttest", "Password123!")
			},
		},
		{
			name:            "fail - username less than min length",
			email:           "johndoe@example.com",
			username:        "jonathandoe",
			password:        "Iamsoreal123!",
			confirmPassword: "Iamsoreal123!",
			success:         false,
			setup:           nil,
		},
	}

	utils.JWTInit("jwtSecret", "wsSecret")

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			repo := &MockAuthRepository{}
			if tc.setup != nil {
				tc.setup(repo)
			}
			s := AuthService{
				repo: repo,
				cfg:  sutils.NewConfig(&sql.DB{}, &sqlc.Queries{}, &utils.RedisClient{}),
			}

			body, err := json.Marshal(map[string]string{
				"email":            tc.email,
				"username":         tc.username,
				"password":         tc.password,
				"confirm_password": tc.confirmPassword,
			})
			if err != nil {
				t.Fatal(err)
			}

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
			rr := httptest.NewRecorder()

			s.HandleRegister(rr, req)

			if rr.Header().Get("Content-Type") != "application/json" {
				t.Errorf("wrong content type")
			}

			if tc.success {
				if rr.Code >= 400 {
					res, err := utils.ParseError(rr.Body)
					if err != nil {
						t.Fatalf("failed to marshal error response")
					}

					t.Fatalf("expected no error, got %s", res.Error.Message)
				}

				res, err := utils.ParseSuccess[UserWithRefresh](rr.Body)
				if err != nil {
					t.Fatalf("failed to marshal success response")
				}

				if !res.Success {
					t.Errorf("recieved success response with success as false")
				}
				if res.Data.Email != tc.email {
					t.Errorf("expected email %q, got %q", tc.email, res.Data.Email)
				}
				if res.Data.Username != tc.username {
					t.Errorf("expected username %q, got %q", tc.username, res.Data.Username)
				}
				if res.Data.Id == uuid.Nil {
					t.Error("expected non-nil user ID")
				}
			} else {
				if rr.Code < 400 {
					t.Fatal("expected error, got success code")
				}

				/*
					res, err := utils.ParseError(rr.Body)
					if err != nil {
						t.Fatalf("failed to marshal error response")
					}
				*/
			}
		})
	}
}
