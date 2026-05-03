package grpc

import (
	"mangahub/proto/user"
)

type GRPCUserService interface {
	user.GRPCUserServiceServer
}