package command

// Store is the minimal interface a queue must satisfy.
type Store interface {
	Push(item string)
	Pop() (string, bool)
}

// Command represents a protocol command that can recognise and execute itself.
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
