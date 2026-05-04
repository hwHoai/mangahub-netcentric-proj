package models

import (
	"github.com/google/uuid"
)

type ChapterModel struct {
	ID            string  `gorm:"primaryKey;type:varchar(36)" json:"id"`
	MangaID       string  `gorm:"type:varchar(36);index;not null" json:"manga_id"`
	ChapterNumber float64 `gorm:"type:decimal(10,2)" json:"chapter_number"`
	Title         string  `gorm:"type:varchar(255)" json:"title"`
	PagesData     string  `gorm:"type:text" json:"pages_data"` // JSON array of page URLs

	BaseModel       `gorm:"embedded"`
	MetaUpdateModel `gorm:"embedded"`
}

func NewChapterModel(mangaID string, chapterNumber float64, title string, pagesData string) *ChapterModel {
	return &ChapterModel{
		ID:            uuid.New().String(),
		MangaID:       mangaID,
		ChapterNumber: chapterNumber,
		Title:         title,
		PagesData:     pagesData,
	}
}

func (ChapterModel) TableName() string {
	return "chapters"
}
