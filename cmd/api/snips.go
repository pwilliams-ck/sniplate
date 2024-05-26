package main

import (
	"fmt"
	"net/http"
)

// For now we simply return a plain-text placeholder response.
func (app *application) createSnipHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Create a new snip.")
}

// "GET /v1/snips/:id" endpoint. For now, we retrieve the "id" parameter from the
// current URL and include it in a placeholder response.
func (app *application) showSnipHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	id, err := app.readIdParam(r)
	if err != nil {
		app.logger.Error(err.Error())
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "snip id: %d\n", id)
}
