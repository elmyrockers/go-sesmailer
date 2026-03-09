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
		{"email with name", "helmi@xeno.com.my", "Helmi Aziz", "\"Helmi Aziz\" <helmi@xeno.com.my>"},
		{"email only", "helmi@xeno.com.my", "", "helmi@xeno.com.my"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatAddress(tt.email, tt.fullname) //set
			if got != tt.expected { //check
				t.Errorf("expected %s, got %s", tt.expected, got)
			}
		})
	}
}

func TestMailer_SetFrom(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		fullname string
	}{
		{"email with name", "helmi@xeno.com.my", "Helmi Aziz"},
		{"email only", "helmi@xeno.com.my", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Mailer{}
			m.SetFrom(tt.email, tt.fullname) //set

			//Check
				if m.From != tt.email {
					t.Errorf("expected From %s but got %s", tt.email, m.From)
				}
				if m.FromName != tt.fullname {
					t.Errorf("expected FromName %s but got %s", tt.fullname, m.FromName)
				}
		})
	}
}