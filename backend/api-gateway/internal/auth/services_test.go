package auth

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/carsondecker/MindSyncr/internal/db/sqlc"
	"github.com/carsondecker/MindSyncr/internal/sutils"
	"github.com/carsondecker/MindSyncr/utils"
	"github.com/google/uuid"
	"github.com/gopherjs/gopherjs/compiler/natives/src/strings"
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

			body := strings.NewReader(fmt.Sprintf(`{
					"email":"%s",
					"username":"%s",
					"password":"%s",
					"confirm_password":"%s",
				}`,
				tc.email,
				tc.username,
				tc.password,
				tc.confirmPassword))
			req := httptest.NewRequest(http.MethodPost, "/", body)
			rr := httptest.NewRecorder()

			s.HandleRegister(rr, req)

			if rr.Header().Get("Content-Type") != "application/json" {
				t.Errorf("wrong content type")
			}

			expected := `{"message":"ok"}`
			if rr.Body.String() != expected {
				t.Errorf("expected %s, got %s", expected, rr.Body.String())
			}

			user, jwtToken, refreshToken, sErr := s.registerService(tc.email, tc.username, tc.password)

			if tc.success {
				if sErr != nil {
					t.Fatalf("expected no error, got %s", sErr.Message)
				}
				if user.Email != tc.email {
					t.Errorf("expected email %q, got %q", tc.email, user.Email)
				}
				if user.Username != tc.username {
					t.Errorf("expected username %q, got %q", tc.username, user.Username)
				}
				if user.Id == uuid.Nil {
					t.Error("expected non-nil user ID")
				}
				if jwtToken == "" {
					t.Error("expected jwt token")
				}
				if refreshToken == "" {
					t.Error("expected refresh token")
				}
			} else {
				if sErr == nil {
					t.Fatal("expected error, got nil")
				}
				emptyUser := UserWithRefresh{}
				if user != emptyUser {
					t.Errorf("expected no user, got %v", user)
				}
				if jwtToken != "" {
					t.Errorf("expected no jwt token, got %s", jwtToken)
				}
				if refreshToken != "" {
					t.Errorf("expected no refresh token, got %s", refreshToken)
				}
			}
		})
	}
}
