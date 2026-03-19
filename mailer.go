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

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/aws/smithy-go/logging"

	// "github.com/davecgh/go-spew/spew"
)


type Attachment struct {
	Filename string
	Data     []byte
}
type Mailer struct {
	From        string
	FromName    string
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

func (m *Mailer) SetFrom(email, name string) *Mailer {
	m.From = email
	m.FromName = name
	return m
}

// helper to format "Name <email>" if name is given
func formatAddress(email, name string) string {
	if name != "" {
		return fmt.Sprintf("\"%s\" <%s>", name, email)
	}
	return email
}

func (m *Mailer) AddAddress(email, name string) *Mailer {
	m.To = append(m.To, formatAddress(email, name))
	return m
}

func (m *Mailer) AddCC(email, name string) *Mailer {
	m.Cc = append(m.Cc, formatAddress(email, name))
	return m
}

func (m *Mailer) AddBCC(email, name string) *Mailer {
	m.Bcc = append(m.Bcc, formatAddress(email, name))
	return m
}

func (m *Mailer) AddReplyTo(email, name string) *Mailer {
	m.ReplyTo = append(m.ReplyTo, formatAddress(email, name))
	return m
}

func (m *Mailer) SetSubject(subject string) *Mailer {
	m.Subject = subject
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
	// Get binary data
		data, err := os.ReadFile(path)
		if err != nil {
			log.Printf("failed to read attachment: %v", err)
			return m
		}

	// Set filename
		if name == "" { name = filepath.Base(path) }

	m.Attachments = append(m.Attachments, Attachment{
		Filename: name,
		Data:     data,
	})

	return m
}

func (m *Mailer) SendRaw(ctx context.Context) error {
	boundary := fmt.Sprintf("NextPartBoundary_%d", time.Now().UnixNano())
	var buf bytes.Buffer

	// Headers 
		// from, to, cc and reply-to
			from := formatAddress(m.From, m.FromName)
			buf.WriteString(fmt.Sprintf("From: %s\r\n", from))
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
			buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\r\n\r\n", boundary ))

	// Main body part
		buf.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		buf.WriteString(fmt.Sprintf("Content-Type: %s; charset=\"UTF-8\"\r\n", m.ContentType))
		buf.WriteString("Content-Transfer-Encoding: 7bit\r\n\r\n")
		buf.WriteString( m.Body + "\r\n" )

		// Attachments
			for _, att := range m.Attachments {
				buf.WriteString(fmt.Sprintf("--%s\r\n", boundary))
				buf.WriteString(fmt.Sprintf("Content-Type: application/octet-stream; name=\"%s\"\r\n", att.Filename))
				buf.WriteString("Content-Transfer-Encoding: base64\r\n")
				buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n\r\n", att.Filename))

					encoded := make([]byte, base64.StdEncoding.EncodedLen(len(att.Data)))
					base64.StdEncoding.Encode(encoded, att.Data)
				buf.WriteString(string(encoded) + "\r\n")
			}
		buf.WriteString(fmt.Sprintf("--%s--", boundary))


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

	_, err := m.client.SendRawEmail(ctx, input)
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
	from := formatAddress(m.From, m.FromName)

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