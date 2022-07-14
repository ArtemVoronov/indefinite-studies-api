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

	TEST_USER_LOGIN_TEMPLATE   string = "Test user "
	TEST_USER_EMAIL_TEMPLATE   string = "user%v@somewhere.com"
	TEST_USER_PASSORD_TEMPLATE string = "Test password "
)

func GenerateUserLogin(template string, id int) string {
	return template + strconv.Itoa(id)
}

func GenerateUserPassword(template string, id int) string {
	return template + strconv.Itoa(id)
}

func GenerateUserEmail(template string, id int) string {
	return fmt.Sprintf(template, id)
}

func GenerateUser(id int) entities.User {
	return entities.User{
		Id:       id,
		Login:    GenerateUserLogin(TEST_USER_LOGIN_TEMPLATE, id),
		Email:    GenerateUserEmail(TEST_USER_EMAIL_TEMPLATE, id),
		Password: GenerateUserPassword(TEST_USER_PASSORD_TEMPLATE, id),
		Role:     TEST_USER_ROLE_1,
		State:    TEST_USER_STATE_1,
	}
}

func AssertEqualUsers(t *testing.T, expected entities.User, actual entities.User) {
	assert.Equal(t, expected.Id, actual.Id)
	assert.Equal(t, expected.Login, actual.Login)
	assert.Equal(t, expected.Email, actual.Email)
	assert.Equal(t, expected.Password, actual.Password)
	assert.Equal(t, expected.State, actual.State)
}

func AssertEqualUserArrays(t *testing.T, expected []entities.User, actual []entities.User) {
	assert.Equal(t, len(expected), len(actual))

	length := len(expected)
	for i := 0; i < length; i++ {
		AssertEqualUsers(t, expected[i], actual[i])
	}
}

func CreateUserInDB(t *testing.T, login string, email string, password string, role string, state string) int {
	userId, err := queries.CreateUser(db.DB, login, email, password, role, state)
	assert.Nil(t, err)
	assert.NotEqual(t, userId, -1)
	return userId
}

func CreateUsersInDB(t *testing.T, count int, loginTemplate string, emailTemplate string, passwordTemplate string, role string, state string) {
	for i := 1; i <= count; i++ {
		CreateUserInDB(t, GenerateUserLogin(loginTemplate, i), GenerateUserEmail(emailTemplate, i), GenerateUserPassword(passwordTemplate, i), role, state)
	}
}

func TestDBUserGet(t *testing.T) {
	t.Run("NotFoundCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		_, actualError := queries.GetUser(db.DB, 1)

		assert.Equal(t, sql.ErrNoRows, actualError)
	})))
	t.Run("BasicCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expected := GenerateUser(1)

		userId, err := queries.CreateUser(db.DB, expected.Login, expected.Email, expected.Password, expected.Role, expected.State)

		assert.Nil(t, err)
		assert.Equal(t, userId, expected.Id)

		actual, err := queries.GetUser(db.DB, userId)

		AssertEqualUsers(t, expected, actual)
	})))
}

func TestDBUserCreate(t *testing.T) {
	t.Run("BasicCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		userId, err := queries.CreateUser(db.DB, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)

		assert.Nil(t, err)
		assert.Equal(t, userId, 1)
	})))
	t.Run("DuplicateCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		userId, err := queries.CreateUser(db.DB, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)

		assert.Nil(t, err)
		assert.NotEqual(t, userId, -1)

		_, err = queries.CreateUser(db.DB, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)

		assert.Equal(t, db.ErrorUserDuplicateKey, err)
	})))
}

func TestDBUserGetAll(t *testing.T) {
	t.Run("ExpectedEmpty", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		users, err := queries.GetUsers(db.DB, 50, 0)

		assert.Nil(t, err)
		assert.Equal(t, 0, len(users))
	})))
	t.Run("BasicCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		var expectedUsers []entities.User
		for i := 1; i <= 10; i++ {
			expectedUsers = append(expectedUsers, GenerateUser(i))
		}
		CreateUsersInDB(t, 10, TEST_USER_LOGIN_TEMPLATE, TEST_USER_EMAIL_TEMPLATE, TEST_USER_PASSORD_TEMPLATE, TEST_USER_ROLE_1, TEST_USER_STATE_1)

		actualUsers, err := queries.GetUsers(db.DB, 50, 0)

		assert.Nil(t, err)
		AssertEqualUserArrays(t, expectedUsers, actualUsers)
	})))
	t.Run("LimitParameterCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		var expectedUsers []entities.User
		for i := 1; i <= 5; i++ {
			expectedUsers = append(expectedUsers, GenerateUser(i))
		}
		CreateUsersInDB(t, 10, TEST_USER_LOGIN_TEMPLATE, TEST_USER_EMAIL_TEMPLATE, TEST_USER_PASSORD_TEMPLATE, TEST_USER_ROLE_1, TEST_USER_STATE_1)

		actualUsers, err := queries.GetUsers(db.DB, 5, 0)

		assert.Nil(t, err)
		AssertEqualUserArrays(t, expectedUsers, actualUsers)
	})))
	t.Run("OffsetParameterCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		var expectedUsers []entities.User
		for i := 6; i <= 10; i++ {
			expectedUsers = append(expectedUsers, GenerateUser(i))
		}
		CreateUsersInDB(t, 10, TEST_USER_LOGIN_TEMPLATE, TEST_USER_EMAIL_TEMPLATE, TEST_USER_PASSORD_TEMPLATE, TEST_USER_ROLE_1, TEST_USER_STATE_1)

		actualUsers, err := queries.GetUsers(db.DB, 50, 5)

		assert.Nil(t, err)
		AssertEqualUserArrays(t, expectedUsers, actualUsers)
	})))
}

func TestDBUserUpdate(t *testing.T) {
	t.Run("NotFoundCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		err := queries.UpdateUser(db.DB, 1, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)

		assert.Equal(t, sql.ErrNoRows, err)
	})))
	t.Run("DeletedCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		userId, err := queries.CreateUser(db.DB, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)

		assert.Nil(t, err)
		assert.NotEqual(t, userId, -1)

		err = queries.DeleteUser(db.DB, userId)

		assert.Nil(t, err)

		err = queries.UpdateUser(db.DB, userId, TEST_USER_LOGIN_2, TEST_USER_EMAIL_2, TEST_USER_PASSWORD_2, TEST_USER_ROLE_2, TEST_USER_STATE_2)

		assert.Equal(t, sql.ErrNoRows, err)
	})))
	t.Run("BasicCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expected := GenerateUser(1)

		userId, err := queries.CreateUser(db.DB, TEST_USER_LOGIN_2, TEST_USER_EMAIL_2, TEST_USER_PASSWORD_2, TEST_USER_ROLE_2, TEST_USER_STATE_2)

		assert.Nil(t, err)
		assert.Equal(t, expected.Id, userId)

		err = queries.UpdateUser(db.DB, expected.Id, expected.Login, expected.Email, expected.Password, expected.Role, expected.State)

		assert.Nil(t, err)

		actual, err := queries.GetUser(db.DB, expected.Id)

		AssertEqualUsers(t, expected, actual)
	})))
	t.Run("DuplicateCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		userId, err := queries.CreateUser(db.DB, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)

		assert.Nil(t, err)
		assert.NotEqual(t, userId, -1)

		userId, err = queries.CreateUser(db.DB, TEST_USER_LOGIN_2, TEST_USER_EMAIL_2, TEST_USER_PASSWORD_2, TEST_USER_ROLE_2, TEST_USER_STATE_2)

		assert.Nil(t, err)
		assert.NotEqual(t, userId, -1)

		actualError := queries.UpdateUser(db.DB, userId, TEST_USER_LOGIN_2, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_2, TEST_USER_ROLE_2, TEST_USER_STATE_1)

		assert.Equal(t, db.ErrorUserDuplicateKey, actualError)
	})))
}

func TestDBUserDelete(t *testing.T) {
	t.Run("NotFoundCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		err := queries.DeleteUser(db.DB, 1)

		assert.Equal(t, sql.ErrNoRows, err)
	})))
	t.Run("AlreadyDeletedCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		userId, err := queries.CreateUser(db.DB, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)

		assert.Nil(t, err)
		assert.NotEqual(t, userId, -1)

		err = queries.DeleteUser(db.DB, userId)

		assert.Nil(t, err)

		err = queries.DeleteUser(db.DB, userId)

		assert.Equal(t, sql.ErrNoRows, err)
	})))
	t.Run("BasicCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		var expectedUsers []entities.User
		expectedUsers = append(expectedUsers, GenerateUser(1))
		expectedUsers = append(expectedUsers, GenerateUser(3))

		userIdToDelete := 2

		CreateUsersInDB(t, 3, TEST_USER_LOGIN_TEMPLATE, TEST_USER_EMAIL_TEMPLATE, TEST_USER_PASSORD_TEMPLATE, TEST_USER_ROLE_1, TEST_USER_STATE_1)

		err := queries.DeleteUser(db.DB, userIdToDelete)

		assert.Nil(t, err)

		users, err := queries.GetUsers(db.DB, 50, 0)

		assert.Nil(t, err)
		AssertEqualUserArrays(t, expectedUsers, users)

		_, err = queries.GetUser(db.DB, userIdToDelete)

		assert.Equal(t, sql.ErrNoRows, err)
	})))
}
