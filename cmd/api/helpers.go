package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

type envelope map[string]any

// writeJSON() helper for sending responses. This takes the destination
// http.ResponseWriter, the HTTP status code, the data to encode, and a
// header map containing any additional HTTP headers we want to include in the response.
func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	// Encode to JSON, return any errors, if any.
	// In benchmarks json.MarshalIndent() takes 65% longer to run and uses around 30% more memory
	// than json.Marshal(), as well as making more heap allocations.
	// We’re talking about a few thousandths of a millisecond. If your service is operating in
	// a very resource-constrained environment, then this is worth being aware of.
	json, err := json.MarshalIndent(data, "", "\t")
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

// setupLogger configures the logging output based on the useLog flag.
func setupLogger(useLog bool) io.Writer {
	var logWriter io.Writer
	if useLog {
		// Open a file for writing logs if useLog is true
		logFile, err := os.OpenFile("sniplate.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Error opening log file:", err)
			os.Exit(1)
		}
		// Use a multiwriter to write logs to both standard output and the log file
		logWriter = io.MultiWriter(os.Stdout, logFile)
	} else {
		// If useLog is false, write logs only to standard output
		logWriter = os.Stdout
	}
	return logWriter
}
