package commands

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
func (c Command) ValidateArgs(args []string) ([]interface{}, error) {

	if len(args) < len(c.ArgsReq) {
		return nil, errors.New("too few arguments")
	}

	converted := make([]interface{}, len(args))
	combinedArgs := append(c.ArgsReq, c.ArgsOpt...)

	for i, arg := range args {

		if i > len(combinedArgs)-1 {
			return nil, fmt.Errorf("too many arguments, '%s'", arg)
		}

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
			if arg == "" {
				return nil, typeMismatch(wantArg, arg)
			}
			converted[i] = arg

		case argument.GreedyString:
			if arg == "" {
				return nil, typeMismatch(wantArg, arg)
			}
			converted[i] = strings.Join(args[i:], " ")
			return converted, nil

		default:
			return nil, fmt.Errorf("validation of type, '%s', unsupported", wantArg.Type)
		}
	}

	return converted, nil
}
