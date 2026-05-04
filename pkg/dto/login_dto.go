package dto

// =================================================================
// Username and Password
// =================================================================
type LoginByUsernameRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginByUsernameResponse struct {
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn int64 `json:"expires_in"`
}