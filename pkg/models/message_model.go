package models

import "github.com/google/uuid"

type MessageModel struct {
	ID string `gorm:"primaryKey;type:varchar(36)" json:"id"`
	SenderID string `gorm:"type:varchar(36);index;" json:"sender_id"`
	RoomID string `gorm:"type:varchar(36);index;" json:"room_id"`
	Content string `gorm:"type:text" json:"content"`

	// FK constraints
	// Sender is the User who sent the message
	Sender UserModel `gorm:"foreignKey:SenderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	// RooomID is the MangaID - Manga is a chat room for its followers
	Room MangaModel `gorm:"foreignKey:RoomID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`

	// BaseModel defines the basic structure and methods for all models.
	BaseModel `gorm:"embedded"`
	MetaUpdateModel	`gorm:"embedded"`
}

func NewMessageModel(senderID string, roomID string, content string) *MessageModel {
	return &MessageModel{
		ID: uuid.New().String(),
		SenderID: senderID,
		RoomID: roomID,
		Content: content,
	}
}

func (MessageModel) TableName() string {
	return "messages"
}