//go:build integration
// +build integration

package queries_test

import (
	"database/sql"
	integrationTesting "github.com/ArtemVoronov/indefinite-studies-api/internal/app/testing"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/db"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/db/entities"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/db/queries"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

const (
	TEST_NOTE_TEXT_1  string = "Test text 1"
	TEST_NOTE_TOPIC_1 string = "Test topic 1"
	TEST_NOTE_STATE_1 string = entities.NOTE_STATE_NEW
	TEST_NOTE_TEXT_2  string = "Test text 2"
	TEST_NOTE_TOPIC_2 string = "Test topic 2"
	TEST_NOTE_STATE_2 string = entities.NOTE_STATE_NEW

	TEST_NOTE_TEXT_TEMPLATE  string = "Test text "
	TEST_NOTE_TOPIC_TEMPLATE string = "Test topic"
)

func TestGetNote(t *testing.T) {
	t.Run("ExpectedNotFoundError", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedError := sql.ErrNoRows

		_, actualError := queries.GetNote(db.DB, 1)

		assert.Equal(t, expectedError, actualError)
	})))
	t.Run("ExpectedResult", integrationTesting.RunWithRecreateDB((func(t *testing.T) {

		userId, err := queries.CreateUser(db.DB, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)
		if err != nil || userId == -1 {
			t.Errorf("Unable to create user: %s", err)
		}

		tagId, err := queries.CreateTag(db.DB, TEST_TAG_NAME_1, TEST_TAG_STATE_1)
		if err != nil || tagId == -1 {
			t.Errorf("Unable to create tag: %s", err)
		}

		expectedNoteId, err := queries.CreateNote(db.DB, TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, tagId, userId, TEST_NOTE_STATE_1)
		if err != nil || expectedNoteId == -1 {
			t.Errorf("Unable to create note: %s", err)
		}

		actual, err := queries.GetNote(db.DB, expectedNoteId)

		assert.Equal(t, expectedNoteId, actual.Id)
		assert.Equal(t, TEST_NOTE_TEXT_1, actual.Text)
		assert.Equal(t, TEST_NOTE_TOPIC_1, actual.Topic)
		assert.Equal(t, tagId, actual.TagId)
		assert.Equal(t, userId, actual.UserId)
		assert.Equal(t, TEST_NOTE_STATE_1, actual.State)
	})))
}

func TestCreateNote(t *testing.T) {
	t.Run("BasicCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedNoteId := 1

		userId, err := queries.CreateUser(db.DB, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)
		if err != nil || userId == -1 {
			t.Errorf("Unable to create user: %s", err)
		}

		tagId, err := queries.CreateTag(db.DB, TEST_TAG_NAME_1, TEST_TAG_STATE_1)
		if err != nil || tagId == -1 {
			t.Errorf("Unable to create tag: %s", err)
		}

		actualNoteId, err := queries.CreateNote(db.DB, TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, tagId, userId, TEST_NOTE_STATE_1)
		if err != nil || actualNoteId == -1 {
			t.Errorf("Unable to create note: %s", err)
		}

		assert.Equal(t, expectedNoteId, actualNoteId)
	})))
}

func TestGetNotes(t *testing.T) {
	t.Run("ExpectedEmpty", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedArrayLength := 0

		notes, err := queries.GetNotes(db.DB, "50", "0")
		if err != nil {
			t.Errorf("Unable to get to notes : %s", err)
		}
		actualArrayLength := len(notes)

		assert.Equal(t, expectedArrayLength, actualArrayLength)
	})))
	t.Run("ExpectedResult", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedArrayLength := 3
		userId, err := queries.CreateUser(db.DB, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)
		if err != nil || userId == -1 {
			t.Errorf("Unable to create user: %s", err)
		}

		tagId, err := queries.CreateTag(db.DB, TEST_TAG_NAME_1, TEST_TAG_STATE_1)
		if err != nil || tagId == -1 {
			t.Errorf("Unable to create tag: %s", err)
		}

		for i := 0; i < 3; i++ {
			noteId, err := queries.CreateNote(db.DB, TEST_NOTE_TEXT_TEMPLATE+strconv.Itoa(i), TEST_NOTE_TOPIC_TEMPLATE+strconv.Itoa(i), tagId, userId, TEST_NOTE_STATE_1)
			if err != nil || noteId == -1 {
				t.Errorf("Unable to create note: %s", err)
			}
		}
		notes, err := queries.GetNotes(db.DB, "50", "0")
		if err != nil {
			t.Errorf("Unable to get to notes : %s", err)
		}
		actualArrayLength := len(notes)

		assert.Equal(t, expectedArrayLength, actualArrayLength)
		for i, note := range notes {
			assert.Equal(t, i+1, note.Id)
			assert.Equal(t, TEST_NOTE_TEXT_TEMPLATE+strconv.Itoa(i), note.Text)
			assert.Equal(t, TEST_NOTE_TOPIC_TEMPLATE+strconv.Itoa(i), note.Topic)
			assert.Equal(t, tagId, note.TagId)
			assert.Equal(t, userId, note.UserId)
			assert.Equal(t, TEST_NOTE_STATE_1, note.State)
		}
	})))
	t.Run("LimitParameterCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedArrayLength := 5
		userId, err := queries.CreateUser(db.DB, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)
		if err != nil || userId == -1 {
			t.Errorf("Unable to create user: %s", err)
		}

		tagId, err := queries.CreateTag(db.DB, TEST_TAG_NAME_1, TEST_TAG_STATE_1)
		if err != nil || tagId == -1 {
			t.Errorf("Unable to create tag: %s", err)
		}

		for i := 0; i < 10; i++ {
			noteId, err := queries.CreateNote(db.DB, TEST_NOTE_TEXT_TEMPLATE+strconv.Itoa(i), TEST_NOTE_TOPIC_TEMPLATE+strconv.Itoa(i), tagId, userId, TEST_NOTE_STATE_1)
			if err != nil || noteId == -1 {
				t.Errorf("Unable to create note: %s", err)
			}
		}

		notes, err := queries.GetNotes(db.DB, "5", "0")
		if err != nil {
			t.Errorf("Unable to get to notes : %s", err)
		}
		actualArrayLength := len(notes)

		assert.Equal(t, expectedArrayLength, actualArrayLength)
		for i, note := range notes {
			assert.Equal(t, i+1, note.Id)
			assert.Equal(t, TEST_NOTE_TEXT_TEMPLATE+strconv.Itoa(i), note.Text)
			assert.Equal(t, TEST_NOTE_TOPIC_TEMPLATE+strconv.Itoa(i), note.Topic)
			assert.Equal(t, tagId, note.TagId)
			assert.Equal(t, userId, note.UserId)
			assert.Equal(t, TEST_NOTE_STATE_1, note.State)
		}
	})))
	t.Run("OffsetParameterCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedArrayLength := 5
		userId, err := queries.CreateUser(db.DB, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)
		if err != nil || userId == -1 {
			t.Errorf("Unable to create user: %s", err)
		}

		tagId, err := queries.CreateTag(db.DB, TEST_TAG_NAME_1, TEST_TAG_STATE_1)
		if err != nil || tagId == -1 {
			t.Errorf("Unable to create tag: %s", err)
		}
		for i := 0; i < 10; i++ {
			noteId, err := queries.CreateNote(db.DB, TEST_NOTE_TEXT_TEMPLATE+strconv.Itoa(i), TEST_NOTE_TOPIC_TEMPLATE+strconv.Itoa(i), tagId, userId, TEST_NOTE_STATE_1)
			if err != nil || noteId == -1 {
				t.Errorf("Unable to create note: %s", err)
			}
		}

		notes, err := queries.GetNotes(db.DB, "50", "5")
		if err != nil {
			t.Errorf("Unable to get to notes : %s", err)
		}
		actualArrayLength := len(notes)

		assert.Equal(t, expectedArrayLength, actualArrayLength)
		for i, note := range notes {
			assert.Equal(t, i+6, note.Id)
			assert.Equal(t, TEST_NOTE_TEXT_TEMPLATE+strconv.Itoa(i+5), note.Text)
			assert.Equal(t, TEST_NOTE_TOPIC_TEMPLATE+strconv.Itoa(i+5), note.Topic)
			assert.Equal(t, tagId, note.TagId)
			assert.Equal(t, userId, note.UserId)
			assert.Equal(t, TEST_NOTE_STATE_1, note.State)
		}
	})))
}

func TestUpdateNote(t *testing.T) {
	t.Run("ExpectedNotFoundError", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedError := sql.ErrNoRows

		userId, err := queries.CreateUser(db.DB, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)
		if err != nil || userId == -1 {
			t.Errorf("Unable to create user: %s", err)
		}

		tagId, err := queries.CreateTag(db.DB, TEST_TAG_NAME_1, TEST_TAG_STATE_1)
		if err != nil || tagId == -1 {
			t.Errorf("Unable to create tag: %s", err)
		}

		actualError := queries.UpdateNote(db.DB, 1, TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, tagId, userId, TEST_NOTE_STATE_1)

		assert.Equal(t, expectedError, actualError)
	})))
	t.Run("DeletedCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedError := sql.ErrNoRows

		userId, err := queries.CreateUser(db.DB, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)
		if err != nil || userId == -1 {
			t.Errorf("Unable to create user: %s", err)
		}

		tagId, err := queries.CreateTag(db.DB, TEST_TAG_NAME_1, TEST_TAG_STATE_1)
		if err != nil || tagId == -1 {
			t.Errorf("Unable to create tag: %s", err)
		}

		noteId, err := queries.CreateNote(db.DB, TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, tagId, userId, TEST_NOTE_STATE_1)
		if err != nil || noteId == -1 {
			t.Errorf("Unable to create note: %s", err)
		}

		err = queries.DeleteNote(db.DB, noteId)
		if err != nil {
			t.Errorf("Unable to delete note: %s", err)
		}

		actualError := queries.UpdateNote(db.DB, 1, TEST_NOTE_TEXT_2, TEST_NOTE_TOPIC_2, tagId, userId, TEST_NOTE_STATE_2)

		assert.Equal(t, expectedError, actualError)
	})))
	t.Run("BasicCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedText := TEST_NOTE_TEXT_2
		expectedTopic := TEST_NOTE_TOPIC_2
		expectedState := TEST_NOTE_STATE_2
		expectedUserId, err := queries.CreateUser(db.DB, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)
		if err != nil || expectedUserId == -1 {
			t.Errorf("Unable to create user: %s", err)
		}

		expectedTagId, err := queries.CreateTag(db.DB, TEST_TAG_NAME_1, TEST_TAG_STATE_1)
		if err != nil || expectedTagId == -1 {
			t.Errorf("Unable to create tag: %s", err)
		}

		expectedId, err := queries.CreateNote(db.DB, TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, expectedUserId, expectedUserId, TEST_NOTE_STATE_1)
		if err != nil || expectedId == -1 {
			t.Errorf("Unable to create note: %s", err)
		}

		err = queries.UpdateNote(db.DB, expectedId, TEST_NOTE_TEXT_2, TEST_NOTE_TOPIC_2, expectedTagId, expectedUserId, TEST_NOTE_STATE_2)
		if err != nil {
			t.Errorf("Unable to update user: %s", err)
		}

		actual, err := queries.GetNote(db.DB, expectedId)

		assert.Equal(t, expectedId, actual.Id)
		assert.Equal(t, expectedText, actual.Text)
		assert.Equal(t, expectedTopic, actual.Topic)
		assert.Equal(t, expectedTagId, actual.TagId)
		assert.Equal(t, expectedUserId, actual.UserId)
		assert.Equal(t, expectedState, actual.State)
	})))
}

func TestDeleteNote(t *testing.T) {
	t.Run("NotFoundCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		notExistentNoteId := 1

		actualError := queries.DeleteNote(db.DB, notExistentNoteId)
		if actualError != nil {
			t.Errorf("Unable to delete note: %s", actualError)
		}
	})))
	t.Run("AlreadyDeletedCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		userId, err := queries.CreateUser(db.DB, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)
		if err != nil || userId == -1 {
			t.Errorf("Unable to create user: %s", err)
		}

		tagId, err := queries.CreateTag(db.DB, TEST_TAG_NAME_1, TEST_TAG_STATE_1)
		if err != nil || tagId == -1 {
			t.Errorf("Unable to create tag: %s", err)
		}

		noteId, err := queries.CreateNote(db.DB, TEST_NOTE_TEXT_1, TEST_NOTE_TOPIC_1, tagId, userId, TEST_NOTE_STATE_1)
		if err != nil || noteId == -1 {
			t.Errorf("Unable to create note: %s", err)
		}

		err = queries.DeleteNote(db.DB, noteId)
		if err != nil {
			t.Errorf("Unable to delete note: %s", err)
		}

		actualError := queries.DeleteNote(db.DB, noteId)
		if actualError != nil {
			t.Errorf("Unable to delete note: %s", actualError)
		}
	})))
	t.Run("BasicCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedFirstNoteId := 1
		expectedSecondNoteId := 3
		expectedState := TEST_USER_STATE_1
		expectedError := sql.ErrNoRows
		expectedArrayLength := 2
		noteIdToDelete := 2
		userId, err := queries.CreateUser(db.DB, TEST_USER_LOGIN_1, TEST_USER_EMAIL_1, TEST_USER_PASSWORD_1, TEST_USER_ROLE_1, TEST_USER_STATE_1)
		if err != nil || userId == -1 {
			t.Errorf("Unable to create user: %s", err)
		}

		tagId, err := queries.CreateTag(db.DB, TEST_TAG_NAME_1, TEST_TAG_STATE_1)
		if err != nil || tagId == -1 {
			t.Errorf("Unable to create tag: %s", err)
		}

		for i := 0; i < 3; i++ {
			noteId, err := queries.CreateNote(db.DB, TEST_NOTE_TEXT_TEMPLATE+strconv.Itoa(i), TEST_NOTE_TOPIC_TEMPLATE+strconv.Itoa(i), tagId, userId, TEST_NOTE_STATE_1)
			if err != nil || noteId == -1 {
				t.Errorf("Unable to create note: %s", err)
			}
		}

		err = queries.DeleteNote(db.DB, noteIdToDelete)
		if err != nil {
			t.Errorf("Unable to delete note: %s", err)
		}

		notes, err := queries.GetNotes(db.DB, "50", "0")
		if err != nil {
			t.Errorf("Unable to get to notes : %s", err)
		}
		actualArrayLength := len(notes)

		assert.Equal(t, expectedArrayLength, actualArrayLength)

		assert.Equal(t, expectedFirstNoteId, notes[0].Id)
		assert.Equal(t, TEST_NOTE_TEXT_TEMPLATE+"0", notes[0].Text)
		assert.Equal(t, TEST_NOTE_TOPIC_TEMPLATE+"0", notes[0].Topic)
		assert.Equal(t, tagId, notes[0].TagId)
		assert.Equal(t, userId, notes[0].UserId)
		assert.Equal(t, expectedState, notes[0].State)

		assert.Equal(t, expectedSecondNoteId, notes[1].Id)
		assert.Equal(t, TEST_NOTE_TEXT_TEMPLATE+"2", notes[1].Text)
		assert.Equal(t, TEST_NOTE_TOPIC_TEMPLATE+"2", notes[1].Topic)
		assert.Equal(t, tagId, notes[1].TagId)
		assert.Equal(t, userId, notes[1].UserId)
		assert.Equal(t, expectedState, notes[1].State)

		_, actualError := queries.GetNote(db.DB, noteIdToDelete)

		assert.Equal(t, expectedError, actualError)
	})))
}
