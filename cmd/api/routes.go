package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	// Initialize a new httprouter router instance.
	mux := http.NewServeMux()

	// Register the relevant methods, URL patterns and handler functions for our
	// endpoints using the HandleFunc() method.
	mux.HandleFunc("GET /v1/healthcheck", app.healthcheckHandler)

	mux.HandleFunc("GET /v1/snips", app.listSnipsHandler)
	mux.HandleFunc("POST /v1/snips", app.createSnipHandler)
	mux.HandleFunc("GET /v1/snips/{id}", app.showSnipHandler)
	mux.HandleFunc("PATCH /v1/snips/{id}", app.updateSnipHandler)
	mux.HandleFunc("DELETE /v1/snips/{id}", app.deleteSnipHandler)

	// Return mux router with middleware.
	return app.gracefulRecovery(app.logRequest((commonHeaders(mux))))
}
