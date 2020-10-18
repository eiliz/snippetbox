package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
)

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)

	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")

		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()

			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}

// can be started with the module name
// go run github.com/eiliz/snippetbox
func main() {
	addr := flag.String("addr", ":4000", "HTTP newtwork address")

	// Call Parse before using any flag var or you'll get the default values
	flag.Parse()

	mux := http.NewServeMux()

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	// This removes the leading / from the URL path of the req and then
	// starts looking for the asset inside the dir
	fs := http.FileServer(neuteredFileSystem{http.Dir("./ui/static")})
	mux.Handle("/static/", http.StripPrefix("/static", fs))

	log.Printf("Starting server on http://localhost%s/<3/", *addr)

	// Start a new web server listening on :4000 and using the mux as its router
	// http.ListenAndServe always returns non-nil errors so there's no need to check

	// ":http" ":http-alt" are called named ports and Go will look up the actual
	// number by reading the /etc/services file
	// Every req gets its own goroutine -> they run concurrently so there are potential
	// race conditions to be careful of if they access shared resources
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}
