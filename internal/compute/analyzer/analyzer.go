package analyzer

import (
	"github.com/alexver/golang_database/internal/compute/parser"
	"go.uber.org/zap"
)

type AnalyzerInterface interface {
	Name() string
	Description() string
	Usage() string
	Supports(name string) bool
	Validate(query parser.Query) error
	Run(query parser.Query) (any, error)
}

type Analyzer struct {
	name        string
	description string
	usage       string

	logger *zap.Logger
}
