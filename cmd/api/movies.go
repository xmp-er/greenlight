package main

import (
	"errors"
	"fmt"
	"net/http"
	_ "time"

	"greenlight.architsproject/internal/data"
	"greenlight.architsproject/internal/validator"
)

func (app *application) CreateMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	err := app.readJSON(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	movie := data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	if data.ValidateMovie(v, &movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Movies.Insert(&movie)

	if err != nil {
		app.serverErrorResponse(w, r)
		return
	}

	headers := make(http.Header)

	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	err = app.convertDataToJson(w, http.StatusCreated, envelope{"movies": movie}, headers)

	if err != nil {
		app.errorResponse(err, http.StatusInternalServerError, w, r)
	}

}

func (app *application) ShowMovieHandler(w http.ResponseWriter, r *http.Request) {
	id := "id"
	id_main, err := app.readParamAsInt(&id, r)
	if err != nil || id_main < 1 {
		app.logError(r, err)
		app.notFoundResponse(w, r)
		return
	}

	movie, err := app.models.Movies.Get(id_main)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.errorResponse(err, http.StatusInternalServerError, w, r)
		}
		return
	}

	err = app.convertDataToJson(w, http.StatusOK, envelope{"movie": movie}, nil)

	if err != nil {
		app.logError(r, err)
		app.logger.Println("Error converting the id to json, please recheck json data" + err.Error())
		app.serverErrorResponse(w, r)
		return
	}

}

func (app *application) UpdateMovieHandler(w http.ResponseWriter, r *http.Request) {
	id := "id"
	id_main, err := app.readParamAsInt(&id, r)

	if err != nil {
		app.logger.Println("Error parsing id from parameter as :" + err.Error())
		app.serverErrorResponse(w, r)
		return
	}

	updated_movie, err := app.models.Movies.Get(id_main)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.logger.Println("The data for given id does not exist in the database," + err.Error())
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r)
		}
		return
	}

	var input struct {
		Title   *string       `json:"title"`
		Year    *int32        `json:"year"`
		Runtime *data.Runtime `json:"runtime"`
		Genres  []string      `json:"genres"`
	}

	err = app.readJSON(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Title != nil {
		updated_movie.Title = *input.Title
	}

	if input.Year != nil {
		updated_movie.Year = *input.Year
	}

	if input.Runtime != nil {
		updated_movie.Runtime = *input.Runtime
	}

	if input.Genres != nil {
		updated_movie.Genres = input.Genres
	}

	v := validator.New()

	if data.ValidateMovie(v, updated_movie); v.Valid() == false {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Movies.Update(updated_movie)

	if err != nil {
		app.logger.Println("Failed to update the movie")
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.errorResponse(err, http.StatusInternalServerError, w, r)
		}
		return
	}

	err = app.convertDataToJson(w, http.StatusOK, envelope{"movie": updated_movie}, nil)
	if err != nil {
		app.logger.Println(err)
		app.serverErrorResponse(w, r)
	}

}

func (app *application) DeleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	id := "id"
	id_main, err := app.readParamAsInt(&id, r)

	if err != nil {
		app.logger.Println("Unable to parse id from request")
		app.serverErrorResponse(w, r)
		return
	}

	err = app.models.Movies.Delete(id_main)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.errorResponse(err, http.StatusInternalServerError, w, r)
		}
		return
	}

	err = app.convertDataToJson(w, http.StatusOK, envelope{"message": "movie successfully deleted"}, nil)
	if err != nil {
		app.errorResponse(err, http.StatusInternalServerError, w, r)
	}
}

func (app *application) ListMoviesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string
		Genres []string
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Title = app.readString(qs, "title", "")

	input.Genres = app.readCSV(qs, "genres", []string{})

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")

	if v.Valid() == false {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}
