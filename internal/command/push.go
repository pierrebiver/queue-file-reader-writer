package command

import (
	"fmt"
	"strings"
)

type PushCommand struct{}

func (c PushCommand) CommandName() string { return "PUSH" }

func (c PushCommand) Is(line string) bool {
	return strings.HasPrefix(line, "PUSH ")
}

func (c PushCommand) Execute(line string, q Store) (string, error) {
	item := strings.TrimPrefix(line, "PUSH ")
	if item == "" {
		return "", fmt.Errorf("PUSH requires a value")
	}
	q.Push(item)
	return "OK", nil
}
