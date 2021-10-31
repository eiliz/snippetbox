package main

import (
	"net/http"
	"testing"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t)

	// We use the httpTest.NewTLSServer func to create a test server usinng the
	// handler returned by app.routes.
	// This starts a new HTTPS server on a random port number for the duration of
	// the test. The server needs to be closed with ts.Close() when the test finishes.
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")

	if code != http.StatusOK {
		t.Errorf("want %d, got %d", http.StatusOK, code)
	}

	if string(body) != "OK" {
		t.Errorf("want body equal to %q", "OK")
	}
}
