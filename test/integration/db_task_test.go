//go:build integration
// +build integration

package integration

import (
	"context"
	"database/sql"
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

func CreateTaskInDB(t *testing.T, database *sql.DB, ctx context.Context, name string, state string) int {
	taskId, err := queries.CreateTask(database, ctx, name, state)
	assert.Nil(t, err)
	assert.NotEqual(t, taskId, -1)
	return taskId
}

func CreateTasksInDB(t *testing.T, database *sql.DB, ctx context.Context, count int, nameTemplate string, state string) {
	for i := 1; i <= count; i++ {
		CreateTaskInDB(t, database, ctx, GenerateTaskName(TEST_TASK_NAME_TEMPLATE, i), state)
	}
}

func TestDBTaskGet(t *testing.T) {
	t.Run("NotFoundCase", RunWithRecreateDB((func(t *testing.T) {
		db.RunWithWithTimeout(func(database *sql.DB, ctx context.Context) {
			_, err := queries.GetTask(database, ctx, 1)

			assert.Equal(t, sql.ErrNoRows, err)
		})()
	})))
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		db.RunWithWithTimeout(func(database *sql.DB, ctx context.Context) {
			expected := GenerateTask(1)

			taskId, err := queries.CreateTask(database, ctx, expected.Name, expected.State)

			assert.Nil(t, err)
			assert.Equal(t, taskId, expected.Id)

			actual, err := queries.GetTask(database, ctx, taskId)

			AssertEqualTasks(t, expected, actual)
		})()
	})))
	t.Run("TimeoutErrorExpected", RunWithRecreateDB((func(t *testing.T) {
		t.Fatalf("not implemented")
		// TODO: add configurable timeout for db queries
		// db.RunWithWithTimeout(func(database *sql.DB, ctx context.Context) {
		//	_, err = tx.ExecContext(ctx, "SELECT pg_sleep(10)") // todo clean
		// })()
	})))
}

func TestDBTaskCreate(t *testing.T) {
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		db.RunWithWithTimeout(func(database *sql.DB, ctx context.Context) {
			taskId, err := queries.CreateTask(database, ctx, TEST_TASK_NAME_1, TEST_TASK_STATE_1)

			assert.Nil(t, err)
			assert.Equal(t, taskId, 1)
		})()
	})))
	t.Run("DuplicateCase", RunWithRecreateDB((func(t *testing.T) {
		db.RunWithWithTimeout(func(database *sql.DB, ctx context.Context) {
			taskId, err := queries.CreateTask(database, ctx, TEST_TASK_NAME_1, TEST_TASK_STATE_1)

			assert.Nil(t, err)
			assert.NotEqual(t, taskId, -1)

			_, err = queries.CreateTask(database, ctx, TEST_TASK_NAME_1, TEST_TASK_STATE_1)

			assert.Equal(t, db.ErrorTaskDuplicateKey, err)
		})()
	})))
	t.Run("TimeoutErrorExpected", RunWithRecreateDB((func(t *testing.T) {
		t.Fatalf("not implemented")
		// TODO: add configurable timeout for db queries
		// db.RunWithWithTimeout(func(database *sql.DB, ctx context.Context) {
		//	_, err = tx.ExecContext(ctx, "SELECT pg_sleep(10)") // todo clean
		// })()
	})))
}

func TestDBTaskGetAll(t *testing.T) {
	t.Run("ExpectedEmpty", RunWithRecreateDB((func(t *testing.T) {
		db.RunWithWithTimeout(func(database *sql.DB, ctx context.Context) {
			tasks, err := queries.GetTasks(database, ctx, 50, 0)

			assert.Nil(t, err)
			assert.Equal(t, 0, len(tasks))
		})()
	})))
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		db.RunWithWithTimeout(func(database *sql.DB, ctx context.Context) {
			var expectedTasks []entities.Task
			for i := 1; i <= 10; i++ {
				expectedTasks = append(expectedTasks, GenerateTask(i))
			}
			CreateTasksInDB(t, database, ctx, 10, TEST_TASK_NAME_TEMPLATE, entities.TASK_STATE_NEW)

			actualTasks, err := queries.GetTasks(database, ctx, 50, 0)

			assert.Nil(t, err)
			AssertEqualTaskArrays(t, expectedTasks, actualTasks)
		})()
	})))
	t.Run("LimitParameterCase", RunWithRecreateDB((func(t *testing.T) {
		db.RunWithWithTimeout(func(database *sql.DB, ctx context.Context) {
			var expectedTasks []entities.Task
			for i := 1; i <= 5; i++ {
				expectedTasks = append(expectedTasks, GenerateTask(i))
			}

			CreateTasksInDB(t, database, ctx, 10, TEST_TASK_NAME_TEMPLATE, entities.TASK_STATE_NEW)

			actualTasks, err := queries.GetTasks(database, ctx, 5, 0)

			assert.Nil(t, err)
			AssertEqualTaskArrays(t, expectedTasks, actualTasks)
		})()
	})))
	t.Run("OffsetParameterCase", RunWithRecreateDB((func(t *testing.T) {
		db.RunWithWithTimeout(func(database *sql.DB, ctx context.Context) {
			var expectedTasks []entities.Task
			for i := 6; i <= 10; i++ {
				expectedTasks = append(expectedTasks, GenerateTask(i))
			}

			CreateTasksInDB(t, database, ctx, 10, TEST_TASK_NAME_TEMPLATE, entities.TASK_STATE_NEW)

			actualTasks, err := queries.GetTasks(database, ctx, 50, 5)

			assert.Nil(t, err)
			AssertEqualTaskArrays(t, expectedTasks, actualTasks)
		})()
	})))
	t.Run("TimeoutErrorExpected", RunWithRecreateDB((func(t *testing.T) {
		t.Fatalf("not implemented")
		// TODO: add configurable timeout for db queries
		// db.RunWithWithTimeout(func(database *sql.DB, ctx context.Context) {
		//	_, err = tx.ExecContext(ctx, "SELECT pg_sleep(10)") // todo clean
		// })()
	})))

}

func TestDBTaskUpdate(t *testing.T) {
	t.Run("NotFoundCase", RunWithRecreateDB((func(t *testing.T) {
		db.RunWithWithTimeout(func(database *sql.DB, ctx context.Context) {
			err := queries.UpdateTask(database, ctx, 1, TEST_TASK_NAME_1, TEST_TASK_STATE_1)

			assert.Equal(t, sql.ErrNoRows, err)
		})()
	})))
	t.Run("DeletedCase", RunWithRecreateDB((func(t *testing.T) {
		db.RunWithWithTimeout(func(database *sql.DB, ctx context.Context) {
			taskId, err := queries.CreateTask(database, ctx, TEST_TASK_NAME_1, TEST_TASK_STATE_1)

			assert.Nil(t, err)
			assert.NotEqual(t, taskId, -1)

			err = queries.DeleteTask(database, ctx, taskId)

			assert.Nil(t, err)

			err = queries.UpdateTask(database, ctx, taskId, TEST_TASK_NAME_2, TEST_TASK_STATE_2)

			assert.Equal(t, sql.ErrNoRows, err)
		})()
	})))
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		db.RunWithWithTimeout(func(database *sql.DB, ctx context.Context) {
			expected := GenerateTask(1)

			taskId, err := queries.CreateTask(database, ctx, TEST_TASK_NAME_2, TEST_TASK_STATE_2)

			assert.Nil(t, err)
			assert.Equal(t, expected.Id, taskId)

			err = queries.UpdateTask(database, ctx, expected.Id, expected.Name, expected.State)

			assert.Nil(t, err)

			actual, err := queries.GetTask(database, ctx, expected.Id)

			AssertEqualTasks(t, expected, actual)
		})()
	})))
	t.Run("DuplicateCase", RunWithRecreateDB((func(t *testing.T) {
		db.RunWithWithTimeout(func(database *sql.DB, ctx context.Context) {
			taskId, err := queries.CreateTask(database, ctx, TEST_TASK_NAME_1, TEST_TASK_STATE_1)

			assert.Nil(t, err)
			assert.NotEqual(t, taskId, -1)

			taskId, err = queries.CreateTask(database, ctx, TEST_TASK_NAME_2, TEST_TASK_STATE_2)

			assert.Nil(t, err)
			assert.NotEqual(t, taskId, -1)

			actualError := queries.UpdateTask(database, ctx, taskId, TEST_TASK_NAME_1, TEST_TASK_STATE_1)

			assert.Equal(t, db.ErrorTaskDuplicateKey, actualError)
		})()
	})))
	t.Run("TimeoutErrorExpected", RunWithRecreateDB((func(t *testing.T) {
		t.Fatalf("not implemented")
		// TODO: add configurable timeout for db queries
		// db.RunWithWithTimeout(func(database *sql.DB, ctx context.Context) {
		//	_, err = tx.ExecContext(ctx, "SELECT pg_sleep(10)") // todo clean
		// })()
	})))
}

func TestDBTaskDelete(t *testing.T) {
	t.Run("NotFoundCase", RunWithRecreateDB((func(t *testing.T) {
		db.RunWithWithTimeout(func(database *sql.DB, ctx context.Context) {
			err := queries.DeleteTask(database, ctx, 1)

			assert.Equal(t, sql.ErrNoRows, err)
		})()
	})))
	t.Run("AlreadyDeletedCase", RunWithRecreateDB((func(t *testing.T) {
		db.RunWithWithTimeout(func(database *sql.DB, ctx context.Context) {
			taskId, err := queries.CreateTask(database, ctx, TEST_TASK_NAME_1, TEST_TASK_STATE_1)

			assert.Nil(t, err)
			assert.NotEqual(t, taskId, -1)

			err = queries.DeleteTask(database, ctx, taskId)

			assert.Nil(t, err)

			err = queries.DeleteTask(database, ctx, taskId)

			assert.Equal(t, sql.ErrNoRows, err)
		})()
	})))
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		db.RunWithWithTimeout(func(database *sql.DB, ctx context.Context) {
			var expectedTasks []entities.Task
			expectedTasks = append(expectedTasks, GenerateTask(1))
			expectedTasks = append(expectedTasks, GenerateTask(3))

			taskIdToDelete := 2

			CreateTasksInDB(t, database, ctx, 3, TEST_TASK_NAME_TEMPLATE, entities.TASK_STATE_NEW)

			err := queries.DeleteTask(database, ctx, taskIdToDelete)

			assert.Nil(t, err)

			tasks, err := queries.GetTasks(database, ctx, 50, 0)

			assert.Nil(t, err)
			AssertEqualTaskArrays(t, expectedTasks, tasks)

			_, err = queries.GetTask(database, ctx, taskIdToDelete)

			assert.Equal(t, sql.ErrNoRows, err)
		})()
	})))
	t.Run("TimeoutErrorExpected", RunWithRecreateDB((func(t *testing.T) {
		t.Fatalf("not implemented")
		// TODO: add configurable timeout for db queries
		// db.RunWithWithTimeout(func(database *sql.DB, ctx context.Context) {
		//	_, err = tx.ExecContext(ctx, "SELECT pg_sleep(10)") // todo clean
		// })()
	})))
}
