package controllers

import (
	"mangahub/internal/scrape"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ScrapeController struct {
	scrapeService scrape.ScrapeService
}

func NewScrapeController(scrapeService scrape.ScrapeService) *ScrapeController {
	return &ScrapeController{
		scrapeService: scrapeService,
	}
}

// GET /api/v1/quotes
// ScrapeQuotes demonstrates educational web scraping from a practice site.
func (sc *ScrapeController) ScrapeQuotes(c *gin.Context) {
	quotes, err := sc.scrapeService.ScrapeQuotes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Quotes scraped successfully (Educational Practice)",
		"source":  "http://quotes.toscrape.com",
		"data":    quotes,
	})
}
