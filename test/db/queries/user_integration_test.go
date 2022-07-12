//go:build integration
// +build integration

package queries_test

import (
	"database/sql"
	"fmt"
	"strconv"
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
	TEST_USER_LOGIN_2    string = "Test user 2"
	TEST_USER_EMAIL_2    string = "user2@somewhere.com"
	TEST_USER_PASSWORD_2 string = "17d99c6502225c7e8ee5c85af1070cbcf04724763836ad29edaedab552a54c63d79fb04f62e7a8b4a4b849a6edc558010a67b9b57a949aaf425c6a0dc821fa2d"
	TEST_USER_ROLE_2     string = entities.USER_ROLE_RESIDENT
	TEST_USER_STATE_2    string = entities.USER_STATE_BLOCKED

	TEST_USER_LOGIN_TEMPLATE string = "Test user "
	TEST_USER_EMAIL_TEMPLATE string = "@somewhere.com"
)

var UserDuplicateKeyConstraintViolationError = fmt.Errorf(DuplicateKeyConstraintViolationError, "users_email_unique")

func TestGetUser(t *testing.T) {
	t.Run("ExpectedNotFoundError", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedError := sql.ErrNoRows

		_, actualError := queries.GetUser(db.DB, 1)

		assert.Equal(t, expectedError, actualError)
	})))
	t.Run("ExpectedResult", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedLogin := TEST_USER_LOGIN_1
		expectedEmail := TEST_USER_EMAIL_1
		expectedPassword := TEST_USER_PASSWORD_1
		expectedRole := TEST_USER_ROLE_1
		expectedState := TEST_USER_STATE_1
		expectedId, err := queries.CreateUser(db.DB, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)

		if err != nil || expectedId == -1 {
			t.Errorf("Unable to create user: %s", err)
		}

		actual, err := queries.GetUser(db.DB, expectedId)

		assert.Equal(t, expectedId, actual.Id)
		assert.Equal(t, expectedLogin, actual.Login)
		assert.Equal(t, expectedEmail, actual.Email)
		assert.Equal(t, expectedPassword, actual.Password)
		assert.Equal(t, expectedRole, actual.Role)
		assert.Equal(t, expectedState, actual.State)
	})))
}

func TestCreateUser(t *testing.T) {
	t.Run("BasicCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedUserId := 1

		actualUserId, err := queries.CreateUser(db.DB, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)
		if err != nil || actualUserId == -1 {
			t.Errorf("Unable to create user: %s", err)
		}

		assert.Equal(t, expectedUserId, actualUserId)
	})))
	t.Run("DuplicateCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedError := fmt.Errorf("error at inserting user (Login: '%s', Email: '%s') into db, case after db.QueryRow.Scan: %s", TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, UserDuplicateKeyConstraintViolationError)

		userId, err := queries.CreateUser(db.DB, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)
		if err != nil || userId == -1 {
			t.Errorf("Unable to create user: %s", err)
		}

		_, actualError := queries.CreateUser(db.DB, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)

		assert.Equal(t, expectedError, actualError)
	})))
}

func TestGetUsers(t *testing.T) {
	t.Run("ExpectedEmpty", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedArrayLength := 0

		users, err := queries.GetUsers(db.DB, "50", "0")
		if err != nil {
			t.Errorf("Unable to get to users : %s", err)
		}
		actualArrayLength := len(users)

		assert.Equal(t, expectedArrayLength, actualArrayLength)
	})))
	t.Run("ExpectedResult", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedArrayLength := 3

		for i := 0; i < 3; i++ {
			userId, err := queries.CreateUser(db.DB, TEST_USER_LOGIN_TEMPLATE+strconv.Itoa(i), "user"+strconv.Itoa(i)+TEST_USER_EMAIL_TEMPLATE, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)
			if err != nil || userId == -1 {
				t.Errorf("Unable to create user: %s", err)
			}
		}
		users, err := queries.GetUsers(db.DB, "50", "0")
		if err != nil {
			t.Errorf("Unable to get to users : %s", err)
		}
		actualArrayLength := len(users)

		assert.Equal(t, expectedArrayLength, actualArrayLength)
		for i, user := range users {
			assert.Equal(t, i+1, user.Id)
			assert.Equal(t, TEST_USER_LOGIN_TEMPLATE+strconv.Itoa(i), user.Login)
			assert.Equal(t, "user"+strconv.Itoa(i)+TEST_USER_EMAIL_TEMPLATE, user.Email)
			assert.Equal(t, TEST_USER_PASSWORD_1, user.Password)
			assert.Equal(t, TEST_USER_ROLE_1, user.Role)
			assert.Equal(t, TEST_USER_STATE_1, user.State)
		}
	})))
	t.Run("LimitParameterCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedArrayLength := 5
		for i := 0; i < 10; i++ {
			userId, err := queries.CreateUser(db.DB, TEST_USER_LOGIN_TEMPLATE+strconv.Itoa(i), "user"+strconv.Itoa(i)+TEST_USER_EMAIL_TEMPLATE, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)
			if err != nil || userId == -1 {
				t.Errorf("Unable to create user: %s", err)
			}
		}

		users, err := queries.GetUsers(db.DB, "5", "0")
		if err != nil {
			t.Errorf("Unable to get to users : %s", err)
		}
		actualArrayLength := len(users)

		assert.Equal(t, expectedArrayLength, actualArrayLength)
		for i, user := range users {
			assert.Equal(t, i+1, user.Id)
			assert.Equal(t, TEST_USER_LOGIN_TEMPLATE+strconv.Itoa(i), user.Login)
			assert.Equal(t, "user"+strconv.Itoa(i)+TEST_USER_EMAIL_TEMPLATE, user.Email)
			assert.Equal(t, TEST_USER_PASSWORD_1, user.Password)
			assert.Equal(t, TEST_USER_ROLE_1, user.Role)
			assert.Equal(t, TEST_USER_STATE_1, user.State)
		}
	})))
	t.Run("OffsetParameterCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedArrayLength := 5
		for i := 0; i < 10; i++ {
			userId, err := queries.CreateUser(db.DB, TEST_USER_LOGIN_TEMPLATE+strconv.Itoa(i), "user"+strconv.Itoa(i)+TEST_USER_EMAIL_TEMPLATE, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)
			if err != nil || userId == -1 {
				t.Errorf("Unable to create user: %s", err)
			}
		}

		users, err := queries.GetUsers(db.DB, "50", "5")
		if err != nil {
			t.Errorf("Unable to get to users : %s", err)
		}
		actualArrayLength := len(users)

		assert.Equal(t, expectedArrayLength, actualArrayLength)
		for i, user := range users {
			assert.Equal(t, i+6, user.Id)
			assert.Equal(t, TEST_USER_LOGIN_TEMPLATE+strconv.Itoa(i+5), user.Login)
			assert.Equal(t, "user"+strconv.Itoa(i+5)+TEST_USER_EMAIL_TEMPLATE, user.Email)
			assert.Equal(t, TEST_USER_PASSWORD_1, user.Password)
			assert.Equal(t, TEST_USER_ROLE_1, user.Role)
			assert.Equal(t, TEST_USER_STATE_1, user.State)
		}
	})))
}

func TestUpdateUser(t *testing.T) {
	t.Run("ExpectedNotFoundError", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedError := sql.ErrNoRows

		actualError := queries.UpdateUser(db.DB, 1, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)

		assert.Equal(t, expectedError, actualError)
	})))
	t.Run("DeletedCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedError := sql.ErrNoRows

		userId, err := queries.CreateUser(db.DB, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)
		if err != nil || userId == -1 {
			t.Errorf("Unable to create user: %s", err)
		}

		err = queries.DeleteUser(db.DB, userId)
		if err != nil {
			t.Errorf("Unable to delete user: %s", err)
		}

		actualError := queries.UpdateUser(db.DB, 1, TEST_USER_LOGIN_2, TEST_USER_EMAIL_2, TEST_USER_PASSWORD_2, TEST_USER_ROLE_2, TEST_USER_STATE_2)

		assert.Equal(t, expectedError, actualError)
	})))
	t.Run("BasicCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedLogin := TEST_USER_LOGIN_2
		expectedEmail := TEST_USER_EMAIL_2
		expectedPassword := TEST_USER_PASSWORD_2
		expectedRole := TEST_USER_ROLE_2
		expectedState := TEST_USER_STATE_2
		expectedId, err := queries.CreateUser(db.DB, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)
		if err != nil || expectedId == -1 {
			t.Errorf("Unable to create user: %s", err)
		}

		err = queries.UpdateUser(db.DB, expectedId, TEST_USER_LOGIN_2, TEST_USER_EMAIL_2, TEST_USER_PASSWORD_2, TEST_USER_ROLE_2, TEST_USER_STATE_2)
		if err != nil {
			t.Errorf("Unable to update user: %s", err)
		}

		actual, err := queries.GetUser(db.DB, expectedId)

		assert.Equal(t, expectedId, actual.Id)
		assert.Equal(t, expectedLogin, actual.Login)
		assert.Equal(t, expectedEmail, actual.Email)
		assert.Equal(t, expectedPassword, actual.Password)
		assert.Equal(t, expectedRole, actual.Role)
		assert.Equal(t, expectedState, actual.State)
	})))
	t.Run("DuplicateCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		userId, err := queries.CreateUser(db.DB, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)
		if err != nil || userId == -1 {
			t.Errorf("Unable to create user: %s", err)
		}
		userId, err = queries.CreateUser(db.DB, TEST_USER_LOGIN_2, TEST_USER_EMAIL_2, TEST_USER_PASSWORD_2, TEST_USER_ROLE_2, TEST_USER_STATE_2)
		if err != nil || userId == -1 {
			t.Errorf("Unable to create user: %s", err)
		}
		expectedError := fmt.Errorf("error at updating user (Id: %d, Login: '%s', Email: '%s', State: '%s'), case after executing statement: %s", userId, TEST_USER_LOGIN_2, TEST_USER_EMAIL_1, TEST_USER_STATE_2, UserDuplicateKeyConstraintViolationError)

		actualError := queries.UpdateUser(db.DB, userId, TEST_USER_LOGIN_2, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_2, TEST_USER_ROLE_2, TEST_USER_STATE_2)

		assert.Equal(t, expectedError, actualError)
	})))
}

func TestDeleteUser(t *testing.T) {
	t.Run("NotFoundCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		notExistentUserId := 1

		actualError := queries.DeleteUser(db.DB, notExistentUserId)
		if actualError != nil {
			t.Errorf("Unable to delete user: %s", actualError)
		}
	})))
	t.Run("AlreadyDeletedCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		userId, err := queries.CreateUser(db.DB, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)
		if err != nil || userId == -1 {
			t.Errorf("Unable to create user: %s", err)
		}

		err = queries.DeleteUser(db.DB, userId)
		if err != nil {
			t.Errorf("Unable to delete user: %s", err)
		}

		actualError := queries.DeleteUser(db.DB, userId)
		if actualError != nil {
			t.Errorf("Unable to delete user: %s", actualError)
		}
	})))
	t.Run("BasicCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedFirstUserId := 1
		expectedSecondUserId := 3
		expectedPassword := TEST_USER_PASSWORD_1
		expectedRole := TEST_USER_ROLE_1
		expectedState := TEST_USER_STATE_1
		expectedError := sql.ErrNoRows
		expectedArrayLength := 2
		userIdToDelete := 2
		for i := 0; i < 3; i++ {
			userId, err := queries.CreateUser(db.DB, TEST_USER_LOGIN_TEMPLATE+strconv.Itoa(i), "user"+strconv.Itoa(i)+TEST_USER_EMAIL_TEMPLATE, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)
			if err != nil || userId == -1 {
				t.Errorf("Unable to create user: %s", err)
			}
		}

		err := queries.DeleteUser(db.DB, userIdToDelete)
		if err != nil {
			t.Errorf("Unable to delete user: %s", err)
		}

		users, err := queries.GetUsers(db.DB, "50", "0")
		if err != nil {
			t.Errorf("Unable to get to users : %s", err)
		}
		actualArrayLength := len(users)

		assert.Equal(t, expectedArrayLength, actualArrayLength)

		assert.Equal(t, expectedFirstUserId, users[0].Id)
		assert.Equal(t, TEST_USER_LOGIN_TEMPLATE+"0", users[0].Login)
		assert.Equal(t, "user"+strconv.Itoa(0)+TEST_USER_EMAIL_TEMPLATE, users[0].Email)
		assert.Equal(t, expectedPassword, users[0].Password)
		assert.Equal(t, expectedRole, users[0].Role)
		assert.Equal(t, expectedState, users[0].State)

		assert.Equal(t, expectedSecondUserId, users[1].Id)
		assert.Equal(t, TEST_USER_LOGIN_TEMPLATE+"2", users[1].Login)
		assert.Equal(t, "user"+strconv.Itoa(2)+TEST_USER_EMAIL_TEMPLATE, users[1].Email)
		assert.Equal(t, expectedPassword, users[1].Password)
		assert.Equal(t, expectedRole, users[1].Role)
		assert.Equal(t, expectedState, users[1].State)

		_, actualError := queries.GetUser(db.DB, userIdToDelete)

		assert.Equal(t, expectedError, actualError)
	})))
}
