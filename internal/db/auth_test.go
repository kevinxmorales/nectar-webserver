//go:build integration

package db

import (
	"context"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/user"
	"testing"
)

func TestAuthDatabase(t *testing.T) {

	t.Run("test get user credentials", func(t *testing.T) {
		db, err := NewDatabase()
		assert.NoError(t, err)

		email := uuid.NewV4().String()
		//Create new user with credentials
		firstName, lastName, email, password := "Kevin", "Morales", email+"@email.com", "my_password"
		insertedUser, err := db.AddUser(context.Background(), user.User{
			FirstName: firstName,
			LastName:  lastName,
			Email:     email,
			Password:  password,
		})
		assert.NoError(t, err)
		log.Infof("%v", insertedUser)
		assert.Equal(t, email, insertedUser.Email)

		//Try to get user and credentials
		usr, err := db.GetCredentialsByEmail(context.Background(), insertedUser.Email)
		assert.NoError(t, err)

		//Assert that user and credentials were store correctly
		assert.Equal(t, email, usr.Email)

		//Assert that the password was not stored in plaintext
		assert.NotEqual(t, password, usr.Password)

		//Clean up the db
		err = db.DeleteUser(context.Background(), usr.Id)
		assert.NoError(t, err)
	})

	t.Run("test get user credentials failure", func(t *testing.T) {
		db, err := NewDatabase()
		assert.NoError(t, err)

		email := uuid.NewV4().String()
		//Create new user with credentials
		firstName, lastName, email, password := "Kevin", "Morales", email+"@yahoo.com", "my_password"
		insertedUser, err := db.AddUser(context.Background(), user.User{
			FirstName: firstName,
			LastName:  lastName,
			Email:     email,
			Password:  password,
		})
		assert.Equal(t, email, insertedUser.Email)

		//Try to get user and credentials with different email
		diffEmail := "kevin@gmail.com"
		_, err = db.GetCredentialsByEmail(context.Background(), diffEmail)
		assert.Error(t, err)
	})

}
