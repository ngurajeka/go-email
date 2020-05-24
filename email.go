package email

import (
	"fmt"
)

type Target struct {
	Name, Email string
}

func (t Target) String() string {
	if t.Name == "" {
		return fmt.Sprintf("<%s>", t.Email)
	}
	return fmt.Sprintf("%s <%s>", t.Name, t.Email)
}

type Message struct {
	PlainText, HTML string
}

type Service interface {
	Send(subject string, target Target, cc, bcc *[]Target, message Message, attachments *[]Attachment) (string, error)
	Status(messageID string) string
}
