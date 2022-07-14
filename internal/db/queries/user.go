package queries

import (
	"database/sql"
	"fmt"
	"time"

	DBService "github.com/ArtemVoronov/indefinite-studies-api/internal/db"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/db/entities"
)

func GetUsers(db *sql.DB, limit int, offset int) ([]entities.User, error) {
	var user []entities.User
	var (
		id             int
		login          string
		email          string
		password       string
		role           string
		state          string
		createDate     time.Time
		lastUpdateDate time.Time
	)

	rows, err := db.Query("SELECT id, login, email, password, role, state, create_date, last_update_date FROM users WHERE state != $3 LIMIT $1 OFFSET $2 ", limit, offset, entities.USER_STATE_DELETED)
	if err != nil {
		return user, fmt.Errorf("error at loading users from db, case after db.Query: %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &login, &email, &password, &role, &state, &createDate, &lastUpdateDate)
		if err != nil {
			return user, fmt.Errorf("error at loading users from db, case iterating and using rows.Scan: %s", err)
		}
		user = append(user, entities.User{Id: id, Login: login, Email: email, Password: password, Role: role, State: state, CreateDate: createDate, LastUpdateDate: lastUpdateDate})
	}
	err = rows.Err()
	if err != nil {
		return user, fmt.Errorf("error at loading user from db, case after iterating: %s", err)
	}

	return user, nil
}

func GetUser(db *sql.DB, id int) (entities.User, error) {
	var user entities.User

	err := db.QueryRow("SELECT id, login, email, password, role, state, create_date, last_update_date FROM users WHERE id = $1 and state != $2 ", id, entities.USER_STATE_DELETED).
		Scan(&user.Id, &user.Login, &user.Email, &user.Password, &user.Role, &user.State, &user.CreateDate, &user.LastUpdateDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, err
		} else {
			return user, fmt.Errorf("error at loading user by id '%d' from db, case after db.QueryRow.Scan: %s", id, err)
		}
	}

	return user, nil
}

func CreateUser(db *sql.DB, login string, email string, password string, role string, state string) (int, error) {
	lastInsertId := -1

	createDate := time.Now()
	lastUpdateDate := time.Now()

	err := db.QueryRow("INSERT INTO users(login, email, password, role, state, create_date, last_update_date) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		login, email, password, role, state, createDate, lastUpdateDate).
		Scan(&lastInsertId) // scan will release the connection
	if err != nil {
		if err.Error() == DBService.ErrorUserDuplicateKey.Error() {
			return -1, DBService.ErrorUserDuplicateKey
		}
		return -1, fmt.Errorf("error at inserting user (Login: '%s', Email: '%s') into db, case after db.QueryRow.Scan: %s", login, email, err)
	}

	return lastInsertId, nil
}

func UpdateUser(db *sql.DB, id int, login string, email string, password string, role string, state string) error {
	lastUpdateDate := time.Now()
	stmt, err := db.Prepare("UPDATE users SET login = $2, email = $3, password = $4, role = $5, state = $6, last_update_date = $7 WHERE id = $1 and state != $8")
	if err != nil {
		return fmt.Errorf("error at updating user, case after preparing statement: %s", err)
	}
	res, err := stmt.Exec(id, login, email, password, role, state, lastUpdateDate, entities.USER_STATE_DELETED)
	if err != nil {
		if err.Error() == DBService.ErrorUserDuplicateKey.Error() {
			return DBService.ErrorUserDuplicateKey
		}
		return fmt.Errorf("error at updating user (Id: %d, Login: '%s', Email: '%s', State: '%s'), case after executing statement: %s", id, login, email, state, err)
	}

	affectedRowsCount, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error at updating user (Id: %d, Login: '%s', Email: '%s', State: '%s'), case after counting affected rows: %s", id, login, email, state, err)
	}
	if affectedRowsCount == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func DeleteUser(db *sql.DB, id int) error {
	// just for keeping the history we will add suffix to name and change state to 'DELETED', because of key constraint (email, state)
	stmt, err := db.Prepare("UPDATE users SET email = email||'_deleted_'||$1, state = $2 WHERE id = $1 and state != $2")
	if err != nil {
		return fmt.Errorf("error at deleting user, case after preparing statement: %s", err)
	}
	res, err := stmt.Exec(id, entities.USER_STATE_DELETED)
	if err != nil {
		return fmt.Errorf("error at deleting user by id '%d', case after executing statement: %s", id, err)
	}
	affectedRowsCount, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error at deleting user by id '%d', case after counting affected rows: %s", id, err)
	}
	if affectedRowsCount == 0 {
		return sql.ErrNoRows
	}
	return nil
}
