package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

var secret []byte
var wsSecret []byte

type Claims struct {
	UserId   uuid.UUID `json:"id" validate:"required,not_nil_uuid"`
	Email    string    `json:"email" validate:"required,email"`
	Username string    `json:"username" validate:"required,min=1"`
	Role     string    `json:"role" validate:"required,min=1"`
	jwt.RegisteredClaims
}

type WSClaims struct {
	UserId    uuid.UUID `json:"id" validate:"required,not_nil_uuid"`
	SessionId uuid.UUID `json:"session_id" validate:"required,not_nil_uuid"`
	jwt.RegisteredClaims
}

func JWTInit(s, ws string) {
	secret = []byte(s)
	wsSecret = []byte(ws)
}

func CreateJWT(userId uuid.UUID, email, username, role string) (string, error) {
	if len(secret) == 0 {
		return "", fmt.Errorf("failed to get jwt secret key")
	}

	claims := Claims{
		UserId:   userId,
		Email:    email,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "MindSyncr",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	ss, err := token.SignedString(secret)

	if err != nil {
		return "", err
	}

	return ss, nil
}

func CreateWSJWT(userId, sessionId uuid.UUID) (string, error) {
	if len(wsSecret) == 0 {
		return "", fmt.Errorf("failed to get jwt secret key")
	}

	claims := WSClaims{
		UserId:    userId,
		SessionId: sessionId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(3 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "MindSyncr",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	ss, err := token.SignedString(wsSecret)

	if err != nil {
		return "", err
	}

	return ss, nil
}

func GetClaimsFromToken(tokenString string) (*Claims, error) {
	if len(secret) == 0 {
		return nil, fmt.Errorf("failed to get jwt secret key")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("could not parse jwt claims")
	}

	if claims.UserId == uuid.Nil || claims.Email == "" || claims.Username == "" {
		return nil, fmt.Errorf("invalid claims: missing required fields")
	}

	return claims, nil
}
