package main

// look into go enums
type argType string

type argDesc struct {
	name string
	Type argType
}

type funcDesc struct {
	name    string
	inputs  []argDesc
	outputs []argType
}

type ioExpect struct {
	value string
	Type  argType
}

type testDesc struct {
	id        string
	desc      string
	protoFunc funcDesc
	inputs    []ioExpect
	outputs   []ioExpect
}
