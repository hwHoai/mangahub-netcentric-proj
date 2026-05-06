package repository_impl

import (
	"mangahub/internal/repository"
	"mangahub/pkg/models"

	"gorm.io/gorm"
)

type MessageRepositoryImpl struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) repository.MessageRepository {
	return &MessageRepositoryImpl{db: db}
}

func (r *MessageRepositoryImpl) SaveMessage(message *models.MessageModel) error {
	return r.db.Create(message).Error
}

func (r *MessageRepositoryImpl) GetChatHistory(roomID string, limit int, offset int) ([]models.MessageModel, error) {
	var messages []models.MessageModel
	err := r.db.Where("room_id = ?", roomID).
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error
	return messages, err
}