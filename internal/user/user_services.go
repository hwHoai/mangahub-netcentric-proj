package user

import "mangahub/pkg/dto"

type MeService interface {
	GetMe(userID string) (*dto.UserProfileResponse, dto.ExceptionDTO)
}
