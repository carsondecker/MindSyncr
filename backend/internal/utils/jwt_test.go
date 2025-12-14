package utils

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateJWTAndGetClaims(t *testing.T) {
	os.Setenv("JWT_SECRET", "testsecret123")

	tcs := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "success",
			test: func(t *testing.T) {
				userId := uuid.New()
				email := "test@gmail.com"
				username := "testuser"
				role := "user"

				tokenString, err := CreateJWT(userId, email, username, role)

				require.NoError(t, err)
				require.NotEmpty(t, tokenString)

				claims, err := GetClaims(tokenString)

				require.NoError(t, err)
				require.NotNil(t, claims)

				require.Equal(t, userId, claims.UserId)
				require.Equal(t, email, claims.Email)
				require.Equal(t, username, claims.Username)
				require.Equal(t, "MindSyncr", claims.Issuer)

				require.WithinDuration(t,
					time.Now().Add(15*time.Minute),
					claims.ExpiresAt.Time,
					2*time.Second,
				)
			},
		},
		{
			name: "failure - no secret key",
			test: func(t *testing.T) {
				secret = []byte("")
				userId := uuid.New()
				email := "test@gmail.com"
				username := "testuser"
				role := "user"

				tokenString, err := CreateJWT(userId, email, username, role)

				require.Error(t, err)
				require.Empty(t, tokenString)

				// no secret key check comes before token check
				claims, err := GetClaims("")

				require.Error(t, err)
				require.Equal(t, "failed to get jwt secret key", err.Error())
				require.Nil(t, claims)
			},
		},
		{
			name: "failure - bad claims",
			test: func(t *testing.T) {
				claims := jwt.RegisteredClaims{}

				token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
				ss, _ := token.SignedString(secret)

				c, err := GetClaims(ss)
				require.Error(t, err)
				require.Nil(t, c)
			},
		},
	}

	for _, tc := range tcs {
		secret = []byte(os.Getenv("JWT_SECRET"))
		t.Run(tc.name, tc.test)
	}
}
