package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		// Send a 404 response to the client
		http.NotFound(w, r)

		// If we didn't return here, the handler would continue to execute the following lines
		return
	}

	// Home goes *first*
	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	// Read the template files into a template set
	ts, err := template.ParseFiles(files...)

	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal server error", 500)
		return
	}

	// Execute will write the template content as the response body
	err = ts.Execute(w, nil)

	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal server error", 500)
	}
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Display a specific snippet with the ID %d\n", id)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)

		http.Error(w, "Method not allowed", 405)
		return
	}

	w.Write([]byte("Creating a new snippet"))
}
