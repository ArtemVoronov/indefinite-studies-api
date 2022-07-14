//go:build integration
// +build integration

package queries_test

import (
	"database/sql"
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

func GenerateTag(id int) entities.Tag {
	return entities.Tag{
		Id:    id,
		Name:  GenerateTagName(TEST_TAG_NAME_TEMPLATE, id),
		State: TEST_TAG_STATE_1,
	}
}

func GenerateTagName(template string, id int) string {
	return template + strconv.Itoa(id)
}

func AssertEqualTags(t *testing.T, expected entities.Tag, actual entities.Tag) {
	assert.Equal(t, expected.Id, actual.Id)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.State, actual.State)
}

func AssertEqualTagArrays(t *testing.T, expected []entities.Tag, actual []entities.Tag) {
	assert.Equal(t, len(expected), len(actual))

	length := len(expected)
	for i := 0; i < length; i++ {
		AssertEqualTags(t, expected[i], actual[i])
	}
}

func CreateTagInDB(t *testing.T, name string, state string) int {
	tagId, err := queries.CreateTag(db.GetInstance().GetDB(), name, state)
	assert.Nil(t, err)
	assert.NotEqual(t, tagId, -1)
	return tagId
}

func CreateTagsInDB(t *testing.T, count int, nameTemplate string, state string) {
	for i := 1; i <= count; i++ {
		CreateTagInDB(t, nameTemplate+strconv.Itoa(i), state)
	}
}

func TestDBTagGet(t *testing.T) {
	t.Run("NotFoundCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		_, err := queries.GetTag(db.GetInstance().GetDB(), 1)

		assert.Equal(t, sql.ErrNoRows, err)
	})))
	t.Run("BasicCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expected := GenerateTag(1)

		tagId, err := queries.CreateTag(db.GetInstance().GetDB(), expected.Name, expected.State)

		assert.Nil(t, err)
		assert.Equal(t, tagId, expected.Id)

		actual, err := queries.GetTag(db.GetInstance().GetDB(), tagId)

		AssertEqualTags(t, expected, actual)
	})))
}

func TestDBTagCreate(t *testing.T) {
	t.Run("BasicCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		tagId, err := queries.CreateTag(db.GetInstance().GetDB(), TEST_TAG_NAME_1, TEST_TAG_STATE_1)

		assert.Nil(t, err)
		assert.Equal(t, tagId, 1)
	})))
	t.Run("DuplicateCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		tagId, err := queries.CreateTag(db.GetInstance().GetDB(), TEST_TAG_NAME_1, TEST_TAG_STATE_1)

		assert.Nil(t, err)
		assert.NotEqual(t, tagId, -1)

		_, err = queries.CreateTag(db.GetInstance().GetDB(), TEST_TAG_NAME_1, TEST_TAG_STATE_1)

		assert.Equal(t, db.ErrorTagDuplicateKey, err)
	})))
}

func TestDBTagGetAll(t *testing.T) {
	t.Run("ExpectedEmpty", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		tags, err := queries.GetTags(db.GetInstance().GetDB(), 50, 0)

		assert.Nil(t, err)
		assert.Equal(t, 0, len(tags))
	})))
	t.Run("ExpectedResult", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		var expectedTags []entities.Tag
		for i := 1; i <= 10; i++ {
			expectedTags = append(expectedTags, GenerateTag(i))
		}
		CreateTagsInDB(t, 10, TEST_TAG_NAME_TEMPLATE, entities.TAG_STATE_NEW)

		actualTags, err := queries.GetTags(db.GetInstance().GetDB(), 50, 0)

		assert.Nil(t, err)
		AssertEqualTagArrays(t, expectedTags, actualTags)
	})))
	t.Run("LimitParameterCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		var expectedTags []entities.Tag
		for i := 1; i <= 5; i++ {
			expectedTags = append(expectedTags, GenerateTag(i))
		}

		CreateTagsInDB(t, 10, TEST_TAG_NAME_TEMPLATE, entities.TAG_STATE_NEW)

		actualTags, err := queries.GetTags(db.GetInstance().GetDB(), 5, 0)

		assert.Nil(t, err)
		AssertEqualTagArrays(t, expectedTags, actualTags)
	})))
	t.Run("OffsetParameterCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		var expectedTags []entities.Tag
		for i := 6; i <= 10; i++ {
			expectedTags = append(expectedTags, GenerateTag(i))
		}

		CreateTagsInDB(t, 10, TEST_TAG_NAME_TEMPLATE, entities.TAG_STATE_NEW)

		actualTags, err := queries.GetTags(db.GetInstance().GetDB(), 50, 5)

		assert.Nil(t, err)
		AssertEqualTagArrays(t, expectedTags, actualTags)
	})))
}

func TestDBTagUpdate(t *testing.T) {
	t.Run("NotFoundCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		err := queries.UpdateTag(db.GetInstance().GetDB(), 1, TEST_TAG_NAME_1, TEST_TAG_STATE_1)

		assert.Equal(t, sql.ErrNoRows, err)
	})))
	t.Run("DeletedCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		tagId, err := queries.CreateTag(db.GetInstance().GetDB(), TEST_TAG_NAME_1, TEST_TAG_STATE_1)

		assert.Nil(t, err)
		assert.NotEqual(t, tagId, -1)

		err = queries.DeleteTag(db.GetInstance().GetDB(), tagId)

		assert.Nil(t, err)

		err = queries.UpdateTag(db.GetInstance().GetDB(), tagId, TEST_TAG_NAME_2, TEST_TAG_STATE_2)

		assert.Equal(t, sql.ErrNoRows, err)
	})))
	t.Run("BasicCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expected := GenerateTag(1)

		tagId, err := queries.CreateTag(db.GetInstance().GetDB(), TEST_TAG_NAME_2, TEST_TAG_STATE_2)

		assert.Nil(t, err)
		assert.Equal(t, expected.Id, tagId)

		err = queries.UpdateTag(db.GetInstance().GetDB(), expected.Id, expected.Name, expected.State)

		assert.Nil(t, err)

		actual, err := queries.GetTag(db.GetInstance().GetDB(), expected.Id)

		AssertEqualTags(t, expected, actual)
	})))
	t.Run("DuplicateCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		tagId, err := queries.CreateTag(db.GetInstance().GetDB(), TEST_TAG_NAME_1, TEST_TAG_STATE_1)

		assert.Nil(t, err)
		assert.NotEqual(t, tagId, -1)

		tagId, err = queries.CreateTag(db.GetInstance().GetDB(), TEST_TAG_NAME_2, TEST_TAG_STATE_2)

		assert.Nil(t, err)
		assert.NotEqual(t, tagId, -1)

		actualError := queries.UpdateTag(db.GetInstance().GetDB(), tagId, TEST_TAG_NAME_1, TEST_TAG_STATE_1)

		assert.Equal(t, db.ErrorTagDuplicateKey, actualError)
	})))
}

func TestDBTagDelete(t *testing.T) {
	t.Run("NotFoundCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		err := queries.DeleteTag(db.GetInstance().GetDB(), 1)

		assert.Equal(t, sql.ErrNoRows, err)
	})))
	t.Run("AlreadyDeletedCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		tagId, err := queries.CreateTag(db.GetInstance().GetDB(), TEST_TAG_NAME_1, TEST_TAG_STATE_1)

		assert.Nil(t, err)
		assert.NotEqual(t, tagId, -1)

		err = queries.DeleteTag(db.GetInstance().GetDB(), tagId)

		assert.Nil(t, err)

		err = queries.DeleteTag(db.GetInstance().GetDB(), tagId)

		assert.Equal(t, sql.ErrNoRows, err)
	})))
	t.Run("BasicCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		var expectedTags []entities.Tag
		expectedTags = append(expectedTags, GenerateTag(1))
		expectedTags = append(expectedTags, GenerateTag(3))

		tagIdToDelete := 2

		CreateTagsInDB(t, 3, TEST_TAG_NAME_TEMPLATE, entities.TAG_STATE_NEW)

		err := queries.DeleteTag(db.GetInstance().GetDB(), tagIdToDelete)

		assert.Nil(t, err)

		tags, err := queries.GetTags(db.GetInstance().GetDB(), 50, 0)

		assert.Nil(t, err)
		AssertEqualTagArrays(t, expectedTags, tags)

		_, err = queries.GetTag(db.GetInstance().GetDB(), tagIdToDelete)

		assert.Equal(t, sql.ErrNoRows, err)
	})))
}
