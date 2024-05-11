package database

import (
	"errors"

	"github.com/alexver/golang_database/internal/compute"
	"github.com/alexver/golang_database/internal/compute/analyzer"
	"github.com/alexver/golang_database/internal/storage"
	"go.uber.org/zap"
)

type Database struct {
	storage storage.StorageInterface
	compute compute.ComputeInterface

	logger *zap.Logger
}

func NewDatabase(storage storage.StorageInterface, compute compute.ComputeInterface, logger *zap.Logger) (*Database, error) {
	if storage == nil {
		return nil, errors.New("storage is undefined")
	}

	if compute == nil {
		return nil, errors.New("compute layer is undefined")
	}

	if logger == nil {
		return nil, errors.New("logger is undefined")
	}

	return &Database{
		storage: storage,
		compute: compute,
		logger:  logger,
	}, nil
}

func (db *Database) GetAnalyzers() []analyzer.AnalyzerInterface {

	return db.compute.GetAnalyzers()
}

func (db *Database) ProcessQuery(cmdString string) (any, error) {

	return db.compute.HandleQuery(cmdString)
}
