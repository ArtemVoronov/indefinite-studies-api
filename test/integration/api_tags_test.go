//go:build integration
// +build integration

package integration

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/ArtemVoronov/indefinite-studies-api/internal/api"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/db/entities"
	"github.com/stretchr/testify/assert"
)

var (
	ERROR_TAG_NAME_IS_REQUIRED string = "{\"errors\":[" +
		"{\"Field\":\"Name\",\"Msg\":\"This field is required\"}" +
		"]}"
	ERROR_TAG_STATE_IS_REQUIRED string = "{\"errors\":[" +
		"{\"Field\":\"State\",\"Msg\":\"This field is required\"}" +
		"]}"
	ERROR_TAG_NAME_AND_STATE_IS_REQUIRED string = "{\"errors\":[" +
		"{\"Field\":\"Name\",\"Msg\":\"This field is required\"}," +
		"{\"Field\":\"State\",\"Msg\":\"This field is required\"}" +
		"]}"
	ERROR_TAG_CREATE_STATE_WRONG_VALUE string = fmt.Sprintf("Unable to create tag. Wrong 'State' value. Possible values: %v", entities.GetPossibleTagStates())
	ERROR_TAG_UPDATE_STATE_WRONG_VALUE string = fmt.Sprintf("Unable to update tag. Wrong 'State' value. Possible values: %v", entities.GetPossibleTagStates())
)

func CreateTag(name any, state any) (int, string, error) {
	nameField := ""
	stateField := ""

	switch nameType := name.(type) {
	case int:
		nameField = "\"Name\": " + strconv.Itoa(name.(int))
	case string:
		nameField = "\"Name\": \"" + name.(string) + "\""
	case nil:
		nameField = ""
	default:
		return -1, "", fmt.Errorf("unkown type for 'name': %v", nameType)
	}

	switch stateType := state.(type) {
	case int:
		stateField = "\"State\": " + strconv.Itoa(state.(int))
	case string:
		stateField = "\"State\": \"" + state.(string) + "\""
	case nil:
		stateField = ""
	default:
		return -1, "", fmt.Errorf("unkown type for 'state': %v", stateType)
	}

	tagCreateDTO := "{"
	if nameField != "" && stateField != "" {
		tagCreateDTO += nameField + ", " + stateField
	} else if nameField != "" {
		tagCreateDTO += nameField
	} else if stateField != "" {
		tagCreateDTO += stateField
	}
	tagCreateDTO += "}"

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/tags", bytes.NewBuffer([]byte(tagCreateDTO)))
	req.Header.Set("Content-Type", "application/json")
	Router.ServeHTTP(w, req)
	return w.Code, w.Body.String(), nil
}

func GetTag(id string) (int, string) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/tags/"+id, nil)
	Router.ServeHTTP(w, req)
	return w.Code, w.Body.String()
}

func GetTags(limit any, offset any) (int, string, error) {
	limitQueryParam := ""
	offsetQueryParam := ""

	switch limitType := limit.(type) {
	case int:
		limitQueryParam = "limit=" + strconv.Itoa(limit.(int))
	case nil:
		limitQueryParam = ""
	default:
		return -1, "", fmt.Errorf("unkown type for 'limit': %v", limitType)
	}

	switch offsetType := offset.(type) {
	case int:
		offsetQueryParam = "offset=" + strconv.Itoa(offset.(int))
	case nil:
		offsetQueryParam = ""
	default:
		return -1, "", fmt.Errorf("unkown type for 'offset': %v", offsetType)
	}

	queryParams := ""
	if limitQueryParam != "" && offsetQueryParam != "" {
		queryParams += "?" + limitQueryParam + "&" + offsetQueryParam
	} else if limitQueryParam != "" {
		queryParams += "?" + limitQueryParam
	} else if offsetQueryParam != "" {
		queryParams += "?" + offsetQueryParam
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/tags"+queryParams, nil)
	Router.ServeHTTP(w, req)
	return w.Code, w.Body.String(), nil
}

func UpdateTag(id any, name any, state any) (int, string, error) {
	idParam := ""
	nameField := ""
	stateField := ""

	switch idType := id.(type) {
	case int:
		idParam = "/" + strconv.Itoa(id.(int))
	case string:
		idParam = "/" + id.(string)
	case nil:
		idParam = ""
	default:
		return -1, "", fmt.Errorf("unkown type for 'id': %v", idType)
	}

	switch nameType := name.(type) {
	case int:
		nameField = "\"Name\": " + strconv.Itoa(name.(int))
	case string:
		nameField = "\"Name\": \"" + name.(string) + "\""
	case nil:
		nameField = ""
	default:
		return -1, "", fmt.Errorf("unkown type for 'name': %v", nameType)
	}

	switch stateType := state.(type) {
	case int:
		stateField = "\"State\": " + strconv.Itoa(state.(int))
	case string:
		stateField = "\"State\": \"" + state.(string) + "\""
	case nil:
		stateField = ""
	default:
		return -1, "", fmt.Errorf("unkown type for 'state': %v", stateType)
	}

	tagUpdateDTO := "{"
	if nameField != "" && stateField != "" {
		tagUpdateDTO += nameField + ", " + stateField
	} else if nameField != "" {
		tagUpdateDTO += nameField
	} else if stateField != "" {
		tagUpdateDTO += stateField
	}
	tagUpdateDTO += "}"

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/tags"+idParam, bytes.NewBuffer([]byte(tagUpdateDTO)))
	req.Header.Set("Content-Type", "application/json")
	Router.ServeHTTP(w, req)
	return w.Code, w.Body.String(), nil
}

func DeleteTag(id any) (int, string, error) {
	idParam := ""

	switch idType := id.(type) {
	case int:
		idParam = "/" + strconv.Itoa(id.(int))
	case string:
		idParam = "/" + id.(string)
	case nil:
		idParam = ""
	default:
		return -1, "", fmt.Errorf("unkown type for 'id': %v", idType)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/tags"+idParam, nil)
	Router.ServeHTTP(w, req)
	return w.Code, w.Body.String(), nil
}

func TestApiTagGet(t *testing.T) {
	t.Run("NotFoundCase", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body := GetTag("1")

		assert.Equal(t, http.StatusNotFound, httpStatusCode)
		assert.Equal(t, "\""+api.PAGE_NOT_FOUND+"\"", body)
	})))
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		expectedId := "1"
		expectedName := "Test Tag 1"
		expectedState := entities.TAG_STATE_NEW
		expectedBody := "{" +
			"\"Id\":" + expectedId + "," +
			"\"Name\":\"" + expectedName + "\"," +
			"\"State\":\"" + expectedState + "\"" +
			"}"

		CreateTag(expectedName, expectedState)
		httpStatusCode, body := GetTag(expectedId)

		assert.Equal(t, http.StatusOK, httpStatusCode)
		assert.Equal(t, expectedBody, body)
	})))
	t.Run("WrongInput: 'Id' is a string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body := GetTag("text")

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_ID_WRONG_FORMAT+"\"", body)
	})))
	t.Run("WrongInput: 'Id' is a float", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body := GetTag("2.15")

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_ID_WRONG_FORMAT+"\"", body)
	})))
	t.Run("WrongInput: 'Id' is a empty string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body := GetTag("")

		assert.Equal(t, http.StatusMovedPermanently, httpStatusCode)
		assert.Equal(t, "<a href=\"/tags\">Moved Permanently</a>.\n\n", body)
	})))
}

func TestApiTagGetAll(t *testing.T) {
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		expectedBody := "{"
		expectedBody += "\"Count\":10,"
		expectedBody += "\"Offset\":0,"
		expectedBody += "\"Limit\":50,"
		expectedBody += "\"Data\":["
		for i := 1; i <= 10; i++ {
			id := strconv.Itoa(i)
			name := "Test Tag " + id
			state := entities.TAG_STATE_NEW

			CreateTag(name, state)
			expectedBody += "{\"Id\":" + id + "," + "\"Name\":\"" + name + "\"," + "\"State\":\"" + state + "\"}"
			if i != 10 {
				expectedBody += ","
			} else {
				expectedBody += "]"
			}
		}
		expectedBody += "}"

		httpStatusCode, body, _ := GetTags(nil, nil)

		assert.Equal(t, http.StatusOK, httpStatusCode)
		assert.Equal(t, expectedBody, body)
	})))
	t.Run("EmptyResult", RunWithRecreateDB((func(t *testing.T) {
		expectedBody := "{"
		expectedBody += "\"Count\":0,"
		expectedBody += "\"Offset\":0,"
		expectedBody += "\"Limit\":50,"
		expectedBody += "\"Data\":[]"
		expectedBody += "}"
		httpStatusCode, body, _ := GetTags(nil, nil)

		assert.Equal(t, http.StatusOK, httpStatusCode)
		assert.Equal(t, expectedBody, body)
	})))
	t.Run("LimitCase", RunWithRecreateDB((func(t *testing.T) {
		expectedBody := "{"
		expectedBody += "\"Count\":5,"
		expectedBody += "\"Offset\":0,"
		expectedBody += "\"Limit\":5,"
		expectedBody += "\"Data\":["
		for i := 1; i <= 5; i++ {
			id := strconv.Itoa(i)
			name := "Test Tag " + id
			state := entities.TAG_STATE_NEW
			expectedBody += "{\"Id\":" + id + "," + "\"Name\":\"" + name + "\"," + "\"State\":\"" + state + "\"}"
			if i != 5 {
				expectedBody += ","
			} else {
				expectedBody += "]"
			}
		}
		expectedBody += "}"

		for i := 1; i <= 10; i++ {
			id := strconv.Itoa(i)
			name := "Test Tag " + id
			state := entities.TAG_STATE_NEW
			CreateTag(name, state)
		}

		httpStatusCode, body, _ := GetTags(5, 0)

		assert.Equal(t, http.StatusOK, httpStatusCode)
		assert.Equal(t, expectedBody, body)
	})))
	t.Run("OffsetCase", RunWithRecreateDB((func(t *testing.T) {
		expectedBody := "{"
		expectedBody += "\"Count\":5,"
		expectedBody += "\"Offset\":5,"
		expectedBody += "\"Limit\":50,"
		expectedBody += "\"Data\":["
		for i := 6; i <= 10; i++ {
			id := strconv.Itoa(i)
			name := "Test Tag " + id
			state := entities.TAG_STATE_NEW
			expectedBody += "{\"Id\":" + id + "," + "\"Name\":\"" + name + "\"," + "\"State\":\"" + state + "\"}"
			if i != 10 {
				expectedBody += ","
			} else {
				expectedBody += "]"
			}
		}
		expectedBody += "}"

		for i := 1; i <= 10; i++ {
			id := strconv.Itoa(i)
			name := "Test Tag " + id
			state := entities.TAG_STATE_NEW
			CreateTag(name, state)
		}

		httpStatusCode, body, _ := GetTags(50, 5)

		assert.Equal(t, http.StatusOK, httpStatusCode)
		assert.Equal(t, expectedBody, body)
	})))
}

func TestApiTagCreate(t *testing.T) {
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := CreateTag("Test Tag 1", entities.TAG_STATE_NEW)

		assert.Equal(t, http.StatusCreated, httpStatusCode)
		assert.Equal(t, "1", body)
	})))
	t.Run("WrongInput: Missed 'Name'", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := CreateTag(nil, entities.TAG_STATE_NEW)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_TAG_NAME_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: Missed 'State'", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := CreateTag("Test Tag 1", nil)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_TAG_STATE_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: Missed 'Name' and 'State'", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := CreateTag(nil, nil)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_TAG_NAME_AND_STATE_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: 'Name' is not a string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := CreateTag(1, entities.TAG_STATE_NEW)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_MESSAGE_PARSING_BODY_JSON+"\"", body)
	})))
	t.Run("WrongInput: 'State' is not a string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := CreateTag("Test Tag 1", 1)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_MESSAGE_PARSING_BODY_JSON+"\"", body)
	})))
	t.Run("WrongInput: 'Name' is empty string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := CreateTag("", entities.TAG_STATE_NEW)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_TAG_NAME_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: 'State' is empty a string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := CreateTag("Test Tag 1", "")

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_TAG_STATE_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: 'State' has a value that not from enum", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := CreateTag("Test Tag 1", "MISSED TEST STATE")

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+ERROR_TAG_CREATE_STATE_WRONG_VALUE+"\"", body)
	})))
	t.Run("DuplicateCase", RunWithRecreateDB((func(t *testing.T) {
		expectedId := "1"

		httpStatusCode, body, _ := CreateTag("Test Tag 1", entities.TAG_STATE_NEW)

		assert.Equal(t, http.StatusCreated, httpStatusCode)
		assert.Equal(t, expectedId, body)

		httpStatusCode, body, _ = CreateTag("Test Tag 1", entities.TAG_STATE_NEW)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.DUPLICATE_FOUND+"\"", body)
	})))
	t.Run("DeletedCase: try to create as deleted", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := CreateTag("Test Tag 1", entities.TAG_STATE_DELETED)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.DELETE_VIA_POST_REQUEST_IS_FODBIDDEN+"\"", body)
	})))
}

func TestApiTagUpdate(t *testing.T) {
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		expectedId := "1"
		expectedName := "Test Tag 2"
		expectedState := entities.TAG_STATE_BLOCKED
		expectedBody := "{" +
			"\"Id\":" + expectedId + "," +
			"\"Name\":\"" + expectedName + "\"," +
			"\"State\":\"" + expectedState + "\"" +
			"}"
		CreateTag("Test Tag 1", entities.TAG_STATE_NEW)

		httpStatusCode, body, _ := UpdateTag(expectedId, "Test Tag 2", entities.TAG_STATE_BLOCKED)

		assert.Equal(t, http.StatusOK, httpStatusCode)
		assert.Equal(t, "\""+api.DONE+"\"", body)

		httpStatusCode, body = GetTag(expectedId)

		assert.Equal(t, http.StatusOK, httpStatusCode)
		assert.Equal(t, expectedBody, body)

	})))
	t.Run("WrongInput: 'Id' is a empty string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := UpdateTag("", "Test Tag 2", entities.TAG_STATE_BLOCKED)

		assert.Equal(t, http.StatusNotFound, httpStatusCode)
		assert.Equal(t, api.PAGE_NOT_FOUND, body)
	})))
	t.Run("WrongInput: 'Id' is a string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := UpdateTag("text", "Test Tag 2", entities.TAG_STATE_BLOCKED)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_ID_WRONG_FORMAT+"\"", body)
	})))
	t.Run("WrongInput: 'Id' is a float", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := UpdateTag("2.15", "Test Tag 2", entities.TAG_STATE_BLOCKED)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_ID_WRONG_FORMAT+"\"", body)
	})))
	t.Run("WrongInput: Missed 'Name'", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := UpdateTag("1", nil, entities.TAG_STATE_BLOCKED)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_TAG_NAME_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: Missed 'State'", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := UpdateTag("1", "Test Tag 2", nil)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_TAG_STATE_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: Missed 'Name' and 'State'", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := UpdateTag("1", nil, nil)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_TAG_NAME_AND_STATE_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: 'Name' is not a string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := UpdateTag("1", 10000, entities.TAG_STATE_BLOCKED)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_MESSAGE_PARSING_BODY_JSON+"\"", body)
	})))
	t.Run("WrongInput: 'State' is not a string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := UpdateTag("1", "Test Tag 2", 10000)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_MESSAGE_PARSING_BODY_JSON+"\"", body)
	})))
	t.Run("WrongInput: 'Name' is empty string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := UpdateTag("1", "", entities.TAG_STATE_BLOCKED)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_TAG_NAME_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: 'State' is empty a string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := UpdateTag("1", "Test Tag 2", "")

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_TAG_STATE_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: 'State' has a value that not from enum", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := UpdateTag("1", "Test Tag 2", "MISSED TEST STATE")

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+ERROR_TAG_UPDATE_STATE_WRONG_VALUE+"\"", body)
	})))
	t.Run("NotFoundCase", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := UpdateTag("1", "Test Tag 2", entities.TAG_STATE_BLOCKED)

		assert.Equal(t, http.StatusNotFound, httpStatusCode)
		assert.Equal(t, "\""+api.PAGE_NOT_FOUND+"\"", body)
	})))
	t.Run("DeletedCase: find deleted", RunWithRecreateDB((func(t *testing.T) {
		expectedId := "1"

		CreateTag("Test Tag 1", entities.TAG_STATE_NEW)
		DeleteTag(expectedId)

		httpStatusCode, body, _ := UpdateTag(expectedId, "Test Tag 2", entities.TAG_STATE_BLOCKED)

		assert.Equal(t, http.StatusNotFound, httpStatusCode)
		assert.Equal(t, "\""+api.PAGE_NOT_FOUND+"\"", body)
	})))
	t.Run("DeletedCase: try to mark as deleted", RunWithRecreateDB((func(t *testing.T) {
		expectedId := "1"

		CreateTag("Test Tag 1", entities.TAG_STATE_NEW)
		httpStatusCode, body, _ := UpdateTag(expectedId, "Test Tag 2", entities.TAG_STATE_DELETED)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.DELETE_VIA_PUT_REQUEST_IS_FODBIDDEN+"\"", body)
	})))
	t.Run("DuplicateCase", RunWithRecreateDB((func(t *testing.T) {
		expectedId := "2"

		CreateTag("Test Tag 1", entities.TAG_STATE_NEW)
		CreateTag("Test Tag 2", entities.TAG_STATE_NEW)

		httpStatusCode, body, _ := UpdateTag(expectedId, "Test Tag 1", entities.TAG_STATE_NEW)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.DUPLICATE_FOUND+"\"", body)
	})))
	t.Run("MultipleUpdateCase", RunWithRecreateDB((func(t *testing.T) {
		expectedId := "1"
		expectedName := "Test Tag 2"
		expectedState := entities.TAG_STATE_BLOCKED
		expectedBody := "{" +
			"\"Id\":" + expectedId + "," +
			"\"Name\":\"" + expectedName + "\"," +
			"\"State\":\"" + expectedState + "\"" +
			"}"

		CreateTag("Test Tag 1", entities.TAG_STATE_NEW)

		for i := 1; i <= 3; i++ {
			httpStatusCode, body, _ := UpdateTag(expectedId, "Test Tag 2", entities.TAG_STATE_BLOCKED)

			assert.Equal(t, http.StatusOK, httpStatusCode)
			assert.Equal(t, "\""+api.DONE+"\"", body)

			httpStatusCode, body = GetTag(expectedId)

			assert.Equal(t, http.StatusOK, httpStatusCode)
			assert.Equal(t, expectedBody, body)
		}
	})))
}

func TestApiTagDelete(t *testing.T) {
	t.Run("BasicCase", RunWithRecreateDB((func(t *testing.T) {
		expectedId := "1"

		CreateTag("Test Tag 1", entities.TAG_STATE_NEW)

		httpStatusCode, body, _ := DeleteTag(expectedId)

		assert.Equal(t, http.StatusOK, httpStatusCode)
		assert.Equal(t, "\""+api.DONE+"\"", body)

		httpStatusCode, body = GetTag(expectedId)

		assert.Equal(t, http.StatusNotFound, httpStatusCode)
		assert.Equal(t, "\""+api.PAGE_NOT_FOUND+"\"", body)
	})))
	t.Run("WrongInput: 'Id' is a string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := DeleteTag("text")

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_ID_WRONG_FORMAT+"\"", body)
	})))
	t.Run("WrongInput: 'Id' is a float", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := DeleteTag("2.15")

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.ERROR_ID_WRONG_FORMAT+"\"", body)
	})))
	t.Run("WrongInput: 'Id' is a empty string", RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := DeleteTag("")

		assert.Equal(t, http.StatusNotFound, httpStatusCode)
		assert.Equal(t, api.PAGE_NOT_FOUND, body)
	})))
	t.Run("MultipleDeleteCase", RunWithRecreateDB((func(t *testing.T) {
		expectedId := "1"

		CreateTag("Test Tag 1", entities.TAG_STATE_NEW)

		httpStatusCode, body, _ := DeleteTag(expectedId)

		assert.Equal(t, http.StatusOK, httpStatusCode)
		assert.Equal(t, "\""+api.DONE+"\"", body)

		httpStatusCode, body, _ = DeleteTag(expectedId)

		assert.Equal(t, http.StatusNotFound, httpStatusCode)
		assert.Equal(t, "\""+api.PAGE_NOT_FOUND+"\"", body)
	})))
}
