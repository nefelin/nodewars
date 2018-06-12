package docs

import (
	"fmt"
	"help"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func RegisterTopics(r *help.Registry) {

	helpTopics := loadDocs()

	for _, t := range helpTopics {
		t.Clean() // order seeAlso
		err := r.AddEntry(t)
		if err != nil {
			panic(err)
		}
	}
}

func loadDocs() []help.Topic {
	raw, err := ioutil.ReadFile("./docs/docs.yaml")
	if err != nil {
		panic(err)
	}

	var h []help.Topic
	yaml.Unmarshal(raw, &h)
	fmt.Printf("Unmasrhalled %v\n", h)

	return h
}
