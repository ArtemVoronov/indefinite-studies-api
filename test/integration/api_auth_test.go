//go:build integration
// +build integration

package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ArtemVoronov/indefinite-studies-api/internal/api/rest/v1/auth"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/db/entities"
	"github.com/stretchr/testify/assert"
)

func TestApiAuthenicate(t *testing.T) {
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		id := "1"
		login := "Test user 1"
		email := "user1@somewhere.com"
		password := "Test password 1"
		role := entities.USER_ROLE_OWNER
		state := entities.USER_STATE_NEW

		httpStatusCode, body, err := testHttpClient.CreateUser(login, email, password, role, state)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, httpStatusCode)
		assert.Equal(t, id, body)

		httpStatusCode, body, err = testHttpClient.Authenicate(email, password)

		var result auth.AuthenicationResultDTO
		err = json.Unmarshal([]byte(body), &result)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, httpStatusCode)

		assert.NotNil(t, result.Token)
		assert.NotNil(t, result.ExpiredAt)
		assert.NotEqual(t, "", result.Token)
		assert.NotEqual(t, "", result.ExpiredAt)
	})))
}
