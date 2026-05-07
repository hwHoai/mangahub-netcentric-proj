package controllers

import (
	"mangahub/internal/manga"
	"mangahub/internal/user"
	"mangahub/internal/tcp"
	"net/http"
	"strconv"
	"log"

	"github.com/gin-gonic/gin"
)

type UserMangaController struct {
	userMangaService     user.UserMangaService
	chapterService       manga.ChapterService
	tcpChapterSyncClient tcp_services.TCPChapterSyncServices
}

func NewUserMangaController(
	userMangaService user.UserMangaService,
	chapterService manga.ChapterService,
	tcpChapterSyncClient tcp_services.TCPChapterSyncServices,
) *UserMangaController {
	return &UserMangaController{
		userMangaService:     userMangaService,
		chapterService:       chapterService,
		tcpChapterSyncClient: tcpChapterSyncClient,
	}
}

// POST /api/v1/user/mangas/:id/follow
func (uc *UserMangaController) FollowManga(c *gin.Context) {
	userID := c.GetString("user_id")
	mangaID := c.Param("id")

	follower, err := uc.userMangaService.FollowManga(userID, mangaID)
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

	err := uc.userMangaService.UnfollowManga(userID, mangaID)
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

	mangas, err := uc.userMangaService.GetFollowingMangas(userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Following mangas retrieved",
		"data":    mangas,
	})
}

// GET /api/v1/user/chapters/:chapter_id
func (uc *UserMangaController) ReadChapterWithDevicesSync(c *gin.Context) {
	userID := c.GetString("user_id")
	chapterID := c.Param("chapter_id")

	chapter, err := uc.chapterService.ReadChapter(chapterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Trigger TCP broadcast to other devices
	go func() {
		if err := uc.tcpChapterSyncClient.SyncReading(userID, chapterID); err != nil {
			log.Printf("Failed to broadcast reading progress: %v", err)
		}
	}()

	c.JSON(http.StatusOK, gin.H{
		"message": "Chapter retrieved and sync triggered",
		"data":    chapter,
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

	history, err := uc.userMangaService.GetReadingHistory(userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Reading history retrieved",
		"data":    history,
	})
}
