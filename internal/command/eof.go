package command

type EOFCommand struct{}

func (c EOFCommand) CommandName() string { return "EOF" }

func (c EOFCommand) Is(line string) bool { return line == "EOF" }

func (c EOFCommand) Execute(line string, q Store) (string, error) {
	q.Push("EOF")
	return "OK", nil
}
