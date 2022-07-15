package queries

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ArtemVoronov/indefinite-studies-api/internal/db"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/db/entities"
)

func GetTasks(database *sql.DB, ctx context.Context, limit int, offset int) ([]entities.Task, error) {
	var tasks []entities.Task
	var (
		id    int
		name  string
		state string
	)

	tx, err := database.BeginTx(ctx, nil)
	if err != nil {
		return tasks, fmt.Errorf("error at loading tasks from db, case after begin tx: %s", err)
	}
	defer tx.Rollback()

	rows, err := tx.QueryContext(ctx, "SELECT id, name, state FROM tasks WHERE state != $3 LIMIT $1 OFFSET $2 ", limit, offset, entities.TASK_STATE_DELETED)
	if err != nil {
		return tasks, fmt.Errorf("error at loading tasks from db, case after Query: %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &name, &state)
		if err != nil {
			return tasks, fmt.Errorf("error at loading tasks from db, case iterating and using rows.Scan: %s", err)
		}
		tasks = append(tasks, entities.Task{Id: id, Name: name, State: state})
	}
	err = rows.Err()
	if err != nil {
		return tasks, fmt.Errorf("error at loading tasks from db, case after iterating: %s", err)
	}

	err = tx.Commit()
	if err != nil {
		return tasks, fmt.Errorf("error at loading tasks from db, case after commit tx: %s", err)
	}

	return tasks, nil
}

func GetTask(database *sql.DB, ctx context.Context, id int) (entities.Task, error) {
	var task entities.Task

	tx, err := database.BeginTx(ctx, nil)
	if err != nil {
		return task, fmt.Errorf("error at loading task by id '%d' from db, case after begin tx: %s", id, err)
	}
	defer tx.Rollback()

	err = tx.QueryRowContext(ctx, "SELECT id, name, state FROM tasks WHERE id = $1 and state != $2 ", id, entities.TASK_STATE_DELETED).Scan(&task.Id, &task.Name, &task.State)
	if err != nil {
		if err == sql.ErrNoRows {
			return task, err
		}
		return task, fmt.Errorf("error at loading task by id '%d' from db, case after QueryRow.Scan: %s", id, err)
	}

	err = tx.Commit()
	if err != nil {
		return task, fmt.Errorf("error at loading task by id '%d' from db, case after commit tx: %s", id, err)
	}

	return task, nil
}

func CreateTask(database *sql.DB, ctx context.Context, name string, state string) (int, error) {
	lastInsertId := -1

	tx, err := database.BeginTx(ctx, nil)
	if err != nil {
		return lastInsertId, fmt.Errorf("error at inserting task (Name: '%s', State: '%s') into db, case after begin tx: %s", name, state, err)
	}
	defer tx.Rollback()

	err = tx.QueryRowContext(ctx, "INSERT INTO tasks(name, state) VALUES($1, $2) RETURNING id", name, state).Scan(&lastInsertId) // scan will release the connection
	if err != nil {
		if err.Error() == db.ErrorTaskDuplicateKey.Error() {
			return -1, db.ErrorTaskDuplicateKey
		}
		return -1, fmt.Errorf("error at inserting task (Name: '%s', State: '%s') into db, case after QueryRow.Scan: %s", name, state, err)
	}

	err = tx.Commit()
	if err != nil {
		return lastInsertId, fmt.Errorf("error at inserting task (Name: '%s', State: '%s') into db, case after commit tx: %s", name, state, err)
	}

	return lastInsertId, nil
}

func UpdateTask(database *sql.DB, ctx context.Context, id int, name string, state string) error {
	tx, err := database.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error at updating task (Id: %d, Name: '%s', State: '%s'), case after begin tx: %s", id, name, state, err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, "UPDATE tasks SET name = $2, state = $3 WHERE id = $1 and state != $4")
	if err != nil {
		return fmt.Errorf("error at updating task, case after preparing statement: %s", err)
	}

	res, err := stmt.ExecContext(ctx, id, name, state, entities.TASK_STATE_DELETED)
	if err != nil {
		if err.Error() == db.ErrorTaskDuplicateKey.Error() {
			return db.ErrorTaskDuplicateKey
		}
		return fmt.Errorf("error at updating task (Id: %d, Name: '%s', State: '%s'), case after executing statement: %s", id, name, state, err)
	}

	affectedRowsCount, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error at updating task (Id: %d, Name: '%s', State: '%s'), case after counting affected rows: %s", id, name, state, err)
	}

	if affectedRowsCount == 0 {
		return sql.ErrNoRows
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error at updating task (Id: %d, Name: '%s', State: '%s'), case after commit tx: %s", id, name, state, err)
	}

	return nil
}

func DeleteTask(database *sql.DB, ctx context.Context, id int) error {
	tx, err := database.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error at deleting task by id '%d', case after begin tx: %s", id, err)
	}
	defer tx.Rollback()

	// just for keeping the history we will add suffix to name and change state to 'DELETED', because of key constraint (name, state)
	stmt, err := tx.PrepareContext(ctx, "UPDATE tasks SET name = name||'_deleted_'||$1, state = $2 WHERE id = $1 and state != $2")
	if err != nil {
		return fmt.Errorf("error at deleting task, case after preparing statement: %s", err)
	}
	res, err := stmt.ExecContext(ctx, id, entities.TASK_STATE_DELETED)
	if err != nil {
		return fmt.Errorf("error at deleting task by id '%d', case after executing statement: %s", id, err)
	}

	affectedRowsCount, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error at deleting task by id '%d', case after counting affected rows: %s", id, err)
	}
	if affectedRowsCount == 0 {
		return sql.ErrNoRows
	}
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error at deleting task by id '%d', case after commit tx: %s", id, err)
	}
	return nil
}
