package oauth

import (
	"context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	oauth2api "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
	"savebite/internal/domain/dto"
	"savebite/internal/domain/env"
	"savebite/pkg/log"
)

type GoogleOAuthItf interface {
	GetGoogleLogin(state string) string
	ExchangeToken(code string) (*oauth2.Token, error)
	GetUserInfo(token *oauth2.Token) (dto.GoogleUserProfileResponse, error)
}

type GoogleOAuthStruct struct {
	googleConfig *oauth2.Config
}

var GoogleOAuth = getGoogleOAuth()

func getGoogleOAuth() GoogleOAuthItf {
	config := &oauth2.Config{
		ClientID:     env.AppEnv.GoogleClientID,
		ClientSecret: env.AppEnv.GoogleClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  env.AppEnv.GoogleRedirectURL,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
	}

	return &GoogleOAuthStruct{googleConfig: config}
}

func (o *GoogleOAuthStruct) GetGoogleLogin(state string) string {
	return o.googleConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (o *GoogleOAuthStruct) ExchangeToken(code string) (*oauth2.Token, error) {
	token, err := o.googleConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[OAuth][ExchangeToken] Failed to exchange token")
		return nil, err
	}

	return token, nil
}

func (o *GoogleOAuthStruct) GetUserInfo(token *oauth2.Token) (dto.GoogleUserProfileResponse, error) {
	userData := dto.GoogleUserProfileResponse{}
	client := o.googleConfig.Client(context.Background(), token)
	oauth2Service, err := oauth2api.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[OAuth][GetUserInfo] Failed to create service")
		return userData, err
	}

	userInfo, err := oauth2Service.Userinfo.Get().Do()
	if err != nil {
		log.Error(log.LogInfo{
			"error": err.Error(),
		}, "[OAuth][GetUserInfo] Failed to retrieve user info")
		return userData, err
	}

	userData.Name = userInfo.Name
	userData.Email = userInfo.Email

	return userData, nil
}
