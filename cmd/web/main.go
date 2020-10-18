package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type neuteredFileSystem struct {
	fs http.FileSystem
}

// Define an application struct to hold app wide dependencies like loggers or
// other items that might otherwise be needed as global vars.
// Because all the handlers are in the same package we can define the functions
// as methods against this struct to make sure they have access to the loggers
// or other deps from this struct.
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
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

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app := application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)

	// This removes the leading / from the URL path of the req and then starts
	// looking for the asset inside the dir

	fs := http.FileServer(neuteredFileSystem{http.Dir("./ui/static")})
	mux.Handle("/static/", http.StripPrefix("/static", fs))

	// Start a new web server listening on :4000 and using the mux as its router
	// http.ListenAndServe always returns non-nil errors so there's no need to check

	// ":http" ":http-alt" are called named ports and Go will look up the actual
	// number by reading the /etc/services file
	// Every req gets its own goroutine -> they run concurrently so there are potential
	// race conditions to be careful of if they access shared resources

	// By default Go's HTTP server uses the std logger for errors.
	// To be able to use our custom errorLog we need to instantiate a Server
	// rather than use the http.ListenAndServe shortcut which creates the server
	// and calls ListenAndServe on it.
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	infoLog.Printf("Starting server on http://localhost%s/<3/", *addr)
	err := srv.ListenAndServe()

	errorLog.Fatal(err)
}
