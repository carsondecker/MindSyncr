package auth

import (
	"net/http"

	"github.com/carsondecker/MindSyncr/internal/sutils"
	"github.com/carsondecker/MindSyncr/utils"
)

type AuthService struct {
	repo AuthRepository
}

func NewAuthService(cfg *sutils.Config) *AuthService {
	return &AuthService{}
}

func (h *AuthService) GetConfig() *sutils.Config {
	return h.cfg
}

func (h *AuthService) HandleRegister(w http.ResponseWriter, r *http.Request) {
	sutils.BaseHandlerFuncWithBody(h, w, r,
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

func (h *AuthService) HandleLogin(w http.ResponseWriter, r *http.Request) {
	sutils.BaseHandlerFuncWithBody(h, w, r,
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

func (h *AuthService) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	sutils.BaseHandlerFunc(h, w, r,
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

func (h *AuthService) HandleLogout(w http.ResponseWriter, r *http.Request) {
	sutils.BaseHandlerFuncWithClaims(h, w, r,
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

func (h *AuthService) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	sutils.BaseHandlerFuncWithClaims(h, w, r,
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
