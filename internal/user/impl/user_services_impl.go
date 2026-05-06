package user_services_impl

import (
	"context"
	"mangahub/internal/user"
	"mangahub/pkg/dto"
	user_proto "mangahub/proto/user"
)

type UserServiceImpl struct {
	GRPCUserClient user_proto.GRPCUserServiceClient
}

func NewUserService(grpcUserClient user_proto.GRPCUserServiceClient) user.UserService {
	return &UserServiceImpl{GRPCUserClient: grpcUserClient}
}

var _ user.UserService = (*UserServiceImpl)(nil)

func (s *UserServiceImpl) GetUserDetails(userID string) (*dto.UserProfileResponse, dto.ExceptionDTO) {
	if userID == "" {
		return nil, dto.ExceptionDTO{Code: 401, Message: "Unauthorized"}
	}

	grpcResponse, err := s.GRPCUserClient.GetUserByID(context.Background(), &user_proto.GetUserByIDRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, dto.ExceptionDTO{
			Code:    404,
			Message: "User not found",
		}
	}

	return &dto.UserProfileResponse{
		UserID:    grpcResponse.UserId,
		Username:  grpcResponse.Username,
		CreatedAt: grpcResponse.CreatedAt,
		UpdatedAt: grpcResponse.UpdatedAt,
	}, dto.ExceptionDTO{}
}
