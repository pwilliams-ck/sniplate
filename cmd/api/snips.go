package main

import (
	"errors"
	"fmt"
	"net/http"

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

	// Call the Insert() method on our movies model, passing in a pointer to the
	// validated snip struct. This will create a record in the database and update the
	// movie struct with the system-generated information.
	err = app.models.Snips.Insert(snip)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// When sending a HTTP response, we want to include a Location header to let the
	// client know which URL they can find the newly-created resource at. We make an
	// empty http.Header map and then use the Set() method to add a new Location header,
	// interpolating the system-generated ID for our new snip in the URL.
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/snips/%d", snip.ID))

	// Write a JSON response with a 201 Created status code, the snip data in the
	// response body, and the Location header.
	err = app.writeJSON(w, http.StatusCreated, envelope{"snip": snip}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// "GET /v1/snips/{id}" endpoint. For now, we retrieve the "id" parameter from the
// current URL and include it in a placeholder response.
func (app *application) showSnipHandler(w http.ResponseWriter, r *http.Request) {
	// Read ID from URL param.
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// Call the Get() method to fetch the data for a specific movie. We also need to
	// use the errors.Is() function to check if it returns a data.ErrRecordNotFound
	// error, in which case we send a 404 Not Found response to the client.
	snip, err := app.models.Snips.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Write the response, passing the envelope defined in helpers.go.
	err = app.writeJSON(w, http.StatusOK, envelope{"snip": snip}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
