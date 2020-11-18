package main

import "github.com/eiliz/snippetbox/pkg/models"

// Define a templateData type to act as the holding structure for any dynamic
// data that we want to pass to our HTML templates.
// This is needed because Go's html/template pkg accepts a single item of
// dynamic data
type templateData struct {
	Snippet *models.Snippet
}
