//go:build integration
// +build integration

package tasks_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ArtemVoronov/indefinite-studies-api/internal/api/rest/v1/tasks"
	integrationTesting "github.com/ArtemVoronov/indefinite-studies-api/internal/app/testing"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var router *gin.Engine

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/tasks", tasks.GetTasks)
	r.GET("/tasks/:id", tasks.GetTask)
	r.POST("/tasks", tasks.CreateTask)
	r.PUT("/tasks", tasks.UpdateTask)
	r.DELETE("/tasks/:id", tasks.DeleteTask)
	return r
}

func TestMain(m *testing.M) {
	integrationTesting.Setup()
	router = setupRouter()
	code := m.Run()
	integrationTesting.Shutdown()
	os.Exit(code)
}

func TestGetTaskRoute(t *testing.T) {
	t.Run("ExpectedNotFoundError", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/tasks/1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "\"NOT_FOUND\"", w.Body.String())
	})))
	t.Run("BasicCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		t.Errorf("not implemented")
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
		taskCreateDTO := tasks.TaskCreateDTO{
			Name:  "Test Task 1",
			State: "NEW",
		}
		jsonValue, _ := json.Marshal(taskCreateDTO)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, "1", w.Body.String())
	})))
	t.Run("WrongInput", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		t.Errorf("not implemented")
	})))
	t.Run("DuplicateCase", integrationTesting.RunWithRecreateDB((func(t *testing.T) {
		t.Errorf("not implemented")
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
