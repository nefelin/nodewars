package validation

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func Validate(typeList []string, args []string) ([]interface{}, error) {
	errList := make([]string, len(typeList))
	converted := make([]interface{}, len(typeList))

	for i, kind := range typeList {
		if i > len(args)-1 {
			errList[i] = fmt.Sprint("Ran out of arguments")
			continue
		}

		switch kind {

		case "Int":
			num, err := strconv.Atoi(args[i])
			if err != nil {
				errList[i] = err.Error()
				continue
			}

			converted[i] = num

		case "Float":
			num, err := strconv.ParseFloat(args[i], 64)
			if err != nil {
				errList[i] = err.Error()
				continue
			}

			converted[i] = num

		case "Bool":
			b, err := strconv.ParseBool(args[i])
			if err != nil {
				errList[i] = err.Error()
				continue
			}

			converted[i] = b

		case "String":
			converted[i] = args[i]

		default:
			errList[i] = fmt.Sprintf("Validation of type, '%s', unsupported\n", kind)
		}
	}

	for _, e := range errList {
		if e != "" {
			return converted, errors.New(strings.Join(errList, ", "))
		}
	}

	return converted, nil
}
