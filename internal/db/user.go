package db

import (
	"context"
	"fmt"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/user"
	"strings"
)

type UserRow struct {
	Id         int    `db:"id"`
	PlantCount uint   `db:"plant_count"`
	Name       string `db:"name"`
	FirstName  string `db:"first_name"`
	LastName   string `db:"last_name"`
	Email      string `db:"email"`
	Password   string `db:"password"`
	Username   string `db:"username"`
	AuthId     string `db:"auth_id"`
	ImageUrl   string `db:"image_url"`
}

func convertUserRowToUser(u UserRow) *user.User {
	return &user.User{
		Id:         u.Id,
		FirstName:  u.FirstName,
		LastName:   u.LastName,
		Email:      u.Email,
		PlantCount: u.PlantCount,
		Username:   u.Username,
		AuthId:     u.AuthId,
	}
}

func (d *Database) GetUser(ctx context.Context, id int) (*user.User, error) {
	tag := "db.user.GetUser"
	query := `SELECT 
       			id,
       			usr_frst_nm, 
       			usr_lst_nm, 
       			usr_email, 
       			(select count(*) from plants where plants.plnt_usr_id = $1) as plant_count 
				FROM users
				WHERE users.id = $1`
	var userRow UserRow
	row := d.Client.QueryRowContext(ctx, query, id)
	if err := row.Scan(&userRow.Id, &userRow.FirstName, &userRow.LastName, &userRow.Email, &userRow.PlantCount); err != nil {
		return nil, fmt.Errorf("QueryRowContext in %s failed for %v", tag, err)
	}
	return convertUserRowToUser(userRow), nil
}

func (d *Database) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	tag := "db.user.GetUserByEmail"
	query := `SELECT id, usr_frst_nm, usr_lst_nm, usr_email
				FROM users
				WHERE usr_email = $1`
	var userRow UserRow
	row := d.Client.QueryRowContext(ctx, query, email)
	if err := row.Scan(&userRow.Id, &userRow.FirstName, &userRow.LastName, &userRow.Email); err != nil {
		return nil, fmt.Errorf("QueryRowContext in %s failed for %v", tag, err)
	}
	return convertUserRowToUser(userRow), nil
}

func (d *Database) GetUserByAuthId(ctx context.Context, email string) (*user.User, error) {
	tag := "db.user.GetUserByAuthId"
	query := `SELECT id, usr_frst_nm, usr_lst_nm, usr_email, usr_auth_id, usr_username
				FROM users
				WHERE usr_auth_id = $1`
	var u UserRow
	row := d.Client.QueryRowContext(ctx, query, email)
	if err := row.Scan(&u.Id, &u.FirstName, &u.LastName, &u.Email, &u.AuthId, &u.Username); err != nil {
		return nil, fmt.Errorf("QueryRowContext in %s failed for %v", tag, err)
	}
	return convertUserRowToUser(u), nil
}

func (d *Database) AddUser(ctx context.Context, u user.User) (*user.User, error) {
	tag := "db.user.AddUser"
	query := `INSERT INTO users 
				(usr_auth_id,
				 usr_name,
				 usr_email, 
				 usr_username
				 ) 
				VALUES 
				(:auth_id,
				 :name,
				 :email, 
				 :username)
				RETURNING id`
	userRow := UserRow{
		Name:     u.Name,
		Email:    u.Email,
		AuthId:   u.AuthId,
		Username: u.Username,
	}
	rows, err := d.Client.NamedQueryContext(ctx, query, userRow)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil, DuplicateKeyError
		}
		return nil, fmt.Errorf("NamedQueryContext in %s failed for %v", tag, err)
	}
	var userID int
	for rows.Next() {
		if err := rows.Scan(&userID); err != nil {
			return nil, fmt.Errorf("rows.Scan in %s failed for %v", tag, err)
		}
	}
	u.Id = userID
	return &u, nil
}

func (d *Database) DeleteUser(ctx context.Context, id int) error {
	tag := "db.user.DeleteUser"
	query := `DELETE FROM users
				WHERE users.id = $1`
	_, err := d.Client.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("ExecContext in %s failed for %v", tag, err)
	}
	return nil
}

func (d *Database) UpdateUser(ctx context.Context, id int, u user.User) (*user.User, error) {
	tag := "db.user.UpdateUser"
	query := `UPDATE users 
				SET 
					usr_frst_nm = :first_name,
				    usr_lst_nm = :last_name,
					usr_email = :email,
				    usr_username = :username,
					usr_image_url = :image_url
				WHERE users.id = :id`
	userRow := UserRow{
		Id:        id,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
	}
	_, err := d.Client.NamedQueryContext(ctx, query, userRow)
	if err != nil {
		return nil, fmt.Errorf("NamedQueryContext in %s failed for %v", tag, err)
	}
	return convertUserRowToUser(userRow), nil
}

func (d *Database) CheckIfUsernameIsTaken(ctx context.Context, username string) (bool, error) {
	tag := "db.user.CheckIfUsernameIsTaken"
	query := `SELECT exists 
				(SELECT 1 
				FROM users 
				WHERE usr_username = $1 
				LIMIT 1)`
	var isUsernameTaken bool
	row := d.Client.QueryRowContext(ctx, query, username)
	err := row.Scan(&isUsernameTaken)
	if err != nil {
		return false, fmt.Errorf("row.Scan in %s failed for %v", tag, err)
	}
	return isUsernameTaken, nil
}
