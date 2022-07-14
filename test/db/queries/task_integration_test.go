//go:build integration
// +build integration

package queries_test

import (
	"database/sql"
	"strconv"
	"testing"

	integrationTesting "github.com/ArtemVoronov/indefinite-studies-api/internal/app/testing"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/db"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/db/entities"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/db/queries"
	"github.com/stretchr/testify/assert"
)

const (
	TEST_TASK_NAME_1        string = "Test task 1"
	TEST_TASK_STATE_1       string = entities.TASK_STATE_NEW
	TEST_TASK_NAME_2        string = "Test task 2"
	TEST_TASK_STATE_2       string = entities.TASK_STATE_DONE
	TEST_TASK_NAME_TEMPLATE string = "Test task "
)

func GenerateTask(id int) entities.Task {
	return entities.Task{
		Id:    id,
		Name:  GenerateTaskName(TEST_TASK_NAME_TEMPLATE, id),
		State: TEST_TASK_STATE_1,
	}
}

func GenerateTaskName(template string, id int) string {
	return template + strconv.Itoa(id)
}

func AssertEqualTasks(t *testing.T, expected entities.Task, actual entities.Task) {
	assert.Equal(t, expected.Id, actual.Id)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.State, actual.State)
}

func AssertEqualTaskArrays(t *testing.T, expected []entities.Task, actual []entities.Task) {
	assert.Equal(t, len(expected), len(actual))

	length := len(expected)
	for i := 0; i < length; i++ {
		AssertEqualTasks(t, expected[i], actual[i])
	}
}

func CreateTaskInDB(t *testing.T, name string, state string) int {
	taskId, err := queries.CreateTask(db.GetInstance().GetDB(), name, state)
	assert.Nil(t, err)
	assert.NotEqual(t, taskId, -1)
	return taskId
}

func CreateTasksInDB(t *testing.T, count int, nameTemplate string, state string) {
	for i := 1; i <= count; i++ {
		CreateTaskInDB(t, GenerateTaskName(TEST_TASK_NAME_TEMPLATE, i), state)
	}
}

func TestDBTaskGet(t *testing.T) {
	t.Run("NotFoundCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		_, err := queries.GetTask(db.GetInstance().GetDB(), 1)

		assert.Equal(t, sql.ErrNoRows, err)
	})))
	t.Run("BasicCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expected := GenerateTask(1)

		taskId, err := queries.CreateTask(db.GetInstance().GetDB(), expected.Name, expected.State)

		assert.Nil(t, err)
		assert.Equal(t, taskId, expected.Id)

		actual, err := queries.GetTask(db.GetInstance().GetDB(), taskId)

		AssertEqualTasks(t, expected, actual)
	})))
}

func TestDBTaskCreate(t *testing.T) {
	t.Run("BasicCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		taskId, err := queries.CreateTask(db.GetInstance().GetDB(), TEST_TASK_NAME_1, TEST_TASK_STATE_1)

		assert.Nil(t, err)
		assert.Equal(t, taskId, 1)
	})))
	t.Run("DuplicateCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		taskId, err := queries.CreateTask(db.GetInstance().GetDB(), TEST_TASK_NAME_1, TEST_TASK_STATE_1)

		assert.Nil(t, err)
		assert.NotEqual(t, taskId, -1)

		_, err = queries.CreateTask(db.GetInstance().GetDB(), TEST_TASK_NAME_1, TEST_TASK_STATE_1)

		assert.Equal(t, db.ErrorTaskDuplicateKey, err)
	})))
}

func TestDBTaskGetAll(t *testing.T) {
	t.Run("ExpectedEmpty", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		tasks, err := queries.GetTasks(db.GetInstance().GetDB(), 50, 0)

		assert.Nil(t, err)
		assert.Equal(t, 0, len(tasks))
	})))
	t.Run("BasicCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		var expectedTasks []entities.Task
		for i := 1; i <= 10; i++ {
			expectedTasks = append(expectedTasks, GenerateTask(i))
		}
		CreateTasksInDB(t, 10, TEST_TASK_NAME_TEMPLATE, entities.TASK_STATE_NEW)

		actualTasks, err := queries.GetTasks(db.GetInstance().GetDB(), 50, 0)

		assert.Nil(t, err)
		AssertEqualTaskArrays(t, expectedTasks, actualTasks)
	})))
	t.Run("LimitParameterCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		var expectedTasks []entities.Task
		for i := 1; i <= 5; i++ {
			expectedTasks = append(expectedTasks, GenerateTask(i))
		}

		CreateTasksInDB(t, 10, TEST_TASK_NAME_TEMPLATE, entities.TASK_STATE_NEW)

		actualTasks, err := queries.GetTasks(db.GetInstance().GetDB(), 5, 0)

		assert.Nil(t, err)
		AssertEqualTaskArrays(t, expectedTasks, actualTasks)
	})))
	t.Run("OffsetParameterCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		var expectedTasks []entities.Task
		for i := 6; i <= 10; i++ {
			expectedTasks = append(expectedTasks, GenerateTask(i))
		}

		CreateTasksInDB(t, 10, TEST_TASK_NAME_TEMPLATE, entities.TASK_STATE_NEW)

		actualTasks, err := queries.GetTasks(db.GetInstance().GetDB(), 50, 5)

		assert.Nil(t, err)
		AssertEqualTaskArrays(t, expectedTasks, actualTasks)
	})))
}

func TestDBTaskUpdate(t *testing.T) {
	t.Run("NotFoundCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		err := queries.UpdateTask(db.GetInstance().GetDB(), 1, TEST_TASK_NAME_1, TEST_TASK_STATE_1)

		assert.Equal(t, sql.ErrNoRows, err)
	})))
	t.Run("DeletedCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		taskId, err := queries.CreateTask(db.GetInstance().GetDB(), TEST_TASK_NAME_1, TEST_TASK_STATE_1)

		assert.Nil(t, err)
		assert.NotEqual(t, taskId, -1)

		err = queries.DeleteTask(db.GetInstance().GetDB(), taskId)

		assert.Nil(t, err)

		err = queries.UpdateTask(db.GetInstance().GetDB(), taskId, TEST_TASK_NAME_2, TEST_TASK_STATE_2)

		assert.Equal(t, sql.ErrNoRows, err)
	})))
	t.Run("BasicCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expected := GenerateTask(1)

		taskId, err := queries.CreateTask(db.GetInstance().GetDB(), TEST_TASK_NAME_2, TEST_TASK_STATE_2)

		assert.Nil(t, err)
		assert.Equal(t, expected.Id, taskId)

		err = queries.UpdateTask(db.GetInstance().GetDB(), expected.Id, expected.Name, expected.State)

		assert.Nil(t, err)

		actual, err := queries.GetTask(db.GetInstance().GetDB(), expected.Id)

		AssertEqualTasks(t, expected, actual)
	})))
	t.Run("DuplicateCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		taskId, err := queries.CreateTask(db.GetInstance().GetDB(), TEST_TASK_NAME_1, TEST_TASK_STATE_1)

		assert.Nil(t, err)
		assert.NotEqual(t, taskId, -1)

		taskId, err = queries.CreateTask(db.GetInstance().GetDB(), TEST_TASK_NAME_2, TEST_TASK_STATE_2)

		assert.Nil(t, err)
		assert.NotEqual(t, taskId, -1)

		actualError := queries.UpdateTask(db.GetInstance().GetDB(), taskId, TEST_TASK_NAME_1, TEST_TASK_STATE_1)

		assert.Equal(t, db.ErrorTaskDuplicateKey, actualError)
	})))
}

func TestDBTaskDelete(t *testing.T) {
	t.Run("NotFoundCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		err := queries.DeleteTask(db.GetInstance().GetDB(), 1)

		assert.Equal(t, sql.ErrNoRows, err)
	})))
	t.Run("AlreadyDeletedCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		taskId, err := queries.CreateTask(db.GetInstance().GetDB(), TEST_TASK_NAME_1, TEST_TASK_STATE_1)

		assert.Nil(t, err)
		assert.NotEqual(t, taskId, -1)

		err = queries.DeleteTask(db.GetInstance().GetDB(), taskId)

		assert.Nil(t, err)

		err = queries.DeleteTask(db.GetInstance().GetDB(), taskId)

		assert.Equal(t, sql.ErrNoRows, err)
	})))
	t.Run("BasicCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		var expectedTasks []entities.Task
		expectedTasks = append(expectedTasks, GenerateTask(1))
		expectedTasks = append(expectedTasks, GenerateTask(3))

		taskIdToDelete := 2

		CreateTasksInDB(t, 3, TEST_TASK_NAME_TEMPLATE, entities.TASK_STATE_NEW)

		err := queries.DeleteTask(db.GetInstance().GetDB(), taskIdToDelete)

		assert.Nil(t, err)

		tasks, err := queries.GetTasks(db.GetInstance().GetDB(), 50, 0)

		assert.Nil(t, err)
		AssertEqualTaskArrays(t, expectedTasks, tasks)

		_, err = queries.GetTask(db.GetInstance().GetDB(), taskIdToDelete)

		assert.Equal(t, sql.ErrNoRows, err)
	})))
}
