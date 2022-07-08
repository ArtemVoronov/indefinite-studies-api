//go:build integration
// +build integration

package queries_test

import (
	"database/sql"
	integrationTesting "github.com/ArtemVoronov/indefinite-studies-api/internal/app/testing"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/db"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/db/entities"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/db/queries"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const (
	TEST_TASK_NAME_1  string = "Test task 1"
	TEST_TASK_STATE_1 string = entities.TASK_STATE_NEW
	TEST_TASK_NAME_2  string = "Test task 2"
	TEST_TASK_STATE_2 string = entities.TASK_STATE_DONE
)

func TestMain(m *testing.M) {
	integrationTesting.Setup()
	code := m.Run()
	integrationTesting.Shutdown()
	os.Exit(code)

}

func TestGetTask(t *testing.T) {
	t.Run("ExpectedNotFoundError", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedError := sql.ErrNoRows

		_, actualError := queries.GetTask(db.DB, 1)

		assert.Equal(t, expectedError, actualError)

	})))
	t.Run("ExpectedResult", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedName := TEST_TASK_NAME_1
		expectedState := TEST_TASK_STATE_1
		expectedId, err := queries.CreateTask(db.DB, expectedName, expectedState)

		if err != nil || expectedId == -1 {
			t.Errorf("Unable to create task: %s", err)
		}

		actual, err := queries.GetTask(db.DB, expectedId)

		assert.Equal(t, expectedId, actual.Id)
		assert.Equal(t, expectedName, actual.Name)
		assert.Equal(t, expectedState, actual.State)

	})))
}

func TestCreateTask(t *testing.T) {
	t.Run("", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedFirstTaskId := 1
		expectedSecondTaskId := 2

		actualFirstTaskId, err := queries.CreateTask(db.DB, TEST_TASK_NAME_1, TEST_TASK_STATE_1)
		if err != nil || actualFirstTaskId == -1 {
			t.Errorf("Unable to create task: %s", err)
		}
		actualSecondTaskId, err := queries.CreateTask(db.DB, TEST_TASK_NAME_1, TEST_TASK_STATE_1)
		if err != nil || actualSecondTaskId == -1 {
			t.Errorf("Unable to create task: %s", err)
		}

		assert.Equal(t, expectedFirstTaskId, actualFirstTaskId)
		assert.Equal(t, expectedSecondTaskId, actualSecondTaskId)
	})))
}

func TestGetTasks(t *testing.T) {
	t.Run("ExpectedEmpty", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedArrayLength := 0

		tasks, err := queries.GetTasks(db.DB, "50", "0")
		if err != nil {
			t.Errorf("Unable to get to tasks : %s", err)
		}
		actualArrayLength := len(tasks)

		assert.Equal(t, expectedArrayLength, actualArrayLength)
	})))
	t.Run("ExpectedResult", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedArrayLength := 3

		for i := 0; i < 3; i++ {
			taskId, err := queries.CreateTask(db.DB, TEST_TASK_NAME_1, TEST_TASK_STATE_1)
			if err != nil || taskId == -1 {
				t.Errorf("Unable to create task: %s", err)
			}
		}
		tasks, err := queries.GetTasks(db.DB, "50", "0")
		if err != nil {
			t.Errorf("Unable to get to tasks : %s", err)
		}
		actualArrayLength := len(tasks)

		assert.Equal(t, expectedArrayLength, actualArrayLength)
	})))
	t.Run("LimitParameterCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedArrayLength := 5
		for i := 0; i < 10; i++ {
			taskId, err := queries.CreateTask(db.DB, TEST_TASK_NAME_1, TEST_TASK_STATE_1)
			if err != nil || taskId == -1 {
				t.Errorf("Unable to create task: %s", err)
			}
		}

		tasks, err := queries.GetTasks(db.DB, "5", "0")
		if err != nil {
			t.Errorf("Unable to get to tasks : %s", err)
		}
		actualArrayLength := len(tasks)

		assert.Equal(t, expectedArrayLength, actualArrayLength)
	})))
	t.Run("OffsetParameterCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedName := TEST_TASK_NAME_1
		expectedState := TEST_TASK_STATE_1
		expectedArrayLength := 5
		for i := 0; i < 10; i++ {
			taskId, err := queries.CreateTask(db.DB, TEST_TASK_NAME_1, TEST_TASK_STATE_1)
			if err != nil || taskId == -1 {
				t.Errorf("Unable to create task: %s", err)
			}
		}

		tasks, err := queries.GetTasks(db.DB, "50", "5")
		if err != nil {
			t.Errorf("Unable to get to tasks : %s", err)
		}
		actualArrayLength := len(tasks)

		assert.Equal(t, expectedArrayLength, actualArrayLength)
		for i, task := range tasks {
			assert.Equal(t, i+6, task.Id)
			assert.Equal(t, expectedName, task.Name)
			assert.Equal(t, expectedState, task.State)
		}
	})))
}

func TestUpdateTask(t *testing.T) {
	t.Run("ExpectedNotFoundError", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedError := sql.ErrNoRows

		actualError := queries.UpdateTask(db.DB, 1, TEST_TASK_NAME_2, TEST_TASK_STATE_2)

		assert.Equal(t, expectedError, actualError)
	})))
	t.Run("DeletedCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedError := sql.ErrNoRows

		taskId, err := queries.CreateTask(db.DB, TEST_TASK_NAME_1, TEST_TASK_STATE_1)
		if err != nil || taskId == -1 {
			t.Errorf("Unable to create task: %s", err)
		}

		err = queries.DeleteTask(db.DB, taskId)
		if err != nil {
			t.Errorf("Unable to delete task: %s", err)
		}

		actualError := queries.UpdateTask(db.DB, taskId, TEST_TASK_NAME_2, TEST_TASK_STATE_2)

		assert.Equal(t, expectedError, actualError)

	})))
	t.Run("BasicCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedName := TEST_TASK_NAME_2
		expectedState := TEST_TASK_STATE_2
		expectedId, err := queries.CreateTask(db.DB, TEST_TASK_NAME_1, TEST_TASK_STATE_1)
		if err != nil || expectedId == -1 {
			t.Errorf("Unable to create task: %s", err)
		}

		err = queries.UpdateTask(db.DB, expectedId, TEST_TASK_NAME_2, TEST_TASK_STATE_2)
		if err != nil {
			t.Errorf("Unable to update task: %s", err)
		}

		actual, err := queries.GetTask(db.DB, expectedId)

		assert.Equal(t, expectedId, actual.Id)
		assert.Equal(t, expectedName, actual.Name)
		assert.Equal(t, expectedState, actual.State)
	})))
}

func TestDeleteTask(t *testing.T) {
	t.Run("NotFoundCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		notExistentTaskId := 1

		actualError := queries.DeleteTask(db.DB, notExistentTaskId)
		if actualError != nil {
			t.Errorf("Unable to delete task: %s", actualError)
		}
	})))
	t.Run("AlreadyDeletedCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {

		taskId, err := queries.CreateTask(db.DB, TEST_TASK_NAME_1, TEST_TASK_STATE_1)
		if err != nil || taskId == -1 {
			t.Errorf("Unable to create task: %s", err)
		}

		err = queries.DeleteTask(db.DB, taskId)
		if err != nil {
			t.Errorf("Unable to delete task: %s", err)
		}

		actualError := queries.DeleteTask(db.DB, taskId)
		if actualError != nil {
			t.Errorf("Unable to delete task: %s", actualError)
		}
	})))
	t.Run("BasicCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedFirstTaskId := 1
		expectedSecondTaskId := 3
		expectedName := TEST_TASK_NAME_1
		expectedState := TEST_TASK_STATE_1
		expectedError := sql.ErrNoRows
		expectedArrayLength := 2
		for i := 0; i < 3; i++ {
			taskId, err := queries.CreateTask(db.DB, TEST_TASK_NAME_1, TEST_TASK_STATE_1)
			if err != nil || taskId == -1 {
				t.Errorf("Unable to create task: %s", err)
			}
		}

		err := queries.DeleteTask(db.DB, 2)
		if err != nil {
			t.Errorf("Unable to delete task: %s", err)
		}

		tasks, err := queries.GetTasks(db.DB, "50", "0")
		if err != nil {
			t.Errorf("Unable to get to tasks : %s", err)
		}
		actualArrayLength := len(tasks)

		assert.Equal(t, expectedArrayLength, actualArrayLength)

		assert.Equal(t, expectedFirstTaskId, tasks[0].Id)
		assert.Equal(t, expectedName, tasks[0].Name)
		assert.Equal(t, expectedState, tasks[0].State)
		assert.Equal(t, expectedSecondTaskId, tasks[1].Id)
		assert.Equal(t, expectedName, tasks[1].Name)
		assert.Equal(t, expectedState, tasks[1].State)

		_, actualError := queries.GetTask(db.DB, 2)

		assert.Equal(t, expectedError, actualError)
	})))
}
