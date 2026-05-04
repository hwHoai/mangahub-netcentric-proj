package seeder

import (
	"encoding/json"
	"fmt"
	"log"
	"mangahub/pkg/clients"
	"mangahub/pkg/models"
	"mangahub/pkg/models/enums"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type MangaSeeder struct {
	db             *gorm.DB
	mangaDexClient *clients.MangaDexClient
}

func NewMangaSeeder(db *gorm.DB) *MangaSeeder {
	return &MangaSeeder{
		db:             db,
		mangaDexClient: clients.NewMangaDexClient(),
	}
}

// SeedMangaData seeds manga data from MangaDex API
func (ms *MangaSeeder) SeedMangaData(limit int, numberOfBatches int) error {
	log.Printf("Starting to seed manga data from MangaDex. Limit: %d, Batches: %d", limit, numberOfBatches)

	// 1. Fetch and seed genres/tags first
	if err := ms.seedGenres(); err != nil {
		log.Printf("Error seeding genres: %v", err)
		return err
	}

	// 2. Fetch and seed manga data in batches
	for batch := 0; batch < numberOfBatches; batch++ {
		offset := batch * limit
		log.Printf("Fetching batch %d (offset: %d, limit: %d)", batch+1, offset, limit)

		mangaList, err := ms.mangaDexClient.GetMangaList(limit, offset)
		if err != nil {
			log.Printf("Error fetching manga list (batch %d): %v", batch+1, err)
			continue
		}

		if len(mangaList.Data) == 0 {
			log.Println("No more manga data to fetch")
			break
		}

		for _, mdManga := range mangaList.Data {
			if err := ms.seedSingleManga(&mdManga); err != nil {
				log.Printf("Error seeding manga %s: %v", mdManga.ID, err)
				continue
			}
		}

		log.Printf("Completed batch %d", batch+1)
	}

	log.Println("Manga seeding completed!")
	return nil
}

// seedGenres fetches all genres from MangaDex and saves them to database
// Uses retry logic to handle connection resets from MangaDex WAF
func (ms *MangaSeeder) seedGenres() error {
	log.Println("Starting to seed genres...")

	var tags map[string]clients.TagAttributes
	var err error

	// Retry up to 3 times with exponential backoff
	for attempt := 1; attempt <= 3; attempt++ {
		tags, err = ms.mangaDexClient.GetTags()
		if err == nil {
			break
		}

		if attempt < 3 {
			waitTime := time.Duration(attempt*2) * time.Second
			log.Printf("Attempt %d failed: %v. Retrying in %v...", attempt, err, waitTime)
			time.Sleep(waitTime)
		} else {
			return fmt.Errorf("failed to fetch tags after 3 attempts: %w", err)
		}
	}

	for tagID, tagAttr := range tags {
		// Get English name for the genre
		name := ""
		if enName, ok := tagAttr.Name["en"]; ok {
			name = enName
		} else if len(tagAttr.Name) > 0 {
			for _, n := range tagAttr.Name {
				name = n
				break
			}
		}

		if name == "" {
			continue
		}

		// Get English description
		description := ""
		if enDesc, ok := tagAttr.Description["en"]; ok {
			description = enDesc
		} else if len(tagAttr.Description) > 0 {
			for _, d := range tagAttr.Description {
				description = d
				break
			}
		}

		// Check if genre already exists
		var existingGenre models.GenresModel
		if err := ms.db.Where("name = ?", name).First(&existingGenre).Error; err == gorm.ErrRecordNotFound {
			genre := &models.GenresModel{
				ID:          tagID,
				Name:        name,
				Description: description,
			}

			if err := ms.db.Create(genre).Error; err != nil {
				log.Printf("Error saving genre %s: %v", name, err)
				continue
			}

			log.Printf("Seeded genre: %s (ID: %s)", name, tagID)
		}
	}

	log.Println("Genre seeding completed!")
	return nil
}

// seedSingleManga maps a MangaDex manga to the local schema and saves it
func (ms *MangaSeeder) seedSingleManga(mdManga *clients.MangaDexManga) error {
	// Check if manga already exists
	var existingManga models.MangaModel
	if err := ms.db.Where("id = ?", mdManga.ID).First(&existingManga).Error; err == nil {
		return nil
	} else if err != gorm.ErrRecordNotFound {
		return err
	}

	// Extract basic information
	title := ms.getTitleInEnglish(mdManga.Attributes.Title)
	if title == "" {
		return fmt.Errorf("no valid title found for manga %s", mdManga.ID)
	}

	description := ms.getDescriptionInEnglish(mdManga.Attributes.Description)
	author := ms.getAuthorName(mdManga.Relationships)
	coverURL := ms.getCoverURL(mdManga.Relationships)
	totalChapters := ms.parseChapterCount(mdManga.Attributes.LastChapter)
	status := ms.mapStatus(mdManga.Attributes.Status)

	// Create manga model
	manga := &models.MangaModel{
		ID:            mdManga.ID,
		Title:         title,
		Author:        author,
		TotalChapters: totalChapters,
		Description:   description,
		CoverURL:      coverURL,
		Status:        status,
	}

	// Save manga to database
	if err := ms.db.Create(manga).Error; err != nil {
		return fmt.Errorf("failed to save manga: %w", err)
	}

	// Add genres/tags to the manga
	if err := ms.attachGenresToManga(manga, mdManga.Attributes.Tags); err != nil {
		log.Printf("Error attaching genres to manga %s: %v", manga.ID, err)
	}

	log.Printf("Seeded manga: %s (ID: %s)", title, manga.ID)

	// Fetch and seed chapters (Limit to 10 to save time and API rate limits)
	if err := ms.seedChaptersForManga(manga.ID, 10); err != nil {
		log.Printf("Error seeding chapters for manga %s: %v", manga.ID, err)
	}

	return nil
}

// attachGenresToManga associates genres with a manga
func (ms *MangaSeeder) attachGenresToManga(manga *models.MangaModel, tags []clients.MangaDexTag) error {
	for _, tag := range tags {
		var genre models.GenresModel
		if err := ms.db.Where("id = ?", tag.ID).First(&genre).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				continue
			}
			return err
		}

		if err := ms.db.Model(manga).Association("Genres").Append(&genre); err != nil {
			log.Printf("Error associating genre %s: %v", genre.Name, err)
			continue
		}
	}

	return nil
}

// getTitleInEnglish extracts English title from MangaDex's multilingual title map
func (ms *MangaSeeder) getTitleInEnglish(titles map[string]string) string {
	if en, ok := titles["en"]; ok && en != "" {
		return en
	}
	if ja, ok := titles["ja"]; ok && ja != "" {
		return ja
	}
	if jaRo, ok := titles["ja-ro"]; ok && jaRo != "" {
		return jaRo
	}
	for _, title := range titles {
		if title != "" {
			return title
		}
	}
	return ""
}

// getDescriptionInEnglish extracts English description
func (ms *MangaSeeder) getDescriptionInEnglish(descriptions map[string]string) string {
	if en, ok := descriptions["en"]; ok && en != "" {
		return en
	}
	for _, desc := range descriptions {
		if desc != "" {
			return desc
		}
	}
	return ""
}

// getAuthorName extracts author name from relationships
func (ms *MangaSeeder) getAuthorName(relationships []clients.MangaDexRelationship) string {
	for _, rel := range relationships {
		if rel.Type == "author" {
			if nameAttr, ok := rel.Attributes["name"]; ok {
				if name, ok := nameAttr.(string); ok {
					return name
				}
			}
		}
	}
	return "Unknown"
}

// getCoverURL constructs the cover image URL from relationships
func (ms *MangaSeeder) getCoverURL(relationships []clients.MangaDexRelationship) string {
	for _, rel := range relationships {
		if rel.Type == "cover_art" {
			if fileName, ok := rel.Attributes["fileName"]; ok {
				if fileNameStr, ok := fileName.(string); ok {
					return fmt.Sprintf("https://uploads.mangadex.org/covers/%s/%s.256.jpg", rel.ID, fileNameStr)
				}
			}
		}
	}
	return ""
}

// parseChapterCount converts the last chapter string to an integer
func (ms *MangaSeeder) parseChapterCount(lastChapter string) int {
	if lastChapter == "" || lastChapter == "null" {
		return 0
	}

	f, err := strconv.ParseFloat(lastChapter, 64)
	if err != nil {
		return 0
	}

	return int(f)
}

// mapStatus maps MangaDex status to project's MangaStatus enum
func (ms *MangaSeeder) mapStatus(mdStatus string) enums.MangaStatus {
	switch strings.ToLower(mdStatus) {
	case "ongoing":
		return enums.MangaStatusInProgress
	case "completed":
		return enums.MangaStatusCompleted
	case "hiatus":
		return enums.MangaStatusInProgress
	case "cancelled":
		return enums.MangaStatusInProgress
	default:
		return enums.MangaStatusInProgress
	}
}

// seedChaptersForManga fetches chapters from MangaDex and saves them
func (ms *MangaSeeder) seedChaptersForManga(mangaID string, limit int) error {
	log.Printf("Fetching chapters for manga %s", mangaID)

	chapterList, err := ms.mangaDexClient.GetMangaChapters(mangaID, limit)
	if err != nil {
		return err
	}

	if len(chapterList.Data) == 0 {
		return nil
	}

	for _, mdChapter := range chapterList.Data {
		// Check if chapter already exists
		var existingChapter models.ChapterModel
		if err := ms.db.Where("id = ?", mdChapter.ID).First(&existingChapter).Error; err == nil {
			continue // Already exists
		} else if err != gorm.ErrRecordNotFound {
			log.Printf("Database error checking chapter %s: %v", mdChapter.ID, err)
			continue
		}

		// Parse chapter number
		chapNumStr := mdChapter.Attributes.Chapter
		var chapterNumber float64
		if chapNumStr != "" && chapNumStr != "null" {
			parsed, err := strconv.ParseFloat(chapNumStr, 64)
			if err == nil {
				chapterNumber = parsed
			}
		}

		// Fetch pages - wait a bit to avoid hitting MangaDex rate limits
		time.Sleep(500 * time.Millisecond)
		pages, err := ms.mangaDexClient.GetChapterPages(mdChapter.ID)
		if err != nil {
			log.Printf("Error fetching pages for chapter %s: %v", mdChapter.ID, err)
			continue
		}

		// Serialize pages to JSON string
		pagesDataBytes, err := json.Marshal(pages)
		if err != nil {
			log.Printf("Error marshaling pages for chapter %s: %v", mdChapter.ID, err)
			continue
		}

		// Create and save ChapterModel
		chapter := models.NewChapterModel(mangaID, chapterNumber, mdChapter.Attributes.Title, string(pagesDataBytes))
		chapter.ID = mdChapter.ID // Use MangaDex ID for consistency

		if err := ms.db.Create(chapter).Error; err != nil {
			log.Printf("Failed to save chapter %s: %v", mdChapter.ID, err)
			continue
		}
		log.Printf("Seeded chapter %s (ID: %s) with %d pages", chapNumStr, mdChapter.ID, len(pages))
	}

	return nil
}
