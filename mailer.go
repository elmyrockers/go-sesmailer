package sesmailer

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/aws/smithy-go/logging"
)

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

// SendContext sends the email using AWS SES
func (m *Mailer) SendContext(ctx context.Context) error {
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