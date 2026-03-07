package sesmailer

import (
	"context"
	"log"
	"fmt"
	"os"
	
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/aws/smithy-go/logging"
)

type Mail struct {
    From     string
    FromName string
    To       []string
    Cc       []string
    Bcc      []string
    ReplyTo  []string
    Subject  string
    Body     string
    AltBody  string
    ContentType string  //"text/plain" or "text/html"
    Debug    int

    client *ses.Client
}



// New initializes Mail and automatically creates SES client
func New() *Mail {
	// Load config
		ctx := context.Background()
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			log.Fatalf("unable to load AWS config: %v", err)
		}

	// Create SES client
		client := ses.NewFromConfig(cfg)
		return &Mail{
			To:      []string{},
			Cc:      []string{},
			Bcc:     []string{},
			ReplyTo: []string{},
			ContentType: "text/plain",
			Debug:   0,
			client:  client,
		}
}

func (m *Mail) SetFrom(email, name string) *Mail {
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

func (m *Mail) AddAddress(email, name string) *Mail {
	m.To = append(m.To, formatAddress(email, name))
	return m
}

func (m *Mail) AddCC(email, name string) *Mail {
	m.Cc = append(m.Cc, formatAddress(email, name))
	return m
}

func (m *Mail) AddBCC(email, name string) *Mail {
	m.Bcc = append(m.Bcc, formatAddress(email, name))
	return m
}

func (m *Mail) AddReplyTo(email, name string) *Mail {
	m.ReplyTo = append(m.ReplyTo, formatAddress(email, name))
	return m
}

func (m *Mail) SetSubject(subject string) *Mail {
    m.Subject = subject
    return m
}

func (m *Mail) SetBody(body string) *Mail {
    m.Body = body
    return m
}

func (m *Mail) SetAltBody(alt string) *Mail {
    m.AltBody = alt
    return m
}







// Set debug level: 0 = none, 1 = errors only, 2 = verbose
// SetDebug sets the debug level
func (m *Mail) SetDebug(level int) *Mail {
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







func (m *Mail) IsHTML(isHtml bool) *Mail {
    if isHtml {
        m.ContentType = "text/html"
    } else {
        m.ContentType = "text/plain"
    }
    return m
}



// Send sends the email using AWS SES
func (m *Mail) Send(ctx context.Context) error {
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
		if m.FromName != "" {
			from = fmt.Sprintf("%s <%s>", m.FromName, m.From)
		}

		input := &ses.SendEmailInput{
			Source:      &from,
			Destination: destination,
			Message:     message,
		}

	// Verbose logging before sending
		if m.Debug >= 1 {
			log.Println("[DEBUG] Preparing to send email")
			log.Printf("[DEBUG] From: %s\nTo: %v\nCC: %v\nBCC: %v\n", m.From, m.To, m.Cc, m.Bcc)
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