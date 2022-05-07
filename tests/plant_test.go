//go:build e2e

package tests

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/golang-jwt/jwt"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

func createToken() string {
	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, err := token.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
	if err != nil {
		fmt.Println(err)
	}
	return tokenString
}

func TestPostPlant(t *testing.T) {
	t.Run("can create a plant", func(t *testing.T) {
		id := uuid.NewV4().String()
		client := resty.New()
		resp, err := client.R().
			SetBody(fmt.Sprintf(`{
				"name": "testPlant"
				"userId": "%s"}`, id)).
			SetHeader("Authorization", fmt.Sprintf("Bearer %s", createToken())).
			Post("http://localhost:8080/api/v1/plant")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode())
	})

	t.Run("cannot add a plant without JWT", func(t *testing.T) {
		id := uuid.NewV4().String()
		client := resty.New()
		resp, err := client.R().
			SetBody(fmt.Sprintf(`{
				"name": "testPlant"
				"userId": "%s"}`, id)).
			Post("http://localhost:8080/api/v1/plant")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode())
	})
}
