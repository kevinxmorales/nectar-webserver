package validation

import (
	"errors"
	"net/mail"
)

var InvalidEmailFormatError = errors.New("not a valid email address format")

func IsValidEmail(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}
