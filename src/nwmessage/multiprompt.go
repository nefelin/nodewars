package nwmessage

import (
	"errors"
	"fmt"
)

type Fn func(*Dialogue, string) Message

// Dialogue allows us to take asynchronous input from user and compose into a single object,
// evaluating and storing the data arbitrarily along the way
type Dialogue struct {
	currentStep int
	stages      []Fn
	props       map[string]string
}

// NewMultiPrompt creates a new generic object
func NewDialogue(fns []Fn) *Dialogue {
	return &Dialogue{
		stages: fns,
		props:  make(map[string]string),
	}
}

func (d *Dialogue) AddStage(fns ...Fn) {
	for _, n := range fns {
		d.stages = append(d.stages, n)
	}
}

func (d *Dialogue) SetProp(k, v string) {
	d.props[k] = v
}

func (d *Dialogue) GetProp(k string) string {
	return d.props[k]
}

func (d *Dialogue) Run(s string) Message {
	if d.currentStep >= len(d.stages) {
		return PsError(errors.New("Dialogue completed"))
	}
	return d.stages[d.currentStep](d, s)
}

func (d *Dialogue) Adv() {
	d.currentStep++
}

func (d *Dialogue) Rew() {
	d.currentStep--
}

func (d Dialogue) String() string {
	return fmt.Sprintf("<Dialogue> currentStep: %d, props: %v", d.currentStep, d.props)
}
