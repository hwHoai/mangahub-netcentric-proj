package repository

import "mangahub/pkg/models"

type SessionRepository interface {
	SaveSession(session *models.SessionModel) (*models.SessionModel, error)
	UpdateSessionByUserID(userID, accessToken, refreshToken string) (*models.SessionModel, error)
	GetSessionByAccessToken(token string) (*models.SessionModel, error)
	GetSessionByUserID(userID string) (*models.SessionModel, error)
	GetSessionByRefreshToken(token string) (*models.SessionModel, error)
}