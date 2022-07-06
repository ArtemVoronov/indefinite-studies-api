package tasks

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/ArtemVoronov/indefinite-studies-api/db"
	"github.com/ArtemVoronov/indefinite-studies-api/db/queries"
	"github.com/gin-gonic/gin"
)

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
		c.IndentedJSON(http.StatusInternalServerError, "Unable to get tasks")
		log.Fatalf("Unable to get to tasks : %s", err) // TODO fix os.Exit(1) problem
	}

	c.IndentedJSON(http.StatusOK, tasks)
}

func GetTask(c *gin.Context) {
	idStr := c.Param("id")

	if idStr == "" {
		c.IndentedJSON(http.StatusBadRequest, "Missed ID")
		return
	}

	var id int
	var parseErr error
	if id, parseErr = strconv.Atoi(idStr); parseErr != nil {
		c.IndentedJSON(http.StatusBadRequest, "Wrong ID format. Expected number")
		return
	}

	task, err := queries.GetTask(db.DB, id)
	if err != nil {

		if err == sql.ErrNoRows {
			c.IndentedJSON(http.StatusOK, "NOT_FOUND")
		} else {
			c.IndentedJSON(http.StatusInternalServerError, "Unable to get task")
			log.Fatalf("Unable to get to tasks : %s", err) // TODO fix os.Exit(1) problem
		}
		return
	}
	c.IndentedJSON(http.StatusOK, task)
}

// router.GET("/api/v1/tasks/", api.GetTasks)
// tasks, err := query.GetTasks(DBService, "50", "0")
// if err != nil {
// 	log.Fatalf("Unable to get to tasks : %s", err)
// }

// log.Println(tasks)

// res, err := createTask(db, "Task4", TASK_STATE_NEW)
// if err != nil {
// 	log.Fatalf("Unable to create task : %s", err)
// }

// log.Println(res)

// res, err := deleteTaskByID(db, 4)
// if err != nil {
// 	log.Fatalf("Unable to create task : %s", err)
// }
