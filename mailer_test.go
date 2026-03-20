// +build unit

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
		{"email with name", "helmi@xeno.com.my", "Helmi Aziz", "Helmi Aziz <helmi@xeno.com.my>"},
		{"email only", "helmi@xeno.com.my", "", "helmi@xeno.com.my"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatAddress(tt.email, tt.fullname) //set
			if got != tt.expected {                     //check
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
		expected string
	}{
		{"email with name", "helmi@xeno.com.my", "Helmi Aziz", "Helmi Aziz <helmi@xeno.com.my>"},
		{"email only", "helmi@xeno.com.my", "", "helmi@xeno.com.my"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Mailer{}
			m.SetFrom(tt.email, tt.fullname) //set

			//Check
			if m.From != tt.expected {
				t.Errorf("expected From %s but got %s", tt.expected, m.From)
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
		{"email with name", "helmi@xeno.com.my", "Helmi Aziz", "Helmi Aziz <helmi@xeno.com.my>"},
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
		{"email with name", "helmi@xeno.com.my", "Helmi Aziz", "Helmi Aziz <helmi@xeno.com.my>"},
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
		{"email with name", "helmi@xeno.com.my", "Helmi Aziz", "Helmi Aziz <helmi@xeno.com.my>"},
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
		{"email with name", "helmi@xeno.com.my", "Helmi Aziz", "Helmi Aziz <helmi@xeno.com.my>"},
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

func TestMailer_SetSubject(t *testing.T) {
	tests := []struct {
		name    string
		subject string
	}{
		{"simple subject", "Hello"},
		{"long subject", "Welcome to our service"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Mailer{}
			m.SetSubject(tt.subject) //Set

			//Check
			if m.Subject != tt.subject {
				t.Errorf("expected %s but got %s", tt.subject, m.Subject)
			}
		})
	}
}

func TestMailer_SetBodyAndAltBody(t *testing.T) {
	tests := []struct {
		name string
		body string
		alt  string
	}{
		{"Body with alt", "<div><p>this is body with alt</p></div>", "this is body with alt"},
		{"Body only", "this is body only", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Mailer{}

			//Set
			m.SetBody(tt.body)
			m.SetAltBody(tt.alt)

			//Check
			if m.Body != tt.body {
				t.Errorf("expected body %s but got %s", tt.body, m.Body)
			}

			if m.AltBody != tt.alt {
				t.Errorf("expected alt body %s but got %s", tt.alt, m.AltBody)
			}
		})
	}
}

func TestMailer_IsHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    bool
		expected string
	}{
		{"html", true, "text/html"},
		{"plain", false, "text/plain"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Mailer{}
			m.IsHTML(tt.input) //Set

			//Check
			if m.ContentType != tt.expected {
				t.Errorf("expected %s but got %s", tt.expected, m.ContentType)
			}
		})
	}
}

func TestMailer_SetDebug(t *testing.T) {
	tests := []struct {
		name  string
		level int
	}{
		{"no debug", 0},
		{"error debug", 1},
		{"verbose debug", 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Mailer{}
			m.SetDebug(tt.level) //Set

			//Check
			if m.Debug != tt.level {
				t.Errorf("expected debug %d but got %d", tt.level, m.Debug)
			}
		})
	}
}