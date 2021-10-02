package main

import (
	"net/http"

	"github.com/eiliz/snippetbox/pkg/nfs"

	"github.com/justinas/alice"

	"github.com/bmizerany/pat"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynamicMiddleware := alice.New(app.session.Enable)

	mux := pat.New()

	mux.Get("/", dynamicMiddleware.Then(http.HandlerFunc(app.home)))
	// Register the exactly matched paths (snippet/create) before the snippet/:id
	// because that one would match first otherwise.
	mux.Get("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).Then(http.HandlerFunc(app.createSnippetForm)))
	mux.Post("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).Then(http.HandlerFunc(app.createSnippet)))
	mux.Get("/snippet/:id", dynamicMiddleware.Then(http.HandlerFunc(app.showSnippet)))

	// User signup, login and logout
	mux.Get("/user/signup", dynamicMiddleware.Then(http.HandlerFunc(app.signupUserForm)))
	mux.Post("/user/signup", dynamicMiddleware.Then(http.HandlerFunc(app.signupUser)))
	mux.Get("/user/login", dynamicMiddleware.Then(http.HandlerFunc(app.loginUserForm)))
	mux.Post("/user/login", dynamicMiddleware.Then(http.HandlerFunc(app.loginUser)))
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).Then(http.HandlerFunc(app.logoutUser)))

	// This removes the leading /static from the URL path of the req and then starts
	// looking for the asset inside the dir
	fs := http.FileServer(nfs.NeuteredFileSystem{Fs: http.Dir(app.config.staticDir)})
	mux.Get("/static/", http.StripPrefix("/static", fs))

	return standardMiddleware.Then(mux)
}
