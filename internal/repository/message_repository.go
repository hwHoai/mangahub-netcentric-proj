package repository

import "mangahub/pkg/models"

type MessageRepository interface {
	SaveMessage(message *models.MessageModel) error
	GetChatHistory(roomID string, limit int, offset int) ([]models.MessageModel, error)
}