package dto

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" form:"refresh_token" binding:"required"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserProfileResponse struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
