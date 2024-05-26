package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pwilliams-ck/sniplate/internal/data"
	"github.com/pwilliams-ck/sniplate/internal/validator"
)

// For now we simply return a plain-text placeholder response.
func (app *application) createSnipHandler(w http.ResponseWriter, r *http.Request) {
	// Declare an anonymous struct to hold the information that we expect to be in the
	// HTTP request body. This struct will be our *target decode destination*.
	var input struct {
		Title   string   `json:"title"`   // Snip title
		Content string   `json:"content"` // Content of the snip
		Tags    []string `json:"tags"`    // Slice of tags for the snip
	}

	// Initialize a new json.Decoder instance which reads from the request body, and
	// then use the Decode() method to decode the body contents into the input struct.
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Copy the values from the input struct into a new Snip struct.
	snip := &data.Snip{
		Title:   input.Title,
		Content: input.Content,
		Tags:    input.Tags,
	}

	// Init new Validator instance.
	v := validator.New()

	// Use ValidateMovie() to return any validation errors.
	if data.ValidateSnip(v, snip); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Dump contents of input struct in a HTTP response.
	fmt.Fprintf(w, "%+v\n", input)
}

// "GET /v1/snips/:id" endpoint. For now, we retrieve the "id" parameter from the
// current URL and include it in a placeholder response.
func (app *application) showSnipHandler(w http.ResponseWriter, r *http.Request) {
	// Read ID from URL param.
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
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

	// Write the response, passing the envelope defined in helpers.go.
	err = app.writeJSON(w, http.StatusOK, envelope{"snip": snip}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
