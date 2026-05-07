package controllers

import (
	"mangahub/internal/manga"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/html"
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

// GET /api/v1/quotes
// ScrapeQuotes demonstrates educational web scraping from a practice site.
// Uses golang.org/x/net/html for proper DOM parsing instead of fragile string matching.
func (mc *MangaController) ScrapeQuotes(c *gin.Context) {
	resp, err := http.Get("http://quotes.toscrape.com")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch quotes"})
		return
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse HTML"})
		return
	}

	type Quote struct {
		Text   string `json:"text"`
		Author string `json:"author"`
	}
	var quotes []Quote

	// Walk the DOM tree to find quote elements
	var walkNode func(*html.Node)
	walkNode = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" && hasClass(n, "quote") {
			q := Quote{}
			// Find text and author within this quote div
			var extractFromQuote func(*html.Node)
			extractFromQuote = func(child *html.Node) {
				if child.Type == html.ElementNode {
					if child.Data == "span" && hasClass(child, "text") {
						q.Text = getTextContent(child)
					}
					if child.Data == "small" && hasClass(child, "author") {
						q.Author = getTextContent(child)
					}
				}
				for c := child.FirstChild; c != nil; c = c.NextSibling {
					extractFromQuote(c)
				}
			}
			extractFromQuote(n)
			if q.Text != "" && q.Author != "" {
				quotes = append(quotes, q)
			}
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			walkNode(child)
		}
	}
	walkNode(doc)

	c.JSON(http.StatusOK, gin.H{
		"message": "Quotes scraped successfully (Educational Practice)",
		"source":  "http://quotes.toscrape.com",
		"data":    quotes,
	})
}

// hasClass checks if an HTML node has a specific CSS class.
func hasClass(n *html.Node, className string) bool {
	for _, attr := range n.Attr {
		if attr.Key == "class" {
			for _, cls := range strings.Split(attr.Val, " ") {
				if cls == className {
					return true
				}
			}
		}
	}
	return false
}

// getTextContent extracts all text content from an HTML node and its children.
func getTextContent(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	var result string
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		result += getTextContent(child)
	}
	return result
}