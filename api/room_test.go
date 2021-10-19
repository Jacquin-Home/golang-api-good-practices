package api

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
)

// Test possible availability states
func TestRoom_IsValid(t *testing.T) {

	tests := map[string]Room{
		"free": Room{
			Availability: "free",
			ID:           uuid.New(),
		},
		"reserved": Room{
			Availability: "reserved",
			ID:           uuid.New(),
		},
		"inuse": Room{
			Availability: "free",
			ID:           uuid.New(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.IsValid()
			if !reflect.DeepEqual(true, got) {
				t.Errorf("wanted %t, got %t", true, got)
			}
		})
	}
}
