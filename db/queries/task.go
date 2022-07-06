package queries

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ArtemVoronov/indefinite-studies-api/db/entities"
)

//TODO: add tests for each case

func GetTasks(db *sql.DB, limit string, offset string) ([]entities.Task, error) {
	var tasks []entities.Task
	var (
		id    int
		name  string
		state string
	)

	rows, err := db.Query("SELECT id, name, state FROM tasks LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		log.Fatal(err) // TODO fix os.Exit(1) problem
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &name, &state)
		if err != nil {
			log.Fatal(err) // TODO fix os.Exit(1) problem
		}
		tasks = append(tasks, entities.Task{Id: id, Name: name, State: state})
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err) // TODO fix os.Exit(1) problem
	}

	return tasks, nil
}

func GetTask(db *sql.DB, id int) (entities.Task, error) {
	var task entities.Task

	err := db.QueryRow("SELECT id, name, state FROM tasks WHERE id = $1", id).Scan(&task.Id, &task.Name, &task.State)
	if err != nil {
		if err == sql.ErrNoRows {
			return task, err
		} else {
			log.Fatal(err) // TODO fix os.Exit(1) problem
		}
	}

	return task, nil
}

func CreateTask(db *sql.DB, name string, state string) (string, error) {
	stmt, err := db.Prepare("INSERT INTO tasks(name, state) VALUES($1, $2)")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(name, state)
	if err != nil {
		log.Fatal(err) // TODO fix os.Exit(1) problem
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err) // TODO fix os.Exit(1) problem
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err) // TODO fix os.Exit(1) problem
	}
	result := fmt.Sprintf("ID = %d, affected = %d\n", lastId, rowCnt)
	return result, nil
}

func UpdateTask(db *sql.DB, id int, name string, state string) (string, error) {
	stmt, err := db.Prepare("UPDATE tasks SET name = $2, state = $3 WHERE id = $1")
	if err != nil {
		log.Fatal(err) // TODO fix os.Exit(1) problem
	}
	res, err := stmt.Exec(id, name, state)
	if err != nil {
		log.Fatal(err) // TODO fix os.Exit(1) problem
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err) // TODO fix os.Exit(1) problem
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err) // TODO fix os.Exit(1) problem
	}
	result := fmt.Sprintf("ID = %d, affected = %d\n", lastId, rowCnt)
	return result, nil
}

func DeleteTask(db *sql.DB, id int) (string, error) {
	stmt, err := db.Prepare("UPDATE tasks SET state = $2 WHERE id = $1")
	if err != nil {
		log.Fatal(err) // TODO fix os.Exit(1) problem
	}
	res, err := stmt.Exec(id, entities.TASK_STATE_DELETED)
	if err != nil {
		log.Fatal(err) // TODO fix os.Exit(1) problem
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err) // TODO fix os.Exit(1) problem
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err) // TODO fix os.Exit(1) problem
	}
	result := fmt.Sprintf("ID = %d, affected = %d\n", lastId, rowCnt)
	return result, nil
}
