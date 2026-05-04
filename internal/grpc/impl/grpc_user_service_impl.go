package grpc_services_impl

import (
	"context"
	"crypto/sha256"
	"fmt"
	"mangahub/internal/grpc"
	"mangahub/pkg/models"
	"mangahub/pkg/utils"
	repository_impl "mangahub/pkg/repository/impl"
	"mangahub/proto/user"

	"gorm.io/gorm"
)

type GRPCUserService struct {
	user.UnimplementedGRPCUserServiceServer
	db *gorm.DB
}

var _ grpc.GRPCUserService = (*GRPCUserService)(nil)

func NewGRPCUserService(db *gorm.DB) *GRPCUserService {
	return &GRPCUserService{db: db}
}

func (s *GRPCUserService) GetUserModelByUsername(ctx context.Context, req *user.GetUserModelByUsernameRequest) (*user.GetUserModelByUsernameResponse, error) {
	// 1. Fetch user data from database using repository
	userRepo := repository_impl.NewUserRepository(s.db)
	userModel, err := userRepo.GetUserByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	// 2. Map userModel to gRPC response
	response := &user.GetUserModelByUsernameResponse{
		UserId:          userModel.ID,
		Username:        userModel.Username,
		CreatedAt:       userModel.CreatedAt.Format(utils.TimeLayout),
		UpdatedAt:       userModel.UpdatedAt.Format(utils.TimeLayout),
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
	userRepo := repository_impl.NewUserRepository(s.db)
	savedUser, err := userRepo.CreateUser(userModel)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Map to gRPC response
	response := &user.CreateNewUserResponse{
		UserId:    savedUser.ID,
		Username:  savedUser.Username,
		CreatedAt: savedUser.CreatedAt.Format(utils.TimeLayout),
		UpdatedAt: savedUser.UpdatedAt.Format(utils.TimeLayout),
	}
	return response, nil
}

func (s *GRPCUserService) GetUserByID(ctx context.Context, req *user.GetUserByIDRequest) (*user.GetUserByIDResponse, error) {
	userRepo := repository_impl.NewUserRepository(s.db)
	userModel, err := userRepo.GetUserByID(req.UserId)
	if err != nil {
		return nil, err
	}

	return &user.GetUserByIDResponse{
		UserId:    userModel.ID,
		Username:  userModel.Username,
		CreatedAt: userModel.CreatedAt.Format(utils.TimeLayout),
		UpdatedAt: userModel.UpdatedAt.Format(utils.TimeLayout),
	}, nil
}
