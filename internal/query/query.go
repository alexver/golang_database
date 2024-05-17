package query

type Query struct {
	command   string
	arguments []string
}

func CreateQuery(command string, arguments []string) *Query {
	return &Query{
		command:   command,
		arguments: arguments,
	}
}

func (q *Query) GetCommand() string {

	return q.command
}

func (q *Query) SetCommand(command string) *Query {
	q.command = command

	return q
}

func (q *Query) GetArguments() []string {

	if q.arguments == nil {
		q.arguments = []string{}
	}

	return q.arguments
}

func (q *Query) GetArgumentCount() int {

	return len(q.arguments)
}
