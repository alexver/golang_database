package analyzer

import (
	"github.com/alexver/golang_database/internal/compute/parser"
	"github.com/alexver/golang_database/internal/storage"
)

type AnalyzerInterface interface {
	Name() string
	Description() string
	Usage() string
	Supports(name string) bool
	Validate(query *parser.Query) error
	NormalizeQuery(query *parser.Query) *parser.Query
}

type Analyzer struct {
	storage storage.StorageInterface
}
