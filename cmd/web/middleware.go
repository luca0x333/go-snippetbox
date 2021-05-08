package main

import (
	"fmt"
	"net/http"
)

// secureHeaders add two Http headers to every response.
// It acts on every request that is received and it will be executed before the request hits servemux.
// secureHeaders <-> servemux <-> application handler
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}

// logRequest <-> secureHeaders <-> servemux <-> application handler
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

// recoverPanic <-> logRequest <-> secureHeaders <-> servemux <-> application handler
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			// Use recover() to check if there has been a panic.
			if err := recover(); err != nil {
				// In case of panic let's make the http server close the current connection
				// after a response has been sent.
				w.Header().Set("Connection", "Close")
				// Call serverError method to return 500 status.
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
