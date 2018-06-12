package docs

import "help"

func RegisterTopics(r *help.Registry) {
	helpTopics := []help.Topic{
		help.NewTopic(
			"doctest",
			"blah blah blah",
			"help", "topics", "other shit",
		),
	}

	for _, t := range helpTopics {
		err := r.AddEntry(t)
		if err != nil {
			panic(err)
		}
	}
	// r.AddEntry(help.NewTopic("doctest", "a test of the topical help system", "help", "system", "topical"))
}
