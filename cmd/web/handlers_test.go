package main

import (
	"net/http"
	"testing"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")

	// Check request status code
	if code != http.StatusOK {
		t.Errorf("want %q; got %q", http.StatusOK, code)
	}

	// Check request body
	if string(body) != "OK" {
		t.Errorf("want body equal %q", "OK")
	}
}
