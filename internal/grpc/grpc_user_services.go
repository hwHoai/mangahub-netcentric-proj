package grpc

import (
	"context"
	"mangahub/proto/user"
)

type GRPCUserService interface {
	GetUserModelByUsername(ctx context.Context, req *user.GetUserModelByUsernameRequest) (*user.GetUserModelByUsernameResponse, error)
	CreateNewUser(ctx context.Context, req *user.CreateNewUserRequest) (*user.CreateNewUserResponse, error)
}