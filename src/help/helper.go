package help

import "room"

type Helper interface {
	Name() string
	LongHelp() string
	ShortHelp() string
	Usage() string
	SupportsContext(room.Type) bool
	Contexts() []room.Type
}
