package commandinfo

import (
	"argument"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func typeMismatch(need argument.Arg, got string) error {
	return fmt.Errorf("type mismatch, '%s' (%s), must be of type '%s'", got, need.Name, strings.ToLower(fmt.Sprint(need.Type))) // TODO couldn't figure out how to better convert type...
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
	combinedArgs := append(i.ArgsReq, i.ArgsOpt...)

	for i, arg := range args {

		wantArg := combinedArgs[i]
		switch wantArg.Type {

		case argument.Int:
			num, err := strconv.Atoi(arg)
			if err != nil {
				return nil, typeMismatch(wantArg, arg)
			}

			converted[i] = num

		case argument.Float:
			num, err := strconv.ParseFloat(arg, 64)
			if err != nil {
				return nil, typeMismatch(wantArg, arg)
			}

			converted[i] = num

		case argument.Bool:
			b, err := strconv.ParseBool(arg)
			if err != nil {
				return nil, typeMismatch(wantArg, arg)
			}

			converted[i] = b

		case argument.String:
			converted[i] = arg

		default:
			return nil, fmt.Errorf("validation of type, '%s', unsupported", wantArg.Type)
		}
	}

	return converted, nil
}
