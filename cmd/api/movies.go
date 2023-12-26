package main

import (
	"fmt"
	"net/http"
	"time"

	"greenlight.architsproject/internal/data"
)

func (app *application) CreateMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string   `json:"title"`
		Year    int32    `json:"year"`
		Runtime int32    `json:"runtime"`
		Genres  []string `json:"genres"`
	}

	err := app.readJSON(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	fmt.Fprintf(w, "%+v \n", input)
}

func (app *application) ShowMovieHandler(w http.ResponseWriter, r *http.Request) {
	id := "id"
	id_main, err := app.readParamAsInt(&id, r)
	if err != nil || id_main < 1 {
		app.logError(r, err)
		app.notFoundResponse(w, r)
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
		app.logError(r, err)
		app.logger.Println("Error converting the id to json, please recheck json data" + err.Error())
		app.serverErrorResponse(w, r)
		return
	}

}
