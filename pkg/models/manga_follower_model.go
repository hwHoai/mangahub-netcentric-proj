package models

/* 
Need to register SetUpJoinTable() for many2many relationship between UserModel and MangaModel
	db.SetupJoinTable(&MangaModel{}, "Followers", &MangaFollowerModel{})
	db.SetupJoinTable(&UserModel{}, "FollowingMangas", &MangaFollowerModel{})
*/

type MangaFollowerModel struct {
	UserModelID string `gorm:"primaryKey;type:varchar(36)" json:"user_id"`
	MangaModelID string `gorm:"primaryKey;type:varchar(36);index" json:"manga_id"`

	//FK constraints
	User UserModel `gorm:"foreignKey:UserModelID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Manga MangaModel `gorm:"foreignKey:MangaModelID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`

	// BaseModel defines the basic structure and methods for all models.
	BaseModel `gorm:"embedded"`
}

func (MangaFollowerModel) TableName() string {
	return "manga_followers"
}