package main

import (
	"net/http"

	"github.com/eiliz/snippetbox/pkg/nfs"
)

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	// This removes the leading / from the URL path of the req and then starts
	// looking for the asset inside the dir
	fs := http.FileServer(nfs.NeuteredFileSystem{Fs: http.Dir(app.config.staticDir)})
	mux.Handle("/static/", http.StripPrefix("/static", fs))

	return mux
}
