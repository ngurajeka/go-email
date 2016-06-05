package email

import (
	"log"

	"gopkg.in/gomail.v2"
)

func (m *Message) MailMessage() (bool, *gomail.Message) {

	message := gomail.NewMessage()

	if ok, errs := m.Check(); !ok {
		for _, err := range errs {
			log.Println(err)
		}
		return false, message
	}

	message.SetHeader("From", m.From.String())
	message.SetHeader("To", m.GetTo()...)
	if len(m.Cc) > 0 {
		for _, cc := range m.Cc {
			message.SetAddressHeader("Cc", cc.Address, cc.Name)
		}
	}
	if len(m.Bcc) > 0 {
		for _, bcc := range m.Bcc {
			message.SetAddressHeader("Cc", bcc.Address, bcc.Name)
		}
	}
	message.SetHeader("Subject", m.Subject)
	plain := ParsingBody(m.Body, m.Params)
	html := ParsingBody(m.HTMLBody, m.Params)
	message.SetBody(MIME_PLAIN, plain)
	message.AddAlternative(MIME_HTML, html)
	if len(m.Attachments) > 0 {
		for _, attach := range m.Attachments {
			message.Attach(attach.Path)
		}
	}

	m.Source = map[string]string{"plain": plain, "html": html}

	return true, message
}
