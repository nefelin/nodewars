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
	value interface{}
	Type  argType
}

type testDesc struct {
	id        string
	desc      string
	protoFunc funcDesc
	inputs    [][]ioExpect // length must be same as len(inputs) of funcDesc
	outputs   []ioExpect
}
