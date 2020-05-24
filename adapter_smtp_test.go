package email_test

import (
	"log"
	"testing"

	"github.com/ngurajeka/go-email"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	_ = viper.ReadInConfig()
}

func TestSendingSMTPEmail(t *testing.T) {
	smtpConfig := email.NewSMTPConfig(
		viper.GetString("smtp.host"),
		viper.GetString("smtp.username"),
		viper.GetString("smtp.password"),
		viper.GetInt("smtp.port"))
	sender := email.Target{
		Name:  viper.GetString("sender.name"),
		Email: viper.GetString("sender.email"),
	}
	emailSvc := email.NewSMTPEmailService(smtpConfig, sender)

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
