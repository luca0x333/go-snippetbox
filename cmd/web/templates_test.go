package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	// Slice of anonymous structs containing the test case name, the input for humanDate function
	// and the output expected.
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2020, 12, 17, 10, 0, 0, 0, time.UTC),
			want: "17 Dec 2020 at 10:00",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2020, 12, 17, 10, 0, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "17 Dec 2020 at 09:00",
		},
	}

	// Loop over the test cases.
	for _, tt := range tests {
		// t.Run() run a sub-test for each of the test case.
		// The first parameter is the name of the test, the second parameter is an anonymous function
		// containing the actual test for each case.
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)

			// Check if the output from humanDate is what we expect.
			if hd != tt.want {
				t.Errorf("want %q; got %q", tt.want, hd)
			}
		})
	}
}
