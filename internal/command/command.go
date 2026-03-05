package command

type Store interface {
	Push(item string)
	Pop() string
}

type Command interface {
	CommandName() string
	Is(line string) bool
	Execute(line string, q Store) (string, error)
}

var EOF = EOFCommand{}

// Registry is the list of all known commands.
// The server iterates over this to find a matching command for each line.
var Registry = []Command{
	PushCommand{},
	PopCommand{},
	EOF,
}
