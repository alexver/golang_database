package storage

import (
	"fmt"

	"go.uber.org/zap"
)

type StorageInterface interface {
	Set(key string, value string) error
	Get(key string) (string, bool, error)
	Del(key string) error
}

type Storage struct {
	engine EngineInterface
	logger *zap.Logger
}

func CreateStorage(engine EngineInterface, logger *zap.Logger) (StorageInterface, error) {
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

func (s *Storage) Set(key string, value string) error {
	s.logger.Debug(fmt.Sprintf("storage SET command: %s = %s", key, value))

	err := s.engine.Set(key, value)
	if err != nil {
		s.logger.Error(fmt.Sprintf("storage SET command error: %s", err))
	}

	return err
}

func (s *Storage) Get(key string) (string, bool, error) {
	s.logger.Debug(fmt.Sprintf("storage GET command: %s", key))

	value, ok, err := s.engine.Get(key)
	if err != nil {
		s.logger.Error(fmt.Sprintf("storage GET command error: %s", err))
	}

	return value, ok, err
}

func (s *Storage) Del(key string) error {
	s.logger.Debug(fmt.Sprintf("storage DEL command: %s", key))

	err := s.engine.Delete(key)
	if err != nil {
		s.logger.Error(fmt.Sprintf("storage DEL command error: %s", err))
	}

	return err
}
