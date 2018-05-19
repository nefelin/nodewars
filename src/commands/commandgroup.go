package commands

import (
	"fmt"
	"nwmessage"
	"sort"
	"strings"
)

// CommandGroup stores a collection of commands
type CommandGroup map[string]Command

// Exec either provides command information via the 'help' comands, or tries to process a command
func (cg CommandGroup) Exec(context interface{}, m nwmessage.ClientMessage) error {
	fullCmd := strings.Split(m.Data, " ")
	cmdString := fullCmd[0]
	args := fullCmd[1:]

	// if players in chatmode and context supports yelling
	cmd, yellingEnabled := cg["yell"]
	if cmdString != "chat" && yellingEnabled && m.Sender.ChatMode() {
		err := cmd.Exec(m.Sender, context, fullCmd)
		if err != nil {
			m.Sender.Outgoing(nwmessage.PsError(err))
		}
		return nil
	}

	// handle help
	if cmdString == "help" {
		if len(fullCmd) == 1 {

			m.Sender.Outgoing(nwmessage.PsNeutral(cg.AllHelp()))

		} else {
			help, err := cg.Help(fullCmd[1:])

			if err != nil {
				m.Sender.Outgoing(nwmessage.PsError(err))
			}
			m.Sender.Outgoing(nwmessage.PsNeutral(help))

		}
		return nil

	}

	// if we find the command, try to execute
	if cmd, ok := cg[cmdString]; ok {
		err := cmd.Exec(m.Sender, context, args)
		if err != nil {
			m.Sender.Outgoing(nwmessage.PsError(err))
		}
		return nil
	}

	// if we don't find the command, pass an error back to caller in case caller wants to do something else
	return unknownCommand(fullCmd[0])
}

// Help composes a help string for the given command
func (cg CommandGroup) Help(args []string) (string, error) {
	if cmd, ok := cg[args[0]]; ok {
		return cmd.LongHelp(), nil
	}
	return "", unknownCommand(args[0])
}

// AllHelp composes help for all commands in the group
func (cg CommandGroup) AllHelp() string {
	cmds := make([]string, len(cg))
	var i int
	for key := range cg {
		cmds[i] = key
		i++
	}

	sort.Strings(cmds)
	// offset := cg.longestKey()
	helpStr := make([]string, len(cmds)+1)
	helpStr[0] = "Available commands:"

	for i, cmd := range cmds {
		helpStr[i+1] = cg[cmd].ShortHelp()
	}

	return strings.Join(helpStr, "\n")
}

func unknownCommand(cmd string) error {
	return fmt.Errorf("Unknown command, '%s'", cmd)
}
