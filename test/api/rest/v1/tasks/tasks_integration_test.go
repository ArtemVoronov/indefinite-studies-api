//go:build integration
// +build integration

package tasks_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
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

//TODO: allow missed parameter and any type
func CreateTask(name string, state string) (int, string) {
	taskCreateDTO := tasks.TaskCreateDTO{
		Name:  name,
		State: state,
	}
	jsonValue, _ := json.Marshal(taskCreateDTO)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	return w.Code, w.Body.String()
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

		CreateTask(expectedName, expectedState)
		httpStatusCode, body := GetTask(expectedId)

		assert.Equal(t, http.StatusOK, httpStatusCode)
		assert.Equal(t, "{\"Id\":"+expectedId+",\"Name\":\""+expectedName+"\",\"State\":\""+expectedState+"\"}", body)
	})))
	t.Run("WrongInput", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		t.Errorf("not implemented")
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
		expectedId := "1"

		httpStatusCode, body := CreateTask("Test Task 1", entities.TAG_STATE_NEW)

		assert.Equal(t, http.StatusCreated, httpStatusCode)
		assert.Equal(t, expectedId, body)
	})))
	t.Run("WrongInput: Missed 'Name'", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		taskCreateDTO := tasks.TaskCreateDTO{
			State: "NEW",
		}
		jsonValue, _ := json.Marshal(taskCreateDTO)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "{\"errors\":[{\"Field\":\"Name\",\"Msg\":\"This field is required\"}]}", w.Body.String())
	})))
	t.Run("WrongInput: Missed 'State'", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		taskCreateDTO := tasks.TaskCreateDTO{
			Name: "Test Task 1",
		}
		jsonValue, _ := json.Marshal(taskCreateDTO)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "{\"errors\":[{\"Field\":\"State\",\"Msg\":\"This field is required\"}]}", w.Body.String())
	})))
	t.Run("WrongInput: Missed 'Name' and 'State'", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		taskCreateDTO := tasks.TaskCreateDTO{}
		jsonValue, _ := json.Marshal(taskCreateDTO)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "{\"errors\":["+
			"{\"Field\":\"Name\",\"Msg\":\"This field is required\"},"+
			"{\"Field\":\"State\",\"Msg\":\"This field is required\"}"+
			"]}", w.Body.String())
	})))
	t.Run("WrongInput: 'Name' is not a string", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		taskCreateDTO := "{\"Name\": 1, \"State\": \"NEW\"}"
		jsonValue, _ := json.Marshal(taskCreateDTO)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "\"Error during parsing of HTTP request body. Please check it format correctness: missed brackets, double quotes, commas, matching of names and data types and etc\"", w.Body.String())
	})))
	t.Run("WrongInput: 'State' is not a string", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		taskCreateDTO := "{\"Name\": \"Test Task 1\", \"State\": 1}"
		jsonValue, _ := json.Marshal(taskCreateDTO)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "\"Error during parsing of HTTP request body. Please check it format correctness: missed brackets, double quotes, commas, matching of names and data types and etc\"", w.Body.String())
	})))
	t.Run("WrongInput: 'Name' is empty string", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		taskCreateDTO := tasks.TaskCreateDTO{
			Name:  "",
			State: "NEW",
		}
		jsonValue, _ := json.Marshal(taskCreateDTO)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "{\"errors\":[{\"Field\":\"Name\",\"Msg\":\"This field is required\"}]}", w.Body.String())
	})))
	t.Run("WrongInput: 'State' is empty a string", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		taskCreateDTO := tasks.TaskCreateDTO{
			Name:  "Test Task 1",
			State: "",
		}
		jsonValue, _ := json.Marshal(taskCreateDTO)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "{\"errors\":[{\"Field\":\"State\",\"Msg\":\"This field is required\"}]}", w.Body.String())
	})))
	t.Run("WrongInput: 'State' has a value that not from enum", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		taskCreateDTO := tasks.TaskCreateDTO{
			Name:  "Test Task 1",
			State: "MISSED TEST STATE",
		}
		jsonValue, _ := json.Marshal(taskCreateDTO)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "\"Unable to create task. Wrong 'State' value. Possible values: "+fmt.Sprintf("%v", entities.GetPossibleTaskStates())+"\"", w.Body.String())
	})))
	t.Run("DuplicateCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		expectedId := "1"

		httpStatusCode, body := CreateTask("Test Task 1", entities.TAG_STATE_NEW)

		assert.Equal(t, http.StatusCreated, httpStatusCode)
		assert.Equal(t, expectedId, body)

		httpStatusCode, body = CreateTask("Test Task 1", entities.TAG_STATE_NEW)

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
