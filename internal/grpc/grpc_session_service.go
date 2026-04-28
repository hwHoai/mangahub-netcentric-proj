package grpc

import (
	"context"
	"mangahub/proto/session"
)

type GRPCSessionService interface {
	SaveSession(ctx context.Context, req *session.SaveSessionRequest) (*session.SaveSessionResponse, error)
	UpdateSession(ctx context.Context, req *session.UpdateSessionRequest) (*session.UpdateSessionResponse, error)
}