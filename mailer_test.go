// +build unit

package sesmailer

import (
	"testing"
	"strings"
	"io"
)

func TestSanitizeHeader(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"normal text", "normal text"},
		{"line\r\nbreak", "linebreak"},
		{"control\x00char", "controlchar"},
		{"  trim spaces  ", "trim spaces"},
		{"utf8✓chars", "utf8✓chars"},
		{strings.Repeat("x", 1000), strings.Repeat("x", 998)}, // truncated
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := sanitizeHeader(tt.input)
			if got != tt.expected {
				t.Errorf("sanitizeHeader(%q) = %q; want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestFormatAddress(t *testing.T) {
		tests := []struct {
			email    string
			name     string
			expected string
		}{
			{"user@example.com", "", "user@example.com"},
			{"USER@Example.COM", "John", "John <user@example.com>"}, 
			{"invalid-email", "John", ""},
			{"alice@example.com", "Alice Smith", "Alice Smith <alice@example.com>"},
			{"bob@example.com", "Jöhn Döe", "=?utf-8?q?J=C3=B6hn_D=C3=B6e?= <bob@example.com>"},
		}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			got := formatAddress(tt.email, tt.name)
			if got != tt.expected {
				t.Errorf("formatAddress(%q, %q) = %q; want %q", tt.email, tt.name, got, tt.expected)
			}
		})
	}
}

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"normal.txt", "normal.txt"},
		{"../unsafe.txt", "unsafe.txt"},
		{"C:\\path\\file.pdf", "file.pdf"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := sanitizeFilename(tt.input)
			if got != tt.expected {
				t.Errorf("sanitizeFilename(%q) = %q; want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestGenerateBoundary(t *testing.T) {
	boundary, err := generateBoundary()
	if err != nil {
		t.Fatalf("generateBoundary() error: %v", err)
	}
	if !strings.HasPrefix(boundary, "NextPartBoundary_") {
		t.Errorf("boundary = %q, expected prefix 'NextPartBoundary_'", boundary)
	}
	if len(boundary) == 0 {
		t.Error("boundary is empty")
	}
}

func TestEncodeBodyQP(t *testing.T) {
	input := "Hello=World\nNew line"
	got := encodeBodyQP(input)
	if !strings.Contains(got, "Hello=3DWorld") {
		t.Errorf("encodeBodyQP(%q) does not contain expected output", input)
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

func TestMailer_AddAttachment(t *testing.T) {
	tests := []struct {
		name     string
		path 	 string
		display  string
		expected string
	}{
		{
			name:     "Cat Image",
			path:     "examples/attachment/cat.webp",
			display:  "my_cat.webp",
			expected: "my_cat.webp",
		},
		{
			name:     "Rabbit with path",
			path:     "examples/attachment/rabbit.jpg",
			display:  "", // Empty name should trigger filepath.Base
			expected: "rabbit.jpg",
		},
		{
			name:     "Dog Image",
			path:     "examples/attachment/dog.jpg",
			display:  "black_dog.jpg", // Empty name should trigger filepath.Base
			expected: "black_dog.jpg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Mailer{}
			m.AddAttachment(tt.path, tt.display)

			// Verify count
				if len(m.Attachments) != 1 {
					t.Fatalf("Expected 1 attachment, got %d", len(m.Attachments))
				}

			// Verify Filename Sanitization
				gotFile := m.Attachments[0].Filename
				if gotFile != tt.expected {
					t.Errorf("Filename = %q; want %q", gotFile, tt.expected)
				}

			// Verify Data is not nil (os.Open worked)
				if m.Attachments[0].Data == nil {
					t.Error("Attachment Data property is nil")
				}

			// Verify 'Data' implements io.Reader
				reader, ok := m.Attachments[0].Data.(io.Reader)
				if !ok {
					t.Errorf("Expected Data to be an io.Reader, but it is not")
				}

			// Can we actually read from it?
				if reader == nil {
					t.Error("Reader is nil, os.Open likely failed inside the method")
				}
		})
	}
}