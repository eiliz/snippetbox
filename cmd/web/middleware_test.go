package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecureHeaders(t *testing.T) {
	// Indicates that this test is safe to run concurrently with other tests.
	// Tests marked with t.Parallel() will be run in parallel with and only with
	// other parallel tests. The max no of tests run in parallel is GOMAXPROCS.
	// Can be set as a flag: “go test -parallel 4 ./...
	t.Parallel()

	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(r)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	secureHeaders(next).ServeHTTP(rr, r)
	rs := rr.Result()

	frameOptions := rs.Header.Get("X-Frame-Options")
	if frameOptions != "deny" {
		t.Errorf("want %q, got %q", "deny", frameOptions)
	}

	if xssProtection := rs.Header.Get("X-XSS-Protection"); xssProtection != "1;mode=block" {
		t.Errorf("want %q, got %q", "1;mode=block", xssProtection)
	}

	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d, got %d", 200, rs.StatusCode)
	}

	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "OK" {
		t.Errorf("want body equal to %q", "OK")
	}
}
