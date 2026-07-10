package sesmailer

import (
	"context"
	"fmt"

	// "io"
	"log"
	
	// "strings"
	// "net/mail"
	// "mime"
	// "os"
	// "path/filepath"

	// "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	// "github.com/aws/smithy-go/logging"
	"github.com/elmyrockers/go-mimebuilder"
)

type Mailer struct {
	debug       int
	errorList 	[]error

	client *ses.Client
	builder *mimebuilder.MimeBuilder
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
			debug:       0,
			errorList:   errorList,

			client:      client,
			builder: 	 mimebuilder.New(),
		}
}

func (m *Mailer) SetFrom(email, name string) *Mailer {
	m.builder.SetFrom( email, name )
	return m
}

func (m *Mailer) AddAddress(email, name string) *Mailer {
	m.AddTo( email, name )
	return m
}

func (m *Mailer) AddTo(email, name string) *Mailer {
	m.builder.AddTo( email, name )
	return m
}

func (m *Mailer) AddCC(email, name string) *Mailer {
	m.builder.AddCC( email, name )
	return m
}

func (m *Mailer) AddBCC(email, name string) *Mailer {
	m.builder.AddBCC( email, name )
	return m
}

func (m *Mailer) AddReplyTo(email, name string) *Mailer {
	m.builder.AddReplyTo( email, name )
	return m
}

func (m *Mailer) SetSubject(subject string) *Mailer {
	m.builder.SetSubject( subject )
	return m
}

func (m *Mailer) SetBody(body string) *Mailer {
	m.builder.SetBody( body )
	return m
}

func (m *Mailer) SetAltBody(alt string) *Mailer {
	m.builder.SetAltBody( alt )
	return m
}

func (m *Mailer) AsHTML() *Mailer {
	m.builder.AsHTML()
	return m
}

func (m *Mailer) Embed(name string, data []byte, cid string) *Mailer {
	m.builder.Embed( name, data, cid )
	return m
}

func (m *Mailer) Attach(filename string, data []byte) *Mailer {
	m.builder.Attach( filename, data )
	return m
}

func (m *Mailer) Dump() {
	mime, _ := m.builder.Build()
	defer m.builder.Release( mime )

	fmt.Println( mime.String() ) 
}

func (m *Mailer) SendWithContext( ctx context.Context ) (*ses.SendRawEmailOutput, error) {
	// Check error
		if len(m.errorList)>0 { return _, m.errorList[0] }

	// Get MIME as buffer
		mime, _ := m.builder.Build()
		defer m.builder.Release( mime )

	// Send email
		output, err := m.client.SendRawEmail( ctx, &ses.SendRawEmailInput{
									RawMessage: &types.RawMessage{
										Data: mime.Bytes(),
									},
								})
		if err != nil { return _, err }
	return output, nil
}

func (m *Mailer) Send() (*ses.SendRawEmailOutput, error) {
	return m.SendWithContext( context.Background() )
}