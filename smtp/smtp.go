package smtp

import (
	"fmt"
	"net/smtp"
	"strings"
)

const (
	messageTemplate = "To: %s\r\nFrom: %s<%s>\r\nSubject: %s\r\n%s\r\n\r\n%s"
	contextType     = "Content-Type: text/plain; charset=UTF-8"
)

type Mail struct {
	Subject     string
	ContentType string
	Detail      string
}

type SMTPCfg struct {
	Enable   bool     `env:"SMTP_ENABLE"`
	Addr     string   `env:"SMTP_ADDR"`
	MailFrom string   `env:"SMTP_MAILFROM"`
	MailTo   []string `env:"SMTP_MAILTO"`
	Auth     string   `env:"SMTP_AUTH"`

	host         string
	mailToString string
}

func (s *SMTPCfg) Init() *SMTPCfg {
	s.host = strings.Split(s.Addr, ":")[0]
	s.mailToString = strings.Join(s.MailTo, ";")
	return s
}

func (s *SMTPCfg) SendMail(m Mail) error {
	if !s.Enable {
		return nil
	}
	if len(m.ContentType) == 0 {
		m.ContentType = contextType
	}
	auth := smtp.PlainAuth("", s.MailFrom, s.Auth, s.host)
	msg := fmt.Sprintf(messageTemplate, s.mailToString, s.MailFrom, s.MailFrom, m.Subject, m.ContentType, m.Detail)

	fmt.Println(s.Addr, auth, s.MailFrom, s.MailTo, msg)
	return smtp.SendMail(s.Addr, auth, s.MailFrom, s.MailTo, []byte(msg))
}
