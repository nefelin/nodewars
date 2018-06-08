package command

import (
	"help"
	"nwmessage"
	"room"
)

type Command interface {
	help.Helper
	Executer
	Validator
}

type Executer interface {
	Exec(cli nwmessage.Client, context room.Room, args []interface{}) error
}

type Validator interface {
	Validate(args []string) ([]interface{}, error)
}
