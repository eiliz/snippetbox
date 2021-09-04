package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/eiliz/snippetbox/pkg/models"
)

// Define a templateData type to act as the holding structure for any dynamic
// data that we want to pass to our HTML templates.
// This is needed because Go's html/template pkg accepts a single item of
// dynamic data.
type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// Get page name from the last element of the path, ie home.page.tmpl
		name := filepath.Base(page)

		// Add the page template
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Then add any layouts
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		// Also add the partials
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		// Store the template set under the page's name
		cache[name] = ts
	}

	return cache, nil
}
