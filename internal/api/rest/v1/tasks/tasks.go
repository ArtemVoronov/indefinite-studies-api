package tasks

import (
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

type TaskEditDTO struct {
	Id    int    `json:"id" binding:"required,numberic"`
	Name  string `json:"name" binding:"required"`
	State string `json:"state" binding:"required"`
}

type TaskCreateDTO struct {
	Name  string `json:"name" binding:"required"`
	State string `json:"state" binding:"required"`
}

func GetTasks(c *gin.Context) {
	limit := c.DefaultQuery("limit", "50")
	offset := c.DefaultQuery("offset", "0")

	if _, err := strconv.Atoi(limit); err != nil {
		limit = "50"
	}

	if _, err := strconv.Atoi(offset); err != nil {
		offset = "0"
	}

	tasks, err := queries.GetTasks(db.DB, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Unable to get tasks")
		log.Printf("Unable to get to tasks : %s", err)
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func GetTask(c *gin.Context) {
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

	task, err := queries.GetTask(db.DB, id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, api.NOT_FOUND)
		} else {
			c.JSON(http.StatusInternalServerError, "Unable to get task")
			log.Printf("Unable to get to task : %s", err)
		}
		return
	}
	c.JSON(http.StatusOK, task)
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

	result, err := queries.CreateTask(db.DB, task.Name, task.State)
	if err != nil || result == -1 {
		if err.Error() == db.ErrorDuplicateKey.Error() {
			c.JSON(http.StatusBadRequest, api.DUPLICATE_FOUND)
			return
		}
		c.JSON(http.StatusInternalServerError, "Unable to create task")
		log.Printf("Unable to create task : %s", err)
		return

	}
	c.JSON(http.StatusCreated, result)
}

func UpdateTask(c *gin.Context) {
	var task TaskEditDTO

	if err := c.ShouldBindJSON(&task); err != nil {
		validation.ProcessAndSendValidationErrorMessage(c, err)
		return
	}

	err := queries.UpdateTask(db.DB, task.Id, task.Name, task.State)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, api.NOT_FOUND)
			return
		}

		c.JSON(http.StatusInternalServerError, "Unable to create task")
		log.Printf("Unable to update tasks : %s", err)
		return
	}
	c.JSON(http.StatusOK, "Done")
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
	err := queries.DeleteTask(db.DB, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Unable to delete task")
		log.Printf("Unable to delete task: %s", err)
		return
	}
	c.JSON(http.StatusOK, "Done")
}
