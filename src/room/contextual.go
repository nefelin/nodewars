package room

type Contextual interface {
	SupportsContext(Type) bool
	Contexts() []Type
}
