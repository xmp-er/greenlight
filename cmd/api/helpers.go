package main

import (
	"encoding/json"
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

type envelope map[string]any

func (app *application) convertDataToJson(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	data_json, err := json.MarshalIndent(data, "", "\t")

	if err != nil {
		return err
	}

	data_json = append(data_json, '\n')

	for k, v := range headers {
		w.Header()[k] = v
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data_json)

	return nil
}
