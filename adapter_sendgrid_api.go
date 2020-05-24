package email

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"io/ioutil"
	"net/http"
	"os"
)

type sendgridErrorResponse struct {
	Errors []struct {
		Message string `json:"message"`
		Field   string `json:"field"`
		Help    string `json:"help"`
	} `json:"errors"`
}

type sendgridApiEmailService struct {
	client *sendgrid.Client
	sender Target
}

func NewSendgridApiEmailService(apiKey string, sender Target) Service {
	return &sendgridApiEmailService{
		client: sendgrid.NewSendClient(apiKey),
		sender: sender,
	}
}

func (svc *sendgridApiEmailService) Send(subject string, target Target, cc, bcc *[]Target, message Message, attachments *[]Attachment) (string, error) {
	mailMessage := mail.NewV3Mail()
	mailMessage.SetFrom(mail.NewEmail(svc.sender.Name, svc.sender.Email))
	mailMessage.AddContent(mail.NewContent("text/plain", message.PlainText))
	mailMessage.AddContent(mail.NewContent("text/html", message.HTML))

	personalization := mail.NewPersonalization()
	personalization.Subject = subject
	personalization.AddTos(mail.NewEmail(target.Name, target.Email))
	if cc != nil {
		personalization.AddCCs(svc.convertTargetsToEmails(*cc)...)
	}
	if bcc != nil {
		personalization.AddBCCs(svc.convertTargetsToEmails(*bcc)...)
	}
	if attachments != nil {
		mailMessage.AddAttachment(svc.convertAttachments(*attachments)...)
	}

	mailMessage.AddPersonalizations(personalization)

	response, err := svc.client.Send(mailMessage)
	if err != nil {
		return "", err
	}

	if response.StatusCode == http.StatusUnauthorized ||
		response.StatusCode == http.StatusForbidden {
		var errorResponse sendgridErrorResponse
		if err := json.Unmarshal([]byte(response.Body), &errorResponse); err != nil {
			return "", err
		}

		return "", errors.New(errorResponse.Errors[0].Message)
	}

	messageIDs := response.Headers["X-Message-Id"]
	if len(messageIDs) > 0 {
		return messageIDs[0], nil
	}

	return "", nil
}

func (svc *sendgridApiEmailService) Status(messageID string) string {
	panic("not implemented")
}

func (svc *sendgridApiEmailService) convertTargetsToEmails(targets []Target) []*mail.Email {
	var emails []*mail.Email
	for _, target := range targets {
		emails = append(emails, mail.NewEmail(target.Name, target.Email))
	}
	return emails
}

func (svc *sendgridApiEmailService) convertAttachments(attachments []Attachment) []*mail.Attachment {
	var mailAttachments []*mail.Attachment
	for _, attachment := range attachments {
		mailAttachment := mail.NewAttachment()
		dat, err := ioutil.ReadFile(attachment.Filepath)
		if err != nil {
			continue
		}
		f, err := os.Open(attachment.Filepath)
		if err != nil {
			continue
		}
		contentType, err := GetFileContentType(f)
		if err != nil {
			continue
		}
		f.Close()
		encoded := base64.StdEncoding.EncodeToString(dat)
		mailAttachment.SetContent(encoded)
		mailAttachment.SetType(contentType)
		mailAttachment.SetFilename(attachment.Filename)
		mailAttachment.SetDisposition("attachment")
		mailAttachment.SetContentID(attachment.Filename)
		mailAttachments = append(mailAttachments, mailAttachment)
	}

	return mailAttachments
}
