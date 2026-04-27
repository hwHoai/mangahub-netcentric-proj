package models

import "github.com/google/uuid"

type WishlistModel struct {
	ID string `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Name string `gorm:"type:varchar(255);index" json:"name"`
	CreatorID string `gorm:"type:varchar(36);index" json:"creator_id"`

	//FK constraints
	Creator UserModel `gorm:"foreignKey:CreatorID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Mangas []MangaModel `gorm:"many2many:wishlist_mangas;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"mangas"`
	Subscribers []UserModel `gorm:"many2many:wishlist_subscribers;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"subscribers"`

	// BaseModel defines the basic structure and methods for all models.
	BaseModel `gorm:"embedded"`
	MetaUpdateModel	`gorm:"embedded"`
}

func NewWishlistModel(name string, creatorID string) *WishlistModel {
	return &WishlistModel{
		ID: uuid.New().String(),
		Name: name,
		CreatorID: creatorID,
	}
}

func (WishlistModel) TableName() string {
	return "wishlists"
}