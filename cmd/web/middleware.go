package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/justinas/nosurf"
	"github.com/luca0x333/go-snippetbox/pkg/models"
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

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If the user is not authenticated, redirect to the login page and return from the middleware chain so the
		// other handlers are not executed.
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		// Set the header "Cache-Control" to "no-store" so that pages require authentication are not stored in the users
		// browser cache.
		w.Header().Add("Cache-Control", "no-store")

		// Call next handler in the chain.
		next.ServeHTTP(w, r)
	})
}

// NoSurf is a middleware function which uses a customized CSRF cookie with
// the Secure, Path and HttpOnly flags set.
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})

	return csrfHandler
}

// authenticate middleware fetches the user's ID from their session, checks the database to see if the ID is valid
// and the user is active, then updates the request context to include this information.
func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if a authenticatedUserID value exists in the session.
		// If it does not exist call the next handler in the chain.
		exists := app.session.Exists(r, "authenticatedUserID")
		if !exists {
			next.ServeHTTP(w, r)
			return
		}

		// Fetch the details of the current user from the database.
		// If no record is found or the user is deactivated, remove "authenticatedUserID" value from their session
		// and call the next handler in the chain.
		user, err := app.users.Get(app.session.GetInt(r, "authenticatedUserID"))
		if errors.Is(err, models.ErrNoRecord) || !user.Active {
			app.session.Remove(r, "authenticatedUserID")
			next.ServeHTTP(w, r)
			return
		} else if err != nil {
			app.serverError(w, err)
			return
		}

		// If the request is coming from an authenticated and active user, we create a new copy of the request adding
		// "contextKeyIsAuthenticated" true and call the next handler in the chain using the new copy of the request.
		ctx := context.WithValue(r.Context(), contextKeyIsAuthenticated, true)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
