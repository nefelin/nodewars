package help

import (
	"room"
)

type Helper interface {
	Name() string
	Type() Type
	ShortHelp() string
	LongHelp() string
	room.Contextual
}

type Type interface {
	implementsType() helpType
}

type helpType string

func (r helpType) implementsType() helpType {
	var h helpType
	return h
}

const (
	CommandType helpType = "Command"
	TopicType   helpType = "Topic"
)
