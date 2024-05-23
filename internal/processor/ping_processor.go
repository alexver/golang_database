package database

import (
	"fmt"

	"github.com/alexver/golang_database/internal/query"
)

const PROCESSOR_NAME_PING = "PING"

const PROCESSOR_ANSWER = "PONG"

type PingProcessor struct {
	Processor
}

func NewPingProcessor() *PingProcessor {
	return &PingProcessor{
		Processor: Processor{},
	}
}

func (p *PingProcessor) Name() string {
	return PROCESSOR_NAME_PING
}

func (p *PingProcessor) Suports(query *query.Query) bool {
	return query.GetCommand() == PROCESSOR_NAME_PING
}

func (p *PingProcessor) Process(query *query.Query) (any, error) {
	return fmt.Sprintf("[ok] %s", PROCESSOR_ANSWER), nil
}
