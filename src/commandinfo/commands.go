package commandinfo

import (
	"argument"
	"fmt"
	"strings"
)

type Info struct {
	Name string // strings that can be used to call this Command

	ShortDesc string
	LongDesc  string

	ArgsReq argument.ArgList
	ArgsOpt argument.ArgList
}

var padding int = 10

// Help provides composed help info for the command
func (i Info) ShortHelp() string {
	// padding := strings.Repeat(" ", longest-len(i.Name))

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

	reqStr := strings.Join(required, " ")
	optStr := strings.Join(optional, " ")

	return fmt.Sprintf("%s %s %s", i.Name, reqStr, optStr)
}
