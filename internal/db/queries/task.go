package queries

import (
	"database/sql"
	"fmt"

	"github.com/ArtemVoronov/indefinite-studies-api/internal/db/entities"
)

//TODO: add tests for each case

func GetTasks(db *sql.DB, limit string, offset string) ([]entities.Task, error) {
	var tasks []entities.Task
	var (
		id    int
		name  string
		state string
	)

	rows, err := db.Query("SELECT id, name, state FROM tasks WHERE state != $3 LIMIT $1 OFFSET $2 ", limit, offset, entities.TASK_STATE_DELETED)
	if err != nil {
		return tasks, fmt.Errorf("error at loading tasks from db, case after db.Query: %s", err)
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

	return tasks, nil
}

func GetTask(db *sql.DB, id int) (entities.Task, error) {
	var task entities.Task

	err := db.QueryRow("SELECT id, name, state FROM tasks WHERE id = $1 and state != $2 ", id, entities.TASK_STATE_DELETED).Scan(&task.Id, &task.Name, &task.State)
	if err != nil {
		if err == sql.ErrNoRows {
			return task, err
		} else {
			return task, fmt.Errorf("error at loading task by id '%d' from db, case after db.QueryRow.Scan: %s", id, err)
		}
	}

	return task, nil
}

func CreateTask(db *sql.DB, name string, state string) (int, error) {
	lastInsertId := -1

	err := db.QueryRow("INSERT INTO tasks(name, state) VALUES($1, $2) RETURNING id", name, state).Scan(&lastInsertId) // scan will release the connection
	if err != nil {
		return -1, fmt.Errorf("error at inserting task (Name: '%s', State: '%s') into db, case after db.QueryRow.Scan: %s", name, state, err)
	}

	return lastInsertId, nil
}

func UpdateTask(db *sql.DB, id int, name string, state string) error {
	stmt, err := db.Prepare("UPDATE tasks SET name = $2, state = $3 WHERE id = $1 and state != $4")
	if err != nil {
		return fmt.Errorf("error at updating task, case after preparing statement: %s", err)
	}
	res, err := stmt.Exec(id, name, state, entities.TASK_STATE_DELETED)
	if err != nil {
		return fmt.Errorf("error at updating task (Id: %d, Name: '%s', State: '%s'), case after executing statement: %s", id, name, state, err)
	}

	affectedRowsCount, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error at updating task (Id: %d, Name: '%s', State: '%s'), case after counting affected rows: %s", id, name, state, err)
	}
	if affectedRowsCount == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func DeleteTask(db *sql.DB, id int) error {
	stmt, err := db.Prepare("UPDATE tasks SET state = $2 WHERE id = $1 and state != $3")
	if err != nil {
		return fmt.Errorf("error at deleting task, case after preparing statement: %s", err)
	}
	_, err = stmt.Exec(id, entities.TASK_STATE_DELETED, entities.TASK_STATE_DELETED)
	if err != nil {
		return fmt.Errorf("error at deleting task by id '%d', case after executing statement: %s", id, err)
	}
	return nil
}
