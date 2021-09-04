package main

import (
	"net/http"

	"github.com/eiliz/snippetbox/pkg/nfs"

	"github.com/justinas/alice"

	"github.com/bmizerany/pat"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	mux := pat.New()

	mux.Get("/", http.HandlerFunc(app.home))
	// Register the exactly matched paths (snippet/create) before the snippet/:id
	// because that one would match first otherwise.
	mux.Get("/snippet/create", http.HandlerFunc(app.createSnippetForm))
	mux.Post("/snippet/create", http.HandlerFunc(app.createSnippet))
	mux.Get("/snippet/:id", http.HandlerFunc(app.showSnippet))

	// This removes the leading / from the URL path of the req and then starts
	// looking for the asset inside the dir
	fs := http.FileServer(nfs.NeuteredFileSystem{Fs: http.Dir(app.config.staticDir)})
	mux.Get("/static/", http.StripPrefix("/static", fs))

	return standardMiddleware.Then(mux)
}
