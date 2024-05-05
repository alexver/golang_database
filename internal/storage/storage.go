package storage

import (
	"fmt"

	"github.com/alexver/golang_database/internal/storage/engine"
	"go.uber.org/zap"
)

type StorageInterface interface {
	Set(string, string)
	Get(string) (string, bool)
	Del(string)
}

type Storage struct {
	engine engine.EngineInterface
	logger *zap.SugaredLogger
}

func CreateStorage(engine engine.EngineInterface, logger *zap.SugaredLogger) (*Storage, error) {
	if engine == nil {
		return nil, fmt.Errorf("engine is required")
	}
	if logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	return &Storage{
		engine: engine,
		logger: logger,
	}, nil
}

func (s *Storage) Set(key string, value string) {
	s.logger.Debugf("storage SET command: %s = %s", key, value)

	err := s.engine.Set(key, value)
	if err != nil {
		s.logger.Errorf("storage SET command error: %s", err)
	}
}

func (s *Storage) Get(key string) (string, bool) {
	s.logger.Debugf("storage GET command: %s", key)

	value, ok, err := s.engine.Get(key)
	if err != nil {
		s.logger.Errorf("storage GET command error: %s", err)
	}

	return value, ok
}

func (s *Storage) Det(key string) {
	s.logger.Debugf("storage DEL command: %s", key)

	err := s.engine.Delete(key)
	if err != nil {
		s.logger.Errorf("storage DEL command error: %s", err)
	}
}
