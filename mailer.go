package sesmailer

import (
	"context"
	"fmt"

	// "io"
	// "log"
	
	// "strings"
	// "net/mail"
	// "mime"
	// "os"
	// "path/filepath"

	// "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	// "github.com/aws/aws-sdk-go-v2/service/ses/types"
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