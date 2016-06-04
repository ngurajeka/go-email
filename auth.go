package email

import (
	"log"

	"gopkg.in/gomail.v2"
)

//Authenticate to the email server
func (a *Account) Authenticate() (bool, []error) {

	if ok, errsAccount := a.Check(); !ok {
		for _, e := range errsAccount {
			log.Println(e)
		}
		return false, errsAccount
	}

	auth := gomail.NewDialer(
		a.GetHost(), a.GetPort(), a.GetUsername(), a.GetPassword(),
	)

	a.Auth = auth
	return true, *new([]error)
}
