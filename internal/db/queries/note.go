package queries

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ArtemVoronov/indefinite-studies-api/internal/db/entities"
)

func GetNotes(tx *sql.Tx, ctx context.Context, limit int, offset int) ([]entities.Note, error) {
	var note []entities.Note
	var (
		id             int
		text           string
		topic          string
		tagId          int
		userId         int
		state          string
		createDate     time.Time
		lastUpdateDate time.Time
	)

	rows, err := tx.QueryContext(ctx, "SELECT id, text, topic, tag_id, user_id, state, create_date, last_update_date FROM notes WHERE state != $3 LIMIT $1 OFFSET $2 ", limit, offset, entities.NOTE_STATE_DELETED)
	if err != nil {
		return note, fmt.Errorf("error at loading note from db, case after Query: %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &text, &topic, &tagId, &userId, &state, &createDate, &lastUpdateDate)
		if err != nil {
			return note, fmt.Errorf("error at loading notes from db, case iterating and using rows.Scan: %s", err)
		}
		note = append(note, entities.Note{Id: id, Text: text, Topic: topic, TagId: tagId, UserId: userId, State: state, CreateDate: createDate, LastUpdateDate: lastUpdateDate})
	}
	err = rows.Err()
	if err != nil {
		return note, fmt.Errorf("error at loading notes from db, case after iterating: %s", err)
	}

	return note, nil
}

func GetNote(tx *sql.Tx, ctx context.Context, id int) (entities.Note, error) {
	var note entities.Note

	err := tx.QueryRowContext(ctx, "SELECT id, text, topic, tag_id, user_id, state, create_date, last_update_date FROM notes WHERE id = $1 and state != $2 ", id, entities.NOTE_STATE_DELETED).
		Scan(&note.Id, &note.Text, &note.Topic, &note.TagId, &note.UserId, &note.State, &note.CreateDate, &note.LastUpdateDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return note, err
		} else {
			return note, fmt.Errorf("error at loading note by id '%d' from db, case after QueryRow.Scan: %s", id, err)
		}
	}

	return note, nil
}

func CreateNote(tx *sql.Tx, ctx context.Context, text string, topic string, tagId int, userId int, state string) (int, error) {
	lastInsertId := -1

	createDate := time.Now()
	lastUpdateDate := time.Now()

	err := tx.QueryRowContext(ctx, "INSERT INTO notes(text, topic, tag_id, user_id, state, create_date, last_update_date) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		text, topic, tagId, userId, state, createDate, lastUpdateDate).
		Scan(&lastInsertId) // scan will release the connection
	if err != nil {
		return -1, fmt.Errorf("error at inserting note (Topic: '%s', UserId: '%d') into db, case after QueryRow.Scan: %s", topic, userId, err)
	}

	return lastInsertId, nil
}

func UpdateNote(tx *sql.Tx, ctx context.Context, id int, text string, topic string, tagId int, userId int, state string) error {
	lastUpdateDate := time.Now()
	stmt, err := tx.PrepareContext(ctx, "UPDATE notes SET text = $2, topic = $3, tag_id = $4, user_id = $5, state = $6, last_update_date = $7 WHERE id = $1 and state != $8")
	if err != nil {
		return fmt.Errorf("error at updating note, case after preparing statement: %s", err)
	}
	res, err := stmt.ExecContext(ctx, id, text, topic, tagId, userId, state, lastUpdateDate, entities.NOTE_STATE_DELETED)
	if err != nil {
		return fmt.Errorf("error at updating note (Id: %d, Topic: '%s', UserId: '%d', State: '%s'), case after executing statement: %s", id, topic, userId, state, err)
	}

	affectedRowsCount, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error at updating note (Id: %d, Topic: '%s', UserId: '%d', State: '%s'), case after counting affected rows: %s", id, topic, userId, state, err)
	}
	if affectedRowsCount == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func DeleteNote(tx *sql.Tx, ctx context.Context, id int) error {
	stmt, err := tx.PrepareContext(ctx, "UPDATE notes SET state = $2 WHERE id = $1 and state != $2")
	if err != nil {
		return fmt.Errorf("error at deleting note, case after preparing statement: %s", err)
	}
	res, err := stmt.ExecContext(ctx, id, entities.NOTE_STATE_DELETED)
	if err != nil {
		return fmt.Errorf("error at deleting note by id '%d', case after executing statement: %s", id, err)
	}
	affectedRowsCount, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error at deleting note by id '%d', case after counting affected rows: %s", id, err)
	}
	if affectedRowsCount == 0 {
		return sql.ErrNoRows
	}
	return nil
}
