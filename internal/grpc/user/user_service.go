package grpc_user_services

import (
	"context"
	"crypto/sha256"
	"fmt"
	"mangahub/pkg/models"
	"mangahub/pkg/repository/impl"
	"mangahub/proto/user"

	"gorm.io/gorm"
)

type GRPCUserService struct {
	user.UnimplementedGRPCUserServiceServer
	DBConn *gorm.DB
}

func (s *GRPCUserService) GetUserModelByUsername(ctx context.Context, req *user.GetUserModelByUsernameRequest) (*user.GetUserModelByUsernameResponse, error) {
	// 1. Fetch user data from database using repository
	userRepo := impl.NewUserRepository(s.DBConn)
	userModel, err := userRepo.GetUserByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	// 2. Map userModel to gRPC response
	response := &user.GetUserModelByUsernameResponse{
		UserId:          userModel.ID,
		Username:        userModel.Username,
		CreatedAt:       userModel.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       userModel.UpdatedAt.Format("2006-01-02 15:04:05"),
		HashedPassword:  userModel.HashedPassword,
	}
	return response, nil
}

func (s *GRPCUserService) CreateNewUser(ctx context.Context, req *user.CreateNewUserRequest) (*user.CreateNewUserResponse, error) {
	if req.Username == "" || req.Password == "" {
		return nil, fmt.Errorf("username and password are required")
	}

	// Hash password using SHA256 (consider bcrypt in production)
	hash := sha256.Sum256([]byte(req.Password))
	hashedPassword := fmt.Sprintf("%x", hash)

	// Create user model
	userModel := models.NewUserModel(req.Username, hashedPassword)

	// Save to database
	userRepo := impl.NewUserRepository(s.DBConn)
	savedUser, err := userRepo.CreateUser(userModel)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Map to gRPC response
	response := &user.CreateNewUserResponse{
		UserId: savedUser.ID,
		Username: savedUser.Username,
		CreatedAt: savedUser.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: savedUser.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	return response, nil
}
