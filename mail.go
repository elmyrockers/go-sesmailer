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

func (m *Mail) SetFrom(email, name string) *Mail {
    m.From = email
    m.FromName = name
    return m
}

func (m *Mail) AddAddress(email string) *Mail {
    m.To = append(m.To, email)
    return m
}

func (m *Mail) AddCC(email string) *Mail {
    m.Cc = append(m.Cc, email)
    return m
}

func (m *Mail) AddBCC(email string) *Mail {
    m.Bcc = append(m.Bcc, email)
    return m
}

func (m *Mail) AddReplyTo(email string) *Mail {
    m.ReplyTo = append(m.ReplyTo, email)
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
func (m *Mail) SetDebug(level int) *Mail {
    m.Debug = level
    return m
}