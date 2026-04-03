package auth

import (
	"net/http"

	"github.com/carsondecker/MindSyncr/utils"
	"github.com/google/uuid"
)

func (s *AuthService) registerService(email, username, password string) (UserWithRefresh, string, string, *utils.ServiceError) {
	passwordHash, err := hashPassword(password)
	if err != nil {
		return UserWithRefresh{}, "", "", &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrHashFail,
			Message:    err.Error(),
		}
	}

	user, sErr := s.repo.Register(email, username, passwordHash)

	if sErr != nil {
		return UserWithRefresh{}, "", "", sErr
	}

	jwtToken, err := utils.CreateJWT(user.Id, user.Email, user.Username, user.Role)
	if err != nil {
		return UserWithRefresh{}, "", "", &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrJwtFail,
			Message:    err.Error(),
		}
	}

	refreshToken, refreshRes, sErr := s.createRefreshToken(user.Id)
	if sErr != nil {
		return UserWithRefresh{}, "", "", sErr
	}

	res := UserWithRefresh{
		User:         user,
		RefreshToken: refreshRes,
	}

	return res, jwtToken, refreshToken, nil
}

func (s *AuthService) loginService(email, password string) (UserWithRefresh, string, string, *utils.ServiceError) {
	internalUser, sErr := s.repo.GetInternalUser(email)
	if sErr != nil {
		return UserWithRefresh{}, "", "", sErr
	}

	err := checkPassword(internalUser.PasswordHash, password)
	if err != nil {
		return UserWithRefresh{}, "", "", &utils.ServiceError{
			StatusCode: http.StatusUnauthorized,
			Code:       utils.ErrInvalidCredentials,
			Message:    err.Error(),
		}
	}

	jwtToken, err := utils.CreateJWT(internalUser.Id, internalUser.Email, internalUser.Username, internalUser.Role)
	if err != nil {
		return UserWithRefresh{}, "", "", &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrJwtFail,
			Message:    err.Error(),
		}
	}

	refreshToken, refreshRes, sErr := s.createRefreshToken(internalUser.Id)
	if sErr != nil {
		return UserWithRefresh{}, "", "", sErr
	}

	res := UserWithRefresh{
		User:         internalUser.User,
		RefreshToken: refreshRes,
	}

	return res, jwtToken, refreshToken, nil
}

func (s *AuthService) refreshService(token string) (string, string, RefreshTokenResponse, *utils.ServiceError) {
	tx, sErr := s.repo.BeginTx()
	if sErr != nil {
		return "", "", RefreshTokenResponse{}, sErr
	}

	defer tx.Rollback()

	txRepo := s.repo.GetTxRepo(tx)

	txAuthService := AuthService{
		repo: txRepo,
	}

	userId, sErr := txAuthService.isValidRefreshToken(token)
	if sErr != nil {
		return "", "", RefreshTokenResponse{}, sErr
	}

	sErr = txRepo.RevokeUserTokens(userId)
	if sErr != nil {
		return "", "", RefreshTokenResponse{}, sErr
	}

	refreshToken, res, sErr := s.createRefreshToken(userId)
	if sErr != nil {
		return "", "", RefreshTokenResponse{}, sErr
	}

	jwtToken, sErr := s.createJWTById(userId)
	if sErr != nil {
		return "", "", RefreshTokenResponse{}, sErr
	}

	err := tx.Commit()
	if err != nil {
		return "", "", RefreshTokenResponse{}, &utils.ServiceError{
			StatusCode: http.StatusInternalServerError,
			Code:       utils.ErrTxCommitFail,
			Message:    err.Error(),
		}
	}

	return jwtToken, refreshToken, res, nil
}

func (s *AuthService) logoutService(userId uuid.UUID, token string) *utils.ServiceError {
	_, sErr := s.isValidRefreshToken(token)
	if sErr != nil {
		return sErr
	}

	sErr = s.repo.RevokeUserTokens(userId)
	if sErr != nil {
		return sErr
	}

	return nil
}

func (s *AuthService) getUserService(userId uuid.UUID) (User, *utils.ServiceError) {
	user, sErr := s.repo.GetUserById(userId)
	if sErr != nil {
		return User{}, sErr
	}

	return user, nil
}
