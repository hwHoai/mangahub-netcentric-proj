package grpc

import (
	"mangahub/proto/chapter"
)

type GRPCChapterService interface {
	chapter.GRPCChapterServiceServer
}
