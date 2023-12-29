package main

import (
	"fmt"
	"net/http"
)

func (app *application) logError(r *http.Request, err any) {
	app.logger.Println(err)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(err.Error(), http.StatusBadRequest, w, r)
}

func (app *application) errorResponse(err any, status int, w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"error": err,
	}

	error := app.convertDataToJson(w, status, env, nil)

	if error != nil {
		app.logError(r, "Could not convert error data to Json because"+error.Error())
		w.WriteHeader(500)
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request) {

	app.errorResponse("the server encountered a problem and could not process your request", http.StatusInternalServerError, w, r)
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {

	app.errorResponse("the requested resource could not be found", http.StatusNotFound, w, r)
}

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(fmt.Sprintf("the %s method is not supported for this resource", r.Method), http.StatusMethodNotAllowed, w, r)
}

func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(errors, http.StatusUnprocessableEntity, w, r)
}

func (app *application) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "Unable to edit record due to update conflict,please try again"
	app.errorResponse(message, http.StatusConflict, w, r)
}
