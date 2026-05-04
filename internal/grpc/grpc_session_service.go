package grpc

import (
	"mangahub/proto/session"
)

type GRPCSessionService interface {
	session.GRPCSessionServiceServer
}