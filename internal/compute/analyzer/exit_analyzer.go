package analyzer

import (
	"fmt"
	"os"

	"github.com/alexver/golang_database/internal/compute/parser"
	"go.uber.org/zap"
)

const COMMAND_EXIT_1 = "EXIT"
const COMMAND_EXIT_2 = "QUIT"

type Exit struct {
	Analyzer
}

func NewExit(logger *zap.Logger) AnalyzerInterface {
	return &Exit{
		Analyzer: Analyzer{
			name:        COMMAND_EXIT_1,
			description: "Command to stop and close test database. You can use QUIT as an alias of EXIT command.",
			usage:       fmt.Sprintf("%s|%s", COMMAND_EXIT_1, COMMAND_EXIT_2),
			logger:      logger,
		}}
}

func (c *Exit) Name() string {

	return c.name
}

func (c *Exit) Description() string {

	return c.description
}

func (c *Exit) Usage() string {

	return c.usage
}

func (c *Exit) Supports(name string) bool {

	return name == COMMAND_EXIT_1 || name == COMMAND_EXIT_2
}

func (c *Exit) Validate(query parser.Query) error {
	if !c.Supports(query.GetCommand()) {
		return fmt.Errorf("analyzer EXIT error: cannot process '%s' command", query.GetCommand())
	}

	if len(query.GetArguments()) > 0 {
		return fmt.Errorf("analyzer EXIT error: invalid argumnet count %d", len(query.GetArguments()))
	}

	return nil
}

func (c *Exit) Run(query parser.Query) (any, error) {
	os.Exit(0)

	return 0, nil
}
