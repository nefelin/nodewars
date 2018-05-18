package protocol

import (
	"argument"
	"commandinfo"
	"fmt"
	"nwmodel"
)

// var commandList = map[string]lobbyCommand{
var commandList = LobbyCommandGroup{
	"test": {
		Info: commandinfo.Info{
			Name:      "test",
			ShortDesc: "Simply testing our new command struct",
			ArgsReq:   argument.ArgList{
			// {Name: "reqd_str", Type: argument.String},
			},
			ArgsOpt: argument.ArgList{
				{Name: "opt_int", Type: argument.Int},
			},
		},
		handler: cmdTest,
	},
}

func cmdTest(p *nwmodel.Player, d *Dispatcher, args []interface{}) error {
	// cArgs, err := lc.ValidateArgs(args)
	// if err != nil {
	// 	return err
	// }

	fmt.Println("TEST")
	return nil
}
