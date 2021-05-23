package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	// Initialize a new time.Time object and pass it to the humanDate function.
	tm := time.Date(2021, 05, 23, 21, 41, 0, 0, time.UTC)
	hd := humanDate(tm)

	// Check if the output from humanDate is what we expect.
	if hd != "23 May 2021 at 21:41" {
		t.Errorf("want %q; got %q", "23 May 2021 at 21:41", hd)
	}
}
