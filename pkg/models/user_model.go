package models

import "github.com/google/uuid"

type UserModel struct {
	ID string `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Username string `gorm:"type:varchar(255);unique;index" json:"username"`
	HashedPassword string `gorm:"type:varchar(255)" json:"-"`

	// Relationships
	Wishlists []WishlistModel `gorm:"foreignKey:CreatorID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"wishlists"`
	FollowingMangas []MangaModel `gorm:"many2many:manga_followers;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"following_mangas"`

	// BaseModel defines the basic structure and methods for all models.
	BaseModel `gorm:"embedded"`
	MetaUpdateModel	
}

func NewUserModel(username string, hashedPassword string) *UserModel {
	return &UserModel{
		ID: uuid.New().String(),
		Username: username,
		HashedPassword: hashedPassword,
	}
}

func (UserModel) TableName() string {
	return "users"
}