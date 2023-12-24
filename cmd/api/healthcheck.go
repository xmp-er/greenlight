package main

import (
	"net/http"
)

func (app *application) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response_content_map := map[string]string{
		"status":     "available",
		"enviroment": app.config.env,
		"version":    version,
	}

	err := app.convertDataToJson(w, http.StatusOK, response_content_map, nil)

	if err != nil {
		app.logger.Print(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
