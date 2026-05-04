package repository_impl

import (
	"mangahub/pkg/models"
	"mangahub/pkg/repository"

	"gorm.io/gorm"
)

type MangaRepositoryImpl struct {
	db *gorm.DB
}

var _ repository.MangaRepository = (*MangaRepositoryImpl)(nil)

func NewMangaRepository(db *gorm.DB) repository.MangaRepository {
	return &MangaRepositoryImpl{db: db}
}


func (r *MangaRepositoryImpl) GetMangas(limit int, offset int) ([]models.MangaModel, error) {
	var mangas []models.MangaModel
	err := r.db.Limit(limit).Offset(offset).Find(&mangas).Error
	if err != nil {
		return nil, err
	}
	return mangas, nil
}

func (r *MangaRepositoryImpl) GetMangaDetail(mangaID string) (*models.MangaModel, error) {
	var manga models.MangaModel
	err := r.db.First(&manga, "id = ?", mangaID).Error
	if err != nil {
		return nil, err
	}
	return &manga, nil
}

func (r *MangaRepositoryImpl) GetMangaChapters(mangaID string) ([]models.ChapterModel, error) {
	var chapters []models.ChapterModel
	err := r.db.Where("manga_id = ?", mangaID).Order("chapter_number ASC").Find(&chapters).Error
	if err != nil {
		return nil, err	
	}
	return chapters, nil
}