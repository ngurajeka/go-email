package email

import (
	"net/mail"

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

	ok, message := m.MailMessage()
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

func (m *Message) SendNewsLetter(
	a Account, targets []mail.Address,
) (bool, []error) {

	var (
		errs []error
	)

	if ok, errs := a.Authenticate(); !ok {
		return ok, errs
	}

	auth := a.Auth.(*gomail.Dialer)
	dial, err := auth.Dial()
	if err != nil {
		errs = append(errs, err)
		return false, errs
	}

	for _, target := range targets {
		m.SetTo(target)
		ok, message := m.MailMessage()
		if !ok {
			continue
		}

		if err := gomail.Send(dial, message); err != nil {
			errs = append(errs, err)
			continue
		}
	}

	if len(errs) > 0 {
		return false, errs
	}

	return true, errs
}
