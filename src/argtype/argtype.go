package argtype

//go:generate stringer -type=argType

// Type ...
type Type interface {
	implementsType()
}

type argType int

func (a argType) implementsType() {}

// Int, Float, String, Bool... used for type checking arguments
const (
	Int argType = iota
	Float
	String
	Bool
)
