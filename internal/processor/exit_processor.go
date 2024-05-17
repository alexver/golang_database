package database

import (
	"os"

	"github.com/alexver/golang_database/internal/query"
)

const PROCESSOR_NAME_EXIT = "EXIT"

type ExitProcessor struct {
	Processor
}

func NewExitProcessor() *ExitProcessor {
	return &ExitProcessor{
		Processor: Processor{},
	}
}

func (p *ExitProcessor) Name() string {
	return PROCESSOR_NAME_EXIT
}

func (p *ExitProcessor) Suports(query *query.Query) bool {
	return query.GetCommand() == PROCESSOR_NAME_EXIT
}

func (p *ExitProcessor) Process(query *query.Query) (any, error) {

	p.exit()

	return "Good Buy!", nil
}

func (p *ExitProcessor) exit() {

	os.Exit(0)
}
