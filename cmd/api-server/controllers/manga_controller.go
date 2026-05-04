package controllers

import (
	manga_services_impl "mangahub/internal/manga/impl"
	"net/http"
	"strconv"

	"mangahub/proto/chapter"
	"mangahub/proto/manga"

	"github.com/gin-gonic/gin"
)

type MangaController struct {
	grpcMangaClient   manga.GRPCMangaServiceClient
	grpcChapterClient chapter.GRPCChapterServiceClient
}

func NewMangaController(grpcMangaClient manga.GRPCMangaServiceClient, grpcChapterClient chapter.GRPCChapterServiceClient) *MangaController {
	return &MangaController{
		grpcMangaClient:   grpcMangaClient,
		grpcChapterClient: grpcChapterClient,
	}
}

func (mc *MangaController) ListMangas(c *gin.Context) {
	mangaService := manga_services_impl.NewMangaService(mc.grpcMangaClient)
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	mangas, err := mangaService.ListMangas(int32(limit), int32(offset))
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
	mangaService := manga_services_impl.NewMangaService(mc.grpcMangaClient)	
	mangaDetail, err := mangaService.GetMangaDetail(c.Param("id"))
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
	mangaService := manga_services_impl.NewMangaService(mc.grpcMangaClient)	
	chapters, err := mangaService.GetMangaChapters(c.Param("id"))
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
	chapterService := manga_services_impl.NewChapterService(mc.grpcChapterClient)
	chapter, err := chapterService.ReadChapter(c.Param("chapter_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Chapter retrieved",
		"data":    chapter,
	})
}