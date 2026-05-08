package impl

import (
	"mangahub/internal/scrape"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type scrapeServiceImpl struct {
}

func NewScrapeService() scrape.ScrapeService {
	return &scrapeServiceImpl{}
}

func (s *scrapeServiceImpl) ScrapeQuotes() ([]scrape.Quote, error) {
	resp, err := http.Get("http://quotes.toscrape.com")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	var quotes []scrape.Quote

	// Walk the DOM tree to find quote elements
	var walkNode func(*html.Node)
	walkNode = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" && hasClass(n, "quote") {
			q := scrape.Quote{}
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

	return quotes, nil
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
