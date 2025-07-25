package usecase

import (
	"github.com/google/uuid"
	"savebite/internal/app/user/repository"
	"savebite/internal/domain/dto"
	"savebite/internal/domain/entity"
	"savebite/pkg/jwt"
	"savebite/pkg/oauth"
)

type AuthUsecaseItf interface {
	Redirect(state string) string
	Callback(data dto.OAuthCallbackRequest) (string, error)
}

type AuthUsecase struct {
	userRepo    repository.UserRepoItf
	googleOAuth oauth.GoogleOAuthItf
	jwt         jwt.CustomJWTItf
}

func NewAuthUsecase(userRepo repository.UserRepoItf, o oauth.GoogleOAuthItf, j jwt.CustomJWTItf) AuthUsecaseItf {
	return &AuthUsecase{
		userRepo:    userRepo,
		googleOAuth: o,
		jwt:         j,
	}
}

func (u *AuthUsecase) Redirect(state string) string {
	return u.googleOAuth.GetGoogleLogin(state)
}

func (u *AuthUsecase) Callback(data dto.OAuthCallbackRequest) (string, error) {
	token, err := u.googleOAuth.ExchangeToken(data.Code)
	if err != nil {
		return "", err
	}

	userInfo, err := u.googleOAuth.GetUserInfo(token)
	if err != nil {
		return "", err
	}

	user, err := u.userRepo.Show(dto.UserParam{Email: userInfo.Email})
	if err != nil {
		return "", err
	}

	if user.ID == uuid.Nil {
		user = &entity.User{
			ID:    uuid.New(),
			Name:  userInfo.Name,
			Email: userInfo.Email,
		}

		if err := u.userRepo.Create(user); err != nil {
			return "", err
		}
	}

	jwtToken, err := u.jwt.Create(user.ID, user.Name, user.Email)
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}
