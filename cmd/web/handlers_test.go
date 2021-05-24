package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {
	// Initialize a new httptest.ResponseRecorder.
	rr := httptest.NewRecorder()

	// Initialize a new dummy http.Request.
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Call the ping function passing in httptest.ResponseRecorder and http.Request.
	ping(rr, r)

	// Call the Result() method on the httptest.ResponseRecorder to get http.Response generated
	// by the ping() handler.
	rs := rr.Result()

	// Check if the status code in the response is 200
	if rs.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, rs.StatusCode)
	}

	// Check the response body written by the ping() handler is "OK"
	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != "OK" {
		t.Errorf("want body to equal %q", "OK")
	}
}
