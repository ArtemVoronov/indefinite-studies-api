//go:build integration
// +build integration

package integration

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/ArtemVoronov/indefinite-studies-api/internal/api"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/api/rest/v1/auth"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/db"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/db/queries"
	"github.com/stretchr/testify/assert"
)

func TestApiAuthCreate(t *testing.T) {
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		user := utils.entityGenerators.GenerateUser(1)

		httpStatusCode, body, err := testHttpClient.CreateUser(user.Login, user.Email, user.Password, user.Role, user.State)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, httpStatusCode)
		assert.Equal(t, strconv.Itoa(user.Id), body)

		httpStatusCode, body, err = testHttpClient.Authenicate(user.Email, user.Password)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, httpStatusCode)

		var result auth.AuthenicationResultDTO
		err = json.Unmarshal([]byte(body), &result)

		assert.Nil(t, err)

		assert.NotNil(t, result.AccessToken)
		assert.NotNil(t, result.RefreshToken)
		assert.NotNil(t, result.AccessTokenExpiredAt)
		assert.NotNil(t, result.RefreshTokenExpiredAt)
		assert.NotEqual(t, "", result.AccessToken)
		assert.NotEqual(t, "", result.RefreshToken)
		assert.NotEqual(t, "", result.AccessTokenExpiredAt)
		assert.NotEqual(t, "", result.RefreshTokenExpiredAt)
		assert.Equal(t, 242, len(result.AccessToken))
		assert.Equal(t, 242, len(result.RefreshToken))

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			record, err := queries.GetRefreshTokenByToken(tx, ctx, result.RefreshToken)

			assert.NotNil(t, record)
			assert.Equal(t, record.Token, result.RefreshToken)
			assert.Equal(t, record.UserId, user.Id)

			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			_, err := queries.GetRefreshTokenByToken(tx, ctx, result.AccessToken)

			assert.Equal(t, sql.ErrNoRows, err)

			return err
		})()

	})))
	t.Run("WrongEmail", RunWithRecreateDB((func(t *testing.T) {
		user := utils.entityGenerators.GenerateUser(1)

		httpStatusCode, body, err := testHttpClient.CreateUser(user.Login, user.Email, user.Password, user.Role, user.State)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, httpStatusCode)
		assert.Equal(t, strconv.Itoa(user.Id), body)

		httpStatusCode, body, err = testHttpClient.Authenicate("some_wrong_prefix"+user.Email, user.Password)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_WRONG_PASSWORD_OR_EMAIL+"\"", body)

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			_, err := queries.GetRefreshTokenByUserId(tx, ctx, user.Id)

			assert.Equal(t, sql.ErrNoRows, err)

			return err
		})()
	})))
	t.Run("WrongPassword", RunWithRecreateDB((func(t *testing.T) {
		user := utils.entityGenerators.GenerateUser(1)

		httpStatusCode, body, err := testHttpClient.CreateUser(user.Login, user.Email, user.Password, user.Role, user.State)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, httpStatusCode)
		assert.Equal(t, strconv.Itoa(user.Id), body)

		httpStatusCode, body, err = testHttpClient.Authenicate(user.Email, "some_wrong_prefix"+user.Password)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_WRONG_PASSWORD_OR_EMAIL+"\"", body)

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			_, err := queries.GetRefreshTokenByUserId(tx, ctx, user.Id)

			assert.Equal(t, sql.ErrNoRows, err)

			return err
		})()
	})))
}

func TestApiAuthVerify(t *testing.T) {
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		user := utils.entityGenerators.GenerateUser(1)

		httpStatusCode, body, err := testHttpClient.CreateUser(user.Login, user.Email, user.Password, user.Role, user.State)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, httpStatusCode)
		assert.Equal(t, strconv.Itoa(user.Id), body)

		httpStatusCode, body, err = testHttpClient.Authenicate(user.Email, user.Password)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, httpStatusCode)

		var result auth.AuthenicationResultDTO
		err = json.Unmarshal([]byte(body), &result)

		httpStatusCode, body, err = testHttpClient.Verify(result.AccessToken)
		httpStatusCode, body, err = testHttpClient.Verify(result.RefreshToken)

	})))
}
