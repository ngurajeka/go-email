package email

import (
	"gopkg.in/gomail.v2"
)

//Authenticate to the email server
func (a *Account) Authenticate() (bool, []error) {
	if ok, errsAccount := a.Validate(); !ok {
		return false, errsAccount
	}

	auth := gomail.NewDialer(a.GetHost(), a.GetPort(), a.GetUsername(), a.GetPassword())

	a.updateDialer(auth)

	return true, nil
}

func (a *Account) updateDialer(dialer *gomail.Dialer) {
	a.Auth = dialer
}
