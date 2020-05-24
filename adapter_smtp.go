package email

import (
	"github.com/google/uuid"
	"gopkg.in/gomail.v2"
)

type SMTPConfig struct {
	Host, Username, Password string
	Port                     int
}

func NewSMTPConfig(host, username, password string, port int) SMTPConfig {
	return SMTPConfig{
		Host:     host,
		Username: username,
		Password: password,
		Port:     port,
	}
}

type smtpEmailService struct {
	dialer *gomail.Dialer
	sender Target
}

func NewSMTPEmailService(config SMTPConfig, sender Target) Service {
	return &smtpEmailService{
		dialer: gomail.NewDialer(config.Host, config.Port, config.Username, config.Password),
		sender: sender,
	}
}

func (svc *smtpEmailService) Send(subject string, target Target, cc, bcc *[]Target, message Message, attachments *[]Attachment) (string, error) {
	messageID := uuid.New().String()
	mailMessage := gomail.NewMessage()
	mailMessage.SetHeader("Subject", subject)
	mailMessage.SetHeader("From", svc.sender.String())
	mailMessage.SetHeader("To", target.String())
	if cc != nil {
		for _, v := range *cc {
			mailMessage.SetAddressHeader("Cc", v.Email, v.Name)
		}
	}
	if bcc != nil {
		for _, v := range *bcc {
			mailMessage.SetAddressHeader("Bcc", v.Email, v.Name)
		}
	}

	mailMessage.SetBody("text/plain", message.PlainText)
	mailMessage.SetBody("text/html", message.HTML)

	if attachments != nil {
		for _, attachment := range *attachments {
			mailMessage.Attach(attachment.Filepath, gomail.Rename(attachment.Filename))
		}
	}

	if err := svc.dialer.DialAndSend(mailMessage); err != nil {
		return "", err
	}

	return messageID, nil
}

func (svc *smtpEmailService) Status(messageID string) string {
	panic("not implemented")
}
