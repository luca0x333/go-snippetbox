package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecureHeaders(t *testing.T) {
	// Initialize a new http.ResponseRecorder and a dummy http.Request
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock http handler that we can pass to out secureHeaders middleware.
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Pass the mock http handler to secureHeaders middleware.
	// It returns a http.Handler so we can call it's ServeHTTP method passing in the http.ResponseRecorder
	// and the dummy http.Request.
	secureHeaders(next).ServeHTTP(rr, r)

	// Call the Result() method
	rs := rr.Result()

	// Check that the middleware has correctly set the X-Frame-Options header.
	frameOptions := rs.Header.Get("X-Frame-Options")
	if frameOptions != "deny" {
		t.Errorf("want %q; got %q", "deny", frameOptions)
	}

	// Check that the middleware has correctly called set the X-XSS-Protection header.
	xssProtection := rs.Header.Get("X-XSS-Protection")
	if xssProtection != "1; mode=block" {
		t.Errorf("want %q; got %q", "1; mode=block", xssProtection)
	}

	// Check that the middleware has correctly called the next handler in line
	// and the response status code and body are expected.
	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %q; got %q", http.StatusOK, rs.StatusCode)
	}

	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "OK" {
		t.Errorf("want body equal %q", "OK")
	}
}
