package queries

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ArtemVoronov/indefinite-studies-api/internal/db"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/db/entities"
)

func GetTags(tx *sql.Tx, ctx context.Context, limit int, offset int) ([]entities.Tag, error) {
	var tags []entities.Tag
	var (
		id    int
		name  string
		state string
	)

	rows, err := tx.QueryContext(ctx, "SELECT id, name, state FROM tags WHERE state != $3 LIMIT $1 OFFSET $2 ", limit, offset, entities.TAG_STATE_DELETED)
	if err != nil {
		return tags, fmt.Errorf("error at loading tags from db, case after Query: %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &name, &state)
		if err != nil {
			return tags, fmt.Errorf("error at loading tags from db, case iterating and using rows.Scan: %s", err)
		}
		tags = append(tags, entities.Tag{Id: id, Name: name, State: state})
	}
	err = rows.Err()
	if err != nil {
		return tags, fmt.Errorf("error at loading tags from db, case after iterating: %s", err)
	}

	return tags, nil
}

func GetTag(tx *sql.Tx, ctx context.Context, id int) (entities.Tag, error) {
	var tag entities.Tag

	err := tx.QueryRowContext(ctx, "SELECT id, name, state FROM tags WHERE id = $1 and state != $2 ", id, entities.TAG_STATE_DELETED).Scan(&tag.Id, &tag.Name, &tag.State)
	if err != nil {
		if err == sql.ErrNoRows {
			return tag, err
		}
		return tag, fmt.Errorf("error at loading tag by id '%d' from db, case after QueryRow.Scan: %s", id, err)
	}

	return tag, nil
}

func CreateTag(tx *sql.Tx, ctx context.Context, name string, state string) (int, error) {
	lastInsertId := -1

	err := tx.QueryRowContext(ctx, "INSERT INTO tags(name, state) VALUES($1, $2) RETURNING id", name, state).Scan(&lastInsertId) // scan will release the connection
	if err != nil {
		if err.Error() == db.ErrorTagDuplicateKey.Error() {
			return -1, db.ErrorTagDuplicateKey
		}
		return -1, fmt.Errorf("error at inserting tag (Name: '%s', State: '%s') into db, case after QueryRow.Scan: %s", name, state, err)
	}

	return lastInsertId, nil
}

func UpdateTag(tx *sql.Tx, ctx context.Context, id int, name string, state string) error {
	stmt, err := tx.PrepareContext(ctx, "UPDATE tags SET name = $2, state = $3 WHERE id = $1 and state != $4")
	if err != nil {
		return fmt.Errorf("error at updating tag, case after preparing statement: %s", err)
	}

	res, err := stmt.ExecContext(ctx, id, name, state, entities.TAG_STATE_DELETED)
	if err != nil {
		if err.Error() == db.ErrorTagDuplicateKey.Error() {
			return db.ErrorTagDuplicateKey
		}
		return fmt.Errorf("error at updating tag (Id: %d, Name: '%s', State: '%s'), case after executing statement: %s", id, name, state, err)
	}

	affectedRowsCount, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error at updating tag (Id: %d, Name: '%s', State: '%s'), case after counting affected rows: %s", id, name, state, err)
	}

	if affectedRowsCount == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func DeleteTag(tx *sql.Tx, ctx context.Context, id int) error {
	// just for keeping the history we will add suffix to name and change state to 'DELETED', because of key constraint (name, state)
	stmt, err := tx.PrepareContext(ctx, "UPDATE tags SET name = name||'_deleted_'||$1, state = $2 WHERE id = $1 and state != $2")
	if err != nil {
		return fmt.Errorf("error at deleting tag, case after preparing statement: %s", err)
	}

	res, err := stmt.ExecContext(ctx, id, entities.TAG_STATE_DELETED)
	if err != nil {
		return fmt.Errorf("error at deleting tag by id '%d', case after executing statement: %s", id, err)
	}

	affectedRowsCount, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error at deleting tag by id '%d', case after counting affected rows: %s", id, err)
	}

	if affectedRowsCount == 0 {
		return sql.ErrNoRows
	}

	return nil
}
