package main

import (
	"fmt"

	database "github.com/alexver/golang_database/internal"
	"github.com/alexver/golang_database/internal/config"
	"github.com/alexver/golang_database/internal/logger"
	"github.com/alexver/golang_database/internal/network"
)

func main() {

	// для простоты хардкод, можно получить и через переменную окружения
	configFileName := "configs/db_server.yml"
	dbConfig := config.LoadConfig(configFileName)

	logger := logger.CreateLogger(dbConfig.GetLoggerConfig())
	defer logger.Sync()

	db, err := database.CreateServerDatabase(dbConfig.GetDbEngineConfig(), logger)
	if err != nil {
		logger.Error(fmt.Sprintf("Database initialization error: %s", err))
	}

	server, err := network.CreateServer(dbConfig.GetNetworkServerConfig(), db, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("Database server run error: %s", err))
	}
	server.StartServer()
}
