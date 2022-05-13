package db

import (
	"context"
	"database/sql"
	"fmt"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/user"
)

type UserRow struct {
	ID        string         `db:"id"`
	FirstName string         `db:"first_name"`
	LastName  sql.NullString `db:"last_name"`
	Email     string         `db:"email"`
	Password  string         `db:"password"`
}

func convertUserRowToUser(u UserRow) user.User {
	return user.User{
		ID:    u.ID,
		Name:  u.FirstName,
		Email: u.Email,
	}
}

func (d *Database) GetUser(ctx context.Context, uuid string) (user.User, error) {
	query := `SELECT usr_id, usr_frst_nm, usr_email
				FROM users
				WHERE usr_id = $1`
	var userRow UserRow
	row := d.Client.QueryRowContext(ctx, query, uuid)
	err := row.Scan(&userRow.ID, &userRow.FirstName, &userRow.Email)
	if err != nil {
		return user.User{}, fmt.Errorf("error fetching user by uuid. %w", err)
	}
	return convertUserRowToUser(userRow), nil
}

func (d *Database) GetUserByEmail(ctx context.Context, email string) (user.User, error) {
	query := `SELECT usr_id, usr_frst_nm, usr_email
				FROM users
				WHERE usr_email = $1`
	log.Info("in Store.GetUser")
	var userRow UserRow
	row := d.Client.QueryRowContext(ctx, query, email)
	err := row.Scan(&userRow.ID, &userRow.FirstName, &userRow.Email)
	if err != nil {
		return user.User{}, fmt.Errorf("error fetching user by email. %w", err)
	}
	return convertUserRowToUser(userRow), nil
}

func (d *Database) AddUser(ctx context.Context, u user.User) (user.User, error) {
	query := `INSERT INTO users 
				(usr_id, 
				 usr_frst_nm, 
				 usr_email, 
				 usr_psswrd) 
				VALUES 
				(:id, 
				 :first_name, 
				 :email, 
				 :password)`
	log.Info("attempting to create user")
	u.ID = uuid.NewV4().String()
	userRow := UserRow{
		FirstName: u.Name,
		ID:        u.ID,
		Email:     u.Email,
		Password:  u.Password,
	}
	log.Info(userRow)
	log.Info(userRow.FirstName)
	rows, err := d.Client.NamedQueryContext(ctx, query, userRow)
	if err != nil {
		return user.User{}, fmt.Errorf("FAILED to insert new user: %w", err)
	}
	if err := rows.Close(); err != nil {
		return user.User{}, fmt.Errorf("FAILED to close rows: %w", err)
	}
	log.Info("leaving AddUser", u)
	return u, nil
}

func (d *Database) DeleteUser(ctx context.Context, id string) error {
	query := `DELETE FROM users
				WHERE usr_id = $1`
	_, err := d.Client.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("FAILED to delete user from database: %w", err)
	}
	return nil
}

func (d *Database) UpdateUser(ctx context.Context, id string, u user.User) (user.User, error) {
	query := `UPDATE users 
				SET 
					usr_frst_nm = :first_name,
					usr_email = :email,
					usr_psswrd = :password
				WHERE usr_id = :id`
	userRow := UserRow{
		ID:        id,
		FirstName: u.Name,
		Email:     u.Email,
		Password:  u.Password,
	}
	rows, err := d.Client.NamedQueryContext(ctx, query, userRow)
	if err != nil {
		return user.User{}, fmt.Errorf("FAILED to update user: %w", err)
	}
	if err := rows.Close(); err != nil {
		return user.User{}, fmt.Errorf("FAILED to close rows: %w", err)
	}
	return convertUserRowToUser(userRow), nil
}
