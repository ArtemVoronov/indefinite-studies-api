//go:build integration
// +build integration

package tasks_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/ArtemVoronov/indefinite-studies-api/internal/api"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/api/rest/v1/tasks"
	integrationTesting "github.com/ArtemVoronov/indefinite-studies-api/internal/app/testing"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/db/entities"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var router *gin.Engine

func TestMain(m *testing.M) {
	integrationTesting.Setup()
	router = SetupRouter()
	code := m.Run()
	integrationTesting.Shutdown()
	os.Exit(code)
}

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/tasks", tasks.GetTasks)
	r.GET("/tasks/:id", tasks.GetTask)
	r.POST("/tasks", tasks.CreateTask)
	r.PUT("/tasks/:id", tasks.UpdateTask)
	r.DELETE("/tasks/:id", tasks.DeleteTask)
	return r
}

var (
	ERROR_MESSAGE_PARSING_BODY_JSON string = "\"Error during parsing of HTTP request body. Please check it format correctness: missed brackets, double quotes, commas, matching of names and data types and etc\""
	ERROR_NAME_IS_REQUIRED          string = "{\"errors\":[" +
		"{\"Field\":\"Name\",\"Msg\":\"This field is required\"}" +
		"]}"
	ERROR_STATE_IS_REQUIRED string = "{\"errors\":[" +
		"{\"Field\":\"State\",\"Msg\":\"This field is required\"}" +
		"]}"
	ERROR_NAME_AND_STATE_IS_REQUIRED string = "{\"errors\":[" +
		"{\"Field\":\"Name\",\"Msg\":\"This field is required\"}," +
		"{\"Field\":\"State\",\"Msg\":\"This field is required\"}" +
		"]}"
	ERROR_CREATE_STATE_WRONG_VALUE string = "\"Unable to create task. Wrong 'State' value. Possible values: " + fmt.Sprintf("%v", entities.GetPossibleTaskStates()) + "\""
	ERROR_UPDATE_STATE_WRONG_VALUE string = "\"Unable to update task. Wrong 'State' value. Possible values: " + fmt.Sprintf("%v", entities.GetPossibleTaskStates()) + "\""
	ERROR_ID_WRONG_FORMAT          string = "\"Wrong ID format. Expected number\""
)

func CreateTask(name any, state any) (int, string, error) {
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

	taskCreateDTO := "{"
	if nameField != "" && stateField != "" {
		taskCreateDTO += nameField + ", " + stateField
	} else if nameField != "" {
		taskCreateDTO += nameField
	} else if stateField != "" {
		taskCreateDTO += stateField
	}
	taskCreateDTO += "}"

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer([]byte(taskCreateDTO)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	return w.Code, w.Body.String(), nil
}

func GetTask(id string) (int, string) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/tasks/"+id, nil)
	router.ServeHTTP(w, req)
	return w.Code, w.Body.String()
}

func GetTasks(limit any, offset any) (int, string, error) {
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
	req, _ := http.NewRequest(http.MethodGet, "/tasks"+queryParams, nil)
	router.ServeHTTP(w, req)
	return w.Code, w.Body.String(), nil
}

func UpdateTask(id any, name any, state any) (int, string, error) {
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

	taskUpdateDTO := "{"
	if nameField != "" && stateField != "" {
		taskUpdateDTO += nameField + ", " + stateField
	} else if nameField != "" {
		taskUpdateDTO += nameField
	} else if stateField != "" {
		taskUpdateDTO += stateField
	}
	taskUpdateDTO += "}"

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/tasks"+idParam, bytes.NewBuffer([]byte(taskUpdateDTO)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	return w.Code, w.Body.String(), nil
}

func DeleteTask(id any) (int, string, error) {
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
	req, _ := http.NewRequest(http.MethodDelete, "/tasks"+idParam, nil)
	router.ServeHTTP(w, req)
	return w.Code, w.Body.String(), nil
}

func TestApiGetTask(t *testing.T) {
	t.Run("ExpectedNotFoundError", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body := GetTask("1")

		assert.Equal(t, http.StatusOK, httpStatusCode)
		assert.Equal(t, "\""+api.PAGE_NOT_FOUND+"\"", body)
	})))
	t.Run("BasicCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedId := "1"
		expectedName := "Test Task 1"
		expectedState := entities.TAG_STATE_NEW
		expectedBody := "{" +
			"\"Id\":" + expectedId + "," +
			"\"Name\":\"" + expectedName + "\"," +
			"\"State\":\"" + expectedState + "\"" +
			"}"

		CreateTask(expectedName, expectedState)
		httpStatusCode, body := GetTask(expectedId)

		assert.Equal(t, http.StatusOK, httpStatusCode)
		assert.Equal(t, expectedBody, body)
	})))
	t.Run("WrongInput: 'Id' is a string", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body := GetTask("text")

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_ID_WRONG_FORMAT, body)
	})))
	t.Run("WrongInput: 'Id' is a float", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body := GetTask("2.15")

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_ID_WRONG_FORMAT, body)
	})))
	t.Run("WrongInput: 'Id' is a empty string", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body := GetTask("")

		assert.Equal(t, http.StatusMovedPermanently, httpStatusCode)
		assert.Equal(t, "<a href=\"/tasks\">Moved Permanently</a>.\n\n", body)
	})))
}

func TestApiGetTasks(t *testing.T) {
	t.Run("BasicCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedBody := "{"
		expectedBody += "\"Count\":10,"
		expectedBody += "\"Offset\":0,"
		expectedBody += "\"Limit\":50,"
		expectedBody += "\"Data\":["
		for i := 1; i <= 10; i++ {
			id := strconv.Itoa(i)
			name := "Test Task " + id
			state := entities.TAG_STATE_NEW

			CreateTask(name, state)
			expectedBody += "{\"Id\":" + id + "," + "\"Name\":\"" + name + "\"," + "\"State\":\"" + state + "\"}"
			if i != 10 {
				expectedBody += ","
			} else {
				expectedBody += "]"
			}
		}
		expectedBody += "}"

		httpStatusCode, body, _ := GetTasks(nil, nil)

		assert.Equal(t, http.StatusOK, httpStatusCode)
		assert.Equal(t, expectedBody, body)
	})))
	t.Run("EmptyResult", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedBody := "{"
		expectedBody += "\"Count\":0,"
		expectedBody += "\"Offset\":0,"
		expectedBody += "\"Limit\":50,"
		expectedBody += "\"Data\":[]"
		expectedBody += "}"
		httpStatusCode, body, _ := GetTasks(nil, nil)

		assert.Equal(t, http.StatusOK, httpStatusCode)
		assert.Equal(t, expectedBody, body)
	})))
	t.Run("LimitCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedBody := "{"
		expectedBody += "\"Count\":5,"
		expectedBody += "\"Offset\":0,"
		expectedBody += "\"Limit\":5,"
		expectedBody += "\"Data\":["
		for i := 1; i <= 5; i++ {
			id := strconv.Itoa(i)
			name := "Test Task " + id
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
			name := "Test Task " + id
			state := entities.TAG_STATE_NEW
			CreateTask(name, state)
		}

		httpStatusCode, body, _ := GetTasks(5, 0)

		assert.Equal(t, http.StatusOK, httpStatusCode)
		assert.Equal(t, expectedBody, body)
	})))
	t.Run("OffsetCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedBody := "{"
		expectedBody += "\"Count\":5,"
		expectedBody += "\"Offset\":5,"
		expectedBody += "\"Limit\":50,"
		expectedBody += "\"Data\":["
		for i := 6; i <= 10; i++ {
			id := strconv.Itoa(i)
			name := "Test Task " + id
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
			name := "Test Task " + id
			state := entities.TAG_STATE_NEW
			CreateTask(name, state)
		}

		httpStatusCode, body, _ := GetTasks(50, 5)

		assert.Equal(t, http.StatusOK, httpStatusCode)
		assert.Equal(t, expectedBody, body)
	})))
}

func TestApiCreateTask(t *testing.T) {
	t.Run("BasicCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := CreateTask("Test Task 1", entities.TAG_STATE_NEW)

		assert.Equal(t, http.StatusCreated, httpStatusCode)
		assert.Equal(t, "1", body)
	})))
	t.Run("WrongInput: Missed 'Name'", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := CreateTask(nil, entities.TAG_STATE_NEW)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_NAME_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: Missed 'State'", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := CreateTask("Test Task 1", nil)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_STATE_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: Missed 'Name' and 'State'", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := CreateTask(nil, nil)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_NAME_AND_STATE_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: 'Name' is not a string", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := CreateTask(1, entities.TAG_STATE_NEW)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_MESSAGE_PARSING_BODY_JSON, body)
	})))
	t.Run("WrongInput: 'State' is not a string", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := CreateTask("Test Task 1", 1)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_MESSAGE_PARSING_BODY_JSON, body)
	})))
	t.Run("WrongInput: 'Name' is empty string", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := CreateTask("", entities.TAG_STATE_NEW)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_NAME_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: 'State' is empty a string", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := CreateTask("Test Task 1", "")

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_STATE_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: 'State' has a value that not from enum", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := CreateTask("Test Task 1", "MISSED TEST STATE")

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_CREATE_STATE_WRONG_VALUE, body)
	})))
	t.Run("DuplicateCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedId := "1"

		httpStatusCode, body, _ := CreateTask("Test Task 1", entities.TAG_STATE_NEW)

		assert.Equal(t, http.StatusCreated, httpStatusCode)
		assert.Equal(t, expectedId, body)

		httpStatusCode, body, _ = CreateTask("Test Task 1", entities.TAG_STATE_NEW)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.DUPLICATE_FOUND+"\"", body)
	})))
	t.Run("DeletedCase: try to create as deleted", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := CreateTask("Test Task 1", entities.TASK_STATE_DELETED)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.DELETE_VIA_POST_REQUEST_IS_FODBIDDEN+"\"", body)
	})))
}

func TestApiUpdateTask(t *testing.T) {
	t.Run("BasicCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedId := "1"
		expectedName := "Test Task 2"
		expectedState := entities.TASK_STATE_DONE
		expectedBody := "{" +
			"\"Id\":" + expectedId + "," +
			"\"Name\":\"" + expectedName + "\"," +
			"\"State\":\"" + expectedState + "\"" +
			"}"
		CreateTask("Test Task 1", entities.TASK_STATE_NEW)

		httpStatusCode, body, _ := UpdateTask(expectedId, "Test Task 2", entities.TASK_STATE_DONE)

		assert.Equal(t, http.StatusOK, httpStatusCode)
		assert.Equal(t, "\""+api.DONE+"\"", body)

		httpStatusCode, body = GetTask(expectedId)

		assert.Equal(t, http.StatusOK, httpStatusCode)
		assert.Equal(t, expectedBody, body)

	})))
	t.Run("WrongInput: 'Id' is a empty string", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := UpdateTask("", "Test Task 2", entities.TASK_STATE_DONE)

		assert.Equal(t, http.StatusNotFound, httpStatusCode)
		assert.Equal(t, api.PAGE_NOT_FOUND, body)
	})))
	t.Run("WrongInput: 'Id' is a string", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := UpdateTask("text", "Test Task 2", entities.TASK_STATE_DONE)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_ID_WRONG_FORMAT, body)
	})))
	t.Run("WrongInput: 'Id' is a float", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := UpdateTask("2.15", "Test Task 2", entities.TASK_STATE_DONE)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_ID_WRONG_FORMAT, body)
	})))
	t.Run("WrongInput: Missed 'Name'", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := UpdateTask("1", nil, entities.TASK_STATE_DONE)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_NAME_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: Missed 'State'", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := UpdateTask("1", "Test Task 2", nil)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_STATE_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: Missed 'Name' and 'State'", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := UpdateTask("1", nil, nil)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_NAME_AND_STATE_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: 'Name' is not a string", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := UpdateTask("1", 10000, entities.TASK_STATE_DONE)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_MESSAGE_PARSING_BODY_JSON, body)
	})))
	t.Run("WrongInput: 'State' is not a string", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := UpdateTask("1", "Test Task 2", 10000)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_MESSAGE_PARSING_BODY_JSON, body)
	})))
	t.Run("WrongInput: 'Name' is empty string", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := UpdateTask("1", "", entities.TASK_STATE_DONE)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_NAME_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: 'State' is empty a string", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := UpdateTask("1", "Test Task 2", "")

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_STATE_IS_REQUIRED, body)
	})))
	t.Run("WrongInput: 'State' has a value that not from enum", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := UpdateTask("1", "Test Task 2", "MISSED TEST STATE")

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_UPDATE_STATE_WRONG_VALUE, body)
	})))
	t.Run("NotFoundCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := UpdateTask("1", "Test Task 2", entities.TASK_STATE_DONE)

		assert.Equal(t, http.StatusNotFound, httpStatusCode)
		assert.Equal(t, "\""+api.PAGE_NOT_FOUND+"\"", body)
	})))
	t.Run("DeletedCase: find deleted", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedId := "1"

		CreateTask("Test Task 1", entities.TASK_STATE_NEW)
		DeleteTask(expectedId)

		httpStatusCode, body, _ := UpdateTask(expectedId, "Test Task 2", entities.TASK_STATE_DONE)

		assert.Equal(t, http.StatusNotFound, httpStatusCode)
		assert.Equal(t, "\""+api.PAGE_NOT_FOUND+"\"", body)
	})))
	t.Run("DeletedCase: try to mark as deleted", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedId := "1"

		CreateTask("Test Task 1", entities.TASK_STATE_NEW)
		httpStatusCode, body, _ := UpdateTask(expectedId, "Test Task 2", entities.TASK_STATE_DELETED)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.DELETE_VIA_PUT_REQUEST_IS_FODBIDDEN+"\"", body)
	})))
	t.Run("DuplicateCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedId := "2"

		CreateTask("Test Task 1", entities.TASK_STATE_NEW)
		CreateTask("Test Task 2", entities.TASK_STATE_NEW)

		httpStatusCode, body, _ := UpdateTask(expectedId, "Test Task 1", entities.TASK_STATE_NEW)

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, "\""+api.DUPLICATE_FOUND+"\"", body)
	})))
	t.Run("MultipleUpdateCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedId := "1"
		expectedName := "Test Task 2"
		expectedState := entities.TASK_STATE_DONE
		expectedBody := "{" +
			"\"Id\":" + expectedId + "," +
			"\"Name\":\"" + expectedName + "\"," +
			"\"State\":\"" + expectedState + "\"" +
			"}"

		CreateTask("Test Task 1", entities.TASK_STATE_NEW)

		for i := 1; i <= 3; i++ {
			httpStatusCode, body, _ := UpdateTask(expectedId, "Test Task 2", entities.TASK_STATE_DONE)

			assert.Equal(t, http.StatusOK, httpStatusCode)
			assert.Equal(t, "\""+api.DONE+"\"", body)

			httpStatusCode, body = GetTask(expectedId)

			assert.Equal(t, http.StatusOK, httpStatusCode)
			assert.Equal(t, expectedBody, body)
		}
	})))
}

func TestApiDeleteTask(t *testing.T) {
	t.Run("BasicCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedId := "1"

		CreateTask("Test Task 1", entities.TASK_STATE_NEW)

		httpStatusCode, body, _ := DeleteTask(expectedId)

		assert.Equal(t, http.StatusOK, httpStatusCode)
		assert.Equal(t, "\""+api.DONE+"\"", body)

		httpStatusCode, body = GetTask(expectedId)

		assert.Equal(t, http.StatusNotFound, httpStatusCode)
		assert.Equal(t, "\""+api.PAGE_NOT_FOUND+"\"", body)
	})))
	t.Run("WrongInput: 'Id' is a string", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := DeleteTask("text")

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_ID_WRONG_FORMAT, body)
	})))
	t.Run("WrongInput: 'Id' is a float", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := DeleteTask("2.15")

		assert.Equal(t, http.StatusBadRequest, httpStatusCode)
		assert.Equal(t, ERROR_ID_WRONG_FORMAT, body)
	})))
	t.Run("WrongInput: 'Id' is a empty string", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body, _ := DeleteTask("")

		assert.Equal(t, http.StatusNotFound, httpStatusCode)
		assert.Equal(t, api.PAGE_NOT_FOUND, body)
	})))
	t.Run("MultipleDeleteCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedId := "1"

		CreateTask("Test Task 1", entities.TASK_STATE_NEW)

		httpStatusCode, body, _ := DeleteTask(expectedId)

		assert.Equal(t, http.StatusOK, httpStatusCode)
		assert.Equal(t, "\""+api.DONE+"\"", body)

		httpStatusCode, body, _ = DeleteTask(expectedId)

		assert.Equal(t, http.StatusNotFound, httpStatusCode)
		assert.Equal(t, "\""+api.PAGE_NOT_FOUND+"\"", body)
	})))
}
