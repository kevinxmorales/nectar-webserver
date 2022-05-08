package db

import (
	"context"
	"errors"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/user"
)

var errNotImplemented = errors.New("not implemented")

func (d *Database) GetUser(context.Context, string) (user.User, error) {
	return user.User{}, errNotImplemented
}
func (d *Database) AddUser(context.Context, user.User) (user.User, error) {
	return user.User{}, errNotImplemented
}
func (d *Database) DeleteUser(context.Context, string) error {
	return errNotImplemented
}
func (d *Database) UpdateUser(context.Context, string, user.User) (user.User, error) {
	return user.User{}, errNotImplemented
}
