package email

import (
	"errors"
	"net/mail"
)

var (
	ErrInvalidName     = errors.New("invalid name")
	ErrInvalidEmail    = errors.New("invalid email address")
	ErrInvalidHost     = errors.New("invalid host")
	ErrInvalidUsername = errors.New("invalid username")
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidPort     = errors.New("invalid port")
)

type Account struct {
	Auth       interface{}
	Credential *Credential
	Mail       mail.Address
}

func NewAccount(name, email string) *Account {
	a := &Account{}
	a.SetMail(name, email)
	a.Credential = &Credential{}

	return a
}

func (a *Account) SetMail(name, email string) {
	a.Mail = mail.Address{Name: name, Address: email}
}

func (a *Account) SetCredential(host, username, password string, port int) {
	a.Credential = NewCredential(host, username, password, port)
}

func (a *Account) Validate() (bool, []error) {
	var errs []error

	if a.Mail.Name == "" {
		errs = append(errs, ErrInvalidName)
	}

	if a.Mail.Address == "" {
		errs = append(errs, ErrInvalidEmail)
	}

	_, errsCredential := a.Credential.Validate()
	if len(errsCredential) > 0 {
		for _, err := range errsCredential {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return false, errs
	}

	return true, nil
}

func (a *Account) GetHost() string {
	return a.Credential.Host
}

func (a *Account) GetPort() int {
	return a.Credential.Port
}

func (a *Account) GetUsername() string {
	return a.Credential.Username
}

func (a *Account) GetPassword() string {
	return a.Credential.Password
}

type Credential struct {
	Host, Username, Password string
	Port                     int
}

func NewCredential(host, username, password string, port int) *Credential {
	return &Credential{
		Host:     host,
		Username: username,
		Password: password,
		Port:     port,
	}
}

func (c *Credential) SetHost(host string) {
	c.Host = host
}

func (c *Credential) SetUsername(username string) {
	c.Username = username
}

func (c *Credential) SetPassword(password string) {
	c.Password = password
}

func (c *Credential) SetPort(port int) {
	c.Port = port
}

func (c *Credential) Validate() (bool, []error) {
	var errs []error

	if c.Host == "" {
		errs = append(errs, ErrInvalidHost)
	}

	if c.Username == "" {
		errs = append(errs, ErrInvalidUsername)
	}

	if c.Password == "" {
		errs = append(errs, ErrInvalidPassword)
	}

	if c.Port == 0 {
		errs = append(errs, ErrInvalidPort)
	}

	if len(errs) > 0 {
		return false, errs
	}

	return true, nil
}
