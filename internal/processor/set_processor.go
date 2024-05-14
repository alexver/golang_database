package database

import (
	"github.com/alexver/golang_database/internal/compute/parser"
	"github.com/alexver/golang_database/internal/storage"
)

const PROCESSOR_NAME_SET = "SET"

type SetProcessor struct {
	Processor
}

func NewSetProcessor(storage storage.StorageInterface) ProcessorInterface {
	return &SetProcessor{
		Processor: Processor{storage: storage},
	}
}

func (p *SetProcessor) Name() string {
	return PROCESSOR_NAME_SET
}

func (p *SetProcessor) Suports(query *parser.Query) bool {
	return query.GetCommand() == PROCESSOR_NAME_SET
}

func (p *SetProcessor) Process(query *parser.Query) (any, error) {
	err := p.storage.Set(query.GetArguments()[0], query.GetArguments()[1])
	if err != nil {
		return "", err
	}

	return "[ok]", nil
}
