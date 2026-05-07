package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

const MangaDexBaseURL = "https://api.mangadex.org"

type MangaDexClientInterface interface {
	GetMangaList(limit int, offset int) (*MangaDexListResponse, error)
	GetTags() (map[string]TagAttributes, error)
	GetMangaChapters(mangaID string, limit int) (*MangaDexChapterListResponse, error)
	GetChapterDetails(chapterID string) (*MangaDexChapter, error)
	GetChapterPages(chapterID string) ([]string, error)
}

type MangaDexClient struct {
	baseURL    string
	httpClient *http.Client
}

var _ MangaDexClientInterface = (*MangaDexClient)(nil)

// MangaDex API Response Structures
type MangaDexListResponse struct {
	Data   []MangaDexManga `json:"data"`
	Limit  int             `json:"limit"`
	Offset int             `json:"offset"`
	Total  int             `json:"total"`
}

type MangaDexManga struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Attributes MangaDexMangaAttribute `json:"attributes"`
	Relationships []MangaDexRelationship `json:"relationships"`
}

type MangaDexMangaAttribute struct {
	Title                    map[string]string `json:"title"`
	AltTitles               []map[string]string `json:"altTitles"`
	Description             map[string]string `json:"description"`
	IsLocked                bool              `json:"isLocked"`
	Links                   map[string]string `json:"links"`
	OriginalLanguage        string            `json:"originalLanguage"`
	LastVolume              string            `json:"lastVolume"`
	LastChapter             string            `json:"lastChapter"`
	PublicationDemographic  string            `json:"publicationDemographic"`
	Status                  string            `json:"status"`
	Year                    *int              `json:"year"`
	ContentRating           string            `json:"contentRating"`
	Tags                    []MangaDexTag     `json:"tags"`
	State                   string            `json:"state"`
	ChapterNumbersResetOnNewVolume bool        `json:"chapterNumbersResetOnNewVolume"`
	CreatedAt               string            `json:"createdAt"`
	UpdatedAt               string            `json:"updatedAt"`
	Version                 int               `json:"version"`
}

type MangaDexTag struct {
	ID         string         `json:"id"`
	Type       string         `json:"type"`
	Attributes TagAttributes  `json:"attributes"`
}

type TagAttributes struct {
	Name        map[string]string `json:"name"`
	Description map[string]string `json:"description"`
	Group       string            `json:"group"`
	Version     int               `json:"version"`
}

type MangaDexRelationship struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Related    string                 `json:"related,omitempty"`
	Attributes map[string]interface{} `json:"attributes"`
}

func NewMangaDexClient() *MangaDexClient {
	// Create a custom dialer that uses Cloudflare DNS (1.1.1.1) to bypass system DNS issues
	dialer := &net.Dialer{
		Resolver: &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				// Use Cloudflare DNS instead of system DNS to resolve MangaDex
				return net.Dial("udp", "1.1.1.1:53")
			},
		},
		Timeout: 15 * time.Second, // Longer timeout for TLS handshake
	}

	// Create transport with custom dialer and longer TLS handshake timeout
	transport := &http.Transport{
		DialContext:           dialer.DialContext,
		TLSHandshakeTimeout:   15 * time.Second,
		ResponseHeaderTimeout: 15 * time.Second,
		ExpectContinueTimeout: 5 * time.Second,
	}

	return &MangaDexClient{
		baseURL: MangaDexBaseURL,
		httpClient: &http.Client{
			Transport: transport,
			Timeout:   45 * time.Second,
		},
	}
}

// addHeaders adds required headers to requests to bypass MangaDex anti-bot protection
func (c *MangaDexClient) addHeaders(req *http.Request) {
	// Use application name instead of Go's default user-agent
	req.Header.Set("User-Agent", "MangaHub/1.0 (+https://github.com/mangahub)")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Connection", "keep-alive")
}

// GetMangaList fetches a list of manga from MangaDex API with pagination
func (c *MangaDexClient) GetMangaList(limit int, offset int) (*MangaDexListResponse, error) {
	if limit > 100 {
		limit = 100
	}

	// Build URL with proper query parameters
	url := fmt.Sprintf("%s/manga?limit=%d&offset=%d&includes[]=author&includes[]=artist&includes[]=cover_art&order[relevance]=desc", c.baseURL, limit, offset)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	c.addHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch manga: status code %d from %s", resp.StatusCode, url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result MangaDexListResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetTags fetches all available tags (genres) from MangaDex
func (c *MangaDexClient) GetTags() (map[string]TagAttributes, error) {
	url := fmt.Sprintf("%s/manga/tag", c.baseURL)

	type TagResponse struct {
		Data []MangaDexTag `json:"data"`
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	c.addHeaders(req)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get tags: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tagResp TagResponse
	if err := json.Unmarshal(body, &tagResp); err != nil {
		return nil, err
	}

	tagMap := make(map[string]TagAttributes)
	for _, tag := range tagResp.Data {
		tagMap[tag.ID] = tag.Attributes
	}

	return tagMap, nil
}

// Chapter and Pages API Response Structures
type MangaDexChapterDetailResponse struct {
	Data MangaDexChapter `json:"data"`
}

type MangaDexChapterListResponse struct {
	Data []MangaDexChapter `json:"data"`
}

type MangaDexChapter struct {
	ID            string                   `json:"id"`
	Attributes    MangaDexChapterAttribute `json:"attributes"`
	Relationships []MangaDexRelationship   `json:"relationships"`
}

type MangaDexChapterAttribute struct {
	Chapter string `json:"chapter"`
	Title   string `json:"title"`
	Pages   int    `json:"pages"`
}

type MangaDexAtHomeResponse struct {
	BaseUrl string `json:"baseUrl"`
	Chapter struct {
		Hash string   `json:"hash"`
		Data []string `json:"data"`
	} `json:"chapter"`
}

// GetMangaChapters fetches chapters for a given manga ID
func (c *MangaDexClient) GetMangaChapters(mangaID string, limit int) (*MangaDexChapterListResponse, error) {
	url := fmt.Sprintf("%s/manga/%s/feed?translatedLanguage[]=en&order[chapter]=asc&limit=%d", c.baseURL, mangaID, limit)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	c.addHeaders(req)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get chapters: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result MangaDexChapterListResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetChapterDetails fetches a specific chapter's metadata
func (c *MangaDexClient) GetChapterDetails(chapterID string) (*MangaDexChapter, error) {
	url := fmt.Sprintf("%s/chapter/%s", c.baseURL, chapterID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "MangaHub-StudentProject/1.0 (github.com/hwHoai)")
	
	c.addHeaders(req)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get chapter details: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result MangaDexChapterDetailResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// GetChapterPages fetches page image URLs for a given chapter ID
func (c *MangaDexClient) GetChapterPages(chapterID string) ([]string, error) {
	url := fmt.Sprintf("%s/at-home/server/%s", c.baseURL, chapterID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	c.addHeaders(req)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get pages: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result MangaDexAtHomeResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	var pages []string
	for _, filename := range result.Chapter.Data {
		pageUrl := fmt.Sprintf("%s/data/%s/%s", result.BaseUrl, result.Chapter.Hash, filename)
		pages = append(pages, pageUrl)
	}

	return pages, nil
}
