package impl

import (
	"mangahub/pkg/models"
	"mangahub/pkg/models/repository"

	"gorm.io/gorm"
)

type SessionRepositoryImpl struct {
	db *gorm.DB
}

var _ repository.SessionRepository = (*SessionRepositoryImpl)(nil)

func NewSessionRepository(db *gorm.DB) *SessionRepositoryImpl {
	return &SessionRepositoryImpl{db: db}
}

func (r *SessionRepositoryImpl) SaveSession(session *models.SessionModel) (*models.SessionModel, error) {
	result := r.db.Create(session)
	if result.Error != nil {
		return nil, result.Error
	}
	return session, nil
}

func (r *SessionRepositoryImpl) GetSessionByAccessToken(token string) (*models.SessionModel, error) {
	var session models.SessionModel
	result := r.db.Where("access_token = ?", token).First(&session)
	if result.Error != nil {
		return nil, result.Error
	}
	return &session, nil
}
