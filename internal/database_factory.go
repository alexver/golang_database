package database

import (
	"fmt"

	"github.com/alexver/golang_database/internal/compute"
	"github.com/alexver/golang_database/internal/compute/analyzer"
	"github.com/alexver/golang_database/internal/compute/parser"
	"github.com/alexver/golang_database/internal/config"
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
	compute.RegisterAnalyzer(analyzer.NewGet())
	compute.RegisterAnalyzer(analyzer.NewSet())
	compute.RegisterAnalyzer(analyzer.NewDel())
	compute.RegisterAnalyzer(analyzer.NewPing())
	compute.RegisterAnalyzer(analyzer.NewExit())

	db := NewDatabase(compute, logger)
	db.RegisterProcessor(database.NewGetProcessor(storage))
	db.RegisterProcessor(database.NewSetProcessor(storage))
	db.RegisterProcessor(database.NewDelProcessor(storage))
	db.RegisterProcessor(database.NewExitProcessor())
	db.RegisterProcessor(database.NewPingProcessor())

	return db, nil
}

func CreateServerDatabase(dbEngineConfig *config.DbEngineConfig, logger *zap.Logger) (*Database, error) {

	var dbEngine storage.EngineInterface
	switch dbEngineConfig.Type {
	case storage.ENGINE_IN_MEMORY:
		dbEngine = engine.CreateEngine()
	default:
		logger.Error(fmt.Sprintf("DB Engine is not valid '%s'", dbEngineConfig.Type))
	}

	storage, err := storage.CreateStorage(dbEngine, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("Storage creation error: %s", err.Error()))

		return nil, err
	}

	compute := compute.CreateComputeLayer(parser.CreatePaser(logger), logger)
	compute.RegisterAnalyzer(analyzer.NewGet())
	compute.RegisterAnalyzer(analyzer.NewSet())
	compute.RegisterAnalyzer(analyzer.NewDel())
	compute.RegisterAnalyzer(analyzer.NewPing())

	db := NewDatabase(compute, logger)
	db.RegisterProcessor(database.NewGetProcessor(storage))
	db.RegisterProcessor(database.NewSetProcessor(storage))
	db.RegisterProcessor(database.NewDelProcessor(storage))
	db.RegisterProcessor(database.NewPingProcessor())

	return db, nil
}
