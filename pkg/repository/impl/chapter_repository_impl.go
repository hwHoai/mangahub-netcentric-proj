package repository_impl

import (
	"mangahub/pkg/models"
	"mangahub/pkg/repository"

	"gorm.io/gorm"
)

type ChapterRepositoryImpl struct {
	db *gorm.DB
}

var _ repository.ChapterRepository = (*ChapterRepositoryImpl)(nil)

func NewChapterRepositoryImpl(db *gorm.DB) repository.ChapterRepository {
	return &ChapterRepositoryImpl{db: db}
}


func (c *ChapterRepositoryImpl) GetChapterDataByMangaID(mangaID string) ([]models.ChapterModel, error) {
	var chapterList []models.ChapterModel
	err := c.db.Where("manga_id = ?", mangaID).Find(&chapterList).Error
	if err != nil {
		return nil, err
	}
	return chapterList, nil	
}

func (c *ChapterRepositoryImpl) GetChapterByID(chapterID string) (*models.ChapterModel, error) {
	var chapter models.ChapterModel
	err := c.db.First(&chapter, "id = ?", chapterID).Error
	if err != nil {
		return nil, err
	}
	return &chapter, nil
}