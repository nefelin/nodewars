package command

import (
	"argument"
	"fmt"
	"room"
	"strings"
)

type Info struct {
	CmdName     string // strings that can be used to call this Command
	CmdContexts []room.Type

	ShortDesc string
	LongDesc  string

	ArgsReq argument.ArgList
	ArgsOpt argument.ArgList

	// Handler func(nwmessage.Client, receiver.Receiver, []interface{}) error
}

func (i Info) Name() string {
	return i.CmdName
}

func (i Info) Contexts() []room.Type {
	return i.CmdContexts
}

func (i Info) SupportsContext(r room.Type) bool {
	if len(i.CmdContexts) == 0 { // not Contexts mean its global
		return true
	}

	for _, c := range i.CmdContexts {
		if r == c {
			return true
		}
	}

	return false
}

func (i Info) ShortHelp() string {
	// padding := strings.Repeat(" ", longest-len(i.Name))
	padding := 7

	return fmt.Sprintf("%-*s - %s", padding, i.Name, i.ShortDesc)
}

func (i Info) LongHelp() string {
	ret := fmt.Sprintf("%s\nusage: %s", i.ShortHelp(), i.Usage())
	if i.LongDesc != "" {
		ret += "\n\n" + i.LongDesc
	}
	return ret
}

func (i Info) Usage() string {
	required := make([]string, len(i.ArgsReq))
	optional := make([]string, len(i.ArgsOpt))

	for i, arg := range i.ArgsReq {
		required[i] = arg.Name
	}

	for i, arg := range i.ArgsOpt {
		optional[i] = fmt.Sprintf("[%s]", arg.Name)
	}

	// prevent unwanted spaces
	var reqStr, optStr string
	if len(required) > 0 {
		reqStr = strings.Join(required, " ")
	}
	if len(optional) > 0 {
		optStr = strings.Join(optional, " ")
	}
	argStr := reqStr + optStr

	return strings.TrimSpace(strings.Join([]string{i.CmdName, argStr}, " "))
}
