// Code generated by "stringer -type paramType argument.go"; DO NOT EDIT.

package param

import "strconv"

const _paramType_name = "IntFloatStringBool"

var _paramType_index = [...]uint8{0, 3, 8, 14, 18}

func (i paramType) String() string {
	if i < 0 || i >= paramType(len(_paramType_index)-1) {
		return "paramType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _paramType_name[_paramType_index[i]:_paramType_index[i+1]]
}
