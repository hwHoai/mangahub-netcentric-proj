package repository

import "mangahub/pkg/models"

type UserRepository interface {
	GetUserByUsername(username string) (*models.UserModel, error)
	CreateUser(user *models.UserModel) (*models.UserModel, error)
}