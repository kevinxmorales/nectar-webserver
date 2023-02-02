package nectar_errors

import (
	"fmt"
)

type NoEntityError struct {
	Message string
}

func (e *NoEntityError) Error() string {
	return fmt.Sprintf("parse %v: internal error", e.Message)
}

type BadRequestError struct {
	Message string
}

func (e BadRequestError) Error() string {
	return e.Message
}

type DuplicateKeyError struct{}

func (DuplicateKeyError) Error() string {
	return "Unable to insert record. A duplicate key was found"
}
