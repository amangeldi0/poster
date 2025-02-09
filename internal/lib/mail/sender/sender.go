package sender

import (
	"gopkg.in/gomail.v2"
)

type Sender struct {
	Email  string
	Dialer *gomail.Dialer
}

func (s Sender) Send(m *gomail.Message) error {
	m.SetHeader("From", s.Email)
	return s.Dialer.DialAndSend(m)
}

func NewSender(email string, dialer *gomail.Dialer) Sender {
	return Sender{
		Email:  email,
		Dialer: dialer,
	}
}
