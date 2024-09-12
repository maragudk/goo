package model_test

import (
	"testing"

	"maragu.dev/is"

	"maragu.dev/goo/model"
)

func TestEmail_IsValid(t *testing.T) {
	tests := []struct {
		address string
		valid   bool
	}{
		{"me@example.com", true},
		{"@example.com", false},
		{"me@", false},
		{"@", false},
		{"", false},
		{"me@example", false},
	}
	t.Run("reports valid email addresses", func(t *testing.T) {
		for _, test := range tests {
			t.Run(test.address, func(t *testing.T) {
				e := model.Email(test.address)
				is.Equal(t, test.valid, e.IsValid())
			})
		}
	})
}
