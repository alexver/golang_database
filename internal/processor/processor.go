package database

import (
	"github.com/alexver/golang_database/internal/compute/parser"
	"github.com/alexver/golang_database/internal/storage"
)

type ProcessorInterface interface {
	Name() string
	Suports(query *parser.Query) bool
	Process(query *parser.Query) (any, error)
}

type Processor struct {
	storage storage.StorageInterface
}
