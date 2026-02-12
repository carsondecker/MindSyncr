package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/carsondecker/MindSyncr/utils"
	"github.com/go-playground/validator/v10"
)

type BaseHandler interface {
	GetConfig() *Config
}

func DecodeAndValidate[T any](r *http.Request, validator *validator.Validate) (T, *utils.ServiceError) {
	var data T
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return data, &utils.ServiceError{
			StatusCode: http.StatusBadRequest,
			Code:       ErrBadRequest,
			Message:    fmt.Sprintf("failed to decode data: %s", err.Error()),
		}
	}

	err := validator.Struct(data)
	if err != nil {
		return data, &utils.ServiceError{
			StatusCode: http.StatusUnprocessableEntity,
			Code:       ErrValidationFailed,
			Message:    err.Error(),
		}
	}

	return data, nil
}

func GetClaims(r *http.Request, validator *validator.Validate) (*Claims, *utils.ServiceError) {
	ctx := r.Context()

	raw := ctx.Value(UserContextKey)
	claims, ok := raw.(*Claims)
	if !ok || claims == nil {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusUnauthorized,
			Code:       ErrGetUserDataFail,
			Message:    "failed to get user claims from context",
		}
	}

	err := validator.Struct(claims)
	if err != nil {
		return nil, &utils.ServiceError{
			StatusCode: http.StatusUnprocessableEntity,
			Code:       ErrValidationFailed,
			Message:    err.Error(),
		}
	}

	return claims, nil
}

func BaseHandlerFuncWithBodyAndClaims[Req, Res any](
	h BaseHandler,
	w http.ResponseWriter,
	r *http.Request,
	successCode int,
	serviceFunc func(data Req, claims *Claims) (Res, *utils.ServiceError),
) {
	data, sErr := DecodeAndValidate[Req](r, h.GetConfig().Validator)
	if sErr != nil {
		utils.SError(w, sErr)
		return
	}

	claims, sErr := GetClaims(r, h.GetConfig().Validator)
	if sErr != nil {
		utils.SError(w, sErr)
		return
	}

	res, sErr := serviceFunc(data, claims)
	if sErr != nil {
		utils.SError(w, sErr)
		return
	}

	utils.Success(w, successCode, res)
}

func BaseHandlerFuncWithClaims[Res any](
	h BaseHandler,
	w http.ResponseWriter,
	r *http.Request,
	successCode int,
	serviceFunc func(claims *Claims) (Res, *utils.ServiceError),
) {
	claims, sErr := GetClaims(r, h.GetConfig().Validator)
	if sErr != nil {
		utils.SError(w, sErr)
		return
	}

	res, sErr := serviceFunc(claims)
	if sErr != nil {
		utils.SError(w, sErr)
		return
	}

	utils.Success(w, successCode, res)
}

func BaseHandlerFuncWithBody[Req, Res any](
	h BaseHandler,
	w http.ResponseWriter,
	r *http.Request,
	successCode int,
	serviceFunc func(data Req) (Res, *utils.ServiceError),
) {
	data, sErr := DecodeAndValidate[Req](r, h.GetConfig().Validator)
	if sErr != nil {
		utils.SError(w, sErr)
		return
	}

	res, sErr := serviceFunc(data)
	if sErr != nil {
		utils.SError(w, sErr)
		return
	}

	utils.Success(w, successCode, res)
}

func BaseHandlerFunc[Res any](
	h BaseHandler,
	w http.ResponseWriter,
	r *http.Request,
	successCode int,
	serviceFunc func() (Res, *utils.ServiceError),
) {
	res, sErr := serviceFunc()
	if sErr != nil {
		utils.SError(w, sErr)
		return
	}

	utils.Success(w, successCode, res)
}
