package db

import (
	"context"
	"fmt"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/user"
)

type UserRowWithCreds struct {
	Id        int
	FirstName string
	LastName  string
	Email     string
	Password  string
}

func convertUserRowWithCredsToUser(row UserRowWithCreds) user.User {
	return user.User{
		Id:        row.Id,
		FirstName: row.FirstName,
		LastName:  row.LastName,
		Email:     row.Email,
		Password:  row.Password,
	}
}

func (d *Database) GetCredentialsByEmail(ctx context.Context, email string) (user.User, error) {
	query := `SELECT id, usr_frst_nm, usr_lst_nm, usr_email, usr_psswrd
				FROM users
				WHERE usr_email = $1`
	var userRow UserRowWithCreds
	row := d.Client.QueryRowContext(ctx, query, email)
	err := row.Scan(&userRow.Id, &userRow.FirstName, &userRow.LastName, &userRow.Email, &userRow.Password)
	if err != nil {
		return user.User{}, fmt.Errorf("error fetching user by email. %w", err)
	}
	return convertUserRowWithCredsToUser(userRow), nil
}
