package help

import (
	"fmt"
	"room"
	"sort"
	"strings"
)

type Registry map[string][]Helper

func NewRegistry() *Registry {
	r := make(Registry)
	return &r
}

func (r Registry) AddEntry(h Helper) error {

	hName := h.Name()

	// If that names in use, confirm it's in a different context
	if helpBucket, exists := r[hName]; exists {
		for _, bh := range helpBucket {
			if contextOverlap(bh.Contexts(), h.Contexts()) {
				return fmt.Errorf("error registering help, multiple helps named, '%s', cannot overlap contexts", hName)
			}
		}
		// helpBucket = append(helpBucket, h)
		r[hName] = append(helpBucket, h)
		return nil
	}

	// if name not in use, create a new slot and add this command to it
	r[hName] = []Helper{h}
	return nil
}

func (r Registry) allHelp(c room.Type) string {
	commandHelps := make([]string, 0)
	topicHelps := make([]string, 0)
	for _, bucket := range r {
		for _, h := range bucket {
			if h.SupportsContext(c) {
				switch h.Type() {
				case CommandType:
					commandHelps = append(commandHelps, h.ShortHelp())
				case TopicType:
					topicHelps = append(topicHelps, h.ShortHelp())
				default:
					fmt.Printf("error, unhandled help type, '%s'\n", h.Type())
				}
			}
		}
	}

	sort.StringSlice(commandHelps).Sort()
	sort.StringSlice(topicHelps).Sort()

	return fmt.Sprintf("Available Commands:\n%s\nOther Topics:\n%s", strings.Join(commandHelps, "\n"), strings.Join(topicHelps, ", "))
}

func (r Registry) Help(c room.Type, args []string) string {
	if len(args) == 0 {
		return r.allHelp(c)
	}

	bucket, ok := r[args[0]]

	if !ok {
		return fmt.Sprintf("No help found for query, '%s'", args[0])
	}

	for _, h := range bucket {
		if h.SupportsContext(c) {
			return h.LongHelp()
		}
	}

	return fmt.Sprintf("Command not available in this context")
}

func contextOverlap(a, b []room.Type) bool {
	// if either is global, must overlap
	if len(a) == 0 || len(b) == 0 {
		return true
	}

	// if they have any of the same contexts, also overlap
	for _, c1 := range a {
		for _, c2 := range b {
			if c1 == c2 {
				return true
			}
		}
	}

	// we're fine
	return false
}
