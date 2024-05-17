package analyzer

import (
	"fmt"

	"github.com/alexver/golang_database/internal/compute/tools"
	"github.com/alexver/golang_database/internal/query"
)

const COMMAND_DEL = "DEL"
const COMMAND_DEL2 = "DELETE"
const COMMAND_DEL_ARG_COUNT = 1

type Del struct{}

func NewDel() *Del {
	return &Del{}
}

func (d *Del) Name() string {
	return COMMAND_DEL
}

func (d *Del) Description() string {
	return fmt.Sprintf("Delete value by key from storage. %s is an alias. Usage: %s.", COMMAND_DEL2, d.Usage())
}

func (d *Del) Usage() string {
	return fmt.Sprintf("%s|%s key", COMMAND_DEL, COMMAND_DEL2)
}

func (d *Del) Supports(name string) bool {
	return name == COMMAND_DEL || name == COMMAND_DEL2
}

func (d *Del) Validate(query *query.Query) error {
	if !d.Supports(query.GetCommand()) {
		return fmt.Errorf("analyzer DEL error: cannot process '%s' command", query.GetCommand())
	}

	if query.GetArgumentCount() != COMMAND_DEL_ARG_COUNT {
		return fmt.Errorf("analyzer DEL error: invalid argument count %d", query.GetArgumentCount())
	}

	if !tools.ValidateArgument(query.GetArguments()[0]) {
		return fmt.Errorf("analyzer DEL error: invalid argument #%d: %s", 1, query.GetArguments()[0])
	}

	return nil
}

func (d *Del) NormalizeQuery(query *query.Query) *query.Query {

	query.SetCommand(COMMAND_DEL)

	return query
}
