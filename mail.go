package helpers

import (
	"net/mail"
	"net/smtp"

	"github.com/scorredoira/email"
)

type mailSender struct {
	FromName, Login, Psw, Domain, Port string
	Bcc                                []string
}

var Mail = &mailSender{}

func (m *mailSender) Init(fromName, login, psw, domain, port string, bcc []string) {
	m.FromName, m.Login, m.Psw, m.Bcc = fromName, login, psw, bcc
	m.Domain, m.Port = domain, port
}

func (a *mailSender) SendEmail(subject, text, fileName string, to []string) {

	m := email.NewHTMLMessage(subject, text)
	m.From = mail.Address{Name: a.FromName, Address: a.Login}
	m.To = to
	m.Bcc = a.Bcc

	if fileName != "" {
		if err := m.Attach(fileName); err != nil {
			ErrorLogger.Println(err)
		}
	}

	if err := email.Send(a.Domain+":"+a.Port, smtp.PlainAuth("", a.Login, a.Psw, a.Domain), m); err != nil {
		ErrorLogger.Println(err)
	}
}
