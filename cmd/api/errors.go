package main

import (
	"net/http"
	"runtime/debug"
)

// `serverError` is a function that handles internal server errors. It logs the error, method,
// URI, and stack trace, then responds to the client with a 500 status code. This helps prevent server
// crashes and supports easier server-side debugging.
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// `clientError` responds with a given HTTP status code and its associated error message.
// This ensures a standardized response for erroneous client requests.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
