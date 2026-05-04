package controllers

import (
	"crypto/rsa"
	auth_service_impl "mangahub/internal/auth/impl"
	user_service_impl "mangahub/internal/user/impl"
	"mangahub/pkg/dto"
	"mangahub/proto/session"
	"mangahub/proto/user"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	grpcUserClient    user.GRPCUserServiceClient
	grpcSessionClient session.GRPCSessionServiceClient
	privateKey        *rsa.PrivateKey
	publicKey         *rsa.PublicKey
}

func NewAuthController(grpcUserClient user.GRPCUserServiceClient, grpcSessionClient session.GRPCSessionServiceClient, privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) *AuthController {
	return &AuthController{
		grpcUserClient:    grpcUserClient,
		grpcSessionClient: grpcSessionClient,
		privateKey:        privateKey,
		publicKey:         publicKey,
	}
}

func (ac *AuthController) LoginByUsername(c *gin.Context) {
	if ac.grpcUserClient == nil {
		c.JSON(500, gin.H{"error": "User service is unavailable"})
		return
	}

	if ac.grpcSessionClient == nil {
		c.JSON(500, gin.H{"error": "Session service is unavailable"})
		return
	}

	loginServices := auth_service_impl.LoginServiceImpl{
		Context:           c,
		GRPCUserClient:    ac.grpcUserClient,
		GRPCSessionClient: ac.grpcSessionClient,
		PrivateKey:        ac.privateKey,
	}
	response, exception := loginServices.LoginByUsername(&dto.LoginByUsernameRequest{
		Username: c.PostForm("username"),
		Password: c.PostForm("password"),
	})

	if exception.Code != 0 {
		c.JSON(exception.Code, gin.H{"error": exception.Message})
		return
	}

	c.JSON(200, gin.H{
		"message": "Login successfully",
		"data":    response,	
	})
}

func (ac *AuthController) SignupByUsername(c *gin.Context) {
	if ac.grpcUserClient == nil {
		c.JSON(500, gin.H{"error": "User service is unavailable"})
		return
	}

	if ac.grpcSessionClient == nil {
		c.JSON(500, gin.H{"error": "Session service is unavailable"})
		return
	}

	signupServices := auth_service_impl.SignupServiceImpl{
		Context:            c.Request.Context(),
		GRPCUserClient:     ac.grpcUserClient,
		GRPCSessionClient:  ac.grpcSessionClient,
		PrivateKey:         ac.privateKey,
	}	
	response, exception := signupServices.SignupByUsername(&dto.SignupByUsernameRequest{
		Username: c.PostForm("username"),
		Password: c.PostForm("password"),
	})

	if exception.Code != 0 {
		c.JSON(exception.Code, gin.H{"error": exception.Message})
		return
	}

	c.JSON(200, gin.H{
		"message": "Signup successfully",
		"data":    response,	
	})
}

func (ac *AuthController) RefreshToken(c *gin.Context) {
	var request dto.RefreshTokenRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(400, gin.H{"error": "Refresh token is required"})
		return
	}

	jwtService := auth_service_impl.NewJWTService(ac.grpcSessionClient)

	response, exception := jwtService.RefreshToken(&request, ac.privateKey, ac.publicKey)
	if exception.Code != 0 {
		c.JSON(exception.Code, gin.H{"error": exception.Message})
		return
	}

	c.JSON(200, gin.H{
		"message": "Refresh token successfully",
		"data":    response,	
	})
}

func (ac *AuthController) GetMe(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	meService := user_service_impl.MeServiceImpl{
		GRPCUserClient: ac.grpcUserClient,
	}

	response, exception := meService.GetMe(userID)
	if exception.Code != 0 {
		c.JSON(exception.Code, gin.H{"error": exception.Message})
		return
	}

	c.JSON(200, gin.H{
		"message": "Get me successfully",
		"data":    response,
	})
}
