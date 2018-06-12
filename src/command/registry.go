package command

import (
	"fmt"
	"help"
	"nwmessage"
	"room"
	"strings"
)

type Registry struct {
	commands map[string][]Command
	helpReg  *help.Registry
}

func NewRegistry(helpReg *help.Registry) *Registry {
	return &Registry{
		commands: make(map[string][]Command),
		helpReg:  helpReg,
	}
}

func (r *Registry) AddEntry(c Command) error {

	// try to register in help
	err := r.helpReg.AddEntry(c)
	if err != nil {
		return err
	}

	cmdName := c.Name()

	// If that names in use, confirm it's in a different context
	if cmdBucket, exists := r.commands[cmdName]; exists {

		for _, bc := range cmdBucket {
			if contextOverlap(bc.Contexts(), c.Contexts()) {
				return fmt.Errorf("error registering command, multiple commands named, '%s', cannot overlap contexts", cmdName)
			}
		}

		r.commands[cmdName] = append(cmdBucket, c)
		return nil
	}

	// if name not in use, create a new slot and add this command to it
	r.commands[cmdName] = []Command{c}
	return nil
}

func (r Registry) Exec(context room.Room, m nwmessage.ClientMessage) error {
	// m.Sender.Outgoing(nwmessage.PsAlert("Heard you!"))
	// fmt.Println("registry is execing")

	fullCmd := strings.Split(strings.TrimSpace(m.Data), " ")
	cmdString := fullCmd[0]
	strArgs := fullCmd[1:]

	// if players in chatmode and context supports yelling
	// cmd, yellingEnabled := cg["yell"]
	// if cmdString != "chat" && yellingEnabled && m.Sender.ChatMode() {
	// 	err := cmd.Exec(m.Sender, context, fullCmd)
	// 	if err != nil {
	// 		m.Sender.Outgoing(nwmessage.PsError(err))
	// 	}
	// 	return nil
	// }

	// handle help // TODO register this as a command
	if cmdString == "help" {
		m.Sender.Outgoing(nwmessage.PsNeutral(r.helpReg.Help(context.Type(), strArgs)))
		return nil
	}

	// does this command exist?
	if bucket, ok := r.commands[cmdString]; ok {
		// fmt.Printf("Found Bucket, has %d entries\n", len(bucket))
		// find version of command suited to this context
		for _, cmd := range bucket {
			// fmt.Printf("Evaluating cmd, '%s', supports contexts %v\n", cmd.Name(), cmd.Contexts())
			if cmd.SupportsContext(context.Type()) {
				// are the args valid?
				args, err := cmd.Validate(strArgs)
				if err != nil {
					return fmt.Errorf("%s\nusage: %s", err.Error(), cmd.Usage())
				}

				// try to execute
				err = cmd.Exec(m.Sender, context, args)
				if err != nil {
					m.Sender.Outgoing(nwmessage.PsError(err))
				}

				return nil
			}
		}

		return fmt.Errorf("Command, '%s', not supported in this context", cmdString)

	}

	// if we don't find the command, pass an error back to caller in case caller wants to do something else
	return fmt.Errorf("Unknown command, '%s'", cmdString)
}

func contextOverlap(a, b []room.Type) bool {
	// if either is global, must overlap
	if len(a) == 0 || len(b) == 0 {
		return true
	}

	// if they have any of the same contexts, also overlap
	for _, c1 := range a {
		for _, c2 := range b {
			if c1 == c2 {
				return true
			}
		}
	}

	// we're fine
	return false
}
