package sender

import (
	"errors"
	"gopkg.in/gomail.v2"
)

type MailSender interface {
	Send(m *gomail.Message) error
}

type Sender struct {
	Email  string
	Dialer *gomail.Dialer
}

func (s *Sender) Send(m *gomail.Message) error {
	if m == nil {
		return errors.New("message cannot be nil")
	}
	if s.Email == "" {
		return errors.New("sender email is empty")
	}
	if s.Dialer == nil {
		return errors.New("dialer is not initialized")
	}

	m.SetHeader("From", s.Email)
	return s.Dialer.DialAndSend(m)
}

func NewSender(email string, dialer *gomail.Dialer) *Sender {
	return &Sender{
		Email:  email,
		Dialer: dialer,
	}
}
