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
		utils.Error(w, http.StatusBadRequest, utils.ErrBadRequest, fmt.Sprintf("failed to decode data: %s", err.Error()))
		return
	}

	err := h.cfg.Validator.Struct(registerRequest)
	if err != nil {
		utils.Error(w, http.StatusUnprocessableEntity, utils.ErrValidationFailed, err.Error())
		return
	}

	res, jwtToken, refreshToken, sErr := h.registerService(r.Context(), registerRequest.Email, registerRequest.Username, registerRequest.Password)
	if sErr != nil {
		utils.SError(w, sErr)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    jwtToken,
		Path:     "/",
		MaxAge:   15 * 60,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
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
		utils.Error(w, http.StatusBadRequest, utils.ErrBadRequest, fmt.Sprintf("failed to decode data: %s", err.Error()))
		return
	}

	err := h.cfg.Validator.Struct(loginRequest)
	if err != nil {
		utils.Error(w, http.StatusUnprocessableEntity, utils.ErrValidationFailed, err.Error())
		return
	}

	res, jwtToken, refreshToken, sErr := h.loginService(r.Context(), loginRequest.Email, loginRequest.Password)
	if sErr != nil {
		utils.SError(w, sErr)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    jwtToken,
		Path:     "/",
		MaxAge:   15 * 60,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		MaxAge:   7 * 24 * 60,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	utils.Success(w, 201, res)
}

func (h *AuthHandler) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	refreshTokenCookie, err := r.Cookie("refresh_token")
	if err != nil || refreshTokenCookie.Value == "" {
		utils.Error(w, http.StatusBadRequest, utils.ErrBadRequest, "no refresh token cookie provided")
		return
	}
	refreshToken := refreshTokenCookie.Value

	jwtToken, newRefreshToken, res, sErr := h.refreshService(r.Context(), refreshToken)
	if sErr != nil {
		utils.SError(w, sErr)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    jwtToken,
		Path:     "/",
		MaxAge:   15 * 60,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		Path:     "/",
		MaxAge:   7 * 24 * 60,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	utils.Success(w, 201, res)
}

func (h *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	refreshTokenCookie, err := r.Cookie("refresh_token")
	if err != nil || refreshTokenCookie.Value == "" {
		utils.Error(w, http.StatusBadRequest, utils.ErrBadRequest, "no refresh token cookie provided")
		return
	}
	refreshToken := refreshTokenCookie.Value

	sErr := h.logoutService(r.Context(), refreshToken)
	if sErr != nil {
		utils.SError(w, sErr)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
	})

	utils.Success(w, 200, struct{}{})
}
