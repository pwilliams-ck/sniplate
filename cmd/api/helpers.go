package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

// writeJSON() helper for sending responses. This takes the destination
// http.ResponseWriter, the HTTP status code, the data to encode, and a
// header map containing any additional HTTP headers we want to include in the response.
func (app *application) writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	// Encode to JSON, return any errors, if any.
	json, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Append new line for easier reading in terminals.
	json = append(json, '\n')

	// Loop through the header map and add each header to the http.ResponseWriter header map.
	// Go doesn't throw an error if you try to range over a nil map.
	for key, value := range headers {
		w.Header()[key] = value
	}

	// Add the "Content-Type: application/json" header, then write the status code and
	// JSON response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(json)

	return nil
}

// Read snip ID URL param.
func (app *application) readIDParam(r *http.Request) (int64, error) {
	// PathValue() is new for Go 1.22 and allows us to read URL params.
	// We also convert the string to base 10 integer with a 64 bit size.
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}
