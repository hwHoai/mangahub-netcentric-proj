package grpc

import (
	"mangahub/proto/message"
)

type GRPCMessageService interface {
	message.GRPCMessageServiceServer
}
