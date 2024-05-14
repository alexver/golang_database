package database

import (
	"os"

	"github.com/alexver/golang_database/internal/compute/parser"
)

const PROCESSOR_NAME_EXIT = "EXIT"

type ExitProcessor struct {
	Processor
}

func NewExitProcessor() ProcessorInterface {
	return &ExitProcessor{
		Processor: Processor{},
	}
}

func (p *ExitProcessor) Name() string {
	return PROCESSOR_NAME_EXIT
}

func (p *ExitProcessor) Suports(query *parser.Query) bool {
	return query.GetCommand() == PROCESSOR_NAME_EXIT
}

func (p *ExitProcessor) Process(query *parser.Query) (any, error) {

	p.exit()

	return "Good Buy!", nil
}

func (p *ExitProcessor) exit() {

	os.Exit(0)
}
