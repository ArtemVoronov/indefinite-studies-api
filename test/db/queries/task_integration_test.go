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
	TEST_TASK_NAME_1        string = "Test task 1"
	TEST_TASK_STATE_1       string = entities.TASK_STATE_NEW
	TEST_TASK_NAME_2        string = "Test task 2"
	TEST_TASK_STATE_2       string = entities.TASK_STATE_DONE
	TEST_TASK_NAME_TEMPLATE string = "Test task "
)

var TaskDuplicateKeyConstraintViolationError = fmt.Errorf(DuplicateKeyConstraintViolationError, "tasks_name_state_unique")

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
	t.Run("BasicCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedTaskId := 1

		actualTaskId, err := queries.CreateTask(db.DB, TEST_TASK_NAME_1, TEST_TASK_STATE_1)
		if err != nil || actualTaskId == -1 {
			t.Errorf("Unable to create task: %s", err)
		}

		assert.Equal(t, expectedTaskId, actualTaskId)
	})))
	t.Run("DuplicateCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedError := fmt.Errorf("error at inserting task (Name: '%s', State: '%s') into db, case after db.QueryRow.Scan: %s", TEST_TASK_NAME_1, TEST_TASK_STATE_1, TaskDuplicateKeyConstraintViolationError)

		taskId, err := queries.CreateTask(db.DB, TEST_TASK_NAME_1, TEST_TASK_STATE_1)
		if err != nil || taskId == -1 {
			t.Errorf("Unable to create task: %s", err)
		}
		_, actualError := queries.CreateTask(db.DB, TEST_TASK_NAME_1, TEST_TASK_STATE_1)

		assert.Equal(t, expectedError, actualError)
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
			taskId, err := queries.CreateTask(db.DB, TEST_TASK_NAME_TEMPLATE+strconv.Itoa(i), entities.TASK_STATE_NEW)
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
		for i, task := range tasks {
			assert.Equal(t, i+1, task.Id)
			assert.Equal(t, TEST_TASK_NAME_TEMPLATE+strconv.Itoa(i), task.Name)
			assert.Equal(t, entities.TASK_STATE_NEW, task.State)
		}
	})))
	t.Run("LimitParameterCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedArrayLength := 5
		for i := 0; i < 10; i++ {
			taskId, err := queries.CreateTask(db.DB, TEST_TASK_NAME_TEMPLATE+strconv.Itoa(i), entities.TASK_STATE_NEW)
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
		for i, task := range tasks {
			assert.Equal(t, i+1, task.Id)
			assert.Equal(t, TEST_TASK_NAME_TEMPLATE+strconv.Itoa(i), task.Name)
			assert.Equal(t, entities.TASK_STATE_NEW, task.State)
		}
	})))
	t.Run("OffsetParameterCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedArrayLength := 5
		for i := 0; i < 10; i++ {
			taskId, err := queries.CreateTask(db.DB, TEST_TASK_NAME_TEMPLATE+strconv.Itoa(i), entities.TASK_STATE_NEW)
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
			assert.Equal(t, TEST_TASK_NAME_TEMPLATE+strconv.Itoa(i+5), task.Name)
			assert.Equal(t, entities.TASK_STATE_NEW, task.State)
		}
	})))
}

func TestUpdateTask(t *testing.T) {
	t.Run("ExpectedNotFoundError", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedError := sql.ErrNoRows

		actualError := queries.UpdateTask(db.DB, 1, TEST_TASK_NAME_1, TEST_TASK_STATE_1)

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
	t.Run("DuplicateCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		taskId, err := queries.CreateTask(db.DB, TEST_TASK_NAME_1, TEST_TASK_STATE_1)
		if err != nil || taskId == -1 {
			t.Errorf("Unable to create task: %s", err)
		}
		taskId, err = queries.CreateTask(db.DB, TEST_TASK_NAME_2, TEST_TASK_STATE_2)
		if err != nil || taskId == -1 {
			t.Errorf("Unable to create task: %s", err)
		}

		expectedError := fmt.Errorf("error at updating task (Id: %d, Name: '%s', State: '%s'), case after executing statement: %s", taskId, TEST_TASK_NAME_1, TEST_TASK_STATE_1, TaskDuplicateKeyConstraintViolationError)

		actualError := queries.UpdateTask(db.DB, taskId, TEST_TASK_NAME_1, TEST_TASK_STATE_1)

		assert.Equal(t, expectedError, actualError)
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
		expectedState := entities.TASK_STATE_NEW
		expectedError := sql.ErrNoRows
		expectedArrayLength := 2
		taskIdToDelete := 2
		for i := 0; i < 3; i++ {
			taskId, err := queries.CreateTask(db.DB, TEST_TASK_NAME_TEMPLATE+strconv.Itoa(i), entities.TASK_STATE_NEW)
			if err != nil || taskId == -1 {
				t.Errorf("Unable to create task: %s", err)
			}
		}

		err := queries.DeleteTask(db.DB, taskIdToDelete)
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
		assert.Equal(t, TEST_TASK_NAME_TEMPLATE+"0", tasks[0].Name)
		assert.Equal(t, expectedState, tasks[0].State)
		assert.Equal(t, expectedSecondTaskId, tasks[1].Id)
		assert.Equal(t, TEST_TASK_NAME_TEMPLATE+"2", tasks[1].Name)
		assert.Equal(t, expectedState, tasks[1].State)

		_, actualError := queries.GetTask(db.DB, taskIdToDelete)

		assert.Equal(t, expectedError, actualError)
	})))
}
