package dto

type OAuthCallbackRequest struct {
	Code  string `json:"code" validate:"required"`
	State string `json:"state" validate:"required"`
	Error string `json:"error"`
}

type GoogleUserProfileResponse struct {
	Email string
	Name  string
}
