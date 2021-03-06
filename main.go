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
	"github.com/ArtemVoronov/indefinite-studies-api/internal/app/utils"
	"github.com/ArtemVoronov/indefinite-studies-api/internal/db"
	"github.com/gin-gonic/gin"
)

func main() {

	app.InitEnv()
	auth.Setup()
	host := app.GetHost()

	router := gin.Default()

	router.Use(app.Cors(utils.EnvVar("CORS")))

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

	v1.GET("/ping", ping.Ping)
	v1.POST("/auth/login", auth.Authenicate)
	v1.POST("/auth/refresh-token", auth.RefreshToken)
	authorized := router.Group("/api/v1")
	authorized.Use(app.AuthReqired())
	{
		authorized.GET("/debug/vars", expvar.Handler())
		authorized.GET("/safe-ping", ping.SafePing)

		// TODO: add signup

		authorized.GET("/tasks/", tasks.GetTasks)
		authorized.GET("/tasks/:id", tasks.GetTask)
		authorized.POST("/tasks/", tasks.CreateTask)
		authorized.PUT("/tasks/:id", tasks.UpdateTask)
		authorized.DELETE("/tasks/:id", tasks.DeleteTask)

		authorized.GET("/tags", tags.GetTags)
		authorized.GET("/tags/:id", tags.GetTag)
		authorized.POST("/tags", tags.CreateTag)
		authorized.PUT("/tags/:id", tags.UpdateTag)
		authorized.DELETE("/tags/:id", tags.DeleteTag)

		authorized.GET("/users", users.GetUsers)
		authorized.GET("/users/:id", users.GetUser)
		authorized.POST("/users", users.CreateUser)
		authorized.PUT("/users/:id", users.UpdateUser)
		authorized.DELETE("/users/:id", users.DeleteUser)

		authorized.GET("/notes", notes.GetNotes)
		authorized.GET("/notes/:id", notes.GetNote)
		authorized.POST("/notes", notes.CreateNote)
		authorized.PUT("/notes/:id", notes.UpdateNote)
		authorized.DELETE("/notes/:id", notes.DeleteNote)
	}

	app.StartServer(host, router)
}
