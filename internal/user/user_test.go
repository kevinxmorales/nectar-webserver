//go:build integration

package user

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

var contextFailureKey = "shouldFail"

type StoreImpl struct {
	mapDB map[int]User
}

func (s *StoreImpl) AddUser(ctx context.Context, u User) (*User, error) {
	shouldFail := ctx.Value(contextFailureKey)
	if shouldFail == "true" {
		return nil, errors.New("could not add user to db")
	}
	u.Id = rand.Intn(1000)
	s.mapDB[u.Id] = u
	return &u, nil
}

func (s *StoreImpl) GetUser(ctx context.Context, id int) (*User, error) {
	user, ok := s.mapDB[id]
	if !ok {
		return nil, errors.New("unable to get user from db")
	}
	return &user, nil
}

func (s *StoreImpl) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	for _, user := range s.mapDB {
		if user.Email == email {
			return &user, nil
		}
	}
	return nil, errors.New("could not find user in db")
}

func (s *StoreImpl) DeleteUser(ctx context.Context, id int) error {
	_, ok := s.mapDB[id]
	if !ok {
		return errors.New("unable to delete user")
	}
	delete(s.mapDB, id)
	return nil
}

func (s StoreImpl) UpdateUser(ctx context.Context, id int, u User) (*User, error) {
	_, present := s.mapDB[id]
	if present {
		return nil, errors.New("cannot add user to db, duplicate ids")
	}
	s.mapDB[u.Id] = u
	return &u, nil
}

var testUser = User{
	FirstName: "Kevin",
	LastName:  "Morales",
	Email:     "kevin@testEmail.com",
	Password:  "myPassword",
}

func TestUserService(t *testing.T) {

	//Get the mocked data store struct
	store := StoreImpl{
		mapDB: make(map[int]User),
	}
	//Initialize a User Service struct with the mock data store
	service := NewService(&store)

	t.Run("test create user", func(t *testing.T) {
		//Db process should succeed
		ctx := context.WithValue(context.Background(), contextFailureKey, "false")

		//Try to add the user through the service method
		insertedUser, err := service.AddUser(ctx, testUser)
		assert.NoError(t, err)
		assert.Equal(t, testUser.Email, insertedUser.Email)
	})

	t.Run("test create user should fail", func(t *testing.T) {
		//DB process should fail for some reason
		ctx := context.WithValue(context.Background(), "shouldFail", "true")

		//Try to add the user through the service method
		_, err := service.AddUser(ctx, testUser)
		assert.Error(t, err)
	})

	t.Run("test delete user success", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "shouldFail", "false")

		//Try to add the user through the service method
		newUser, err := service.AddUser(ctx, testUser)
		assert.NoError(t, err)
		//Now try to delete, should succeed
		err = service.DeleteUser(ctx, newUser.Id)
		assert.NoError(t, err)
	})

	t.Run("test fail to delete user, not in db", func(t *testing.T) {
		//Now try to delete, should fail
		err := service.DeleteUser(context.Background(), 9999)
		assert.Error(t, err)
	})

}
