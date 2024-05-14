package database

import (
	"github.com/alexver/golang_database/internal/compute/parser"
	"github.com/alexver/golang_database/internal/storage"
)

const PROCESSOR_NAME_DEL = "DEL"

type DelProcessor struct {
	Processor
}

func NewDelProcessor(storage storage.StorageInterface) ProcessorInterface {
	return &DelProcessor{
		Processor: Processor{storage: storage},
	}
}

func (p *DelProcessor) Name() string {
	return PROCESSOR_NAME_DEL
}

func (p *DelProcessor) Suports(query *parser.Query) bool {
	return query.GetCommand() == PROCESSOR_NAME_DEL
}

func (p *DelProcessor) Process(query *parser.Query) (any, error) {
	err := p.storage.Del(query.GetArguments()[0])
	if err != nil {
		return "", err
	}

	return "[ok]", nil
}
