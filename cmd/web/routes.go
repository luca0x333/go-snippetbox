package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	// Middleware chain using alice package.
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Create a new middleware chain containing the middleware specific to
	// our dynamic application routes.
	dynamicMiddleware := alice.New(app.session.Enable)

	// Initialize a new mux using pat package.
	// Pat matches patterns in the order that they are registered.
	// We need to register GET "/snippet/create/" before GET "/snippet/:id"
	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/snippet/create", dynamicMiddleware.ThenFunc(app.createSnippetForm))
	mux.Post("/snippet/create", dynamicMiddleware.ThenFunc(app.createSnippet))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// There is no need for this route to have stateful behaviour.
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	// recoverPanic <-> logRequest <-> secureHeaders <-> servemux <-> dynamicMiddleware <-> application handler
	// When the request comes in, it will be passed to m1, then m2, then m3
	// and finally, the given handler
	return standardMiddleware.Then(mux)
}
