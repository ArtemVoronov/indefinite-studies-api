package main

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/expvar"

	"github.com/ArtemVoronov/indefinite-studies-api/internal/api/rest/v1/auth"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/api/rest/v1/notes"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/api/rest/v1/ping"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/api/rest/v1/tags"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/api/rest/v1/tasks"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/api/rest/v1/users"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/app"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/db"
	"github.com/gin-gonic/gin"
)

func main() {

	app.InitEnv()
	auth.Setup()
	// apiUsers := app.GetApiUsers() // TODO clean
	host := app.GetHost()

	router := gin.Default()

	router.Use(app.Cors())

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	router.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	db.GetInstance()

	// TODO: add permission controller by user role and user state
	// v1 := router.Group("/api/v1", gin.BasicAuth(apiUsers)) // TODO: add auth via jwt, update model accordingly
	v1 := router.Group("/api/v1") // TODO: add auth via jwt, update model accordingly
	{
		v1.GET("/debug/vars", expvar.Handler())
		v1.GET("/ping", ping.Ping)

		v1.POST("/auth/login", auth.Authenicate)
		v1.POST("/auth/verify", auth.Verify)
		// TODO: add logout, signup, refresh token

		v1.GET("/tasks/", tasks.GetTasks)
		v1.GET("/tasks/:id", tasks.GetTask)
		v1.POST("/tasks/", tasks.CreateTask)
		v1.PUT("/tasks/:id", tasks.UpdateTask)
		v1.DELETE("/tasks/:id", tasks.DeleteTask)

		v1.GET("/tags", tags.GetTags)
		v1.GET("/tags/:id", tags.GetTag)
		v1.POST("/tags", tags.CreateTag)
		v1.PUT("/tags/:id", tags.UpdateTag)
		v1.DELETE("/tags/:id", tags.DeleteTag)

		v1.GET("/users", users.GetUsers)
		v1.GET("/users/:id", users.GetUser)
		v1.POST("/users", users.CreateUser)
		v1.PUT("/users/:id", users.UpdateUser)
		v1.DELETE("/users/:id", users.DeleteUser)

		v1.GET("/notes", notes.GetNotes)
		v1.GET("/notes/:id", notes.GetNote)
		v1.POST("/notes", notes.CreateNote)
		v1.PUT("/notes/:id", notes.UpdateNote)
		v1.DELETE("/notes/:id", notes.DeleteNote)
	}

	app.StartServer(host, router)
}
