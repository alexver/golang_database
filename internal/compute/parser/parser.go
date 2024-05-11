package parser

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
)

type ParserInterface interface {
	ParseStringToQuery(command string) Query
}

type Parser struct {
	logger *zap.Logger
}

func CreatePaser(logger *zap.Logger) ParserInterface {

	return &Parser{logger: logger}
}

func (p *Parser) ParseStringToQuery(command string) Query {
	if strings.Trim(command, " \t\n") == "" {
		return CreateQuery("", []string{})
	}

	items := strings.Fields(command)

	args := []string{}
	if len(items) > 1 {
		args = items[1:]
	}
	result := CreateQuery(items[0], args)

	p.logger.Debug(fmt.Sprintf("Parsed command '%s' into query structure '%v'", command, result))

	return result
}
