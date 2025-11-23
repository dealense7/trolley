package main

import (
	"storePrices/internal/domain/parser"
)

func GetTargets() []parser.TargetStore {
	return []parser.TargetStore{
		{
			Name:    "Dia Madrid (Calle Laguna)",
			Country: parser.CountryES,
			City:    "MAD",
			URL:     "https://glovoapp.com/es/es/madrid/stores/supermercado-dia-2",
		},
	}
}
