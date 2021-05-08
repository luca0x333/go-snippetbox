package forms

import (
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"
)

type Form struct {
	// url.Values is the same underlying type as r.PostForm map we used in createSnippet method.
	// url.Values is embedded anonymously.
	url.Values
	Errors errors
}

// New initialize a custom Form struct.
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Required method checks that specific fields in the form data are present and not blank.
// If any field fails this check, add the message to the Form errors.
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// MaxLength method checks that a specific field in the form contains a maximum number of characters.
// If the check fails then add the message to the form errors.
func (f *Form) MaxLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}

	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("This field is too long (maximum is %d characters)", d))
	}
}

// PermittedValues method check that a specific field in the form matches one of a set of specific permitted values.
// If the check fails then add the appropriate message to the form errors.
func (f *Form) PermittedValues(field string, opts ...string) {
	value := f.Get(field)
	if value == "" {
		return
	}

	for _, opt := range opts {
		if value == opt {
			return
		}
	}

	f.Errors.Add(field, "This field is invalid")
}

// Valid method returns true if there are no errors.
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
