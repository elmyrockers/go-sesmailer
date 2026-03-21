package sesmailer

import (
	"context"
	"fmt"
	"log"
	"os"
	"encoding/base64"
	"strings"
	"bytes"
	"path/filepath"
	"time"
	"crypto/rand"
	"net/mail"
	"mime/quotedprintable"
	"mime"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/aws/smithy-go/logging"
)


type Attachment struct {
	Filename string
	Data     io.Reader
	// Data     []byte
}
type Mailer struct {
	From        string
	To          []string
	Cc          []string
	Bcc         []string
	ReplyTo     []string
	Subject     string
	Body        string
	AltBody     string
	ContentType string //"text/plain" or "text/html"
	Attachments []Attachment
	Debug       int

	client *ses.Client
}

// New initializes Mailer and automatically creates SES client
func New() *Mailer {
	// Load config
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}

	// Create SES client
	client := ses.NewFromConfig(cfg)
	return &Mailer{
		To:          []string{},
		Cc:          []string{},
		Bcc:         []string{},
		ReplyTo:     []string{},
		ContentType: "text/plain",
		Debug:       0,
		client:      client,
	}
}


//---------------------------------------------------------------------------------------------------- HELPERS
func sanitizeHeader(s string) string {
	// Clean control characters (Rune-aware)
		s = strings.Map(func(r rune) rune {
			if r == '\r' || r == '\n' || r == 0 || r == 127 {
				return -1
			}
			if r < 32 && r != '\t' {
				return -1
			}
			return r
		}, s)
		s = strings.TrimSpace(s)

	// Truncate by BYTES (RFC 5322) but stay UTF-8 safe
		if len(s) > 998 {
			idx := 998
			// Step back if inside a UTF-8 continuation byte
			for idx > 0 && (s[idx]&0xC0 == 0x80) {
				idx--
			}
			s = s[:idx]
		}

	return s
}

func formatAddress(email, name string) string {
    email = strings.ToLower(strings.TrimSpace(email))
    email = sanitizeHeader(email)

    // Validate the email first
    if _, err := mail.ParseAddress(email); err != nil {
        return "" 
    }

    if name == "" {
        return email
    }

    name = sanitizeHeader(name)
    name = mime.QEncoding.Encode("utf-8", name)

    return fmt.Sprintf("%s <%s>", name, email)
}

// Sanitize filenames
func sanitizeFilename(name string) string {
	name = sanitizeHeader(name)
	return filepath.Base(name)
}

func generateBoundary() (string, error) {
	b := make([]byte, 8) // 8 random bytes
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	boundary := fmt.Sprintf("NextPartBoundary_%d_%x", time.Now().UnixNano(), b)
	return boundary, nil
}

func encodeBodyQP(body string) string {
	var buf bytes.Buffer
	w := quotedprintable.NewWriter(&buf)
	_, err := w.Write([]byte(body))
	if err != nil {
		log.Printf("encodeBodyQP: failed to write body: %v", err)
	}
	w.Close()
	return buf.String()
}
//----------------------------------------------------------------------------------------------------

func (m *Mailer) SetFrom(email, name string) *Mailer {
	address := formatAddress(email, name)
	if address == "" { return m }

	m.From = address
	return m
}

func (m *Mailer) AddAddress(email, name string) *Mailer {
	address := formatAddress(email, name)
	if address == "" { return m }

	m.To = append(m.To, address)
	return m
}

func (m *Mailer) AddCC(email, name string) *Mailer {
	address := formatAddress(email, name)
	if address == "" { return m }

	m.Cc = append(m.Cc, address)
	return m
}

func (m *Mailer) AddBCC(email, name string) *Mailer {
	address := formatAddress(email, name)
	if address == "" { return m }

	m.Bcc = append(m.Bcc, address)
	return m
}

func (m *Mailer) AddReplyTo(email, name string) *Mailer {
	address := formatAddress(email, name)
	if address == "" { return m }

	m.ReplyTo = append(m.ReplyTo, address)
	return m
}


func (m *Mailer) SetSubject(subject string) *Mailer {
	subject = sanitizeHeader(subject)
	m.Subject = mime.QEncoding.Encode("utf-8",subject)
	return m
}

func (m *Mailer) SetBody(body string) *Mailer {
	m.Body = body
	return m
}

func (m *Mailer) SetAltBody(alt string) *Mailer {
	m.AltBody = alt
	return m
}

// Set debug level: 0 = none, 1 = errors only, 2 = verbose
// SetDebug sets the debug level
func (m *Mailer) SetDebug(level int) *Mailer {
	m.Debug = level

	if level == 0 {
		return m
	}

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}

	// Use smithy logger (required by AWS SDK v2)
	cfg.Logger = logging.NewStandardLogger(os.Stderr)

	if level == 1 {
		// Minimal logging (retries/errors)
		cfg.ClientLogMode = aws.LogRetries
	}

	if level >= 2 {
		// Verbose logging (request + response)
		cfg.ClientLogMode =
			aws.LogRetries |
				aws.LogRequest |
				aws.LogResponse |
				aws.LogRequestWithBody |
				aws.LogResponseWithBody
	}

	// Recreate client with logging enabled
	m.client = ses.NewFromConfig(cfg)

	return m
}

func (m *Mailer) IsHTML(isHtml bool) *Mailer {
	if isHtml {
		m.ContentType = "text/html"
	} else {
		m.ContentType = "text/plain"
	}
	return m
}

func (m *Mailer) AddAttachment(path string, name string) *Mailer {
	// Get file pointer
		file, err := os.Open( path )
		if err != nil {
			log.Printf("failed to open attachment: %v", err)
			return m
		}

	// Set filename
		if name == "" { name = filepath.Base(path) }
		
	// Sanitize it first
		name = sanitizeFilename(name)
		name = mime.BEncoding.Encode("utf-8", name)
	m.Attachments = append(m.Attachments, Attachment{
		Filename: name,
		Data:     file, // stream data
	})

	return m
}

func (m *Mailer) SendRaw(ctx context.Context) error {
	// Automatically close the opened attachments
		defer func() {
			for _, att := range m.Attachments {
				if closer, ok := att.Data.(io.Closer); ok {
					closer.Close()
				}
			}
		}()

	// Validate inputs first
		if m.From == "" {
			return fmt.Errorf("Invalid 'from' address")
		}
		if len(m.To)+len(m.Cc)+len(m.Bcc) == 0 {
			return fmt.Errorf("No recipients")
		}

	// Generate mixed boundary (because there are attachments)
		mixedBoundary, err := generateBoundary()
		if err != nil {
			return fmt.Errorf("Failed to generate mixed boundary: %w", err)
		}

	// Headers
		var buf bytes.Buffer

		// from, to, cc and reply-to
			buf.WriteString(fmt.Sprintf("From: %s\r\n", m.From))
			buf.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(m.To, ",")))
			if len(m.Cc) > 0 {
				buf.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(m.Cc, ",")))
			}
			if len(m.ReplyTo) > 0 {
				buf.WriteString(fmt.Sprintf("Reply-To: %s\r\n", strings.Join(m.ReplyTo, ",")))
			}
		// subject
			buf.WriteString(fmt.Sprintf("Subject: %s\r\n", m.Subject))
			buf.WriteString("MIME-Version: 1.0\r\n")
			buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\r\n\r\n", mixedBoundary ))

	// Main body part
		buf.WriteString(fmt.Sprintf("--%s\r\n", mixedBoundary))
		
		// Plaintext-only or Html-only
			if m.ContentType == "text/plain" || len(m.AltBody) == 0 {
				buf.WriteString(fmt.Sprintf("Content-Type: %s; charset=\"UTF-8\"\r\n", m.ContentType))
				buf.WriteString("Content-Transfer-Encoding: quoted-printable\r\n\r\n")
				buf.WriteString( encodeBodyQP(m.Body) + "\r\n" )

		// Html with Altbody
			} else {
				altBoundary, err := generateBoundary()
				if err != nil {
					return fmt.Errorf("Failed to generate alt boundary: %w", err)
				}
				buf.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=\"%s\"\r\n\r\n", altBoundary))

				// Plaintext
					buf.WriteString(fmt.Sprintf("--%s\r\n", altBoundary))
					buf.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
					buf.WriteString("Content-Transfer-Encoding: quoted-printable\r\n\r\n")
					buf.WriteString( encodeBodyQP(m.AltBody) + "\r\n" )

				// HTML
					buf.WriteString(fmt.Sprintf("--%s\r\n", altBoundary))
					buf.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
					buf.WriteString("Content-Transfer-Encoding: quoted-printable\r\n\r\n")
					buf.WriteString( encodeBodyQP(m.Body) + "\r\n" )

				buf.WriteString(fmt.Sprintf("--%s--\r\n", altBoundary)) // close alt boundary
			}

	// Attachments
		for _, att := range m.Attachments {
			extension := filepath.Ext( att.Filename )
			contentType := mime.TypeByExtension( extension )
			if contentType == "" { contentType = "application/octet-stream" }

			buf.WriteString(fmt.Sprintf("--%s\r\n", mixedBoundary))
			buf.WriteString(fmt.Sprintf("Content-Type: %s; name=\"%s\"\r\n", contentType, att.Filename))
			buf.WriteString("Content-Transfer-Encoding: base64\r\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n\r\n", att.Filename))

			// High Performance: Pipe the data from Reader -> Base64 Encoder -> Final Buffer
				encoder := base64.NewEncoder(base64.StdEncoding, &buf)
				_, err := io.Copy(encoder, att.Data) // Streams in 32KB chunks
				if err != nil {
					return fmt.Errorf("Failed to stream attachment %s: %w", att.Filename, err)
				}
				err = encoder.Close() // Flush the encoder
				if err != nil {
					return fmt.Errorf("Failed to finalize attachment %s: %w", att.Filename, err)
				}
				buf.WriteString("\r\n")
		}
		buf.WriteString(fmt.Sprintf("--%s--\r\n", mixedBoundary)) // close mixed boundary


	// Prepare input for SendRawEmail()
		allRecipients := append([]string{}, m.To...)
		allRecipients = append(allRecipients, m.Cc...)
		allRecipients = append(allRecipients, m.Bcc...)
		input := &ses.SendRawEmailInput{
			RawMessage: &types.RawMessage{
				Data: buf.Bytes(),
			},
			Destinations: allRecipients,
		}

	_, err = m.client.SendRawEmail(ctx, input)
	if err != nil {
		if m.Debug > 0 {
			fmt.Printf("\n\nSES SendRawEmail error: %v", err)
		}
		return err
	}

	if m.Debug > 0 {
		fmt.Println("\n\nEmail sent successfully")
	}


	return err
}

// SendContext sends the email using AWS SES
func (m *Mailer) SendContext(ctx context.Context) error {
	// If there is an attachment, then send with raw email input
		if len(m.Attachments) > 0 { return m.SendRaw( ctx ) }

	// Validate inputs
		if m.From == "" {
			return fmt.Errorf("no 'From' address specified")
		}
		if len(m.To)+len(m.Cc)+len(m.Bcc) == 0 {
			return fmt.Errorf("no recipients specified")
		}

	// Prepare destination
	destination := &types.Destination{
		ToAddresses:  m.To,
		CcAddresses:  m.Cc,
		BccAddresses: m.Bcc,
	}

	// Prepare message body
	var body *types.Body
	if m.ContentType == "text/html" {
		// HTML body (for future use)
		body = &types.Body{
			Html: &types.Content{
				Data:    &m.Body,
				Charset: aws.String("UTF-8"),
			},
		}
		if m.AltBody != "" {
			// optional plain text fallback
			body.Text = &types.Content{
				Data:    &m.AltBody,
				Charset: aws.String("UTF-8"),
			}
		}
	} else {
		// Plain text only
		body = &types.Body{
			Text: &types.Content{
				Data:    &m.Body,
				Charset: aws.String("UTF-8"),
			},
		}
	}

	// Prepare message
	message := &types.Message{
		Subject: &types.Content{
			Data:    &m.Subject,
			Charset: aws.String("UTF-8"),
		},
		Body: body,
	}

	// Prepare SES input
	from := m.From

	input := &ses.SendEmailInput{
		Source:           &from,
		Destination:      destination, // To/Cc/Bcc
		Message:          message,
		ReplyToAddresses: m.ReplyTo,
	}

	// Verbose logging before sending
	if m.Debug >= 1 {
		log.Println("[DEBUG] Preparing to send email")
		log.Printf("[DEBUG] From: %s\nReply-To: %v\nTo: %v\nCC: %v\nBCC: %v\n", m.From, m.ReplyTo, m.To, m.Cc, m.Bcc)
		log.Printf("[DEBUG] Subject: %s\nBody: %s\nAltBody: %s\nContentType: %s\n\n", m.Subject, m.Body, m.AltBody, m.ContentType)
	}

	// Send email
	_, err := m.client.SendEmail(ctx, input)
	if err != nil {
		if m.Debug > 0 {
			fmt.Printf("\n\nSES SendEmail error: %v", err)
		}
		return err
	}

	if m.Debug > 0 {
		fmt.Println("\n\nEmail sent successfully")
	}

	return nil
}

func (m *Mailer) Send() error {
	ctx := context.Background()
	return m.SendContext(ctx)
}