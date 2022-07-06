package util

import (
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"net/mail"
)

var InvalidEmailFormatError = errors.New("not a valid email address format")

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func IsValidUUID(u string) bool {
	_, err := uuid.FromString(u)
	return err == nil
}

func CreateInvalidUuidError(uuid string) error {
	return fmt.Errorf("invalid uuid: %s", uuid)
}
