package database

import (
	"fmt"

	"github.com/alexver/golang_database/internal/compute/parser"
	"github.com/alexver/golang_database/internal/storage"
)

const PROCESSOR_NAME_GET = "GET"

type GetProcessor struct {
	Processor
}

func NewGetProcessor(storage storage.StorageInterface) ProcessorInterface {
	return &GetProcessor{
		Processor: Processor{storage: storage},
	}
}

func (p *GetProcessor) Name() string {
	return PROCESSOR_NAME_GET
}

func (p *GetProcessor) Suports(query *parser.Query) bool {
	return query.GetCommand() == PROCESSOR_NAME_GET
}

func (p *GetProcessor) Process(query *parser.Query) (any, error) {
	value, ok, err := p.storage.Get(query.GetArguments()[0])
	if err != nil {
		return "", err
	}

	if ok {
		return fmt.Sprintf("[ok] %s", value), nil
	}

	return "[not found]", nil
}
