package analyzer

import (
	"fmt"

	"github.com/alexver/golang_database/internal/compute/parser"
	"github.com/alexver/golang_database/internal/compute/tools"
)

const COMMAND_SET = "SET"
const COMMAND_SET_ARG_COUNT = 2

type Set struct {
}

func NewSet() AnalyzerInterface {
	return &Set{}
}

func (s *Set) Name() string {

	return COMMAND_SET
}

func (s *Set) Description() string {

	return fmt.Sprintf("Set value by key. Usage: %s.", s.Usage())
}

func (s *Set) Usage() string {

	return fmt.Sprintf("%s key value", COMMAND_SET)
}

func (s *Set) Supports(name string) bool {
	return name == COMMAND_SET
}

func (s *Set) Validate(query *parser.Query) error {
	if !s.Supports(query.GetCommand()) {
		return fmt.Errorf("analyzer SET error: cannot process '%s' command", query.GetCommand())
	}

	if query.GetArgumentCount() != COMMAND_SET_ARG_COUNT {
		return fmt.Errorf("analyzer SET error: invalid argument count %d", query.GetArgumentCount())
	}

	for i := 0; i < query.GetArgumentCount(); i++ {
		if !tools.ValidateArgument(query.GetArguments()[i]) {
			return fmt.Errorf("analyzer SET error: invalid argument #%d: %s", i+1, query.GetArguments()[i])
		}
	}

	return nil
}

func (s *Set) NormalizeQuery(query *parser.Query) *parser.Query {

	query.SetCommand(COMMAND_SET)

	return query
}
