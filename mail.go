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

func (m *Mail) SetFrom(email, name string) {
    m.From = email
    m.FromName = name
}

func (m *Mail) AddAddress(email string) {
    m.To = append(m.To, email)
}

func (m *Mail) AddCC(email string) {
    m.Cc = append(m.Cc, email)
}

func (m *Mail) AddBCC(email string) {
    m.Bcc = append(m.Bcc, email)
}

func (m *Mail) AddReplyTo(email string) {
    m.ReplyTo = append(m.ReplyTo, email)
}

func (m *Mail) SetSubject(subject string) {
    m.Subject = subject
}

func (m *Mail) SetBody(body string) {
    m.Body = body
}

func (m *Mail) SetAltBody(alt string) {
    m.AltBody = alt
}

// Set debug level: 0 = none, 1 = errors only, 2 = verbose
func (m *Mail) SetDebug(level int) {
    m.Debug = level
}