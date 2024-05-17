package database

import (
	"github.com/alexver/golang_database/internal/storage"
)

type Processor struct {
	storage storage.StorageInterface
}
