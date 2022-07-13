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
	r.PUT("/tasks", tasks.UpdateTask)
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
	ERROR_STATE_WRONG_VALUE string = "\"Unable to create task. Wrong 'State' value. Possible values: " + fmt.Sprintf("%v", entities.GetPossibleTaskStates()) + "\""
	ERROR_ID_WRONG_FORMAT   string = "\"Wrong ID format. Expected number\""
)

func CreateTask(name any, state any) (int, string, error) {
	nameField := ""
	stateField := ""

	switch v := name.(type) {
	case int:
		nameField = "\"Name\": " + strconv.Itoa(name.(int))
	case string:
		nameField = "\"Name\": \"" + name.(string) + "\""
	case nil:
		nameField = ""
	default:
		return -1, "", fmt.Errorf("unkown type for 'name': %v", v)
	}

	switch v := state.(type) {
	case int:
		stateField = "\"State\": " + strconv.Itoa(state.(int))
	case string:
		stateField = "\"State\": \"" + state.(string) + "\""
	case nil:
		stateField = ""
	default:
		return -1, "", fmt.Errorf("unkown type for 'state': %v", v)
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

// func UpdateTask() (int, string) {
//	TODO
// }

// func DeleteTask() (int, string) {
//	TODO
// }

func TestGetTaskRoute(t *testing.T) {
	t.Run("ExpectedNotFoundError", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		httpStatusCode, body := GetTask("1")

		assert.Equal(t, http.StatusOK, httpStatusCode)
		assert.Equal(t, "\""+api.NOT_FOUND+"\"", body)
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
}

func TestGetTasksRoute(t *testing.T) {
	t.Run("BasicCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		t.Errorf("not implemented")
	})))
	t.Run("EmptyResult", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		t.Errorf("not implemented")
	})))
	t.Run("WrongInput", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		t.Errorf("not implemented")
	})))
}

func TestCreateTaskRoute(t *testing.T) {
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
		assert.Equal(t, ERROR_STATE_WRONG_VALUE, body)
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
}

func TestUpdateTaskRoute(t *testing.T) {
	t.Run("BasicCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		t.Errorf("not implemented")
	})))
	t.Run("WrongInput", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		t.Errorf("not implemented")
	})))
	t.Run("DeletedCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		t.Errorf("not implemented")
	})))
	t.Run("DuplicateCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		t.Errorf("not implemented")
	})))
	t.Run("TwiceUpdateCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		t.Errorf("not implemented")
	})))
}

func TestDeleteTaskRoute(t *testing.T) {
	t.Run("BasicCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		t.Errorf("not implemented")
	})))
	t.Run("WrongInput", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		t.Errorf("not implemented")
	})))
	t.Run("DeletedCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		t.Errorf("not implemented")
	})))
	t.Run("DuplicateCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		t.Errorf("not implemented")
	})))
}
