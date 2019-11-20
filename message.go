package email

import (
	"gopkg.in/gomail.v2"
)

func (m *Message) MailMessage() (*gomail.Message, bool, []error) {
	message := gomail.NewMessage()

	if ok, errs := m.Validate(); !ok {
		return message, false, errs
	}

	message.SetHeader("Subject", m.Subject)
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

	if m.Body != nil {
		plain, _ := ParsingBody(m.Body, m.Params)
		message.SetBody(MimePlain, plain)
	}
	if m.HTMLBody != nil {
		html, _ := ParsingBody(m.HTMLBody, m.Params)
		message.AddAlternative(MimeHtml, html)
	}

	if len(m.Attachments) > 0 {
		for _, attach := range m.Attachments {
			message.Attach(attach.Path)
		}
	}

	return message, true, nil
}
