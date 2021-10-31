package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/eiliz/snippetbox/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
)

type config struct {
	addr      string
	staticDir string
	dsn       string
	secret    string
}

// Define an application struct to hold app wide dependencies like loggers or
// other items that might otherwise be needed as global vars.
// Because all the handlers are in the same package we can define the functions
// as methods against this struct to make sure they have access to the loggers
// or other deps from this struct.
type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *sessions.Session
	snippets      *mysql.SnippetModel
	users         *mysql.UserModel
	templateCache map[string]*template.Template
}

type contextKey string

const contextKeyIsAuthenticated = contextKey("isAuthenticated")

// can be started with the module name
// go run github.com/eiliz/snippetbox
func main() {
	cfg := new(config)
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.dsn, "dsn", "web:testing@/snippets?parseTime=true", "MySQL data source name")
	flag.StringVar(&cfg.secret, "secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key to encrypt cookies - 32 bytes long")

	// The SQL driver requires '?parseTime=true' in the DSN to be able to
	// automatically transform TIME and DATE fields to time.Time objects.

	// Call Parse before using any flag var or you'll get the default values
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(cfg.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(cfg.secret))
	session.Lifetime = 12 * time.Hour
	// The default value for the SameSite attribute of the session cookie is
	// "Lax". If we made it "Strict", a logged in user that's being redirected
	// to our app from a 3rd party would initially be treated as not logged in.
	// After navigating to another page they'd be treated as logged in.
	// session.SameSite=http.SameSiteStrictMode

	app := application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		session:       session,
		snippets:      &mysql.SnippetModel{DB: db},
		users:         &mysql.UserModel{DB: db},
		templateCache: templateCache,
	}

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true, // Prefer Go's cyper suite to the user's
		// These 2 elliptic curves have assembly implementations which makes them
		// more efficient and therefore the server performance will be better
		// under heavy loads
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

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
		Addr:         cfg.addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on http://localhost%s/<3/", cfg.addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")

	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
