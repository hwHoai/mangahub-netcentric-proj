package grpc_services_impl

import (
	"context"
	"fmt"
	"mangahub/pkg/models"
	"mangahub/pkg/repository/impl"
	"mangahub/proto/session"

	"gorm.io/gorm"
)

type GRPCSessionService struct {
	session.UnimplementedGRPCSessionServiceServer
	DBConn *gorm.DB
}

var _ session.GRPCSessionServiceServer = (*GRPCSessionService)(nil)

func (s *GRPCSessionService) SaveSession(ctx context.Context, req *session.SaveSessionRequest) (*session.SaveSessionResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request is required")
	}

	if req.UserId == "" || req.AccessToken == "" || req.RefreshToken == "" {
		return nil, fmt.Errorf("user_id, access_token, and refresh_token are required")
	}

	// Create session model
	sessionModel := models.NewSessionModel(req.UserId, req.AccessToken, req.RefreshToken)

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

func (s *GRPCSessionService) UpdateSession(ctx context.Context, req *session.UpdateSessionRequest) (*session.UpdateSessionResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request is required")
	}

	if req.UserId == "" || req.AccessToken == "" || req.RefreshToken == "" {
		return nil, fmt.Errorf("user_id, access_token, and refresh_token are required")
	}

	sessionRepo := impl.NewSessionRepository(s.DBConn)
	updatedSession, err := sessionRepo.UpdateSessionByUserID(req.UserId, req.AccessToken, req.RefreshToken)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			newSession := models.NewSessionModel(req.UserId, req.AccessToken, req.RefreshToken)
			savedSession, saveErr := sessionRepo.SaveSession(&newSession)
			if saveErr != nil {
				return nil, fmt.Errorf("failed to create session: %w", saveErr)
			}

			return &session.UpdateSessionResponse{
				SessionId:    savedSession.ID,
				UserId:       savedSession.UserID,
				AccessToken:  savedSession.AccessToken,
				RefreshToken: savedSession.RefreshToken,
				UpdatedAt:    savedSession.UpdatedAt.Format("2006-01-02 15:04:05"),
			}, nil
		}

		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	return &session.UpdateSessionResponse{
		SessionId:    updatedSession.ID,
		UserId:       updatedSession.UserID,
		AccessToken:  updatedSession.AccessToken,
		RefreshToken: updatedSession.RefreshToken,
		UpdatedAt:    updatedSession.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
