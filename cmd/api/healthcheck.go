package main

import (
	"fmt"
	"net/http"
)

func (app *application) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "status:available")
	fmt.Fprintln(w, "enviroment:", app.config.env)
	fmt.Fprintln(w, "version:", version)
}
