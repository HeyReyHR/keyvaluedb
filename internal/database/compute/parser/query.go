package parser

type Query struct {
	commandId int
	arguments []string
}

func NewQuery(commandId int, args []string) Query {
	return Query{
		commandId: commandId,
		arguments: args,
	}
}

func (q *Query) CommandId() int {
	return q.commandId
}
func (q *Query) Arguments() []string {
	return q.arguments
}
