package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var invalidRuntimeFormat = errors.New("invalid Runtime Format")

type Runtime int32

func (r Runtime) MarshalJSON() ([]byte, error) {
	data := fmt.Sprintf("%d mins", r)

	quoted_data := strconv.Quote(data)

	return []byte(quoted_data), nil
}

func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {
	unquotedJsonValue, err := strconv.Unquote(string(jsonValue))

	if err != nil {
		return invalidRuntimeFormat
	}

	parts := strings.Split(unquotedJsonValue, " ")

	if len(parts) != 2 || parts[1] != "mins" {
		return invalidRuntimeFormat
	}

	final, err := strconv.ParseInt(parts[0], 10, 32)

	if err != nil {
		return invalidRuntimeFormat
	}

	*r = Runtime(final) //converting to Runtime datatype

	return nil
}
