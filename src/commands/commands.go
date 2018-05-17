package commands

import (
	"fmt"
	"nwmodel"
)

type command struct {
	names []string // strings that can be used to call this command

	usage       string
	description string

	argsReq []argType
	argsOpt []argType
	argEval func(args []string, p *nwmodel.Player, gm *nwmodel.GameModel) ([]interface{}, error)

	handler func(p *nwmodel.Player, gm *nwmodel.GameModel, args []interface{}) error
}

func newCommand(names []string) *command {
	return &command{
		names:   names,
		argsReq: make([]argType, 0),
		argsOpt: make([]argType, 0),
	}
}

type CommandGroup struct {
	commands []*command
}

func NewCommandGroup(c []*command) *CommandGroup {
	group := &CommandGroup{
		commands: make([]*command, 0),
	}

	for _, command := range c {
		group.AddCommand(command)
	}

	return group
}

func (cg *CommandGroup) AddCommand(newCommand *command) {
	cg.commands = append(cg.commands, newCommand)
}

func (cg *CommandGroup) Validate() error {
	// make sure the namespace is clear
	seen := make(map[string]bool)

	for _, command := range cg.commands {
		for _, name := range command.names {
			if _, ok := seen[name]; ok {
				return fmt.Errorf("CommandGroup corrupt! Duplicate entries for '%s'", name)
			}
			seen[name] = true
		}
	}
	return nil
}

func (cg *CommandGroup) String() string {
	return fmt.Sprintf("Command group contains %d commands", len(cg.commands))
}
