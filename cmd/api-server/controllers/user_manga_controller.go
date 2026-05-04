package controllers

import (
	user_services_impl "mangahub/internal/user/impl"
	"mangahub/proto/user_manga"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserMangaController struct {
	grpcUserMangaClient user_manga.GRPCUserMangaServiceClient
}

func NewUserMangaController(grpcUserMangaClient user_manga.GRPCUserMangaServiceClient) *UserMangaController {
	return &UserMangaController{
		grpcUserMangaClient: grpcUserMangaClient,
	}
}

// POST /api/v1/user/mangas/:id/follow
func (uc *UserMangaController) FollowManga(c *gin.Context) {
	userID := c.GetString("user_id")
	mangaID := c.Param("id")

	service := user_services_impl.NewUserMangaService(uc.grpcUserMangaClient)

	follower, err := service.FollowManga(userID, mangaID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Successfully followed manga",
		"data":    follower,
	})
}

// DELETE /api/v1/user/mangas/:id/follow
func (uc *UserMangaController) UnfollowManga(c *gin.Context) {
	userID := c.GetString("user_id")
	mangaID := c.Param("id")

	service := user_services_impl.NewUserMangaService(uc.grpcUserMangaClient)

	err := service.UnfollowManga(userID, mangaID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully unfollowed manga",
	})
}

// GET /api/v1/user/mangas/following
func (uc *UserMangaController) GetFollowingMangas(c *gin.Context) {
	userID := c.GetString("user_id")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	service := user_services_impl.NewUserMangaService(uc.grpcUserMangaClient)

	mangas, err := service.GetFollowingMangas(userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Following mangas retrieved",
		"data":    mangas,
	})
}

// POST /api/v1/user/chapters/:chapter_id/read
func (uc *UserMangaController) StoreReadingProgress(c *gin.Context) {
	userID := c.GetString("user_id")
	chapterID := c.Param("chapter_id")

	service := user_services_impl.NewUserMangaService(uc.grpcUserMangaClient)

	progress, err := service.StoreReadingProgress(userID, chapterID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Reading progress stored",
		"data":    progress,
	})
}
// GET /api/v1/user/history
func (uc *UserMangaController) GetReadingHistory(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: user_id is empty"})
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	service := user_services_impl.NewUserMangaService(uc.grpcUserMangaClient)

	history, err := service.GetReadingHistory(userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Reading history retrieved",
		"data":    history,
	})
}
