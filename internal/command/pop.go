package command

type PopCommand struct{}

func (c PopCommand) CommandName() string { return "POP" }

func (c PopCommand) Is(line string) bool { return line == "POP" }

func (c PopCommand) Execute(line string, q Store) (string, error) {
	item := q.Pop()

	return "DATA " + item, nil
}
