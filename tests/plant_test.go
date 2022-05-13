//go:build e2e

package tests

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestPostPlant(t *testing.T) {
	url := BaseUrl + "/" + Version + "/" + plantEndpoint
	t.Run("can create a plant", func(t *testing.T) {
		id := uuid.NewV4().String()
		client := resty.New()
		resp, err := client.R().
			SetBody(fmt.Sprintf(`{
				"name": "testPlant",
				"userId": "%s"}`, id)).
			SetHeader("Authorization", fmt.Sprintf("Bearer %s", CreateToken())).
			Post(url)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode())
	})

	t.Run("cannot add a plant without JWT", func(t *testing.T) {
		id := uuid.NewV4().String()
		client := resty.New()
		resp, err := client.R().
			SetBody(fmt.Sprintf(`{
				"name": "testPlant",
				"userId": "%s"}`, id)).
			Post(url)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode())
	})
}
