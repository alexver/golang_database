package database

import (
	"fmt"

	"github.com/alexver/golang_database/internal/compute"
	"github.com/alexver/golang_database/internal/compute/analyzer"
	"github.com/alexver/golang_database/internal/compute/parser"
	database "github.com/alexver/golang_database/internal/processor"
	"github.com/alexver/golang_database/internal/storage"
	engine "github.com/alexver/golang_database/internal/storage/engine/in-memory"
	"go.uber.org/zap"
)

func CreateCLIDatabase() (*Database, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	// in-memory engine
	engine := engine.CreateEngine()

	storage, err := storage.CreateStorage(engine, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("Storage creation error: %s", err.Error()))

		return nil, err
	}

	compute := compute.CreateComputeLayer(parser.CreatePaser(logger), logger)
	compute.RegisterAnalyzer(analyzer.NewGet(storage))
	compute.RegisterAnalyzer(analyzer.NewSet(storage))
	compute.RegisterAnalyzer(analyzer.NewDel(storage))
	compute.RegisterAnalyzer(analyzer.NewExit())

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
