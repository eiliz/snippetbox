package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/eiliz/snippetbox/pkg/models"
	"github.com/justinas/nosurf"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// X-Frame-Options does the same thing as Content-Security-Policy:
		// frame-ancestors 'none'. But it's supported in older browsers
		// where CSP isn't.
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("Content-Security-Policy", "frame-ancestors 'none'")

		// The HTTP X-XSS-Protection response header stops pages from loading
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

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		// Add the the "Cache-Control: no-store" header so that pages that require
		// authentication are not stores in the user's browser (or other
		// intermediary cache).
		w.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		exists := app.session.Exists(r, "authenticatedUserID")
		if !exists {
			next.ServeHTTP(w, r)
			return
		}

		user, err := app.users.Get(app.session.GetInt(r, "authenticatedUserID"))
		if errors.Is(err, models.ErrNoRecord) || !user.Active {
			app.session.Remove(r, "authenticatedUserID")
			next.ServeHTTP(w, r)
			return
		} else if err != nil {
			app.serverError(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyIsAuthenticated, true)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

// This is a middleware used to add CSRF support. A special CSRF token is
// generated and sent as a cookie but also included as a hidden field in
// the forms that require protection. Later on the cookie's value is compared
// against the form input hidden value to see if they're a match.

// It uses a customized CSRF cookie with the Secure, Path and HttpOnly flags set.
func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})

	return csrfHandler
}
