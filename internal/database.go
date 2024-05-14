package database

import (
	"errors"
	"fmt"

	"github.com/alexver/golang_database/internal/compute"
	"github.com/alexver/golang_database/internal/compute/analyzer"
	database "github.com/alexver/golang_database/internal/processor"
	"github.com/alexver/golang_database/internal/storage"
	"go.uber.org/zap"
)

type Database struct {
	compute compute.ComputeInterface

	processors map[string]database.ProcessorInterface

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

	db := Database{
		compute:    compute,
		processors: make(map[string]database.ProcessorInterface),
		logger:     logger,
	}
	db.RegisterProcessor(database.NewGetProcessor(storage))
	db.RegisterProcessor(database.NewSetProcessor(storage))
	db.RegisterProcessor(database.NewDelProcessor(storage))
	db.RegisterProcessor(database.NewExitProcessor())

	return &db, nil
}

func (db *Database) RegisterProcessor(processor database.ProcessorInterface) error {
	if processor == nil {
		db.logger.Error("Register DB Processor error: Processor is not defined")

		return errors.New("register DB Processor error: Processor is not defined")
	}

	_, ok := db.processors[processor.Name()]
	if ok {
		db.logger.Error(fmt.Sprintf("Register DB Processor error: '%s' processor already registered", processor.Name()))

		return fmt.Errorf("processor '%s' already registered", processor.Name())
	}

	db.processors[processor.Name()] = processor

	return nil
}

func (db *Database) GetAnalyzers() []analyzer.AnalyzerInterface {

	return db.compute.GetAnalyzers()
}

func (db *Database) ProcessQuery(cmdString string) (any, error) {
	query, err := db.compute.HandleQuery(cmdString)
	if err != nil {
		return nil, err
	}

	for _, processor := range db.processors {
		if !processor.Suports(query) {
			continue
		}

		return processor.Process(query)
	}

	db.logger.Error(fmt.Sprintf("Query %#v cannot be processed: processor is not found", query))

	return nil, fmt.Errorf("query %s cannot be processed: processor is not found", query.GetCommand())
}
