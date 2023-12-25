package main

import (
	"fmt"
	"net/http"
	"time"

	"greenlight.architsproject/internal/data"
)

func (app *application) CreateMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new movie")
}

func (app *application) ShowMovieHandler(w http.ResponseWriter, r *http.Request) {
	id := "id"
	id_main, err := app.readParamAsInt(&id, r)
	if err != nil || id_main < 1 {
		http.NotFound(w, r)
		return
	}

	movie := data.Movie{
		ID:        id_main,
		CreatedAt: time.Now(),
		Title:     "Casablanca",
		Runtime:   102,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}

	err = app.convertDataToJson(w, http.StatusOK, envelope{"movie": movie}, nil)

	if err != nil {
		app.logger.Print(err)
		http.Error(w, "Error converting the id to json, please recheck json data", http.StatusInternalServerError)
	}

}
