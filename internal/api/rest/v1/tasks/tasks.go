package tasks

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/ArtemVoronov/indefinite-studies-api/internal/api"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/api/validation"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/app/utils"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/db"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/db/entities"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/db/queries"
	"github.com/gin-gonic/gin"
)

type TaskDTO struct {
	Id    int
	Name  string
	State string
}

type TaskListDTO struct {
	Count  int
	Offset int
	Limit  int
	Data   []TaskDTO
}

type TaskEditDTO struct {
	Name  string `json:"name" binding:"required"`
	State string `json:"state" binding:"required"`
}

type TaskCreateDTO struct {
	Name  string `json:"name" binding:"required"`
	State string `json:"state" binding:"required"`
}

func convertTasks(tasks []entities.Task) []TaskDTO {
	if tasks == nil {
		return make([]TaskDTO, 0)
	}
	var result []TaskDTO
	for _, task := range tasks {
		result = append(result, convertTask(task))
	}
	return result
}

func convertTask(task entities.Task) TaskDTO {
	return TaskDTO{Id: task.Id, Name: task.Name, State: task.State}
}

func GetTasks(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	db.RunWithWithTimeout(func(database *sql.DB, ctx context.Context) {
		tasks, err := queries.GetTasks(database, ctx, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "Unable to get tasks")
			log.Printf("Unable to get to tasks : %s", err)
			return
		}
		result := &TaskListDTO{Data: convertTasks(tasks), Count: len(tasks), Offset: offset, Limit: limit}
		c.JSON(http.StatusOK, result)
	})()
}

func GetTask(c *gin.Context) {
	taskIdStr := c.Param("id")

	if taskIdStr == "" {
		c.JSON(http.StatusBadRequest, "Missed ID")
		return
	}

	var taskId int
	var parseErr error
	if taskId, parseErr = strconv.Atoi(taskIdStr); parseErr != nil {
		c.JSON(http.StatusBadRequest, "Wrong ID format. Expected number")
		return
	}

	db.RunWithWithTimeout(func(database *sql.DB, ctx context.Context) {
		task, err := queries.GetTask(database, ctx, taskId)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, api.PAGE_NOT_FOUND)
			} else {
				c.JSON(http.StatusInternalServerError, "Unable to get task")
				log.Printf("Unable to get to task : %s", err)
			}
			return
		}
		c.JSON(http.StatusOK, task)
	})()
}

func CreateTask(c *gin.Context) {
	var task TaskCreateDTO

	if err := c.ShouldBindJSON(&task); err != nil {
		validation.ProcessAndSendValidationErrorMessage(c, err)
		return
	}

	possibleTaskStates := entities.GetPossibleTaskStates()
	if !utils.Contains(possibleTaskStates, task.State) {
		errMsg := fmt.Sprintf("Unable to create task. Wrong 'State' value. Possible values: %v", possibleTaskStates)
		c.JSON(http.StatusBadRequest, errMsg)
		return
	}

	if task.State == entities.TASK_STATE_DELETED {
		c.JSON(http.StatusBadRequest, api.DELETE_VIA_POST_REQUEST_IS_FODBIDDEN)
		return
	}

	db.RunWithWithTimeout(func(database *sql.DB, ctx context.Context) {
		result, err := queries.CreateTask(database, ctx, task.Name, task.State)
		if err != nil || result == -1 {
			if err.Error() == db.ErrorTaskDuplicateKey.Error() {
				c.JSON(http.StatusBadRequest, api.DUPLICATE_FOUND)
				return
			}
			c.JSON(http.StatusInternalServerError, "Unable to create task")
			log.Printf("Unable to create task : %s", err)
			return

		}
		c.JSON(http.StatusCreated, result)
	})()
}

func UpdateTask(c *gin.Context) {
	taskIdStr := c.Param("id")

	if taskIdStr == "" {
		c.JSON(http.StatusBadRequest, "Missed ID")
		return
	}

	var taskId int
	var parseErr error
	if taskId, parseErr = strconv.Atoi(taskIdStr); parseErr != nil {
		c.JSON(http.StatusBadRequest, "Wrong ID format. Expected number")
		return
	}

	var task TaskEditDTO

	if err := c.ShouldBindJSON(&task); err != nil {
		validation.ProcessAndSendValidationErrorMessage(c, err)
		return
	}

	if task.State == entities.TASK_STATE_DELETED {
		c.JSON(http.StatusBadRequest, api.DELETE_VIA_PUT_REQUEST_IS_FODBIDDEN)
		return
	}

	possibleTaskStates := entities.GetPossibleTaskStates()
	if !utils.Contains(possibleTaskStates, task.State) {
		errMsg := fmt.Sprintf("Unable to update task. Wrong 'State' value. Possible values: %v", possibleTaskStates)
		c.JSON(http.StatusBadRequest, errMsg)
		return
	}

	db.RunWithWithTimeout(func(database *sql.DB, ctx context.Context) {
		err := queries.UpdateTask(database, ctx, taskId, task.Name, task.State)

		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, api.PAGE_NOT_FOUND)
				return
			}
			if err.Error() == db.ErrorTaskDuplicateKey.Error() {
				c.JSON(http.StatusBadRequest, api.DUPLICATE_FOUND)
				return
			}
			c.JSON(http.StatusInternalServerError, "Unable to update task")
			log.Printf("Unable to update task : %s", err)
			return
		}
		c.JSON(http.StatusOK, api.DONE)
	})()
}

func DeleteTask(c *gin.Context) {
	idStr := c.Param("id")

	if idStr == "" {
		c.JSON(http.StatusBadRequest, "Missed ID")
		return
	}

	var id int
	var parseErr error
	if id, parseErr = strconv.Atoi(idStr); parseErr != nil {
		c.JSON(http.StatusBadRequest, "Wrong ID format. Expected number")
		return
	}

	db.RunWithWithTimeout(func(database *sql.DB, ctx context.Context) {
		err := queries.DeleteTask(database, ctx, id)

		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, api.PAGE_NOT_FOUND)
				return
			}
			c.JSON(http.StatusInternalServerError, "Unable to delete task")
			log.Printf("Unable to delete task: %s", err)
			return
		}
		c.JSON(http.StatusOK, api.DONE)
	})()
}
