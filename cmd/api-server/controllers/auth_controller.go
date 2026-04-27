package controllers

import (
	auth_service_impl "mangahub/internal/auth/impl"
	"mangahub/pkg/dto"
	"mangahub/proto/session"
	"mangahub/proto/user"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	grpcUserClient    user.GRPCUserServiceClient
	grpcSessionClient session.GRPCSessionServiceClient
}

func NewAuthController(grpcUserClient user.GRPCUserServiceClient, grpcSessionClient session.GRPCSessionServiceClient) *AuthController {
	return &AuthController{
		grpcUserClient:    grpcUserClient,
		grpcSessionClient: grpcSessionClient,
	}
}

func (ac *AuthController) LoginByUsername(c *gin.Context) {
	if ac.grpcUserClient == nil {
		c.JSON(500, gin.H{
			"error": "User service is unavailable",
		})
		return
	}

	if ac.grpcSessionClient == nil {
		c.JSON(500, gin.H{
			"error": "Session service is unavailable",
		})
		return
	}

	loginServices := auth_service_impl.LoginServiceImpl{
		Context:           c,
		GRPCUserClient:    ac.grpcUserClient,
		GRPCSessionClient: ac.grpcSessionClient,
	}
	response, exception := loginServices.LoginByUsername(&dto.LoginByUsernameRequest{
		Username: c.PostForm("username"),
		Password: c.PostForm("password"),
	})

	if exception.Code != 0 {
		c.JSON(exception.Code, gin.H{
			"error": exception.Message,
		})
		return
	}

	c.JSON(200, response)
}

func (ac *AuthController) SignupByUsername(c *gin.Context) {
	if ac.grpcUserClient == nil {
		c.JSON(500, gin.H{
			"error": "User service is unavailable",
		})
		return
	}

	if ac.grpcSessionClient == nil {
		c.JSON(500, gin.H{
			"error": "Session service is unavailable",
		})
		return
	}

	signupServices := auth_service_impl.SignupServiceImpl{
		Context:            c.Request.Context(),
		GRPCUserClient:     ac.grpcUserClient,
		GRPCSessionClient:  ac.grpcSessionClient,
	}	
	response, exception := signupServices.SignupByUsername(&dto.SignupByUsernameRequest{
		Username: c.PostForm("username"),
		Password: c.PostForm("password"),
	})

	if exception.Code != 0 {
		c.JSON(exception.Code, gin.H{
			"error": exception.Message,
		})
		return
	}

	c.JSON(200, response)
}
