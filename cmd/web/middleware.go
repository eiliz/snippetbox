package main

import (
	"fmt"
	"net/http"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// X-Frame-Options does the same thing as Content-Security-Policy:
		// frame-ancestors 'none'. But it's supported in older browsers
		// where CSP isn't.
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("Content-Security-Policy", "frame-ancestors 'none'")

		// The HTTP X-XSS-Protection response header  stops pages from loading
		// when they detect reflected cross-site scripting (XSS) attacks.
		// Although these protections are largely unnecessary in modern browsers
		// when sites implement a strong Content-Security-Policy that disables
		// the use of inline JavaScript ('unsafe-inline'), they can still provide
		// protections for users of older web browsers that don't yet support CSP.
		w.Header().Set("X-XSS-Protection", "1;mode=block")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// A deferred func will always be called, even in the event of a panic
		// as Go unwinds the stack.
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
