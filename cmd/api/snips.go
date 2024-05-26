package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pwilliams-ck/sniplate/internal/data"
)

// For now we simply return a plain-text placeholder response.
func (app *application) createSnipHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Create a new snip.")
}

// "GET /v1/snips/:id" endpoint. For now, we retrieve the "id" parameter from the
// current URL and include it in a placeholder response.
func (app *application) showSnipHandler(w http.ResponseWriter, r *http.Request) {
	// Read ID from URL param.
	id, err := app.readIDParam(r)
	if err != nil {
		app.logger.Error(err.Error())
		http.NotFound(w, r)
		return
	}

	// New instance of snip struct, containing the ID we extracted from the URL.
	snip := data.Snip{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Test",
		Content:   "Content test",
		Tags:      []string{"test", "api-dev"},
		Version:   1,
	}

	// Write the response.
	err = app.writeJSON(w, http.StatusOK, snip, nil)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "The server encontered a problem and could not process your request", http.StatusInternalServerError)
	}
}
