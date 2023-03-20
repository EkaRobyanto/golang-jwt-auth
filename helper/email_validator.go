package helper

import (
	"net/mail"
)

func ValidateEmail(e string) (string, bool) {
	addr, err := mail.ParseAddress(e)
	if err != nil {
		return "", false
	}
	return addr.Address, true
}
