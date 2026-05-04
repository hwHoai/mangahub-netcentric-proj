package grpc

import (
	"mangahub/proto/manga"
)

type GRPCMangaService interface {
	manga.GRPCMangaServiceServer
}