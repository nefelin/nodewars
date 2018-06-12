package help

import (
	"fmt"
	"room"
	"sort"
	"strings"
)

type Topic struct {
	TopicName string   `yaml:"name"`
	Desc      string   `yaml:"desc"`
	SeeAlso   []string `yaml:"seeAlso"`
}

func (t Topic) Name() string {
	return t.TopicName
}

func (t Topic) Type() Type {
	return TopicType
}

func (t Topic) ShortHelp() string {
	return t.TopicName
}

func (t Topic) LongHelp() string {
	return fmt.Sprintf("-%s-\n%s\nSee also: %s", t.TopicName, strings.Replace(t.Desc, "\n", "", -1), strings.Join(t.SeeAlso, ", "))
}

func (t Topic) Contexts() []room.Type {
	// Global
	return []room.Type{}
}

func (t Topic) SupportsContext(r room.Type) bool {
	return true
}

func (t *Topic) Clean() {
	sort.StringSlice(t.SeeAlso).Sort()
}
