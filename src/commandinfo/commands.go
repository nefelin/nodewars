package commandinfo

import (
	"argument"
	"fmt"
	"strings"
)

type Info struct {
	Name string // strings that can be used to call this Command

	// Usage     string
	ShortDesc string
	LongDesc  string

	ArgsReq argument.ArgList
	ArgsOpt argument.ArgList

	// ArgEval func(args []string, p *nwmodel.Player, gm *nwmodel.GameModel) ([]interface{}, error)
	//handler func(p *nwmodel.Player, gm *nwmodel.GameModel, args []interface{}) error
}

// Help provides composed help info for the command
func (i Info) Help() string {
	return fmt.Sprintf("%s - %s", i.Usage, i.ShortDesc)
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
