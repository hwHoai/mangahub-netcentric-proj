package grpc_services_impl

import (
	"context"
	"fmt"
	"mangahub/internal/grpc"
	repository_impl "mangahub/internal/repository/impl"
	"mangahub/pkg/models"
	"mangahub/pkg/utils"
	"mangahub/proto/session"

	"gorm.io/gorm"
)

type GRPCSessionService struct {
	session.UnimplementedGRPCSessionServiceServer
	db *gorm.DB
}
var _ grpc.GRPCSessionService = (*GRPCSessionService)(nil)

func NewGRPCSessionService(db *gorm.DB) *GRPCSessionService {
	return &GRPCSessionService{db: db}
}

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
	sessionRepo := repository_impl.NewSessionRepository(s.db)
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
		CreatedAt:    savedSession.CreatedAt.Format(utils.TimeLayout),
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

	sessionRepo := repository_impl.NewSessionRepository(s.db)
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
				UpdatedAt:    savedSession.UpdatedAt.Format(utils.TimeLayout),
			}, nil
		}

		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	return &session.UpdateSessionResponse{
		SessionId:    updatedSession.ID,
		UserId:       updatedSession.UserID,
		AccessToken:  updatedSession.AccessToken,
		RefreshToken: updatedSession.RefreshToken,
		UpdatedAt:    updatedSession.UpdatedAt.Format(utils.TimeLayout),
	}, nil
}
func (s *GRPCSessionService) GetSessionByRefreshToken(ctx context.Context, req *session.GetSessionByRefreshTokenRequest) (*session.GetSessionByRefreshTokenResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request is required")
	}

	if req.RefreshToken == "" {
		return nil, fmt.Errorf("refresh_token is required")
	}

	sessionRepo := repository_impl.NewSessionRepository(s.db)
	sessionModel, err := sessionRepo.GetSessionByRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get session by refresh token: %w", err)
	}

	return &session.GetSessionByRefreshTokenResponse{
		SessionId:    sessionModel.ID,
		UserId:       sessionModel.UserID,
		AccessToken:  sessionModel.AccessToken,
		RefreshToken: sessionModel.RefreshToken,
		CreatedAt:    sessionModel.CreatedAt.Format(utils.TimeLayout),
		UpdatedAt:    sessionModel.UpdatedAt.Format(utils.TimeLayout),
	}, nil
}
