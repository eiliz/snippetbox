package main

import (
	"fmt"
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		// send a 404 response to the client
		http.NotFound(w, r)

		// if we didn't return here, the handler would continue to execute the following lines
		return
	}

	// Go will sniff the content type without setting this header
	// But in the case of JSON we need to explicitly do it as the sniffer
	// mistakenly sets it to text/plain

	// The header name is auto canonicalized by Go
	// for HTTP2 conns, Go sets header names and values to lowercase as per the HTTP2 spec
	w.Header().Set("content-type", "application/json")

	// Set the header directly into the header map to escape the auto canon. behavior
	w.Header()["X-XSS-Protection"] = []string{"1;mode=block"}

	// Go auto sets the content-type, content-length and date headers
	// You cannot remove them with Del() but you can directly set them to nil inside http.Header()
	w.Header()["Date"] = nil

	fmt.Printf("%#v\t%v", w.Header(), w.Header())
	w.Write([]byte(`{"message": "Hello from Snippetbox"}`))
}

func showSnippet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a specific snippet page"))
}

func createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// Go only allows you to call WriteHeader once per response
		// The first call to Write will call WriteHeader with a 200 OK automatically
		// before writing any data, if WriteHeader has not been called yet
		// So WriteHeader is usually present when you want to send errors and must come before Write

		// Header gives you access to the response header map to add new headers
		// and writing to it must always come before calls to WriteHeader and Write
		// or else the headers added will be ignored
		w.Header().Set("Allow", http.MethodPost)

		// // WriteHeader will actually send the client the status code header
		// w.WriteHeader(405)
		// // Then the client gets his body data as well
		// w.Write([]byte("Method not allowed"))

		http.Error(w, "Method not allowed", 405)
		return
	}

	w.Write([]byte("Creating a new snippet"))
}

func test(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("test"))
}

func testa(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("test a"))
}

// can be started with the module name
// go run github.com/eiliz/snippetbox
func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)
	mux.HandleFunc("/test/", test)

	log.Println("Starting server on http://localhost:4000/")

	// start a new web server listening on :4000 and using the mux as its router
	// http.ListenAndServe always returns non-nil errors so there's no need to check

	// ":http" ":http-alt" are called named ports and Go will look up the actual
	// number by reading the /etc/services file
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
