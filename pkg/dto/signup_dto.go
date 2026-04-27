package dto

type SignupByUsernameRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SignupByUsernameResponse struct {
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn int64 `json:"expires_in"`
}