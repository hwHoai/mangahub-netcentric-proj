package auth

import "mangahub/pkg/dto"

type SignupService interface {
	// SignupByUsername registers a new user using their username and password.
	SignupByUsername(request *dto.SignupByUsernameRequest) (*dto.SignupByUsernameResponse, dto.ExceptionDTO)
}