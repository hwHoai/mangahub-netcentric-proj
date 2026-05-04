package models

import "github.com/google/uuid"

type GenresModel struct {
	ID string `gorm:"primaryKey" json:"id"`

	// Genre fields
	Name string `gorm:"type:varchar(100);uniqueIndex" json:"name"`
	Description string `gorm:"type:text" json:"description"`

	// Relationships
	Mangas []MangaModel `gorm:"many2many:manga_genres;" json:"-"`

	BaseModel
	MetaUpdateModel	`gorm:"embedded"`
}

func NewGenresModel(name string, description string) *GenresModel {
	return &GenresModel{
		ID: uuid.New().String(),
		Name: name,
		Description: description,
	}
}	

func (GenresModel) TableName() string {
	return "genres"
}