package sender

import (
	"errors"
	assert2 "github.com/stretchr/testify/assert"
	"gopkg.in/gomail.v2"
	"testing"
)

type MockSender struct {
	ShouldFail bool
}

func (m *MockSender) Send(msg *gomail.Message) error {

	if m.ShouldFail {
		return errors.New("failed to send email")
	}
	return nil
}

func TestSender_Send_Success(t *testing.T) {
	assert := assert2.New(t)

	t.Run("success send", func(t *testing.T) {
		mock := MockSender{ShouldFail: false}
		mail := gomail.NewMessage()
		mail.SetHeader("To", "test@example.com")
		mail.SetHeader("Subject", "Test Email")
		mail.SetBody("text/plain", "This is a test email.")

		err := mock.Send(mail)

		assert.NoError(err, "Expected email to be sent successfully")
	})

	t.Run("fail send", func(t *testing.T) {
		mock := MockSender{ShouldFail: true}
		mail := gomail.NewMessage()
		mail.SetHeader("To", "test@example.com")
		mail.SetHeader("Subject", "Test Email")
		mail.SetBody("text/plain", "This is a test email.")

		err := mock.Send(mail)

		assert.Error(err, "Expected email send to fail")
		assert.Equal("failed to send email", err.Error(), "Expected error message to match")
	})

}

func TestNewSender(t *testing.T) {
	const (
		port = 587
		host = "smtp.example.com"
		user = "test@example.com"
		pass = "password"
	)
	dialer := gomail.NewDialer(host, port, user, pass)
	assert := assert2.New(t)

	t.Run("passed empty email should be error", func(t *testing.T) {
		sender := NewSender("", dialer)
		mail := gomail.NewMessage()
		mail.SetHeader("To", "test@example.com")
		mail.SetHeader("Subject", "Test Email")
		mail.SetBody("text/plain", "This is a test email.")

		err := sender.Send(mail)

		assert.Error(err, "Expected an error due to missing sender email")
		assert.Equal("sender email is empty", err.Error(), "Expected specific error message")
	})

	t.Run("check email", func(t *testing.T) {
		sender := NewSender("test@example.com", dialer)
		assert.Equal(sender.Email, "test@example.com", "Expected email to match")
	})

	t.Run("check fields", func(t *testing.T) {
		assert.Equal(dialer.Host, host, "Expected match")
		assert.Equal(dialer.Port, port, "Expected match")
		assert.Equal(dialer.Username, user, "Expected match")
		assert.Equal(dialer.Password, pass, "Expected match")
	})

}
