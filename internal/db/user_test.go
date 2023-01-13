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

		firstName, email, username := "Kevin", "kevin@email.com", "kevin"
		insertedUser, err := db.AddUser(context.Background(), user.User{
			Id:       uuid.NewV4().String(),
			Name:     firstName,
			Email:    email,
			Username: username,
		})
		assert.NoError(t, err)
		newUser, err := db.GetUser(context.Background(), insertedUser.Id)
		assert.NoError(t, err)
		assert.Equal(t, firstName, newUser.Name)
		err = db.DeleteUser(context.Background(), newUser.Id)
		assert.NoError(t, err)
	})

	t.Run("test delete user", func(t *testing.T) {
		db, err := NewDatabase()
		assert.NoError(t, err)
		firstName, username, email := "Kevin", "kevin_m", "kevin@email.com"
		newUser, err := db.AddUser(context.Background(), user.User{
			Id:       uuid.NewV4().String(),
			Name:     firstName,
			Email:    email,
			Username: username,
		})
		assert.NoError(t, err)

		err = db.DeleteUser(context.Background(), newUser.Id)
		assert.NoError(t, err)

		_, err = db.GetUser(context.Background(), newUser.Id)
		assert.Error(t, err)
	})

	t.Run("test updating user", func(t *testing.T) {
		db, err := NewDatabase()
		assert.NoError(t, err)

		firstName, username, email := "Kevin", "kevin_m", "kevin@email.com"
		newUser, err := db.AddUser(context.Background(), user.User{
			Name:     firstName,
			Email:    email,
			Username: username,
			Id:       uuid.NewV4().String(),
		})
		assert.NoError(t, err)
		newEmail := "kevin@protonmail.com"
		updatedUser, err := db.UpdateUser(context.Background(), newUser.Id, user.User{
			Name:     newUser.Name,
			Email:    newEmail,
			Username: newUser.Username,
		})
		assert.NoError(t, err)

		assert.Equal(t, newEmail, updatedUser.Email)
		err = db.DeleteUser(context.Background(), updatedUser.Id)
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

		firstName, username, email := "Kevin", "kevin_m", "kevin1234@email.com"
		insertedUser, err := db.AddUser(context.Background(), user.User{
			Name:     firstName,
			Email:    email,
			Username: username,
		})
		assert.NoError(t, err)

		//This should fail, expecting an error
		_, err = db.AddUser(context.Background(), *insertedUser)
		assert.Error(t, err)
	})
}
