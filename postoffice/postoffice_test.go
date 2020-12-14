package main

import (
	"testing"
)

var readRouteCases = []struct {
	name           string
	route          string
	expectError    bool
	expectedName   string
	expectedLength int
}{
	{"unknown route", "non-existent",
		true, "", 0},
	{"known fast route", "fast-and-easy",
		false, "fast-and-easy", 2},
	{"known slow route", "slow-and-risky",
		false, "slow-and-risky", 8},
}

func TestReadRoute(t *testing.T) {
	for _, tt := range readRouteCases {
		t.Run(tt.name, func(t *testing.T) {
			r, err := readRoute(tt.route)

			if err != nil && !tt.expectError {
				t.Fatal(err)
			}

			if tt.expectError {
				if err == nil {
					t.Fatal("Expected error, got nil")
				}
				return
			}

			if r.Name != tt.expectedName {
				t.Errorf("Incorrect name. Have %q, want %q.",
					r.Name, tt.expectedName)
			}

			if r.Position != 0 {
				t.Errorf("Incorrect position. Have %d, want %d.",
					r.Position, 0)
			}

			if len(r.Slip) != tt.expectedLength {
				t.Errorf("Incorrect slip length. Have %d, want %d.",
					len(r.Slip), tt.expectedLength)
			}
		})
	}
}
