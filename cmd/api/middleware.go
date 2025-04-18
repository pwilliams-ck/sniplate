package main

import (
	"fmt"
	"net/http"

	"golang.org/x/time/rate"
)

// Below is an example of how these middleware functions works.
// func (app *application) exampleMiddleware(next http.Handler) http.Handler {
//
//     // Any code here will run only once, when we wrap something with the middleware.
//
//     return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//
//         // Any code here will run for every request that the middleware handles.
//
//         next.ServeHTTP(w, r)
//     })
// }

// Sets headers for incoming requests, we can set these as environment variables for the server
// config if needed.
func commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Content Security Policy
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self'")
		// Deny use of page in <frame>, <iframe>, or <object>
		w.Header().Set("X-Frame-Options", "deny")
		// Block pages from loading when they detect reflected cross-site scripting (XSS) attacks
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		// Force communication using HTTPS, preventing HTTP use
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		// Configure same-origin policy
		w.Header().Set("Access-Control-Allow-Origin", "same-origin")
		// Referrer-Policy
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		// Disallow the browser from MIME-sniffing a response away from the declared content-type
		w.Header().Set("X-Content-Type-Options", "nosniff")
		// Cross-Origin Opener Policy
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		next.ServeHTTP(w, r)
	})
}

// Logs the details of each incoming request. It captures and logs the client's IP address,
// the protocol used, the HTTP method, and the requested URI.
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)

		app.logger.Info("Request received", "ip", ip, "proto", proto, "method", method, "uri", uri)

		next.ServeHTTP(w, r)
	})
}

// This function handles unexpected errors during request processing. If a panic occurs, the function
// intercepts it, recovers normal execution flow, closes the connection, and sends an error response
// to the client, ensuring that the server can continue to handle other requests gracefully.
// If you don't close the connection after a panic, the client that sent the request might hang or
// wait indefinitely for a response. Basically , this closes the connection and sends the error message.
func (app *application) gracefulRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Defer ensures the following function is executed after the surrounding function.
		defer func() {
			// The recover built-in function allows a program to manage behavior of a
			// panicking goroutine. Executing a call to recover inside a deferred
			// function (but not any function called by it) stops the panicking sequence
			// by restoring normal execution and retrieves the error value passed to the
			// call of panic().
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// rateLimit limits the number of requests a client is allowed to make during a specific set of time.
func (app *application) rateLimit(next http.Handler) http.Handler {
	// Init a new rate limiter which allows an average of 2 requests per second, with a max of 4
	// requests in a single burst.
	limiter := rate.NewLimiter(2, 4)

	// The function we are returning is a closure, which 'closes over' the limiter variable.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Call limiter.Allow() to see if the request is permitted, and if it's not,
		// then we call the rateLimitExceededResponse() helper to return a 429 Too Many
		// Requests response (we will create this helper in a minute).
		if !limiter.Allow() {
			app.rateLimitExceededResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
