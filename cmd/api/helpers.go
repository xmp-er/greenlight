package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"

	"greenlight.architsproject/internal/validator"
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

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	max_bytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(max_bytes)) //setting limit of max_header size

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&data)

	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError
		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}

	err = decoder.Decode(&struct{}{}) //if after JSON there is garbage value in request

	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

// read string from url parameter
func (app *application) readString(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	return s
}

func (app *application) readCSV(qs url.Values, key string, defaultValue []string) []string {
	ret := qs.Get(key)

	if ret == "" {
		return defaultValue
	}

	return strings.Split(ret, ",")
}

func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be a string")
		return defaultValue
	}
	return i
}

func (app *application) background(fn func()) {
	app.wg.Add(1)
	go func() {
		defer app.wg.Done()
		defer func() {
			if err := recover(); err != nil {
				app.logger.PrintError(fmt.Errorf("%s", err), nil)
			}
		}()

		fn()
	}()
}
