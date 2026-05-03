package grpc

import (
	"mangahub/proto/user_manga"
)

type GRPCUserMangaService interface {
	user_manga.GRPCUserMangaServiceServer
}
