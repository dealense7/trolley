package strategies

import (
	"fmt"
	"regexp"
	"storePrices/internal/domain/parser"
	"storePrices/internal/domain/retailer"
	"storePrices/internal/platform/logger"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

type GlovoStrategy struct {
	log *zap.Logger
}

// Ensure interface compliance
var _ parser.Strategy = (*GlovoStrategy)(nil)

func NewGlovoStrategy(logFactory logger.Factory) *GlovoStrategy {
	return &GlovoStrategy{
		log: logFactory.For(logger.Parser).With(zap.String("strategy", "glovo_final")),
	}
}

func (s *GlovoStrategy) Name() string { return "Glovo" }

func (s *GlovoStrategy) CanParse(u string) bool {
	return strings.Contains(u, "glovoapp.com")
}

func (s *GlovoStrategy) Parse(target retailer.Store) (*[]parser.ScrapedProduct, error) {
	items := make([]parser.ScrapedProduct, 0)
	s.fetchData(&items, "/v4/stores/726/addresses/164859/content/main?nodeType=DEEP_LINK&link=ofertas-black-friday-sc.43611987/drogueria-higiene-y-salsas-c.43611795\\", target)

	//urls, err := s.getLink(target)
	//if err != nil || len(urls) < 1 {
	//	s.log.Info("No parsing urls found", zap.String("url", targetURL))
	//	return nil, nil
	//}

	//for _, path := range urls {
	//	fmt.Println(path)
	//	s.fetchData(&items, path, target)
	//}

	return &items, nil
}

func (s *GlovoStrategy) getLink(target retailer.Store) ([]string, error) {
	var rawBuilder strings.Builder

	c := NewCollector([]string{"glovoapp.com", "www.glovoapp.com"})

	c.OnHTML("script", func(e *colly.HTMLElement) {
		if strings.Contains(e.Text, "self.__next_f.push") {
			rawBuilder.WriteString(e.Text)
			rawBuilder.WriteString("\n")
		}
	})

	s.log.Info("visiting", zap.String("url", target.Url))
	if err := c.Visit(target.Url); err != nil {
		return nil, err
	}

	fullData := rawBuilder.String()
	if len(fullData) == 0 {
		return nil, fmt.Errorf("no script data found")
	}

	// 1. Regex to find the path.
	// Note: Added handling to ensure we capture enough context
	re := regexp.MustCompile(`[^"]*nodeType=DEEP_LINK[^"]*-sc[^"]*`)
	rawMatches := re.FindAllString(fullData, -1)

	var cleanMatches []string
	seen := make(map[string]bool)

	for _, m := range rawMatches {
		// 2. CRITICAL FIX: Unescape the JSON string
		// This converts "\u0026" -> "&" and "\/" -> "/"
		cleaned, err := strconv.Unquote(`"` + m + `"`)
		if err != nil {
			// If standard unquote fails (because m isn't perfectly quoted),
			// fallback to manual string replacement which is often safer for regex partials
			cleaned = strings.ReplaceAll(m, `\u0026`, "&")
			cleaned = strings.ReplaceAll(cleaned, `\/`, "/")
		}

		// Ensure the path starts with /
		if !strings.HasPrefix(cleaned, "/") {
			cleaned = "/" + cleaned
		}

		if !seen[cleaned] {
			cleanMatches = append(cleanMatches, cleaned)
			seen[cleaned] = true
		}
	}

	return cleanMatches, nil
}

func (g *GlovoStrategy) fetchData(items *[]parser.ScrapedProduct, path string, target retailer.Store) {
	var url string

	// Ensure a path is clean before appending
	path = strings.TrimSpace(path)

	if strings.HasPrefix(path, "/v4") {
		url = "https://api.glovoapp.com" + path
	} else {
		url = "https://api.glovoapp.com/v3" + path
	}

	url = strings.ReplaceAll(url, "v4", "v3")
	url = strings.ReplaceAll(url, "content/main", "content")
	url = strings.ReplaceAll(url, "content/", "content")

	url, _ = strings.CutSuffix(url, "\\")
	url, _ = strings.CutSuffix(url, "/")

	// Initialize the collector
	c := NewCollector([]string{"api.glovoapp.com"})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("glovo-api-version", "14")
		r.Headers.Set("glovo-app-platform", "web")
		r.Headers.Set("glovo-app-type", "customer")
		r.Headers.Set("glovo-location-city-code", target.City)
		r.Headers.Set("glovo-location-country-code", target.Country.Code)
	})

	c.OnResponse(func(r *colly.Response) {
		// Parse the raw bytes with gjson
		result := gjson.ParseBytes(r.Body)

		// Debug: Check if we got an error in the JSON response even with 200 OK
		if result.Get("error").Exists() {
			g.log.Error("Glovo API Error", zap.String("msg", result.Get("error.message").String()))
			return
		}

		// Iterate over data.body array
		bodyItems := result.Get("data.body").Array()

		for _, item := range bodyItems {
			elementsRaw := item.Get("data.elements")
			if !elementsRaw.Exists() {
				g.log.Info("Elements empty for", zap.String("msg", item.Get("id").String()))
				continue
			}
			g.transformItems(items, elementsRaw.Array())
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(r.Request.Headers)
		g.log.Error("Request failed",
			zap.String("url", r.Request.URL.String()),
			zap.Error(err),
		)
	})

	// Execute the request
	err := c.Visit(url)
	if err != nil {
		g.log.Error("Error visiting url", zap.Error(err))
	}
}

func (g *GlovoStrategy) transformItems(items *[]parser.ScrapedProduct, elements []gjson.Result) {
	for _, raw := range elements {
		m := raw.Get("data")

		if !m.Get("imageUrl").Exists() || m.Get("imageUrl").String() == "" {
			continue
		}

		price := int(m.Get("priceInfo.amount").Float() * 100)
		oldPrice := int(m.Get("price").Float() * 100)

		newItem := parser.ScrapedProduct{
			ExternalID: m.Get("externalId").String(),
			Name:       m.Get("name").String(),
			Price:      int64(price),
			OldPrice:   int64(oldPrice),
			ImageURL:   m.Get("imageUrl").String(),
			ScrapedAt:  time.Now(),
		}

		*items = append(*items, newItem)
	}
}
