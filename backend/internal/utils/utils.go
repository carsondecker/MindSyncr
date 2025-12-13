package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type SuccessResponse struct {
	Success bool `json:"success"`
	Data    any  `json:"data"`
}

type ErrorResponse struct {
	Success bool `json:"success"`
	Error   Err  `json:"error"`
}

type Err struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ServiceError struct {
	StatusCode int
	Code       string
	Message    string
}

func WriteSuccess[Data any](w http.ResponseWriter, statusCode int, data Data) error {
	if statusCode >= 400 {
		return fmt.Errorf("cannot use success with an error code")
	}

	res := SuccessResponse{
		Success: true,
		Data:    data,
	}

	b, err := json.Marshal(res)

	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(b)

	if err != nil {
		return fmt.Errorf("failed to write response: %w", err)
	}

	return nil
}

func WriteError(w http.ResponseWriter, statusCode int, code string, incErr string) error {
	if statusCode < 300 {
		return fmt.Errorf("cannot use error with a success code")
	}

	res := ErrorResponse{
		Success: false,
		Error: Err{
			Code:    code,
			Message: incErr,
		},
	}

	b, err := json.Marshal(res)

	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(b)

	if err != nil {
		return fmt.Errorf("failed to write response: %w", err)
	}

	return nil
}

func Success[Data any](w http.ResponseWriter, statusCode int, data Data) {
	err := WriteSuccess(w, statusCode, data)
	if err != nil {
		Error(w, 500, "BAD_RESPONSE", err.Error())
	}
}

func Error(w http.ResponseWriter, statusCode int, code string, msg string) {
	err := WriteError(w, statusCode, code, msg)
	if err != nil {
		log.Printf("error: %v", err)
	}
}

func SError(w http.ResponseWriter, se *ServiceError) {
	Error(w, se.StatusCode, se.Code, se.Message)
}
