package user

import "mangahub/pkg/dto"

type UserService interface {
	GetUserDetails(userID string) (*dto.UserProfileResponse, dto.ExceptionDTO)
}
