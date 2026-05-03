package repository_impl

import (
	"mangahub/pkg/models"
	"mangahub/pkg/repository"

	"gorm.io/gorm"
)

type SessionRepositoryImpl struct {
	db *gorm.DB
}

var _ repository.SessionRepository = (*SessionRepositoryImpl)(nil)

func NewSessionRepository(db *gorm.DB) repository.SessionRepository {
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

func (r *SessionRepositoryImpl) UpdateSessionByUserID(userID, accessToken, refreshToken string) (*models.SessionModel, error) {
	var session models.SessionModel
	result := r.db.Where("user_id = ?", userID).First(&session)
	if result.Error != nil {
		return nil, result.Error
	}

	session.AccessToken = accessToken
	session.RefreshToken = refreshToken

	if err := r.db.Save(&session).Error; err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *SessionRepositoryImpl) GetSessionByUserID(userID string) (*models.SessionModel, error) {
	var session models.SessionModel
	result := r.db.Where("user_id = ?", userID).First(&session)
	if result.Error != nil {
		return nil, result.Error
	}
	return &session, nil
}
