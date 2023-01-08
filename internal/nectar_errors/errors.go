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
