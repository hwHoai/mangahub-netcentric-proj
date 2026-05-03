package models

import (
	"mangahub/pkg/models/enums"

	"github.com/google/uuid"
)

type MangaModel struct {
	ID string `gorm:"primaryKey;type:varchar(36)" json:"id"`

	// Manga fields
	Title string `gorm:"type:varchar(255);index" json:"title"`
	Author string `gorm:"type:varchar(255);index" json:"author"`
	TotalChapters int `gorm:"type:int" json:"total_chapters"`
	Description string `gorm:"type:text" json:"description"`
	CoverURL string `gorm:"type:varchar(255)" json:"cover_url"`

	// Status of the manga (e.g., "comming_soon", "in_progress", "completed")
	Status enums.MangaStatus `gorm:"type:varchar(50)" json:"status"`

	// Relationships
	Reviews []ReviewModel `gorm:"foreignKey:MangaID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Genres []GenresModel `gorm:"many2many:manga_genres;" json:"genres"`
	Followers []UserModel `gorm:"many2many:manga_followers;" json:"-"`
	Chapters  []ChapterModel `gorm:"foreignKey:MangaID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"chapters"`

	// Messages is the chat history for this manga's followers
	Messages []MessageModel `gorm:"foreignKey:RoomID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`

	// Meta fields
	BaseModel `gorm:"embedded"`
	MetaUpdateModel	`gorm:"embedded"`
}

func NewMangaModel(title string, author string, totalChapters int, description string, coverURL string, status enums.MangaStatus) *MangaModel {
	return &MangaModel{
		ID: uuid.New().String(),
		Title: title,
		Author: author,
		TotalChapters: totalChapters,
		Description: description,
		CoverURL: coverURL,
		Status: status,
	}
}

func (MangaModel) TableName() string {
	return "mangas"
}