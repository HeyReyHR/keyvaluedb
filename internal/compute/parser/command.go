package parser

const (
	UnknownCommandId = iota
	SetCommandId
	GetCommandId
	DelCommandId
)
const (
	setCommandArgs = 2
	getCommandArgs = 1
	delCommandArgs = 1
)
const (
	UnknownCommand = "UNKNOWN"
	SetCommand     = "SET"
	GetCommand     = "GET"
	DelCommand     = "DEL"
)

var namesToId = map[string]int{
	SetCommand: SetCommandId,
	GetCommand: GetCommandId,
	DelCommand: DelCommandId,
}

var idsToArgNumber = map[int]int{
	SetCommandId: setCommandArgs,
	GetCommandId: getCommandArgs,
	DelCommandId: delCommandArgs,
}

func commandNameToId(command string) int {
	id, ok := namesToId[command]
	if !ok {
		return UnknownCommandId
	}
	return id
}

func commandArgumentsNumber(commandId int) int {
	return idsToArgNumber[commandId]
}
