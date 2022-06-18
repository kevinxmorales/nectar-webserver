//go:build integration

package db

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/user"
	"testing"
)

func TestUserDatabase(t *testing.T) {
	t.Run("test create user", func(t *testing.T) {
		db, err := NewDatabase()
		assert.NoError(t, err)

		firstName, lastName, email, password := "Kevin", "Morales", "kevin@email.com", "my_password"
		insertedUser, err := db.AddUser(context.Background(), user.User{
			FirstName: firstName,
			LastName:  lastName,
			Email:     email,
			Password:  password,
		})
		assert.NoError(t, err)
		newUser, err := db.GetUser(context.Background(), insertedUser.ID)
		assert.Equal(t, firstName, newUser.FirstName)
		db.DeleteUser(context.Background(), newUser.ID)
	})

	t.Run("test delete user", func(t *testing.T) {
		db, err := NewDatabase()
		assert.NoError(t, err)
		firstName, lastName, email, password := "Kevin", "Morales", "kevin@email.com", "my_password"
		newUser, err := db.AddUser(context.Background(), user.User{
			FirstName: firstName,
			LastName:  lastName,
			Email:     email,
			Password:  password,
		})
		assert.NoError(t, err)

		err = db.DeleteUser(context.Background(), newUser.ID)
		assert.NoError(t, err)

		_, err = db.GetUser(context.Background(), newUser.ID)
		assert.Error(t, err)
	})

	t.Run("test updating user", func(t *testing.T) {
		db, err := NewDatabase()
		assert.NoError(t, err)

		firstName, lastName, email, password := "Kevin", "Morales", "kevin@email.com", "my_password"
		newUser, err := db.AddUser(context.Background(), user.User{
			FirstName: firstName,
			LastName:  lastName,
			Email:     email,
			Password:  password,
		})
		assert.NoError(t, err)
		newEmail := "kevin@protonmail.com"
		updatedUser, err := db.UpdateUser(context.Background(), newUser.ID, user.User{
			FirstName: newUser.FirstName,
			LastName:  newUser.LastName,
			Email:     newEmail,
			Password:  newUser.Password,
		})
		assert.NoError(t, err)

		assert.Equal(t, newEmail, updatedUser.Email)
		err = db.DeleteUser(context.Background(), updatedUser.ID)
		assert.NoError(t, err)
	})

	t.Run("test getting a user that does not exist", func(t *testing.T) {
		db, err := NewDatabase()
		assert.NoError(t, err)

		idNotInDB := uuid.NewV4().String()
		_, err = db.GetUser(context.Background(), idNotInDB)
		assert.Error(t, err)

	})

	t.Run("test adding a user with an already registered email", func(t *testing.T) {
		db, err := NewDatabase()
		assert.NoError(t, err)

		firstName, lastName, email, password := "Kevin", "Morales", "kevin1234@email.com", "my_password"
		insertedUser, err := db.AddUser(context.Background(), user.User{
			FirstName: firstName,
			LastName:  lastName,
			Email:     email,
			Password:  password,
		})
		assert.NoError(t, err)

		//This should fail, expecting an error
		_, err = db.AddUser(context.Background(), *insertedUser)
		assert.Error(t, err)

		err = db.DeleteUser(context.Background(), insertedUser.ID)
		assert.NoError(t, err)
	})
}
