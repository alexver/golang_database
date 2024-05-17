package parser

import (
	"fmt"
	"strings"

	"github.com/alexver/golang_database/internal/query"
	"go.uber.org/zap"
)

type Parser struct {
	logger *zap.Logger
}

func CreatePaser(logger *zap.Logger) *Parser {

	return &Parser{logger: logger}
}

func (p *Parser) ParseStringToQuery(command string) *query.Query {
	if strings.Trim(command, " \t\n") == "" {
		return query.CreateQuery("", []string{})
	}

	items := strings.Fields(command)

	args := []string{}
	if len(items) > 1 {
		args = items[1:]
	}
	result := query.CreateQuery(items[0], args)

	p.logger.Debug(fmt.Sprintf("Parsed command '%s' into query structure '%v'", command, result))

	return result
}
