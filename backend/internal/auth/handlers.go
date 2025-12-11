package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/carsondecker/MindSyncr/internal/config"
	"github.com/carsondecker/MindSyncr/internal/utils"
)

type AuthHandler struct {
	cfg *config.Config
}

func NewAuthHandler(cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		cfg,
	}
}

func (h *AuthHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var registerRequest RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&registerRequest); err != nil {
		utils.Error(w, 400, "BAD_REQUEST", fmt.Sprintf("failed to decode data: %s", err.Error()))
		return
	}

	err := h.cfg.Validator.Struct(registerRequest)
	if err != nil {
		utils.Error(w, 422, "VALIDATION_FAIL", err.Error())
	}

	row, sErr := h.registerService(r.Context(), registerRequest.Email, registerRequest.Username, registerRequest.Password)
	if sErr != nil {
		utils.SError(w, sErr)
	}

	utils.Success(w, 201, row)
}
