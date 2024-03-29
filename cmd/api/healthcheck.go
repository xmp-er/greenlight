package main

import (
	"net/http"
)

func (app *application) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"enviroment": app.config.env,
			"version":    version,
		},
	}

	err := app.convertDataToJson(w, http.StatusOK, env, nil)

	if err != nil {
		app.logger.PrintError(err, map[string]string{})
		app.serverErrorResponse(w, r)
	}
}
