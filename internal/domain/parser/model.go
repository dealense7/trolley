package parser

import "time"

type CountryCode string

const (
	CountryES CountryCode = "ES"
	CountryUK CountryCode = "UK"
	CountryIT CountryCode = "IT"
	CountryFR CountryCode = "FR"
)

type ScrapedProduct struct {
	ExternalID     string
	Name           string
	NormalizedName *string
	Price          int64
	OldPrice       int64
	ImageURL       string
	ScrapedAt      time.Time
}

type Strategy interface {
	Name() string
	CanParse(url string) bool
	Parse(target TargetStore) (*[]ScrapedProduct, error)
}

type TargetStore struct {
	Name    string
	Country CountryCode
	City    string
	URL     string
}
