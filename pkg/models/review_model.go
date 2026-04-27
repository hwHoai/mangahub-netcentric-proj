package models

import "github.com/google/uuid"

type ReviewModel struct {
	ID string `gorm:"primaryKey;type:varchar(36)" json:"id"`

	// Foreign keys
	UserID string `gorm:"type:varchar(36);index;" json:"user_id"`
	MangaID string `gorm:"type:varchar(36);index;" json:"manga_id"`

	// Review fields
	Rating int `gorm:"type:int" json:"rating"`
	Content string `gorm:"type:text" json:"content"`

	//FK constraint
	User UserModel `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Manga MangaModel `gorm:"foreignKey:MangaID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	
	// BaseModel defines the basic structure and methods for all models.
	BaseModel `gorm:"embedded"`
	MetaUpdateModel	`gorm:"embedded"`
}

func NewReviewModel(userID string, mangaID string, rating int, content string) *ReviewModel {
	return &ReviewModel{
		ID: uuid.New().String(),
		UserID: userID,
		MangaID: mangaID,
		Rating: rating,
		Content: content,
	}
}

func (ReviewModel) TableName() string {
	return "reviews"
}