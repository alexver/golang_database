package analyzer

import (
	"fmt"

	"github.com/alexver/golang_database/internal/compute/tools"
	"github.com/alexver/golang_database/internal/query"
)

const COMMAND_GET = "GET"
const COMMAND_GET_ARG_COUNT = 1

type Get struct {
}

func NewGet() *Get {
	return &Get{}
}

func (g *Get) Name() string {
	return COMMAND_GET
}

func (g *Get) Description() string {
	return fmt.Sprintf("Get saved value by key. Usage: %s.", g.Usage())
}

func (g *Get) Usage() string {
	return fmt.Sprintf("%s key", COMMAND_GET)
}

func (g *Get) Supports(name string) bool {
	return name == COMMAND_GET
}

func (g *Get) Validate(query *query.Query) error {
	if !g.Supports(query.GetCommand()) {
		return fmt.Errorf("analyzer GET error: cannot process '%s' command", query.GetCommand())
	}

	if query.GetArgumentCount() != COMMAND_GET_ARG_COUNT {
		return fmt.Errorf("analyzer GET error: invalid argument count %d", query.GetArgumentCount())
	}

	if !tools.ValidateArgument(query.GetArguments()[0]) {
		return fmt.Errorf("analyzer GET error: invalid argument #%d: %s", 1, query.GetArguments()[0])
	}

	return nil
}

func (g *Get) NormalizeQuery(query *query.Query) *query.Query {

	query.SetCommand(COMMAND_GET)

	return query
}
