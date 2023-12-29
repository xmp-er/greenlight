package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.HealthCheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.CreateMovieHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.ShowMovieHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.UpdateMovieHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.DeleteMovieHandler)

	return router
}
