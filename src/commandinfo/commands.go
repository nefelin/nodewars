package commandinfo

import "argtype"

type Info struct {
	// Name string // strings that can be used to call this Command

	Usage     string
	ShortDesc string
	LongDesc  string

	ArgsReq []argtype.Type
	ArgsOpt []argtype.Type

	// ArgEval func(args []string, p *nwmodel.Player, gm *nwmodel.GameModel) ([]interface{}, error)
	//handler func(p *nwmodel.Player, gm *nwmodel.GameModel, args []interface{}) error
}
