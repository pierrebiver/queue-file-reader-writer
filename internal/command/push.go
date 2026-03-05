package command

import (
	"strings"
)

type PushCommand struct{}

func (c PushCommand) CommandName() string { return "PUSH" }

func (c PushCommand) Is(line string) bool {
	return strings.HasPrefix(line, "PUSH ")
}

func (c PushCommand) Execute(line string, q Store) (string, error) {
	item := strings.TrimPrefix(line, "PUSH ")

	q.Push(item)
	return "OK", nil
}
