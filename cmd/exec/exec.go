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
	MarvelCmdExecutor MarvelCmdExecutor
	TableMigrator TableMigrator
}

func NewCommandExecutor(executor MarvelCmdExecutor, migrator TableMigrator) CommandExecutor {
	return &commandExecImpl{
		MarvelCmdExecutor: executor,
		TableMigrator: migrator,
	}
}

func (c *commandExecImpl) Exec(cmd string) error {
	switch cmd {
	case "characters":
		return c.MarvelCmdExecutor.InitializeMarvelCharacter(context.Background())
	case "migration":
		return c.TableMigrator.MigrateTable(context.Background())
	case "delete":
		return c.TableMigrator.DeleteAll(context.Background())
	default:
		fmt.Println("Command not supported")
		os.Exit(1)
	}
	return nil
}
