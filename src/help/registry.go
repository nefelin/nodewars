package help

import (
	"fmt"
	"room"
)

type Registry map[string][]Helper

func NewRegistry() Registry {
	return make(Registry)
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
		helpBucket = append(helpBucket, h)
		return nil
	}

	// if name not in use, create a new slot and add this command to it
	r[hName] = []Helper{h}
	return nil
}

func (r Registry) allHelp(c room.Type) string {
	str := ""
	for _, bucket := range r {
		for _, h := range bucket {
			if h.SupportsContext(c) {
				str += h.ShortHelp()
			}
		}
	}
	return str
}

func (r Registry) Help(c room.Type, subj string) string {
	if subj == "" {
		return r.allHelp(c)
	}

	bucket, ok := r[subj]

	if !ok {
		return fmt.Sprintf("No help found for query, '%s'", subj)
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
