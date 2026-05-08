package grpc_services_impl

import (
	"context"
	"mangahub/internal/grpc"
	"mangahub/internal/repository"
	repository_impl "mangahub/internal/repository/impl"
	"mangahub/pkg/models"
	"mangahub/pkg/utils"
	"mangahub/proto/message"

	"gorm.io/gorm"
)

type GRPCMessageService struct {
	message.UnimplementedGRPCMessageServiceServer
	repo repository.MessageRepository
}

var _ grpc.GRPCMessageService = (*GRPCMessageService)(nil)

func NewGRPCMessageService(db *gorm.DB) *GRPCMessageService {
	return &GRPCMessageService{
		repo: repository_impl.NewMessageRepository(db),
	}
}

func (s *GRPCMessageService) SaveMessage(ctx context.Context, req *message.SaveMessageRequest) (*message.SaveMessageResponse, error) {
	msg := models.NewMessageModel(req.SenderId, req.RoomId, req.Content)
	if err := s.repo.SaveMessage(msg); err != nil {
		return nil, err
	}

	return &message.SaveMessageResponse{
		Message: &message.MessageResponse{
			Id:        msg.ID,
			SenderId:  msg.SenderID,
			RoomId:    msg.RoomID,
			Content:   msg.Content,
			CreatedAt: msg.CreatedAt.Format(utils.TimeLayout),
		},
	}, nil
}

func (s *GRPCMessageService) GetChatHistory(ctx context.Context, req *message.GetChatHistoryRequest) (*message.GetChatHistoryResponse, error) {
	limit := int(req.Limit)
	if limit <= 0 {
		limit = 50
	}
	offset := int(req.Offset)

	messages, err := s.repo.GetChatHistory(req.RoomId, limit, offset)
	if err != nil {
		return nil, err
	}

	var protoMessages []*message.MessageResponse
	for _, m := range messages {
		protoMessages = append(protoMessages, &message.MessageResponse{
			Id:        m.ID,
			SenderId:  m.SenderID,
			RoomId:    m.RoomID,
			Content:   m.Content,
			CreatedAt: m.CreatedAt.Format(utils.TimeLayout),
		})
	}

	return &message.GetChatHistoryResponse{
		Messages: protoMessages,
	}, nil
}
