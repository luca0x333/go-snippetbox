package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	// Middleware chain using alice package.
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Initialize a new mux using pat package.
	// Pat matches patterns in the order that they are registered.
	// We need to register GET "/snippet/create/" before GET "/snippet/:id"
	mux := pat.New()
	mux.Get("/", http.HandlerFunc(app.home))
	mux.Get("/snippet/create", http.HandlerFunc(app.createSnippetForm))
	mux.Post("/snippet/create", http.HandlerFunc(app.createSnippet))
	mux.Get("/snippet/:id", http.HandlerFunc(app.showSnippet))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	// recoverPanic <-> logRequest <-> secureHeaders <-> servemux <-> application handler
	// When the request comes in, it will be passed to m1, then m2, then m3
	// and finally, the given handler
	return standardMiddleware.Then(mux)
}
