package analyzer

import (
	"fmt"

	"github.com/alexver/golang_database/internal/compute/parser"
)

const COMMAND_EXIT_1 = "EXIT"
const COMMAND_EXIT_2 = "QUIT"

type Exit struct {
}

func NewExit() AnalyzerInterface {
	return &Exit{}
}

func (c *Exit) Name() string {

	return COMMAND_EXIT_1
}

func (c *Exit) Description() string {

	return "Command to stop and close test database. You can use QUIT as an alias of EXIT command."
}

func (c *Exit) Usage() string {

	return fmt.Sprintf("%s|%s", COMMAND_EXIT_1, COMMAND_EXIT_2)
}

func (c *Exit) Supports(name string) bool {

	return name == COMMAND_EXIT_1 || name == COMMAND_EXIT_2
}

func (c *Exit) Validate(query *parser.Query) error {
	if !c.Supports(query.GetCommand()) {
		return fmt.Errorf("analyzer EXIT error: cannot process '%s' command", query.GetCommand())
	}

	if len(query.GetArguments()) != 0 {
		return fmt.Errorf("analyzer EXIT error: invalid argument count %d", len(query.GetArguments()))
	}

	return nil
}

func (c *Exit) NormalizeQuery(query *parser.Query) *parser.Query {
	query.SetCommand(COMMAND_EXIT_1)

	return query
}
