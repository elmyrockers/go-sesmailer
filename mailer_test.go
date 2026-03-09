package sesmailer

import (
	"testing"
	// "github.com/davecgh/go-spew/spew"
)

func TestFormatAddress(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		fullname string
		expected string
	}{
		{"with name", "helmi@xeno.com.my", "Helmi Aziz", "\"Helmi Aziz\" <helmi@xeno.com.my>"},
		{"without name", "helmi@xeno.com.my", "", "helmi@xeno.com.my"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatAddress(tt.email, tt.fullname)
			if got != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, got)
			}
		})
	}
}