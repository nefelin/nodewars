package param

//go:generate stringer -type=paramType

// Type ...
type Type interface {
	implementsType()
}

type paramType int

func (a paramType) implementsType() {}

// Int, Float, String, Bool... used for type checking params
const (
	Int paramType = iota
	Float
	String
	GreedyString
	Bool
)
