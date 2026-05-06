package controllers

import (
	"mangahub/internal/manga"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MangaController struct {
	mangaService   manga.MangaService
	chapterService manga.ChapterService
}

func NewMangaController(mangaService manga.MangaService, chapterService manga.ChapterService) *MangaController {
	return &MangaController{
		mangaService:   mangaService,
		chapterService: chapterService,
	}
}

func (mc *MangaController) ListMangas(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	mangas, err := mc.mangaService.ListMangas(int32(limit), int32(offset))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Mangas retrieved",
		"data":    mangas,
	})
}

func (mc *MangaController) GetMangaDetail(c *gin.Context) {
	mangaDetail, err := mc.mangaService.GetMangaDetail(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Manga detail retrieved",
		"data":    mangaDetail,
	})
}

func (mc *MangaController) GetMangaChapters(c *gin.Context) {
	chapters, err := mc.mangaService.GetMangaChapters(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Chapters retrieved",
		"data":    chapters,
	})
}

func (mc *MangaController) ReadChapter(c *gin.Context) {
	chapter, err := mc.chapterService.ReadChapter(c.Param("chapter_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Chapter retrieved",
		"data":    chapter,
	})
}

func (mc *MangaController) CreateNewChapter(c *gin.Context) {
	mangaID := c.Param("id")

	var reqBody struct {
		MangaDexChapterID string `json:"mangadex_chapter_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chapterID, err := mc.chapterService.CreateNewChapter(c.Request.Context(), mangaID, reqBody.MangaDexChapterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Chapter created and synced successfully",
		"data": gin.H{
			"id": chapterID,
		},
	})
}