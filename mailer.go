package sesmailer

import (
	"context"
	"io"
	"log"
	"fmt"
	"strings"
	"net/mail"
	"mime"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/aws/smithy-go/logging"
)

type Attachment struct {
	Filename string
	Data     io.Reader
}
type Mailer struct {
	from        string
	to          []string
	cc          []string
	bcc         []string
	replyTo     []string
	subject     string
	body        string
	altBody     string
	contentType string //"text/plain" or "text/html"
	attachments []Attachment

	debug       int
	errorList 	[]error

	client *ses.Client
}


// New initializes Mailer and automatically creates SES client
func New() *Mailer {
	// Load config
		var errorList []error
		ctx := context.Background()
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			err = fmt.Errorf("SESMailer: Unable to load AWS config: %w", err)
			errorList = append( errorList, err )
		}

	// Create SES client
		client := ses.NewFromConfig(cfg)
		return &Mailer{
			to:          []string{},
			cc:          []string{},
			bcc:         []string{},
			replyTo:     []string{},
			contentType: "text/plain",
			attachments: []Attachment{},

			debug:       0,
			errorList:   errorList,

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

func (m *Mailer) formatAddress(email, name string) string {
	email = strings.ToLower(strings.TrimSpace(email))
	email = sanitizeHeader(email)

	// Validate the email first
	_, err := mail.ParseAddress(email)
	if err != nil {
		err = fmt.Errorf("SESMailer: invalid email address format: %w", err)
		m.errorList = append( m.errorList, err )
		return ""
	}
	if name == "" { return email }

	name = sanitizeHeader(name)
	name = mime.QEncoding.Encode("utf-8", name)

	return name + " <" + email + ">"
}

func (m *Mailer) encodeBodyQP(body string) string {
	var buf bytes.Buffer
	w := quotedprintable.NewWriter(&buf)
	_, err := w.Write([]byte(body))
	if err != nil {
		err = fmt.Errorf("SESMailer: failed to encode with quoted printable: %w", err)
		m.errorList = append( m.errorList, err )
	}
	w.Close()
	return buf.String()
}

// Sanitize filenames
func sanitizeFilename(name string) string {
	name = sanitizeHeader(name)
	name = filepath.Base(name)
	name = mime.BEncoding.Encode("utf-8", name)
	return name
}
//------------------------------------------------------------------------------------------------------------

func (m *Mailer) SetFrom(email, name string) *Mailer {
	address := m.formatAddress(email, name)
	if address == "" { return m }

	m.from = address
	return m
}

func (m *Mailer) AddAddress(email, name string) *Mailer {
	address := m.formatAddress(email, name)
	if address == "" { return m }

	m.to = append(m.to, address)
	return m
}

func (m *Mailer) AddCC(email, name string) *Mailer {
	address := m.formatAddress(email, name)
	if address == "" { return m }

	m.cc = append(m.cc, address)
	return m
}

func (m *Mailer) AddBCC(email, name string) *Mailer {
	address := m.formatAddress(email, name)
	if address == "" { return m }

	m.bcc = append(m.bcc, address)
	return m
}

func (m *Mailer) AddReplyTo(email, name string) *Mailer {
	address := m.formatAddress(email, name)
	if address == "" { return m }

	m.replyTo = append(m.replyTo, address)
	return m
}

func (m *Mailer) SetSubject(subject string) *Mailer {
	subject = sanitizeHeader(subject)
	m.subject = mime.QEncoding.Encode("utf-8",subject)
	return m
}

func (m *Mailer) SetBody(body string) *Mailer {
	m.body = m.encodeBodyQP(body)
	return m
}

func (m *Mailer) SetAltBody(alt string) *Mailer {
	m.altBody = m.encodeBodyQP(alt)
	return m
}

func (m *Mailer) IsHTML(isHtml bool) *Mailer {
	if isHtml {
		m.contentType = "text/html"
	} else {
		m.contentType = "text/plain"
	}
	return m
}

func (m *Mailer) AddAttachment(path string, name string) *Mailer {
	// Get file pointer
		file, err := os.Open( path )
		if err != nil {
			err = fmt.Errorf("SESMailer - failed to open attachment: %w", err)
			m.errorList = append( m.errorList, err )
			return m
		}

	// Set filename
		if name == "" { name = filepath.Base(path) }
		
	// Sanitize filename
		name = sanitizeFilename(name)
	m.attachments = append(m.attachments, Attachment{
		Filename: name,
		Data:     file, // stream data
	})

	return m
}

// SendContext sends the email using AWS SES
func (m *Mailer) SendContext(ctx context.Context) (*ses.SendEmailOutput, error) {
	// Validate inputs
		var err error
		if m.from == "" {
			err = fmt.Errorf("SESMailer: No 'from' address specified")
			m.errorList = append( m.errorList, err )
			return nil, err
		}
		if len(m.to)+len(m.cc)+len(m.bcc) == 0 {
			err = fmt.Errorf("SESMailer: No recipients specified")
			m.errorList = append( m.errorList, err )
			return nil, err
		}

	// Prepare destination (to, cc and bcc)
		destination := &types.Destination{
			ToAddresses:  m.to,
			CcAddresses:  m.cc,
			BccAddresses: m.bcc,
		}

	// Prepare content (subject, body and altbody)
		var body *types.Body

		// HTML body (for future use)
			if m.contentType == "text/html" {
				body = &types.Body{
					Html: &types.Content{
						Data:    &m.body,
						Charset: aws.String("UTF-8"),
					},
				}

			// Plain text fallback
				if m.altBody != "" {
					body.Text = &types.Content{
						Data:    &m.AltBody,
						Charset: aws.String("UTF-8"),
					}
				}
		// Plain text only
			} else {
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