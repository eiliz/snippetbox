package models

import (
	"errors"
	"time"
)

// ErrNoRecord represents a custom error for a missing snippet
var ErrNoRecord = errors.New("models: no matching record found")

// Snippet represents the snippet object
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}
