package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
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

	fmt.Fprintln(w, "Showing the details of the movie ", id)

}
