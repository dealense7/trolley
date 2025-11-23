package parser

import (
	"fmt"
	"storePrices/internal/domain/retailer"
	"storePrices/internal/platform/logger"

	"go.uber.org/zap"
)

type ParserService struct {
	log        *zap.Logger
	strategies []Strategy
}

func NewParserService(logFactory logger.Factory) *ParserService {
	return &ParserService{
		log:        logFactory.For(logger.Parser),
		strategies: []Strategy{},
	}
}

func (s *ParserService) AddStrategy(st Strategy) {
	s.strategies = append(s.strategies, st)
}

// ScrapeAndPrint just logs the data, no saving
func (s *ParserService) ScrapeAndPrint(target retailer.Store) error {
	var str Strategy
	for _, st := range s.strategies {
		if st.CanParse(target.Url) {
			str = st
			break
		}
	}
	if str == nil {
		return fmt.Errorf("no strategy found for %s", target.Url)
	}

	products, err := str.Parse(target)
	if err != nil {
		return err
	}

	type nameStruct struct {
		name string
		id   string
	}

	var namesToNormalize []nameStruct

	for _, p := range *products {
		namesToNormalize = append(namesToNormalize, nameStruct{
			name: p.Name,
			id:   p.ExternalID,
		})

	}
	fmt.Printf("%+v\n", namesToNormalize)

	s.log.Info("--- END BATCH ---", zap.Int("total_items", len(*products)))

	return nil
}
