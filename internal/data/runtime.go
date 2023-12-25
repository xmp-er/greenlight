package data

import (
	"fmt"
	"strconv"
)

type Runtime int32

func (r Runtime) MarshalJSON() ([]byte, error) {
	data := fmt.Sprintf("%d mins", r)

	quoted_data := strconv.Quote(data)

	return []byte(quoted_data), nil
}
