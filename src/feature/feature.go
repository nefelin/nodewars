package feature

import "errors"

// Type ...
type Type interface {
	implementsType()
}

type featureType string

func (f featureType) implementsType() {}

// FromString ...
func FromString(s string) (Type, error) {
	t, ok := stringMap[s]
	if !ok {
		return nil, errors.New("Unknown feature type")
	}
	return t, nil
}

const (
	POE       featureType = "poe"
	Cloak     featureType = "cloak"
	Firewall  featureType = "firewall"
	Overclock featureType = "overclock"
	None      featureType = "none"
	// NA        featureType = ""
)

var stringMap = map[string]featureType{
	"poe":       POE,
	"cloak":     Cloak,
	"firewall":  Firewall,
	"overclock": Overclock,
	"none":      None,
}
