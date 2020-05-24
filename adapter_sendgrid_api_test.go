package email_test

import (
	"github.com/ngurajeka/go-email"
	"github.com/spf13/viper"
	"log"
	"testing"
)

func TestSendingSendgridAPIEmail(t *testing.T) {
	sender := email.Target{
		Name:  viper.GetString("sendgrid.sender.name"),
		Email: viper.GetString("sendgrid.sender.email"),
	}
	emailSvc := email.NewSendgridApiEmailService(viper.GetString("sendgrid.key"), sender)

	receiver := email.Target{
		Name:  viper.GetString("receiver.name"),
		Email: viper.GetString("receiver.email"),
	}
	message := email.Message{
		PlainText: "Hi, this email only visible when email client cannot view html content",
		HTML:      "<h1>Hi there</h1>",
	}
	var attachments []email.Attachment
	attachments = append(attachments, email.Attachment{
		Filepath: "text.txt",
		Filename: "Sample File.txt",
	})

	messageID, err := emailSvc.Send("Sending email using SMTP", receiver, nil, nil, message, &attachments)
	if err != nil {
		t.Error(err)
	}

	log.Println(messageID)
}
