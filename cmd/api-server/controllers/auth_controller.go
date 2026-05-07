package controllers

import (
	"crypto/rsa"
	"net/http"

	auth_service_impl "mangahub/internal/auth/impl"
	jwt_impl "mangahub/pkg/utils/jwt/impl"
	user_internal "mangahub/internal/user"
	"mangahub/pkg/dto"
	"mangahub/proto/session"
	"mangahub/proto/user"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	grpcUserClient    user.GRPCUserServiceClient
	grpcSessionClient session.GRPCSessionServiceClient
	userService       user_internal.UserService
	privateKey        *rsa.PrivateKey
	publicKey         *rsa.PublicKey
}

func NewAuthController(
	grpcUserClient user.GRPCUserServiceClient,
	grpcSessionClient session.GRPCSessionServiceClient,
	UserService user_internal.UserService,
	privateKey *rsa.PrivateKey,
	publicKey *rsa.PublicKey,
) *AuthController {
	return &AuthController{
		grpcUserClient:    grpcUserClient,
		grpcSessionClient: grpcSessionClient,
		userService:         UserService,
		privateKey:        privateKey,
		publicKey:         publicKey,
	}
}

func (ac *AuthController) LoginByUsername(c *gin.Context) {
	if ac.grpcUserClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "User service is unavailable"})
		return
	}

	if ac.grpcSessionClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Session service is unavailable"})
		return
	}

	var request dto.LoginByUsernameRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	loginServices := auth_service_impl.LoginServiceImpl{
		Context:           c,
		GRPCUserClient:    ac.grpcUserClient,
		GRPCSessionClient: ac.grpcSessionClient,
		PrivateKey:        ac.privateKey,
	}
	response, exception := loginServices.LoginByUsername(&request)

	if exception.Code != 0 {
		c.JSON(exception.Code, gin.H{"error": exception.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successfully",
		"data":    response,
	})
}

func (ac *AuthController) SignupByUsername(c *gin.Context) {
	if ac.grpcUserClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "User service is unavailable"})
		return
	}

	if ac.grpcSessionClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Session service is unavailable"})
		return
	}

	var request dto.SignupByUsernameRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	signupServices := auth_service_impl.SignupServiceImpl{
		Context:            c.Request.Context(),
		GRPCUserClient:     ac.grpcUserClient,
		GRPCSessionClient:  ac.grpcSessionClient,
		PrivateKey:         ac.privateKey,
	}
	response, exception := signupServices.SignupByUsername(&request)

	if exception.Code != 0 {
		c.JSON(exception.Code, gin.H{"error": exception.Message})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Signup successfully",
		"data":    response,
	})
}

func (ac *AuthController) RefreshToken(c *gin.Context) {
	var request dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token is required"})
		return
	}

	jwtUtil := jwt_impl.NewJWTUtil(ac.grpcSessionClient)

	response, exception := jwtUtil.RefreshToken(&request, ac.privateKey, ac.publicKey)
	if exception.Code != 0 {
		c.JSON(exception.Code, gin.H{"error": exception.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Refresh token successfully",
		"data":    response,
	})
}

func (ac *AuthController) GetMe(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	response, exception := ac.userService.GetUserDetails(userID)
	if exception.Code != 0 {
		c.JSON(exception.Code, gin.H{"error": exception.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Get me successfully",
		"data":    response,
	})
}
