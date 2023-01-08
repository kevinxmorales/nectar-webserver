package db

import (
	"context"
	"database/sql"
	"fmt"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/user"
	"strings"
)

type UserRow struct {
	Id         string         `db:"id"`
	PlantCount uint           `db:"plant_count"`
	Name       sql.NullString `db:"name"`
	Email      string         `db:"email"`
	Username   string         `db:"username"`
	ImageUrl   sql.NullString `db:"profile_image"`
}

func convertUserRowToUser(u UserRow) *user.User {
	return &user.User{
		Id:         u.Id,
		Email:      u.Email,
		PlantCount: u.PlantCount,
		Username:   u.Username,
		Name:       u.Name.String,
		ImageUrl:   u.ImageUrl.String,
	}
}

func (d *Database) GetUserById(ctx context.Context, id string) (*user.User, error) {
	tag := "db.user.GetUserById"
	query := `SELECT 
    				nectar_users.id, 
    				nectar_users.first_name as name, 
    				nectar_users.email,
    				nectar_users.username,
    				nectar_users.profile_image
				FROM nectar_users
				WHERE nectar_users.id = $1`
	var rows []UserRow
	if err := d.Client.SelectContext(ctx, &rows, query, id); err != nil {
		return nil, fmt.Errorf("sqlx.SelectContext in %s failed for %v", tag, err)
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("no user with the given auth id: %s", id)
	}
	return convertUserRowToUser(rows[0]), nil
}

func (d *Database) AddUser(ctx context.Context, u user.User) (*user.User, error) {
	tag := "db.user.AddUser"
	query := `INSERT INTO nectar_users 
				(id,
				 email,
				 first_name,
				 username) 
				VALUES 
				(:id,
				 :email,
				 :name, 
				 :username)
				RETURNING id`
	userRow := UserRow{
		Name:     sql.NullString{String: u.Name, Valid: true},
		Email:    u.Email,
		Id:       u.Id,
		Username: u.Username,
	}
	rows, err := d.Client.NamedQueryContext(ctx, query, userRow)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil, DuplicateKeyError
		}
		return nil, fmt.Errorf("NamedQueryContext in %s failed for %v", tag, err)
	}
	defer closeDbRows(rows, query)
	var userID string
	for rows.Next() {
		if err := rows.Scan(&userID); err != nil {
			return nil, fmt.Errorf("rows.Scan in %s failed for %v", tag, err)
		}
	}
	u.Id = userID
	return &u, nil
}

func (d *Database) UpdateUser(ctx context.Context, id string, u user.User) (*user.User, error) {
	tag := "db.user.UpdateUser"
	query := `UPDATE nectar_users
					SET
						first_name = $1,
					    username = $2,
					    email = $3,
					    profile_image = $4
					WHERE nectar_users.id = $5
					RETURNING
						nectar_users.id, 
						nectar_users.first_name as name, 
						nectar_users.email,
						nectar_users.username,
						nectar_users.profile_image`
	row := d.Client.QueryRowContext(ctx, query, u.Name, u.Username, u.Email, u.ImageUrl, id)
	var ur UserRow
	if err := row.Scan(&ur.Id, &ur.Name, &ur.Email, &ur.Username, &ur.ImageUrl); err != nil {
		return nil, fmt.Errorf("rows.Scan in %s failed for %v", tag, err)
	}
	return convertUserRowToUser(ur), nil
}

func (d *Database) CheckIfUsernameIsTaken(ctx context.Context, username string) (bool, error) {
	tag := "db.user.CheckIfUsernameIsTaken"
	query := `SELECT exists 
				(SELECT 1 
				FROM nectar_users 
				WHERE nectar_users.username = $1 
				LIMIT 1)`
	var isUsernameTaken bool
	row := d.Client.QueryRowContext(ctx, query, username)
	err := row.Scan(&isUsernameTaken)
	if err != nil {
		return false, fmt.Errorf("row.Scan in %s failed for %v", tag, err)
	}
	return isUsernameTaken, nil
}

func (d *Database) GetUser(ctx context.Context, id string) (*user.User, error) {
	tag := "db.user.GetUser"
	query := `SELECT 
    				nectar_users.id, 
    				nectar_users.first_name as name, 
    				nectar_users.email,
    				nectar_users.username,
    				nectar_users.profile_image
				FROM nectar_users
				WHERE nectar_users.id = $1`
	var rows []UserRow
	if err := d.Client.SelectContext(ctx, &rows, query, id); err != nil {
		return nil, fmt.Errorf("sqlx.SelectContext in %s failed for %v", tag, err)
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("no user with the given id: %s", id)
	}
	return convertUserRowToUser(rows[0]), nil
}

// DeleteUser - Update the deletion date to today
func (d *Database) DeleteUser(ctx context.Context, id string) error {
	tag := "db.user.DeleteUser"
	query := `UPDATE nectar_users
				SET account_deletion_date = current_timestamp
				WHERE nectar_users.id = $1`
	_, err := d.Client.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("ExecContext in %s failed for %v", tag, err)
	}
	return nil
}
