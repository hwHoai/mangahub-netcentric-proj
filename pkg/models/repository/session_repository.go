package repository

import "mangahub/pkg/models"

type SessionRepository interface {
	SaveSession(session *models.SessionModel) (*models.SessionModel, error)
	GetSessionByAccessToken(token string) (*models.SessionModel, error)
}