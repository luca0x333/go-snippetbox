package main

import (
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	// Middleware chain using alice package.
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// recoverPanic <-> logRequest <-> secureHeaders <-> servemux <-> application handler
	// When the request comes in, it will be passed to m1, then m2, then m3
	// and finally, the given handler
	return standardMiddleware.Then(mux)
}
