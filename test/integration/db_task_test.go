//go:build integration
// +build integration

package integration

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"testing"

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

func CreateTaskInDB(t *testing.T, tx *sql.Tx, ctx context.Context, name string, state string) (int, error) {
	taskId, err := queries.CreateTask(tx, ctx, name, state)
	assert.Nil(t, err)
	assert.NotEqual(t, taskId, -1)
	return taskId, err
}

func CreateTasksInDB(t *testing.T, tx *sql.Tx, ctx context.Context, count int, nameTemplate string, state string) error {
	var lastErr error
	for i := 1; i <= count; i++ {
		_, err := CreateTaskInDB(t, tx, ctx, GenerateTaskName(TEST_TASK_NAME_TEMPLATE, i), state)
		if err != nil {
			lastErr = err
		}
	}
	return lastErr
}

func TestDBTaskGet(t *testing.T) {
	t.Run("NotFoundCase", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			_, err := queries.GetTask(tx, ctx, 1)

			assert.Equal(t, sql.ErrNoRows, err)
			return err
		})()
	})))
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		expected := GenerateTask(1)
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			taskId, err := queries.CreateTask(tx, ctx, expected.Name, expected.State)
			assert.Nil(t, err)
			assert.Equal(t, taskId, expected.Id)
			return err
		})()
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			actual, err := queries.GetTask(tx, ctx, expected.Id)
			AssertEqualTasks(t, expected, actual)
			return err
		})()
	})))
	t.Run("TimeoutError", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			expectedError := errors.New("error at loading task by id '1' from db, case after QueryRow.Scan: context deadline exceeded")
			_, err := tx.ExecContext(ctx, "SELECT pg_sleep(10)")
			_, err = queries.GetTask(tx, ctx, 1)

			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
	t.Run("ContextCancelled", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			expectedError := errors.New("error at loading task by id '1' from db, case after QueryRow.Scan: context canceled")
			cancel()
			_, err := queries.GetTask(tx, ctx, 1)

			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
}

func TestDBTaskCreate(t *testing.T) {
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			taskId, err := queries.CreateTask(tx, ctx, TEST_TASK_NAME_1, TEST_TASK_STATE_1)

			assert.Nil(t, err)
			assert.Equal(t, taskId, 1)
			return err
		})()
	})))
	t.Run("DuplicateCase", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			taskId, err := queries.CreateTask(tx, ctx, TEST_TASK_NAME_1, TEST_TASK_STATE_1)

			assert.Nil(t, err)
			assert.NotEqual(t, taskId, -1)
			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			_, err := queries.CreateTask(tx, ctx, TEST_TASK_NAME_1, TEST_TASK_STATE_1)

			assert.Equal(t, db.ErrorTaskDuplicateKey, err)
			return err
		})()
	})))
	t.Run("TimeoutError", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			expectedError := fmt.Errorf("error at inserting task (Name: '%s', State: '%s') into db, case after QueryRow.Scan: %s", TEST_TASK_NAME_1, TEST_TASK_STATE_1, "context deadline exceeded")
			_, err := tx.ExecContext(ctx, "SELECT pg_sleep(10)")
			_, err = queries.CreateTask(tx, ctx, TEST_TASK_NAME_1, TEST_TASK_STATE_1)

			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
	t.Run("ContextCancelled", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			expectedError := fmt.Errorf("error at inserting task (Name: '%s', State: '%s') into db, case after QueryRow.Scan: %s", TEST_TASK_NAME_1, TEST_TASK_STATE_1, "context canceled")
			cancel()
			_, err := queries.CreateTask(tx, ctx, TEST_TASK_NAME_1, TEST_TASK_STATE_1)

			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
}

func TestDBTaskGetAll(t *testing.T) {
	t.Run("ExpectedEmpty", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			tasks, err := queries.GetTasks(tx, ctx, 50, 0)

			assert.Nil(t, err)
			assert.Equal(t, 0, len(tasks))
			return err
		})()
	})))
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		var expectedTasks []entities.Task
		for i := 1; i <= 10; i++ {
			expectedTasks = append(expectedTasks, GenerateTask(i))
		}
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := CreateTasksInDB(t, tx, ctx, 10, TEST_TASK_NAME_TEMPLATE, entities.TASK_STATE_NEW)
			return err
		})()
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			actualTasks, err := queries.GetTasks(tx, ctx, 50, 0)

			assert.Nil(t, err)
			AssertEqualTaskArrays(t, expectedTasks, actualTasks)
			return err
		})()
	})))
	t.Run("LimitParameterCase", RunWithRecreateDB((func(t *testing.T) {
		var expectedTasks []entities.Task
		for i := 1; i <= 5; i++ {
			expectedTasks = append(expectedTasks, GenerateTask(i))
		}
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := CreateTasksInDB(t, tx, ctx, 10, TEST_TASK_NAME_TEMPLATE, entities.TASK_STATE_NEW)
			return err
		})()
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			actualTasks, err := queries.GetTasks(tx, ctx, 5, 0)

			assert.Nil(t, err)
			AssertEqualTaskArrays(t, expectedTasks, actualTasks)
			return err
		})()
	})))
	t.Run("OffsetParameterCase", RunWithRecreateDB((func(t *testing.T) {
		var expectedTasks []entities.Task
		for i := 6; i <= 10; i++ {
			expectedTasks = append(expectedTasks, GenerateTask(i))
		}
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := CreateTasksInDB(t, tx, ctx, 10, TEST_TASK_NAME_TEMPLATE, entities.TASK_STATE_NEW)
			return err
		})()
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			actualTasks, err := queries.GetTasks(tx, ctx, 50, 5)

			assert.Nil(t, err)
			AssertEqualTaskArrays(t, expectedTasks, actualTasks)
			return err
		})()
	})))
	t.Run("TimeoutError", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			expectedError := fmt.Errorf("error at loading tasks from db, case after Query: context deadline exceeded")
			_, err := tx.ExecContext(ctx, "SELECT pg_sleep(10)")
			_, err = queries.GetTasks(tx, ctx, 50, 0)

			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
	t.Run("ContextCancelled", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			expectedError := fmt.Errorf("error at loading tasks from db, case after Query: context canceled")
			cancel()
			_, err := queries.GetTasks(tx, ctx, 50, 0)

			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
}

func TestDBTaskUpdate(t *testing.T) {
	t.Run("NotFoundCase", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := queries.UpdateTask(tx, ctx, 1, TEST_TASK_NAME_1, TEST_TASK_STATE_1)

			assert.Equal(t, sql.ErrNoRows, err)
			return err
		})()
	})))
	t.Run("DeletedCase", RunWithRecreateDB((func(t *testing.T) {
		expectedTaskId := 1
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			taskId, err := queries.CreateTask(tx, ctx, TEST_TASK_NAME_1, TEST_TASK_STATE_1)

			assert.Nil(t, err)
			assert.Equal(t, expectedTaskId, taskId)
			return err
		})()
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := queries.DeleteTask(tx, ctx, expectedTaskId)

			assert.Nil(t, err)
			return err
		})()
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := queries.UpdateTask(tx, ctx, expectedTaskId, TEST_TASK_NAME_2, TEST_TASK_STATE_2)

			assert.Equal(t, sql.ErrNoRows, err)
			return err
		})()
	})))
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		expected := GenerateTask(1)
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			taskId, err := queries.CreateTask(tx, ctx, TEST_TASK_NAME_2, TEST_TASK_STATE_2)

			assert.Nil(t, err)
			assert.Equal(t, expected.Id, taskId)
			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := queries.UpdateTask(tx, ctx, expected.Id, expected.Name, expected.State)

			assert.Nil(t, err)
			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			actual, err := queries.GetTask(tx, ctx, expected.Id)

			AssertEqualTasks(t, expected, actual)
			return err
		})()
	})))
	t.Run("DuplicateCase", RunWithRecreateDB((func(t *testing.T) {
		expectedTaskId1 := 1
		expectedTaskId2 := 2
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			taskId, err := queries.CreateTask(tx, ctx, TEST_TASK_NAME_1, TEST_TASK_STATE_1)

			assert.Nil(t, err)
			assert.Equal(t, expectedTaskId1, taskId)
			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			taskId, err := queries.CreateTask(tx, ctx, TEST_TASK_NAME_2, TEST_TASK_STATE_2)

			assert.Nil(t, err)
			assert.Equal(t, expectedTaskId2, taskId)
			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := queries.UpdateTask(tx, ctx, expectedTaskId2, TEST_TASK_NAME_1, TEST_TASK_STATE_1)

			assert.Equal(t, db.ErrorTaskDuplicateKey, err)
			return err
		})()
	})))
	t.Run("TimeoutError", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			expectedError := fmt.Errorf("error at updating task, case after preparing statement: %s", "context deadline exceeded")
			_, err := tx.ExecContext(ctx, "SELECT pg_sleep(10)")
			err = queries.UpdateTask(tx, ctx, 1, TEST_TASK_NAME_1, TEST_TASK_STATE_1)

			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
	t.Run("ContextCancelled", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			expectedError := fmt.Errorf("error at updating task, case after preparing statement: %s", "context canceled")
			cancel()
			err := queries.UpdateTask(tx, ctx, 1, TEST_TASK_NAME_1, TEST_TASK_STATE_1)
			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
}

func TestDBTaskDelete(t *testing.T) {
	t.Run("NotFoundCase", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := queries.DeleteTask(tx, ctx, 1)

			assert.Equal(t, sql.ErrNoRows, err)
			return err
		})()
	})))
	t.Run("AlreadyDeletedCase", RunWithRecreateDB((func(t *testing.T) {
		expectedTaskId := 1
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			taskId, err := queries.CreateTask(tx, ctx, TEST_TASK_NAME_1, TEST_TASK_STATE_1)

			assert.Nil(t, err)
			assert.Equal(t, expectedTaskId, taskId)
			return err
		})()
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := queries.DeleteTask(tx, ctx, expectedTaskId)

			assert.Nil(t, err)
			return err
		})()
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := queries.DeleteTask(tx, ctx, expectedTaskId)

			assert.Equal(t, sql.ErrNoRows, err)
			return err
		})()
	})))
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		var expectedTasks []entities.Task
		expectedTasks = append(expectedTasks, GenerateTask(1))
		expectedTasks = append(expectedTasks, GenerateTask(3))

		taskIdToDelete := 2
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {

			err := CreateTasksInDB(t, tx, ctx, 3, TEST_TASK_NAME_TEMPLATE, entities.TASK_STATE_NEW)
			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			err := queries.DeleteTask(tx, ctx, taskIdToDelete)

			assert.Nil(t, err)
			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			tasks, err := queries.GetTasks(tx, ctx, 50, 0)

			assert.Nil(t, err)
			AssertEqualTaskArrays(t, expectedTasks, tasks)
			return err
		})()

		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			_, err := queries.GetTask(tx, ctx, taskIdToDelete)

			assert.Equal(t, sql.ErrNoRows, err)
			return err
		})()
	})))
	t.Run("TimeoutError", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			taskId := 1
			expectedError := fmt.Errorf("error at deleting task, case after preparing statement: %s", "context deadline exceeded")
			_, err := tx.ExecContext(ctx, "SELECT pg_sleep(10)")
			err = queries.DeleteTask(tx, ctx, taskId)

			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
	t.Run("ContextCancelled", RunWithRecreateDB((func(t *testing.T) {
		db.TxVoid(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) error {
			taskId := 1
			expectedError := fmt.Errorf("error at deleting task, case after preparing statement: %s", "context canceled")
			cancel()
			err := queries.DeleteTask(tx, ctx, taskId)
			assert.Equal(t, expectedError, err)
			return err
		})()
	})))
}
