package help

import (
	"fmt"
	"room"
	"strings"
)

type topicName = string

type Topic struct {
	name    topicName
	desc    string
	seeAlso []topicName
}

func NewTopic(name, desc string, seeAlso ...string) Topic {
	return Topic{
		name:    name,
		desc:    desc,
		seeAlso: seeAlso,
	}
}

func (t Topic) Name() string {
	return t.name
}

func (t Topic) Type() Type {
	return TopicType
}

func (t Topic) ShortHelp() string {
	return t.name
}

func (t Topic) LongHelp() string {
	return fmt.Sprintf("-%s-\n%s\nSee also:%s", t.name, t.desc, strings.Join(t.seeAlso, ", "))
}

func (t Topic) Contexts() []room.Type {
	// Global
	return []room.Type{}
}

func (t Topic) SupportsContext(r room.Type) bool {
	return true
}
