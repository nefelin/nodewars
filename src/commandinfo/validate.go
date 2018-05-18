package commandinfo

import (
	"argtype"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func typeMismatch(need argtype.Type, got string) error {
	return fmt.Errorf("type mismatch, need %s, got '%s'", strings.ToLower(fmt.Sprint(need)), got) // TODO couldn't figure out how to better convert type...
}

// ValidateArgs tarkes the command info and a slice of strings and checks to ensure that arguement requirements are not violated
func (i Info) ValidateArgs(args []string) ([]interface{}, error) {

	if len(args) < len(i.ArgsReq) {
		return nil, errors.New("too few arguments")
	}

	if len(args) > len(i.ArgsReq)+len(i.ArgsOpt) {
		return nil, errors.New("too many arguments")
	}

	converted := make([]interface{}, len(args))
	combinedTypes := append(i.ArgsReq, i.ArgsOpt...)

	for i, arg := range args {

		kind := combinedTypes[i]
		switch kind {

		case argtype.Int:
			num, err := strconv.Atoi(arg)
			if err != nil {
				return nil, typeMismatch(kind, arg)
			}

			converted[i] = num

		case argtype.Float:
			num, err := strconv.ParseFloat(arg, 64)
			if err != nil {
				return nil, typeMismatch(kind, arg)
			}

			converted[i] = num

		case argtype.Bool:
			b, err := strconv.ParseBool(arg)
			if err != nil {
				return nil, typeMismatch(kind, arg)
			}

			converted[i] = b

		case argtype.String:
			converted[i] = arg

		default:
			return nil, fmt.Errorf("validation of type, '%s', unsupported", kind)
		}
	}

	return converted, nil
}
