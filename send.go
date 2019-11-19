package email

import (
	"gopkg.in/gomail.v2"
)

func (m *Message) Send(a *Account) (bool, []error) {
	var (
		message *gomail.Message
		errs    []error
	)

	if ok, errs := a.Authenticate(); !ok {
		return ok, errs
	}

	message, ok, errs := m.MailMessage()
	if !ok {
		return ok, errs
	}

	auth := a.Auth.(*gomail.Dialer)
	if err := auth.DialAndSend(message); err != nil {
		errs = append(errs, err)
		return false, errs
	}

	return true, errs
}
