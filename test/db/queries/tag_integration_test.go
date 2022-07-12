//go:build integration
// +build integration

package queries_test

import (
	"database/sql"
	"fmt"
	"strconv"
	"testing"

	integrationTesting "github.com/ArtemVoronov/indefinite-studies-api/internal/app/testing"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/db"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/db/entities"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/db/queries"
	"github.com/stretchr/testify/assert"
)

const (
	TEST_TAG_NAME_1        string = "Test tag 1"
	TEST_TAG_STATE_1       string = entities.TAG_STATE_NEW
	TEST_TAG_NAME_2        string = "Test tag 2"
	TEST_TAG_STATE_2       string = entities.TAG_STATE_BLOCKED
	TEST_TAG_NAME_TEMPLATE string = "Test tag "
)

var TagDuplicateKeyConstraintViolationError = fmt.Errorf(DuplicateKeyConstraintViolationError, "tags_name_state_unique")

func TestGetTag(t *testing.T) {
	t.Run("ExpectedNotFoundError", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedError := sql.ErrNoRows

		_, actualError := queries.GetTag(db.DB, 1)

		assert.Equal(t, expectedError, actualError)
	})))
	t.Run("ExpectedResult", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedName := TEST_TAG_NAME_1
		expectedState := TEST_TAG_STATE_1
		expectedId, err := queries.CreateTag(db.DB, expectedName, expectedState)
		if err != nil || expectedId == -1 {
			t.Errorf("Unable to create tag: %s", err)
		}

		actual, err := queries.GetTag(db.DB, expectedId)

		assert.Equal(t, expectedId, actual.Id)
		assert.Equal(t, expectedName, actual.Name)
		assert.Equal(t, expectedState, actual.State)
	})))
}

func TestCreateTag(t *testing.T) {
	t.Run("BasicCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedTagId := 1

		actualTagId, err := queries.CreateTag(db.DB, TEST_TAG_NAME_1, TEST_TAG_STATE_1)
		if err != nil || actualTagId == -1 {
			t.Errorf("Unable to create tag: %s", err)
		}

		assert.Equal(t, expectedTagId, actualTagId)
	})))
	t.Run("DuplicateCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedError := fmt.Errorf("error at inserting tag (Name: '%s', State: '%s') into db, case after db.QueryRow.Scan: %s", TEST_TAG_NAME_1, TEST_TAG_STATE_1, TagDuplicateKeyConstraintViolationError)

		tagId, err := queries.CreateTag(db.DB, TEST_TAG_NAME_1, TEST_TAG_STATE_1)
		if err != nil || tagId == -1 {
			t.Errorf("Unable to create tag: %s", err)
		}
		_, actualError := queries.CreateTag(db.DB, TEST_TAG_NAME_1, TEST_TAG_STATE_1)

		assert.Equal(t, expectedError, actualError)
	})))
}

func TestGetTags(t *testing.T) {
	t.Run("ExpectedEmpty", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedArrayLength := 0

		tags, err := queries.GetTags(db.DB, "50", "0")
		if err != nil {
			t.Errorf("Unable to get to tags : %s", err)
		}
		actualArrayLength := len(tags)

		assert.Equal(t, expectedArrayLength, actualArrayLength)
	})))
	t.Run("ExpectedResult", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedArrayLength := 3

		for i := 0; i < 3; i++ {
			tagId, err := queries.CreateTag(db.DB, TEST_TAG_NAME_TEMPLATE+strconv.Itoa(i), entities.TAG_STATE_NEW)
			if err != nil || tagId == -1 {
				t.Errorf("Unable to create tag: %s", err)
			}
		}
		tags, err := queries.GetTags(db.DB, "50", "0")
		if err != nil {
			t.Errorf("Unable to get to tags : %s", err)
		}
		actualArrayLength := len(tags)

		assert.Equal(t, expectedArrayLength, actualArrayLength)
		for i, tag := range tags {
			assert.Equal(t, i+1, tag.Id)
			assert.Equal(t, TEST_TAG_NAME_TEMPLATE+strconv.Itoa(i), tag.Name)
			assert.Equal(t, entities.TAG_STATE_NEW, tag.State)
		}
	})))
	t.Run("LimitParameterCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedArrayLength := 5
		for i := 0; i < 10; i++ {
			tagId, err := queries.CreateTag(db.DB, TEST_TAG_NAME_TEMPLATE+strconv.Itoa(i), entities.TAG_STATE_NEW)
			if err != nil || tagId == -1 {
				t.Errorf("Unable to create tag: %s", err)
			}
		}

		tags, err := queries.GetTags(db.DB, "5", "0")
		if err != nil {
			t.Errorf("Unable to get to tags : %s", err)
		}
		actualArrayLength := len(tags)

		assert.Equal(t, expectedArrayLength, actualArrayLength)
		for i, tag := range tags {
			assert.Equal(t, i+1, tag.Id)
			assert.Equal(t, TEST_TAG_NAME_TEMPLATE+strconv.Itoa(i), tag.Name)
			assert.Equal(t, entities.TAG_STATE_NEW, tag.State)
		}
	})))
	t.Run("OffsetParameterCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedArrayLength := 5
		for i := 0; i < 10; i++ {
			tagId, err := queries.CreateTag(db.DB, TEST_TAG_NAME_TEMPLATE+strconv.Itoa(i), entities.TAG_STATE_NEW)
			if err != nil || tagId == -1 {
				t.Errorf("Unable to create tag: %s", err)
			}
		}

		tags, err := queries.GetTags(db.DB, "50", "5")
		if err != nil {
			t.Errorf("Unable to get to tags : %s", err)
		}
		actualArrayLength := len(tags)

		assert.Equal(t, expectedArrayLength, actualArrayLength)
		for i, tag := range tags {
			assert.Equal(t, i+6, tag.Id)
			assert.Equal(t, TEST_TAG_NAME_TEMPLATE+strconv.Itoa(i+5), tag.Name)
			assert.Equal(t, entities.TAG_STATE_NEW, tag.State)
		}
	})))
}

func TestUpdateTag(t *testing.T) {
	t.Run("ExpectedNotFoundError", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedError := sql.ErrNoRows

		actualError := queries.UpdateTag(db.DB, 1, TEST_TAG_NAME_1, TEST_TAG_STATE_1)

		assert.Equal(t, expectedError, actualError)
	})))
	t.Run("DeletedCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedError := sql.ErrNoRows

		tagId, err := queries.CreateTag(db.DB, TEST_TAG_NAME_1, TEST_TAG_STATE_1)
		if err != nil || tagId == -1 {
			t.Errorf("Unable to create tag: %s", err)
		}

		err = queries.DeleteTag(db.DB, tagId)
		if err != nil {
			t.Errorf("Unable to delete tag: %s", err)
		}

		actualError := queries.UpdateTag(db.DB, tagId, TEST_TAG_NAME_2, TEST_TAG_STATE_2)

		assert.Equal(t, expectedError, actualError)
	})))
	t.Run("BasicCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedName := TEST_TAG_NAME_2
		expectedState := TEST_TAG_STATE_2
		expectedId, err := queries.CreateTag(db.DB, TEST_TAG_NAME_1, TEST_TAG_STATE_1)
		if err != nil || expectedId == -1 {
			t.Errorf("Unable to create tag: %s", err)
		}

		err = queries.UpdateTag(db.DB, expectedId, TEST_TAG_NAME_2, TEST_TAG_STATE_2)
		if err != nil {
			t.Errorf("Unable to update tag: %s", err)
		}

		actual, err := queries.GetTag(db.DB, expectedId)

		assert.Equal(t, expectedId, actual.Id)
		assert.Equal(t, expectedName, actual.Name)
		assert.Equal(t, expectedState, actual.State)
	})))
	t.Run("DuplicateCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		tagId, err := queries.CreateTag(db.DB, TEST_TAG_NAME_1, TEST_TAG_STATE_1)
		if err != nil || tagId == -1 {
			t.Errorf("Unable to create tag: %s", err)
		}
		tagId, err = queries.CreateTag(db.DB, TEST_TAG_NAME_2, TEST_TAG_STATE_2)
		if err != nil || tagId == -1 {
			t.Errorf("Unable to create tag: %s", err)
		}

		expectedError := fmt.Errorf("error at updating tag (Id: %d, Name: '%s', State: '%s'), case after executing statement: %s", tagId, TEST_TAG_NAME_1, TEST_TAG_STATE_1, TagDuplicateKeyConstraintViolationError)

		actualError := queries.UpdateTag(db.DB, tagId, TEST_TAG_NAME_1, TEST_TAG_STATE_1)

		assert.Equal(t, expectedError, actualError)
	})))
}

func TestDeleteTag(t *testing.T) {
	t.Run("NotFoundCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		notExistentTagId := 1

		actualError := queries.DeleteTag(db.DB, notExistentTagId)
		if actualError != nil {
			t.Errorf("Unable to delete tag: %s", actualError)
		}
	})))
	t.Run("AlreadyDeletedCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		tagId, err := queries.CreateTag(db.DB, TEST_TAG_NAME_1, TEST_TAG_STATE_1)
		if err != nil || tagId == -1 {
			t.Errorf("Unable to create tag: %s", err)
		}

		err = queries.DeleteTag(db.DB, tagId)
		if err != nil {
			t.Errorf("Unable to delete tag: %s", err)
		}

		actualError := queries.DeleteTag(db.DB, tagId)
		if actualError != nil {
			t.Errorf("Unable to delete tag: %s", actualError)
		}
	})))
	t.Run("BasicCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedFirstTagId := 1
		expectedSecondTagId := 3
		expectedState := entities.TAG_STATE_NEW
		expectedError := sql.ErrNoRows
		expectedArrayLength := 2
		tagIdToDelete := 2
		for i := 0; i < 3; i++ {
			tagId, err := queries.CreateTag(db.DB, TEST_TAG_NAME_TEMPLATE+strconv.Itoa(i), entities.TAG_STATE_NEW)
			if err != nil || tagId == -1 {
				t.Errorf("Unable to create tag: %s", err)
			}
		}

		err := queries.DeleteTag(db.DB, tagIdToDelete)
		if err != nil {
			t.Errorf("Unable to delete tag: %s", err)
		}

		tags, err := queries.GetTags(db.DB, "50", "0")
		if err != nil {
			t.Errorf("Unable to get to tags : %s", err)
		}
		actualArrayLength := len(tags)

		assert.Equal(t, expectedArrayLength, actualArrayLength)

		assert.Equal(t, expectedFirstTagId, tags[0].Id)
		assert.Equal(t, TEST_TAG_NAME_TEMPLATE+"0", tags[0].Name)
		assert.Equal(t, expectedState, tags[0].State)
		assert.Equal(t, expectedSecondTagId, tags[1].Id)
		assert.Equal(t, TEST_TAG_NAME_TEMPLATE+"2", tags[1].Name)
		assert.Equal(t, expectedState, tags[1].State)

		_, actualError := queries.GetTag(db.DB, tagIdToDelete)

		assert.Equal(t, expectedError, actualError)
	})))
}
