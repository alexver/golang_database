package analyzer

import (
	"fmt"

	"github.com/alexver/golang_database/internal/query"
)

const COMMAND_PING = "PING"

type Ping struct {
}

func NewPing() *Ping {
	return &Ping{}
}

func (c *Ping) Name() string {

	return COMMAND_PING
}

func (c *Ping) Description() string {

	return "Command to ping test database. Standard answer of the server is PONG."
}

func (c *Ping) Usage() string {

	return COMMAND_PING
}

func (c *Ping) Supports(name string) bool {

	return name == COMMAND_PING
}

func (c *Ping) Validate(query *query.Query) error {
	if !c.Supports(query.GetCommand()) {
		return fmt.Errorf("analyzer PING error: cannot process '%s' command", query.GetCommand())
	}

	if len(query.GetArguments()) != 0 {
		return fmt.Errorf("analyzer PING error: invalid argument count %d", len(query.GetArguments()))
	}

	return nil
}

func (c *Ping) NormalizeQuery(query *query.Query) *query.Query {
	query.SetCommand(COMMAND_PING)

	return query
}
