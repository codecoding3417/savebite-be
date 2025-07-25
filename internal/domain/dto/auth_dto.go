package dto

//type OAuthCallbackRequest struct {
//	Code  string `json:"code" validate:"required"`
//	State string `json:"state" validate:"required"`
//	Error string `json:"error"`
//}

type OAuthCallbackRequest struct {
	Code  string `json:"code" query:"code" validate:"required"`
	State string `json:"state" query:"state" validate:"required"`
	Error string `json:"error" query:"error"`
}

type GoogleUserProfileResponse struct {
	Email string
	Name  string
}
