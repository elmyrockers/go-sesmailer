package sesmailer

import (
	"github.com/aws/aws-sdk-go-v2/service/ses"
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
    Debug    int

    client *ses.Client
}



func New() *Mail {
    return &Mail{
        To:      []string{},
        Cc:      []string{},
        Bcc:     []string{},
        ReplyTo: []string{},
        Debug:   0, // 0 = none
    }
}