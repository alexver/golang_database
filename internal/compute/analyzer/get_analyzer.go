package analyzer

import (
	"fmt"

	"github.com/alexver/golang_database/internal/compute/parser"
	"github.com/alexver/golang_database/internal/compute/tools"
	"github.com/alexver/golang_database/internal/storage"
)

const COMMAND_GET = "GET"
const COMMAND_GET_ARG_COUNT = 1

type Get struct {
	Analyzer
}

func NewGet(storage storage.StorageInterface) AnalyzerInterface {
	return &Get{
		Analyzer: Analyzer{
			storage: storage,
		}}
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

func (g *Get) Validate(query parser.Query) error {
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

func (g *Get) Run(query parser.Query) (any, error) {
	value, ok, err := g.storage.Get(query.GetArguments()[0])
	if err != nil {
		return "", err
	}

	if ok {
		return value, nil
	}

	return "[not found]", nil
}
