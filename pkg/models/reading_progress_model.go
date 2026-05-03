package models

import (
	"mangahub/pkg/models/enums"
)

type ReadingProgressModel struct {
	// Foreign keys
	UserID string `gorm:"primaryKey;type:varchar(36);" json:"user_id"`
	MangaID string `gorm:"primaryKey;type:varchar(36);index" json:"manga_id"`

	//FK constraints
	User UserModel `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Manga MangaModel `gorm:"foreignKey:MangaID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	
	// Reading progress fields
	Status enums.ReadingStatus `gorm:"type:varchar(20)" json:"status"`
	CurrentChapter int `gorm:"type:int" json:"current_chapter"`

	// Meta fields
	BaseModel 
	MetaUpdateModel `gorm:"embedded"`
}

func NewReadingProgressModel(userID string, mangaID string, status enums.ReadingStatus, currentChapter int) *ReadingProgressModel {
	return &ReadingProgressModel{
		UserID: userID,
		MangaID: mangaID,
		Status: status,
		CurrentChapter: currentChapter,
	}
}

func (ReadingProgressModel) TableName() string {
	return "reading_progress"
}