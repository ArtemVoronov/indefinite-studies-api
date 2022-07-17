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
	TEST_NOTE_TEXT_1  string = "Test text 1"
	TEST_NOTE_TOPIC_1 string = "Test topic 1"
	TEST_NOTE_STATE_1 string = entities.NOTE_STATE_NEW
	TEST_NOTE_TEXT_2  string = "Test text 2"
	TEST_NOTE_TOPIC_2 string = "Test topic 2"
	TEST_NOTE_STATE_2 string = entities.NOTE_STATE_BLOCKED

	TEST_NOTE_TEXT_TEMPLATE  string = "Test text "
	TEST_NOTE_TOPIC_TEMPLATE string = "Test topic "
)

func GenerateNoteText(template string, id int) string {
	return template + strconv.Itoa(id)
}

func GenerateNoteTopic(template string, id int) string {
	return template + strconv.Itoa(id)
}

func GenerateNote(noteId int, userId int, tagId int) entities.Note {
	return entities.Note{
		Id:     noteId,
		Text:   GenerateNoteText(TEST_NOTE_TEXT_TEMPLATE, noteId),
		Topic:  GenerateNoteTopic(TEST_NOTE_TOPIC_TEMPLATE, noteId),
		TagId:  tagId,
		UserId: userId,
		State:  TEST_USER_STATE_1,
	}
}

func AssertEqualNotes(t *testing.T, expected entities.Note, actual entities.Note) {
	assert.Equal(t, expected.Id, actual.Id)
	assert.Equal(t, expected.Text, actual.Text)
	assert.Equal(t, expected.Topic, actual.Topic)
	assert.Equal(t, expected.TagId, actual.TagId)
	assert.Equal(t, expected.UserId, actual.UserId)
	assert.Equal(t, expected.State, actual.State)
}

func AssertEqualNoteArrays(t *testing.T, expected []entities.Note, actual []entities.Note) {
	assert.Equal(t, len(expected), len(actual))

	length := len(expected)
	for i := 0; i < length; i++ {
		AssertEqualNotes(t, expected[i], actual[i])
	}
}

func CreateNoteInDB(t *testing.T, text string, topic string, tagId int, userId int, state string) int {
	noteId, err := queries.CreateNote(db.GetInstance().GetDB(), text, topic, tagId, userId, state)
	assert.Nil(t, err)
	assert.NotEqual(t, noteId, -1)
	return noteId
}

func CreateNotesInDB(t *testing.T, count int, textTemplate string, topicTemplate string, tagId int, userId int, state string) {
	for i := 1; i <= count; i++ {
		CreateNoteInDB(t, GenerateNoteText(textTemplate, i), GenerateNoteTopic(topicTemplate, i), tagId, userId, state)
	}
}

func TestDBNoteGet(t *testing.T) {
	t.Run("NotFoundCase", RunWithRecreateDB((func(t *testing.T) {
		_, actualError := queries.GetNote(db.GetInstance().GetDB(), 1)

		assert.Equal(t, sql.ErrNoRows, actualError)
	})))
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		userId := CreateUserInDB(t, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)

		result, err := db.Tx(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) (any, error) {
			tagId, err := CreateTagInDB(t, tx, ctx, TEST_TAG_NAME_1, TEST_TAG_STATE_1)
			return tagId, err
		})()
		tagId, ok := result.(int)
		assert.True(t, ok)
		expected := GenerateNote(1, userId, tagId)

		noteId, err := queries.CreateNote(db.GetInstance().GetDB(), expected.Text, expected.Topic, expected.TagId, expected.UserId, expected.State)

		assert.Nil(t, err)
		assert.Equal(t, expected.Id, noteId)

		actual, err := queries.GetNote(db.GetInstance().GetDB(), noteId)

		AssertEqualNotes(t, expected, actual)
	})))
}

func TestDBNoteCreate(t *testing.T) {
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		userId := CreateUserInDB(t, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)
		result, err := db.Tx(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) (any, error) {
			tagId, err := CreateTagInDB(t, tx, ctx, TEST_TAG_NAME_1, TEST_TAG_STATE_1)
			return tagId, err
		})()
		tagId, ok := result.(int)
		assert.True(t, ok)

		expected := GenerateNote(1, userId, tagId)

		noteId, err := queries.CreateNote(db.GetInstance().GetDB(), expected.Text, expected.Topic, expected.TagId, expected.UserId, expected.State)

		assert.Nil(t, err)
		assert.Equal(t, expected.Id, noteId)
	})))
}

func TestDBNoteGetAll(t *testing.T) {
	t.Run("ExpectedEmpty", RunWithRecreateDB((func(t *testing.T) {
		notes, err := queries.GetNotes(db.GetInstance().GetDB(), 50, 0)

		assert.Nil(t, err)
		assert.Equal(t, 0, len(notes))
	})))
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		userId := CreateUserInDB(t, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)
		result, err := db.Tx(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) (any, error) {
			tagId, err := CreateTagInDB(t, tx, ctx, TEST_TAG_NAME_1, TEST_TAG_STATE_1)
			return tagId, err
		})()
		tagId, ok := result.(int)
		assert.True(t, ok)
		var expectedNotes []entities.Note
		for i := 1; i <= 10; i++ {
			expectedNotes = append(expectedNotes, GenerateNote(i, userId, tagId))
		}
		CreateNotesInDB(t, 10, TEST_NOTE_TEXT_TEMPLATE, TEST_NOTE_TOPIC_TEMPLATE, tagId, userId, TEST_NOTE_STATE_1)
		actualNotes, err := queries.GetNotes(db.GetInstance().GetDB(), 50, 0)

		assert.Nil(t, err)
		AssertEqualNoteArrays(t, expectedNotes, actualNotes)
	})))
	t.Run("LimitParameterCase", RunWithRecreateDB((func(t *testing.T) {
		userId := CreateUserInDB(t, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)
		result, err := db.Tx(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) (any, error) {
			tagId, err := CreateTagInDB(t, tx, ctx, TEST_TAG_NAME_1, TEST_TAG_STATE_1)
			return tagId, err
		})()
		tagId, ok := result.(int)
		assert.True(t, ok)
		var expectedNotes []entities.Note
		for i := 1; i <= 5; i++ {
			expectedNotes = append(expectedNotes, GenerateNote(i, userId, tagId))
		}
		CreateNotesInDB(t, 10, TEST_NOTE_TEXT_TEMPLATE, TEST_NOTE_TOPIC_TEMPLATE, tagId, userId, TEST_NOTE_STATE_1)
		actualNotes, err := queries.GetNotes(db.GetInstance().GetDB(), 5, 0)

		assert.Nil(t, err)
		AssertEqualNoteArrays(t, expectedNotes, actualNotes)
	})))
	t.Run("OffsetParameterCase", RunWithRecreateDB((func(t *testing.T) {
		userId := CreateUserInDB(t, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)
		result, err := db.Tx(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) (any, error) {
			tagId, err := CreateTagInDB(t, tx, ctx, TEST_TAG_NAME_1, TEST_TAG_STATE_1)
			return tagId, err
		})()
		tagId, ok := result.(int)
		assert.True(t, ok)
		var expectedNotes []entities.Note
		for i := 6; i <= 10; i++ {
			expectedNotes = append(expectedNotes, GenerateNote(i, userId, tagId))
		}
		CreateNotesInDB(t, 10, TEST_NOTE_TEXT_TEMPLATE, TEST_NOTE_TOPIC_TEMPLATE, tagId, userId, TEST_NOTE_STATE_1)
		actualNotes, err := queries.GetNotes(db.GetInstance().GetDB(), 50, 5)

		assert.Nil(t, err)
		AssertEqualNoteArrays(t, expectedNotes, actualNotes)
	})))
}

func TestDBNoteUpdate(t *testing.T) {
	t.Run("NotFoundCase", RunWithRecreateDB((func(t *testing.T) {
		userId := CreateUserInDB(t, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)
		result, err := db.Tx(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) (any, error) {
			tagId, err := CreateTagInDB(t, tx, ctx, TEST_TAG_NAME_1, TEST_TAG_STATE_1)
			return tagId, err
		})()
		tagId, ok := result.(int)
		assert.True(t, ok)

		err = queries.UpdateNote(db.GetInstance().GetDB(), 1, TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, tagId, userId, TEST_NOTE_STATE_1)

		assert.Equal(t, sql.ErrNoRows, err)
	})))
	t.Run("DeletedCase", RunWithRecreateDB((func(t *testing.T) {
		userId := CreateUserInDB(t, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)
		result, err := db.Tx(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) (any, error) {
			tagId, err := CreateTagInDB(t, tx, ctx, TEST_TAG_NAME_1, TEST_TAG_STATE_1)
			return tagId, err
		})()
		tagId, ok := result.(int)
		assert.True(t, ok)

		noteId, err := queries.CreateNote(db.GetInstance().GetDB(), TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, tagId, userId, TEST_NOTE_STATE_1)

		assert.Nil(t, err)
		assert.NotEqual(t, noteId, -1)

		err = queries.DeleteNote(db.GetInstance().GetDB(), noteId)

		assert.Nil(t, err)

		err = queries.UpdateNote(db.GetInstance().GetDB(), noteId, TEST_NOTE_TEXT_2, TEST_NOTE_TOPIC_2, tagId, userId, TEST_NOTE_STATE_2)

		assert.Equal(t, sql.ErrNoRows, err)
	})))
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		userId := CreateUserInDB(t, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)
		result, err := db.Tx(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) (any, error) {
			tagId, err := CreateTagInDB(t, tx, ctx, TEST_TAG_NAME_1, TEST_TAG_STATE_1)
			return tagId, err
		})()
		tagId, ok := result.(int)
		assert.True(t, ok)

		expected := GenerateNote(1, userId, tagId)

		noteId, err := queries.CreateNote(db.GetInstance().GetDB(), TEST_NOTE_TEXT_2, TEST_NOTE_TOPIC_2, userId, tagId, TEST_NOTE_STATE_2)

		assert.Nil(t, err)
		assert.Equal(t, expected.Id, noteId)

		err = queries.UpdateNote(db.GetInstance().GetDB(), expected.Id, expected.Text, expected.Topic, expected.TagId, expected.UserId, expected.State)

		assert.Nil(t, err)

		actual, err := queries.GetNote(db.GetInstance().GetDB(), expected.Id)

		AssertEqualNotes(t, expected, actual)
	})))
}

func TestDBNoteDelete(t *testing.T) {
	t.Run("NotFoundCase", RunWithRecreateDB((func(t *testing.T) {
		err := queries.DeleteNote(db.GetInstance().GetDB(), 1)

		assert.Equal(t, sql.ErrNoRows, err)
	})))
	t.Run("AlreadyDeletedCase", RunWithRecreateDB((func(t *testing.T) {
		userId := CreateUserInDB(t, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)
		result, err := db.Tx(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) (any, error) {
			tagId, err := CreateTagInDB(t, tx, ctx, TEST_TAG_NAME_1, TEST_TAG_STATE_1)
			return tagId, err
		})()
		tagId, ok := result.(int)
		assert.True(t, ok)
		noteId := CreateNoteInDB(t, TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, tagId, userId, TEST_NOTE_STATE_1)

		err = queries.DeleteNote(db.GetInstance().GetDB(), noteId)

		assert.Nil(t, err)

		err = queries.DeleteNote(db.GetInstance().GetDB(), noteId)

		assert.Equal(t, sql.ErrNoRows, err)
	})))
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		userId := CreateUserInDB(t, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)
		result, err := db.Tx(func(tx *sql.Tx, ctx context.Context, cancel context.CancelFunc) (any, error) {
			tagId, err := CreateTagInDB(t, tx, ctx, TEST_TAG_NAME_1, TEST_TAG_STATE_1)
			return tagId, err
		})()
		tagId, ok := result.(int)
		assert.True(t, ok)

		var expectedNotes []entities.Note
		expectedNotes = append(expectedNotes, GenerateNote(1, userId, tagId))
		expectedNotes = append(expectedNotes, GenerateNote(3, userId, tagId))

		noteIdToDelete := 2

		CreateNotesInDB(t, 3, TEST_NOTE_TEXT_TEMPLATE, TEST_NOTE_TOPIC_TEMPLATE, tagId, userId, TEST_NOTE_STATE_1)

		err = queries.DeleteNote(db.GetInstance().GetDB(), noteIdToDelete)

		assert.Nil(t, err)

		notes, err := queries.GetNotes(db.GetInstance().GetDB(), 50, 0)

		assert.Nil(t, err)
		AssertEqualNoteArrays(t, expectedNotes, notes)

		_, err = queries.GetNote(db.GetInstance().GetDB(), noteIdToDelete)

		assert.Equal(t, sql.ErrNoRows, err)
	})))
}
