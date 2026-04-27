package models

import "github.com/google/uuid"

type SessionModel struct {
	ID string `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UserID string `gorm:"type:varchar(36);index;" json:"user_id"`
	AccessToken string `gorm:"type:varchar(255);unique;index" json:"access_token"`
	RefreshToken string `gorm:"type:varchar(255);unique;index" json:"refresh_token"`
	PublicKey string `gorm:"type:varchar(255)" json:"public_key"`

	//FK constraint
	User UserModel `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`

	// BaseModel defines the basic structure and methods for all models.
	BaseModel 
	MetaUpdateModel	`gorm:"embedded"`
}

func NewSessionModel(userID string, AccessToken string, RefreshToken string, PublicKey string) SessionModel {
	return SessionModel{
		ID: uuid.New().String(),
		UserID: userID,
		RefreshToken: RefreshToken,
		AccessToken: AccessToken,
		PublicKey: PublicKey,
		BaseModel: BaseModel{},
		MetaUpdateModel: MetaUpdateModel{},
	}
}

func (SessionModel) TableName() string {
	return "sessions"
}