package config

import (
	"log"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
)

type LoggerConfig struct {
	Level  string `mapstructure:"level" default:"info"`
	Output string `mapstructure:"output" default:"/tmp/db_server.log"`
}

type NetworkServerConfig struct {
	Network        string `mapstructure:"network" default:"tcp"`
	Host           string `mapstructure:"host" default:"127.0.0.1"`
	Port           int    `mapstructure:"port" default:"8080"`
	MaxConnections int    `mapstructure:"max_connections" default:"10"`
}

type DbEngineConfig struct {
	Type string `mapstructure:"type" default:"in-memory"`
}

type DBConfig struct {
	loggerConfig        LoggerConfig
	networkServerConfig NetworkServerConfig
	dbEngineConfig      DbEngineConfig
}

func LoadConfig(configFile string) *DBConfig {
	config.AddDriver(yaml.Driver)
	config.WithOptions(config.ParseDefault)

	err := config.LoadFiles(configFile)
	if err != nil {
		log.Fatalf("Cannot load configuration file %v", err)
	}

	loggerConfig := LoggerConfig{}
	err = config.BindStruct("logger", &loggerConfig)
	if err != nil {
		log.Fatalf("Cannot bind logger configuration %v", err)
	}

	networkServerConfig := NetworkServerConfig{}
	err = config.BindStruct("network_server", &networkServerConfig)
	if err != nil {
		log.Fatalf("Cannot bind network server configuration %v", err)
	}

	dbEngineConfig := DbEngineConfig{}
	err = config.BindStruct("db_engine", &dbEngineConfig)
	if err != nil {
		log.Fatalf("Cannot bind DB engine configuration %v", err)
	}

	return &DBConfig{
		loggerConfig:        loggerConfig,
		networkServerConfig: networkServerConfig,
		dbEngineConfig:      dbEngineConfig,
	}
}

func (d *DBConfig) GetLoggerConfig() *LoggerConfig {
	return &d.loggerConfig
}

func (d *DBConfig) GetNetworkServerConfig() *NetworkServerConfig {
	return &d.networkServerConfig
}

func (d *DBConfig) GetDbEngineConfig() *DbEngineConfig {
	return &d.dbEngineConfig
}
