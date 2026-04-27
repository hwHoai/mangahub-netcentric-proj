package impl

import (
	"mangahub/pkg/models"
	"mangahub/pkg/models/repository"

	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

var _ repository.UserRepository = (*UserRepositoryImpl)(nil)

func NewUserRepository(db *gorm.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) GetUserByUsername(username string) (*models.UserModel, error) {
	var userModel models.UserModel
	result := r.db.Where("username = ?", username).First(&userModel)
	if result.Error != nil {
		return nil, result.Error
	}
	return &userModel, nil
}

func (r *UserRepositoryImpl) CreateUser(user *models.UserModel) (*models.UserModel, error) {
	result := r.db.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}