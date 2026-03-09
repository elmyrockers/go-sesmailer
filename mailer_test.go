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

func TestMailer_AddAddress(t *testing.T) {
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
			m := &Mailer{To: []string{}}
			m.AddAddress(tt.email, tt.fullname) //set

			//Check
				if len(m.To) != 1 {
					t.Fatalf("expected 1 address but got %d", len(m.To))
				}

				if m.To[0] != tt.expected {
					t.Errorf("expected %s but got %s", tt.expected, m.To[0])
				}
		})
	}
}

func TestMailer_AddCC(t *testing.T) {
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
			m := &Mailer{Cc: []string{}}
			m.AddCC(tt.email, tt.fullname) //set

			//Check
				if len(m.Cc) != 1 {
					t.Fatalf("expected 1 address but got %d", len(m.Cc))
				}

				if m.Cc[0] != tt.expected {
					t.Errorf("expected %s but got %s", tt.expected, m.Cc[0])
				}
		})
	}
}

func TestMailer_AddBCC(t *testing.T) {
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
			m := &Mailer{Bcc: []string{}}
			m.AddBCC(tt.email, tt.fullname) //set

			//Check
				if len(m.Bcc) != 1 {
					t.Fatalf("expected 1 address but got %d", len(m.Bcc))
				}

				if m.Bcc[0] != tt.expected {
					t.Errorf("expected %s but got %s", tt.expected, m.Bcc[0])
				}
		})
	}
}

func TestMailer_AddReplyTo(t *testing.T) {
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
			m := &Mailer{ReplyTo: []string{}}
			m.AddReplyTo(tt.email, tt.fullname) //set

			//Check
				if len(m.ReplyTo) != 1 {
					t.Fatalf("expected 1 address but got %d", len(m.ReplyTo))
				}

				if m.ReplyTo[0] != tt.expected {
					t.Errorf("expected %s but got %s", tt.expected, m.ReplyTo[0])
				}
		})
	}
}