package main

import (
	"bytes"
	"fmt"
	"github.com/justinas/nosurf"
	"net/http"
	"runtime/debug"
	"time"
)

// The serverError helper writes an error message and stack trace to the errorLog,
// then sends a generic 500 Internal Server Error response to the user.
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	// StatusText returns a text for the HTTP status code. It returns the empty
	// string if the code is unknown.
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	// Retrieve the appropriate template set from the cache based on the page name.
	// ex: 'home.base.tmpl'. If no entry exists calls serverError method.
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("the template %s does not exist", name))
	}

	// The new built-in function allocates memory.
	// Initialize a new buffer.
	buf := new(bytes.Buffer)

	// Write the template set to the buffer instead of http.ResponseWriter.
	// Call serverError method and return in case of an error.
	// Also inject the default data into the templateData.
	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Write the contents of the buffer to the http.ResponseWriter.
	// WriteTo() takes an io.Writer.
	buf.WriteTo(w)
}

// addDefaultData takes a pointer to templateData struct, add some default data to the struct and return it.
// and returns the pointer.
func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}

	// Add the CSRF token to the templateData struct.
	td.CSRFToken = nosurf.Token(r)

	td.CurrentYear = time.Now().Year()

	// PopString returns the string value for a given key and then deletes it from the
	// session data. One-time fetch.
	td.Flash = app.session.PopString(r, "flash")

	// Authentication status
	td.IsAuthenticated = app.isAuthenticated(r)

	return td
}

// isAuthenticated return true if "authenticatedUserID" exists in the request context.
func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(contextKeyIsAuthenticated).(bool)
	if !ok {
		return false
	}

	return isAuthenticated
}
