package scrape

type Quote struct {
	Text   string `json:"text"`
	Author string `json:"author"`
}

type ScrapeService interface {
	ScrapeQuotes() ([]Quote, error)
}
