package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDBConfig_GetLoggerConfig(t *testing.T) {
	level := "debug"
	file := "tmp/my_log.log"

	t.Run("run GetLoggerConfig", func(t *testing.T) {
		d := &DBConfig{
			loggerConfig: LoggerConfig{
				Level:  level,
				Output: file,
			},
			networkServerConfig: NetworkServerConfig{},
			dbEngineConfig:      DbEngineConfig{},
		}
		got := d.GetLoggerConfig()

		require.Equal(t, level, got.Level)
		require.Equal(t, file, got.Output)
	})
}

func TestDBConfig_GetNetworkServerConfig(t *testing.T) {
	network := "http"
	host := "test.com"
	port := 4321
	maxCon := 250

	t.Run("run GetNetworkServerConfig", func(t *testing.T) {
		d := &DBConfig{
			loggerConfig: LoggerConfig{},
			networkServerConfig: NetworkServerConfig{
				Network:        network,
				Host:           host,
				Port:           port,
				MaxConnections: maxCon,
			},
			dbEngineConfig: DbEngineConfig{},
		}
		got := d.GetNetworkServerConfig()

		require.Equal(t, network, got.Network)
		require.Equal(t, host, got.Host)
		require.Equal(t, port, got.Port)
		require.Equal(t, maxCon, got.MaxConnections)
	})
}

func TestDBConfig_GetDbEngineConfig(t *testing.T) {
	engine := "test_engine"

	t.Run("get DbEngineConfig", func(t *testing.T) {
		d := &DBConfig{
			loggerConfig:        LoggerConfig{},
			networkServerConfig: NetworkServerConfig{},
			dbEngineConfig: DbEngineConfig{
				Type: engine,
			},
		}
		got := d.GetDbEngineConfig()

		require.Equal(t, engine, got.Type)
	})
}
