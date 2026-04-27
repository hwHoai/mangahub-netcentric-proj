package grpc_session_services

import (
	"context"
	"fmt"
	"mangahub/pkg/models"
	"mangahub/pkg/models/repository/impl"
	"mangahub/proto/session"

	"gorm.io/gorm"
)

type GRPCSessionService struct {
	session.UnimplementedGRPCSessionServiceServer
	DBConn *gorm.DB
}

func (s *GRPCSessionService) SaveSession(ctx context.Context, req *session.SaveSessionRequest) (*session.SaveSessionResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request is required")
	}

	if req.UserId == "" || req.AccessToken == "" || req.RefreshToken == "" {
		return nil, fmt.Errorf("user_id, access_token, and refresh_token are required")
	}

	// Create session model
	sessionModel := models.NewSessionModel(req.UserId, req.AccessToken, req.RefreshToken, req.PublicKey)

	// Save to database
	sessionRepo := impl.NewSessionRepository(s.DBConn)
	savedSession, err := sessionRepo.SaveSession(&sessionModel)
	if err != nil {
		return nil, fmt.Errorf("failed to save session: %w", err)
	}

	// Map to gRPC response
	response := &session.SaveSessionResponse{
		SessionId:    savedSession.ID,
		UserId:       savedSession.UserID,
		AccessToken:  savedSession.AccessToken,
		RefreshToken: savedSession.RefreshToken,
		CreatedAt:    savedSession.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	return response, nil
}
