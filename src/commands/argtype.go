package commands

// import ()

//go:generate stringer -type=argType
type argType int

const (
	Int argType = iota
	Float
	String
	Bool
)
