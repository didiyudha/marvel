package exec

import (
	"context"
	"fmt"
	"os"
)

type CommandExecutor interface {
	Exec(cmd string) error
}

type commandExecImpl struct {
	MarvelCmdExecutor
}

func NewCommandExecutor(executor MarvelCmdExecutor) CommandExecutor {
	return &commandExecImpl{
		MarvelCmdExecutor: executor,
	}
}

func (c *commandExecImpl) Exec(cmd string) error {
	switch cmd {
	case "init":
		return c.InitializeMarvelCharacter(context.Background())
	default:
		fmt.Println("Command not supported")
		os.Exit(1)
	}
	return nil
}
