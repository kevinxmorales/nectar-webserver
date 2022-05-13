package db

import (
	"context"
	"fmt"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/user"
)

type UserRowWithCreds struct {
	ID        string
	FirstName string
	Email     string
	Password  string
}

func convertUserRowWithCredsToUser(row UserRowWithCreds) user.User {
	return user.User{
		ID:       row.ID,
		Name:     row.FirstName,
		Email:    row.Email,
		Password: row.Password,
	}
}

func (d *Database) GetCredentialsByEmail(ctx context.Context, email string) (user.User, error) {
	query := `SELECT usr_id, usr_frst_nm, usr_email, usr_psswrd
				FROM users
				WHERE usr_email = $1`
	var userRow UserRowWithCreds
	row := d.Client.QueryRowContext(ctx, query, email)
	err := row.Scan(&userRow.ID, &userRow.FirstName, &userRow.Email, &userRow.Password)
	if err != nil {
		return user.User{}, fmt.Errorf("error fetching user by email. %w", err)
	}
	return convertUserRowWithCredsToUser(userRow), nil
}
