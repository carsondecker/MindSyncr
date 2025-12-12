package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

var secret []byte

type Claims struct {
	Id       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
	Role     string    `json:"role"`
	jwt.RegisteredClaims
}

func JWTInit() {
	secret = []byte(os.Getenv("JWT_SECRET"))
}

func CreateJWT(id uuid.UUID, email, username string) (string, error) {
	if secret == nil || len(secret) == 0 {
		return "", fmt.Errorf("failed to get jwt secret key")
	}

	claims := Claims{
		Id:       id,
		Email:    email,
		Username: username,
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

func GetClaims(tokenString string) (*Claims, error) {
	if secret == nil || len(secret) == 0 {
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

	if claims.Id == uuid.Nil || claims.Email == "" || claims.Username == "" {
		return nil, fmt.Errorf("invalid claims: missing required fields")
	}

	return claims, nil
}
