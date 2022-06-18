//go:build e2e

package tests

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"gitlab.com/kevinmorales/nectar-rest-api/internal/user"
	"net/http"
	"testing"
)

const url = BaseUrl + "/" + Version + "/" + userEndpoint

func TestPostUserBadEmail(t *testing.T) {

	t.Run("cannot create a user with a bad email address", func(t *testing.T) {
		first, last, email, password := "kevin", "Morales", "kevinEmail.com", "password123"
		client := resty.New()
		resp, err := client.R().
			SetBody(fmt.Sprintf(`{
				"firstName": "%s",
				"lastName": "%s",
				"email": "%s",
				"password": "%s"}`, first, last, email, password)).
			Post(url)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode())
	})
}

func TestPostUserThenGetUser(t *testing.T) {
	t.Run("can get a user from the db", func(t *testing.T) {
		first, last, email, password := "kevin", "Morales", "kevin@email.com", "password123"
		client := resty.New()
		resp, err := client.R().
			SetBody(fmt.Sprintf(`{
				"firstName": "%s",
				"lastName": "%s",
				"email": "%s",
				"password": "%s"}`, first, last, email, password)).
			Post(url)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode())

		body := resp.Body()
		var usr user.User
		err = json.Unmarshal(body, &usr)
		assert.NoError(t, err)
		assert.NotNil(t, usr.ID)

		resp, err = client.R().
			SetHeader("Authorization", fmt.Sprintf("Bearer %s", GetToken())).
			Get(url + "/" + usr.ID)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode())
	})
}

func TestPostUserThenDeleteUser(t *testing.T) {
	t.Run("can delete a user", func(t *testing.T) {
		// create the user
		first, last, email, password := "kevin", "Morales", "kevin@yahoo.com", "password123"
		client := resty.New()
		resp, err := client.R().
			SetBody(fmt.Sprintf(`{
				"firstName": "%s",
				"lastName": "%s",
				"email": "%s",
				"password": "%s"}`, first, last, email, password)).
			SetHeader("Authorization", fmt.Sprintf("Bearer %s", GetToken())).
			Post(url)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode())
		body := resp.Body()
		var usr user.User
		err = json.Unmarshal(body, &usr)
		assert.NoError(t, err)
		assert.NotNil(t, usr.ID)

		// delete the user
		resp, err = client.R().
			SetHeader("Authorization", fmt.Sprintf("Bearer %s", GetToken())).
			Delete(url + "/" + usr.ID)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode())

		// try to get deleted user, expected to fail
		resp, err = client.R().
			SetHeader("Authorization", fmt.Sprintf("Bearer %s", GetToken())).
			Get(url + "/" + usr.ID)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode())
	})
}

func TestPostUserThenPut(t *testing.T) {
	t.Run("can edit a user", func(t *testing.T) {
		// create the user
		first, last, email, password := "kevin", "Morales", "kevin@protonmail.com", "password123"
		client := resty.New()
		resp, err := client.R().
			SetBody(fmt.Sprintf(`{
				"firstName": "%s",
				"lastName": "%s",
				"email": "%s",
				"password": "%s"}`, first, last, email, password)).
			SetHeader("Authorization", fmt.Sprintf("Bearer %s", GetToken())).
			Post(url)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode())

		var usr user.User
		body := resp.Body()
		err = json.Unmarshal(body, &usr)
		assert.NoError(t, err)
		assert.NotNil(t, usr.ID)

		// edit the user with PUT request
		newFirst, newLast, newEmail := "Joe", "Biden", "joe@email.com"
		editUserJson := fmt.Sprintf(`{
				"firstName": "%s",
				"lastName": "%s",
				"email": "%s"
				}`, newFirst, newLast, newEmail)
		resp, err = client.R().
			SetHeader("Authorization", fmt.Sprintf("Bearer %s", GetToken())).
			SetBody(editUserJson).
			Put(url + "/" + usr.ID)
		assert.NoError(t, err)

		var newUser user.User
		body = resp.Body()
		err = json.Unmarshal(body, &newUser)
		assert.NoError(t, err)

		assert.Equal(t, newFirst, newUser.FirstName)
	})
}
