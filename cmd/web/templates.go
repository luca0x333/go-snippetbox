package main

import (
	"github.com/luca0x333/go-snippetbox/pkg/models"
	"html/template"
	"path/filepath"
)

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates.
type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}

	// Glob returns the names of all files matching pattern or nil
	// if there is no matching file.
	// Get a slice of all filepaths with .page.tmpl
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	// Iterate over the pages.
	for _, page := range pages {
		// // Base returns the last element of path.
		name := filepath.Base(page)

		// Parse the page template file into a template set.
		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Add 'layout' templates to the template set.
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		// Add 'partial' templates to the template set.
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		// Add the template set to the cache map using the name of the page
		// ex: 'home.page.tmpl' as key.
		cache[name] = ts
	}

	return cache, nil
}
