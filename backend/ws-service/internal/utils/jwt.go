package utils

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

var secret []byte

type Claims struct {
	UserId    uuid.UUID `json:"id" validate:"required,not_nil_uuid"`
	SessionId uuid.UUID `json:"session_id" validate:"required,not_nil_uuid"`
	jwt.RegisteredClaims
}

func JWTInit(s string) {
	secret = []byte(s)
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

	if claims.UserId == uuid.Nil || claims.SessionId == uuid.Nil {
		return nil, fmt.Errorf("invalid claims: missing required fields")
	}

	return claims, nil
}

func GetClaims(r *http.Request, validator *validator.Validate) (*Claims, *ServiceError) {
	ctx := r.Context()

	raw := ctx.Value(UserContextKey)
	claims, ok := raw.(*Claims)
	if !ok || claims == nil {
		return nil, &ServiceError{
			StatusCode: http.StatusUnauthorized,
			Code:       ErrGetUserDataFail,
			Message:    "failed to get user claims from context",
		}
	}

	err := validator.Struct(claims)
	if err != nil {
		return nil, &ServiceError{
			StatusCode: http.StatusUnprocessableEntity,
			Code:       ErrValidationFailed,
			Message:    err.Error(),
		}
	}

	return claims, nil
}
