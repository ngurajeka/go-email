package email

import (
	"errors"
	"net/mail"
	"net/textproto"
)

const (
	MIME_HTML  = "text/html"
	MIME_PLAIN = "text/plain"

	EMPTY_SUBJECT = "Empty Subject"
	EMPTY_FROM    = "Empty Sender"
	EMPTY_TO      = "Empty Target"
	EMPTY_BODY    = "Empty Body"
	EMPTY_HTML    = "Empty Html Body"
)

type Param map[string]interface{}

type Attachments struct {
	Name, Extension, Path string
}

func NewAttachment(file string) *Attachments {

	return &Attachments{Path: file}
}

type Message struct {
	Subject        string
	To, Cc, Bcc    []mail.Address
	From, ReplyTo  mail.Address
	Headers        textproto.MIMEHeader
	Attachments    []*Attachments
	Body, HTMLBody []byte
	Params         Param
	Source         map[string]string
}

func Default() *Message {

	return &Message{
		To:     []mail.Address{},
		From:   mail.Address{},
		Params: Param{},
	}
}

func (m *Message) SetFrom(name, email string) {
	m.From = mail.Address{Name: name, Address: email}
}

func (m *Message) SetReplyTo(name, email string) {
	m.ReplyTo = mail.Address{Name: name, Address: email}
}

func (m *Message) AddTo(name, email string) {
	m.To = append(m.To, mail.Address{Name: name, Address: email})
}

func (m *Message) GetTo() []string {
	var list []string
	for _, to := range m.To {
		list = append(list, to.String())
	}
	return list
}

func (m *Message) SetTo(target mail.Address) {
	m.To = []mail.Address{target}
}

func (m *Message) AddCc(name, email string) {
	m.Cc = append(m.Cc, mail.Address{Name: name, Address: email})
}

func (m *Message) GetCc() []string {
	var list []string
	for _, cc := range m.Cc {
		list = append(list, cc.String())
	}
	return list
}

func (m *Message) AddBcc(name, email string) {
	m.Bcc = append(m.Bcc, mail.Address{Name: name, Address: email})
}

func (m *Message) GetBcc() []string {
	var list []string
	for _, bcc := range m.Bcc {
		list = append(list, bcc.String())
	}
	return list
}

func (m *Message) SetSubject(subject string) {
	m.Subject = subject
}

func (m *Message) SetBody(body []byte) {
	m.Body = body
}

func (m *Message) SetHTMLBody(body []byte) {
	m.HTMLBody = body
}

func (m *Message) AddHeader() {
}

func (m *Message) AddParam(key string, value interface{}) {

	if m.Params == nil {
		m.ResetParam()
	}

	m.Params[key] = value
}

func (m *Message) AddParams(params Param) {

	if m.Params == nil {
		m.ResetParam()
	}

	for k, v := range params {
		m.Params[k] = v
	}
}

func (m *Message) ResetParam() {
	m.Params = make(Param)
}

func (m *Message) Attach(file string) {
	m.Attachments = append(m.Attachments, &Attachments{Path: file})
}

func (m *Message) Check() (bool, []error) {

	var (
		errs []error
	)

	if m.Subject == "" {
		errs = append(errs, errors.New(EMPTY_SUBJECT))
	}

	if len(m.To) == 0 {
		errs = append(errs, errors.New(EMPTY_TO))
	}

	if m.From == *new(mail.Address) {
		errs = append(errs, errors.New(EMPTY_FROM))
	}

	if len(errs) > 0 {
		return false, errs
	}

	return true, errs
}
