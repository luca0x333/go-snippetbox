package main

import (
	"github.com/luca0x333/go-snippetbox/pkg/models"
	"html/template"
	"path/filepath"
	"time"
)

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates.
type templateData struct {
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
	CurrentYear int
}

// humanDate returns a nicely formatted string containing time.Time object.
func humanDate(t time.Time) string {
	return t.Format("07 May 2021 at 15:39")
}

// FuncMap is the type of the map defining the mapping from names to
// functions. Each function must have either a single return value, or two
// return values of which the second has type error.
// String-keyed map which acts as a lookup between the names of our custom template functions (names in template files)
// and the name of the functions themselves.
var functions = template.FuncMap{
	"humanDate": humanDate,
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

		// The template.FuncMap must be registered with the template set before you call ParseFiles().
		// Create a new empty template with template.New(), use Funcs() to register the template.FuncMap
		// and then parse the file.
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
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
