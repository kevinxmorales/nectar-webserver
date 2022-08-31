//go:build integration

package validation

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidationService(t *testing.T) {

	t.Run("test valid email function success", func(t *testing.T) {
		//Create valid email
		email := "kevin@email.com"

		//Call validate email function
		isValid := IsValidEmail(email)

		//Assert that this is a valid email
		assert.True(t, isValid)
	})

	t.Run("test valid email function failure", func(t *testing.T) {
		//Create invalid email
		email := "kevin@email@yahoo.com.edu"

		//Call validate email function
		isValid := IsValidEmail(email)

		//Assert that this is not a valid email
		assert.False(t, isValid)
	})
}
