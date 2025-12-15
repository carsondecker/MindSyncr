package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type BaseHandler interface {
	GetConfig() *Config
}

func DecodeAndValidate[T any](r *http.Request, validator *validator.Validate) (T, *ServiceError) {
	var data T
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return data, &ServiceError{
			StatusCode: http.StatusBadRequest,
			Code:       ErrBadRequest,
			Message:    fmt.Sprintf("failed to decode data: %s", err.Error()),
		}
	}

	err := validator.Struct(data)
	if err != nil {
		return data, &ServiceError{
			StatusCode: http.StatusUnprocessableEntity,
			Code:       ErrValidationFailed,
			Message:    err.Error(),
		}
	}

	return data, nil
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

func BaseHandlerWithClaimsFunc[T, R any](
	h BaseHandler,
	w http.ResponseWriter,
	r *http.Request,
	successCode int,
	serviceFunc func(data T, claims *Claims) (R, *ServiceError),
) {
	data, sErr := DecodeAndValidate[T](r, h.GetConfig().Validator)
	if sErr != nil {
		SError(w, sErr)
		return
	}

	claims, sErr := GetClaims(r, h.GetConfig().Validator)
	if sErr != nil {
		SError(w, sErr)
		return
	}

	res, sErr := serviceFunc(data, claims)
	if sErr != nil {
		SError(w, sErr)
		return
	}

	Success(w, successCode, res)
}
