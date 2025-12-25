package auth

import (
	"net/http"

	"github.com/carsondecker/MindSyncr/internal/utils"
)

type AuthHandler struct {
	cfg *utils.Config
}

func NewAuthHandler(cfg *utils.Config) *AuthHandler {
	return &AuthHandler{
		cfg,
	}
}

func (h *AuthHandler) GetConfig() *utils.Config {
	return h.cfg
}

func (h *AuthHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	utils.BaseHandlerFuncWithBody(h, w, r,
		http.StatusCreated,
		func(data RegisterRequest) (UserWithRefresh, *utils.ServiceError) {
			res, jwtToken, refreshToken, sErr := h.registerService(r.Context(), data.Email, data.Username, data.Password)
			if sErr != nil {
				return UserWithRefresh{}, sErr
			}

			// TODO: change samesite to strict and secure to true
			http.SetCookie(w, &http.Cookie{
				Name:     "access_token",
				Value:    jwtToken,
				Path:     "/",
				MaxAge:   15 * 60,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteNoneMode,
			})

			http.SetCookie(w, &http.Cookie{
				Name:     "refresh_token",
				Value:    refreshToken,
				Path:     "/",
				MaxAge:   7 * 24 * 60,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteNoneMode,
			})

			return res, nil
		},
	)
}

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	utils.BaseHandlerFuncWithBody(h, w, r,
		http.StatusOK,
		func(data LoginRequest) (UserWithRefresh, *utils.ServiceError) {
			res, jwtToken, refreshToken, sErr := h.loginService(r.Context(), data.Email, data.Password)
			if sErr != nil {
				return UserWithRefresh{}, sErr
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "access_token",
				Value:    jwtToken,
				Path:     "/",
				MaxAge:   15 * 60,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteNoneMode,
			})

			http.SetCookie(w, &http.Cookie{
				Name:     "refresh_token",
				Value:    refreshToken,
				Path:     "/",
				MaxAge:   7 * 24 * 60,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteNoneMode,
			})

			return res, nil
		},
	)
}

func (h *AuthHandler) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	utils.BaseHandlerFunc(h, w, r,
		http.StatusCreated,
		func() (RefreshTokenResponse, *utils.ServiceError) {
			refreshToken, sErr := getRefreshToken(r)
			if sErr != nil {
				return RefreshTokenResponse{}, sErr
			}

			jwtToken, newRefreshToken, res, sErr := h.refreshService(r.Context(), refreshToken)
			if sErr != nil {
				return RefreshTokenResponse{}, sErr
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "access_token",
				Value:    jwtToken,
				Path:     "/",
				MaxAge:   15 * 60,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteNoneMode,
			})

			http.SetCookie(w, &http.Cookie{
				Name:     "refresh_token",
				Value:    newRefreshToken,
				Path:     "/",
				MaxAge:   7 * 24 * 60,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteNoneMode,
			})

			return res, nil
		},
	)
}

func (h *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	utils.BaseHandlerFuncWithClaims(h, w, r,
		http.StatusOK,
		func(claims *utils.Claims) (struct{}, *utils.ServiceError) {
			refreshToken, sErr := getRefreshToken(r)
			if sErr != nil {
				return struct{}{}, sErr
			}

			sErr = h.logoutService(r.Context(), claims.UserId, refreshToken)
			if sErr != nil {
				return struct{}{}, sErr
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

			return struct{}{}, nil
		},
	)
}

func (h *AuthHandler) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	utils.BaseHandlerFuncWithClaims(h, w, r,
		http.StatusOK,
		func(claims *utils.Claims) (User, *utils.ServiceError) {
			res, sErr := h.getUserService(r.Context(), claims.UserId)
			if sErr != nil {
				return User{}, sErr
			}

			return res, nil
		},
	)
}
