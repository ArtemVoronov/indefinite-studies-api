//go:build integration
// +build integration

package queries_test

import (
	"testing"

	integrationTesting "github.com/ArtemVoronov/indefinite-studies-api/internal/app/testing"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/db"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/db/entities"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/db/queries"
	"github.com/stretchr/testify/assert"
)

const (
	TEST_USER_LOGIN_1    string = "Test user 1"
	TEST_USER_EMAIL_1    string = "user1@somewhere.com"
	TEST_USER_PASSWORD_1 string = "16d99c6502225c7e8ee5c85af1070cbcf04724763836ad29edaedab552a54c63d79fb04f62e7a8b4a4b849a6edc558010a67b9b57a949aaf425c6a0dc821fa2d"
	TEST_USER_ROLE_1     string = entities.USER_ROLE_OWNER
	TEST_USER_STATE_1    string = entities.USER_STATE_NEW
)

func TestGetUser(t *testing.T) {
	t.Run("ExpectedNotFoundError", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		t.Errorf("Not implemented")
	})))
	t.Run("ExpectedResult", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		t.Errorf("Not implemented")
	})))
}

func TestCreateUser(t *testing.T) {
	t.Run("BasicCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedUserId := 1

		actualUserId, err := queries.CreateUser(db.DB, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)
		if err != nil || actualUserId == -1 {
			t.Errorf("Unable to create task: %s", err)
		}

		assert.Equal(t, expectedUserId, actualUserId)
	})))
}

func TestGetUsers(t *testing.T) {
	t.Run("ExpectedEmpty", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		t.Errorf("Not implemented")
	})))
	t.Run("ExpectedResult", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		t.Errorf("Not implemented")
	})))
	t.Run("LimitParameterCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		t.Errorf("Not implemented")
	})))
	t.Run("OffsetParameterCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		t.Errorf("Not implemented")
	})))
}

func TestUpdateUser(t *testing.T) {
	t.Run("ExpectedNotFoundError", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		t.Errorf("Not implemented")
	})))
	t.Run("DeletedCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		t.Errorf("Not implemented")
	})))
	t.Run("BlockedCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		t.Errorf("Not implemented")
	})))
	t.Run("BasicCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		t.Errorf("Not implemented")
	})))
}

func TestDeleteUser(t *testing.T) {
	t.Run("NotFoundCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		t.Errorf("Not implemented")
	})))
	t.Run("AlreadyDeletedCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		t.Errorf("Not implemented")
	})))
	t.Run("BasicCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		t.Errorf("Not implemented")
	})))
}
