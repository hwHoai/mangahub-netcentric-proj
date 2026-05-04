package auth

import "mangahub/pkg/dto"

type LoginService interface {
	// LoginByUsername authenticates a user using their username and password.
	LoginByUsername(request *dto.LoginByUsernameRequest) (*dto.LoginByUsernameResponse, dto.ExceptionDTO)
}