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
		return
	}

	res, jwtToken, refreshToken, sErr := h.registerService(r.Context(), registerRequest.Email, registerRequest.Username, registerRequest.Password)
	if sErr != nil {
		utils.SError(w, sErr)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwtToken",
		Value:    jwtToken,
		Path:     "/",
		MaxAge:   15 * 60,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		Path:     "/",
		MaxAge:   7 * 24 * 60,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	utils.Success(w, 201, res)
}

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var loginRequest LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		utils.Error(w, 400, "BAD_REQUEST", fmt.Sprintf("failed to decode data: %s", err.Error()))
		return
	}

	err := h.cfg.Validator.Struct(loginRequest)
	if err != nil {
		utils.Error(w, 422, "VALIDATION_FAIL", err.Error())
		return
	}

	res, jwtToken, refreshToken, sErr := h.loginService(r.Context(), loginRequest.Email, loginRequest.Password)
	if sErr != nil {
		utils.SError(w, sErr)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwtToken",
		Value:    jwtToken,
		Path:     "/",
		MaxAge:   15 * 60,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		Path:     "/",
		MaxAge:   7 * 24 * 60,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	utils.Success(w, 201, res)
}
