package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type SuccessResponse struct {
	Success bool `json:"success"`
	Data    any  `json:"data"`
}

func Success[Data any](w http.ResponseWriter, statusCode int, data Data) error {
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
