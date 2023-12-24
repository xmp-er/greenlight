package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (app *application) readParamAsInt(s *string, r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	*s = params.ByName(*s)

	param_int, err := strconv.ParseInt(*s, 10, 64)

	if err != nil || param_int < 1 {
		err_desc := "Unable to convert parameter to int_64 format because " + err.Error() + " ,please check if paramter is greater than 1, it is currently :" + *s
		return 0, errors.New(err_desc)
	}

	return param_int, nil
}
